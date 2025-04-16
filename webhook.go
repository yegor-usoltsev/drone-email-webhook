package main

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/99designs/httpsignatures-go"
	"github.com/drone/drone-go/plugin/webhook"
)

type WebhookHandler struct {
	secret      string
	emailSender AsyncEmailSender
}

type AsyncEmailSender interface {
	SendAsync(req *webhook.Request)
}

func NewWebhookHandler(settings Settings, emailSender AsyncEmailSender) *WebhookHandler {
	return &WebhookHandler{secret: settings.Secret, emailSender: emailSender}
}

func (h *WebhookHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	signature, err := httpsignatures.FromRequest(r)
	if err != nil {
		slog.Error("webhook handler received invalid or missing signature", "error", err)
		http.Error(w, "Invalid or Missing Signature", http.StatusBadRequest)
		return
	}
	if !signature.IsValid(h.secret, r) {
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
		h.emailSender.SendAsync(&req)
	}

	w.WriteHeader(http.StatusNoContent)
}
