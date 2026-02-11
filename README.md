# Fibox

Fibox adalah library Go yang menyediakan utilitas untuk membangun aplikasi web dengan [Fiber](https://github.com/gofiber/fiber). Library ini mencakup modul untuk response handling, JWT authentication, HTTP client, middleware, caching, dan logging.

## Fitur

- **Response** - Response handler standar untuk API dengan format JSON yang konsisten
- **JWT** - Service untuk generate dan validate JWT token
- **HTTP Client** - Wrapper untuk resty dengan retry, timeout, dan konfigurasi yang mudah
- **Middleware** - Authentication dan rate limiting middleware
- **Cache** - LRU (Least Recently Used) cache dengan TTL support
- **Logging** - Structured logging menggunakan zerolog

## Instalasi

```bash
go get github.com/username/fibox
```

## Penggunaan

### Response Handler

```go
import "fibox/response"

// Success response
response.Success(c, data)

// Success with message
response.SuccessWithMessage(c, "Data berhasil disimpan")

// Created response (201)
response.Created(c, data)

// Error responses
response.BadRequest(c, "Invalid input")
response.Unauthorized(c, "Token tidak valid")
response.Forbidden(c, "Akses ditolak")
response.NotFound(c, "Data tidak ditemukan")
response.InternalError(c, "Terjadi kesalahan server")
response.ValidationError(c, map[string]string{"email": "Invalid email format"})
```

### JWT Authentication

```go
import "fibox/jwt"

// Buat JWT service
jwtSvc := jwt.NewJWTService("secret-key", 24) // 24 jam expiry

// Generate token
token, err := jwtSvc.GenerateToken("user123", "user@example.com", "admin")

// Generate refresh token
refreshToken, err := jwtSvc.GenerateRefreshToken("user123", "user@example.com", "admin")

// Validate token
claims, err := jwtSvc.ValidateToken(token)
```

### HTTP Client

```go
import "fibox/client"

// Buat HTTP client dengan konfigurasi
http := client.NewHTTPClient(client.HTTPClientConfig{
    BaseURL:          "https://api.example.com",
    Timeout:          30 * time.Second,
    RetryCount:       3,
    RetryWaitTime:    1 * time.Second,
    RetryMaxWaitTime: 30 * time.Second,
})

// Atau gunakan default client
http := client.GetDefaultHTTPClient("https://api.example.com")

// Set auth token
http.SetBearerToken("your-token")

// GET request
var result map[string]interface{}
err := http.Get("/users", map[string]string{"page": "1"}, &result)

// POST request
body := map[string]string{"name": "John"}
err := http.Post("/users", body, &result)

// PUT request
err := http.Put("/users/123", body, &result)

// DELETE request
err := http.Delete("/users/123", nil, &result)
```

### Middleware

```go
import "fibox/middleware"

// Auth middleware
app.Use(middleware.AuthMiddleware(jwtSvc))

// Rate limiter untuk auth endpoints (max 5 request/menit)
app.Post("/login", middleware.NewAuthRateLimiter(5), handler)

// Rate limiter untuk general endpoints (max 60 request/menit)
app.Get("/api/*", middleware.NewGeneralRateLimiter(60), handler)

// Get auth info dari context
func handler(c fiber.Ctx) error {
    auth := middleware.GetAuthInfo(c)
    fmt.Println(auth.UserID, auth.Email, auth.Role)
}
```

### Cache

```go
import "fibox/cache"
import "context"

// Buat LRU cache dengan kapasitas 1000 items
cache := cache.NewLRUCache(1000)

ctx := context.Background()

// Set value
cache.Set(ctx, "user:123", userData)

// Set dengan TTL (5 menit)
cache.SetWithTTL(ctx, "user:123", userData, 5*time.Minute)

// Get value
if value, found := cache.Get(ctx, "user:123"); found {
    fmt.Println(value)
}

// Delete value
cache.Delete(ctx, "user:123")

// Delete by pattern
cache.DeleteByPattern(ctx, "user:*")

// Clear semua cache
cache.Clear(ctx)

// Get cache stats
stats := cache.Stats()
fmt.Printf("Hits: %d, Misses: %d, Size: %d\n", stats.Hits, stats.Misses, stats.Size)
```

### Logging

```go
import "fibox/logging"

// Set log level
logging.SetLogLevel("debug") // debug, info, warn, error, fatal, panic

// Simple logging
logging.Info("Application started")
logging.Error("Something went wrong", err)
logging.Debug("Debug info")
logging.Warn("Warning message")

// Logging dengan fields
logging.InfoWithFields("User logged in", map[string]interface{}{
    "user_id": "123",
    "ip":      "127.0.0.1",
})

logging.ErrorWithFields("Database error", err, map[string]interface{}{
    "query": "SELECT * FROM users",
})
```

## Contoh Aplikasi Lengkap

```go
package main

import (
    "fibox/cache"
    "fibox/client"
    "fibox/jwt"
    "fibox/logging"
    "fibox/middleware"
    "fibox/response"
    "time"

    "github.com/gofiber/fiber/v3"
)

func main() {
    // Initialize services
    jwtSvc := jwt.NewJWTService("your-secret-key", 24)
    cache := cache.NewLRUCache(1000)
    logging.SetLogLevel("info")

    app := fiber.New()

    // Public route
    app.Post("/login", func(c fiber.Ctx) error {
        // Login logic...
        token, _ := jwtSvc.GenerateToken("user123", "user@example.com", "admin")
        return response.Success(c, fiber.Map{"token": token})
    })

    // Protected route
    api := app.Group("/api")
    api.Use(middleware.AuthMiddleware(jwtSvc))
    api.Use(middleware.NewGeneralRateLimiter(60))

    api.Get("/profile", func(c fiber.Ctx) error {
        auth := middleware.GetAuthInfo(c)
        return response.Success(c, auth)
    })

    app.Listen(":3000")
}
```

## Lisensi

MIT
