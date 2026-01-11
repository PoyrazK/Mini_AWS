package domain

import (
	"github.com/google/uuid"
)

type ProvisionJob struct {
	InstanceID uuid.UUID          `json:"instance_id"`
	UserID     uuid.UUID          `json:"user_id"`
	Volumes    []VolumeAttachment `json:"volumes"`
}
