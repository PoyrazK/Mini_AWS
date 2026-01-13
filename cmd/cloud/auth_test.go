package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadConfig_WhenFileMissing_ReturnsEmptyString(t *testing.T) {
	t.Setenv("HOME", t.TempDir())

	got := loadConfig()
	if got != "" {
		t.Fatalf("expected empty string when config missing, got %q", got)
	}
}

func TestSaveAndLoadConfig_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("HOME", dir)

	saveConfig("abc123")

	got := loadConfig()
	if got != "abc123" {
		t.Fatalf("expected API key %q, got %q", "abc123", got)
	}

	// Ensure the file was written where we expect.
	path := filepath.Join(dir, ".cloud", "config.json")
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("expected config file at %s: %v", path, err)
	}
}

func TestLoadConfig_WhenInvalidJSON_ReturnsEmptyString(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("HOME", dir)

	path := filepath.Join(dir, ".cloud", "config.json")
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	if err := os.WriteFile(path, []byte("{not json"), 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}

	got := loadConfig()
	if got != "" {
		t.Fatalf("expected empty string on invalid json, got %q", got)
	}
}
