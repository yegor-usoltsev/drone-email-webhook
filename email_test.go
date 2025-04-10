package main

import (
	"context"
	"strings"
	"testing"

	"github.com/drone/drone-go/drone"
	"github.com/drone/drone-go/plugin/webhook"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func setupMailHog(t *testing.T) (string, int, func()) {
	ctx := context.Background()

	req := testcontainers.ContainerRequest{
		Image:        "mailhog/mailhog:latest",
		ExposedPorts: []string{"1025/tcp", "8025/tcp"},
		WaitingFor:   wait.ForListeningPort("1025/tcp"),
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	require.NoError(t, err)

	host, err := container.Host(ctx)
	require.NoError(t, err)

	smtpPort, err := container.MappedPort(ctx, "1025/tcp")
	require.NoError(t, err)

	return host, smtpPort.Int(), func() {
		_ = container.Terminate(ctx)
	}
}

func TestEmailSender_Send(t *testing.T) {
	host, port, cleanup := setupMailHog(t)
	defer cleanup()

	sender := NewEmailSender(Settings{
		EmailSmtpHost: host,
		EmailSmtpPort: port,
		EmailFrom:     "test@example.com",
	})

	tests := []struct {
		name    string
		req     *webhook.Request
		wantErr bool
	}{
		{
			name: "successful_email_with_author_name",
			req: &webhook.Request{
				Build: &drone.Build{
					Number:      1,
					AuthorName:  "Test User",
					Author:      "test",
					AuthorEmail: "test@example.com",
					Message:     "Test commit\nMore details",
					After:       "abcdef1234567890",
					Ref:         "refs/heads/main",
				},
				Repo: &drone.Repo{
					Slug: "test/repo",
				},
				System: &drone.System{
					Host: "drone.example.com",
					Link: "https://drone.example.com",
				},
			},
			wantErr: false,
		},
		{
			name: "successful_email_without_author_name",
			req: &webhook.Request{
				Build: &drone.Build{
					Number:      2,
					Author:      "test",
					AuthorEmail: "test@example.com",
					Message:     "Test commit",
					After:       "abcdef1234567890",
					Ref:         "refs/heads/feature",
				},
				Repo: &drone.Repo{
					Slug: "test/repo",
				},
				System: &drone.System{
					Host: "drone.example.com",
					Link: "https://drone.example.com",
				},
			},
			wantErr: false,
		},
		{
			name: "empty_commit_message",
			req: &webhook.Request{
				Build: &drone.Build{
					Number:      3,
					Author:      "test",
					AuthorEmail: "test@example.com",
					Message:     "",
					After:       "abcdef1234567890",
					Ref:         "refs/heads/main",
				},
				Repo: &drone.Repo{
					Slug: "test/repo",
				},
				System: &drone.System{
					Host: "drone.example.com",
					Link: "https://drone.example.com",
				},
			},
			wantErr: false,
		},
		{
			name: "very_long_commit_message",
			req: &webhook.Request{
				Build: &drone.Build{
					Number:      4,
					Author:      "test",
					AuthorEmail: "test@example.com",
					Message:     strings.Repeat("Very long commit message. ", 100),
					After:       "abcdef1234567890",
					Ref:         "refs/heads/main",
				},
				Repo: &drone.Repo{
					Slug: "test/repo",
				},
				System: &drone.System{
					Host: "drone.example.com",
					Link: "https://drone.example.com",
				},
			},
			wantErr: false,
		},
		{
			name: "missing_author_email",
			req: &webhook.Request{
				Build: &drone.Build{
					Number:  5,
					Author:  "test",
					Message: "Test commit",
					After:   "abcdef1234567890",
					Ref:     "refs/heads/main",
				},
				Repo: &drone.Repo{
					Slug: "test/repo",
				},
				System: &drone.System{
					Host: "drone.example.com",
					Link: "https://drone.example.com",
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := sender.Send(tt.req)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

		})
	}
}

func TestEmailSender_TemplateExecution(t *testing.T) {
	host, port, cleanup := setupMailHog(t)
	defer cleanup()

	s := NewEmailSender(Settings{
		EmailSmtpHost: host,
		EmailSmtpPort: port,
		EmailFrom:     "test@example.com",
	})

	assert.NotNil(t, s.html, "HTML template should be parsed")
	assert.NotNil(t, s.text, "Text template should be parsed")
}

func TestEmailSender_InvalidSMTP(t *testing.T) {
	sender := NewEmailSender(Settings{
		EmailSmtpHost: "nonexistent.example.com",
		EmailSmtpPort: 1025,
		EmailFrom:     "test@example.com",
	})

	req := &webhook.Request{
		Build: &drone.Build{
			Number:      1,
			Author:      "test",
			AuthorEmail: "test@example.com",
			Message:     "Test commit",
			After:       "abcdef1234567890",
			Ref:         "refs/heads/main",
		},
		Repo: &drone.Repo{
			Slug: "test/repo",
		},
		System: &drone.System{
			Host: "drone.example.com",
			Link: "https://drone.example.com",
		},
	}

	err := sender.Send(req)
	assert.Error(t, err)
}
