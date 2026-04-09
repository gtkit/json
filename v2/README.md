# gtkit/json/v2

[![Go Reference](https://pkg.go.dev/badge/github.com/gtkit/json/v2.svg)](https://pkg.go.dev/github.com/gtkit/json/v2)
[![Go Version](https://img.shields.io/badge/go-%3E%3D1.26-blue)](https://go.dev)

通过 build tags 无缝切换 JSON 编解码后端的 Go 库。零修改业务代码，一个 `-tags` 参数即可获得数倍性能提升。

## 支持的后端

| Build Tag | 后端库 | 适用场景 |
|-----------|--------|---------|
| _(默认)_ | `encoding/json` | 零依赖、最大兼容性 |
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

如果你的目标平台不在此列表中，建议使用 `go_json` 作为替代：

```bash
GOOS=freebsd GOARCH=amd64 go build -tags=go_json ./...
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

### 常量

```go
json.Package  // 当前后端库名，如 "encoding/json"、"github.com/bytedance/sonic"
json.Version  // 包版本号，如 "v2.0.0"
```

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

func (mockJSON) Marshal(v any) ([]byte, error)                          { return []byte(`{}`), nil }
func (mockJSON) Unmarshal(data []byte, v any) error                     { return nil }
func (mockJSON) MarshalIndent(v any, prefix, indent string) ([]byte, error) { return []byte(`{}`), nil }
func (mockJSON) MarshalToString(v any) (string, error)                  { return "{}", nil }
func (mockJSON) NewEncoder(w io.Writer) json.Encoder                    { return nil }
func (mockJSON) NewDecoder(r io.Reader) json.Decoder                    { return nil }
func (mockJSON) Valid(data []byte) bool                                 { return true }

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
