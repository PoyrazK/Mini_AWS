package services

import (
	"context"
	"time"

	"github.com/poyrazk/thecloud/internal/core/ports"
)

type Checkable interface {
	Ping(ctx context.Context) error
}

type HealthServiceImpl struct {
	db     Checkable
	docker ports.DockerClient
}

func NewHealthServiceImpl(db Checkable, docker ports.DockerClient) *HealthServiceImpl {
	return &HealthServiceImpl{
		db:     db,
		docker: docker,
	}
}

func (s *HealthServiceImpl) Check(ctx context.Context) ports.HealthCheckResult {
	checks := make(map[string]string)
	overall := "UP"

	// Check DB
	if err := s.db.Ping(ctx); err != nil {
		checks["database"] = "DISCONNECTED: " + err.Error()
		overall = "DEGRADED"
	} else {
		checks["database"] = "CONNECTED"
	}

	// Check Docker
	if err := s.docker.Ping(ctx); err != nil {
		checks["docker"] = "DISCONNECTED: " + err.Error()
		overall = "DEGRADED"
	} else {
		checks["docker"] = "CONNECTED"
	}

	return ports.HealthCheckResult{
		Status: overall,
		Checks: checks,
		Time:   time.Now(),
	}
}
