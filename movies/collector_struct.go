package movies

import (
	"fengqi/kodi-metadata-tmdb-cli/config"
	"github.com/fsnotify/fsnotify"
)

type Collector struct {
	config  *config.Config
	watcher *fsnotify.Watcher
	channel chan *Movie
}

// Movie 电影目录详情，从名字分析
// Fortress.2021.BluRay.1080p.AVC.DTS-HD.MA5.1-MTeam
type Movie struct {
	Dir             string `json:"dir"`
	OriginTitle     string `json:"origin_title"` // 原始文件名
	Title           string `json:"title"`        // 名称 Hawkeye
	MovieId         int    `json:"tv_id"`        // 电影id
	Year            int    `json:"year"`         // 年份：2020、2021
	Format          string `json:"format"`       // 格式：720p、1080p
	Source          string `json:"source"`       // 来源
	Studio          string `json:"studio"`       // 媒体
	IsFile          bool   `json:"is_file"`      // 是否是单文件，而不是目录
	Suffix          string `json:"suffix"`       // 单文件时，文件的后缀
	IsBluray        bool   `json:"is_bluray"`    // 蓝光目录
	IsDvd           bool   `json:"is_dvd"`       // DVD目录
	IdCacheFile     string `json:"id_cache_file"`
	DetailCacheFile string `json:"detail_cache_file"`
}
