package node

import (
	"os"
	"path/filepath"
	"sync"
)

// LocalStore manages file storage on the local disk.
type LocalStore struct {
	rootDir string
	mu      sync.RWMutex
}

// NewLocalStore initializes a new local storage backend.
func NewLocalStore(dataDir string) (*LocalStore, error) {
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return nil, err
	}
	return &LocalStore{rootDir: dataDir}, nil
}

// Write saves data to disk. Overwrites if exists.
func (s *LocalStore) Write(bucket, key string, data []byte) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	path := s.getObjectPath(bucket, key)
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}

// Read retrieves data from disk.
func (s *LocalStore) Read(bucket, key string) ([]byte, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return os.ReadFile(s.getObjectPath(bucket, key))
}

// Delete removes data from disk.
func (s *LocalStore) Delete(bucket, key string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	return os.Remove(s.getObjectPath(bucket, key))
}

func (s *LocalStore) getObjectPath(bucket, key string) string {
	return filepath.Join(s.rootDir, bucket, key)
}
