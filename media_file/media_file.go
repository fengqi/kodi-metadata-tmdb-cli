package media_file

import (
	"fengqi/kodi-metadata-tmdb-cli/utils"
	"fmt"
	"path/filepath"
	"regexp"
	"strings"
)

var (
	trailerCompile, _ = regexp.Compile("(?i).*[\\[\\]\\(\\)_.-]+trailer[\\[\\]\\(\\)_.-]?(\\d)*$")
	sampleCompile, _  = regexp.Compile("(?i).*[\\[\\]\\(\\)_.-]+sample[\\[\\]\\(\\)_.-]?$")
	dvdCompile, _     = regexp.Compile("(video_ts|vts_\\d\\d_\\d)\\.(vob|bup|ifo)")
	blurayCompile, _  = regexp.Compile("(index\\.bdmv|movieobject\\.bdmv|\\d{5}\\.m2ts|\\d{5}\\.clpi|\\d{5}\\.mpls)")
)

func NewMediaFile(path, filename string) *mediaFile {
	if filename[0:1] == "." {
		return nil
	}

	mediaType := parseMediaType(path, filename)

	return &mediaFile{
		Path:     path,
		Filename: filename,
		Type:     mediaType,
	}
}

func parseMediaType(path, filename string) MediaType {
	folderName := strings.ToLower(filepath.Base(path))
	ext := filepath.Ext(filename)
	basename := strings.ToLower(strings.Replace(filename, ext, "", 1))

	fmt.Println(filename, folderName, basename)

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

func isDiscFile(filename, path string) bool {
	return isDVDFile(filename, path) || isBlurayFile(filename, path) || isHDDVDFile(filename, path)
}

func isDVDFile(filename, path string) bool {
	if filename == "video_ts" || utils.EndsWith(path, "video_ts") {
		return true
	}

	return dvdCompile.FindString(filename) != ""
}

func isBlurayFile(filename, path string) bool {
	if filename == "bdmv" || utils.EndsWith(path, "bdmv") {
		return true
	}

	return blurayCompile.FindString(filename) != ""
}

func isHDDVDFile(filename, path string) bool {
	return filename == "hvdvd_ts" || utils.EndsWith(path, "hvdvd_ts")
}
