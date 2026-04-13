package shows

import (
	"errors"
	"fengqi/kodi-metadata-tmdb-cli/common/ai"
	"fengqi/kodi-metadata-tmdb-cli/config"
	"fengqi/kodi-metadata-tmdb-cli/kodi"
	"fengqi/kodi-metadata-tmdb-cli/media_file"
	"fengqi/kodi-metadata-tmdb-cli/tmdb"
	"fengqi/kodi-metadata-tmdb-cli/utils"
	"fmt"
	"log"
	"path/filepath"
	"strings"

	"github.com/fengqi/lrace"
)

// Process 处理扫描到的电视剧文件
func Process(mf *media_file.MediaFile) error {
	show, detail, episodeDetail, err := loadShowCache(mf)
	if err != nil {
		return err
	}

	if show == nil || detail == nil || episodeDetail == nil {
		show, err = parseShowFile(mf)
		if err != nil {
			return err
		}
	}
	if show == nil {
		return errors.New("show file empty")
	}

	utils.Logger.DebugF("receive shows task: %v", show)

	show.checkTvCacheDir()
	if detail == nil {
		detail, err = show.getTvDetail()
		if err != nil {
			return err
		}
	}
	if detail == nil {
		return errors.New("get show detail empty")
	}

	_ = show.SaveTvNfo(detail)
	show.downloadTvImage(detail)

	show.checkCacheDir()
	if episodeDetail == nil {
		episodeDetail, err = show.getEpisodeDetail()
		if err != nil {
			return err
		}
	}
	if episodeDetail == nil {
		return errors.New("get show episode detail empty")
	}

	_ = show.SaveEpisodeNfo(episodeDetail)
	show.downloadEpisodeImage(episodeDetail)

	if !detail.FromCache {
		kodi.Rpc.AddRefreshTask(kodi.TaskRefreshTVShow, detail.Name)
	}

	if !episodeDetail.FromCache {
		taskVal := fmt.Sprintf("%s|-|%d|-|%d", detail.Name, episodeDetail.SeasonNumber, episodeDetail.EpisodeNumber)
		kodi.Rpc.AddRefreshTask(kodi.TaskRefreshEpisode, taskVal)
	}

	kodi.Rpc.AddScanTaskByName(1, detail.Name)
	return nil
}

func parseShowFile(mf *media_file.MediaFile) (*Show, error) {
	if mf == nil {
		return nil, errors.New("show media file nil")
	}

	if !ai.Enabled() {
		return parseShowFileByRule(mf)
	}

	if config.Ai.MatchMode == config.AiMatchModeRuleThenAi {
		ruleShow, ruleErr := parseShowFileByRule(mf)
		if ruleErr == nil && ruleShow != nil && ruleShow.Title != "" && ruleShow.Episode > 0 {
			return ruleShow, nil
		}

		aiShow, aiErr := parseShowFileByAI(mf, ruleShow)
		if aiErr == nil {
			return aiShow, nil
		}
		return ruleShow, ruleErr
	}

	if config.Ai.MatchMode == config.AiMatchModeAiThenRule {
		aiShow, aiErr := parseShowFileByAI(mf, nil)
		if aiErr == nil {
			return aiShow, nil
		}
		return parseShowFileByRule(mf)
	}

	ruleShow, ruleErr := parseShowFileByRule(mf)
	if ruleErr != nil {
		return nil, ruleErr
	}
	aiShow, aiErr := parseShowFileByAI(mf, ruleShow)
	if aiErr == nil {
		return aiShow, nil
	}
	return ruleShow, nil
}

func parseShowFileByRule(mf *media_file.MediaFile) (*Show, error) {
	if mf == nil {
		return nil, errors.New("show media file nil")
	}
	show := &Show{MediaFile: mf}
	fillShowPathMeta(show)
	return show, ParseShowFile(show, show.MediaFile.Path)
}

