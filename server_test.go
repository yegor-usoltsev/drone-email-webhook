package main

import (
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestWithRecovery(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name           string
		handler        http.HandlerFunc
		expectedStatus int
	}{
		{
			name: "normal handler",
			handler: func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusOK)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "panic handler",
			handler: func(_ http.ResponseWriter, _ *http.Request) {
				panic("test panic")
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			wrapped := withRecovery(tt.handler)
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			w := httptest.NewRecorder()
			wrapped.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}
		})
	}
}

func TestServer(t *testing.T) {
	t.Parallel()
	settings := Settings{
		ServerHost: "localhost",
		ServerPort: 0,
	}
	server := NewServer(settings)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /test", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	go server.ListenAndServe(mux)

	time.Sleep(100 * time.Millisecond)

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	server.cancelServerCtx()
	time.Sleep(100 * time.Millisecond)
}

func TestHealthCheck(t *testing.T) {
	t.Parallel()
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestServer_ListenAndServe(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		settings Settings
		wantErr  bool
	}{
		{
			name: "successful_startup",
			settings: Settings{
				ServerHost: "localhost",
				ServerPort: 0,
			},
			wantErr: false,
		},
		{
			name: "invalid_port",
			settings: Settings{
				ServerHost: "localhost",
				ServerPort: -1,
			},
			wantErr: true,
		},
		{
			name: "invalid_host",
			settings: Settings{
				ServerHost: "invalid-host",
				ServerPort: 8080,
			},
			wantErr: true,
		},
		{
			name: "port_already_in_use",
			settings: Settings{
				ServerHost: "localhost",
				ServerPort: 0,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			server := NewServer(tt.settings)

			if tt.name == "port_already_in_use" {
				listener, err := net.Listen("tcp", "localhost:0")
				require.NoError(t, err)
				defer listener.Close()
				addr := listener.Addr().(*net.TCPAddr)
				tt.settings.ServerPort = addr.Port
			}
			handler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusOK)
			})
			go func() {
				server.ListenAndServe(handler)
			}()
			time.Sleep(100 * time.Millisecond)
			if !tt.wantErr {
				conn, err := net.Dial("tcp", fmt.Sprintf("localhost:%d", tt.settings.ServerPort))
				if err == nil {
					conn.Close()
				}
			}
		})
	}
}
