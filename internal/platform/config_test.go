package platform

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewConfig(t *testing.T) {
	// Save original env and restore after test
	originalPort := os.Getenv("PORT")
	defer os.Setenv("PORT", originalPort)

	t.Run("Default values", func(t *testing.T) {
		os.Unsetenv("PORT")
		cfg, err := NewConfig()
		assert.NoError(t, err)
		assert.Equal(t, "8080", cfg.Port)
	})

	t.Run("Env override", func(t *testing.T) {
		os.Setenv("PORT", "9090")
		cfg, err := NewConfig()
		assert.NoError(t, err)
		assert.Equal(t, "9090", cfg.Port)
	})
}

func TestGetEnv(t *testing.T) {
	t.Run("Existing env", func(t *testing.T) {
		os.Setenv("TEST_KEY", "test_value")
		defer os.Unsetenv("TEST_KEY")
		assert.Equal(t, "test_value", getEnv("TEST_KEY", "fallback"))
	})

	t.Run("Fallback value", func(t *testing.T) {
		os.Unsetenv("NON_EXISTENT_KEY")
		assert.Equal(t, "fallback", getEnv("NON_EXISTENT_KEY", "fallback"))
	})
}
