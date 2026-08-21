# Changelog

本文件记录本仓库所有值得下游关注的变更，格式遵循 [Keep a Changelog 1.1.0](https://keepachangelog.com/zh-CN/1.1.0/)，版本号遵循 [语义化版本 2.0.0](https://semver.org/lang/zh-CN/)。

## [Unreleased]

## [2.1.1] - 2026-08-21

### Fixed

- `v2`：修正 README 中 `jsonv2` 后端性能说明的归因。原文把慢于默认后端的原因归给 `encoding/json` 的实现基础，实际来自本后端为跨后端一致而启用的三项对齐选项与缩进的额外拷贝；直接对比标准库两个包，`encoding/json/v2` 反而更快

### Changed

- `v2`：README 新增「标准库 struct tag 与接口的后端支持」小节，给出 `omitzero`、`omitzero` 认 `IsZero()`、`case:strict`、`case:ignore`、`embed` 兜底字段、`MarshalerTo` / `UnmarshalerFrom` 在五个后端下的实测支持情况——依赖这些能力时需选择默认后端、`jsonv2` 或 `sonic`

## [2.1.0] - 2026-08-21

### Added

- `v2`：新增 `jsonv2` build tag，使用 Go 1.27 标准库的 `encoding/json/v2` 作为编解码后端，零外部依赖。`-tags=jsonv2` 下 `json.Package` 为 `"encoding/json/v2"`，顶层函数、`Core` / `Encoder` / `Decoder` 接口与类型别名的签名保持不变，业务代码无需修改
- `v2`：README 新增「encoding/json/v2 后端」小节，列出与默认后端的行为差异与选型建议

### Changed

- `v2`：模块要求的 Go 版本提升到 1.27（`encoding/json/v2` 的最低版本）。默认后端仍为 `encoding/json`，不加 build tag 时行为不变

### 使用 `-tags=jsonv2` 前请注意

该后端采用 RFC 7493 语义，与默认后端存在可观察差异，完整清单见 [v2/README.md](v2/README.md#与默认后端的行为差异)。唯一需要改代码的一条：**结构体字段名严格区分大小写**，JSON 成员名与字段名大小写不一致时字段保持零值且不返回错误——为字段补上 `json` tag 即可，该写法在所有后端下都正确。

以下三项不在差异之列，本后端显式与默认后端及 `sonic` / `go_json` / `jsoniter` 对齐，输出逐字节一致：map 键顺序按键排序、`time.Duration` 编解码为纳秒整数、`omitempty` 省略零值。若采用 `encoding/json/v2` 的默认值，这三项分别会导致输出字节不稳定、含 `Duration` 字段的结构体编解码直接报错、响应体多出零值字段。
