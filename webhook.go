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
	emailSender *EmailSender
}

func NewWebhookHandler(settings Settings, emailSender *EmailSender) *WebhookHandler {
	return &WebhookHandler{secret: settings.Secret, emailSender: emailSender}
}

func (h *WebhookHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	signature, err := httpsignatures.FromRequest(r)
	if err != nil {
		log.Println("webhook: invalid or missing signature")
		http.Error(w, "Invalid or Missing Signature", http.StatusBadRequest)
		return
	}
	if !signature.IsValid(h.secret, r) {
		log.Println("webhook: invalid signature")
		http.Error(w, "Invalid Signature", http.StatusBadRequest)
		return
	}

	req := &webhook.Request{}
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Println("webhook: cannot unmarshal request body")
		http.Error(w, "Invalid Input", http.StatusBadRequest)
		return
	}

	if req.Event == webhook.EventBuild && req.Action == webhook.ActionUpdated && req.Build.Status == "failure" {
		log.Printf("webhook: processing event for build #%d in repo %s\n", req.Build.ID, req.Repo.Slug)
		go func() {
			if err := h.emailSender.Send(req); err != nil {
				log.Printf("webhook: failed to send email for build #%d: %v\n", req.Build.ID, err)
			}
		}()
	}

	w.WriteHeader(http.StatusNoContent)
}
