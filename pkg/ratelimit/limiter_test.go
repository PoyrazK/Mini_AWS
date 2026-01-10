package ratelimit_test

import (
	"log/slog"
	"testing"

	"github.com/poyrazk/thecloud/pkg/ratelimit"
	"golang.org/x/time/rate"
)

func BenchmarkRateLimiter_GetLimiter(b *testing.B) {
	limiter := ratelimit.NewIPRateLimiter(rate.Limit(100), 10, slog.Default())

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = limiter.GetLimiter("192.168.1.1")
	}
}

func BenchmarkRateLimiter_GetLimiterParallel(b *testing.B) {
	limiter := ratelimit.NewIPRateLimiter(rate.Limit(1000), 100, slog.Default())

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = limiter.GetLimiter("192.168.1.1")
		}
	})
}

func BenchmarkRateLimiter_GetLimiterParallel_MultiKey(b *testing.B) {
	limiter := ratelimit.NewIPRateLimiter(rate.Limit(1000), 100, slog.Default())

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			// Simulate different clients
			key := "key-" + string(rune(i%100))
			_ = limiter.GetLimiter(key)
			i++
		}
	})
}
