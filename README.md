# kodi-metadata-tmdb-cli

电影、电视剧刮削器命令行版本，使用TMDB数据源生成Kodi兼容的NFO文件和相关图片，可用来代替Kodi自带以及tinyMediaManager等其他第三方的刮削器。

有定时扫描扫描、实时监听新增文件两种模式，可配置有新增时触发Kodi更新媒体库。

# 使用

1. 打开 Kodi 设置 - 媒体 - 视频  - 更改内容（仅限电影和剧集类型） - 信息提供者改为：Local information only
2. 根据平台[下载](https://github.com/fengqi/kodi-metadata-tmdb-cli/releases)对应的文件，配置 `config.json`并后台运行。

> 本程序必须和下载软件（如Transmission、µTorrent等）运行在同一个环境，不然实时监听模式不生效。
> 详细配置参考 [配置总览](https://github.com/fengqi/kodi-metadata-tmdb-cli/wiki/%E9%85%8D%E7%BD%AE%E6%96%87%E4%BB%B6)

# 功能

- [x] 从TMDB获取电视剧、电影信息，支持指定代理
- [x] 从TMDB获取演员列表、封面图片、海报图片、内容分级、logo图片
- [x] 定时扫描电影、电视剧、音乐视频
- [x] 实时监听新添加的电影、电视剧、音乐视频
- [x] 命名不规范或有歧义的电影、电视剧手动指定tmdb id
- [x] 命名不规范的电视剧手动指定season
- [x] 更新NFO文件后触发Kodi更新媒体库
- [x] 音乐视频使用ffmpeg提取缩略图、视频音频信息

# 代办

- [ ] 多个搜索结果根据特征信息排序
- [ ] 手动指定电视剧分组

# 参考

- Kodi v19 (Matrix) JSON-RPC API/V12 https://kodi.wiki/view/JSON-RPC_API/v12
- Kodi v19 (Matrix) NFO files https://kodi.wiki/view/NFO_files
- Kodi Artwork types https://kodi.wiki/view/Artwork_types
- TMDB Api Overview https://www.themoviedb.org/documentation/api
- TMDB Api V3 https://developers.themoviedb.org/3/getting-started/introduction
- File system notifications for Go https://github.com/fsnotify/fsnotify
- tinyMediaManager https://gitlab.com/tinyMediaManager/tinyMediaManager
