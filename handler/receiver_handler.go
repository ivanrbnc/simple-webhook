package handler

import (
	"encoding/json"
	"io"
	"net/http"
	"simple-webhook/service"
)

type ReceiverHandler struct {
	receiverService *service.ReceiverService
}

func NewReceiverHandler(receiverService *service.ReceiverService) *ReceiverHandler {
	return &ReceiverHandler{
		receiverService: receiverService,
	}
}

func (h *ReceiverHandler) ReceiveWebhook(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading body", http.StatusBadRequest)
		return
	}

	signature := r.Header.Get("X-Webhook-Signature")
	
	event, err := h.receiverService.ProcessWebhook(body, signature)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status":   "received",
		"event_id": event.ID,
	})
}