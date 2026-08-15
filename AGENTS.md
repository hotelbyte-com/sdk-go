# HotelByte Go SDK Agent Guide

本文件覆盖整个独立仓库。不要假设调用方同时检出了 `hotel-be`，也不要依赖父仓库命令或约定。

## 职责

- 维护公开 Go module `github.com/hotelbyte-com/sdk-go`：客户端配置、认证、传输、酒店与房型映射 API、`protocol/` 的线协议类型以及可运行示例。
- 保持现有消费者可升级：公开标识符、方法签名、JSON 字段、端点路径、日期/ID/金额编码和错误语义都是兼容性边界。
- 完成变更后在 feature branch 提交并 push；创建 tag、GitHub Release 或其他发布动作必须有明确发布任务。

## 先读入口

1. `go.mod`：模块路径、实际 Go 版本和依赖真相。
2. `client.go`、`auth.go`、`transport.go`：配置、认证、重试和 HTTP 生命周期。
3. `hotels.go`、`room_mapping.go`：公开 API 与端点绑定。
4. `protocol/` 及相邻 `*_test.go`：请求/响应线协议和兼容性样例。
5. `examples/quickstart/main.go`：真实调用流程。README 示例可能滞后，修改前要与导出 API 对照。

## 实现与兼容性边界

- 使用 `go.mod` 声明的 Go 版本；Go 文件必须 `gofmt`。不要仅为文档声称的旧版本降低语言或标准库用法。
- 协议类型尽量增量演进。删除或重命名公开字段、修改 JSON 表示、端点或默认 header/retry 行为前，必须给出迁移路径和覆盖旧行为的测试。
- 认证刷新和自动重试必须保持并发安全、有界且可观察；不得把真实错误转成空成功，也不得无条件重放可能产生副作用的请求。
- 金额与币种必须来自同一个响应对象；不要在 SDK 中猜测、补零或跨对象拼接财务值。
- `room_mapping_live_test.go` 只有显式提供 `HOTELBYTE_ROOM_MAPPING_LIVE_*` 环境变量时才访问网络。默认验证不得设置这些变量，也不得运行会创建/取消预订的 quickstart。
- 不新增真实凭据、token、客户数据或未脱敏请求/响应到源码、测试、日志和文档。

## 最小验证

- 仅 Markdown/指导变更：`git diff --check`，并人工核对示例名称与当前导出 API。
- Go 源码、协议或单元测试变更：先运行相关 package 测试，再运行 `go test ./...`。
- 认证、重试或并发变更：额外运行 `go test -race ./...`。
- live smoke 只在任务明确要求、目标环境与凭据已授权时运行；单元测试通过不等于 live API 已验证。

## Code Review Rules

- 阻止无迁移路径的公开 API、JSON wire contract、端点路径或默认行为破坏；安全路径是保留旧契约或增加向后兼容入口并测试两条路径。
- 阻止认证/重试改动造成数据竞争、无限重试、header 丢失、错误吞噬或副作用请求的隐式重复执行。
- 阻止金额与币种跨来源组合、空值伪装为零，以及文档宣称不存在或未验证的 SDK 能力。
- 阻止测试依赖公网、时间偶然性或真实凭据；默认测试应使用本地 fake/`httptest`，live 测试必须显式 opt-in。
