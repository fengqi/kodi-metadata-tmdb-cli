package config

const (
	CollectorRunModeDaemon = 1 // 运行模式：守护进程模式
	CollectorRunModeOnce   = 2 // 运行模式：单次执行模式
	CollectorRunModeSpec   = 3 // 运行模式：临时指定目录模式
)

const (
	CollectorMoviesNfoModeMovieNfo = 1 // 电影 NFO 模式：movie.nfo
	CollectorMoviesNfoModeVideoNfo = 2 // 电影 NFO 模式：<VideoFileName>.nfo
)

const (
	LogModeStdout  = 1 // 日志输出模式：仅标准输出
	LogModeLogfile = 2 // 日志输出模式：仅日志文件
	LogModeBoth    = 3 // 日志输出模式：标准输出和日志文件
)

const (
	LogLevelDebug   = 0 // 日志等级：debug
	LogLevelInfo    = 1 // 日志等级：info
	LogLevelWarning = 2 // 日志等级：warning
	LogLevelError   = 3 // 日志等级：error
	LogLevelFatal   = 4 // 日志等级：fatal
)
