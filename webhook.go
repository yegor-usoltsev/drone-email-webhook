package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/99designs/httpsignatures-go"
	"github.com/drone/drone-go/plugin/webhook"
)

type WebhookHandler struct {
	secret      string
	emailSender EmailSenderInterface
}

type EmailSenderInterface interface {
	Send(req *webhook.Request) error
}

func NewWebhookHandler(settings Settings, emailSender EmailSenderInterface) *WebhookHandler {
	return &WebhookHandler{secret: settings.Secret, emailSender: emailSender}
}

func (h *WebhookHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	signature, err := httpsignatures.FromRequest(r)
	if err != nil {
		log.Printf("[ERROR] webhook: invalid or missing signature for request from %s", r.RemoteAddr)
		http.Error(w, "Invalid or Missing Signature", http.StatusBadRequest)
		return
	}
	if !signature.IsValid(h.secret, r) {
		log.Printf("[ERROR] webhook: invalid signature for request from %s", r.RemoteAddr)
		http.Error(w, "Invalid Signature", http.StatusBadRequest)
		return
	}

	var req webhook.Request
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Printf("[ERROR] webhook: cannot unmarshal request body from %s: %v", r.RemoteAddr, err)
		http.Error(w, "Invalid Input", http.StatusBadRequest)
		return
	}

	if req.Event == webhook.EventBuild && req.Action == webhook.ActionUpdated && req.Build != nil && req.Build.Status == "failure" {
		log.Printf("[INFO] webhook: processing failure event for build #%d in repo %s", req.Build.ID, req.Repo.Slug)
		go func() {
			if err := h.emailSender.Send(&req); err != nil {
				log.Printf("[ERROR] webhook: failed to send email for build #%d in repo %s: %v", req.Build.ID, req.Repo.Slug, err)
			}
		}()
	}

	w.WriteHeader(http.StatusNoContent)
}
