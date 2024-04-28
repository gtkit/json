### 默认使用 "encoding/json", 运行
```
CheckJSON() // 检查json是否符合规范
TrackTime() // 跟踪json解析时间
```

### Build jsoniter json
```
go build -tags=jsoniter -ldflags "-s -w" -gcflags="-m"  -o app main.go
```

### Build go-json
```
go build -tags=go_json -ldflags "-s -w" -gcflags="-m"  -o app main.go
```

### Build sonic json
```
go build -tags="sonic,avx,darwin,amd64" -ldflags "-s -w" -gcflags="-m"  -o app main.go
```
