# kodi-metadata-tmdb-cli

[English](README.md) · [简体中文](README.zh-CN.md)

[![License](https://img.shields.io/badge/license-GPL--3.0-blue.svg)](LICENSE)
[![Go](https://img.shields.io/badge/go-1.27%2B-00ADD8.svg)](go.mod)
[![Release](https://github.com/fengqi/kodi-metadata-tmdb-cli/actions/workflows/release.yml/badge.svg)](https://github.com/fengqi/kodi-metadata-tmdb-cli/releases)

电影、电视剧、音乐视频刮削器的命令行版本。使用 [TMDB](https://www.themoviedb.org/) 作为数据源，生成 Kodi
兼容的 NFO 文件和相关图片，可用来代替 Kodi 自带刮削器以及 tinyMediaManager 等第三方工具。

提供**定时扫描**和**实时监听新增文件**两种模式，刮削完成后可以自动触发 Kodi 更新媒体库。

## 功能

- 从 TMDB 获取电影、电视剧信息，支持配置代理
- 图片资源：海报、艺术图（fanart）、logo、季海报、剧集缩略图
- 演员列表、分类、标签、国家、制片公司、内容分级、TMDB 唯一 ID 写入 NFO
- 定时扫描与实时文件监听
- 命名不规范或有歧义时，可手动指定 TMDB ID、季、剧集分组、剧集合并
- 刮削完成后通过 JSON-RPC 触发 Kodi 刷新 / 清理媒体库
- 音乐视频使用 ffmpeg 提取缩略图和音视频流信息
- 支持 AI 匹配：解析文件名、从搜索结果中挑选最合适的结果（兼容 OpenAI 接口）
- 识别蓝光（`BDMV`）、DVD（`VIDEO_TS`）原盘目录
- 支持跳过目录和文件名关键字
- 电影、电视剧、音乐视频均可配置多个根目录
- 单一静态二进制：Linux（amd64/arm64/arm）、macOS（amd64/arm64）、Windows（amd64）

## 工作原理

```
movies_dir / shows_dir / music_videos_dir
        │
        ├── 定时扫描 ───────┐
        └── fsnotify 监听 ──┤
                            ▼
                   解析文件名（规则，可选 AI）
                            ▼
              TMDB 搜索 ──► TMDB 详情（并下载图片）
                            ▼
              <文件名>.nfo + poster.jpg + fanart.jpg + clearlogo.png + …
                            ▼
              Kodi JSON-RPC：刷新电影 / 剧集 / 单集
```

每个媒体文件旁边会生成一个 `tmdb/` 缓存目录，用于保存匹配到的 TMDB ID、接口返回内容，以及你需要手动
指定的覆盖文件。

## 环境要求

- TMDB API Key（[v3](https://www.themoviedb.org/settings/api)）
- Kodi 19 (Matrix) 或更高版本；需要自动刷新媒体库时，必须开启 Kodi 的 Web 服务 / JSON-RPC
- `ffmpeg` / `ffprobe` —— 仅音乐视频需要
- Go 1.27+ —— 仅从源码编译时需要

## 使用

### 1. 配置 Kodi

对所有由本程序管理的视频源：

```
设置 → 媒体 → 视频 → 媒体库 → <你的视频源> → 更改内容
  该目录包含   ：电影（或 剧集）
  信息提供者   ：Local information only
```

选择 `Local information only` 后，Kodi 会直接读取 NFO，而不再自行联网刮削。电影和电视剧请使用不同的
内容源。

### 2. 下载

从 [Releases](https://github.com/fengqi/kodi-metadata-tmdb-cli/releases) 下载对应平台的压缩包，包内包含
可执行文件和 `example.config.json`。

### 3. 配置并运行

```bash
cp example.config.json config.json
vi config.json                      # 配置 tmdb.api_key、媒体目录、kodi.json_rpc 等
./kodi-tmdb-linux-amd64 -config config.json
```

> 本程序必须和下载软件（Transmission、µTorrent、qBittorrent 等）运行在同一个环境，并且能看到相同的路径，
> 否则实时监听模式不生效。

后台常驻方式可自行选择：systemd、`nohup`、容器，或 NAS 自带的任务计划。

## 命令行参数

```
Usage of kodi-tmdb:
  -config string   config file, read from working dir first, then binary dir (default "config.json")
  -mode int        run mode: 1: daemon, 2: once, 3: spec
  -version         display version
```

| 运行模式 | 值 | 说明 |
| --- | --- | --- |
| daemon | `1` | 守护进程：可选的启动扫描 + 定时扫描 + 实时监听 |
| once | `2` | 扫描所有配置目录后退出 |
| spec | `3` | 只扫描当前工作目录后退出，适合作为一次性钩子调用 |

`-mode` 参数会覆盖配置文件中的 `collector.run_mode`。

## 配置说明

程序优先从当前工作目录读取 `config.json`，不存在时回退到可执行文件所在目录。路径分隔符 `/` 和 `\` 均可
使用，同一份配置可以跨平台使用。

### `log` 日志

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| `mode` | int | `1` 仅标准输出，`2` 仅日志文件，`3` 两者都输出。默认 `1` |
| `level` | int | `0` debug，`1` info，`2` warning，`3` error，`4` fatal。默认 `1` |
| `file` | string | 日志文件路径，`mode` 为 `2` 或 `3` 时生效 |

### `tmdb`

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| `api_host` | string | TMDB 接口地址，默认 `https://api.themoviedb.org` |
| `image_host` | string | TMDB 图片地址，默认 `https://image.tmdb.org` |
| `api_key` | string | TMDB v3 API Key |
| `language` | string | 元数据语言，如 `zh-CN`、`en-US` |
| `rating` | string | 内容分级使用的地区，如 `US` |
| `proxy` | string | 请求 TMDB 的代理，支持 `http`、`https`、`socks5`、`socks5h`，格式 `scheme://[user:password@]host:port` |
| `timeout_seconds` | int | 请求超时时间（秒），默认 `30` |
| `retry_count` | int | 请求失败重试次数，`0` 表示不重试 |

### `collector` 刮削

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| `run_mode` | int | `1` daemon，`2` once，`3` spec |
| `watcher` | bool | 是否开启实时文件监听（守护进程模式） |
| `cron_scan` | bool | 是否开启定时扫描 |
| `cron_seconds` | int | 定时扫描间隔，单位秒 |
| `cron_scan_boot` | bool | 守护进程启动后立即执行一次扫描 |
| `cron_scan_kodi` | bool | 每次定时扫描结束后触发 Kodi 扫描媒体库 |
| `skip_folders` | []string | 完全跳过的目录名（按目录名匹配） |
| `skip_keywords` | []string | 解析标题时丢弃的关键字，如 `纯享`、`pure` |
| `nfo_field.tag` | bool | 是否将标签写入 NFO |
| `nfo_field.genre` | bool | 是否将分类写入 NFO |
| `movies_dir` | []string | 电影根目录，可多个 |
| `shows_dir` | []string | 电视剧根目录，可多个 |
| `music_videos_dir` | []string | 音乐视频根目录，可多个 |
| `tmp_suffix` | []string | 预留字段，当前版本未生效 |
| `movies_nfo_mode` | int | 预留字段，当前版本固定写入 `<视频文件名>.nfo` |

补充说明：

- 以 `.` 开头的文件和目录始终会被跳过。
- `skip_folders` 匹配的是单个目录名，而不是路径。默认值已覆盖 `@eaDir`、`Sample`、`CERTIFICATE`、
  `$RECYCLE.BIN` 等 NAS / PT 常见目录。

### `kodi`

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| `enable` | bool | 是否向 Kodi 发送刷新 / 扫描请求 |
| `clean_library` | bool | 扫描结束后是否让 Kodi 清理媒体库 |
| `json_rpc` | string | JSON-RPC 地址，如 `http://192.168.1.123:8080/jsonrpc` |
| `timeout` | int | 连接超时时间（秒） |
| `username` / `password` | string | Kodi Web 服务开启鉴权时使用 |

需要在 Kodi 中开启 Web 服务：`设置 → 服务 → 控制 → 允许通过 HTTP 远程控制`。

### `ffmpeg`

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| `max_worker` | int | ffmpeg 最大并发进程数，建议设置为逻辑 CPU 个数 |
| `ffmpeg_path` | string | `ffmpeg` 可执行文件路径 |
| `ffprobe_path` | string | `ffprobe` 可执行文件路径 |

### `ai`

AI 是可选项，**只用于**解析文件名和从 TMDB 搜索结果中挑选最合适的结果，元数据本身始终来自 TMDB。

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| `enable` | bool | 是否开启 AI |
| `base_url` | string | OpenAI 兼容的对话接口地址 |
| `api_key` | string | API Key |
| `model` | string | 模型名称 |
| `temperature` | float | 采样温度 |
| `timeout_seconds` | int | 请求超时时间（秒），默认 `15` |
| `confidence_threshold` | float | 可接受的最低置信度，默认 `0.7` |
| `match_mode` | int | `1` 规则优先、匹配不到再 AI 介入；`2` AI 优先、匹配不到再规则介入；`3` 规则优先，结果给 AI 参考，最终使用 AI 结果 |
| `search_mode` | int | 搜索结果选择方式：`1` 取第一个结果，`2` 使用算法打分选择，`3` 由 AI 决策（失败时回退到 `2`） |

## 手动指定匹配

规则解析不出来的名称，可以在 `tmdb/` 缓存目录中放一个小的文本文件来手动指定。所有缓存都在媒体文件旁边，
修改和排查都很方便。

### 电影

```
/movies/鹰眼 Hawkeye 2021/
├── 鹰眼 Hawkeye 2021.mkv             # 视频文件
├── 鹰眼 Hawkeye 2021.nfo             # 自动生成
├── 鹰眼 Hawkeye 2021-poster.jpg      # 自动生成
├── 鹰眼 Hawkeye 2021-fanart.jpg      # 自动生成
├── 鹰眼 Hawkeye 2021-clearlogo.png   # 自动生成
└── tmdb/
    ├── 鹰眼 Hawkeye 2021.id.txt      # 手动指定的 TMDB 电影 ID
    └── 鹰眼 Hawkeye 2021.movie.json  # TMDB 详情缓存
```

蓝光目录按目录名刮削，NFO 写入 `index.nfo`；DVD 目录 NFO 写入 `VIDEO_TS/VIDEO_TS.nfo`。

### 电视剧

```
/shows/鹰眼 Hawkeye (2021)/
├── poster.jpg / fanart.jpg / clearlogo.png / season01-poster.jpg
├── tmdb/
│   ├── id.txt        # 手动指定的 TMDB 剧集 ID
│   ├── tv.json       # TMDB 详情缓存
│   ├── group.txt     # TMDB 剧集分组 ID（可选）
│   └── <单集文件名>.episode.json
└── Season 01/
    ├── tmdb/
    │   ├── season.txt   # 手动指定季
    │   └── join.txt     # 剧集合并
    ├── 鹰眼 Hawkeye S01E01.mkv
    ├── 鹰眼 Hawkeye S01E01.nfo
    └── 鹰眼 Hawkeye S01E01-thumb.jpg
```

### 覆盖文件一览

| 文件 | 位置 | 作用 |
| --- | --- | --- |
| `<视频文件名>.id.txt` | 电影旁的 `tmdb/` | 指定 TMDB **电影** ID，内容为纯数字 |
| `id.txt` | `<剧集根目录>/tmdb/` | 指定 TMDB **剧集** ID |
| `season.txt` | `<季目录>/tmdb/` 或 `<剧集根目录>/tmdb/` | 指定季号，`0` 表示特别篇 |
| `group.txt` | `<剧集根目录>/tmdb/` 或 `<季目录>/tmdb/` | 使用 TMDB **剧集分组**（`.../tv/{id}/episode_group/{group_id}`），分组的排序会被当作季号 |
| `join.txt` | `<季目录>/tmdb/` | 剧集合并，格式 `当前季,新的季,集数偏移`，如 `1,1,13` 会把 `S01E01` 映射为 `S01E13` |

修改 ID 之后，请删除同目录下的 `*.movie.json` / `tv.json` / `*.episode.json` 缓存并重新扫描，让程序按新的
ID 重新拉取数据。

因为缓存全部集中在 `tmdb/` 目录，删除该目录下的 `*.json` 即可强制重新刮削该条目。

## 从源码编译

```bash
git clone https://github.com/fengqi/kodi-metadata-tmdb-cli.git
cd kodi-metadata-tmdb-cli

make linux-amd64        # 单个目标：linux-amd64、linux-arm64、linux-arm、darwin-amd64、darwin-arm64、windows-amd64
make all                # 编译以上全部平台
make release            # 编译并打包，产物在 ./release

go test ./...           # 运行测试
```

编译使用 `CGO_ENABLED=0` 与 `-trimpath`，产物为静态二进制。

## 参考

- 项目 Wiki（部署方式、配置项细节、问题排查，遇到问题先看这里）https://github.com/fengqi/kodi-metadata-tmdb-cli/wiki
- Kodi v19 (Matrix) JSON-RPC API/V12 https://kodi.wiki/view/JSON-RPC_API/v12
- Kodi v19 (Matrix) NFO files https://kodi.wiki/view/NFO_files
- Kodi Artwork types https://kodi.wiki/view/Artwork_types
- TMDB Api Overview https://www.themoviedb.org/documentation/api
- TMDB Api V3 https://developer.themoviedb.org/docs
- File system notifications for Go https://github.com/fsnotify/fsnotify
- tinyMediaManager https://gitlab.com/tinyMediaManager/tinyMediaManager

## 许可证

[GPL-3.0](LICENSE)
