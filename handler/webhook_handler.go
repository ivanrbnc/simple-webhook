package handler

import (
	"encoding/json"
	"net/http"
	"simple-webhook/domain"
	"simple-webhook/service"
)

type WebhookHandler struct {
	webhookService *service.WebhookService
}

func NewWebhookHandler(webhookService *service.WebhookService) *WebhookHandler {
	return &WebhookHandler{
		webhookService: webhookService,
	}
}

func (h *WebhookHandler) Subscribe(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var sub domain.Subscription
	if err := json.NewDecoder(r.Body).Decode(&sub); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if err := h.webhookService.Subscribe(sub); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Subscription created successfully",
		"url":     sub.URL,
	})
}

func (h *WebhookHandler) TriggerEvent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Simple trigger with hardcoded data
	simpleData := json.RawMessage(`"New user registered: john@example.com"`)
	eventID, err := h.webhookService.TriggerEvent("user.created", simpleData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message":  "Event triggered",
		"event_id": eventID,
	})
}

func (h *WebhookHandler) TriggerCustomEvent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req domain.TriggerEventRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if req.Type == "" {
		http.Error(w, "Event type is required", http.StatusBadRequest)
		return
	}

	eventID, err := h.webhookService.TriggerEvent(req.Type, req.Data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message":  "Custom event triggered",
		"event_id": eventID,
	})
}
