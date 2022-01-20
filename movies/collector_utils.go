package movies

import (
	"fengqi/kodi-metadata-tmdb-cli/tmdb"
	"fengqi/kodi-metadata-tmdb-cli/utils"
	"io/fs"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

// 解析目录, 返回详情
// TODO 跳过电视剧，放错目录了
func parseMoviesDir(baseDir string, file fs.FileInfo) *Movie {
	suffix := utils.IsVideo(file.Name())
	if !file.IsDir() && suffix == "" {
		return nil
	}

	// 使用目录或者没有后缀的文件名
	movieName := file.Name()
	if suffix != "" {
		movieName = strings.Replace(movieName, "."+suffix, "", 1)
	}

	// 用点号.或者空格分割，这里假定文件或者目录命名时规范的
	formatName := strings.Replace(movieName, " ", ".", -1)
	split := strings.Split(formatName, ".")
	if split == nil || len(split) < 3 {
		utils.Logger.WarningF("file name: %s syntax err, skipped", file.Name())
		return nil
	}

	movieDir := &Movie{Dir: baseDir, OriginTitle: movieName, IsFile: !file.IsDir(), Suffix: suffix}

	// 文件名识别
	nameStart := false
	nameStop := false
	for _, item := range split {
		if year := utils.IsYear(item); year > 0 {
			movieDir.Year = year
			nameStop = true
			continue
		}

		if format := utils.IsFormat(item); len(format) > 0 {
			movieDir.Format = format
			nameStop = true
			continue
		}

		if source := utils.IsSource(item); len(source) > 0 {
			movieDir.Source = source
			nameStop = true
			continue
		}

		if studio := utils.IsStudio(item); len(studio) > 0 {
			movieDir.Studio = studio
			nameStop = true
			continue
		}

		if !nameStart {
			nameStart = true
			nameStop = false
		}

		if !nameStop {
			movieDir.Title += item + " "
		}
	}
	movieDir.Title = utils.CleanTitle(movieDir.Title)
	if len(movieDir.Title) == 0 {
		utils.Logger.WarningF("file: %s parse title empty: %v", file.Name(), movieDir)
		return nil
	}

	// 通过文件指定id
	// todo all use baseDir + "/tmdb/"
	idFile := baseDir + "/" + file.Name() + "/tmdb/id.txt"
	if !file.IsDir() {
		idFile = baseDir + "/tmdb/" + movieName + ".id.txt"
	}
	if _, err := os.Stat(idFile); err == nil {
		bytes, err := ioutil.ReadFile(idFile)
		if err == nil {
			movieDir.MovieId, _ = strconv.Atoi(strings.Trim(string(bytes), "\r\n "))
		} else {
			utils.Logger.WarningF("read movies id specially file: %s err: %v", idFile, err)
		}
	}

	//识别是否时蓝光或dvd目录
	if file.IsDir() {
		fileInfo, err := ioutil.ReadDir(baseDir + "/" + file.Name())
		if err == nil {
			audioTs := false
			videoTs := false
			for _, item := range fileInfo {
				if item.IsDir() && item.Name() == "BDMV" || item.Name() == "CERTIFICATE" {
					movieDir.IsBluray = true
					break
				}

				if item.IsDir() && item.Name() == "AUDIO_TS" {
					audioTs = true
				}
				if item.IsDir() && item.Name() == "VIDEO_TS" {
					videoTs = true
				}
				if videoTs && audioTs {
					movieDir.IsDvd = true
					break
				}

				if suffix := utils.IsVideo(item.Name()); suffix != "" {
					movieDir.IsSingleFile = true
					movieDir.VideoFileName = item.Name()
					break
				}
			}
		}
	}

	return movieDir
}

// tmdb 缓存目录
// TODO 统一使用一个目录
func (d *Movie) checkCacheDir() {
	dir := d.GetCacheDir()
	if _, err := os.Stat(dir); err != nil && os.IsNotExist(err) {
		err := os.Mkdir(dir, 0755)
		if err != nil {
			utils.Logger.ErrorF("create cache: %s dir err: %v", dir, err)
		}
	}
}

func (d *Movie) GetCacheDir() string {
	if d.IsFile {
		return d.Dir + "/tmdb"
	}
	return d.GetFullDir() + "/tmdb"
}

func (d *Movie) GetFullDir() string {
	return d.Dir + "/" + d.OriginTitle
}

func (d *Movie) downloadImage(detail *tmdb.MovieDetail) error {
	utils.Logger.DebugF("download %s images", d.Title)

	var err error
	if len(detail.PosterPath) > 0 {
		posterFile := d.GetFullDir() + "/poster.jpg"
		if d.IsFile {
			posterFile = d.GetFullDir() + "-poster.jpg"
		}
		err = utils.DownloadFile(tmdb.ImageOriginal+detail.PosterPath, posterFile)
	}

	if len(detail.BackdropPath) > 0 {
		fanArtFile := d.GetFullDir() + "/fanart.jpg"
		if d.IsFile {
			fanArtFile = d.GetFullDir() + "-fanart.jpg"
		}
		err = utils.DownloadFile(tmdb.ImageOriginal+detail.BackdropPath, fanArtFile)
	}

	return err
}

// maybe <VideoFileName>.nfo
func (m *Movie) getNfoFile() string {
	if m.IsFile {
		return m.GetFullDir() + ".nfo"
	}

	if m.IsBluray {
		return m.GetFullDir() + "/BDMV/index.nfo"
	}

	if m.IsDvd {
		return m.GetFullDir() + "/VIDEO_TS/VIDEO_TS.nfo"
	}

	suffix := utils.IsVideo(m.VideoFileName)
	return m.GetFullDir() + "/" + strings.Replace(m.VideoFileName, "."+suffix, "", 1) + ".nfo"
}
