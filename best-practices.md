# HotelByte SDK 最佳实践指南

本指南提供了使用 HotelByte SDK 的最佳实践，帮助您构建高效、可靠、可维护的应用程序。

## 1. 客户端配置和管理

### 1.1 客户端实例化

**推荐做法：应用全局使用单一客户端实例**

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
        hotelbyte.WithBaseURL(os.Getenv("HOTELBYTE_BASE_URL")),
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

**避免：为每个请求创建新客户端**

```go
// ❌ 避免：频繁创建客户端
func searchHotels() {
    client, _ := hotelbyte.NewClient(...)  // 资源浪费
    defer client.Close()
    // ...
}
```

### 1.2 配置最佳实践

**使用环境变量管理配置**

```go
type Config struct {
    AppKey        string `envconfig:"HOTELBYTE_APP_KEY" required:"true"`
    AppSecret     string `envconfig:"HOTELBYTE_APP_SECRET" required:"true"`
    BaseURL       string `envconfig:"HOTELBYTE_BASE_URL" default:"https://api.hotelbyte.com"`
    Timeout       time.Duration `envconfig:"HOTELBYTE_TIMEOUT" default:"30s"`
    MaxRetries    int    `envconfig:"HOTELBYTE_MAX_RETRIES" default:"3"`
    Debug         bool   `envconfig:"HOTELBYTE_DEBUG" default:"false"`
}

func NewHotelClient() (*hotelbyte.Client, error) {
    var cfg Config
    if err := envconfig.Process("", &cfg); err != nil {
        return nil, err
    }
    
    options := []hotelbyte.ClientOption{
        hotelbyte.WithCredentials(cfg.AppKey, cfg.AppSecret),
        hotelbyte.WithBaseURL(cfg.BaseURL),
        hotelbyte.WithTimeout(cfg.Timeout),
        hotelbyte.WithRetryConfig(cfg.MaxRetries, time.Second, 30*time.Second),
    }
    
    return hotelbyte.NewClient(options...)
}
```

## 2. 认证管理

### 2.1 自动认证

**推荐：依赖 SDK 的自动认证机制**

```go
// ✅ SDK 会自动处理认证和令牌刷新
func searchHotels(ctx context.Context) (*hotelbyte.SearchHotelsResponse, error) {
    client := GetHotelClient()
    
    // 无需手动认证，SDK 会自动处理
    return client.SearchHotels(ctx, &hotelbyte.SearchHotelsRequest{
        DestinationID: "beijing",
        CheckIn:       time.Now().AddDate(0, 0, 7),
        CheckOut:      time.Now().AddDate(0, 0, 9),
        AdultCount:    2,
        RoomCount:     1,
    })
}
```

### 2.2 预认证策略

**对于关键操作，提前验证认证状态**

```go
func criticalOperation(ctx context.Context) error {
    client := GetHotelClient()
    
    // 对于关键操作，预先确保认证状态
    if !client.Auth().IsTokenValid() {
        if err := client.Auth().Authenticate(ctx); err != nil {
            return fmt.Errorf("预认证失败: %w", err)
        }
    }
    
    // 执行关键操作
    return performCriticalTask(ctx, client)
}
```

## 3. 错误处理

### 3.1 分层错误处理

**在不同层次处理不同类型的错误**

