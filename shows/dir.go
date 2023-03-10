package shows

import (
	"fengqi/kodi-metadata-tmdb-cli/tmdb"
	"fengqi/kodi-metadata-tmdb-cli/utils"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func (d *Dir) ReadTvId() {
	idFile := d.Dir + "/" + d.OriginTitle + "/tmdb/id.txt"
	if _, err := os.Stat(idFile); err == nil {
		bytes, err := os.ReadFile(idFile)
		if err == nil {
			d.TvId, _ = strconv.Atoi(strings.Trim(string(bytes), "\r\n "))
		} else {
			utils.Logger.WarningF("read tv id specially file: %s err: %v", idFile, err)
		}
	}
}

func (d *Dir) ReadSeason() {
	seasonFile := d.Dir + "/" + d.OriginTitle + "/tmdb/season.txt"
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

func (d *Dir) ReadGroupId() {
	groupFile := d.Dir + "/" + d.OriginTitle + "/tmdb/group.txt"
	if _, err := os.Stat(groupFile); err == nil {
		bytes, err := os.ReadFile(groupFile)
		if err == nil {
			d.GroupId = strings.Trim(string(bytes), "\r\n ")
		} else {
			utils.Logger.WarningF("read group id specially file: %s err: %v", groupFile, err)
		}
	}
}

func (d *Dir) GetCacheDir() string {
	return d.GetFullDir() + "/tmdb"
}

func (d *Dir) GetFullDir() string {
	return d.Dir + "/" + d.OriginTitle
}

func (d *Dir) getNfoFile() string {
	return d.GetFullDir() + "/tvshow.nfo"
}

func (d *Dir) NfoExist() bool {
	nfo := d.getNfoFile()

	if info, err := os.Stat(nfo); err == nil && info.Size() > 0 {
		return true
	}

	return false
}

// CheckCacheDir tmdb 缓存目录
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
		_ = utils.DownloadFile(tmdb.Api.GetImageOriginal(detail.PosterPath), d.GetFullDir()+"/poster.jpg")
	}

	if len(detail.BackdropPath) > 0 {
		_ = utils.DownloadFile(tmdb.Api.GetImageOriginal(detail.BackdropPath), d.GetFullDir()+"/fanart.jpg")
	}

	// TODO group的信息里可能 season poster不全
	if len(detail.Seasons) > 0 {
		for _, item := range detail.Seasons {
			if !d.IsCollection && item.SeasonNumber != d.Season || item.PosterPath == "" {
				continue
			}
			seasonPoster := fmt.Sprintf("season%02d-poster.jpg", item.SeasonNumber)
			_ = utils.DownloadFile(tmdb.Api.GetImageOriginal(item.PosterPath), d.GetFullDir()+"/"+seasonPoster)
		}
	}
}
