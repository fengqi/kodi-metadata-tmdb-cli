package shows

import (
	"fengqi/kodi-metadata-tmdb-cli/utils"
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
		utils.Logger.DebugF("pass : %s", file.Name())
		return nil
	}

	fileName = utils.ReplaceChsNumber(fileName)

	// 提取季和集
	se, snum, enum := utils.MatchEpisode(fileName)
	if dir.Season > 0 {
		snum = dir.Season
	}
	utils.Logger.InfoF("find season: %d episode: %d %s", snum, enum, file.Name())
	if len(se) == 0 || snum == 0 || enum == 0 {
		utils.Logger.WarningF("seaon or episode not find: %s", file.Name())
		return nil
	}

	return &File{
		Dir:           dir.Dir + "/" + dir.OriginTitle,
		OriginTitle:   file.Name(),
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
	showName = utils.SeasonCorrecting(showName)

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
				if season != item { // TODO 这里假定只有名字和season在一起，没有其他特殊字符的情况，如：黄石S01，否则可能不适合这样处理
					showsDir.Title += strings.TrimRight(item, season) + " "
				}
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

	// 读特殊指定的值
	showsDir.ReadSeason()
	showsDir.ReadTvId()
	showsDir.ReadGroupId()

	return showsDir
}
