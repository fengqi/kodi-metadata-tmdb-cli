package shows

import (
	"fengqi/kodi-metadata-tmdb-cli/tmdb"
	"fengqi/kodi-metadata-tmdb-cli/utils"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

// Dir 电视剧目录详情，从名字分析
// World.Heritage.In.China.E01-E38.2008.CCTVHD.x264.AC3.720p-CMCT
type Dir struct {
	Dir          string `json:"dir"`
	OriginTitle  string `json:"origin_title"`  // 原始文件名
	Title        string `json:"title"`         // 从视频提取的文件名 鹰眼 Hawkeye
	AliasTitle   string `json:"alias_title"`   // 别名，通常没有用
	ChsTitle     string `json:"chs_title"`     // 分离出来的中文名称 鹰眼
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
	PartMode     int    `json:"part_mode"`     // 分卷模式: 0不使用分卷, 1-自动, 2以上为手动指定分卷数量
}

// ReadTvId 从文件读取tvId
func (d *Dir) ReadTvId() {
	idFile := d.GetCacheDir() + "/id.txt"
	if _, err := os.Stat(idFile); err == nil {
		bytes, err := os.ReadFile(idFile)
		if err == nil {
			d.TvId, _ = strconv.Atoi(strings.Trim(string(bytes), "\r\n "))
		} else {
			utils.Logger.WarningF("read tv id specially file: %s err: %v", idFile, err)
		}
	}
}

// CacheTvId 缓存tvId到文件
func (d *Dir) CacheTvId() {
	idFile := d.GetCacheDir() + "/id.txt"
	err := os.WriteFile(idFile, []byte(strconv.Itoa(d.TvId)), 0664)
	if err != nil {
		utils.Logger.ErrorF("save tvId %d to %s err: %v", d.TvId, idFile, err)
	}
}

// ReadSeason 从文件读取季
func (d *Dir) ReadSeason() {
	seasonFile := d.GetCacheDir() + "/season.txt"
	if _, err := os.Stat(seasonFile); err == nil {
		bytes, err := os.ReadFile(seasonFile)
		if err == nil {
			d.Season, _ = strconv.Atoi(strings.Trim(string(bytes), "\r\n "))
		} else {
			utils.Logger.WarningF("read season specially file: %s err: %v", seasonFile, err)
		}
	}

	if d.Season == 0 && len(d.YearRange) == 0 {
		d.Season = 1
	}
}

// ReadGroupId 从文件读取剧集分组
func (d *Dir) ReadGroupId() {
	groupFile := d.GetCacheDir() + "/group.txt"
	if _, err := os.Stat(groupFile); err == nil {
		bytes, err := os.ReadFile(groupFile)
		if err == nil {
			d.GroupId = strings.Trim(string(bytes), "\r\n ")
		} else {
			utils.Logger.WarningF("read group id specially file: %s err: %v", groupFile, err)
		}
	}
}

// GetCacheDir 获取TMDB信息缓存目录, 通常是在每部电视剧的根目录下
func (d *Dir) GetCacheDir() string {
	return d.GetFullDir() + "/tmdb"
}

// GetFullDir 获取电视剧的完整目录
func (d *Dir) GetFullDir() string {
	return d.Dir + "/" + d.OriginTitle
}

// GetNfoFile 获取电视剧的NFO文件路径
func (d *Dir) GetNfoFile() string {
	return d.GetFullDir() + "/tvshow.nfo"
}

// NfoExist 判断NFO文件是否存在
func (d *Dir) NfoExist() bool {
	nfo := d.GetNfoFile()

	if info, err := os.Stat(nfo); err == nil && info.Size() > 0 {
		return true
	}

	return false
}

// CheckCacheDir 检查并创建缓存目录
func (d *Dir) checkCacheDir() {
	dir := d.GetCacheDir()
	if _, err := os.Stat(dir); err != nil && os.IsNotExist(err) {
		err := os.Mkdir(dir, 0755)
		if err != nil {
			utils.Logger.ErrorF("create cache: %s dir err: %v", dir, err)
		}
	}
}

// 下载电视剧的相关图片
// TODO 下载失败后，没有重复以及很长一段时间都不会再触发下载
func (d *Dir) downloadImage(detail *tmdb.TvDetail) {
	utils.Logger.DebugF("download %s images", d.Title)

	if len(detail.PosterPath) > 0 {
		_ = tmdb.DownloadFile(tmdb.Api.GetImageOriginal(detail.PosterPath), d.GetFullDir()+"/poster.jpg")
	}

	if len(detail.BackdropPath) > 0 {
		_ = tmdb.DownloadFile(tmdb.Api.GetImageOriginal(detail.BackdropPath), d.GetFullDir()+"/fanart.jpg")
	}

	// TODO group的信息里可能 season poster不全
	if len(detail.Seasons) > 0 {
		for _, item := range detail.Seasons {
			if !d.IsCollection && item.SeasonNumber != d.Season || item.PosterPath == "" {
				continue
			}
			seasonPoster := fmt.Sprintf("season%02d-poster.jpg", item.SeasonNumber)
			_ = tmdb.DownloadFile(tmdb.Api.GetImageOriginal(item.PosterPath), d.GetFullDir()+"/"+seasonPoster)
		}
	}

	if detail.Images != nil && len(detail.Images.Logos) > 0 {
		sort.SliceStable(detail.Images.Logos, func(i, j int) bool {
			return detail.Images.Logos[i].VoteAverage > detail.Images.Logos[j].VoteAverage
		})
		image := detail.Images.Logos[0]
		for _, item := range detail.Images.Logos {
			if image.FilePath == "" && item.FilePath != "" {
				image = item
			}
			if item.Iso6391 == "zh" && image.Iso6391 != "zh" {
				image = item
				break
			}
		}
		if image.FilePath != "" {
			logoFile := d.GetFullDir() + "/clearlogo.png"
			_ = tmdb.DownloadFile(tmdb.Api.GetImageOriginal(image.FilePath), logoFile)
		}
	}
}

// ReadPart 读取分卷模式
func (d *Dir) ReadPart() {
	partFile := d.GetCacheDir() + "/part.txt"
	if _, err := os.Stat(partFile); err == nil {
		bytes, err := os.ReadFile(partFile)
		if err == nil {
			d.PartMode, _ = strconv.Atoi(strings.Trim(string(bytes), "\r\n "))
		} else {
			utils.Logger.WarningF("read part specially file: %s err: %v", partFile, err)
		}
	}
}
