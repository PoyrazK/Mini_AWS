package domain

import (
	"time"

	"github.com/google/uuid"
)

type DeploymentStatus string

const (
	DeploymentStatusScaling  DeploymentStatus = "SCALING"
	DeploymentStatusReady    DeploymentStatus = "READY"
	DeploymentStatusDegraded DeploymentStatus = "DEGRADED"
	DeploymentStatusDeleting DeploymentStatus = "DELETING"
)

type Deployment struct {
	ID           uuid.UUID        `json:"id"`
	UserID       uuid.UUID        `json:"user_id"`
	Name         string           `json:"name"`
	Image        string           `json:"image"`
	Replicas     int              `json:"replicas"`
	CurrentCount int              `json:"current_count"`
	Ports        string           `json:"ports"`
	Status       DeploymentStatus `json:"status"`
	CreatedAt    time.Time        `json:"created_at"`
	UpdatedAt    time.Time        `json:"updated_at"`
}

type DeploymentContainer struct {
	ID           uuid.UUID `json:"id"`
	DeploymentID uuid.UUID `json:"deployment_id"`
	InstanceID   uuid.UUID `json:"instance_id"`
}
