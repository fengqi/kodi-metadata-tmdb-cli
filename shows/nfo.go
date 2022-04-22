package shows

import (
	"fengqi/kodi-metadata-tmdb-cli/tmdb"
	"fengqi/kodi-metadata-tmdb-cli/utils"
	"strconv"
)

func (d *Dir) saveToNfo(detail *tmdb.TvDetail) error {
	utils.Logger.InfoF("save tvshow.nfo to: %s", d.getNfoFile())

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
			actor = append(actor, Actor{
				Name:  item.Name,
				Role:  item.Roles[0].Character,
				Order: item.Order,
				Thumb: tmdb.ImageW500 + item.ProfilePath,
			})
		}
	}

	namedSeason := make([]NamedSeason, 0)
	for _, item := range detail.Seasons {
		namedSeason = append(namedSeason, NamedSeason{
			Number: item.SeasonNumber,
			Value:  item.Name,
		})
	}

	top := &TvShowNfo{
		Title:         detail.Name,
		OriginalTitle: detail.OriginalName,
		ShowTitle:     detail.Name,
		SortTitle:     detail.Name,
		Plot:          detail.Overview,
		UniqueId: UniqueId{
			Type:    strconv.Itoa(detail.Id),
			Default: true,
		},
		Id:          detail.Id,
		Premiered:   detail.FirstAirDate,
		Ratings:     Ratings{Rating: rating},
		MPaa:        "TV-14",
		Status:      detail.Status,
		Genre:       genre,
		Studio:      studio,
		Season:      detail.NumberOfSeasons,
		Episode:     detail.NumberOfEpisodes,
		UserRating:  detail.VoteAverage,
		Actor:       actor,
		NamedSeason: namedSeason,
		FanArt: FanArt{
			Thumb: []ShowThumb{
				{
					Preview: tmdb.ImageW500 + detail.BackdropPath,
				},
			},
		},
	}

	return utils.SaveNfo(d.getNfoFile(), top)
}

// SaveTvEpisodeNFO 保存每集的信息到独立的NFO文件
func (f *File) saveToNfo(episode *tmdb.TvEpisodeDetail) error {
	utils.Logger.InfoF("save episode nfo to: %s", f.getNfoFile())

	actor := make([]Actor, 0)
	for _, item := range episode.GuestStars {
		actor = append(actor, Actor{
			Name:      item.Name,
			Role:      item.Character,
			Order:     item.Order,
			Thumb:     tmdb.ImageW500 + item.ProfilePath,
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
		MPaa:           "TV-14",
		Season:         episode.SeasonNumber,
		Episode:        episode.EpisodeNumber,
		DisplaySeason:  episode.SeasonNumber,
		DisplayEpisode: episode.EpisodeNumber,
		UserRating:     episode.VoteAverage,
		//Tagline:        "111",
		TmdbId:  "tmdm" + strconv.Itoa(episode.Id),
		Runtime: 6,
		Status:  "ok",
		Actor:   actor,
		Thumb: Thumb{
			Aspect:  "thumb",
			Preview: tmdb.ImageOriginal + episode.StillPath,
		},
		Ratings: rating,
		Aired:   episode.AirDate,
	}

	return utils.SaveNfo(f.getNfoFile(), top)
}
