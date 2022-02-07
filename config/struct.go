package config

type Config struct {
	LogLevel int    `json:"log_level"`
	LogFile  string `json:"log_file"`

	Rating   string `json:"rating"`   // 内容分级
	ApiKey   string `json:"api_key"`  // api key
	Language string `json:"language"` // 语言

	MoviesNfoMode int      `json:"movies_nfo_mode"` // 电影NFO写入模式：1 movie.nfo， 2 <VideoFileName>.nfo
	ShowsDir      []string `json:"shows_dir"`
	MoviesDir     []string `json:"movies_dir"`
	CronSeconds   int      `json:"cron_seconds"` // todo、shows、movies 分别设置

	Kodi    KodiConfig `json:"kodi"`
	Exclude []string   `json:"exclude`
}

type KodiConfig struct {
	JsonRpc  string `json:"json_rpc"`
	Timeout  int    `json:"timeout"`
	Username string `json:"username"`
	Password string `json:"password"`
}
