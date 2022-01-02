package shows

import (
	"errors"
	"fengqi/kodi-metadata-tmdb-cli/config"
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

			files, err := dir.scanShowsFile()
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
			err := c.watcher.Add(item + "/" + showDir.OriginTitle)
			utils.Logger.DebugF("runWatcher add shows dir: %s to watcher", item+"/"+showDir.OriginTitle)
			if err != nil {
				utils.Logger.FatalF("add shows dir: %s to watcher err :%v", item+"/"+showDir.OriginTitle, err)
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
			if fileInfo == nil || !fileInfo.IsDir() || event.Op&fsnotify.Create != fsnotify.Create {
				continue
			}

			utils.Logger.InfoF("created file: %s", event.Name)

			showsDir := parseShowsDir(filepath.Dir(event.Name), fileInfo)
			if showsDir != nil {
				c.dirChan <- showsDir
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
					utils.Logger.FatalF("add shows dir: %s to err: %v err: %v", item, err)
				}

				showDirs, err := c.scanDir(item)
				if err != nil {
					utils.Logger.FatalF("scan shows dir: %s err :%v", item, err)
				}

				for _, showDir := range showDirs {
					err := c.watcher.Add(item + "/" + showDir.OriginTitle)
					utils.Logger.DebugF("runCronScan add shows dir: %s to watcher", item+"/"+showDir.OriginTitle)
					if err != nil {
						utils.Logger.FatalF("add shows dir: %s to err: %v err: %v", showDir, err)
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
