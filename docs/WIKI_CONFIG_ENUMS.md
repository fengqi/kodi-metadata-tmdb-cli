# 配置枚举值说明（数字配置 + 代码常量）

本文档用于 GitHub Wiki，说明项目里采用数字枚举的配置字段，以及默认值和非法值回退行为。

## 设计原则

- 配置文件保持简洁：使用 `1/2/3` 这类数字。
- 代码保持可读：业务逻辑中统一使用常量名，不直接写魔法值。
- 配置保持稳健：遇到非法枚举值时自动回退到默认值并输出 warning。

## 枚举字段一览

| 字段 | 可选值 | 默认值 | 非法值回退 | 代码常量 |
|---|---|---|---|---|
| `collector.run_mode` | `1=daemon`、`2=once`、`3=spec` | `1` | 回退到 `1` | `CollectorRunModeDaemon/Once/Spec` |
| `log.mode` | `1=stdout`、`2=logfile`、`3=both` | `1` | 回退到 `1` | `LogModeStdout/Logfile/Both` |
| `log.level` | `0=debug`、`1=info`、`2=warning`、`3=error`、`4=fatal` | `1` | 回退到 `1` | `LogLevelDebug/Info/Warning/Error/Fatal` |
| `collector.movies_nfo_mode` | `1=movie.nfo`、`2=<VideoFileName>.nfo` | `2` | 回退到 `2` | `CollectorMoviesNfoModeMovieNfo/VideoNfo` |

## 示例配置

```json
{
  "log": {
    "mode": 1,
    "level": 1,
    "file": "./tmdb-collector.log"
  },
  "collector": {
    "run_mode": 1,
    "movies_nfo_mode": 2
  }
}
```

## 行为说明

- `collector.run_mode`
  - `1`：守护模式，监听 + 定时扫描。
  - `2`：单次模式，扫描一次后退出。
  - `3`：临时模式，按当前目录识别并处理一次。
- `log.mode`
  - `1`：仅标准输出。
  - `2`：仅写日志文件。
  - `3`：同时写标准输出和日志文件。
- `log.level`
  - `0~4` 从 `debug` 到 `fatal`，数字越大级别越高。
- `collector.movies_nfo_mode`
  - 目前仅保留配置和枚举定义，当前版本尚未接入实际写入策略切换。

## 建议值

- 普通用户建议：
  - `collector.run_mode = 1`
  - `log.mode = 1`
  - `log.level = 1`
  - `collector.movies_nfo_mode = 2`

## 变更说明（实现层）

- 所有枚举常量定义集中在 `config/constants.go`。
- `collector`、`main`、`logger` 中涉及模式判断全部改为常量判断。
- `LoadConfig` 增加轻量枚举校验与回退逻辑。
