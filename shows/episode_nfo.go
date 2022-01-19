package shows

import (
	"encoding/xml"
	"fengqi/kodi-metadata-tmdb-cli/tmdb"
	"fengqi/kodi-metadata-tmdb-cli/utils"
	"strconv"
)

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
	Tagline        string   `xml:"tagline"`
	Runtime        int      `xml:"runtime"`
	Thumb          Thumb    `xml:"thumb"`

	UniqueId UniqueId `xml:"uniqueid"`
	Year     string   `xml:"year"`
	TmdbId   string   `xml:"tmdbid"`

	MPaa      string   `xml:"mpaa"`
	Premiered string   `xml:"premiered"`
	Actor     []Actor  `xml:"actor"`
	Status    string   `xml:"status"`
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
	Season  int    `xml:"season"`
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
