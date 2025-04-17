package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/99designs/httpsignatures-go"
	"github.com/drone/drone-go/drone"
	"github.com/drone/drone-go/plugin/webhook"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

var webhookRequest = &webhook.Request{
	Event:  webhook.EventBuild,
	Action: webhook.ActionUpdated,
	Build:  &drone.Build{Status: "failure", ID: 42},
	Repo:   &drone.Repo{Slug: "test/repo"},
}

type MockEmailSender struct {
	mock.Mock
}

func NewMockEmailSender() *MockEmailSender {
	return &MockEmailSender{}
}

func (m *MockEmailSender) SendAsync(req *webhook.Request) {
	m.Called(req)
}

func assertHTTPStatusCode(t *testing.T, handler http.HandlerFunc, method, url string, body any, statuscode int) {
	t.Helper()
	jsonBody, err := json.Marshal(body)
	require.NoError(t, err)
	req := httptest.NewRequest(method, url, bytes.NewReader(jsonBody))
	err = httpsignatures.DefaultSha256Signer.SignRequest("test-key-id", "test-secret", req)
	require.NoError(t, err)
	w := httptest.NewRecorder()
	handler(w, req)
	assert.Equal(t, statuscode, w.Code)
}

func TestNewHandler(t *testing.T) {
	t.Parallel()
	emailSender := NewMockEmailSender()
	emailSender.On("SendAsync", mock.Anything).Return()
	defer emailSender.AssertExpectations(t)

	handler := NewHandler(Config{Secret: "test-secret"}, emailSender).ServeHTTP

	assert.HTTPSuccess(t, handler, http.MethodGet, "/health", nil)
	assertHTTPStatusCode(t, handler, http.MethodPost, "/", webhookRequest, http.StatusNoContent)
}

func TestHealthHandler(t *testing.T) {
	t.Parallel()
	url := "/health"
	assert.HTTPStatusCode(t, healthHandler, http.MethodGet, url, nil, http.StatusOK)
	assert.HTTPBodyContains(t, healthHandler, http.MethodGet, url, nil, "OK")
}

func TestWebhookHandler(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		t.Parallel()
		emailSender := NewMockEmailSender()
		emailSender.On("SendAsync", mock.Anything).Return()
		defer emailSender.AssertExpectations(t)

		handler := webhookHandler("test-secret", emailSender).ServeHTTP

		assertHTTPStatusCode(t, handler, http.MethodPost, "/", webhookRequest, http.StatusNoContent)
	})

	t.Run("missing signature", func(t *testing.T) {
		t.Parallel()
		emailSender := NewMockEmailSender()

		handler := webhookHandler("test-secret", emailSender).ServeHTTP

		jsonBody, err := json.Marshal(webhookRequest)
		require.NoError(t, err)
		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(jsonBody))
		w := httptest.NewRecorder()
		handler(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("invalid signature", func(t *testing.T) {
		t.Parallel()
		emailSender := NewMockEmailSender()

		handler := webhookHandler("test-secret", emailSender).ServeHTTP

		jsonBody, err := json.Marshal(webhookRequest)
		require.NoError(t, err)
		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(jsonBody))
		err = httpsignatures.DefaultSha256Signer.SignRequest("test-key-id", "invalid-secret", req)
		require.NoError(t, err)
		w := httptest.NewRecorder()
		handler(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("invalid json", func(t *testing.T) {
		t.Parallel()
		emailSender := NewMockEmailSender()

		handler := webhookHandler("test-secret", emailSender).ServeHTTP

		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader([]byte("invalid json")))
		err := httpsignatures.DefaultSha256Signer.SignRequest("test-key-id", "test-secret", req)
		require.NoError(t, err)
		w := httptest.NewRecorder()
		handler(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestWithRecovery(t *testing.T) {
	t.Run("normal handler", func(t *testing.T) {
		t.Parallel()
		handler := withRecovery(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusOK)
		})).ServeHTTP
		assert.HTTPStatusCode(t, handler, http.MethodGet, "/", nil, http.StatusOK)
	})

	t.Run("panic handler", func(t *testing.T) {
		t.Parallel()
		handler := withRecovery(http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {
			panic("test panic")
		})).ServeHTTP
		assert.HTTPStatusCode(t, handler, http.MethodGet, "/", nil, http.StatusInternalServerError)
	})
}
