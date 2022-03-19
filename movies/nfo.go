package movies

import (
	"fengqi/kodi-metadata-tmdb-cli/tmdb"
	"fengqi/kodi-metadata-tmdb-cli/utils"
	"strconv"
)

func (d *Movie) saveToNfo(detail *tmdb.MovieDetail, mode int) error {
	nfoFile := d.getNfoFile(mode)
	if nfoFile == "" {
		utils.Logger.InfoF("movie nfo empty %v", d)
		return nil
	}

	utils.Logger.InfoF("save movie nfo to: %s", nfoFile)

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
			actor = append(actor, Actor{
				Name:      item.Name,
				Role:      item.Character,
				Order:     item.Order,
				Thumb:     tmdb.ImageW500 + item.ProfilePath,
				SortOrder: item.CastId,
			})
		}
	}

	mpaa := "NR"
	if detail.Releases.Countries != nil {
		mpaa = detail.Releases.Countries[0].Certification
		for _, item := range detail.Releases.Countries {
			if item.ISO31661 == collector.config.Rating {
				mpaa = item.Certification
			}
		}
	}

	top := &MovieNfo{
		Title:         detail.Title,
		OriginalTitle: detail.OriginalTitle,
		SortTitle:     detail.Title,
		Plot:          detail.Overview,
		UniqueId: UniqueId{
			Type:    "tmdb",
			Default: true,
			Value:   strconv.Itoa(detail.Id),
		},
		Id:         detail.Id,
		Premiered:  detail.ReleaseDate,
		Ratings:    Ratings{Rating: rating},
		MPaa:       mpaa,
		Status:     detail.Status,
		Genre:      genre,
		Tag:        genre,
		Studio:     studio,
		UserRating: detail.VoteAverage,
		Actor:      actor,
		FanArt: FanArt{
			Thumb: []MovieThumb{
				{
					Preview: tmdb.ImageW500 + detail.BackdropPath,
				},
			},
		},
	}

	return utils.SaveNfo(nfoFile, top)
}
