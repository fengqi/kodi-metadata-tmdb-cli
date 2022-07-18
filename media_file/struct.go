package media_file

type mediaFile struct {
	Path     string
	Filename string
	Type     MediaType
}

type MediaType int

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

func (mf *mediaFile) IsNFO() bool {
	return mf.Type == NFO
}

func (mf *mediaFile) IsVideo() bool {
	return mf.Type == VIDEO
}
