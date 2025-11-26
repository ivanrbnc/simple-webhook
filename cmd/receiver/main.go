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
	receiverService := service.NewReceiverService(hmacSig)
	receiverHandler := handler.NewReceiverHandler(receiverService)

	// Setup routes
	http.HandleFunc("/webhook", receiverHandler.ReceiveWebhook)

	fmt.Println("Webhook Receiver running on :8081")
	fmt.Println("- POST /webhook - Receive webhook events")

	log.Fatal(http.ListenAndServe(":8081", nil))
}
