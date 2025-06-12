package movies

import (
	"fengqi/kodi-metadata-tmdb-cli/config"
	"fengqi/kodi-metadata-tmdb-cli/tmdb"
	"fengqi/kodi-metadata-tmdb-cli/utils"
	"strconv"
	"strings"
)

// 保存NFO文件
func (m *Movie) saveToNfo(detail *tmdb.MovieDetail) error {
	utils.Logger.InfoF("save movie nfo to: %s", m.NfoFile)

	genre := make([]string, 0)
	for _, item := range detail.Genres {
		genre = append(genre, item.Name)
	}

	studio := make([]string, 0)
	for _, item := range detail.ProductionCompanies {
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
	if detail.Credits != nil {
		for _, item := range detail.Credits.Cast {
			if item.ProfilePath == "" {
				continue
			}

			actor = append(actor, Actor{
				Name:      item.Name,
				Role:      item.Character,
				Order:     item.Order,
				Thumb:     tmdb.Api.GetImageW500(item.ProfilePath),
				SortOrder: item.CastId,
			})
		}
	}

	mpaa := "NR"
	contentRating := strings.ToUpper(config.Tmdb.Rating)
	if detail.Releases.Countries != nil && len(detail.Releases.Countries) > 0 {
		mpaa = detail.Releases.Countries[0].Certification
		for _, item := range detail.Releases.Countries {
			if strings.ToUpper(item.ISO31661) == contentRating {
				mpaa = item.Certification
				break
			}
		}
	}

	var fanArt *FanArt
	if detail.BackdropPath != "" {
		fanArt = &FanArt{
			Thumb: []MovieThumb{
				{
					Preview: tmdb.Api.GetImageW500(detail.BackdropPath),
				},
			},
		}
	}

	year := ""
	if detail.ReleaseDate != "" {
		year = detail.ReleaseDate[:4]
	}

	country := make([]string, 0)
	for _, item := range detail.ProductionCountries {
		country = append(country, item.Name) // todo 使用 iso_3166_1 匹配中文
	}

	languages := make([]string, 0)
	for _, item := range detail.SpokenLanguages {
		languages = append(languages, item.Name) // todo 使用 iso_639_1 匹配中文
	}

	top := &MovieNfo{
		Title:         detail.Title,
		OriginalTitle: detail.OriginalTitle,
		SortTitle:     detail.Title,
		Plot:          detail.Overview,
		UniqueId: UniqueId{
			Default: true,
			Type:    "tmdb",
			Value:   strconv.Itoa(detail.Id),
		},
		Id:         detail.Id,
		Premiered:  detail.ReleaseDate,
		Ratings:    Ratings{Rating: rating},
		MPaa:       mpaa,
		Year:       year,
		Status:     detail.Status,
		Genre:      genre,
		Tag:        genre,
		Country:    country,
		Languages:  languages,
		Studio:     studio,
		UserRating: detail.VoteAverage,
		Actor:      actor,
		FanArt:     fanArt,
	}

	return utils.SaveNfo(m.NfoFile, top)
}
