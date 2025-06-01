package shows

import (
	"fengqi/kodi-metadata-tmdb-cli/config"
	"fengqi/kodi-metadata-tmdb-cli/tmdb"
	"fengqi/kodi-metadata-tmdb-cli/utils"
	"github.com/fengqi/lrace"
	"strconv"
	"strings"
)

func (s *Show) SaveTvNfo(detail *tmdb.TvDetail) error {
	tvNfo := s.TvRoot + "/tvshow.nfo"
	if detail.FromCache && lrace.FileExist(tvNfo) {
		return nil
	}

	utils.Logger.InfoF("save tvshow.nfo to: %s", tvNfo)

	genre := make([]string, 0)
	for _, item := range detail.Genres {
		genre = append(genre, item.Name)
	}

	studio := make([]string, 0)
	for _, item := range detail.Networks {
		studio = append(studio, item.Name)
	}

	rating := make([]Rating, 1)
	rating[0] = Rating{
		Name:  "tmdb",
		Max:   10,
		Value: detail.VoteAverage,
		Votes: detail.VoteCount,
	}

	actor := make([]Actor, 0)
	if detail.AggregateCredits != nil {
		for _, item := range detail.AggregateCredits.Cast {
			if item.ProfilePath == "" {
				continue
			}

			actor = append(actor, Actor{
				Name:  item.Name,
				Role:  item.Roles[0].Character,
				Order: item.Order,
				Thumb: tmdb.Api.GetImageW500(item.ProfilePath),
			})
		}
	}

	episodeCount := 0
	namedSeason := make([]NamedSeason, 0)
	for _, item := range detail.Seasons {
		//if !s.IsCollection && item.SeasonNumber != s.Season {
		//	continue
		//}
		namedSeason = append(namedSeason, NamedSeason{
			Number: item.SeasonNumber,
			Value:  item.Name,
		})
		episodeCount = item.EpisodeCount
	}

	mpaa := "NR"
	contentRating := strings.ToUpper(config.Tmdb.Rating)
	if detail.ContentRatings != nil && len(detail.ContentRatings.Results) > 0 {
		mpaa = detail.ContentRatings.Results[0].Rating
		for _, item := range detail.ContentRatings.Results {
			if strings.ToUpper(item.ISO31661) == contentRating {
				mpaa = item.Rating
				break
			}
		}
	}

	var fanArt *FanArt
	if detail.BackdropPath != "" {
		fanArt = &FanArt{
			Thumb: []ShowThumb{
				{
					Preview: tmdb.Api.GetImageW500(detail.BackdropPath),
				},
			},
		}
	}

	top := &TvShowNfo{
		Title:         detail.Name,
		OriginalTitle: detail.OriginalName,
		ShowTitle:     detail.Name,
		SortTitle:     detail.Name,
		Plot:          detail.Overview,
		UniqueId: UniqueId{
			Type:    "tmdb",
			Default: true,
			Value:   strconv.Itoa(detail.Id),
		},
		Id:          detail.Id,
		Premiered:   detail.FirstAirDate,
		Ratings:     Ratings{Rating: rating},
		MPaa:        mpaa,
		Status:      detail.Status,
		Genre:       genre,
		Studio:      studio,
		Season:      s.Season,
		Episode:     episodeCount,
		UserRating:  detail.VoteAverage,
		Actor:       actor,
		NamedSeason: namedSeason,
		FanArt:      fanArt,
	}

	// 使用分组信息
	if s.GroupId != "" && detail.TvEpisodeGroupDetail != nil {
		top.Season = detail.TvEpisodeGroupDetail.GroupCount
		top.Episode = detail.TvEpisodeGroupDetail.EpisodeCount

		namedSeason = make([]NamedSeason, 0)
		for _, item := range detail.TvEpisodeGroupDetail.Groups {
			namedSeason = append(namedSeason, NamedSeason{
				Number: item.Order,
				Value:  item.Name,
			})
		}
		top.NamedSeason = namedSeason
	}

	return utils.SaveNfo(tvNfo, top)
}

// SaveEpisodeNfo 保存每集的信息到独立的NFO文件
func (s *Show) SaveEpisodeNfo(episode *tmdb.TvEpisodeDetail) error {
	episodeNfo := strings.Replace(s.MediaFile.Path, s.MediaFile.Suffix, ".nfo", 1)
	if episode.FromCache && lrace.FileExist(episodeNfo) {
		return nil
	}

	utils.Logger.InfoF("save episode nfo to: %s", episodeNfo)

	actor := make([]Actor, 0)
	for _, item := range episode.GuestStars {
		actor = append(actor, Actor{
			Name:      item.Name,
			Role:      item.Character,
			Order:     item.Order,
			Thumb:     tmdb.Api.GetImageW500(item.ProfilePath),
			SortOrder: item.Order,
		})
	}

	// 评分
	rating := make([]Rating, 1)
	rating[0] = Rating{
		Name:  "tmdb",
		Max:   10,
		Value: episode.VoteAverage,
		Votes: episode.VoteCount,
	}

	top := &TvEpisodeNfo{
		Title:         episode.Name,
		ShowTitle:     episode.Name,
		OriginalTitle: episode.Name,
		Plot:          episode.Overview,
		UniqueId: UniqueId{
			Type:    strconv.Itoa(episode.Id),
			Default: true,
		},
		Premiered:      episode.AirDate,
		Season:         episode.SeasonNumber,
		Episode:        episode.EpisodeNumber,
		DisplaySeason:  episode.SeasonNumber,
		DisplayEpisode: episode.EpisodeNumber,
		UserRating:     episode.VoteAverage,
		TmdbId:         "tmdb" + strconv.Itoa(episode.Id),
		Runtime:        6,
		Actor:          actor,
		Thumb: Thumb{
			Aspect:  "thumb",
			Preview: tmdb.Api.GetImageOriginal(episode.StillPath),
		},
		Ratings: rating,
		Aired:   episode.AirDate,
	}

	return utils.SaveNfo(episodeNfo, top)
}
