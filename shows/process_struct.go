package shows

import (
	"fengqi/kodi-metadata-tmdb-cli/common/constants"
	"fengqi/kodi-metadata-tmdb-cli/media_file"
	"fengqi/kodi-metadata-tmdb-cli/utils"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// Show 电视剧
type Show struct {
	MediaFile    *media_file.MediaFile `json:"media_file"`  // 媒体文件
	TvRoot       string                `json:"tv_root"`     // 电视剧跟目录
	SeasonRoot   string                `json:"season_root"` // 季目录
	TvId         int                   `json:"tv_id"`       // TMDb tv id
	GroupId      string                `json:"group_id"`    // TMDB Episode Group
	Season       int                   `json:"season"`      // 第几季 ，电影类 -1
	Episode      int                   `json:"episode"`     // 第几集，电影类 -1
	Title        string                `json:"title"`       // 从视频提取的文件名 鹰眼 Hawkeye
	AliasTitle   string                `json:"alias_title"` // 别名，通常没有用
	ChsTitle     string                `json:"chs_title"`   // 分离出来的中文名称 鹰眼
	EngTitle     string                `json:"eng_title"`   // 分离出来的英文名称 Hawkeye
	Year         int                   `json:"year"`        // 年份：2020、2021
	Format       string                `json:"format"`
	VideoCoding  string                `json:"video_coding"`
	AudioCoding  string                `json:"audio_coding"`
	Source       string                `json:"source"`
	Studio       string                `json:"studio"`
	Channel      string                `json:"channel"`
	Crew         string                `json:"crew"`
	DynamicRange string                `json:"dynamic_range"`
}

func (s *Show) checkCacheDir() {
	dir := s.GetCacheDir()
	if _, err := os.Stat(dir); err != nil && os.IsNotExist(err) {
		if err := os.Mkdir(dir, 0755); err != nil {
			utils.Logger.ErrorF("create cache: %s dir err: %v", dir, err)
		}
	}
}

func (s *Show) checkTvCacheDir() {
	dir := s.GetTvCacheDir()
	if _, err := os.Stat(dir); err != nil && os.IsNotExist(err) {
		if err := os.Mkdir(dir, 0755); err != nil {
			utils.Logger.ErrorF("create cache: %s dir err: %v", dir, err)
		}
	}
}

func (s *Show) GetTvCacheDir() string {
	return s.TvRoot + "/" + constants.TmdbCacheDir
}

func (s *Show) GetCacheDir() string {
	base := filepath.Dir(s.MediaFile.Path)
	return base + "/" + constants.TmdbCacheDir
}

func (s *Show) GetFullDir() string {
	return s.MediaFile.Path
}

func (s *Show) ReadSeason() {
	seasonFile := s.GetCacheDir() + "/season.txt"
	if _, err := os.Stat(seasonFile); err == nil {
		bytes, err := os.ReadFile(seasonFile)
		if err == nil {
			s.Season, _ = strconv.Atoi(strings.Trim(string(bytes), "\r\n "))
		} else {
			utils.Logger.WarningF("read season specially file: %s err: %v", seasonFile, err)
		}
	}
}

// ReadTvId 从文件读取tvId
func (s *Show) ReadTvId() {
	idFile := s.TvRoot + "/tmdb/id.txt"
	if _, err := os.Stat(idFile); err == nil {
		bytes, err := os.ReadFile(idFile)
		if err == nil {
			s.TvId, _ = strconv.Atoi(strings.Trim(string(bytes), "\r\n "))
		} else {
			utils.Logger.WarningF("read tv id specially file: %s err: %v", idFile, err)
		}
	}
}

// CacheTvId 缓存tvId到文件
func (s *Show) CacheTvId() {
	idFile := s.TvRoot + "/tmdb/id.txt"
	err := os.WriteFile(idFile, []byte(strconv.Itoa(s.TvId)), 0664)
	if err != nil {
		utils.Logger.ErrorF("save tvId %d to %s err: %v", s.TvId, idFile, err)
	}
}

// ReadGroupId 从文件读取剧集分组
func (s *Show) ReadGroupId() {
	groupFile := s.SeasonRoot + "/tmdb/group.txt"
	if _, err := os.Stat(groupFile); err == nil {
		bytes, err := os.ReadFile(groupFile)
		if err == nil {
			s.GroupId = strings.Trim(string(bytes), "\r\n ")
		} else {
			utils.Logger.WarningF("read group id specially file: %s err: %v", groupFile, err)
		}
	}
}
