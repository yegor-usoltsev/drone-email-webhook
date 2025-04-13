package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/99designs/httpsignatures-go"
	"github.com/drone/drone-go/drone"
	"github.com/drone/drone-go/plugin/webhook"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockEmailSender struct {
	mock.Mock
}

func NewMockEmailSender() *MockEmailSender {
	return &MockEmailSender{}
}

func (m *MockEmailSender) Send(req *webhook.Request) error {
	args := m.Called(req)
	return args.Error(0)
}

var (
	errEmailSendingFailed = errors.New("email sending failed")
)

func TestWebhookHandler_ServeHTTP(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name           string
		secret         string
		request        *webhook.Request
		invalidJSON    bool
		signRequest    bool
		expectedStatus int
		emailSent      bool
		emailError     error
	}{
		{
			name:   "valid_failed_build",
			secret: "secret123",
			request: &webhook.Request{
				Event:  webhook.EventBuild,
				Action: webhook.ActionUpdated,
				Build: &drone.Build{
					Status:     "failure",
					ID:         1,
					Number:     1,
					After:      "abcdef1234567890",
					Author:     "test@example.com",
					AuthorName: "Test User",
					Message:    "Test commit",
				},
				Repo: &drone.Repo{
					Slug: "test/repo",
				},
				System: &drone.System{
					Host: "drone.example.com",
					Link: "https://drone.example.com",
				},
			},
			signRequest:    true,
			expectedStatus: http.StatusNoContent,
			emailSent:      true,
			emailError:     nil,
		},
		{
			name:   "successful_build_no_email",
			secret: "secret123",
			request: &webhook.Request{
				Event:  webhook.EventBuild,
				Action: webhook.ActionUpdated,
				Build: &drone.Build{
					Status: "success",
				},
			},
			signRequest:    true,
			expectedStatus: http.StatusNoContent,
			emailSent:      false,
		},
		{
			name:           "missing_signature",
			secret:         "secret123",
			request:        &webhook.Request{},
			signRequest:    false,
			expectedStatus: http.StatusBadRequest,
			emailSent:      false,
		},
		{
			name:           "invalid_signature",
			secret:         "secret123",
			request:        &webhook.Request{},
			signRequest:    true,
			expectedStatus: http.StatusBadRequest,
			emailSent:      false,
		},
		{
			name:           "invalid_json_request",
			secret:         "secret123",
			request:        nil,
			invalidJSON:    true,
			signRequest:    true,
			expectedStatus: http.StatusBadRequest,
			emailSent:      false,
		},
		{
			name:   "different_event_type",
			secret: "secret123",
			request: &webhook.Request{
				Event:  "push",
				Action: webhook.ActionUpdated,
				Build: &drone.Build{
					Status: "failure",
				},
			},
			signRequest:    true,
			expectedStatus: http.StatusNoContent,
			emailSent:      false,
		},
		{
			name:   "different_action_type",
			secret: "secret123",
			request: &webhook.Request{
				Event:  webhook.EventBuild,
				Action: "created",
				Build: &drone.Build{
					Status: "failure",
				},
			},
			signRequest:    true,
			expectedStatus: http.StatusNoContent,
			emailSent:      false,
		},
		{
			name:   "email_sending_error",
			secret: "secret123",
			request: &webhook.Request{
				Event:  webhook.EventBuild,
				Action: webhook.ActionUpdated,
				Build: &drone.Build{
					Status:     "failure",
					ID:         1,
					Number:     1,
					After:      "abcdef1234567890",
					Author:     "test@example.com",
					AuthorName: "Test User",
					Message:    "Test commit",
				},
				Repo: &drone.Repo{
					Slug: "test/repo",
				},
				System: &drone.System{
					Host: "drone.example.com",
					Link: "https://drone.example.com",
				},
			},
			signRequest:    true,
			expectedStatus: http.StatusNoContent,
			emailSent:      true,
			emailError:     errEmailSendingFailed,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			mockEmailSender := NewMockEmailSender()
			if tt.emailSent {
				mockEmailSender.On("Send", tt.request).Return(tt.emailError)
			}

			handler := NewWebhookHandler(Settings{Secret: tt.secret}, mockEmailSender)

			var body []byte
			if tt.invalidJSON {
				body = []byte(`{invalid json`)
			} else {
				body, _ = json.Marshal(tt.request)
			}
			req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(body))

			if tt.signRequest {
				keyID := "drone-email-webhook"
				signer := httpsignatures.NewSigner(httpsignatures.AlgorithmHmacSha256, "date", "content-length")
				req.Header.Set("Date", time.Now().UTC().Format(http.TimeFormat))
				req.Header.Set("Content-Length", "0")
				secret := tt.secret
				if tt.name == "invalid_signature" {
					secret = "wrong_secret"
				}
				_ = signer.SignRequest(keyID, secret, req)
			}

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)

			if tt.emailSent {
				time.Sleep(100 * time.Millisecond)
				mockEmailSender.AssertExpectations(t)
			}
		})
	}
}
