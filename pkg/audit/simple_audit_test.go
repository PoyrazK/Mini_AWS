package simpleaudit

import (
	"bytes"
	"context"
	"log/slog"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/poyrazk/thecloud/internal/core/domain"
	"github.com/stretchr/testify/assert"
)

func TestSimpleAuditLogger_Log(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(&buf, nil))
	audit := NewSimpleAuditLogger(logger)

	userID := uuid.New()
	entry := &domain.AuditLog{
		ID:           uuid.New(),
		UserID:       userID,
		Action:       "CREATE",
		ResourceType: "instance",
		ResourceID:   "123",
		IPAddress:    "127.0.0.1",
		UserAgent:    "test-agent",
		Details:      map[string]interface{}{"foo": "bar"},
		CreatedAt:    time.Now(),
	}

	err := audit.Log(context.Background(), entry)
	assert.NoError(t, err)

	logOutput := buf.String()
	assert.Contains(t, logOutput, "AUDIT_LOG")
	assert.Contains(t, logOutput, "security_audit")
	assert.Contains(t, logOutput, "CREATE")
	assert.Contains(t, logOutput, userID.String())
	assert.Contains(t, logOutput, "instance")
	assert.Contains(t, logOutput, "123")
	assert.Contains(t, logOutput, "127.0.0.1")
	assert.Contains(t, logOutput, "test-agent")
	assert.Contains(t, logOutput, "foo")
	assert.Contains(t, logOutput, "bar")
}

func TestSimpleAuditLogger_Log_Anonymous(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(&buf, nil))
	audit := NewSimpleAuditLogger(logger)

	entry := &domain.AuditLog{
		ID:           uuid.New(),
		UserID:       uuid.Nil,
		Action:       "LOGIN",
		ResourceType: "auth",
		CreatedAt:    time.Now(),
	}

	err := audit.Log(context.Background(), entry)
	assert.NoError(t, err)

	logOutput := buf.String()
	assert.Contains(t, logOutput, "anonymous")
}

func TestSimpleAuditLogger_Log_EmptyDetails(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(&buf, nil))
	audit := NewSimpleAuditLogger(logger)

	entry := &domain.AuditLog{
		ID:           uuid.New(),
		Action:       "TEST",
		ResourceType: "test",
		Details:      map[string]interface{}{},
	}

	err := audit.Log(context.Background(), entry)
	assert.NoError(t, err)

	logOutput := buf.String()
	// slog with JSON handler handles empty map gracefully
	assert.Contains(t, logOutput, "details")
}
