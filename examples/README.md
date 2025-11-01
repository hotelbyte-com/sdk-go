# HotelByte SDK 使用示例

本目录包含了使用 HotelByte SDK 的各种示例代码，帮助您快速上手和学习最佳实践。

## 示例列表

### 基础示例

1. **[快速开始](./quickstart/)** - 最简单的 SDK 使用示例
   - 客户端初始化
   - 基本的酒店搜索
   - 简单的错误处理

2. **[认证管理](./authentication/)** - 认证相关示例
   - 手动认证
   - 令牌管理
   - 认证错误处理

3. **[酒店搜索](./hotel-search/)** - 酒店搜索功能示例
   - 基本搜索
   - 高级搜索参数
   - 分页处理
   - 搜索结果处理

### 高级示例

4. **[预订管理](./booking-management/)** - 预订相关操作
   - 创建预订
   - 查询预订
   - 取消预订
   - 预订状态追踪

5. **[错误处理](./error-handling/)** - 完整的错误处理策略
   - 错误类型识别
   - 重试机制
   - 降级策略
   - 错误日志记录

6. **[并发处理](./concurrent-processing/)** - 并发请求处理
   - 并发搜索
   - 批量处理
   - 限流控制
   - 结果聚合

### 性能优化示例

7. **[缓存策略](./caching/)** - 缓存实现示例
   - 内存缓存
   - Redis 缓存
   - 缓存失效策略
   - 缓存穿透防护

8. **[性能监控](./monitoring/)** - 性能监控和指标收集
   - 请求耗时统计
   - 错误率监控
   - 吞吐量统计
   - 健康检查

### 实际应用示例

9. **[Web 服务集成](./web-service/)** - 在 Web 服务中使用 SDK
   - HTTP API 封装
   - 中间件集成
   - 请求链路追踪
   - 响应格式化

10. **[微服务架构](./microservice/)** - 微服务环境下的使用
    - 服务发现集成
    - 配置管理
    - 服务间通信
    - 分布式追踪

11. **[命令行工具](./cli-tool/)** - 命令行工具示例
    - 参数解析
    - 配置文件支持
    - 交互式操作
    - 结果输出格式化

## 如何使用这些示例

### 1. 环境准备

首先确保您有有效的 HotelByte API 凭据：

```bash
export HOTELBYTE_APP_KEY="your-app-key"
export HOTELBYTE_APP_SECRET="your-app-secret"
export HOTELBYTE_BASE_URL="https://api.hotelbyte.com"
```

### 2. 安装依赖

每个示例目录都包含一个 `go.mod` 文件，运行以下命令安装依赖：

```bash
cd examples/quickstart
go mod download
```

### 3. 运行示例

```bash
go run main.go
```

### 4. 查看代码

每个示例都包含详细的注释，解释了代码的作用和最佳实践。

## 示例结构

每个示例目录通常包含：

```
example-name/
├── main.go              # 主程序入口
├── go.mod              # Go 模块文件
├── go.sum              # 依赖锁定文件
├── README.md           # 示例说明
├── config.yaml         # 配置文件（如果需要）
└── internal/           # 内部包（如果需要）
    ├── handler/        # 处理器
    ├── service/        # 服务层
    └── model/          # 数据模型
```

## 学习路径建议

### 初学者
1. 从 **快速开始** 示例开始
2. 学习 **认证管理** 和 **错误处理**
3. 掌握 **酒店搜索** 和 **预订管理**

### 进阶用户
1. 学习 **并发处理** 和 **缓存策略**
2. 了解 **性能监控** 和优化技巧
3. 参考 **Web 服务集成** 实现自己的服务

### 高级用户
1. 研究 **微服务架构** 示例
2. 自定义扩展和优化
3. 贡献更多示例

## 常见问题

### Q: 示例运行失败怎么办？

A: 请检查：
1. 环境变量是否正确设置
2. 网络连接是否正常
3. API 凭据是否有效
4. Go 版本是否兼容（建议 1.19+）

### Q: 如何修改示例适应自己的需求？

A: 每个示例都是独立的，您可以：
1. 复制示例代码到您的项目
2. 根据注释修改相关参数
3. 添加您自己的业务逻辑

### Q: 示例中的最佳实践是强制的吗？

A: 不是强制的，但强烈建议遵循，因为它们是基于实际生产经验总结的。

## 贡献示例

我们欢迎社区贡献更多有用的示例！如果您有好的使用案例，请：

1. Fork 项目仓库
2. 在 `examples/` 目录下创建新的示例
3. 确保代码质量和文档完整
4. 提交 Pull Request

### 示例贡献指南

1. **命名规范**：使用清晰描述性的目录名
2. **代码质量**：确保代码可以正常运行
3. **文档完整**：包含详细的 README 和代码注释
4. **最佳实践**：体现 SDK 的最佳使用方式
5. **错误处理**：包含完善的错误处理逻辑

## 相关文档

- [SDK 文档](../README.md)
- [API 参考](../api-reference.md)
- [最佳实践](../best-practices.md)
- [配置指南](../configuration.md)
- [故障排除](../troubleshooting.md)

## 支持

如果您在使用示例时遇到问题：

- 📧 Email: developers@hotelbyte.com
- 🐛 GitHub Issues: https://github.com/hotelbyte/hotel-be/issues
- 💬 讨论区: https://github.com/hotelbyte/hotel-be/discussions
- 📖 文档: https://docs.hotelbyte.com