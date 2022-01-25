package shows

import (
	"errors"
	"fengqi/kodi-metadata-tmdb-cli/config"
	"fengqi/kodi-metadata-tmdb-cli/kodi"
	"fengqi/kodi-metadata-tmdb-cli/utils"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

func RunCollector(config *config.Config) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		utils.Logger.FatalF("new shows watcher err: %v", err)
	}

	collector := &Collector{
		config:   config,
		watcher:  watcher,
		dirChan:  make(chan *Dir, 100),
		fileChan: make(chan *File, 100),
	}

	go collector.runWatcher()
	go collector.showsDirProcess()
	go collector.showsFileProcess()
	collector.runCronScan()
}

// 目录处理队列消费
func (c *Collector) showsDirProcess() {
	utils.Logger.Debug("run shows dir process")

	for {
		select {
		case dir := <-c.dirChan:
			dir.checkCacheDir()
			detail, err := dir.getTvDetail()
			if err != nil || detail == nil {
				continue
			}

			_ = dir.saveToNfo(detail)
			dir.downloadImage(detail)

			// 通知kodi刷新媒体库，电视可能没开机，所以先ping一下
			// 电视剧还要刷新分集信息，所以这里放到后台
			// TODO 电影和电视剧的Kodi通知可以使用队列，等电视机开机时立即通知，然后可以关掉电视的开机刷新媒体库
			go func() {
				if kodi.Ping() {
					utils.Logger.DebugF("ping kodi success, starting refresh shows of library")

					videoLibrary := kodi.NewVideoLibrary()
					kodiTvShowsReq := &kodi.GetTVShowsRequest{
						Filter: &kodi.Filter{
							Field:    "originaltitle",
							Operator: "is",
							Value:    detail.OriginalName,
						},
						Limit: &kodi.Limits{
							Start: 0,
							End:   1,
						},
						Properties: []string{"title", "originaltitle", "year"},
					}
					kodiShowsResp := videoLibrary.GetTVShows(kodiTvShowsReq)
					if kodiTvShowsReq == nil || kodiShowsResp.Limits.Total == 0 {
						utils.Logger.DebugF("maybe new shows, scan video library")
						videoLibrary.Scan(nil)
					} else {
						utils.Logger.DebugF("maybe existing shows, refresh video library")
						kodiRefreshReq := &kodi.RefreshTVShowRequest{
							TvShowId:        kodiShowsResp.TvShows[0].TvShowId,
							IgnoreNfo:       false,
							RefreshEpisodes: true,
						}
						videoLibrary.RefreshTVShow(kodiRefreshReq)
					}
				}
			}()

			files := make([]*File, 0)
			if dir.IsCollection {
				subDir, err := c.scanDir(dir.GetFullDir())
				if err != nil {
					utils.Logger.ErrorF("scan collection dir: %s err: %v", dir.OriginTitle, err)
					continue
				}

				for _, item := range subDir {
					item.checkCacheDir()
					subFiles, err := item.scanShowsFile()
					if err != nil {
						utils.Logger.ErrorF("scan collection sub dir: %s err: %v", item.OriginTitle, err)
						continue
					}
					files = append(files, subFiles...)
				}
			} else {
				files, err = dir.scanShowsFile()
			}

			if err != nil || len(files) == 0 {
				continue
			}

			for _, file := range files {
				file.TvId = detail.Id
				c.fileChan <- file
			}
		}
	}
}

// 文件处理队列消费
func (c *Collector) showsFileProcess() {
	utils.Logger.Debug("run shows file process")

	for {
		select {
		case file := <-c.fileChan:
			detail, err := file.getTvEpisodeDetail()
			if err != nil || detail == nil {
				continue
			}

			_ = file.saveToNfo(detail)
			file.downloadImage(detail)
		}
	}
}

