package movies

import (
	"fengqi/kodi-metadata-tmdb-cli/config"
	"fengqi/kodi-metadata-tmdb-cli/tmdb"
	"fengqi/kodi-metadata-tmdb-cli/utils"
	"os"
	"strconv"
	"strings"
)

// 保存NFO文件
func (m *Movie) saveToNfo(detail *tmdb.MovieDetail, mode int) error {
	nfoFile := m.getNfoFile(mode)
	if nfoFile == "" {
		utils.Logger.InfoF("movie nfo empty %v", m)
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

	return utils.SaveNfo(nfoFile, top)
}

// maybe <VideoFileName>.nfo
// Kodi比较推荐 <VideoFileName>.nfo 但是存在一种情况就是，使用inotify监听文件变动，可能电影目录先创建
// 里面的视频文件会迟一点，这个时候 VideoFileName 就会为空，导致NFO写入失败
// 部分资源可能存在 <VideoFileName>.nfo 且里面写入了一些无用的信息，会产生冲突
// 如果使用 movie.nfo 就不需要考虑这个情况，但是需要打开媒体源的 "电影在以片名命名的单独目录中"
func (m *Movie) getNfoFile(mode int) string {
	if m.MediaFile.IsBluRay() {
		//if utils.FileExist(m.GetFullDir() + "/MovieObject.bdmv") {
		//	return m.GetFullDir() + "/MovieObject.nfo"
		//}
		return m.GetFullDir() + "/index.nfo"
	}

	if m.MediaFile.IsDvd() {
		return m.GetFullDir() + "/VIDEO_TS/VIDEO_TS.nfo"
	}

	//if mode == 2 { // todo movie.nfo作为可选，<VideoFileName>.nfo为固定生成
	//	return m.GetFullDir() + "/movie.nfo"
	//}

	suffix := utils.IsVideo(m.MediaFile.Filename)
	return strings.Replace(m.MediaFile.Path, "."+suffix, "", 1) + ".nfo"
}

// NfoExist 判断NFO文件是否存在
func (m *Movie) NfoExist(mode int) bool {
	nfo := m.getNfoFile(mode)

	if info, err := os.Stat(nfo); err == nil && info.Size() > 0 {
		return true
	}

	return false
}