```go
// 业务层：处理业务逻辑错误
func SearchAndBookHotel(ctx context.Context, req BookingRequest) (*BookingResult, error) {
    client := GetHotelClient()
    
    // 搜索酒店
    hotels, err := client.SearchHotels(ctx, req.SearchRequest)
    if err != nil {
        return nil, handleSearchError(err)
    }
    
    if len(hotels.Hotels) == 0 {
        return nil, &BusinessError{
            Code:    "NO_HOTELS_FOUND",
            Message: "未找到符合条件的酒店",
        }
    }
    
    // 预订酒店
    booking, err := client.CreateBooking(ctx, req.BookingRequest)
    if err != nil {
        return nil, handleBookingError(err)
    }
    
    return &BookingResult{Booking: booking}, nil
}

// 错误处理函数
func handleSearchError(err error) error {
    switch {
    case hotelbyte.IsAuthenticationError(err):
        return &BusinessError{
            Code:    "AUTH_FAILED",
            Message: "认证失败，请检查凭据",
            Cause:   err,
        }
    case hotelbyte.IsValidationError(err):
        return &BusinessError{
            Code:    "INVALID_SEARCH_PARAMS",
            Message: "搜索参数无效",
            Cause:   err,
        }
    case hotelbyte.IsRateLimitError(err):
        return &BusinessError{
            Code:    "RATE_LIMITED",
            Message: "请求过于频繁，请稍后重试",
            Cause:   err,
            Retryable: true,
        }
    case hotelbyte.IsSystemError(err):
        return &BusinessError{
            Code:    "SYSTEM_ERROR",
            Message: "系统暂时不可用，请稍后重试",
            Cause:   err,
            Retryable: true,
        }
    default:
        return &BusinessError{
            Code:    "UNKNOWN_ERROR",
            Message: "未知错误",
            Cause:   err,
        }
    }
}
```

### 3.2 错误重试策略

**实现自定义重试逻辑**

```go
func retryableRequest[T any](
    ctx context.Context,
    operation func(context.Context) (T, error),
    isRetryable func(error) bool,
    maxRetries int,
) (T, error) {
    var zero T
    var lastErr error
    
    for attempt := 0; attempt <= maxRetries; attempt++ {
        if attempt > 0 {
            // 指数退避
            delay := time.Duration(1<<uint(attempt-1)) * time.Second
            select {
            case <-ctx.Done():
                return zero, ctx.Err()
            case <-time.After(delay):
            }
        }
        
        result, err := operation(ctx)
        if err == nil {
            return result, nil
        }
        
        lastErr = err
        if !isRetryable(err) {
            break
        }
    }
    
    return zero, fmt.Errorf("操作失败，已重试 %d 次: %w", maxRetries, lastErr)
}

// 使用示例
func searchHotelsWithRetry(ctx context.Context, req *hotelbyte.SearchHotelsRequest) (*hotelbyte.SearchHotelsResponse, error) {
    return retryableRequest(
        ctx,
        func(ctx context.Context) (*hotelbyte.SearchHotelsResponse, error) {
            return GetHotelClient().SearchHotels(ctx, req)
        },
        func(err error) bool {
            return hotelbyte.IsRateLimitError(err) || 
                   hotelbyte.IsSystemError(err) ||
                   hotelbyte.IsTimeoutError(err)
        },
        3,
    )
}
```

## 4. 并发处理

### 4.1 并发安全

**SDK 是并发安全的，可以在多个 goroutine 中使用**

```go
func ConcurrentHotelSearch(ctx context.Context, destinations []string) map[string]*hotelbyte.SearchHotelsResponse {
    client := GetHotelClient()
    results := make(map[string]*hotelbyte.SearchHotelsResponse)
    var mu sync.Mutex
    var wg sync.WaitGroup
    
    // 并发搜索多个目的地
    for _, dest := range destinations {
        wg.Add(1)
        go func(destination string) {
            defer wg.Done()
            
            resp, err := client.SearchHotels(ctx, &hotelbyte.SearchHotelsRequest{
                DestinationID: destination,
                CheckIn:       time.Now().AddDate(0, 0, 7),
                CheckOut:      time.Now().AddDate(0, 0, 9),
                AdultCount:    2,
                RoomCount:     1,
            })
            
            mu.Lock()
            if err == nil {
                results[destination] = resp
            }
            mu.Unlock()
        }(dest)
    }
    
    wg.Wait()
    return results
}
```

### 4.2 控制并发数量

