package movies

import "encoding/xml"

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
	Runtime       int      `xml:"-"` // Kodi会从视频文件提取，考虑到版本（剪辑版、完整版等）问题，这里不提供
	Thumb         []Thumb  `xml:"-"`
	FanArt        *FanArt  `xml:"fanart"`
	MPaa          string   `xml:"mpaa"`
	PlayCount     int      `xml:"-"`
	LastPlayed    string   `xml:"-"`
	Id            int      `xml:"id"`
	UniqueId      UniqueId `xml:"uniqueid"`
	Genre         []string `xml:"genre"`
	Tag           []string `xml:"tag"`
	Set           Set      `xml:"-"`
	Country       []string `xml:"country"`
	Languages     []string `xml:"languages"`
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
	Default bool     `xml:"default,attr"`
	Type    string   `xml:"type,attr"`
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
