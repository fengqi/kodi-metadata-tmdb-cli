package shows

import (
	"fengqi/kodi-metadata-tmdb-cli/tmdb"
	"fengqi/kodi-metadata-tmdb-cli/utils"
	"os"
	"strings"
)

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
		_ = utils.DownloadFile(tmdb.Api.GetImageOriginal(d.StillPath), f.Dir+"/"+file+"-thumb.jpg")
	}
}
