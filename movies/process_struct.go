package movies

import (
	"fengqi/kodi-metadata-tmdb-cli/media_file"
	"fengqi/kodi-metadata-tmdb-cli/utils"
	"os"
	"path/filepath"
)

// Movie 电影详情，从名字分析
// Fortress.2021.BluRay.1080p.AVC.DTS-HD.MA5.1-MTeam
type Movie struct {
	MediaFile     *media_file.MediaFile `json:"media_file"`     // 媒体文件
	PosterFile    string                `json:"poster_file"`    // 海报图
	FanArtFile    string                `json:"fanart_file"`    // 艺术图
	ClearLogoFile string                `json:"clearlogo_file"` // Logo图
	NfoFile       string                `json:"nfo_file"`       // NFO文件
	MovieId       int                   `json:"tv_id"`          // 电影id
	Title         string                `json:"title"`          // 从视频提取的完整文件名 鹰眼 Hawkeye
	AliasTitle    string                `json:"alias_title"`    // 别名，通常没有用
	ChsTitle      string                `json:"chs_title"`      // 分离出来的中午名称 鹰眼
	EngTitle      string                `json:"eng_title"`      // 分离出来的英文名称 Hawkeye
	Year          int                   `json:"year"`           // 年份：2020、2021
	//Dir             string                `json:"dir"`
	//OriginTitle     string                `json:"origin_title"`   // 原始目录名
	//VideoFileName   string                `json:"file_name"`      // 视频文件名，仅限：IsSingleFile=true
	//IsFile          bool                  `json:"is_file"`        // 是否是单文件，而不是目录
	//Suffix          string                `json:"suffix"`         // 单文件时，文件的后缀
	//IsBluRay        bool                  `json:"is_bluray"`      // 蓝光目录
	//IsDvd           bool                  `json:"is_dvd"`         // DVD目录
	//IsSingleFile    bool                  `json:"is_single_file"` // 普通的单文件视频
	//IdCacheFile     string                `json:"id_cache_file"`
	//DetailCacheFile string                `json:"detail_cache_file"`
}

// tmdb 缓存目录
func (m *Movie) checkCacheDir() {
	dir := m.GetCacheDir()
	if _, err := os.Stat(dir); err != nil && os.IsNotExist(err) {
		if err := os.Mkdir(dir, 0755); err != nil {
			utils.Logger.ErrorF("create cache: %s dir err: %v", dir, err)
		}
	}
}

func (m *Movie) GetCacheDir() string {
	base := filepath.Dir(m.MediaFile.Path)
	return base + "/tmdb"
}