**使用信号量控制并发数量**

```go
type ConcurrentSearcher struct {
    client    *hotelbyte.Client
    semaphore chan struct{}
}

func NewConcurrentSearcher(maxConcurrency int) *ConcurrentSearcher {
    return &ConcurrentSearcher{
        client:    GetHotelClient(),
        semaphore: make(chan struct{}, maxConcurrency),
    }
}

func (s *ConcurrentSearcher) SearchHotels(ctx context.Context, req *hotelbyte.SearchHotelsRequest) (*hotelbyte.SearchHotelsResponse, error) {
    // 获取信号量
    select {
    case s.semaphore <- struct{}{}:
        defer func() { <-s.semaphore }()
    case <-ctx.Done():
        return nil, ctx.Err()
    }
    
    return s.client.SearchHotels(ctx, req)
}
```

## 5. 性能优化

### 5.1 请求批处理

**批量处理相似请求**

```go
type BatchProcessor struct {
    client    *hotelbyte.Client
    batchSize int
    timeout   time.Duration
}

func (bp *BatchProcessor) BatchSearchHotels(ctx context.Context, requests []*hotelbyte.SearchHotelsRequest) ([]*hotelbyte.SearchHotelsResponse, error) {
    results := make([]*hotelbyte.SearchHotelsResponse, len(requests))
    errors := make([]error, len(requests))
    
    // 分批处理
    for i := 0; i < len(requests); i += bp.batchSize {
        end := i + bp.batchSize
        if end > len(requests) {
            end = len(requests)
        }
        
        batch := requests[i:end]
        if err := bp.processBatch(ctx, batch, results[i:end], errors[i:end]); err != nil {
            return nil, err
        }
    }
    
    return results, nil
}

func (bp *BatchProcessor) processBatch(ctx context.Context, batch []*hotelbyte.SearchHotelsRequest, results []*hotelbyte.SearchHotelsResponse, errors []error) error {
    var wg sync.WaitGroup
    
    for i, req := range batch {
        wg.Add(1)
        go func(index int, request *hotelbyte.SearchHotelsRequest) {
            defer wg.Done()
            
            resp, err := bp.client.SearchHotels(ctx, request)
            if err != nil {
                errors[index] = err
            } else {
                results[index] = resp
            }
        }(i, req)
    }
    
    wg.Wait()
    return nil
}
```

### 5.2 结果缓存

**缓存频繁访问的数据**

```go
type CachedSearcher struct {
    client *hotelbyte.Client
    cache  *sync.Map // 生产环境建议使用 Redis
}

func NewCachedSearcher() *CachedSearcher {
    return &CachedSearcher{
        client: GetHotelClient(),
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

func (cs *CachedSearcher) generateCacheKey(req *hotelbyte.SearchHotelsRequest) string {
    return fmt.Sprintf("%s:%d:%d:%d:%d",
        req.DestinationID,
        req.CheckIn.Unix(),
        req.CheckOut.Unix(),
        req.AdultCount,
        req.RoomCount,
    )
}
```

## 6. 超时和取消

### 6.1 适当的超时设置

**为不同类型的操作设置合适的超时**

```go
// 搜索操作：较短超时
func QuickSearch(ctx context.Context, req *hotelbyte.SearchHotelsRequest) (*hotelbyte.SearchHotelsResponse, error) {
    ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
    defer cancel()
    
    return GetHotelClient().SearchHotels(ctx, req)
}

// 预订操作：较长超时
func BookHotel(ctx context.Context, req *hotelbyte.CreateBookingRequest) (*hotelbyte.BookingResponse, error) {
    ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
    defer cancel()
    
    return GetHotelClient().CreateBooking(ctx, req)
}
```

### 6.2 优雅的取消处理

**正确处理请求取消**

