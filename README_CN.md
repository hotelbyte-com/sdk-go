# HotelByte Go SDK

[![Go Version](https://img.shields.io/badge/Go-1.21+-blue.svg)](https://golang.org)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![Documentation](https://img.shields.io/badge/Documentation-GoDoc-blue.svg)](https://pkg.go.dev/github.com/hotelbyte-com/sdk-go)

HotelByte Go SDK 是酒店 API 分销平台的官方 Go 语言软件开发工具包，提供简单、高效、类型安全的 API 接口访问能力。

## ✨ 核心特性

- **🔗 统一接口**：聚合多家酒店供应商（HotelBeds、Dida、DerbySoft 等），提供一致的 API 体验
- **🚀 高性能**：基于 Go 语言的并发特性，支持高并发请求处理
- **🛡️ 类型安全**：完整的类型定义，编译时错误检查
- **🔄 自动重试**：内置智能重试机制，自动处理网络异常和临时故障
- **📝 完善文档**：详细的 API 文档和使用示例
- **🧪 测试覆盖**：高测试覆盖率，保证 SDK 稳定性
- **🌐 多语言支持**：支持中英文国际化

## 📦 安装

使用 Go modules 安装 SDK：

```bash
go get github.com/hotelbyte-com/sdk-go
```

在你的代码中导入：

```go
import "github.com/hotelbyte-com/sdk-go"
```

## 🚀 快速开始

### 基础配置

```go
package main

import (
    "context"
    "fmt"
    "time"

    "github.com/hotelbyte-com/sdk-go"
)

func main() {
    // 创建客户端
    client, err := hotelbyte.NewClient(
        hotelbyte.WithCredentials("your-app-key", "your-app-secret"),
        hotelbyte.WithBaseURL("https://api.hotelbyte.com"),
        hotelbyte.WithTimeout(30*time.Second),
    )
    if err != nil {
        panic(err)
    }
    defer client.Close()

    // 搜索酒店
    ctx := context.Background()
    resp, err := client.SearchHotels(ctx, &hotelbyte.SearchHotelsRequest{
        DestinationID: "beijing",
        CheckIn:       time.Now().AddDate(0, 0, 7),
        CheckOut:      time.Now().AddDate(0, 0, 9),
        AdultCount:    2,
        RoomCount:     1,
    })
    if err != nil {
        panic(err)
    }

    fmt.Printf("找到 %d 家酒店\n", len(resp.Hotels))
    for _, hotel := range resp.Hotels {
        fmt.Printf("- %s: %s\n", hotel.Name, hotel.Address)
    }
}
```

### 环境变量配置

```bash
export HOTELBYTE_APP_KEY="your-app-key"
export HOTELBYTE_APP_SECRET="your-app-secret"
export HOTELBYTE_BASE_URL="https://api.hotelbyte.com"
```

```go
// 使用环境变量创建客户端
func NewClientFromEnv() (*hotelbyte.Client, error) {
    return hotelbyte.NewClient(
        hotelbyte.WithCredentials(
            os.Getenv("HOTELBYTE_APP_KEY"),
            os.Getenv("HOTELBYTE_APP_SECRET"),
        ),
        hotelbyte.WithBaseURL(os.Getenv("HOTELBYTE_BASE_URL")),
    )
}
```

## 📖 主要功能

### 🔍 酒店搜索

```go
// 基础搜索
resp, err := client.SearchHotels(ctx, &hotelbyte.SearchHotelsRequest{
    DestinationID: "beijing",
    CheckIn:       time.Now().AddDate(0, 0, 7),
    CheckOut:      time.Now().AddDate(0, 0, 9),
    AdultCount:    2,
    RoomCount:     1,
})

// 高级搜索参数
resp, err := client.SearchHotels(ctx, &hotelbyte.SearchHotelsRequest{
    DestinationID: "beijing",
    CheckIn:       time.Now().AddDate(0, 0, 7),
    CheckOut:      time.Now().AddDate(0, 0, 9),
    AdultCount:    2,
    RoomCount:     1,
    HotelIDs:      []string{"hotel-1", "hotel-2"},
    MinRating:     4,
    MaxPrice:      1000.00,
    Currency:      "CNY",
    Language:      "zh-CN",
})
```

### 🏨 酒店详情

```go
// 获取酒店详细信息
details, err := client.GetHotelDetails(ctx, "hotel-id")
if err != nil {
    return err
}

fmt.Printf("酒店名称: %s\n", details.Name)
fmt.Printf("酒店地址: %s\n", details.Address)
fmt.Printf("酒店设施: %v\n", details.Facilities)
```

### 📋 预订管理

```go
// 创建预订
booking, err := client.CreateBooking(ctx, &hotelbyte.CreateBookingRequest{
    HotelID:    "hotel-id",
    RoomTypeID: "room-type-id",
    CheckIn:    time.Now().AddDate(0, 0, 7),
    CheckOut:   time.Now().AddDate(0, 0, 9),
    GuestInfo: hotelbyte.GuestInfo{
        Name:    "张三",
        Email:   "zhangsan@example.com",
        Phone:   "+86-13800138000",
    },
    PaymentInfo: hotelbyte.PaymentInfo{
        Method: "credit_card",
        Card:   "****-****-****-1234",
    },
})

// 查询预订状态
status, err := client.GetBookingStatus(ctx, booking.BookingID)

// 取消预订
err = client.CancelBooking(ctx, booking.BookingID)
```

### 🔐 认证管理

```go
// SDK 会自动处理认证，但也可以手动管理
auth := client.Auth()

// 检查认证状态
if auth.IsTokenValid() {
    fmt.Println("认证有效")
}

// 手动刷新认证
err := auth.RefreshToken(ctx)
if err != nil {
    fmt.Printf("认证刷新失败: %v\n", err)
}
```

## ⚙️ 客户端配置

### 基础配置选项

```go
client, err := hotelbyte.NewClient(
    // 认证信息
    hotelbyte.WithCredentials("app-key", "app-secret"),

    // 基础配置
    hotelbyte.WithBaseURL("https://api.hotelbyte.com"),
    hotelbyte.WithTimeout(30 * time.Second),
    hotelbyte.WithUserAgent("my-app/1.0"),

    // 重试配置
    hotelbyte.WithRetryConfig(3, time.Second, 30*time.Second),

    // HTTP 配置
    hotelbyte.WithHTTPConfig(hotelbyte.HTTPConfig{
        MaxIdleConns:    100,
        MaxConnsPerHost: 10,
        IdleConnTimeout: 90 * time.Second,
    }),

    // 调试模式
    hotelbyte.WithDebug(true),
)
```

### 高级配置

```go
// 自定义 HTTP 客户端
httpClient := &http.Client{
    Timeout: 30 * time.Second,
    Transport: &http.Transport{
        MaxIdleConns:        100,
        MaxIdleConnsPerHost: 10,
        IdleConnTimeout:     90 * time.Second,
    },
}

client, err := hotelbyte.NewClient(
    hotelbyte.WithHTTPClient(httpClient),
    hotelbyte.WithCredentials("app-key", "app-secret"),
)

// 自定义日志记录器
logger := logrus.New()
logger.SetLevel(logrus.InfoLevel)

client, err := hotelbyte.NewClient(
    hotelbyte.WithLogger(logger),
    hotelbyte.WithCredentials("app-key", "app-secret"),
)
```

## 🔧 错误处理

### 错误类型检查

```go
resp, err := client.SearchHotels(ctx, req)
if err != nil {
    switch {
    case hotelbyte.IsAuthenticationError(err):
        fmt.Println("认证失败，请检查凭据")
    case hotelbyte.IsValidationError(err):
        fmt.Println("请求参数验证失败")
    case hotelbyte.IsRateLimitError(err):
        fmt.Println("请求频率过高，请稍后重试")
    case hotelbyte.IsSystemError(err):
        fmt.Println("系统错误，请稍后重试")
    case hotelbyte.IsTimeoutError(err):
        fmt.Println("请求超时")
    default:
        fmt.Printf("未知错误: %v\n", err)
    }
    return
}
```

### 自定义重试策略

```go
// 实现自定义重试逻辑
func searchWithRetry(ctx context.Context, client *hotelbyte.Client, req *hotelbyte.SearchHotelsRequest) (*hotelbyte.SearchHotelsResponse, error) {
    var lastErr error

    for attempt := 0; attempt < 3; attempt++ {
        if attempt > 0 {
            // 指数退避
            delay := time.Duration(1<<uint(attempt-1)) * time.Second
            select {
            case <-ctx.Done():
                return nil, ctx.Err()
            case <-time.After(delay):
            }
        }

        resp, err := client.SearchHotels(ctx, req)
        if err == nil {
            return resp, nil
        }

        lastErr = err

        // 判断是否可重试
        if !hotelbyte.IsRateLimitError(err) && !hotelbyte.IsSystemError(err) {
            break
        }
    }

    return nil, fmt.Errorf("搜索失败，已重试3次: %w", lastErr)
}
```

## 🚀 并发处理

### 并发搜索

```go
func concurrentSearch(ctx context.Context, client *hotelbyte.Client, destinations []string) map[string]*hotelbyte.SearchHotelsResponse {
    results := make(map[string]*hotelbyte.SearchHotelsResponse)
    var mu sync.Mutex
    var wg sync.WaitGroup

    // 限制并发数
    semaphore := make(chan struct{}, 5)

    for _, dest := range destinations {
        wg.Add(1)
        go func(destination string) {
            defer wg.Done()

            // 获取信号量
            semaphore <- struct{}{}
            defer func() { <-semaphore }()

            req := &hotelbyte.SearchHotelsRequest{
                DestinationID: destination,
                CheckIn:       time.Now().AddDate(0, 0, 7),
                CheckOut:      time.Now().AddDate(0, 0, 9),
                AdultCount:    2,
                RoomCount:     1,
            }

            resp, err := client.SearchHotels(ctx, req)
            if err == nil {
                mu.Lock()
                results[destination] = resp
                mu.Unlock()
            }
        }(dest)
    }

    wg.Wait()
    return results
}
```

## 📊 性能优化

### 缓存策略

```go
type CachedSearcher struct {
    client *hotelbyte.Client
    cache  *sync.Map // 生产环境建议使用 Redis
}

func NewCachedSearcher(client *hotelbyte.Client) *CachedSearcher {
    return &CachedSearcher{
        client: client,
        cache:  &sync.Map{},
    }
}

func (cs *CachedSearcher) SearchHotels(ctx context.Context, req *hotelbyte.SearchHotelsRequest) (*hotelbyte.SearchHotelsResponse, error) {
    // 生成缓存键
    key := cs.generateCacheKey(req)

    // 检查缓存
    if cached, ok := cs.cache.Load(key); ok {
        if entry, ok := cached.(*cacheEntry); ok && !entry.isExpired() {
            return entry.data, nil
        }
    }

    // 缓存未命中，调用 API
    resp, err := cs.client.SearchHotels(ctx, req)
    if err != nil {
        return nil, err
    }

    // 缓存结果
    cs.cache.Store(key, &cacheEntry{
        data:      resp,
        expiresAt: time.Now().Add(10 * time.Minute),
    })

    return resp, nil
}

type cacheEntry struct {
    data      *hotelbyte.SearchHotelsResponse
    expiresAt time.Time
}

func (e *cacheEntry) isExpired() bool {
    return time.Now().After(e.expiresAt)
}
```

## 📝 示例代码

查看 [examples](./examples) 目录获取更多详细示例：

- **[快速开始](./examples/quickstart/)** - 最简单的使用示例
- **[认证管理](./examples/authentication/)** - 认证相关操作
- **[酒店搜索](./examples/hotel-search/)** - 搜索功能示例
- **[预订管理](./examples/booking-management/)** - 预订流程示例
- **[错误处理](./examples/error-handling/)** - 完整的错误处理策略
- **[并发处理](./examples/concurrent-processing/)** - 并发请求处理
- **[性能优化](./examples/caching/)** - 缓存和性能优化
- **[监控集成](./examples/monitoring/)** - 监控和指标收集

## 🔧 最佳实践

### 1. 客户端管理

```go
// ✅ 推荐：全局单例模式
var hotelClient *hotelbyte.Client

func init() {
    var err error
    hotelClient, err = hotelbyte.NewClient(
        hotelbyte.WithCredentials(
            os.Getenv("HOTELBYTE_APP_KEY"),
            os.Getenv("HOTELBYTE_APP_SECRET"),
        ),
        hotelbyte.WithTimeout(30*time.Second),
        hotelbyte.WithRetryConfig(3, time.Second, 30*time.Second),
    )
    if err != nil {
        log.Fatal("Failed to initialize HotelByte client:", err)
    }
}

func GetHotelClient() *hotelbyte.Client {
    return hotelClient
}
```

### 2. 环境变量配置

```go
type Config struct {
    AppKey        string        `envconfig:"HOTELBYTE_APP_KEY" required:"true"`
    AppSecret     string        `envconfig:"HOTELBYTE_APP_SECRET" required:"true"`
    BaseURL       string        `envconfig:"HOTELBYTE_BASE_URL" default:"https://api.hotelbyte.com"`
    Timeout       time.Duration `envconfig:"HOTELBYTE_TIMEOUT" default:"30s"`
    MaxRetries    int           `envconfig:"HOTELBYTE_MAX_RETRIES" default:"3"`
    Debug         bool          `envconfig:"HOTELBYTE_DEBUG" default:"false"`
}

func NewHotelClient() (*hotelbyte.Client, error) {
    var cfg Config
    if err := envconfig.Process("", &cfg); err != nil {
        return nil, err
    }

    return hotelbyte.NewClient(
        hotelbyte.WithCredentials(cfg.AppKey, cfg.AppSecret),
        hotelbyte.WithBaseURL(cfg.BaseURL),
        hotelbyte.WithTimeout(cfg.Timeout),
        hotelbyte.WithRetryConfig(cfg.MaxRetries, time.Second, 30*time.Second),
    )
}
```

### 3. 结构化日志

```go
func loggedSearch(ctx context.Context, client *hotelbyte.Client, req *hotelbyte.SearchHotelsRequest) (*hotelbyte.SearchHotelsResponse, error) {
    logger := logrus.WithFields(logrus.Fields{
        "operation":     "hotel_search",
        "destination":   req.DestinationID,
        "check_in":      req.CheckIn.Format("2006-01-02"),
        "check_out":     req.CheckOut.Format("2006-01-02"),
        "adult_count":   req.AdultCount,
        "room_count":    req.RoomCount,
    })

    logger.Info("开始搜索酒店")
    start := time.Now()

    resp, err := client.SearchHotels(ctx, req)
    duration := time.Since(start)

    if err != nil {
        logger.WithFields(logrus.Fields{
            "error":    err.Error(),
            "duration": duration,
        }).Error("酒店搜索失败")
        return nil, err
    }

    logger.WithFields(logrus.Fields{
        "hotel_count": len(resp.Hotels),
        "total":       resp.Total,
        "duration":    duration,
    }).Info("酒店搜索成功")

    return resp, nil
}
```

## 📚 文档

- [API 参考](./api-reference.md) - 完整的 API 文档
- [配置指南](./configuration.md) - 配置选项详解
- [故障排除](./troubleshooting.md) - 常见问题解决方案
- [示例代码](./examples/README.md) - 综合使用示例

## 🧪 测试

运行测试：

```bash
# 运行所有测试
go test ./...

# 运行测试并生成覆盖率报告
go test -cover ./...

# 运行基准测试
go test -bench=. ./...

# 运行特定包的测试
go test ./protocol/...
```

## 🤝 贡献

我们欢迎社区贡献！请查看 [贡献指南](./CONTRIBUTING.md) 了解如何参与项目开发。

### 贡献流程

1. Fork 项目仓库
2. 创建功能分支 (`git checkout -b feature/amazing-feature`)
3. 提交更改 (`git commit -m 'Add some amazing feature'`)
4. 推送到分支 (`git push origin feature/amazing-feature`)
5. 创建 Pull Request

### 代码规范

- 遵循 Go 官方代码规范
- 添加适当的测试用例
- 更新相关文档
- 确保所有测试通过

## 📄 许可证

本项目采用 MIT 许可证 - 查看 [LICENSE](LICENSE) 文件了解详情。

## 🆘 支持

如果您在使用过程中遇到问题：

- 📧 **邮件支持**: developers@hotelbyte.com
- 🐛 **问题反馈**: [GitHub Issues](https://github.com/hotelbyte-com/hotel-be/issues)
- 💬 **讨论交流**: [GitHub Discussions](https://github.com/hotelbyte-com/hotel-be/discussions)
- 📖 **项目文档**: [HotelByte 技术文档](https://docs.hotelbyte.com)

## 🌟 致谢

感谢所有为 HotelByte SDK 做出贡献的开发者和用户！

---

**注意**: 本 SDK 是 HotelByte 酒店API分销平台的官方 Go 语言实现。如果您对其他语言的 SDK 有需求，请联系我们或查看我们的其他项目。