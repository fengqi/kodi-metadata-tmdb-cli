package movies

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
		utils.Logger.FatalF("new movies watcher err: %v", err)
	}

	collector := &Collector{
		config:  config,
		watcher: watcher,
		channel: make(chan *Movie, 100),
		nfoMode: config.MoviesNfoMode,
	}

	go collector.runWatcher()
	go collector.runMoviesProcess()
	collector.runCronScan()
}

// 开启文件夹监听
func (c *Collector) runWatcher() {
	utils.Logger.Debug("run movies watcher")

	// 监听顶级目录
	for _, item := range c.config.MoviesDir {
		utils.Logger.DebugF("add movies dir: %s to watcher", item)

		err := c.watcher.Add(item)
		if err != nil {
			utils.Logger.FatalF("add movies dir: %s to watcher err :%v", item, err)
		}
	}

	for {
		select {
		// 接受事件，增删改查都会收到，需要过滤，部分情况下可能收不到create而是chmod
		case event, ok := <-c.watcher.Events:
			if !ok {
				return
			}

			if event.Op&fsnotify.Create != fsnotify.Create {
				continue
			}

			fileInfo, _ := os.Stat(event.Name)
			if fileInfo == nil || (!fileInfo.IsDir() && utils.IsVideo(event.Name) == "") {
				continue
			}

			utils.Logger.InfoF("created file: %s", event.Name)

			moviesDir := parseMoviesDir(filepath.Dir(event.Name), fileInfo)
			if moviesDir != nil {
				c.channel <- moviesDir
			}

		case err, ok := <-c.watcher.Errors:
			if !ok {
				return
			}

			utils.Logger.ErrorF("movies watcher error: %v", err)
		}
	}
}

// 电影信息处理：来源包括cron和inotify监听的
func (c *Collector) runMoviesProcess() {
	utils.Logger.Debug("run movies process")

	for {
		select {
		case dir := <-c.channel:
			utils.Logger.DebugF("receive movies task: %v", dir)

			dir.checkCacheDir()
			detail, err := dir.getMovieDetail()
			if err != nil || detail == nil {
				continue
			}

			_ = dir.saveToNfo(detail, c.nfoMode)
			_ = dir.downloadImage(detail)

			// 通知kodi刷新媒体库，电视可能没开机，所以先ping一下
			if kodi.Ping() {
				utils.Logger.DebugF("ping kodi success, starting refresh movies of library")
				videoLibrary := kodi.NewVideoLibrary()
				kodiMoviesReq := &kodi.GetMoviesRequest{
					Filter: &kodi.Filter{
						Field:    "originaltitle",
						Operator: "is",
						Value:    detail.OriginalTitle,
					},
					Limit: &kodi.Limits{
						Start: 0,
						End:   1,
					},
					Properties: []string{"title", "originaltitle", "year"},
				}

				kodiMoviesResp := videoLibrary.GetMovies(kodiMoviesReq)
				if kodiMoviesResp.Limits.Total == 0 {
					utils.Logger.DebugF("maybe new movies, scan video library")
					videoLibrary.Scan(nil)
				} else {
					utils.Logger.DebugF("maybe existing movies, refresh video library")
					videoLibrary.RefreshMovie(&kodi.RefreshMovieRequest{MovieId: kodiMoviesResp.Movies[0].MovieId, IgnoreNfo: false})
				}
			}
		}
	}
}

// 运行定时扫描
func (c *Collector) runCronScan() {
	utils.Logger.DebugF("run movies scan cron_seconds: %d", c.config.CronSeconds)

	ticker := time.NewTicker(time.Second * time.Duration(c.config.CronSeconds))
	for {
		select {
		case <-ticker.C:
			for _, item := range c.config.MoviesDir {
				utils.Logger.DebugF("movies scan ticker trigger")

				movieDirs, err := c.scanDir(item)
				if err != nil {
					utils.Logger.FatalF("scan movies dir: %s err :%v", item, err)
					continue
				}

				for _, movieDir := range movieDirs {
					c.channel <- movieDir
				}
			}
		}
	}
}

// 扫描普通目录，返回其中的电影
func (c *Collector) scanDir(dir string) ([]*Movie, error) {
	if f, err := os.Stat(dir); err != nil || !f.IsDir() {
		return nil, errors.New(fmt.Sprintf("scan err: %v or %s is not dir", err, dir))
	}

	movieDirs := make([]*Movie, 0)
	fileInfo, err := ioutil.ReadDir(dir)
	if err != nil {
		utils.Logger.ErrorF("scan dir: %s err: %v", dir, err)
		return nil, err
	}

	for _, file := range fileInfo {
		// 合集，以 Iron.Man.2008-2013.Blu-ray.x264.MiniBD1080P-CMCT 为例，暂定使用 2008-2013 做为判断特征
		if yearRange := utils.IsYearRangeLike(file.Name()); yearRange != "" {
			movieDir, err := c.scanDir(dir + "/" + file.Name())
			if err != nil {
				utils.Logger.ErrorF("scan collection dir: %s err: %v", dir+"/"+file.Name(), err)
				continue
			}
			movieDirs = append(movieDirs, movieDir...)
			continue
		}

		movieDir := parseMoviesDir(dir, file)
		if movieDir == nil {
			continue
		}

		movieDirs = append(movieDirs, movieDir)
	}

	return movieDirs, nil
}
