package main

import (
	"net"
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
	t.Helper()
	ctx := t.Context()

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

func getFreePort(t *testing.T) int {
	t.Helper()
	listener, err := net.Listen("tcp", "localhost:0")
	require.NoError(t, err)
	defer listener.Close()
	return listener.Addr().(*net.TCPAddr).Port
}

func TestEmailSender_Send(t *testing.T) {
	t.Parallel()
	host, port, cleanup := setupMailHog(t)
	t.Cleanup(cleanup)

	tests := []struct {
		name     string
		settings Settings
		req      *webhook.Request
		wantErr  bool
		useAuth  bool
	}{
		{
			name: "successful_email_with_auth",
			settings: Settings{
				EmailSMTPHost:     host,
				EmailSMTPPort:     port,
				EmailFrom:         "test@example.com",
				EmailSMTPUsername: "test-user",
				EmailSMTPPassword: "test-pass",
			},
			useAuth: true,
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
			name: "successful_email_without_auth",
			settings: Settings{
				EmailSMTPHost: host,
				EmailSMTPPort: port,
				EmailFrom:     "test@example.com",
			},
			useAuth: false,
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
			settings: Settings{
				EmailSMTPHost: host,
				EmailSMTPPort: port,
				EmailFrom:     "test@example.com",
			},
			useAuth: false,
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
			settings: Settings{
				EmailSMTPHost: host,
				EmailSMTPPort: port,
				EmailFrom:     "test@example.com",
			},
			useAuth: false,
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
			settings: Settings{
				EmailSMTPHost: host,
				EmailSMTPPort: port,
				EmailFrom:     "test@example.com",
			},
			useAuth: false,
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
		{
			name: "short_commit_hash",
			settings: Settings{
				EmailSMTPHost: host,
				EmailSMTPPort: port,
				EmailFrom:     "test@example.com",
			},
			useAuth: false,
			req: &webhook.Request{
				Build: &drone.Build{
					Number:      6,
					Author:      "test",
					AuthorEmail: "test@example.com",
					Message:     "Test commit",
					After:       "abc",
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
			name: "with_author_avatar",
			settings: Settings{
				EmailSMTPHost: host,
				EmailSMTPPort: port,
				EmailFrom:     "test@example.com",
			},
			useAuth: false,
			req: &webhook.Request{
				Build: &drone.Build{
					Number:       7,
					Author:       "test",
					AuthorEmail:  "test@example.com",
					AuthorAvatar: "https://example.com/avatar.jpg",
					Message:      "Test commit",
					After:        "abcdef1234567890",
					Ref:          "refs/heads/main",
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			sender := NewEmailSender(tt.settings)
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
	t.Parallel()
	host, port, cleanup := setupMailHog(t)
	t.Cleanup(cleanup)

	s := NewEmailSender(Settings{
		EmailSMTPHost: host,
		EmailSMTPPort: port,
		EmailFrom:     "test@example.com",
	})

	assert.NotNil(t, s.html, "HTML template should be parsed")
	assert.NotNil(t, s.text, "Text template should be parsed")
}

func TestEmailSender_InvalidSMTP(t *testing.T) {
	t.Parallel()
	sender := NewEmailSender(Settings{
		EmailSMTPHost: "localhost",
		EmailSMTPPort: getFreePort(t),
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
