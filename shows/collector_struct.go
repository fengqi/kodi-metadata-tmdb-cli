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
	Title        string `json:"title"`         // 从视频提取的完整文件名 鹰眼 Hawkeye
	ChsTitle     string `json:"chs_title"`     // 分离出来的中午名称 鹰眼
	EngTitle     string `json:"eng_title"`     // 分离出来的英文名称 Hawkeye
	TvId         int    `json:"tv_id"`         // TMDb tv id
	GroupId      string `json:"group_id"`      // TMDB Episode Group
	Season       int    `json:"season"`        // 第几季 ，电影类 -1
	SeasonRange  string `json:"season_range"`  // 合集：S01-S05
	Year         int    `json:"year"`          // 年份：2020、2021
	YearRange    string `json:"year_range"`    // 年份：2010-2015
	Format       string `json:"format"`        // 格式：720p、1080p
	Source       string `json:"source"`        // 来源
	Studio       string `json:"studio"`        // 媒体
	IsCollection bool   `json:"is_collection"` // 是否是合集目录
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
