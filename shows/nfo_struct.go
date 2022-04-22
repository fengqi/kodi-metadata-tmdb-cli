package shows

import "encoding/xml"

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
	Premiered      string        `xml:"premiered"`
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

type TvEpisodeNfo struct {
	XMLName        xml.Name `xml:"episodedetails"`
	Title          string   `xml:"title"`
	OriginalTitle  string   `xml:"originaltitle"`
	ShowTitle      string   `xml:"showtitle"`
	Ratings        []Rating `xml:"ratings"`
	UserRating     float32  `xml:"userrating"`
	Top250         string   `xml:"top250"`
	Season         int      `xml:"season"`
	Episode        int      `xml:"episode"`
	DisplayEpisode int      `xml:"displayepisode"`
	DisplaySeason  int      `xml:"displayseason"`
	Outline        string   `xml:"outline"`
	Plot           string   `xml:"plot"`
	Tagline        string   `xml:"-"`
	Runtime        int      `xml:"runtime"`
	Thumb          Thumb    `xml:"thumb"`

	UniqueId UniqueId `xml:"uniqueid"`
	Year     string   `xml:"year"`
	TmdbId   string   `xml:"tmdbid"`

	MPaa      string   `xml:"-"`
	Premiered string   `xml:"premiered"`
	Actor     []Actor  `xml:"actor"`
	Status    string   `xml:"-"`
	Aired     string   `xml:"aired"`
	Genre     []string `xml:"genre"`
	Studio    []string `xml:"studio"`

	FileInfo FileInfo `xml:"fileinfo"`
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
	Number int    `xml:"number,attr"`
	Value  string `xml:",chardata"`
}

type Resume struct {
	Position string `xml:"position"`
	Total    int    `xml:"total"`
}
