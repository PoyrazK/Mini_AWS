package k8s

import (
	"context"
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/poyrazk/thecloud/internal/core/ports"
	"golang.org/x/crypto/ssh"
)

// NodeExecutor defines an interface for executing commands on a cluster node.
type NodeExecutor interface {
	Run(ctx context.Context, cmd string) (string, error)
}

// ServiceExecutor uses the InstanceService.Exec (for Docker backend).
type ServiceExecutor struct {
	svc    ports.InstanceService
	instID uuid.UUID
}

func NewServiceExecutor(svc ports.InstanceService, instID uuid.UUID) *ServiceExecutor {
	return &ServiceExecutor{svc: svc, instID: instID}
}

func (e *ServiceExecutor) Run(ctx context.Context, cmd string) (string, error) {
	return e.svc.Exec(ctx, e.instID.String(), []string{"sh", "-c", cmd})
}

// SSHExecutor uses SSH to run commands on a node.
type SSHExecutor struct {
	ip   string
	user string
	key  string
}

func NewSSHExecutor(ip, user, key string) *SSHExecutor {
	return &SSHExecutor{ip: ip, user: user, key: key}
}

func (e *SSHExecutor) Run(ctx context.Context, cmd string) (string, error) {
	signer, err := ssh.ParsePrivateKey([]byte(e.key))
	if err != nil {
		return "", fmt.Errorf("failed to parse private key: %w", err)
	}

	config := &ssh.ClientConfig{
		User: e.user,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         10 * time.Second,
	}

	addr := net.JoinHostPort(e.ip, "22")
	client, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		return "", fmt.Errorf("failed to dial ssh: %w", err)
	}
	defer client.Close()

	session, err := client.NewSession()
	if err != nil {
		return "", fmt.Errorf("failed to create session: %w", err)
	}
	defer session.Close()

	var stdout, stderr strings.Builder
	session.Stdout = &stdout
	session.Stderr = &stderr

	err = session.Run(cmd)
	if err != nil {
		return stdout.String() + stderr.String(), fmt.Errorf("command failed: %w (stderr: %s)", err, stderr.String())
	}

	return stdout.String(), nil
}
