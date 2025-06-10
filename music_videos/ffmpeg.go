package music_videos

import (
	"encoding/json"
	"fengqi/kodi-metadata-tmdb-cli/ffmpeg"
	"fengqi/kodi-metadata-tmdb-cli/utils"
	"github.com/fengqi/lrace"
	"os"
	"strconv"
	"strings"
)

// 获取视频信息
func (mv *MusicVideo) getProbe() (*ffmpeg.ProbeData, error) {
	var probe = new(ffmpeg.ProbeData)
	fileMd5 := mv.nameMd5()
	cacheFile := mv.cacheDir() + "/" + fileMd5 + ".json"
	if _, err := os.Stat(cacheFile); err == nil {
		utils.Logger.DebugF("get video probe from cache: %s", cacheFile)
		if bytes, err := os.ReadFile(cacheFile); err == nil {
			if err = json.Unmarshal(bytes, probe); err == nil {
				return probe, nil
			}
		}
	}

	// 保存缓存
	probe, err := ffmpeg.Probe(mv.MediaFile.Path)
	if err == nil {
		utils.Logger.DebugF("save video probe to cache: %s", cacheFile)
		f, err := os.OpenFile(cacheFile, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
		if err == nil {
			bytes, _ := json.MarshalIndent(probe, "", "    ")
			_, err = f.Write(bytes)
			_ = f.Close()
		}
	}

	return probe, err
}

// 缩略图提取
func (mv *MusicVideo) drawThumb() error {
	thumb := mv.thumbPath()

	// 如果有视频文件同名后缀的图片，尝试直接使用
	filenameWithoutSuffix := strings.TrimRight(mv.MediaFile.Path, mv.MediaFile.Suffix)
	for _, i := range ThumbImagesFormat {
		check := filenameWithoutSuffix + "." + i
		if lrace.FileExist(check) {
			n, err := lrace.CopyFile(check, thumb)
			if n > 0 && err == nil {
				return nil
			}
		}
	}

	// 对于大文件，尝试偏移30秒，防止读到的是黑屏白屏或者logo
	ss := "00:00:00"
	second, _ := strconv.ParseFloat(mv.VideoStream.Duration, 10)
	if mv.VideoStream != nil && second > 30 {
		ss = "00:00:30"
	}

	// 特殊处理
	filename := mv.MediaFile.Filename
	if (len(filename) > 2 && filename[0:2] == "03") || (len(filename) > 5 && strings.ToLower(filename[0:5]) == "heyzo") {
		ss = "00:01:10"
	}

	utils.Logger.InfoF("draw thumb start: %s to %s", ss, thumb)

	return ffmpeg.Frame(mv.MediaFile.Path, thumb, "-ss", ss)
}
