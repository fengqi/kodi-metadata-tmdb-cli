package media_file

import (
	"fengqi/kodi-metadata-tmdb-cli/utils"
	"path/filepath"
	"regexp"
	"strings"
)

var (
	trailerCompile, _ = regexp.Compile("(?i).*[\\[\\]\\(\\)_.-]+trailer[\\[\\]\\(\\)_.-]?(\\d)*$")
	sampleCompile, _  = regexp.Compile("(?i).*[\\[\\]\\(\\)_.-]+sample[\\[\\]\\(\\)_.-]?$")
	dvdCompile, _     = regexp.Compile("(video_ts|vts_\\d\\d_\\d)\\.(vob|bup|ifo)")
	bluRayCompile, _  = regexp.Compile("(index\\.bdmv|movieobject\\.bdmv|\\d{5}\\.m2ts|\\d{5}\\.clpi|\\d{5}\\.mpls)")
)

// NewMediaFile 实例化媒体类型
func NewMediaFile(path, filename string, videoType VideoType) *MediaFile {
	if filename[0:1] == "." {
		return nil
	}

	return &MediaFile{
		Path:      path,
		Dir:       filepath.Dir(path),
		Filename:  filename,
		MediaType: parseMediaType(path, filename),
		VideoType: videoType,
		Suffix:    filepath.Ext(filename),
	}
}

// MediaType 解析文件类型
func parseMediaType(path, filename string) MediaType {
	folderName := strings.ToLower(filepath.Base(path))
	ext := filepath.Ext(filename)
	basename := strings.ToLower(strings.Replace(filename, ext, "", 1))

	if folderName == "extras" || folderName == "extra" {
		return EXTRA
	}

	if ext == ".nfo" {
		return NFO
	}

	if ext == ".vsmeta" {
		return VSMETA
	}

	// 图片
	for _, v := range ArtworkFileTypes { // todo map
		if strings.Contains(filename, v) {
			return GRAPHIC
		}
	}

	for _, v := range AudioFileTypes {
		if strings.Contains(filename, v) {
			return AUDIO
		}
	}

	for _, v := range SubtitleFileTypes {
		if strings.Contains(filename, v) {
			return SUBTITLE
		}
	}

	if isDiscFile(filename, path) {
		return DISC
	}

	for _, v := range VideoFileTypes {
		if strings.Contains(filename, v) {
			if basename == "movie-trailer" ||
				folderName == "trailer" ||
				folderName == "trailers" ||
				trailerCompile.FindString(basename) != "" {
				return TRAILER
			}

			if basename == "sample" || folderName == "sample" || sampleCompile.FindString(basename) != "" {
				return SAMPLE
			}

			return VIDEO
		}
	}

	if isDiscFile(filename, folderName) {
		return VIDEO
	}

	if ext == ".txt" {
		return VSMETA
	}

	return UNKNOWN
}

// 是否是光盘文件
func isDiscFile(filename, path string) bool {
	return isDVDFile(filename, path) || isBluRayFile(filename, path) || isHDDVDFile(filename, path)
}

// 是否是DVD光盘文件
func isDVDFile(filename, path string) bool {
	if filename == "VIDEO_TS" || utils.EndsWith(path, "VIDEO_TS") {
		return true
	}

	return dvdCompile.FindString(filename) != ""
}

// 是否是蓝光文件
func isBluRayFile(filename, path string) bool {
	if filename == "BDMV" || utils.EndsWith(path, "BDMV") {
		return true
	}

	return bluRayCompile.FindString(filename) != ""
}

// 是否是HD DVD文件
func isHDDVDFile(filename, path string) bool {
	return filename == "HVDVD_TS" || utils.EndsWith(path, "HVDVD_TS")
}
