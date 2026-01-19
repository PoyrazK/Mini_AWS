package sshutil

import (
	"bytes"
	"fmt"
	"net"
	"time"

	"golang.org/x/crypto/ssh"
)

// Client represents an SSH client for remote execution.
type Client struct {
	Host string
	User string
	Auth []ssh.AuthMethod
}

// NewClientWithKey constructs an SSH client using a private key.
func NewClientWithKey(host, user, privateKey string) (*Client, error) {
	signer, err := ssh.ParsePrivateKey([]byte(privateKey))
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %w", err)
	}

	return &Client{
		Host: host,
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
	}, nil
}

// Run executes a command and returns its output.
func (c *Client) Run(cmd string) (string, error) {
	config := &ssh.ClientConfig{
		User:            c.User,
		Auth:            c.Auth,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // For cloud instances, might change for production
		Timeout:         10 * time.Second,
	}

	addr := c.Host
	if _, _, err := net.SplitHostPort(addr); err != nil {
		addr = net.JoinHostPort(addr, "22")
	}

	client, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		return "", fmt.Errorf("failed to dial: %w", err)
	}
	defer client.Close()

	session, err := client.NewSession()
	if err != nil {
		return "", fmt.Errorf("failed to create session: %w", err)
	}
	defer session.Close()

	var b bytes.Buffer
	session.Stdout = &b
	session.Stderr = &b

	if err := session.Run(cmd); err != nil {
		return b.String(), fmt.Errorf("failed to run command %q: %w", cmd, err)
	}

	return b.String(), nil
}

// WriteFile writes content to a remote file.
func (c *Client) WriteFile(path string, content []byte, mode string) error {
	config := &ssh.ClientConfig{
		User:            c.User,
		Auth:            c.Auth,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         10 * time.Second,
	}

	addr := c.Host
	if _, _, err := net.SplitHostPort(addr); err != nil {
		addr = net.JoinHostPort(addr, "22")
	}

	client, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		return fmt.Errorf("failed to dial: %w", err)
	}
	defer client.Close()

	session, err := client.NewSession()
	if err != nil {
		return fmt.Errorf("failed to create session: %w", err)
	}
	defer session.Close()

	go func() {
		w, _ := session.StdinPipe()
		defer w.Close()
		fmt.Fprintf(w, "C%s %d %s\n", mode, len(content), path)
		if _, err := w.Write(content); err != nil {
			// This might happen if the pipe closes early
			_ = err
		}
		fmt.Fprint(w, "\x00")
	}()

	if err := session.Run("/usr/bin/scp -t " + path); err != nil {
		return fmt.Errorf("failed to scp: %w", err)
	}

	return nil
}

// WaitForSSH waits for the SSH port to be open and accepting connections.
func (c *Client) WaitForSSH(timeout time.Duration) error {
	start := time.Now()
	for {
		if time.Since(start) > timeout {
			return fmt.Errorf("timed out waiting for SSH on %s", c.Host)
		}

		conn, err := net.DialTimeout("tcp", net.JoinHostPort(c.Host, "22"), 2*time.Second)
		if err == nil {
			conn.Close()
			return nil
		}
		time.Sleep(2 * time.Second)
	}
}
