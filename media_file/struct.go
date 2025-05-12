package media_file

import "strings"

type MediaFile struct {
	Path      string    // 完整路径
	Filename  string    // 文件名
	Suffix    string    // 后缀
	MediaType MediaType // 文件类型
	VideoType VideoType // 视频类型
}

type VideoType int

type MediaType int

const (
	Movies = iota + 1
	TvShows
	MusicVideo
)

const (
	UNKNOWN MediaType = iota
	VIDEO
	TRAILER
	SAMPLE
	AUDIO
	SUBTITLE
	NFO
	POSTER
	FANART
	BANNER
	CLEARART
	DISC
	LOGO
	CLEARLOGO
	THUMB
	CHARACTERART
	KEYART
	SEASON_POSTER
	SEASON_FANART
	SEASON_BANNER
	SEASON_THUMB
	EXTRAFANART
	EXTRATHUMB
	EXTRA
	GRAPHIC
	MEDIAINFO
	VSMETA
	THEME
	TEXT
	DOUBLE_EXT
)

var (
	VIDEO_TS = "VIDEO_TS"
	BDMV     = "BDMV"
	HVDVD_TS = "HVDVD_TS"

	ArtworkFileTypes = []string{
		"jpg", "jpeg,", "png", "tbn", "gif", "bmp", "webp",
	}
	VideoFileTypes = []string{
		".3gp", ".asf", ".asx", ".avc", ".avi", ".bdmv", ".bin", ".bivx", ".braw", ".dat", ".divx", ".dv", ".dvr-ms",
		".disc", ".evo", ".fli", ".flv", ".h264", ".ifo", ".img", ".iso", ".mts", ".mt2s", ".m2ts", ".m2v", ".m4v",
		".mkv", ".mk3d", ".mov", ".mp4", ".mpeg", ".mpg", ".nrg", ".nsv", ".nuv", ".ogm", ".pva", ".qt", ".rm", ".rmvb",
		".strm", ".svq3", ".ts", ".ty", ".viv", ".vob", ".vp3", ".wmv", ".webm", ".xvid",
	}
	AudioFileTypes = []string{
		".a52", ".aa3", ".aac", ".ac3", ".adt", ".adts", ".aif", ".aiff", ".alac", ".ape", ".at3", ".atrac", ".au",
		".dts", ".flac", ".m4a", ".m4b", ".m4p", ".mid", ".midi", ".mka", ".mp3", ".mpa", ".mlp", ".oga", ".ogg",
		".pcm", ".ra", ".ram", ".tta", ".thd", ".wav", ".wave", ".wma",
	}
	SubtitleFileTypes = []string{
		".aqt", ".cvd", ".dks", ".jss", ".sub", ".sup", ".ttxt", ".mpl", ".pjs", ".psb", ".rt", ".srt", ".smi",
		".ssf", ".ssa", ".svcd", ".usf", ".ass", ".pgs", ".vobsub",
	}
)

// IsNFO 是否是NFO文件
func (mf *MediaFile) IsNFO() bool {
	return mf.MediaType == NFO
}

// IsVideo 是否是视频
func (mf *MediaFile) IsVideo() bool {
	return mf.MediaType == VIDEO
}

// IsBluRay 是否是蓝光目录
func (mf *MediaFile) IsBluRay() bool {
	return mf.MediaType == DISC && (mf.Filename == BDMV || mf.Filename == HVDVD_TS)
}

// IsDvd 是否是DVD目录
func (mf *MediaFile) IsDvd() bool {
	return mf.IsDisc() && (mf.Filename == VIDEO_TS || mf.Filename == "DVD")
}

// IsDisc 判断是否是光盘目录
func (mf *MediaFile) IsDisc() bool {
	return mf.MediaType == DISC
}

// PathWithoutSuffix 完整路径，去掉后缀，用于生成NFO、海报等
func (mf *MediaFile) PathWithoutSuffix() string {
	return strings.Replace(mf.Path, mf.Suffix, "", 1)
}
