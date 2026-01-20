package node

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLocalStore(t *testing.T) {
	// Setup temp dir
	tmpDir, err := os.MkdirTemp("", "store-test")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	store, err := NewLocalStore(tmpDir)
	require.NoError(t, err)

	bucket := "mybucket"
	key := "folder/myobject.txt"
	data := []byte("hello world")

	// 1. Write
	err = store.Write(bucket, key, data)
	assert.NoError(t, err)

	// Verify file on disk
	expectedPath := filepath.Join(tmpDir, bucket, key)
	_, err = os.Stat(expectedPath)
	assert.NoError(t, err)

	// 2. Read
	readData, err := store.Read(bucket, key)
	assert.NoError(t, err)
	assert.Equal(t, data, readData)

	// 3. Delete
	err = store.Delete(bucket, key)
	assert.NoError(t, err)

	// Verify deletion
	_, err = os.Stat(expectedPath)
	assert.True(t, os.IsNotExist(err))

	// 4. Read Non-existent
	_, err = store.Read(bucket, key)
	assert.Error(t, err)
	assert.True(t, os.IsNotExist(err))
}
