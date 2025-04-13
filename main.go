package main

import (
	"log"
	"net/http"
)

func main() {
	settings := NewSettingsFromEnv()
	server := NewServer(settings)
	emailSender := NewEmailSender(settings)
	webhookHandler := NewWebhookHandler(settings, emailSender)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	mux.Handle("POST /", withRecovery(webhookHandler.ServeHTTP))

	server.ListenAndServe(mux)
}

func withRecovery(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("[ERROR] server: panic recovered: %v, path: %s, method: %s", err, r.URL.Path, r.Method)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	}
}
