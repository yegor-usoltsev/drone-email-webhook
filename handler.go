package main

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/99designs/httpsignatures-go"
	"github.com/drone/drone-go/plugin/webhook"
)

type AsyncEmailSender interface {
	SendAsync(req *webhook.Request)
}

type Handler struct {
	http.Handler
}

func NewHandler(cfg Config, emailSender AsyncEmailSender) *Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", healthHandler)
	mux.Handle("POST /", webhookHandler(cfg.Secret, emailSender))
	return &Handler{Handler: withRecovery(mux)}
}

func healthHandler(w http.ResponseWriter, _ *http.Request) {
	_, _ = fmt.Fprint(w, "OK")
}

func webhookHandler(secret string, emailSender AsyncEmailSender) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		signature, err := httpsignatures.FromRequest(r)
		if err != nil {
			slog.Error("webhook handler received invalid or missing signature", "error", err)
			http.Error(w, "Invalid or Missing Signature", http.StatusBadRequest)
			return
		}
		if !signature.IsValid(secret, r) {
			slog.Error("webhook handler received invalid signature")
			http.Error(w, "Invalid Signature", http.StatusBadRequest)
			return
		}
		var req webhook.Request
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			slog.Error("webhook handler cannot unmarshal request body", "error", err)
			http.Error(w, "Invalid Input", http.StatusBadRequest)
			return
		}
		if req.Event == webhook.EventBuild && req.Action == webhook.ActionUpdated && req.Build != nil && req.Build.Status == "failure" {
			slog.Info("webhook handler processing build failure event", "build_id", req.Build.ID, "repo_slug", req.Repo.Slug)
			emailSender.SendAsync(&req)
		}
		w.WriteHeader(http.StatusNoContent)
	})
}

func withRecovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				slog.Error("http handler panic recovered", "method", r.Method, "path", r.URL.Path, "error", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