// 目录监听，新增的增加到队列，删除的移除监听
func (c *Collector) runWatcher() {
	utils.Logger.Debug("run shows watcher")

	for _, item := range c.config.ShowsDir {
		err := c.watcher.Add(item)
		utils.Logger.DebugF("runWatcher add shows dir: %s to watcher", item)
		if err != nil {
			utils.Logger.FatalF("add shows dir: %s to watcher err :%v", item, err)
		}

		showDirs, err := c.scanDir(item)
		if err != nil {
			utils.Logger.FatalF("scan shows dir for watcher err :%v", err)
		}

		for _, showDir := range showDirs {
			err := c.watcher.Add(showDir.Dir + "/" + showDir.OriginTitle)
			utils.Logger.DebugF("runWatcher add shows dir: %s to watcher", showDir.Dir+"/"+showDir.OriginTitle)
			if err != nil {
				utils.Logger.FatalF("add shows dir: %s to watcher err :%v", showDir.Dir+"/"+showDir.OriginTitle, err)
			}
		}
	}

	for {
		select {
		case event, ok := <-c.watcher.Events:
			if !ok {
				return
			}

			fileInfo, _ := os.Stat(event.Name)
			if fileInfo == nil || (!fileInfo.IsDir() && utils.IsVideo(event.Name) == "") || event.Op&fsnotify.Create != fsnotify.Create {
				continue
			}

			utils.Logger.InfoF("created file: %s", event.Name)

			if fileInfo.IsDir() {
				showsDir := parseShowsDir(filepath.Dir(event.Name), fileInfo)
				if showsDir != nil {
					c.dirChan <- showsDir
				}
			} else {
				// 刷新剧集
				filePath := filepath.Dir(event.Name)
				basePath := filepath.Dir(filePath)
				dirInfo, _ := os.Stat(filePath)
				dir := parseShowsDir(basePath, dirInfo)
				if dir != nil {
					c.dirChan <- dir
				}

				// 刷新单集
				tvDetail, _ := dir.getTvDetail()
				showsFile := parseShowsFile(dir, fileInfo)
				showsFile.TvId = tvDetail.Id
				if showsFile != nil {
					c.fileChan <- showsFile
				}
			}

		case err, ok := <-c.watcher.Errors:
			if !ok {
				return
			}

			utils.Logger.ErrorF("shows watcher error: %v", err)
		}
	}
}

// 目录扫描，定时任务，扫描到的目录和文件增加到队列
func (c *Collector) runCronScan() {
	utils.Logger.DebugF("run shows scan cron_seconds: %d", c.config.CronSeconds)

	ticker := time.NewTicker(time.Second * time.Duration(c.config.CronSeconds))
	for {
		select {
		case <-ticker.C:
			for _, item := range c.config.ShowsDir {
				// 扫描到的每个目录都添加到watcher，因为还不能只监听根目录
				err := c.watcher.Add(item)
				utils.Logger.DebugF("runCronScan add shows dir: %s to watcher", item)
				if err != nil {
					utils.Logger.FatalF("add shows dir: %s to watcher err: %v", item, err)
				}

				showDirs, err := c.scanDir(item)
				if err != nil {
					utils.Logger.FatalF("scan shows dir: %s err :%v", item, err)
				}

				for _, showDir := range showDirs {
					err := c.watcher.Add(showDir.Dir + "/" + showDir.OriginTitle)
					utils.Logger.DebugF("runCronScan add shows dir: %s to watcher", showDir.Dir+"/"+showDir.OriginTitle)
					if err != nil {
						utils.Logger.FatalF("add shows dir: %s to watcher err: %v", showDir.Dir+"/"+showDir.OriginTitle, err)
					}

					c.dirChan <- showDir
				}
			}
		}
	}
}

// 扫描普通目录，返回其中的电视剧
func (c *Collector) scanDir(dir string) ([]*Dir, error) {
	if f, err := os.Stat(dir); err != nil || !f.IsDir() {
		return nil, errors.New(fmt.Sprintf("scan err: %v or %s is not dir", err, dir))
	}

	movieDirs := make([]*Dir, 0)
	fileInfo, err := ioutil.ReadDir(dir)
	if err != nil {
		utils.Logger.ErrorF("scan dir: %s err: %v", dir, err)
		return nil, err
	}

	for _, file := range fileInfo {
		if !file.IsDir() {
			continue
		}

		movieDir := parseShowsDir(dir, file)
		if movieDir == nil {
			continue
		}

		movieDirs = append(movieDirs, movieDir)
	}

	return movieDirs, nil
}

// ScanMovieFile 扫描可以确定的单个电影、电视机目录，返回其中的视频文件信息
func (d *Dir) scanShowsFile() ([]*File, error) {
	fileInfo, err := ioutil.ReadDir(d.Dir + "/" + d.OriginTitle)
	if err != nil {
		return nil, err
	}

	movieFiles := make([]*File, 0)
	for _, file := range fileInfo {
		movieFile := parseShowsFile(d, file)
		if movieFile == nil {
			continue
		}

		movieFiles = append(movieFiles, movieFile)
	}

	return movieFiles, nil
}
