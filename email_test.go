package main

import (
	"encoding/json"
	"fmt"
	"math/rand/v2"
	"net"
	"net/http"
	"net/mail"
	"net/url"
	"testing"

	"github.com/docker/go-connections/nat"
	"github.com/drone/drone-go/drone"
	"github.com/drone/drone-go/plugin/webhook"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	tc "github.com/testcontainers/testcontainers-go"
	tcWait "github.com/testcontainers/testcontainers-go/wait"
)

func setupMailpit(t *testing.T) *MailpitClient {
	t.Helper()
	container, err := tc.GenericContainer(t.Context(), tc.GenericContainerRequest{
		ContainerRequest: tc.ContainerRequest{
			Image: "axllent/mailpit:latest",
			Env: map[string]string{
				"MP_SMTP_AUTH_ACCEPT_ANY":     "true",
				"MP_SMTP_AUTH_ALLOW_INSECURE": "true",
			},
			ExposedPorts: []string{"1025/tcp", "8025/tcp"},
			WaitingFor:   tcWait.ForHTTP("/readyz").WithPort("8025"),
		},
		Started: true,
	})
	tc.CleanupContainer(t, container)
	require.NoError(t, err)

	host, err := container.Host(t.Context())
	require.NoError(t, err)

	smtpPort, err := container.MappedPort(t.Context(), "1025/tcp")
	require.NoError(t, err)

	httpPort, err := container.MappedPort(t.Context(), "8025/tcp")
	require.NoError(t, err)

	return NewMailpitClient(t, host, smtpPort, httpPort)
}

func buildConfig(mailpit *MailpitClient, fns ...func(*Config)) Config {
	cfg := Config{
		EmailSMTPHost:     mailpit.host,
		EmailSMTPPort:     uint16(mailpit.smtpPort),
		EmailSMTPUsername: "drone@example.com",
		EmailSMTPPassword: "password123",
		EmailFrom:         "ci@example.com",
		EmailCC:           []string{"admin@example.com"},
		EmailBCC:          []string{"security@example.com"},
	}
	for _, fn := range fns {
		fn(&cfg)
	}
	return cfg
}

func buildWebhookRequest(fns ...func(*webhook.Request)) *webhook.Request {
	req := &webhook.Request{
		Event:  webhook.EventBuild,
		Action: webhook.ActionUpdated,
		Repo: &drone.Repo{
			Slug: "test/repo",
		},
		Build: &drone.Build{
			ID:           rand.Int64(),
			Number:       rand.Int64(),
			Status:       "failure",
			Ref:          "refs/heads/main",
			Message:      "test commit",
			After:        "e92d9f39abe709d90e8072b8ec992f2c3a02a07a",
			Author:       "test",
			AuthorName:   "Test User",
			AuthorEmail:  "test@example.com",
			AuthorAvatar: "https://example.com/avatar.png",
		},
		System: &drone.System{
			Host: "drone.example.com",
			Link: "https://drone.example.com",
		},
	}
	for _, fn := range fns {
		fn(req)
	}
	return req
}

