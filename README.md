# gtkit/json

通过 build tags 无缝切换 JSON 编解码后端的 Go 库。

## 模块版本

### v1

```bash
go get github.com/gtkit/json@latest
```

```go
import "github.com/gtkit/json"
```

### v2

`v2` 是独立 Go module，引用路径必须带 `/v2`。

```bash
go get github.com/gtkit/json/v2@latest
```

指定版本时：

```bash
go get github.com/gtkit/json/v2@v2.0.2
```

```go
import "github.com/gtkit/json/v2"
```

更完整的 `v2` 文档见 [v2/README.md](v2/README.md)。
维护者发布说明见 [RELEASING.md](RELEASING.md)。

## 默认使用 encoding/json

```go
CheckJSON() // 检查当前实际使用的 JSON 库
```

## Build jsoniter

```bash
go build -tags=jsoniter -ldflags "-s -w" -gcflags="-m" -o app main.go
```

## Build go-json

```bash
go build -tags=go_json -ldflags "-s -w" -gcflags="-m" -o app main.go
```

## Build sonic

```bash
go build -tags="sonic,avx,darwin,amd64" -ldflags "-s -w" -gcflags="-m" -o app main.go
```

`v2` 的 build tags 用法相同，只是 import 路径改为 `github.com/gtkit/json/v2`。

## License

MIT，见 [LICENSE](LICENSE)。
