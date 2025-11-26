package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"simple-webhook/domain"
	"time"
)

type WebhookService struct {
	subscriptions []domain.Subscription
	signatureGen  SignatureGenerator
}

func NewWebhookService(signatureGen SignatureGenerator) *WebhookService {
	return &WebhookService{
		subscriptions: make([]domain.Subscription, 0),
		signatureGen:  signatureGen,
	}
}

func (s *WebhookService) Subscribe(sub domain.Subscription) error {
	// Check if URL already exists
	for i, existing := range s.subscriptions {
		if existing.URL == sub.URL {
			s.subscriptions[i] = sub // UPDATE instead of adding
			return nil
		}
	}
	s.subscriptions = append(s.subscriptions, sub) // Only add if new
	return nil
}

func (s *WebhookService) TriggerEvent(eventType string, data json.RawMessage) (string, error) {
	event := domain.Event{
		ID:        fmt.Sprintf("evt_%d", time.Now().Unix()),
		Type:      eventType,
		Data:      data,
		Timestamp: time.Now(),
	}

	// Send webhook only to subscribers interested in this event type
	for _, sub := range s.subscriptions {
		if s.isSubscribedToEvent(sub, eventType) {
			go s.sendWebhook(sub.URL, event)
		}
	}

	return event.ID, nil
}

func (s *WebhookService) sendWebhook(url string, event domain.Event) {
	payload, err := json.Marshal(event)
	if err != nil {
		log.Printf("Error marshaling event: %v", err)
		return
	}

	signature := s.signatureGen.Generate(payload)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		log.Printf("Error creating request: %v", err)
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Webhook-Signature", signature)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error sending webhook to %s: %v", url, err)
		return
	}
	defer resp.Body.Close()

	log.Printf("Webhook sent to %s - Status: %d", url, resp.StatusCode)
}

func (s *WebhookService) isSubscribedToEvent(sub domain.Subscription, eventType string) bool {
	if len(sub.Events) == 0 {
		return true
	}

	for _, subscribedEvent := range sub.Events {
		if subscribedEvent == eventType || subscribedEvent == "*" {
			return true
		}
	}

	return false
}
