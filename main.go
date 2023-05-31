package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	settings := NewSettingsFromEnv()
	server := NewServer(settings)
	emailSender := NewEmailSender(settings)
	webhookHandler := NewWebhookHandler(settings, emailSender)

	router := chi.NewRouter()
	router.Use(middleware.Recoverer)
	router.Use(middleware.Timeout(HandlerTimeout))
	router.Use(middleware.Heartbeat("/health"))
	router.Post("/", webhookHandler.ServeHTTP)

	server.ListenAndServe(router)
}
