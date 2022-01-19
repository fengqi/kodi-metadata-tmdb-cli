package shows

import (
	"encoding/xml"
	"fengqi/kodi-metadata-tmdb-cli/tmdb"
	"fengqi/kodi-metadata-tmdb-cli/utils"
	"strconv"
)

// TvShowNfo tvshow.nfo
//
// https://kodi.wiki/view/NFO_files/TV_shows
// NFO files for TV Shows are a little bit more complex as they require the following NFO files:
//
// One nfo file for the TV Show. This file holds the overall TV show information
// One nfo file for each Episode. This file holds information specific to that episode
// For one TV Show with 10 episodes, 11 nfo files are required.
type TvShowNfo struct {
	XMLName        xml.Name      `xml:"tvshow"`
	Title          string        `xml:"title"`
	OriginalTitle  string        `xml:"originaltitle"`
	ShowTitle      string        `xml:"showtitle"` // 	Not in common use, but some skins may display an alternate title
	SortTitle      string        `xml:"sorttitle"`
	Ratings        Ratings       `xml:"ratings"`
	UserRating     float32       `xml:"userrating"`
	Top250         string        `xml:"-"`
	Season         int           `xml:"season"`
	Episode        int           `xml:"episode"`
	DisplayEpisode int           `xml:"-"`
	DisplaySeason  int           `xml:"-"`
	Outline        string        `xml:"-"`
	Plot           string        `xml:"plot"`
	Tagline        string        `xml:"-"`
	Runtime        int           `xml:"-"`
	Thumb          []Thumb       `xml:"-"`
	FanArt         FanArt        `xml:"fanart"`
	MPaa           string        `xml:"mpaa"`
	PlayCount      int           `xml:"-"`
	LastPlayed     string        `xml:"-"`
	EpisodeGuide   EpisodeGuide  `xml:"-"`
	Id             int           `xml:"id"`
	UniqueId       UniqueId      `xml:"uniqueid"`
	Genre          []string      `xml:"genre"`
	Tag            []string      `xml:"tag"`
	Premiered      string        `json:"premiered"`
	Year           string        `xml:"-"`
	Status         string        `xml:"status"`
	Aired          string        `xml:"-"`
	Studio         []string      `xml:"studio"`
	Trailer        string        `xml:"trailer"`
	Actor          []Actor       `xml:"actor"`
	NamedSeason    []NamedSeason `xml:"namedseason"`
	Resume         Resume        `xml:"-"`
	DateAdded      int           `xml:"-"`
}

type UniqueId struct {
	XMLName xml.Name `xml:"uniqueid"`
	Type    string   `xml:"type,attr"`
	Default bool     `xml:"default,attr"`
}

type Actor struct {
	Name      string `xml:"name"`
	Role      string `xml:"role"`
	Order     int    `xml:"order"`
	SortOrder int    `xml:"sortorder"`
	Thumb     string `xml:"thumb"`
}

type FanArt struct {
	XMLName xml.Name    `xml:"fanart"`
	Thumb   []ShowThumb `xml:"thumb"`
}

type ShowThumb struct {
	Preview string `xml:"preview,attr"`
}

type EpisodeGuide struct {
	Url Url `xml:"url"`
}

type Url struct {
	Cache string `xml:"cache"`
}

type NamedSeason struct {
	Number string `xml:"number"`
}

type Resume struct {
	Position string `xml:"position"`
	Total    int    `xml:"total"`
}

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
		Id:         detail.Id,
		Premiered:  detail.FirstAirDate,
		Ratings:    Ratings{Rating: rating},
		MPaa:       "TV-14",
		Status:     detail.Status,
		Genre:      genre,
		Studio:     studio,
		Season:     detail.NumberOfSeasons,
		Episode:    detail.NumberOfEpisodes,
		UserRating: detail.VoteAverage,
		Actor:      actor,
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
