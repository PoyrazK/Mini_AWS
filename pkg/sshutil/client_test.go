package sshutil

import (
	"bufio"
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io"
	"net"
	"strings"
	"sync"
	"testing"
	"time"

	"golang.org/x/crypto/ssh"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// generateTestKey generates a private key for testing
func generateTestKey(t *testing.T) string {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)

	keyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(key),
	})

	return string(keyPEM)
}

const (
	testLocalhostSSH = "localhost:22"
	testLoopbackAddr = "127.0.0.1:0"
)

func TestNewClientWithKey(t *testing.T) {
	// 1. Valid Key
	privKey := generateTestKey(t)
	client, err := NewClientWithKey(testLocalhostSSH, "user", privKey)
	require.NoError(t, err)
	assert.NotNil(t, client)
	assert.Equal(t, testLocalhostSSH, client.Host)
	assert.Equal(t, "user", client.User)
	assert.NotEmpty(t, client.Auth)

	// 2. Invalid Key
	_, err = NewClientWithKey(testLocalhostSSH, "user", "invalid-key")
	require.Error(t, err)
}

func TestWaitForSSH(t *testing.T) {
	// Start a dummy TCP server
	l, err := net.Listen("tcp", testLoopbackAddr)
	require.NoError(t, err)
	defer l.Close()

	port := l.Addr().(*net.TCPAddr).Port
	host := fmt.Sprintf("127.0.0.1:%d", port)

	client := &Client{Host: host}

	// Should connect successfully
	err = client.WaitForSSH(context.Background(), 2*time.Second)
	require.NoError(t, err)
}

func TestWaitForSSHTimeout(t *testing.T) {
	// Pick a random port (hopefully unused)
	client := &Client{Host: "127.0.0.1:54321"} // Unlikely to be a valid SSH server immediately

	// Should timeout
	// Use small timeout for test speed
	err := client.WaitForSSH(context.Background(), 100*time.Millisecond)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "timed out")
}

func TestWaitForSSHContextCanceled(t *testing.T) {
	client := &Client{Host: "127.0.0.1:54321"}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := client.WaitForSSH(ctx, 2*time.Second)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "timed out")
}

func TestWaitForSSHHostWithoutPortTimeout(t *testing.T) {
	client := &Client{Host: "127.0.0.1"}

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	err := client.WaitForSSH(ctx, 2*time.Second)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "timed out")
}

func TestRunConnectionRefused(t *testing.T) {
	// Ensure we pick a port that rejects connection
	l, err := net.Listen("tcp", testLoopbackAddr)
	require.NoError(t, err)
	port := l.Addr().(*net.TCPAddr).Port
	l.Close() // Close immediately to ensure connection refused

	host := fmt.Sprintf("127.0.0.1:%d", port)
	privKey := generateTestKey(t)
	client, _ := NewClientWithKey(host, "user", privKey)

	_, err = client.Run(context.Background(), "echo hello")
	require.Error(t, err)
	// Error message format depends on OS, usually "connection refused" or "dial tcp"
}

func TestRunContextTimeout(t *testing.T) {
	privKey := generateTestKey(t)
	client, err := NewClientWithKey("127.0.0.1:0", "user", privKey)
	require.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	_, err = client.Run(ctx, "echo hello")
	require.Error(t, err)
}

func TestRunHostWithoutPort(t *testing.T) {
	privKey := generateTestKey(t)
	client, err := NewClientWithKey("127.0.0.1", "user", privKey)
	require.NoError(t, err)

	_, err = client.Run(context.Background(), "echo hello")
	require.Error(t, err)
}

// TODO: A full SSH server mock for Run and WriteFile would be better but significantly more complex.
// For "Phase 1 Quick Wins", validating the Client logic, Key parsing, and Network dialing is a good start.

func TestWriteFileConnectionRefused(t *testing.T) {
	// Pick a random port (hopefully unused)
	client := &Client{Host: testLoopbackAddr}

	err := client.WriteFile(context.Background(), "/tmp/test", []byte("data"), "0644")
	require.Error(t, err)
	// Expect dial error
}

func TestWriteFileContextTimeout(t *testing.T) {
	client := &Client{Host: "127.0.0.1:0"}

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	err := client.WriteFile(ctx, "/tmp/test", []byte("data"), "0644")
	require.Error(t, err)
}

func TestWriteFileHostWithoutPort(t *testing.T) {
	privKey := generateTestKey(t)
	client, err := NewClientWithKey("127.0.0.1", "user", privKey)
	require.NoError(t, err)

	err = client.WriteFile(context.Background(), "/tmp/test", []byte("data"), "0644")
	require.Error(t, err)
}

