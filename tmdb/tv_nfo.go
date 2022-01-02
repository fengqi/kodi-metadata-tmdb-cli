package tmdb

import (
	"encoding/xml"
	"fengqi/kodi-metadata-tmdb-cli/utils"
	"strconv"
)

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
	Year           string        `xml:"year"`
	Status         string        `xml:"status"`
	Aired          string        `xml:"-"`
	Studio         []string      `xml:"studio"`
	Trailer        string        `xml:"trailer"`
	Actor          []Actor       `xml:"actor"`
	NamedSeason    []NamedSeason `xml:"namedseason"`
	Resume         Resume        `xml:"-"`
	DateAdded      int           `xml:"-"`
}

type Ratings struct {
	Rating []Rating `xml:"rating"`
}

type Rating struct {
	Name  string  `xml:"name,attr"`
	Max   int     `xml:"max,attr"`
	Value float32 `xml:"value"`
	Votes int     `xml:"votes"`
}

type Thumb struct {
	Aspect  string `xml:"aspect,attr"`
	Preview string `xml:"preview,attr"`
	Season  int    `xml:"season"`
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

// SaveToNfo 保存剧集汇总信息到nfo文件
func (d *TvDetail) SaveToNfo(nfo string) error {
	utils.Logger.InfoF("save tvshow.nfo to %s", nfo)

	genre := make([]string, 0)
	for _, item := range d.Genres {
		genre = append(genre, item.Name)
	}

	studio := make([]string, 0)
	for _, item := range d.Networks {
		studio = append(studio, item.Name)
	}

	rating := make([]Rating, 1)
	rating[0] = Rating{
		Name:  "tmdb",
		Max:   10,
		Value: d.VoteAverage,
		Votes: d.VoteCount,
	}

	actor := make([]Actor, 0)
	if d.AggregateCredits != nil {
		for _, item := range d.AggregateCredits.Cast {
			actor = append(actor, Actor{
				Name:  item.Name,
				Role:  item.Roles[0].Character,
				Order: item.Order,
				Thumb: ImageW500 + item.ProfilePath,
			})
		}
	}

	year := ""
	if len(d.FirstAirDate) > 3 {
		year = d.FirstAirDate[0:4]
	}

	top := &TvShowNfo{
		Title:         d.Name,
		OriginalTitle: d.OriginalName,
		ShowTitle:     d.Name,
		SortTitle:     d.Name,
		Plot:          d.Overview,
		UniqueId: UniqueId{
			Type:    strconv.Itoa(d.Id),
			Default: true,
		},
		Id:         d.Id,
		Year:       year,
		Ratings:    Ratings{Rating: rating},
		MPaa:       "TV-14",
		Status:     d.Status,
		Genre:      genre,
		Studio:     studio,
		Season:     d.NumberOfSeasons,
		Episode:    d.NumberOfEpisodes,
		UserRating: d.VoteAverage,
		Actor:      actor,
		FanArt: FanArt{
			Thumb: []ShowThumb{
				{
					Preview: ImageOriginal + d.BackdropPath,
				},
			},
		},
	}

	return utils.SaveNfo(nfo, top)
}
