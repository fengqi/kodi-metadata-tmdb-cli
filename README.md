# kodi-metadata-tmdb-cli

[English](README.md) · [简体中文](README.zh-CN.md)

[![License](https://img.shields.io/badge/license-GPL--3.0-blue.svg)](LICENSE)
[![Go](https://img.shields.io/badge/go-1.27%2B-00ADD8.svg)](go.mod)
[![Release](https://github.com/fengqi/kodi-metadata-tmdb-cli/actions/workflows/release.yml/badge.svg)](https://github.com/fengqi/kodi-metadata-tmdb-cli/releases)

A command-line scraper for **movies**, **TV shows** and **music videos**. It fetches metadata from
[TMDB](https://www.themoviedb.org/) and writes Kodi-compatible NFO files and artwork, so it can replace
Kodi's built-in scraper as well as third-party tools such as tinyMediaManager.

Two collecting strategies are available: **scheduled scanning** and **real-time watching** of newly added
files. After metadata is written, Kodi can be asked to refresh its library automatically.

## Features

- Movie / TV show metadata from TMDB, with proxy support for restricted networks
- Artwork: poster, fanart, clearlogo, season poster, episode thumb
- Cast, genres, tags, country, studio, content rating and TMDB unique id written into NFO
- Scheduled scanning and real-time file watching
- Manual overrides for ambiguous names: TMDB id, season number, episode group, split episodes
- Trigger Kodi library refresh / clean through JSON-RPC after scraping
- Music videos: thumbnail and video/audio stream details extracted with `ffmpeg`
- AI-assisted filename parsing and candidate selection (any OpenAI-compatible endpoint)
- Blu-ray (`BDMV`) and DVD (`VIDEO_TS`) disc-folder detection
- Skip lists for folders and filename keywords
- Multiple media root directories per media type
- Single static binary: Linux (amd64/arm64/arm), macOS (amd64/arm64), Windows (amd64)

## How it works

```
movies_dir / shows_dir / music_videos_dir
        │
        ├── cron scan ──────┐
        └── fsnotify watch ─┤
                            ▼
                   parse file name (rules, optionally AI)
                            ▼
              TMDB search ──► TMDB detail (+ artwork download)
                            ▼
              <name>.nfo + poster.jpg + fanart.jpg + clearlogo.png + …
                            ▼
              Kodi JSON-RPC: refresh movie / TV show / episode
```

A `tmdb/` cache folder is created next to each media file to store the matched TMDB id, the API response
and the files you use for manual overrides.

## Requirements

- A TMDB API key ([v3](https://www.themoviedb.org/settings/api)) — or request it by mail if you cannot
  reach the developer portal
- Kodi 19 (Matrix) or newer; the web server / JSON-RPC must be enabled if you want automatic library refresh
- `ffmpeg` / `ffprobe` — only needed for music videos
- Go 1.27+ — only needed when building from source

## Installation

### 1. Configure Kodi

For every video source that this tool manages:

```
Settings → Media → Videos → Library → <your source> → Change content
  This directory contains : Movies  (or TV shows)
  Information provider    : Local information only
```

`Local information only` makes Kodi read the NFO files instead of scraping online by itself. Keep movies
and TV shows in separate sources.

### 2. Download

Grab the archive for your platform from
[Releases](https://github.com/fengqi/kodi-metadata-tmdb-cli/releases). Each archive contains the binary and
`example.config.json`.

### 3. Configure and run

```bash
cp example.config.json config.json
vi config.json                      # set tmdb.api_key, directories, kodi.json_rpc …
./kodi-tmdb-linux-amd64 -config config.json
```

> The process must run in the same environment as your download client (Transmission, µTorrent, qBittorrent,
> …) and see the same paths, otherwise real-time watching will not notice new files.

Keep it running in the background — systemd, `nohup`, a container, or your NAS' task scheduler.

## Usage

```
Usage of kodi-tmdb:
  -config string   config file, read from working dir first, then binary dir (default "config.json")
  -mode int        run mode: 1: daemon, 2: once, 3: spec
  -version         display version
```

| Run mode | Value | Description |
| --- | --- | --- |
| daemon | `1` | Long-running process: optional startup scan + scheduled scan + file watching |
| once | `2` | Scan all configured directories once, then exit |
| spec | `3` | Scan only the current working directory, then exit — handy as a one-shot hook |

`-mode` overrides `collector.run_mode` from the config file.

## Configuration

`config.json` is looked up in the working directory first; if it does not exist there, the directory of the
binary is used. Both `/` and `\` path separators are accepted, so the same file works on Linux and Windows.

### `log`

| Key | Type | Description |
| --- | --- | --- |
| `mode` | int | `1` stdout, `2` log file, `3` both. Default `1` |
| `level` | int | `0` debug, `1` info, `2` warning, `3` error, `4` fatal. Default `1` |
| `file` | string | Log file path, used when `mode` is `2` or `3` |

### `tmdb`

| Key | Type | Description |
| --- | --- | --- |
| `api_host` | string | TMDB API host, default `https://api.themoviedb.org` |
| `image_host` | string | TMDB image host, default `https://image.tmdb.org` |
| `api_key` | string | TMDB v3 API key |
| `language` | string | Metadata language, e.g. `zh-CN`, `en-US` |
| `rating` | string | Region used for content rating, e.g. `US` |
| `proxy` | string | Proxy for TMDB requests: `http`, `https`, `socks5`, `socks5h`. Format `scheme://[user:password@]host:port` |
| `timeout_seconds` | int | Request timeout in seconds, default `30` |
| `retry_count` | int | Retry attempts on failure, `0` means no retry |

### `collector`

| Key | Type | Description |
| --- | --- | --- |
| `run_mode` | int | `1` daemon, `2` once, `3` spec |
| `watcher` | bool | Watch the media directories for new files (daemon mode) |
| `cron_scan` | bool | Enable scheduled scanning |
| `cron_seconds` | int | Interval of the scheduled scan, in seconds |
| `cron_scan_boot` | bool | Run a full scan immediately after the daemon starts |
| `cron_scan_kodi` | bool | Ask Kodi to scan its library after every scheduled scan |
| `skip_folders` | []string | Folder names skipped entirely (matched by folder name) |
| `skip_keywords` | []string | Keywords dropped while parsing titles, e.g. `纯享`, `pure` |
| `nfo_field.tag` | bool | Write tags into the NFO |
| `nfo_field.genre` | bool | Write genres into the NFO |
| `movies_dir` | []string | Movie root directories |
| `shows_dir` | []string | TV show root directories |
| `music_videos_dir` | []string | Music video root directories |
| `tmp_suffix` | []string | Reserved — not effective in the current version |
| `movies_nfo_mode` | int | Reserved — `<VideoFileName>.nfo` is always written |

Notes:

- Directories and files whose name starts with `.` are always skipped.
- `skip_folders` matches a single directory name, not a path. The defaults cover NAS/PT leftovers such as
  `@eaDir`, `Sample`, `CERTIFICATE`, `$RECYCLE.BIN`.

### `kodi`

| Key | Type | Description |
| --- | --- | --- |
| `enable` | bool | Send refresh / scan requests to Kodi |
| `clean_library` | bool | Ask Kodi to clean its library after scanning |
| `json_rpc` | string | JSON-RPC endpoint, e.g. `http://192.168.1.123:8080/jsonrpc` |
| `timeout` | int | Connection timeout in seconds |
| `username` / `password` | string | Credentials, if Kodi's web server requires authentication |

Enable the web server in Kodi under `Settings → Services → Control → Allow remote control via HTTP`.

### `ffmpeg`

| Key | Type | Description |
| --- | --- | --- |
| `max_worker` | int | Number of concurrent ffmpeg processes; the logical CPU count is a good value |
| `ffmpeg_path` | string | Path to the `ffmpeg` executable |
| `ffprobe_path` | string | Path to the `ffprobe` executable |

### `ai`

AI is optional and used **only** to parse file names and to choose between TMDB search results. All metadata
itself always comes from TMDB.

| Key | Type | Description |
| --- | --- | --- |
| `enable` | bool | Turn AI assistance on |
| `base_url` | string | OpenAI-compatible chat completion endpoint |
| `api_key` | string | API key |
| `model` | string | Model name |
| `temperature` | float | Sampling temperature |
| `timeout_seconds` | int | Request timeout, default `15` |
| `confidence_threshold` | float | Minimum confidence accepted for an AI result, default `0.7` |
| `match_mode` | int | `1` rules first, AI as fallback; `2` AI first, rules as fallback; `3` rules first, then AI result overrides |
| `search_mode` | int | How to pick from TMDB results: `1` first result, `2` score-based algorithm, `3` AI decides (falls back to `2`) |

## Manual matching

Names that cannot be parsed reliably can be pinned down with small text files inside the `tmdb/` cache
folder. The tool keeps everything it needs next to the media, so a wrong match is easy to repair.

### Movies

```
/movies/Hawkeye.2021.1080p.BluRay.x264/
├── Hawkeye.2021.1080p.BluRay.x264.mkv             # video
├── Hawkeye.2021.1080p.BluRay.x264.nfo             # generated
├── Hawkeye.2021.1080p.BluRay.x264-poster.jpg      # generated
├── Hawkeye.2021.1080p.BluRay.x264-fanart.jpg      # generated
├── Hawkeye.2021.1080p.BluRay.x264-clearlogo.png   # generated
└── tmdb/
    ├── Hawkeye.2021.1080p.BluRay.x264.id.txt      # TMDB movie id (override / cache)
    └── Hawkeye.2021.1080p.BluRay.x264.movie.json  # cached TMDB detail
```

Blu-ray folders are scraped by folder name and produce `index.nfo`; DVD folders produce
`VIDEO_TS/VIDEO_TS.nfo`.

### TV shows

```
/shows/Hawkeye (2021)/
├── poster.jpg / fanart.jpg / clearlogo.png / season01-poster.jpg
├── tmdb/
│   ├── id.txt        # TMDB tv id
│   ├── tv.json       # cached TMDB detail
│   ├── group.txt     # TMDB episode group id (optional)
│   └── <episode>.episode.json
└── Season 01/
    ├── tmdb/
    │   ├── season.txt   # force the season number
    │   └── join.txt     # merge split episodes
    ├── Hawkeye.S01E01.1080p.WEB-DL.mkv
    ├── Hawkeye.S01E01.1080p.WEB-DL.nfo
    └── Hawkeye.S01E01.1080p.WEB-DL-thumb.jpg
```

### Override files

| File | Location | Purpose |
| --- | --- | --- |
| `<VideoFileName>.id.txt` | `tmdb/` next to the movie | Force the TMDB **movie** id. Put the numeric id as plain text |
| `id.txt` | `<show root>/tmdb/` | Force the TMDB **TV show** id |
| `season.txt` | `<season dir>/tmdb/` or `<show root>/tmdb/` | Force the season number. `0` means specials |
| `group.txt` | `<show root>/tmdb/` or `<season dir>/tmdb/` | Use a TMDB **episode group** (`.../tv/{id}/episode_group/{group_id}`); the group order is treated as the season number |
| `join.txt` | `<season dir>/tmdb/` | Merge split episodes. Format `currentSeason,newSeason,episodeOffset`, e.g. `1,1,13` maps `S01E01` to `S01E13` |

When you change an id, delete the cached `*.movie.json` / `tv.json` / `*.episode.json` next to it and re-run
the tool (or just re-scan) so the new id is fetched.

Because the whole cache lives in `tmdb/`, deleting `tmdb/*.json` forces a full re-scrape of that title.

## Building from source

```bash
git clone https://github.com/fengqi/kodi-metadata-tmdb-cli.git
cd kodi-metadata-tmdb-cli

make linux-amd64        # one target: linux-amd64, linux-arm64, linux-arm, darwin-amd64, darwin-arm64, windows-amd64
make all                # every target above
make release            # build + zip, output in ./release

go test ./...           # run the test suite
```

The binary is built with `CGO_ENABLED=0` and `-trimpath`, so the artifacts are static.

## References

- Project wiki (deployment, options, troubleshooting — start here when something does not work) — https://github.com/fengqi/kodi-metadata-tmdb-cli/wiki
- Kodi v19 (Matrix) JSON-RPC API v12 — https://kodi.wiki/view/JSON-RPC_API/v12
- Kodi v19 (Matrix) NFO files — https://kodi.wiki/view/NFO_files
- Kodi artwork types — https://kodi.wiki/view/Artwork_types
- TMDB API overview — https://www.themoviedb.org/documentation/api
- TMDB API v3 — https://developer.themoviedb.org/docs
- fsnotify — https://github.com/fsnotify/fsnotify
- tinyMediaManager — https://gitlab.com/tinyMediaManager/tinyMediaManager

## License

[GPL-3.0](LICENSE)
