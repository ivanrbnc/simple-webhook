package main

import (
	"fmt"
	"log"
	"net/http"
	"simple-webhook/handler"
	"simple-webhook/service"
)

func main() {
	// Initialize dependencies
	secret := "your-secret-key-here"
	hmacSig := service.NewHMACSignature(secret)
	webhookService := service.NewWebhookService(hmacSig)
	webhookHandler := handler.NewWebhookHandler(webhookService)

	// Setup routes
	http.HandleFunc("/subscribe", webhookHandler.Subscribe)
	http.HandleFunc("/trigger", webhookHandler.TriggerEvent)
	http.HandleFunc("/trigger-custom", webhookHandler.TriggerCustomEvent)

	fmt.Println("Webhook Server running on :8080")
	fmt.Println("- POST /subscribe - Subscribe to webhooks")
	fmt.Println("- POST /trigger - Trigger a simple test event")
	fmt.Println("- POST /trigger-custom - Trigger event with custom JSON data")

	log.Fatal(http.ListenAndServe(":8080", nil))
}
