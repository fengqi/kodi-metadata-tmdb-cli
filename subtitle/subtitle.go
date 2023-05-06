package subtitle

import (
	"fengqi/kodi-metadata-tmdb-cli/utils"
	"path/filepath"
)

type FileInfo struct {
	Dir              string
	VideoFileName    string
	IsFile           bool
	IsSingleFile     bool
	TmdbId           int
	ImdbId           string
	OriginalTitle    string
	OriginalLanguage string
	Title            string
}

type SubtitleInfo struct {
	Url      string
	Filename string
}

type Subtitle interface {
	GetSubtitleList(fileInfo *FileInfo) (SubtitleInfoList []*SubtitleInfo, err error)
}

func DownloadSubtitle(fileInfo *FileInfo) error {
	sub := NewOpenSubtitles()

	list, err := sub.GetSubtitleList(fileInfo)
	if err != nil {
		utils.Logger.ErrorF("GetSubtitleList: %v", err)
		return err
	}

	for _, v := range list {
		utils.DownloadFile(v.Url, filepath.Join(fileInfo.Dir, v.Filename))
	}

	return nil
}
