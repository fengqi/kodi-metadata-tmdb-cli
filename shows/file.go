package shows

import (
	"fengqi/kodi-metadata-tmdb-cli/tmdb"
	"os"
	"strings"
)

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
	Part          int    `json:"part"` // 分卷模式下，第几部分
	//TvDetail      *tmdb.TvDetail `json:"tv_detail"`
}

func (f *File) getNfoFile() string {
	return f.Dir + "/" + f.getTitleWithoutSuffix() + ".nfo"
}

func (f *File) NfoExist() bool {
	nfo := f.getNfoFile()

	if info, err := os.Stat(nfo); err == nil && info.Size() > 0 {
		return true
	}

	return false
}

func (f *File) getTitleWithoutSuffix() string {
	return strings.Replace(f.OriginTitle, "."+f.Suffix, "", 1)
}

func (f *File) getCacheDir() string {
	return f.Dir + "/tmdb"
}

// 下载剧集的相关图片
func (f *File) downloadImage(d *tmdb.TvEpisodeDetail) {
	file := f.getTitleWithoutSuffix()
	if len(d.StillPath) > 0 {
		_ = tmdb.DownloadFile(tmdb.Api.GetImageOriginal(d.StillPath), f.Dir+"/"+file+"-thumb.jpg")
	}
}
