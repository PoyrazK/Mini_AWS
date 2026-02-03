// Package domain contains the core business entities.
package domain

import (
	"time"

	"github.com/google/uuid"
)

// GlobalLoadBalancer represents a multi-region load balancer that routes traffic
// based on policies like latency, geolocation, or failover.
type GlobalLoadBalancer struct {
	ID        uuid.UUID         `json:"id"`
	UserID    uuid.UUID         `json:"user_id"`
	TenantID  uuid.UUID         `json:"tenant_id"`
	Name      string            `json:"name"`
	Hostname  string            `json:"hostname"` // e.g., "api.myapp.thecloud.io"
	Policy    RoutingPolicy     `json:"routing_policy"`
	Endpoints []*GlobalEndpoint `json:"endpoints,omitempty"`
	Status    string            `json:"status"` // CREATING, ACTIVE, DELETING

	// Health check configuration for all endpoints
	HealthCheck GlobalHealthCheckConfig `json:"health_check"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// RoutingPolicy defines how traffic is distributed.
type RoutingPolicy string

const (
	// RoutingLatency routes traffic to the region with the lowest latency.
	RoutingLatency RoutingPolicy = "LATENCY"
	// RoutingGeolocation routes traffic based on the user's geographic location.
	RoutingGeolocation RoutingPolicy = "GEOLOCATION"
	// RoutingWeighted distributes traffic based on relative weights.
	RoutingWeighted RoutingPolicy = "WEIGHTED"
	// RoutingFailover uses priority-based failover.
	RoutingFailover RoutingPolicy = "FAILOVER"
)

// GlobalEndpoint is a destination for traffic, usually a regional load balancer or static IP.
type GlobalEndpoint struct {
	ID              uuid.UUID  `json:"id"`
	GlobalLBID      uuid.UUID  `json:"global_lb_id"`
	Region          string     `json:"region"`              // e.g., "us-east-1"
	TargetType      string     `json:"target_type"`         // "LB" or "IP"
	TargetID        *uuid.UUID `json:"target_id,omitempty"` // Reference to regional LB
	TargetIP        *string    `json:"target_ip,omitempty"` // Static IP if TargetType == "IP"
	Weight          int        `json:"weight"`              // 1-100, default 1
	Priority        int        `json:"priority"`            // Lower is higher priority (for FAILOVER)
	Healthy         bool       `json:"healthy"`
	LastHealthCheck time.Time  `json:"last_health_check"`

	CreatedAt time.Time `json:"created_at"`
}

// GlobalHealthCheckConfig defines how endpoints are probed.
type GlobalHealthCheckConfig struct {
	Protocol       string `json:"protocol"` // HTTP, HTTPS, TCP
	Port           int    `json:"port"`
	Path           string `json:"path,omitempty"` // For HTTP/HTTPS
	IntervalSec    int    `json:"interval_sec"`
	TimeoutSec     int    `json:"timeout_sec"`
	HealthyCount   int    `json:"healthy_count"`
	UnhealthyCount int    `json:"unhealthy_count"`
}
