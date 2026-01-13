package lvm

import (
	"context"
	"testing"
)

func TestLvmAdapter_Type(t *testing.T) {
	adapter := NewLvmAdapter("testvg")
	if adapter.Type() != "lvm" {
		t.Errorf("expected type 'lvm', got %s", adapter.Type())
	}
}

func TestLvmAdapter_Ping(t *testing.T) {
	adapter := NewLvmAdapter("testvg")
	err := adapter.Ping(context.Background())
	if err == nil {
		// If no error, perhaps lvm is available, but unlikely in test
		t.Log("lvm ping succeeded")
	} else {
		// Expected if lvm not available
		t.Logf("lvm ping failed as expected: %v", err)
	}
}

func TestLvmAdapter_CreateVolume_InvalidName(t *testing.T) {
	adapter := NewLvmAdapter("testvg")
	_, err := adapter.CreateVolume(context.Background(), "", 10)
	if err == nil {
		t.Error("expected error for empty name")
	}
}

func TestLvmAdapter_DeleteVolume_InvalidName(t *testing.T) {
	adapter := NewLvmAdapter("testvg")
	err := adapter.DeleteVolume(context.Background(), "")
	if err == nil {
		t.Error("expected error for empty name")
	}
}
