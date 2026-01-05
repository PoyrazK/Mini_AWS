package ports

import (
	"context"
	"time"
)

type HealthCheckResult struct {
	Status string            `json:"status"`
	Checks map[string]string `json:"checks"`
	Time   time.Time         `json:"time"`
}

type HealthService interface {
	Check(ctx context.Context) HealthCheckResult
}
