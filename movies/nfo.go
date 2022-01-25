package movies

import (
	"encoding/xml"
	"fengqi/kodi-metadata-tmdb-cli/tmdb"
	"fengqi/kodi-metadata-tmdb-cli/utils"
	"strconv"
)

// MovieNfo movie.nfo
//
// https://kodi.wiki/view/NFO_files/Movies
// NFO files to be scraped into the movie library are relatively simple and require only a single nfo file per title.
type MovieNfo struct {
	XMLName       xml.Name `xml:"movie"`
	Title         string   `xml:"title"`
	OriginalTitle string   `xml:"originaltitle"`
	SortTitle     string   `xml:"sorttitle"`
	Ratings       Ratings  `xml:"ratings"`
	UserRating    float32  `xml:"userrating"`
	Top250        string   `xml:"-"`
	Outline       string   `xml:"-"`
	Plot          string   `xml:"plot"`
	Tagline       string   `xml:"-"`
	Runtime       int      `xml:"-"`
	Thumb         []Thumb  `xml:"-"`
	FanArt        FanArt   `xml:"fanart"`
	MPaa          string   `xml:"-"`
	PlayCount     int      `xml:"-"`
	LastPlayed    string   `xml:"-"`
	Id            int      `xml:"id"`
	UniqueId      UniqueId `xml:"uniqueid"`
	Genre         []string `xml:"genre"`
	Tag           []string `xml:"tag"`
	Set           Set      `xml:"-"`
	Country       []string `xml:"country"`
	Credits       []string `xml:"credits"`
	Director      []string `xml:"director"`
	Premiered     string   `xml:"premiered"`
	Year          string   `xml:"-"`
	Status        string   `xml:"status"`
	Aired         string   `xml:"-"`
	Studio        []string `xml:"studio"`
	Trailer       string   `xml:"-"`
	FileInfo      FileInfo `xml:"-"`
	Actor         []Actor  `xml:"actor"`
	ShowLink      string   `xml:"-"`
	Resume        Resume   `xml:"-"`
	DateAdded     int      `xml:"-"`
}

type Set struct {
	Name     string `xml:"name"`
	Overview string `xml:"overview"`
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

type FileInfo struct {
	StreamDetails StreamDetails `xml:"streamdetails"`
}

type StreamDetails struct {
	Video    []Video    `xml:"video"`
	Audio    []Audio    `xml:"audio"`
	Subtitle []Subtitle `xml:"subtitle"`
}

type Video struct {
	Codec             string `xml:"codec"`
	Aspect            string `xml:"aspect"`
	Width             int    `xml:"width"`
	Height            int    `xml:"height"`
	DurationInSeconds int    `xml:"durationinseconds"`
	StereoMode        int    `xml:"stereomode"`
}

type Audio struct {
	Codec    string `xml:"codec"`
	Language string `xml:"language"`
	Channels int    `xml:"channels"`
}

type Subtitle struct {
	Codec    string `xml:"codec"`
	Micodec  string `xml:"micodec"`
	Language string `xml:"language"`
	ScanType string `xml:"scantype"`
	Default  bool   `xml:"default"`
	Forced   bool   `xml:"forced"`
}

type Thumb struct {
	Aspect  string `xml:"aspect,attr"`
	Preview string `xml:"preview,attr"`
}

type UniqueId struct {
	XMLName xml.Name `xml:"uniqueid"`
	Type    string   `xml:"type,attr"`
	Default bool     `xml:"default,attr"`
	Value   string   `xml:",chardata"`
}

type Actor struct {
	Name      string `xml:"name"`
	Role      string `xml:"role"`
	Order     int    `xml:"order"`
	SortOrder int    `xml:"sortorder"`
	Thumb     string `xml:"thumb"`
}

type FanArt struct {
	XMLName xml.Name     `xml:"fanart"`
	Thumb   []MovieThumb `xml:"thumb"`
}

type MovieThumb struct {
	Preview string `xml:"preview,attr"`
}

type Resume struct {
	Position string `xml:"position"`
	Total    int    `xml:"total"`
}

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