func TestRunSuccess(t *testing.T) {
	addr, stop := startTestSSHServer(t)
	defer stop()

	privKey := generateTestKey(t)
	client, err := NewClientWithKey(addr, "user", privKey)
	require.NoError(t, err)

	out, err := client.Run(context.Background(), "echo hello")
	require.NoError(t, err)
	assert.Contains(t, out, "hello")
}

func TestWriteFileSuccess(t *testing.T) {
	addr, stop := startTestSSHServer(t)
	defer stop()

	privKey := generateTestKey(t)
	client, err := NewClientWithKey(addr, "user", privKey)
	require.NoError(t, err)

	err = client.WriteFile(context.Background(), "/tmp/test.txt", []byte("data"), "0644")
	require.NoError(t, err)
}

func TestWriteFileScpError(t *testing.T) {
	addr, stop := startTestSSHServerWithHandler(t, func(cmd string, ch ssh.Channel) error {
		if strings.HasPrefix(cmd, "/usr/bin/scp -t ") {
			return fmt.Errorf("scp failed")
		}
		if cmd == "echo hello" {
			_, err := ch.Write([]byte("hello\n"))
			return err
		}
		return fmt.Errorf("unsupported command: %s", cmd)
	})
	defer stop()

	privKey := generateTestKey(t)
	client, err := NewClientWithKey(addr, "user", privKey)
	require.NoError(t, err)

	err = client.WriteFile(context.Background(), "/tmp/test.txt", []byte("data"), "0644")
	require.Error(t, err)
}

func startTestSSHServer(t *testing.T) (string, func()) {
	return startTestSSHServerWithHandler(t, handleExecCommand)
}

func startTestSSHServerWithHandler(t *testing.T, handler func(cmd string, ch ssh.Channel) error) (string, func()) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)

	signer, err := newTestSigner()
	require.NoError(t, err)

	config := &ssh.ServerConfig{NoClientAuth: true}
	config.AddHostKey(signer)

	var wg sync.WaitGroup
	stop := make(chan struct{})

	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				select {
				case <-stop:
					return
				default:
					return
				}
			}
			wg.Add(1)
			go func(c net.Conn) {
				defer wg.Done()
				handleSSHConn(t, c, config, handler)
			}(conn)
		}
	}()

	return listener.Addr().String(), func() {
		close(stop)
		listener.Close()
		wg.Wait()
	}
}

func newTestSigner() (ssh.Signer, error) {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, err
	}
	return ssh.NewSignerFromKey(key)
}

func handleSSHConn(t *testing.T, conn net.Conn, config *ssh.ServerConfig, handler func(cmd string, ch ssh.Channel) error) {
	sshConn, chans, reqs, err := ssh.NewServerConn(conn, config)
	if err != nil {
		return
	}
	defer sshConn.Close()
	go ssh.DiscardRequests(reqs)

	for newChannel := range chans {
		if newChannel.ChannelType() != "session" {
			newChannel.Reject(ssh.UnknownChannelType, "unknown channel type")
			continue
		}

		channel, requests, err := newChannel.Accept()
		if err != nil {
			continue
		}

		go func(ch ssh.Channel, in <-chan *ssh.Request) {
			defer ch.Close()
			for req := range in {
				if req.Type != "exec" {
					req.Reply(false, nil)
					continue
				}

				var payload struct{ Command string }
				ssh.Unmarshal(req.Payload, &payload)
				req.Reply(true, nil)

				status := uint32(0)
				if err := handler(payload.Command, ch); err != nil {
					status = 1
				}

				_, _ = ch.SendRequest("exit-status", false, ssh.Marshal(struct{ Status uint32 }{Status: status}))
				return
			}
		}(channel, requests)
	}
}

func handleExecCommand(cmd string, ch ssh.Channel) error {
	if cmd == "echo hello" {
		_, err := ch.Write([]byte("hello\n"))
		return err
	}

	if len(cmd) >= len("/usr/bin/scp -t ") && cmd[:len("/usr/bin/scp -t ")] == "/usr/bin/scp -t " {
		return handleScp(ch)
	}

	return fmt.Errorf("unsupported command: %s", cmd)
}

func handleScp(r io.Reader) error {
	br := bufio.NewReader(r)
	line, err := br.ReadString('\n')
	if err != nil {
		return err
	}

	var mode string
	var size int
	var filename string
	_, err = fmt.Sscanf(line, "C%s %d %s", &mode, &size, &filename)
	if err != nil {
		return err
	}

	if size > 0 {
		if _, err := io.CopyN(io.Discard, br, int64(size)); err != nil {
			return err
		}
	}

	_, err = br.ReadByte() // null byte
	return err
}
