package main

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
)

func TestVersionCommandPrintsVersion(t *testing.T) {
	t.Cleanup(func() {
		rootCmd.SetArgs(nil)
		rootCmd.SetOut(nil)
		rootCmd.SetErr(nil)
	})

	buf := &bytes.Buffer{}
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"version"})

	output := captureStdout(t, func() {
		if err := rootCmd.Execute(); err != nil {
			t.Fatalf("execute: %v", err)
		}
	})
	expected := "cloud v" + version
	if !strings.Contains(output, expected) {
		t.Fatalf("expected output to contain %q, got %q", expected, output)
	}
}

func captureStdout(t *testing.T, fn func()) string {
	t.Helper()

	old := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("pipe: %v", err)
	}
	os.Stdout = w

	fn()

	if err := w.Close(); err != nil {
		t.Fatalf("close stdout pipe: %v", err)
	}

	os.Stdout = old

	var buf bytes.Buffer
	if _, err := io.Copy(&buf, r); err != nil {
		t.Fatalf("copy stdout: %v", err)
	}

	return buf.String()
}
