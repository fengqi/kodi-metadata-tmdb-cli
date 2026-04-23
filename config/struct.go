package config

type Config struct {
	Log       *LogConfig       `json:"log"`       // 日志配置
	Ffmpeg    *FfmpegConfig    `json:"ffmpeg"`    // ffmpeg 配置
	Tmdb      *TmdbConfig      `json:"tmdb"`      // TMDB 配置
	Kodi      *KodiConfig      `json:"kodi"`      // kodi 配置
	Collector *CollectorConfig `json:"collector"` // 刮削配置
	Ai        *AiConfig        `json:"ai"`        // AI 配置
}

type KodiConfig struct {
	Enable       bool   `json:"enable"`        // 是否开启 kodi 通知
	CleanLibrary bool   `json:"clean_library"` // 是否清理媒体库
	JsonRpc      string `json:"json_rpc"`      // kodi rpc 路径
	Timeout      int    `json:"timeout"`       // 连接 kodi 超时时间
	Username     string `json:"username"`      // rpc 用户名
	Password     string `json:"password"`      // rpc 密码
}

type LogConfig struct {
	Mode  int    `json:"mode"`  // 日志模式：1 stdout，2 logfile，3 both
	Level int    `json:"level"` // 日志等级：0 debug，1 info，2 warning，3 error，4 fatal
	File  string `json:"file"`  // 日志文件路径
}

type FfmpegConfig struct {
	MaxWorker   int    `json:"max_worker"`   // 最大进程数，建议为逻辑 CPU 个数
	FfmpegPath  string `json:"ffmpeg_path"`  // ffmpeg 可执行文件路径
	FfprobePath string `json:"ffprobe_path"` // ffprobe 可执行文件路径
}

type TmdbConfig struct {
	ApiHost   string `json:"api_host"`   // TMDB 接口地址
	ApiKey    string `json:"api_key"`    // api key
	ImageHost string `json:"image_host"` // 图片地址
	Language  string `json:"language"`   // 语言
	Rating    string `json:"rating"`     // 内容分级
	Proxy     string `json:"proxy"`      // 请求 TMDB 代理，支持 http、https、socks5、socks5h
}

type CollectorConfig struct {
	RunMode        int      `json:"run_mode"`         // 运行模式：1 daemon，2 once，3 spec
	Watcher        bool     `json:"watcher"`          // 是否开启文件监听
	CronSeconds    int      `json:"cron_seconds"`     // 定时扫描频率
	CronScan       bool     `json:"cron_scan"`        // 是否开启定时扫描
	CronScanKodi   bool     `json:"cron_scan_kodi"`   // 定时扫描后触发 kodi 扫描
	TmpSuffix      []string `json:"tmp_suffix"`       // 临时文件后缀列表
	NfoField       NfoField `json:"nfo_field"`        // NFO 字段
	SkipFolders    []string `json:"skip_folders"`     // 跳过目录，可多个
	SkipKeywords   []string `json:"skip_keywords"`    // 跳过文件名中的关键字，可多个
	MoviesNfoMode  int      `json:"movies_nfo_mode"`  // 电影 NFO 写入模式：1 movie.nfo，2 <VideoFileName>.nfo（当前版本未启用）
	MoviesDir      []string `json:"movies_dir"`       // 电影文件根目录，可多个
	ShowsDir       []string `json:"shows_dir"`        // 电视剧文件根目录，可多个
	MusicVideosDir []string `json:"music_videos_dir"` // 音乐视频文件根目录，可多个
}

type NfoField struct {
	Tag   bool `json:"tag"`   // 开启标签
	Genre bool `json:"genre"` // 开启分类
}

type AiConfig struct {
	Enable              bool    `json:"enable"`               // 是否开启 AI
	BaseURL             string  `json:"base_url"`             // OpenAI 兼容接口地址
	ApiKey              string  `json:"api_key"`              // API Key
	Model               string  `json:"model"`                // 模型名称
	Temperature         float64 `json:"temperature"`          // 采样温度
	TimeoutSeconds      int     `json:"timeout_seconds"`      // 超时时间（秒）
	ConfidenceThreshold float64 `json:"confidence_threshold"` // 最低可接受置信度
	MatchMode           int     `json:"match_mode"`           // 匹配模式：1 规则优先，2 AI 优先，3 规则后 AI 覆盖
	SearchMode          int     `json:"search_mode"`          // 搜索模式：1 首个结果，2 算法选择，3 AI 决策
}