func parseShowFileByAI(mf *media_file.MediaFile, ruleShow *Show) (*Show, error) {
	if !ai.Enabled() {
		return nil, errors.New("ai disabled")
	}
	if mf == nil {
		return nil, errors.New("show media file nil")
	}
	show := &Show{MediaFile: mf}
	fillShowPathMeta(show)

	rule := map[string]any{}
	if ruleShow != nil {
		rule = map[string]any{
			"title":       ruleShow.Title,
			"alias_title": ruleShow.AliasTitle,
			"chs_title":   ruleShow.ChsTitle,
			"eng_title":   ruleShow.EngTitle,
			"year":        ruleShow.Year,
			"season":      ruleShow.Season,
			"episode":     ruleShow.Episode,
		}
	}

	result, err := ai.ParseMedia(&ai.ParseInput{
		MediaType: "tv",
		Path:      mf.Path,
		Filename:  mf.Filename,
		Rule:      rule,
	})
	if err != nil {
		return nil, err
	}
	if !ai.ParseUsable(result) {
		return nil, errors.New("ai parse unusable")
	}

	show.Title = result.Title
	show.AliasTitle = result.AliasTitle
	show.ChsTitle = result.ChsTitle
	show.EngTitle = result.EngTitle
	show.Year = result.Year
	show.Season = result.Season
	show.Episode = result.Episode

	if show.Title == "" {
		return nil, errors.New("ai parse title empty")
	}

	utils.Logger.DebugF("ai parse show success: %s S%02dE%02d confidence %.2f", show.Title, show.Season, show.Episode, result.Confidence)
	return show, nil
}

func fillShowPathMeta(show *Show) {
	if show == nil || show.MediaFile == nil {
		return
	}

	cursor := show.MediaFile.Path
	for {
		parent := filepath.Dir(cursor)

		for _, showsDir := range config.Collector.ShowsDir {
			if utils.PathEqual(showsDir, parent) {
				show.TvRoot = cursor
				goto foundRoot
			}
		}

		if utils.PathEqual(parent, cursor) {
			break
		}
		cursor = parent
	}

foundRoot:
	if show.TvRoot == "" {
		show.TvRoot = filepath.Dir(show.MediaFile.Path)
	}
	if show.SeasonRoot == "" {
		show.SeasonRoot = show.TvRoot
	}
}

func ParseShowFile(show *Show, parse string) error {
	if show == nil {
		return errors.New("show nil")
	}

	// 递归到根目录
	parseNorm := utils.NormalizePath(parse)
	for _, showsDir := range config.Collector.ShowsDir {
		if utils.NormalizePath(showsDir) == parseNorm {
			if show.Episode > 0 && show.Season == 0 {
				show.Season = 1
			}

			if show.SeasonRoot == "" {
				show.SeasonRoot = show.TvRoot
			}

			// 读特殊指定的值
			show.checkCacheDir()
			show.checkTvCacheDir()
			show.ReadSeason()
			show.ReadTvId()
			show.ReadGroupId()
			// show.ReadPart()
			show.ReadJoin()

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
	split := utils.Split(strings.TrimRight(filename, show.MediaFile.Suffix))

	nameStart := false
	nameStop := false

	if show.Title != "" {
		nameStop = true
	}

	if lrace.IsDir(parse) {
		if source, season := utils.IsSeason(filename); len(season) > 0 && source == filename {
			split = split[0:0]
			show.Season = utils.StrToInt(season)
			show.SeasonRoot = parse
		}
	}

	for _, item := range split {
		if lrace.InArray(config.Collector.SkipKeywords, item) {
			continue
		}

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
				show.Season = utils.StrToInt(season)
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
				show.Episode = utils.StrToInt(episode)
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
	show.Title = strings.TrimSpace(show.Title)
	show.Title, show.AliasTitle = utils.SplitTitleAlias(show.Title)
	show.ChsTitle, show.EngTitle = utils.SplitChsEngTitle(show.Title)

	parent := filepath.Dir(parse)
	if utils.NormalizePath(parent) == parseNorm {
		return nil
	}

	return ParseShowFile(show, parent)
}

// loadShowCache 从缓存中加载电视剧信息
func loadShowCache(mf *media_file.MediaFile) (*Show, *tmdb.TvDetail, *tmdb.TvEpisodeDetail, error) {
	if mf == nil {
		return nil, nil, nil, errors.New("show media file nil")
	}

	show := &Show{MediaFile: mf}
	fillShowPathMeta(show)
	show.ReadTvId()
	if show.TvId == 0 {
		return nil, nil, nil, nil
	}

	show.checkTvCacheDir()
	detail, err := show.loadTvDetailFromCache()
	if err != nil {
		return nil, nil, nil, err
	}
	if detail == nil {
		return nil, nil, nil, nil
	}

	show.checkCacheDir()
	episodeDetail, err := show.loadEpisodeDetailFromCache()
	if err != nil {
		return nil, nil, nil, err
	}
	if episodeDetail == nil {
		return show, detail, nil, nil
	}

	return show, detail, episodeDetail, nil
}