```go
func InterruptibleSearch(ctx context.Context, req *hotelbyte.SearchHotelsRequest) (*hotelbyte.SearchHotelsResponse, error) {
    // 创建可取消的上下文
    ctx, cancel := context.WithCancel(ctx)
    defer cancel()
    
    // 在单独的 goroutine 中执行请求
    resultChan := make(chan searchResult, 1)
    go func() {
        resp, err := GetHotelClient().SearchHotels(ctx, req)
        resultChan <- searchResult{resp: resp, err: err}
    }()
    
    // 等待结果或取消
    select {
    case result := <-resultChan:
        return result.resp, result.err
    case <-ctx.Done():
        return nil, ctx.Err()
    }
}

type searchResult struct {
    resp *hotelbyte.SearchHotelsResponse
    err  error
}
```

## 7. 日志和监控

### 7.1 结构化日志

**使用结构化日志记录重要事件**

```go
import (
    "github.com/sirupsen/logrus"
)

func loggedSearch(ctx context.Context, req *hotelbyte.SearchHotelsRequest) (*hotelbyte.SearchHotelsResponse, error) {
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
    
    resp, err := GetHotelClient().SearchHotels(ctx, req)
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

### 7.2 性能监控

**收集性能指标**

```go
import (
    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promauto"
)

var (
    requestDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
        Name: "hotelbyte_request_duration_seconds",
        Help: "Duration of HotelByte API requests",
    }, []string{"operation", "status"})
    
    requestCount = promauto.NewCounterVec(prometheus.CounterOpts{
        Name: "hotelbyte_requests_total",
        Help: "Total number of HotelByte API requests",
    }, []string{"operation", "status"})
)

func monitoredSearch(ctx context.Context, req *hotelbyte.SearchHotelsRequest) (*hotelbyte.SearchHotelsResponse, error) {
    start := time.Now()
    operation := "search_hotels"
    
    resp, err := GetHotelClient().SearchHotels(ctx, req)
    
    duration := time.Since(start).Seconds()
    status := "success"
    if err != nil {
        status = "error"
        if hotelbyte.IsRateLimitError(err) {
            status = "rate_limited"
        } else if hotelbyte.IsTimeoutError(err) {
            status = "timeout"
        }
    }
    
    requestDuration.WithLabelValues(operation, status).Observe(duration)
    requestCount.WithLabelValues(operation, status).Inc()
    
    return resp, err
}
```

## 8. 测试最佳实践

### 8.1 模拟测试

**创建可测试的接口**

```go
//go:generate mockgen -source=hotel_service.go -destination=mocks/hotel_service_mock.go

type HotelService interface {
    SearchHotels(ctx context.Context, req *hotelbyte.SearchHotelsRequest) (*hotelbyte.SearchHotelsResponse, error)
    GetHotelDetails(ctx context.Context, hotelID string) (*hotelbyte.HotelDetails, error)
    CreateBooking(ctx context.Context, req *hotelbyte.CreateBookingRequest) (*hotelbyte.BookingResponse, error)
}

type hotelServiceImpl struct {
    client *hotelbyte.Client
}

func NewHotelService(client *hotelbyte.Client) HotelService {
    return &hotelServiceImpl{client: client}
}

func (s *hotelServiceImpl) SearchHotels(ctx context.Context, req *hotelbyte.SearchHotelsRequest) (*hotelbyte.SearchHotelsResponse, error) {
    return s.client.SearchHotels(ctx, req)
}

