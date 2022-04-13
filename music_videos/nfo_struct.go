package music_videos

import "encoding/xml"

type MusicVideoNfo struct {
	XMLName    xml.Name  `xml:"musicvideo"`
	Title      string    `xml:"title"`
	UserRating float32   `xml:"userrating"`
	Album      string    `xml:"-"`
	Plot       string    `xml:"-"`
	RunTime    int       `xml:"-"`
	Thumb      []Thumb   `xml:"thumb"`
	Poster     string    `xml:"poster"`
	PlayCount  int       `xml:"playcount"`
	LastPlayed string    `xml:"-"`
	Genre      []string  `xml:"genre"`
	Tag        []string  `xml:"tag"`
	Director   []string  `xml:"director"`
	Year       int       `xml:"-"`
	Studio     []string  `xml:"studio"`
	FileInfo   *FileInfo `xml:"fileinfo"`
	Actor      []Actor   `xml:"actor"`
	Artist     string    `xml:"-"`
	DateAdded  string    `xml:"dateadded"`
}

type Thumb struct {
	Aspect  string `xml:"aspect,attr"`
	Preview string `xml:"preview,attr"`
}

type FileInfo struct {
	StreamDetails StreamDetails `xml:"streamdetails"`
}

type StreamDetails struct {
	Video []Video `xml:"video"`
	Audio []Audio `xml:"audio"`
}

type Video struct {
	Codec             string `xml:"codec"`
	Aspect            string `xml:"aspect"`
	Width             int    `xml:"width"`
	Height            int    `xml:"height"`
	DurationInSeconds string `xml:"durationinseconds"`
	StereoMode        string `xml:"-"`
}

type Audio struct {
	Codec    string `xml:"codec"`
	Language string `xml:"language"`
	Channels int    `xml:"channels"`
}

type Actor struct {
	Name      string `xml:"name"`
	Role      string `xml:"role"`
	Order     int    `xml:"order"`
	SortOrder int    `xml:"sortorder"`
	Thumb     string `xml:"thumb"`
}
