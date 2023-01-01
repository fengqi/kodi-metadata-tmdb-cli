package music_videos

import (
	"crypto/md5"
	"encoding/json"
	"fengqi/kodi-metadata-tmdb-cli/ffmpeg"
	"fengqi/kodi-metadata-tmdb-cli/utils"
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

func (m *MusicVideo) getFullPath() string {
	return m.Dir + "/" + m.OriginTitle
}

func (m *MusicVideo) getNfoThumb() string {
	return m.Dir + "/" + m.Title + "-thumb.jpg"
}

func (m *MusicVideo) getNfoFile() string {
	return m.Dir + "/" + m.Title + ".nfo"
}

func (m *MusicVideo) NfoExist() bool {
	nfo := m.getNfoFile()

	if info, err := os.Stat(nfo); err == nil && info.Size() > 0 {
		return true
	}

	return false
}

func (m *MusicVideo) ThumbExist() bool {
	thumb := m.getNfoThumb()
	if info, err := os.Stat(thumb); err == nil && info.Size() > 0 {
		return true
	}

	return false
}

func (m *MusicVideo) getProbe() (*ffmpeg.ProbeData, error) {
	// 读取缓存
	var probe = new(ffmpeg.ProbeData)

	fileMd5 := m.GetNameMd5()
	cacheFile := m.BaseDir + "/tmdb/" + fileMd5 + ".json"
	if _, err := os.Stat(cacheFile); err == nil {
		utils.Logger.DebugF("get video probe from cache: %s", cacheFile)
		if bytes, err := ioutil.ReadFile(cacheFile); err == nil {
			if err = json.Unmarshal(bytes, probe); err == nil {
				return probe, nil
			}
		}
	}

	// 保存缓存
	probe, err := ffmpeg.Probe(m.Dir + "/" + m.OriginTitle)
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

func (m *MusicVideo) GetNameMd5() string {
	h := md5.New()
	_, _ = io.WriteString(h, m.getFullPath())
	sum := fmt.Sprintf("%x", h.Sum(nil))
	return sum
}
