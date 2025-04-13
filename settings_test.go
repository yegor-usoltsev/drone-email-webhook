package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewSettingsFromEnv(t *testing.T) {
	tests := []struct {
		name     string
		envVars  map[string]string
		wantErr  bool
		expected Settings
	}{
		{
			name: "all_settings_provided",
			envVars: map[string]string{
				"DRONE_SECRET":              "test-secret",
				"DRONE_SERVER_HOST":         "0.0.0.0",
				"DRONE_SERVER_PORT":         "3000",
				"DRONE_EMAIL_SMTP_HOST":     "smtp.example.com",
				"DRONE_EMAIL_SMTP_PORT":     "587",
				"DRONE_EMAIL_SMTP_USERNAME": "test-user",
				"DRONE_EMAIL_SMTP_PASSWORD": "test-pass",
				"DRONE_EMAIL_FROM":          "drone@example.com",
			},
			wantErr: false,
			expected: Settings{
				Secret:            "test-secret",
				ServerHost:        "0.0.0.0",
				ServerPort:        3000,
				EmailSMTPHost:     "smtp.example.com",
				EmailSMTPPort:     587,
				EmailSMTPUsername: "test-user",
				EmailSMTPPassword: "test-pass",
				EmailFrom:         "drone@example.com",
			},
		},
		{
			name: "use_defaults_without_auth",
			envVars: map[string]string{
				"DRONE_SECRET": "test-secret",
			},
			wantErr: false,
			expected: Settings{
				Secret:        "test-secret",
				ServerHost:    "0.0.0.0",
				ServerPort:    3000,
				EmailSMTPHost: "localhost",
				EmailSMTPPort: 25,
				EmailFrom:     "drone@localhost",
			},
		},
		{
			name: "with_auth_but_no_password",
			envVars: map[string]string{
				"DRONE_SECRET":              "test-secret",
				"DRONE_EMAIL_SMTP_USERNAME": "test-user",
			},
			wantErr: false,
			expected: Settings{
				Secret:            "test-secret",
				ServerHost:        "0.0.0.0",
				ServerPort:        3000,
				EmailSMTPHost:     "localhost",
				EmailSMTPPort:     25,
				EmailSMTPUsername: "test-user",
				EmailFrom:         "drone@localhost",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for k, v := range tt.envVars {
				t.Setenv(k, v)
			}

			if tt.wantErr {
				assert.Panics(t, func() {
					NewSettingsFromEnv()
				})
			} else {
				settings := NewSettingsFromEnv()
				assert.Equal(t, tt.expected, settings)
			}
		})
	}
}
