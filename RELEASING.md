# Releasing

这个仓库里有两个 Go modules：

- 根模块：`github.com/gtkit/json`
- `v2` 子模块：`github.com/gtkit/json/v2`

## Tag 规则

- 根模块继续使用普通语义化版本 tag：`v0.x.y`，以及未来可能的 `v1.x.y`。
- `v2` 虽然放在仓库根目录下的 `v2/` 子目录，但它属于 Go 官方定义的 major version subdirectory。
- 因此，`v2` 的正确 tag 形态是 `v2.x.y`，不是 `v2/v2.x.y`。

参考 Go 官方模块文档：

- https://go.dev/ref/mod

官方规则的关键点是：

- 模块如果定义在仓库根目录，或者定义在仓库根目录下的 major version 子目录（例如 `v2/`），版本 tag 名称直接等于版本号本身。
- 只有普通子目录模块才需要使用 `subdir/vX.Y.Z` 这种带目录前缀的 tag。

## 发布 v2

在 `v2/` 目录执行：

```bash
make tag
```

这个目标会：

1. 根据 `v2/version.go` 自动递增 patch 版本
2. 提交版本变更
3. 推送当前分支
4. 创建并推送正确的 `v2.x.y` tag

查看当前最新 `v2` tag：

```bash
make gittag
```

## 历史说明

仓库里已经出现过 `v2/v2.0.0`、`v2/v2.0.1`、`v2/v2.0.2` 这类错误 tag 形式。

对于当前这个仓库布局，Go 工具链识别的是 `v2.0.2` 这种不带目录前缀的 tag。后续发布不要再继续使用 `v2/v2.x.y`。

如需补齐历史版本，可基于现有错误 tag 重新打出正确 tag，例如：

```bash
git tag -a v2.0.0 'v2/v2.0.0^{}' -m "release v2.0.0"
git tag -a v2.0.1 'v2/v2.0.1^{}' -m "release v2.0.1"
git push gtkit v2.0.0 v2.0.1
```
