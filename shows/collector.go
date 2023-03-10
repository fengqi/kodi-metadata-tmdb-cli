package shows

import (
	"encoding/json"
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

var collector *Collector

func RunCollector(config *config.Config) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		utils.Logger.FatalF("new shows watcher err: %v", err)
	}

	collector = &Collector{
		config:  config,
		watcher: watcher,
		dirChan: make(chan *Dir, 100),
	}

	go collector.runWatcher()
	go collector.showsDirProcess()
	collector.runCronScan()
}

// 目录处理队列消费
func (c *Collector) showsDirProcess() {
	utils.Logger.Debug("run shows dir process")

	for {
		select {
		case dir := <-c.dirChan: // todo dir处理挪到独立的方法
			utils.Logger.DebugF("shows dir process receive task: %v", dir.OriginTitle)

			dir.checkCacheDir()
			detail, err := dir.getTvDetail()
			if err != nil || detail == nil {
				continue
			}

			if !detail.FromCache || !dir.NfoExist() {
				_ = dir.saveToNfo(detail)
				kodi.Rpc.AddRefreshTask(kodi.TaskRefreshTVShow, detail.OriginalName)
			}

			dir.downloadImage(detail)

			files := make(map[int]map[string]*File, 0)
			if dir.IsCollection { // 合集
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

					if len(subFiles) > 0 {
						files[item.Season] = subFiles
					}
				}
			} else { // 普通剧集
				subFiles, err := dir.scanShowsFile()
				if err != nil {
					utils.Logger.ErrorF("scan shows dir: %s err: %v", dir.OriginTitle, err)
					continue
				}

				if len(subFiles) > 0 {
					files[dir.Season] = subFiles
				}
			}

			if len(files) == 0 {
				utils.Logger.WarningF("scan shows file empty: %s", dir.OriginTitle)
				continue
			}

			// 剧集组的分集信息写入缓存, 供后面处理分集信息使用
			if dir.GroupId != "" && detail.TvEpisodeGroupDetail != nil {
				for _, group := range detail.TvEpisodeGroupDetail.Groups {
					group.SortEpisode()
					for k, episode := range group.Episodes {
						se := fmt.Sprintf("s%02de%02d", group.Order, k+1)
						file, ok := files[group.Order][se]
						if !ok {
							continue
						}

						cacheFile := fmt.Sprintf("%s/tmdb/%s.json", file.Dir, se)
						f, err := os.OpenFile(cacheFile, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
						if err != nil {
							utils.Logger.ErrorF("save tv to cache, open_file err: %v", err)
							return
						}

						episode.EpisodeNumber = k + 1
						episode.SeasonNumber = group.Order
						bytes, err := json.MarshalIndent(episode, "", "    ")
						if err != nil {
							utils.Logger.ErrorF("save tv to cache, marshal struct errr: %v", err)
							return
						}

						_, err = f.Write(bytes)
						_ = f.Close()
					}
				}
			}

			for _, file := range files {
				for _, subFile := range file {
					c.showsFileProcess(detail.OriginalName, subFile)
				}
			}
		}
	}
}

// 单个剧集处理
func (c *Collector) showsFileProcess(originalName string, showsFile *File) bool {
	episodeDetail, err := showsFile.getTvEpisodeDetail()
	if err != nil || episodeDetail == nil || episodeDetail.FromCache && showsFile.NfoExist() {
		return false
	}

	_ = showsFile.saveToNfo(episodeDetail)
	showsFile.downloadImage(episodeDetail)

	taskVal := fmt.Sprintf("%s|-|%d|-|%d", originalName, episodeDetail.SeasonNumber, episodeDetail.EpisodeNumber)
	kodi.Rpc.AddRefreshTask(kodi.TaskRefreshEpisode, taskVal)

	return true
}

// 目录监听，新增的增加到队列，删除的移除监听
func (c *Collector) runWatcher() {
	if !c.config.Collector.Watcher {
		return
	}

	utils.Logger.Debug("run shows watcher")

	for {
		select {
		case event, ok := <-c.watcher.Events:
			if !ok {
				continue
			}

			fileInfo, err := os.Stat(event.Name)
			if fileInfo == nil || err != nil {
				utils.Logger.WarningF("get shows stat err: %v", err)
				continue
			}

			// 删除文件夹
			if event.Has(fsnotify.Remove) && fileInfo.IsDir() {
				utils.Logger.InfoF("removed dir: %s", event.Name)

				err := c.watcher.Remove(event.Name)
				if err != nil {
					utils.Logger.WarningF("remove shows watcher: %s error: %v", event.Name, err)
				}
				continue
			}

			// 新增文件夹
			if event.Has(fsnotify.Create) && fileInfo.IsDir() {
				utils.Logger.InfoF("created dir: %s", event.Name)

				showsDir := parseShowsDir(filepath.Dir(event.Name), fileInfo)
				if showsDir != nil {
					c.dirChan <- showsDir
				}

				err = c.watcher.Add(event.Name)
				if err != nil {
					utils.Logger.FatalF("add shows dir: %s to watcher err: %v", event.Name, err)
				}
				continue
			}

			// 新增剧集文件
			if event.Has(fsnotify.Create) && utils.IsVideo(event.Name) != "" {
				utils.Logger.InfoF("created file: %s", event.Name)

				filePath := filepath.Dir(event.Name)
				basePath := filepath.Dir(filePath)
				dirInfo, _ := os.Stat(filePath)
				dir := parseShowsDir(basePath, dirInfo)
				if dir != nil {
					c.dirChan <- dir
				}
			}

		case err, ok := <-c.watcher.Errors:
			utils.Logger.ErrorF("shows watcher error: %v", err)

			if !ok {
				return
			}
		}
	}
}

// 目录扫描，定时任务，扫描到的目录和文件增加到队列
func (c *Collector) runCronScan() {
	utils.Logger.DebugF("run shows scan cron_seconds: %d", c.config.Collector.CronSeconds)

	task := func() {
		for _, item := range c.config.Collector.ShowsDir {
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

		if c.config.Kodi.CleanLibrary {
			kodi.Rpc.AddCleanTask("")
		}
	}

	task() // TODO 启动后立即运行可控
	ticker := time.NewTicker(time.Second * time.Duration(c.config.Collector.CronSeconds))
	for {
		select {
		case <-ticker.C:
			task()
		}
	}
}

// 扫描普通目录，返回其中的电视剧
func (c *Collector) scanDir(dir string) ([]*Dir, error) {
	movieDirs := make([]*Dir, 0)

	if f, err := os.Stat(dir); err != nil || !f.IsDir() {
		return movieDirs, nil
	}

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
func (d *Dir) scanShowsFile() (map[string]*File, error) {
	fileInfo, err := ioutil.ReadDir(d.Dir + "/" + d.OriginTitle)
	if err != nil {
		return nil, err
	}

	movieFiles := make(map[string]*File, 0)
	for _, file := range fileInfo {
		movieFile := parseShowsFile(d, file)
		if movieFile == nil {
			continue
		}

		movieFiles[movieFile.SeasonEpisode] = movieFile
	}

	return movieFiles, nil
}
