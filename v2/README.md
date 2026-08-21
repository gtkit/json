# gtkit/json/v2

[![Go Reference](https://pkg.go.dev/badge/github.com/gtkit/json/v2.svg)](https://pkg.go.dev/github.com/gtkit/json/v2)
[![Go Version](https://img.shields.io/badge/go-%3E%3D1.27-blue)](https://go.dev)

通过 build tags 无缝切换 JSON 编解码后端的 Go 库。零修改业务代码，一个 `-tags` 参数即可获得数倍性能提升。

## 支持的后端

| Build Tag | 后端库 | 适用场景 |
|-----------|--------|---------|
| _(默认)_ | `encoding/json` | 零依赖、最大兼容性 |
| `jsonv2` | `encoding/json/v2` | 零依赖，需要 RFC 7493 严格语义（拒绝重复键、拒绝无效 UTF-8） |
| `sonic` | [bytedance/sonic](https://github.com/bytedance/sonic) | 追求极致性能（Linux/macOS/Windows，amd64/arm64） |
| `go_json` | [goccy/go-json](https://github.com/goccy/go-json) | 高性能且全平台兼容 |
| `jsoniter` | [json-iterator/go](https://github.com/json-iterator/go) | PHP 兼容模式、私有字段支持 |

## 安装

```bash
go get github.com/gtkit/json/v2
```

## 快速开始

```go
package main

import (
    "fmt"

    "github.com/gtkit/json/v2"
)

type User struct {
    Name  string `json:"name"`
    Email string `json:"email"`
    Age   int    `json:"age"`
}

func main() {
    // 查看当前使用的 JSON 后端
    fmt.Println("JSON backend:", json.Package)

    u := User{Name: "张三", Email: "zhangsan@example.com", Age: 30}

    // Marshal
    data, err := json.Marshal(u)
    if err != nil {
        panic(err)
    }
    fmt.Println(string(data))

    // Unmarshal
    var u2 User
    if err := json.Unmarshal(data, &u2); err != nil {
        panic(err)
    }
    fmt.Printf("%+v\n", u2)

    // MarshalToString — 零拷贝返回 string
    s, err := json.MarshalToString(u)
    if err != nil {
        panic(err)
    }
    fmt.Println(s)

    // Valid — 校验 JSON 合法性
    fmt.Println("valid:", json.Valid(data))
}
```

## 构建方式

```bash
# 默认 encoding/json
go build ./...

# 使用 encoding/json/v2（标准库，需 Go 1.27+）
go build -tags=jsonv2 ./...

# 使用 sonic（推荐生产环境，Linux/macOS/Windows）
go build -tags=sonic ./...

# 使用 go-json
go build -tags=go_json ./...

# 使用 jsoniter
go build -tags=jsoniter ./...
```

### 交叉编译注意事项

sonic 依赖 JIT 或汇编优化，仅支持以下平台：

- `linux/amd64`、`linux/arm64`
- `darwin/amd64`、`darwin/arm64`
- `windows/amd64`

如果你的目标平台不在此列表中，可以使用 `go_json`，或用标准库的 `jsonv2`（实测 `freebsd/amd64`、`linux/arm`、`linux/riscv64`、`windows/386`、`js/wasm` 均可构建）：

```bash
GOOS=freebsd GOARCH=amd64 go build -tags=go_json ./...
GOOS=js GOARCH=wasm go build -tags=jsonv2 ./...
```

## API 一览

### 顶层函数

```go
json.Marshal(v any) ([]byte, error)
json.Unmarshal(data []byte, v any) error
json.MarshalIndent(v any, prefix, indent string) ([]byte, error)
json.MarshalToString(v any) (string, error)
json.NewEncoder(w io.Writer) json.Encoder
json.NewDecoder(r io.Reader) json.Decoder
json.Valid(data []byte) bool
json.Compact(dst *bytes.Buffer, src []byte) error
json.Indent(dst *bytes.Buffer, src []byte, prefix, indent string) error
json.HTMLEscape(dst *bytes.Buffer, src []byte)
```

### 兼容类型

```go
type RawMessage = encoding/json.RawMessage
type Number = encoding/json.Number
type Marshaler = encoding/json.Marshaler
type Unmarshaler = encoding/json.Unmarshaler
```

### Core 接口

所有顶层函数都委托给 `json.API`（类型为 `json.Core`），可用于依赖注入和测试 mock：

```go
type Core interface {
    Marshal(v any) ([]byte, error)
    Unmarshal(data []byte, v any) error
    MarshalIndent(v any, prefix, indent string) ([]byte, error)
    MarshalToString(v any) (string, error)
    NewEncoder(writer io.Writer) Encoder
    NewDecoder(reader io.Reader) Decoder
    Valid(data []byte) bool
}
```

### Encoder / Decoder 接口

```go
type Encoder interface {
    SetEscapeHTML(on bool)
    SetIndent(prefix, indent string)
    Encode(v any) error
}

type Decoder interface {
    UseNumber()
    DisallowUnknownFields()
    Decode(v any) error
    Buffered() io.Reader
    More() bool
}
```

### 常量

```go
json.Package  // 当前后端库名，如 "encoding/json"、"encoding/json/v2"、"github.com/bytedance/sonic"
json.Version  // 包版本号，如 "v2.0.0"
```

## encoding/json/v2 后端

`-tags=jsonv2` 使用 Go 1.27 标准库的 `encoding/json/v2`，编解码全部走 v2，`json.Package` 为 `"encoding/json/v2"`。顶层函数、`Core` / `Encoder` / `Decoder` 接口与类型别名的签名与其他后端完全一致，业务代码不需要修改。

该后端需要 `jsonv2` GOEXPERIMENT 处于开启状态。它在 Go 1.27 中默认开启，无需任何设置；只有显式 `GOEXPERIMENT=nojsonv2` 构建时 `encoding/json/v2` 不可导入，此时请改用其他 tag。

### 与默认后端的行为差异

v2 采用 RFC 7493 语义，以下差异在切换前必须确认。表中每一行都有对应测试锁定（`jsonv2_test.go`），实测于 Go 1.27。

| 行为 | 默认后端 | `-tags=jsonv2` |
|------|---------|----------------|
| **字段名匹配** | 大小写不敏感回落 | **严格区分大小写，不匹配时留零值且不报错** |
| HTML 字符 `<` `>` `&` | 转义为 `\u003c` 等 | 原样输出 |
| 重复对象键 | 取最后一个 | `Unmarshal` 报错、`Valid` 返回 `false` |
| 无效 UTF-8 | 替换为 U+FFFD | 报错 |
| `nil` slice / `nil` map | `null` | `[]` / `{}` |

唯一需要改代码的是第一行：结构体字段没写 `json` tag 而 JSON 里是小写名时，v2 不会匹配上，字段保持零值且 `Unmarshal` 返回 `nil`。给字段补上 `json:"name"` tag 即可，这在所有后端下都正确。需要让拼错的成员名暴露出来时，用 `Decoder.DisallowUnknownFields()`。

### 与其他后端保持一致的三处

`encoding/json/v2` 的下列默认值会让本后端成为五个后端里唯一行为不同的那个，因此显式覆盖，与默认后端及 `sonic` / `go_json` / `jsoniter` 逐字节一致：

| 行为 | v2 单独的默认值 | 本后端 |
|------|----------------|--------|
| map 键顺序 | Go 随机化迭代顺序，同一个值每次输出不同字节 | 按键排序（`Deterministic`） |
| `time.Duration` | 报错 `no default representation`，编解码都失败 | 纳秒整数（`FormatDurationAsNano`） |
| `omitempty` | 只省略编码为 `null` / `""` / `{}` / `[]` 的值，`0` 与 `false` 保留 | 零值一并省略（`OmitEmptyWithLegacySemantics`） |

前两项的性质值得说明：map 顺序随机会让签名、缓存键、ETag、快照比对**间歇性**失效（实测同一个 8 键 map 连续 marshal 300 次产生 8 种字节序列）；`time.Duration` 则是直接的运行时失败——任何含 `Duration` 字段的结构体都无法编解码。`omitempty` 一项若不覆盖，历史上一直省略零值的响应体会突然多出 `"n":0`、`"b":false` 这类字段。

对 `omitempty` 语义本身有要求时，用 v2 推荐的 `omitzero` tag：它在所有后端下都表示"Go 零值则省略"。

### 关于性能

这个后端不是性能选项，但原因不是 v2 慢。直接调用标准库两个包对比，`encoding/json/v2` 的 `Marshal` 比 `encoding/json` 快 10~15%、`Unmarshal` 快 5~20%，分配次数相同——`encoding/json` 每次都要处理一整套 v1 兼容选项，这就是差距来源。

本后端把这点优势花掉了：上一节列出的三项对齐各有成本（`Deterministic` 要排序 map 键、另两个是额外的选项处理），`MarshalIndent` 与 `SetIndent` 又多一次缩进拷贝。本机 `go test -bench=. -count=3` 实测的净结果是结构体 `Marshal` 比默认后端慢约 13%、`Valid` 慢约 40%（v2 的校验包含重复键检测），`Unmarshal` 与小对象 `Marshal` 基本持平。

选它的理由是 RFC 7493 严格语义与零外部依赖；要吞吐用 `sonic` 或 `go_json`。

## 标准库 struct tag 与接口的后端支持

`omitzero` 是 `encoding/json` 自身的 tag；`case:ignore` / `case:strict` / `embed` 与 `MarshalerTo` / `UnmarshalerFrom` 接口来自 `encoding/json/v2`，Go 1.27 起在 `encoding/json` 下同样生效（v1 的实现已改用 v2 引擎）。

这些能力在默认后端、`jsonv2` 与 `sonic` 下均可用。`go_json` 与 `jsoniter` 各自实现了独立的 tag 解析器，按自己的规则处理这些字段。实测于 Go 1.27：

| 能力 | _(默认)_ | `jsonv2` | `sonic` | `go_json` | `jsoniter` |
|------|:--------:|:--------:|:-------:|:---------:|:----------:|
| `omitzero` | ✅ | ✅ | ✅ | — | — |
| `omitzero` 认 `IsZero() bool` 方法 | ✅ | ✅ | ✅ | — | — |
| `case:strict` | ✅ | ✅ | ✅ | — | — |
| `case:ignore` | ✅ | ✅ | ✅ | 注 | 注 |
| `embed` 兜底字段（`jsontext.Value` 收集未匹配成员） | ✅ | ✅ | ✅ | — | — |
| `MarshalerTo` / `UnmarshalerFrom`（`MarshalJSONTo` / `UnmarshalJSONFrom`） | ✅ | ✅ | ✅ | — | — |

标 — 的格子表示该 tag 或接口不参与决策，字段回到默认处理：`omitzero` 字段照常出现在输出里（如 `{"ns":null,"z":0}`）、`case:strict` 不改变成员名匹配、`embed` 兜底字段收不到未匹配成员、`MarshalJSONTo` 不被调用而走反射的默认表示。

标「注」的两格结果与左侧相同，但源于 `go_json` 与 `jsoniter` 默认就以大小写不敏感的方式匹配成员名，并非 `case:ignore` 生效。

依赖上表任一能力时，请选择默认后端、`jsonv2` 或 `sonic`。

## 后端特有功能

### sonic: SetFastest

```go
import "github.com/gtkit/json/v2"

func main() {
    // 切换到最高性能模式（禁用 map key 排序、HTML 转义等兼容特性）
    json.SetFastest()

    // 后续所有 Marshal/Unmarshal 使用 ConfigFastest
    data, _ := json.Marshal(myStruct)
}
```

> 注意：`SetFastest()` 仅在 `sonic` build tag 下可用。应在 `main()` 初始化阶段调用，不要在并发环境中动态切换。

### jsoniter: 私有字段 & PHP 兼容模式

```go
import "github.com/gtkit/json/v2"

func main() {
    // 启用非导出字段的编解码
    json.SupportPrivateFields()

    // 启用 PHP 兼容的模糊解码（"123" → int, "true" → bool）
    json.RegisterFuzzyDecoders()
}
```

> 注意：这两个函数仅在 `jsoniter` build tag 下可用。同样应在初始化阶段调用。

## 在 Gin 中使用

将 Gin 的 JSON 编解码切换为本库的后端，只需在构建时加上对应 tag：

```bash
go build -tags=sonic -o server ./cmd/server
```

如果你自定义了 Gin 的 JSON codec，可以直接使用 `json.API`：

```go
import "github.com/gtkit/json/v2"

// 直接使用顶层函数
data, err := json.Marshal(gin.H{"code": 0, "msg": "ok"})
```

## 测试 mock 示例

v2 的 interface 设计使得测试 mock 变得简单：

```go
type mockJSON struct{}

func (mockJSON) Marshal(v any) ([]byte, error)                              { return []byte(`{}`), nil }
func (mockJSON) Unmarshal(data []byte, v any) error                         { return nil }
func (mockJSON) MarshalIndent(v any, prefix, indent string) ([]byte, error) { return []byte(`{}`), nil }
func (mockJSON) MarshalToString(v any) (string, error)                      { return "{}", nil }
func (mockJSON) NewEncoder(w io.Writer) json.Encoder                        { return nil }
func (mockJSON) NewDecoder(r io.Reader) json.Decoder                        { return nil }
func (mockJSON) Valid(data []byte) bool                                     { return true }

func TestWithMockJSON(t *testing.T) {
    original := json.API
    json.API = mockJSON{}
    t.Cleanup(func() { json.API = original })

    data, err := json.Marshal("anything")
    // data == []byte(`{}`)
}
```

## Benchmark

```bash
# 标准库
go test -bench=. -benchmem

# encoding/json/v2
go test -bench=. -benchmem -tags=jsonv2

# sonic
go test -bench=. -benchmem -tags=sonic

# go-json
go test -bench=. -benchmem -tags=go_json

# jsoniter
go test -bench=. -benchmem -tags=jsoniter
```

## 从 v1 迁移

v2 相比 v1 的变化：

| 变更项 | v1 | v2 |
|--------|----|----|
| 导入路径 | `github.com/gtkit/json` | `github.com/gtkit/json/v2` |
| 调用方式 | `json.Marshal(v)` | `json.Marshal(v)`（不变） |
| 底层机制 | 包级函数变量 | `Core` interface + 顶层便捷函数 |
| 新增 API | — | `MarshalToString`、`Valid` |
| 可测试性 | 需要替换包级变量 | 替换 `json.API` 即可 mock |
| `CheckJSON()` | 打印日志 | 改用 `json.Package` 常量 |
| `SetFastest()` | 所有后端都暴露（空实现） | 仅 sonic 后端可用 |

迁移步骤：

1. 更新 import 路径为 `github.com/gtkit/json/v2`
2. `CheckJSON()` 改为读取 `json.Package`
3. 其他调用方式完全兼容，无需修改

## License

MIT，见 [../LICENSE](../LICENSE)。
