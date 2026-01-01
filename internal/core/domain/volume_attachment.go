package domain

// VolumeAttachment represents a request to attach a volume to an instance.
type VolumeAttachment struct {
	VolumeIDOrName string `json:"volume_id"`
	MountPath      string `json:"mount_path"`
}
