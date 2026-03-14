package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewConfigFromEnv(t *testing.T) {
	t.Setenv("DRONE_SECRET", "test-secret")
	t.Setenv("DRONE_SERVER_HOST", "127.0.0.1")
	t.Setenv("DRONE_SERVER_PORT", "8080")
	t.Setenv("DRONE_EMAIL_SMTP_HOST", "smtp.example.com")
	t.Setenv("DRONE_EMAIL_SMTP_PORT", "587")
	t.Setenv("DRONE_EMAIL_SMTP_USERNAME", "test@example.com")
	t.Setenv("DRONE_EMAIL_SMTP_PASSWORD", "password123")
	t.Setenv("DRONE_EMAIL_FROM", "drone@example.com")
	t.Setenv("DRONE_EMAIL_CC", "admin1@example.com,admin2@example.com")
	t.Setenv("DRONE_EMAIL_BCC", "security1@example.com,security2@example.com")

	actual, err := NewConfigFromEnv()

	require.NoError(t, err)
	assert.Equal(t, Config{
		Secret:            "test-secret",
		ServerHost:        "127.0.0.1",
		ServerPort:        8080,
		EmailSMTPHost:     "smtp.example.com",
		EmailSMTPPort:     587,
		EmailSMTPUsername: "test@example.com",
		EmailSMTPPassword: "password123",
		EmailFrom:         "drone@example.com",
		EmailCC:           []string{"admin1@example.com", "admin2@example.com"},
		EmailBCC:          []string{"security1@example.com", "security2@example.com"},
	}, actual)
}

func TestNewConfigFromEnv_Defaults(t *testing.T) {
	t.Setenv("DRONE_SECRET", "test-secret")

	cfg, err := NewConfigFromEnv()

	require.NoError(t, err)
	assert.Equal(t, "0.0.0.0", cfg.ServerHost)
	assert.Equal(t, uint16(3000), cfg.ServerPort)
	assert.Equal(t, "localhost", cfg.EmailSMTPHost)
	assert.Equal(t, uint16(25), cfg.EmailSMTPPort)
	assert.Equal(t, "drone@localhost", cfg.EmailFrom)
}

func TestNewConfigFromEnv_Errors(t *testing.T) {
	t.Run("missing required field", func(t *testing.T) {
		t.Parallel()
		_, err := NewConfigFromEnv()
		assert.Error(t, err)
	})

	t.Run("invalid server port", func(t *testing.T) {
		t.Setenv("DRONE_SECRET", "test-secret")
		t.Setenv("DRONE_SERVER_PORT", "invalid")
		_, err := NewConfigFromEnv()
		assert.Error(t, err)
	})

	t.Run("invalid email SMTP port", func(t *testing.T) {
		t.Setenv("DRONE_SECRET", "test-secret")
		t.Setenv("DRONE_EMAIL_SMTP_PORT", "invalid")
		_, err := NewConfigFromEnv()
		assert.Error(t, err)
	})
}
