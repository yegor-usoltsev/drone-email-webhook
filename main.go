package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
)

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))
	defer func() {
		if r := recover(); r != nil {
			slog.Error("panic", "error", r)
			os.Exit(1) //nolint:forbidigo
		}
	}()

	settings := NewSettingsFromEnv()
	server := NewServer(settings)
	emailSender := NewEmailSender(settings)
	defer emailSender.Wait()
	webhookHandler := NewWebhookHandler(settings, emailSender)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, _ *http.Request) {
		_, _ = fmt.Fprint(w, "OK")
	})
	mux.Handle("POST /", withRecovery(webhookHandler.ServeHTTP))

	server.ListenAndServe(mux)
}

func withRecovery(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				slog.Error("http handler panic recovered", "method", r.Method, "path", r.URL.Path, "error", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	}
}