// 测试
func TestBookingFlow(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()
    
    mockService := mocks.NewMockHotelService(ctrl)
    
    // 设置期望
    mockService.EXPECT().
        SearchHotels(gomock.Any(), gomock.Any()).
        Return(&hotelbyte.SearchHotelsResponse{
            Hotels: []hotelbyte.HotelSummary{{ID: "hotel-1"}},
            Total:  1,
        }, nil)
    
    // 测试业务逻辑
    result, err := BookingFlow(context.Background(), mockService, BookingRequest{})
    assert.NoError(t, err)
    assert.NotNil(t, result)
}
```

### 8.2 集成测试

**编写端到端测试**

```go
func TestIntegration(t *testing.T) {
    if testing.Short() {
        t.Skip("跳过集成测试")
    }
    
    // 使用测试环境配置
    client, err := hotelbyte.NewClient(
        hotelbyte.WithCredentials(
            os.Getenv("TEST_APP_KEY"),
            os.Getenv("TEST_APP_SECRET"),
        ),
        hotelbyte.WithBaseURL(os.Getenv("TEST_BASE_URL")),
    )
    require.NoError(t, err)
    defer client.Close()
    
    ctx := context.Background()
    
    // 测试搜索功能
    resp, err := client.SearchHotels(ctx, &hotelbyte.SearchHotelsRequest{
        DestinationID: "test-destination",
        CheckIn:       time.Now().AddDate(0, 0, 7),
        CheckOut:      time.Now().AddDate(0, 0, 9),
        AdultCount:    2,
        RoomCount:     1,
    })
    
    assert.NoError(t, err)
    assert.NotNil(t, resp)
    
    if len(resp.Hotels) > 0 {
        // 测试酒店详情
        details, err := client.GetHotelDetails(ctx, resp.Hotels[0].ID)
        assert.NoError(t, err)
        assert.NotNil(t, details)
    }
}
```

## 9. 部署和运维

### 9.1 健康检查

**实现健康检查端点**

```go
func HealthCheck(w http.ResponseWriter, r *http.Request) {
    ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
    defer cancel()
    
    client := GetHotelClient()
    
    // 检查认证状态
    if !client.Auth().IsTokenValid() {
        if err := client.Auth().Authenticate(ctx); err != nil {
            http.Error(w, fmt.Sprintf("Authentication failed: %v", err), http.StatusServiceUnavailable)
            return
        }
    }
    
    // 可以添加更多健康检查...
    
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]string{
        "status": "healthy",
        "timestamp": time.Now().UTC().Format(time.RFC3339),
    })
}
```

### 9.2 优雅关闭

**实现优雅关闭**

```go
func main() {
    client := GetHotelClient()
    defer client.Close()
    
    // 设置信号处理
    c := make(chan os.Signal, 1)
    signal.Notify(c, os.Interrupt, syscall.SIGTERM)
    
    server := &http.Server{
        Addr:    ":8080",
        Handler: setupRoutes(),
    }
    
    go func() {
        if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            log.Fatal("Server failed:", err)
        }
    }()
    
    // 等待中断信号
    <-c
    log.Println("Shutting down gracefully...")
    
    // 关闭服务器
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()
    
    if err := server.Shutdown(ctx); err != nil {
        log.Printf("Server shutdown error: %v", err)
    }
    
    log.Println("Server stopped")
}
```

## 10. 性能调优指南

### 10.1 连接池优化

```go
client, err := hotelbyte.NewClient(
    hotelbyte.WithCredentials(appKey, appSecret),
    // 根据预期负载调整连接池大小
    hotelbyte.WithHTTPConfig(hotelbyte.HTTPConfig{
        MaxIdleConns:    200,  // 总的空闲连接数
        MaxConnsPerHost: 50,   // 每个主机的连接数
        Timeout:         30 * time.Second,
    }),
)
```

### 10.2 内存使用优化

```go
// 定期清理缓存
func periodicCacheCleanup(cache *sync.Map) {
    ticker := time.NewTicker(1 * time.Hour)
    defer ticker.Stop()
    
    for range ticker.C {
        now := time.Now()
        cache.Range(func(key, value interface{}) bool {
            if entry, ok := value.(*cacheEntry); ok && entry.isExpired() {
                cache.Delete(key)
            }
            return true
        })
    }
}
```

这些最佳实践将帮助您构建高质量、高性能的应用程序，充分利用 HotelByte SDK 的功能。