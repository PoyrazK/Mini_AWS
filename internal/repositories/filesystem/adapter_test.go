package filesystem

import (
	"bytes"
	"context"
	"io"
	"os"
	"testing"

	"github.com/poyrazk/thecloud/internal/errors"
)

const (
	pathTraversalKey          = "../../../etc/passwd"
	pathTraversalExpectedErr  = "expected error for path traversal"
	pathTraversalInvalidInput = "expected InvalidInput, got %v"
)

func TestLocalFileStoreWriteReadDelete(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "filestore_test")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.RemoveAll(tempDir) }()

	store, err := NewLocalFileStore(tempDir)
	if err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()
	bucket := "testbucket"
	key := "testkey"
	data := []byte("hello world")

	// Write
	n, err := store.Write(ctx, bucket, key, bytes.NewReader(data))
	if err != nil {
		t.Fatal(err)
	}
	if n != int64(len(data)) {
		t.Fatalf("expected %d bytes written, got %d", len(data), n)
	}

	// Read
	r, err := store.Read(ctx, bucket, key)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = r.Close() }()
	readData, err := io.ReadAll(r)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(readData, data) {
		t.Fatalf("expected %s, got %s", data, readData)
	}

	// Delete
	err = store.Delete(ctx, bucket, key)
	if err != nil {
		t.Fatal(err)
	}

	// Read after delete should fail
	_, err = store.Read(ctx, bucket, key)
	if err == nil {
		t.Fatal("expected error after delete")
	}
	if !errors.Is(err, errors.ObjectNotFound) {
		t.Fatalf("expected ObjectNotFound, got %v", err)
	}
}

func TestLocalFileStorePathTraversal(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "filestore_test")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.RemoveAll(tempDir) }()

	store, err := NewLocalFileStore(tempDir)
	if err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()

	// Try path traversal
	_, err = store.Write(ctx, "bucket", pathTraversalKey, bytes.NewReader([]byte("bad")))
	if err == nil {
		t.Fatal(pathTraversalExpectedErr)
	}
	if !errors.Is(err, errors.InvalidInput) {
		t.Fatalf(pathTraversalInvalidInput, err)
	}

	_, err = store.Read(ctx, "bucket", pathTraversalKey)
	if err == nil {
		t.Fatal(pathTraversalExpectedErr)
	}
	if !errors.Is(err, errors.InvalidInput) {
		t.Fatalf(pathTraversalInvalidInput, err)
	}

	err = store.Delete(ctx, "bucket", pathTraversalKey)
	if err == nil {
		t.Fatal(pathTraversalExpectedErr)
	}
	if !errors.Is(err, errors.InvalidInput) {
		t.Fatalf(pathTraversalInvalidInput, err)
	}
}

func TestLocalFileStoreAssemble(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "filestore_test")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.RemoveAll(tempDir) }()

	store, err := NewLocalFileStore(tempDir)
	if err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()
	bucket := "testbucket"

	// Create parts
	p1 := "part1"
	p2 := "part2"
	_, _ = store.Write(ctx, bucket, p1, bytes.NewReader([]byte("hello ")))
	_, _ = store.Write(ctx, bucket, p2, bytes.NewReader([]byte("world")))

	// Assemble
	target := "target"
	size, err := store.Assemble(ctx, bucket, target, []string{p1, p2})
	if err != nil {
		t.Fatal(err)
	}
	if size != 11 {
		t.Fatalf("expected 11 bytes, got %d", size)
	}

	// Verify content
	r, _ := store.Read(ctx, bucket, target)
	data, _ := io.ReadAll(r)
	if string(data) != "hello world" {
		t.Fatalf("expected 'hello world', got '%s'", string(data))
	}
}
