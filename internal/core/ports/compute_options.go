package ports

// CreateInstanceOptions contains parameters for creating a new compute instance
type CreateInstanceOptions struct {
	Name        string
	ImageName   string
	Ports       []string
	NetworkID   string
	VolumeBinds []string
	Env         []string
	Cmd         []string
}
