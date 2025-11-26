package service

import (
	"encoding/json"
	"fmt"
	"log"
	"simple-webhook/domain"
	"time"
)

type ReceiverService struct {
	signatureVerifier SignatureVerifier
}

func NewReceiverService(verifier SignatureVerifier) *ReceiverService {
	return &ReceiverService{
		signatureVerifier: verifier,
	}
}

func (s *ReceiverService) ProcessWebhook(payload []byte, signature string) (*domain.Event, error) {
	// Verify signature
	if !s.signatureVerifier.Verify(payload, signature) {
		log.Println("⚠️  Webhook rejected: Invalid signature")
		return nil, fmt.Errorf("invalid signature")
	}

	// Parse event
	var event domain.Event
	if err := json.Unmarshal(payload, &event); err != nil {
		return nil, fmt.Errorf("invalid JSON: %w", err)
	}

	// Log received event
	log.Printf("✅ Webhook received!")
	log.Printf("   Event ID: %s", event.ID)
	log.Printf("   Type: %s", event.Type)
	log.Printf("   Data: %s", event.Data)
	log.Printf("   Timestamp: %s", event.Timestamp.Format(time.RFC3339))

	return &event, nil
}