func TestEmailSender(t *testing.T) {
	mailpit := setupMailpit(t)

	t.Run("send async", func(t *testing.T) {
		t.Parallel()
		cfg := buildConfig(mailpit)
		emailSender := NewEmailSender(cfg)
		req := buildWebhookRequest()

		emailSender.SendAsync(req)
		emailSender.Shutdown()

		msg := mailpit.FindByBuildNumber(req.Build.Number)
		require.NotNil(t, msg)
		assert.Equal(t, mail.Address{Address: cfg.EmailFrom}, msg.From)
		assert.Equal(t, []mail.Address{{Name: req.Build.AuthorName, Address: req.Build.AuthorEmail}}, msg.To)
		assert.Equal(t, []mail.Address{{Address: cfg.EmailCC[0]}}, msg.Cc)
		assert.Equal(t, []mail.Address{{Address: cfg.EmailBCC[0]}}, msg.Bcc)
		assert.Equal(t, fmt.Sprintf("[%s] Failed build #%d for %s (%s)", req.Repo.Slug, req.Build.Number, req.Build.Ref, req.Build.After[:8]), msg.Subject)
	})

	t.Run("send async with closed sender", func(t *testing.T) {
		t.Parallel()
		cfg := buildConfig(mailpit)
		emailSender := NewEmailSender(cfg)
		req := buildWebhookRequest()

		emailSender.Shutdown()
		emailSender.SendAsync(req)

		msg := mailpit.FindByBuildNumber(req.Build.Number)
		assert.Nil(t, msg)
	})

	t.Run("send", func(t *testing.T) {
		t.Parallel()
		cfg := buildConfig(mailpit)
		emailSender := NewEmailSender(cfg)
		req := buildWebhookRequest()

		err := emailSender.Send(req)
		require.NoError(t, err)

		msg := mailpit.FindByBuildNumber(req.Build.Number)
		assert.NotNil(t, msg)
	})

	t.Run("send with empty author name", func(t *testing.T) {
		t.Parallel()
		cfg := buildConfig(mailpit)
		emailSender := NewEmailSender(cfg)
		req := buildWebhookRequest(func(req *webhook.Request) {
			req.Build.AuthorName = ""
		})

		err := emailSender.Send(req)
		require.NoError(t, err)

		msg := mailpit.FindByBuildNumber(req.Build.Number)
		require.NotNil(t, msg)
		assert.Equal(t, []mail.Address{{Name: req.Build.Author, Address: req.Build.AuthorEmail}}, msg.To)
	})

	t.Run("send with invalid SMTP addr", func(t *testing.T) {
		t.Parallel()
		cfg := buildConfig(mailpit, func(cfg *Config) {
			var lc net.ListenConfig
			l, err := lc.Listen(t.Context(), "tcp", "127.0.0.1:0")
			require.NoError(t, err)
			defer l.Close()
			cfg.EmailSMTPHost = "127.0.0.1"
			cfg.EmailSMTPPort = uint16(l.Addr().(*net.TCPAddr).Port)
		})
		emailSender := NewEmailSender(cfg)
		req := buildWebhookRequest()

		err := emailSender.Send(req)
		require.Error(t, err)

		msg := mailpit.FindByBuildNumber(req.Build.Number)
		assert.Nil(t, msg)
	})

	t.Run("shutdown", func(t *testing.T) {
		t.Parallel()
		cfg := buildConfig(mailpit)
		emailSender := NewEmailSender(cfg)
		assert.NotPanics(t, func() { emailSender.Shutdown() })
		assert.NotPanics(t, func() { emailSender.Shutdown() })
	})
}

type MailpitClient struct {
	t         *testing.T
	host      string
	smtpPort  int
	httpPort  int
	searchURL string
}

func NewMailpitClient(t *testing.T, host string, smtpPort, httpPort nat.Port) *MailpitClient {
	t.Helper()
	return &MailpitClient{
		t:         t,
		host:      host,
		smtpPort:  smtpPort.Int(),
		httpPort:  httpPort.Int(),
		searchURL: "http://" + net.JoinHostPort(host, httpPort.Port()) + "/api/v1/search",
	}
}

func (m *MailpitClient) FindByBuildNumber(buildNumber int64) *MessageSummary {
	req, err := http.NewRequestWithContext(m.t.Context(), http.MethodGet, m.searchURL, http.NoBody)
	require.NoError(m.t, err)
	req.URL.RawQuery = url.Values{
		"query": []string{fmt.Sprintf(`subject:"%d"`, buildNumber)},
	}.Encode()

	res, err := http.DefaultClient.Do(req)
	require.NoError(m.t, err)
	defer res.Body.Close()

	var body MessagesSummaryResponse
	err = json.NewDecoder(res.Body).Decode(&body)
	require.NoError(m.t, err)
	if len(body.Messages) == 0 {
		return nil
	}
	return &body.Messages[0]
}

type MessagesSummaryResponse struct {
	Messages []MessageSummary `json:"messages"`
}

//nolint:tagliatelle
type MessageSummary struct {
	From    mail.Address   `json:"From"`
	To      []mail.Address `json:"To"`
	Cc      []mail.Address `json:"Cc"`
	Bcc     []mail.Address `json:"Bcc"`
	Subject string         `json:"Subject"`
}
