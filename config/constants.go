package config

// 运行模式
const (
	CollectorRunModeDaemon = 1 // 守护进程模式
	CollectorRunModeOnce   = 2 // 单次执行模式
	CollectorRunModeSpec   = 3 // 临时指定目录模式
)

// 电影 NFO 模式
const (
	CollectorMoviesNfoModeMovieNfo = 1 // movie.nfo
	CollectorMoviesNfoModeVideoNfo = 2 // <VideoFileName>.nfo
)

// 日志输出模式
const (
	LogModeStdout  = 1 // 仅标准输出
	LogModeLogfile = 2 // 仅日志文件
	LogModeBoth    = 3 // 标准输出和日志文件
)

// 日志等级
const (
	LogLevelDebug   = 0 // debug
	LogLevelInfo    = 1 // info
	LogLevelWarning = 2 // warning
	LogLevelError   = 3 // error
	LogLevelFatal   = 4 // fatal
)

// AI 匹配模式
const (
	AiMatchModeRuleThenAi         = 1 // 规则优先，匹配不到再 AI 介入
	AiMatchModeAiThenRule         = 2 // AI 优先，匹配不到再规则介入
	AiMatchModeRuleWithAiOverride = 3 // 规则优先，结果给 AI 参考，最终使用 AI 结果
)

// AI 搜索模式
const (
	AiSearchModeFirstResult = 1 // 取第一个结果
	AiSearchModeAlgorithm   = 2 // 使用算法选择最佳结果
	AiSearchModeAiDecision  = 3 // 由 AI 决策最佳结果
)
