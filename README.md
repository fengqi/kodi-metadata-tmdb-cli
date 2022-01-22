# kodi-metadata-tmdb-cli

Kodi 刮削器命令行版本，使用TMDB数据源匹配剧集和电影并生成Kodi兼容的nfo文件。

有定时扫描刮削、实时监听新增文件刮削两种模式。

# 怎么使用

1. 打开 Kodi 设置 - 媒体 - 视频  - 更改内容（仅限电影和剧集类型） - 信息提供者改为：Local information only
2. 根据平台[下载](https://github.com/fengqi/kodi-metadata-tmdb-cli/releases)对应的文件，配置 `config.json`并后台运行。

> 本程序必须和下载软件（如迅雷、µTorrent等）运行在同一个环境，不然实时监听不生效。

# 配置字段说明

- log_level 日志等级，0-4分别对应：DEBUG、INFO、WARNING、ERROR、FATAL
- log_file 日志文件路径
- cron_seconds 定时扫描间隔，单位(秒)
- rating 内容分级国家
- api_key TMDB 开发者token，请参考[Wiki](https://github.com/fengqi/kodi-metadata-tmdb-cli/wiki)申请
- language 刮削语言，中文可以填：zh-CN
- shows_dir 电视剧、电视节目目录，可以多个
- movies_dir 电影目录，可以多个
- stock_dir 不可描述的目录，功能暂未实现
- music_dir 音乐目录，功能暂未实现

# 功能列表

- [x] 从TMDB获取电影、电视剧、电视剧分集信息
- [x] 从TMDB获取电影、电视剧演员列表信息
- [x] 从TMDB下载封面、海报等图片
- [x] 定时扫描电影、电视剧
- [x] 实时监听新添加的电影、电视剧
- [x] 命名不规范或有歧义的电影、电视剧支持手动指定id
- [x] 命名不规范的电视剧支持指定season
- [ ] 命名不规范的电视剧支持指定season和episode的对应，如：E21-34其实是第2季
- [ ] 从TMDB获取电影、电视剧内容分级信息
- [ ] 多个搜索结果尝试根据特征信息确定
- [ ] 适配其他数据源，如：imdb、tvdb等以补全部分tmdb没有的数据
- [x] 支持电影合集
- [x] 支持电视剧合集
- [ ] 更新NFO文件后触发kodi更新数据
- [x] 支持单个电影文件和目录
- [x] 识别蓝光电影目录
- [x] 支持 .part 和 .!qb 文件
- [ ] 支持音乐刮削
- [ ] 支持AV刮削

# 参考

- Kodi v19 (Matrix) JSON-RPC API/V12 https://kodi.wiki/view/JSON-RPC_API/v12
- Kodi v19 (Matrix) NFO files https://kodi.wiki/view/NFO_files
- TMDB Api Overview https://www.themoviedb.org/documentation/api
- TMDB Api V3 https://developers.themoviedb.org/3/getting-started/introduction
- File system notifications for Go https://github.com/fsnotify/fsnotify
