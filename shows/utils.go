package shows

import (
	"fengqi/kodi-metadata-tmdb-cli/tmdb"
	"fengqi/kodi-metadata-tmdb-cli/utils"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// 通过一个绝对路径的文件分析剧集详情
func getFileInfoByFile(file string) *File {
	f, _ := os.Stat(file)
	return parseShowsFile(getDirInfoByFile(file), f)
}

// 通过一个绝对路径的文件分析剧集目录详情
func getDirInfoByFile(file string) *Dir {
	baseDir := filepath.Dir(file)
	fileInfo, _ := os.Stat(baseDir)
	return parseShowsDir(filepath.Dir(baseDir), fileInfo)
}

// 解析文件, 返回详情
func parseShowsFile(dir *Dir, file fs.FileInfo) *File {
	fileName := utils.FilterTmpSuffix(file.Name())

	// 判断是视频, 并获取后缀
	suffix := utils.IsVideo(fileName)
	if len(suffix) == 0 {
		utils.Logger.DebugF("pass : %s", fileName)
		return nil
	}

	// 提取季和集
	se, snum, enum := utils.MatchEpisode(fileName)
	if dir.Season > 0 {
		snum = dir.Season
	}
	utils.Logger.InfoF("find season: %d episode: %d %s", snum, enum, fileName)
	if len(se) == 0 || snum == 0 || enum == 0 {
		utils.Logger.WarningF("seaon or episode not find: %s", fileName)
		return nil
	}

	return &File{
		Dir:           dir.Dir + "/" + dir.OriginTitle,
		OriginTitle:   fileName,
		Season:        snum,
		Episode:       enum,
		SeasonEpisode: se,
		Suffix:        suffix,
		TvId:          dir.TvId,
	}
}

// 解析目录, 返回详情
func parseShowsDir(baseDir string, file fs.FileInfo) *Dir {
	showName := file.Name()

	// 过滤无用文件
	if showName == "@eaDir" || showName == "tmdb" || showName == "metadata" || showName[0:1] == "." {
		utils.Logger.DebugF("pass file: %s", showName)
		return nil
	}

	// 过滤可选字符
	showName = utils.FilterOptionals(showName)

	// 过滤掉或替换歧义的内容
	showName = utils.FilterCorrecting(showName)

	// 过滤掉分段的干扰
	if subEpisodes := utils.IsSubEpisodes(showName); subEpisodes != "" {
		showName = strings.Replace(showName, subEpisodes, "", 1)
	}

	showsDir := &Dir{
		Dir:          baseDir,
		OriginTitle:  file.Name(),
		IsCollection: utils.IsCollection(file.Name()),
	}

	// 年份范围
	if yearRange := utils.IsYearRange(showName); len(yearRange) > 0 {
		showsDir.YearRange = yearRange
		showName = strings.Replace(showName, yearRange, "", 1)
	}

	// 使用自定义方法切割
	split := utils.Split(showName)

	nameStart := false
	nameStop := false
	for _, item := range split {
		if year := utils.IsYear(item); year > 0 {
			// 名字带年的，比如 reply 1994
			if showsDir.Year == 0 {
				showsDir.Year = year
			} else {
				showsDir.Title += strconv.Itoa(showsDir.Year)
				showsDir.Year = year
			}
			nameStop = true
			continue
		}

		if season := utils.IsSeason(item); len(season) > 0 {
			if !showsDir.IsCollection {
				s := season[1:]
				i, err := strconv.Atoi(s)
				if err == nil {
					showsDir.Season = i
					nameStop = true
				}
			}
			continue
		}

		if format := utils.IsFormat(item); len(format) > 0 {
			showsDir.Format = format
			nameStop = true
			continue
		}

		if source := utils.IsSource(item); len(source) > 0 {
			showsDir.Source = source
			nameStop = true
			continue
		}

		if studio := utils.IsStudio(item); len(studio) > 0 {
			showsDir.Studio = studio
			nameStop = true
			continue
		}

		if channel := utils.IsChannel(item); len(channel) > 0 {
			nameStop = true
			continue
		}

		if !nameStart {
			nameStart = true
			nameStop = false
		}

		if !nameStop {
			showsDir.Title += item + " "
		}
	}

	// 文件名清理
	showsDir.Title, showsDir.AliasTitle = utils.SplitTitleAlias(showsDir.Title)
	showsDir.ChsTitle, showsDir.EngTitle = utils.SplitChsEngTitle(showsDir.Title)
	if len(showsDir.Title) == 0 {
		utils.Logger.WarningF("file: %s parse title empty: %v", file.Name(), showsDir)
		return nil
	}

	showsDir.ReadSeason()
	showsDir.ReadTvId()
	showsDir.ReadGroupId()

	return showsDir
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

func (f *File) getTitleWithoutSuffix() string {
	return strings.Replace(f.OriginTitle, "."+f.Suffix, "", 1)
}

func (d *Dir) GetCacheDir() string {
	return d.GetFullDir() + "/tmdb"
}

func (d *Dir) GetFullDir() string {
	return d.Dir + "/" + d.OriginTitle
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

func (f *File) getCacheDir() string {
	return f.Dir + "/tmdb"
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

// 下载剧集的相关图片
func (f *File) downloadImage(d *tmdb.TvEpisodeDetail) {
	file := f.getTitleWithoutSuffix()
	if len(d.StillPath) > 0 {
		_ = utils.DownloadFile(tmdb.Api.GetImageOriginal(d.StillPath), f.Dir+"/"+file+"-thumb.jpg")
	}
}
