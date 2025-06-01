package shows

import (
	"fengqi/kodi-metadata-tmdb-cli/config"
	"fengqi/kodi-metadata-tmdb-cli/kodi"
	"fengqi/kodi-metadata-tmdb-cli/media_file"
	"fengqi/kodi-metadata-tmdb-cli/utils"
	"github.com/fengqi/lrace"
	"github.com/spf13/cast"
	"log"
	"path/filepath"
	"strings"
)

// Process 处理扫描到的电视剧文件
func Process(mf *media_file.MediaFile) error {
	show := &Show{MediaFile: mf}
	ParseShowFile(show, show.MediaFile.Path)

	show.checkTvCacheDir()
	detail, err := show.getTvDetail()
	if err != nil {
		return err
	}

	show.SaveTvNfo(detail)
	show.downloadTvImage(detail)

	show.checkCacheDir()
	episodeDetail, err := show.getEpisodeDetail()
	show.SaveEpisodeNfo(episodeDetail)
	show.downloadEpisodeImage(episodeDetail)

	if !detail.FromCache {
		kodi.Rpc.AddRefreshTask(kodi.TaskRefreshTVShow, detail.OriginalName)
	}

	return nil
}

func ParseShowFile(show *Show, parse string) error {
	// 递归到根目录
	for _, showsDir := range config.Collector.ShowsDir {
		if strings.TrimRight(showsDir, "/") == parse {
			// 读特殊指定的值
			show.checkCacheDir()
			show.checkTvCacheDir()
			show.ReadSeason()
			show.ReadTvId()
			//show.ReadGroupId()
			//show.ReadPart()

			return nil
		}
	}

	filename := filepath.Base(parse)

	// 过滤可选字符
	filename = utils.FilterOptionals(filename)

	// 过滤和替换中文
	filename = utils.ReplaceChsNumber(filename)

	// 过滤掉或替换歧义的内容
	filename = utils.SeasonCorrecting(filename)
	filename = utils.EpisodeCorrecting(filename)

	// 使用自定义方法切割
	split := utils.Split(filename)

	nameStart := false
	nameStop := false

	if show.Title != "" {
		nameStop = true
	}

	if source, season := utils.IsSeason(filename); len(season) > 0 && source == filename {
		split = split[0:0]
		show.Season = cast.ToInt(season)
		show.SeasonRoot = parse
	}

	for _, item := range split {
		if !nameStart && !nameStop {
			nameStart = true
			nameStop = false
		}

		if format := utils.IsFormat(item); len(format) > 0 {
			if show.Format == "" {
				show.Format = format
			}
			nameStop = true
			continue
		}

		if source := utils.IsSource(item); len(source) > 0 {
			if show.Source == "" {
				show.Source = source
			}
			nameStop = true
			continue
		}

		if studio := utils.IsStudio(item); len(studio) > 0 {
			show.Studio = lrace.Ternary(show.Studio == "", studio, show.Studio)
			nameStop = true
			continue
		}

		if channel := utils.IsChannel(item); len(channel) > 0 {
			show.Channel = lrace.Ternary(show.Channel == "", channel, show.Channel)
			nameStop = true
			continue
		}

		if coding := utils.IsVideoCoding(item); len(coding) > 0 {
			if show.VideoCoding == "" {
				show.VideoCoding = coding
			}
			nameStop = true
			continue
		}

		if coding := utils.IsAudioCoding(item); len(coding) > 0 {
			if show.AudioCoding == "" {
				show.AudioCoding = coding
			}
			nameStop = true
			continue
		}

		if crew := utils.IsCrew(item); len(crew) > 0 {
			show.Crew = lrace.Ternary(show.Crew == "", crew, show.Crew)
			nameStop = true
			continue
		}

		if dm := utils.IsDynamicRange(item); len(dm) > 0 {
			show.DynamicRange = lrace.Ternary(show.DynamicRange == "", dm, show.DynamicRange)
			nameStop = true
			continue
		}

		if year := utils.IsYear(item); year > 0 {
			if show.Year == 0 {
				show.Year = year
			}
			nameStop = true
			continue
		}

		if s, e := utils.MatchEpisode(item + show.MediaFile.Suffix); s > 0 && e > 0 {
			if show.Season == 0 {
				show.Season = s
			}
			if show.Episode == 0 {
				show.Episode = e
			}
			nameStop = true
			continue
		}

		if source, season := utils.IsSeason(item); len(season) > 0 {
			if show.Season == 0 {
				show.Season = cast.ToInt(season)
			}
			if source == filename { // 目录是季，如第x季、s02
				show.SeasonRoot = parse
				break
			}
			nameStop = true
			continue
		}

		if source, episode := utils.IsEpisode(item + show.MediaFile.Suffix); len(episode) > 0 {
			if show.Episode == 0 {
				show.Episode = cast.ToInt(episode)
			}
			if show.Episode > 100 {
				log.Println("what episode 50?")
			}
			if source == filename { // 文件名是集，如第x集、e02
				break
			}
			nameStop = true
			continue
		}

		if nameStart && !nameStop {
			show.Title += item + " "
		}
	}

	// 文件名清理
	show.Title, show.AliasTitle = utils.SplitTitleAlias(show.Title)
	show.ChsTitle, show.EngTitle = utils.SplitChsEngTitle(show.Title)

	show.TvRoot = filepath.Dir(parse)
	if lrace.InArray(config.Collector.ShowsDir, show.TvRoot) {
		show.TvRoot = parse
	}

	return ParseShowFile(show, filepath.Dir(parse))
}
