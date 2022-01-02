package shows

import (
	"fengqi/kodi-metadata-tmdb-cli/config"
	"github.com/fsnotify/fsnotify"
)

type Collector struct {
	config   *config.Config
	watcher  *fsnotify.Watcher
	dirChan  chan *Dir
	fileChan chan *File
}

// Dir 电视剧目录详情，从名字分析
// World.Heritage.In.China.E01-E38.2008.CCTVHD.x264.AC3.720p-CMCT
type Dir struct {
	Dir          string `json:"dir"`
	OriginTitle  string `json:"origin_title"`  // 原始文件名
	Title        string `json:"title"`         // 名称 Hawkeye
	TvId         int    `json:"tv_id"`         // TMDV tv id
	Season       int    `json:"season"`        // 第几季 ，电影类 -1
	Year         int    `json:"year"`          // 年份：2020、2021
	Format       string `json:"format"`        // 格式：720p、1080p
	Source       string `json:"source"`        // 来源
	Studio       string `json:"studio"`        // 媒体
	IsCollection bool   `json:"is_collection"` // 是否是合集
}

// File 电视剧目录内文件详情，从名字分析
// Dexter.New.Blood.S01E04.H.is.for.Hero.1080p.AMZN.WEB-DL.DDP5.1.H.264-NTb.mkv
type File struct {
	Dir           string `json:"dir"`
	OriginTitle   string `json:"origin_title"` // 原始文件名
	Season        int    `json:"season"`       // 第几季 ，电影类 -1
	Episode       int    `json:"episode"`      // 第几集，电影类 -1
	SeasonEpisode string `json:"season_episode"`
	Suffix        string `json:"suffix"`
	TvId          int    `json:"tv_id"`
	//TvDetail      *tmdb.TvDetail `json:"tv_detail"`
}
