# HotelByte Go SDK

[![Go Version](https://img.shields.io/badge/Go-1.21+-blue.svg)](https://golang.org)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![Documentation](https://img.shields.io/badge/Documentation-GoDoc-blue.svg)](https://pkg.go.dev/github.com/hotelbyte-com/sdk-go)

HotelByte Go SDK is the official Go language software development kit for the Hotel API Distribution Platform, providing simple, efficient, and type-safe API interface access capabilities.

## ✨ Core Features

- **🔗 Unified Interface**: Aggregates multiple hotel suppliers (HotelBeds, Dida, DerbySoft, etc.) with consistent API experience
- **🚀 High Performance**: Built on Go's concurrency features, supporting high-concurrency request processing
- **🛡️ Type Safety**: Complete type definitions with compile-time error checking
- **🔄 Auto Retry**: Built-in intelligent retry mechanism for automatic handling of network anomalies and temporary failures
- **📝 Comprehensive Documentation**: Detailed API documentation and usage examples
- **🧪 Test Coverage**: High test coverage ensuring SDK stability
- **🌐 Internationalization**: Support for Chinese and English languages

## 📦 Installation

Install the SDK using Go modules:

```bash
go get github.com/hotelbyte-com/sdk-go
```

Import in your code:

```go
import "github.com/hotelbyte-com/sdk-go"
```

## 🚀 Quick Start

### Basic Configuration

```go
package main

import (
    "context"
    "fmt"
    "time"

    "github.com/hotelbyte-com/sdk-go"
)

func main() {
    // Create client
    client, err := hotelbyte.NewClient(
        hotelbyte.WithCredentials("your-app-key", "your-app-secret"),
        hotelbyte.WithBaseURL("https://api.hotelbyte.com"),
        hotelbyte.WithTimeout(30*time.Second),
    )
    if err != nil {
        panic(err)
    }
    defer client.Close()

    // Search hotels
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

    fmt.Printf("Found %d hotels\n", len(resp.Hotels))
    for _, hotel := range resp.Hotels {
        fmt.Printf("- %s: %s\n", hotel.Name, hotel.Address)
    }
}
```

### Environment Variable Configuration

```bash
export HOTELBYTE_APP_KEY="your-app-key"
export HOTELBYTE_APP_SECRET="your-app-secret"
export HOTELBYTE_BASE_URL="https://api.hotelbyte.com"
```

```go
// Create client using environment variables
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

## 📖 Main Features

### 🔍 Hotel Search

```go
// Basic search
resp, err := client.SearchHotels(ctx, &hotelbyte.SearchHotelsRequest{
    DestinationID: "beijing",
    CheckIn:       time.Now().AddDate(0, 0, 7),
    CheckOut:      time.Now().AddDate(0, 0, 9),
    AdultCount:    2,
    RoomCount:     1,
})

// Advanced search parameters
resp, err := client.SearchHotels(ctx, &hotelbyte.SearchHotelsRequest{
    DestinationID: "beijing",
    CheckIn:       time.Now().AddDate(0, 0, 7),
    CheckOut:      time.Now().AddDate(0, 0, 9),
    AdultCount:    2,
    RoomCount:     1,
    HotelIDs:      []string{"hotel-1", "hotel-2"},
    MinRating:     4,
    MaxPrice:      1000.00,
    Currency:      "USD",
    Language:      "en-US",
})
```

### 🏨 Hotel Details

```go
// Get detailed hotel information
details, err := client.GetHotelDetails(ctx, "hotel-id")
if err != nil {
    return err
}

fmt.Printf("Hotel Name: %s\n", details.Name)
fmt.Printf("Hotel Address: %s\n", details.Address)
fmt.Printf("Hotel Facilities: %v\n", details.Facilities)
```

### 📋 Booking Management

```go
// Create booking
booking, err := client.CreateBooking(ctx, &hotelbyte.CreateBookingRequest{
    HotelID:    "hotel-id",
    RoomTypeID: "room-type-id",
    CheckIn:    time.Now().AddDate(0, 0, 7),
    CheckOut:   time.Now().AddDate(0, 0, 9),
    GuestInfo: hotelbyte.GuestInfo{
        Name:    "John Doe",
        Email:   "john.doe@example.com",
        Phone:   "+1-555-1234-5678",
    },
    PaymentInfo: hotelbyte.PaymentInfo{
        Method: "credit_card",
        Card:   "****-****-****-1234",
    },
})

// Query booking status
status, err := client.GetBookingStatus(ctx, booking.BookingID)

// Cancel booking
err = client.CancelBooking(ctx, booking.BookingID)
```

### 🔐 Authentication Management

```go
// SDK handles authentication automatically, but manual management is also possible
auth := client.Auth()

// Check authentication status
if auth.IsTokenValid() {
    fmt.Println("Authentication is valid")
}

// Manually refresh authentication
err := auth.RefreshToken(ctx)
if err != nil {
    fmt.Printf("Authentication refresh failed: %v\n", err)
}
```

## ⚙️ Client Configuration

### Basic Configuration Options

```go
client, err := hotelbyte.NewClient(
    // Authentication information
    hotelbyte.WithCredentials("app-key", "app-secret"),

    // Basic configuration
    hotelbyte.WithBaseURL("https://api.hotelbyte.com"),
    hotelbyte.WithTimeout(30 * time.Second),
    hotelbyte.WithUserAgent("my-app/1.0"),

    // Retry configuration
    hotelbyte.WithRetryConfig(3, time.Second, 30*time.Second),

    // HTTP configuration
    hotelbyte.WithHTTPConfig(hotelbyte.HTTPConfig{
        MaxIdleConns:    100,
        MaxConnsPerHost: 10,
        IdleConnTimeout: 90 * time.Second,
    }),

    // Debug mode
    hotelbyte.WithDebug(true),
)
```

### Advanced Configuration

```go
// Custom HTTP client
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

// Custom logger
logger := logrus.New()
logger.SetLevel(logrus.InfoLevel)

client, err := hotelbyte.NewClient(
    hotelbyte.WithLogger(logger),
    hotelbyte.WithCredentials("app-key", "app-secret"),
)
```

## 🔧 Error Handling

### Error Type Checking

```go
resp, err := client.SearchHotels(ctx, req)
if err != nil {
    switch {
    case hotelbyte.IsAuthenticationError(err):
        fmt.Println("Authentication failed, please check credentials")
    case hotelbyte.IsValidationError(err):
        fmt.Println("Request parameter validation failed")
    case hotelbyte.IsRateLimitError(err):
        fmt.Println("Request rate too high, please try again later")
    case hotelbyte.IsSystemError(err):
        fmt.Println("System error, please try again later")
    case hotelbyte.IsTimeoutError(err):
        fmt.Println("Request timeout")
    default:
        fmt.Printf("Unknown error: %v\n", err)
    }
    return
}
```

### Custom Retry Strategy

```go
// Implement custom retry logic
func searchWithRetry(ctx context.Context, client *hotelbyte.Client, req *hotelbyte.SearchHotelsRequest) (*hotelbyte.SearchHotelsResponse, error) {
    var lastErr error

    for attempt := 0; attempt < 3; attempt++ {
        if attempt > 0 {
            // Exponential backoff
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

        // Check if retryable
        if !hotelbyte.IsRateLimitError(err) && !hotelbyte.IsSystemError(err) {
            break
        }
    }

    return nil, fmt.Errorf("search failed after 3 retries: %w", lastErr)
}
```

## 🚀 Concurrent Processing

### Concurrent Search

```go
func concurrentSearch(ctx context.Context, client *hotelbyte.Client, destinations []string) map[string]*hotelbyte.SearchHotelsResponse {
    results := make(map[string]*hotelbyte.SearchHotelsResponse)
    var mu sync.Mutex
    var wg sync.WaitGroup

    // Limit concurrent count
    semaphore := make(chan struct{}, 5)

    for _, dest := range destinations {
        wg.Add(1)
        go func(destination string) {
            defer wg.Done()

            // Acquire semaphore
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

## 📊 Performance Optimization

### Caching Strategy

```go
type CachedSearcher struct {
    client *hotelbyte.Client
    cache  *sync.Map // Use Redis in production
}

func NewCachedSearcher(client *hotelbyte.Client) *CachedSearcher {
    return &CachedSearcher{
        client: client,
        cache:  &sync.Map{},
    }
}

func (cs *CachedSearcher) SearchHotels(ctx context.Context, req *hotelbyte.SearchHotelsRequest) (*hotelbyte.SearchHotelsResponse, error) {
    // Generate cache key
    key := cs.generateCacheKey(req)

    // Check cache
    if cached, ok := cs.cache.Load(key); ok {
        if entry, ok := cached.(*cacheEntry); ok && !entry.isExpired() {
            return entry.data, nil
        }
    }

    // Cache miss, call API
    resp, err := cs.client.SearchHotels(ctx, req)
    if err != nil {
        return nil, err
    }

    // Cache result
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

## 📝 Example Code

Check the [examples](./examples) directory for more detailed examples:

- **[Quick Start](./examples/quickstart/)** - Simplest usage example
- **[Authentication](./examples/authentication/)** - Authentication operations
- **[Hotel Search](./examples/hotel-search/)** - Search functionality examples
- **[Booking Management](./examples/booking-management/)** - Booking process examples
- **[Error Handling](./examples/error-handling/)** - Complete error handling strategies
- **[Concurrent Processing](./examples/concurrent-processing/)** - Concurrent request processing
- **[Performance Optimization](./examples/caching/)** - Caching and performance optimization
- **[Monitoring Integration](./examples/monitoring/)** - Monitoring and metrics collection

## 🔧 Best Practices

### 1. Client Management

```go
// ✅ Recommended: Global singleton pattern
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

### 2. Environment Variable Configuration

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

### 3. Structured Logging

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

    logger.Info("Starting hotel search")
    start := time.Now()

    resp, err := client.SearchHotels(ctx, req)
    duration := time.Since(start)

    if err != nil {
        logger.WithFields(logrus.Fields{
            "error":    err.Error(),
            "duration": duration,
        }).Error("Hotel search failed")
        return nil, err
    }

    logger.WithFields(logrus.Fields{
        "hotel_count": len(resp.Hotels),
        "total":       resp.Total,
        "duration":    duration,
    }).Info("Hotel search successful")

    return resp, nil
}
```

## 📚 Documentation

- [API Reference](./api-reference.md) - Complete API documentation
- [Configuration Guide](./configuration.md) - Configuration options details
- [Troubleshooting](./troubleshooting.md) - Common problem solutions
- [Example Code](./examples/README.md) - Comprehensive usage examples

## 🧪 Testing

Run tests:

```bash
# Run all tests
go test ./...

# Run tests with coverage report
go test -cover ./...

# Run benchmarks
go test -bench=. ./...

# Run tests for specific package
go test ./protocol/...
```

## 🤝 Contributing

We welcome community contributions! Please check the [Contributing Guide](./CONTRIBUTING.md) to learn how to participate in project development.

### Contribution Process

1. Fork the project repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Create a Pull Request

### Code Standards

- Follow Go official code standards
- Add appropriate test cases
- Update relevant documentation
- Ensure all tests pass

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🆘 Support

If you encounter problems during use:

- 📧 **Email Support**: developers@hotelbyte.com
- 🐛 **Issue Reporting**: [GitHub Issues](https://github.com/hotelbyte-com/hotel-be/issues)
- 💬 **Discussions**: [GitHub Discussions](https://github.com/hotelbyte-com/hotel-be/discussions)
- 📖 **Project Documentation**: [HotelByte Technical Documentation](https://docs.hotelbyte.com)

## 🌟 Acknowledgments

Thanks to all developers and users who have contributed to the HotelByte SDK!

---

**Note**: This SDK is the official Go language implementation of the HotelByte Hotel API Distribution Platform. If you need SDKs in other languages, please contact us or check our other projects.