package main

import (
	"context"
	"fmt"
	"log"
	"os"

	dataedge "gitlab.com/dataedgemd/go-sdk"
)

func main() {
	// Get token from environment variable
	token := os.Getenv("DATAEDGE_TOKEN")
	if token == "" {
		log.Fatal("DATAEDGE_TOKEN environment variable is required")
	}

	// Initialize client
	client := dataedge.NewClient(token)
	smsService := dataedge.NewSmsService(client)
	viberService := dataedge.NewViberService(client)

	ctx := context.Background()

	// 1. Check Balance
	fmt.Println("--- Checking Balance ---")
	balance, err := smsService.GetBalance(ctx)
	if err != nil {
		log.Printf("Error getting balance: %v\n", err)
	} else {
		fmt.Printf("Current Balance: %.2f\n", balance)
	}

	// 2. Check Limit
	fmt.Println("\n--- Checking SMS Limit ---")
	limit, err := smsService.GetLimit(ctx)
	if err != nil {
		log.Printf("Error getting limit: %v\n", err)
	} else {
		fmt.Printf("Current Limit: %d\n", limit)
	}

	// 3. Send a Quick SMS (Example)
	// Uncomment to test real sending. WARNING: This costs money/credits.
	/*
	   fmt.Println("\n--- Sending Quick SMS ---")
	   phone := "373xxxxxxxxx" // Replace with real number
	   msg := "Hello from Go SDK!"
	   resp, err := smsService.SendQuickSms(ctx, msg, phone, true)
	   if err != nil {
	       log.Printf("Error sending SMS: %v\n", err)
	   } else {
	       fmt.Printf("SMS Sent Response: %v\n", resp)
	   }
	*/

	// 4. Send Viber Message (Example)
	// Uncomment to test real sending.
	/*
	   fmt.Println("\n--- Sending Viber Message ---")
	   viberPhone := "373xxxxxxxxx" // Replace
	   viberID := 123 // Replace with valid Viber Name ID
	   viberMsg := "Hello from Go SDK via Viber"
	   vResp, err := viberService.SendQuickViberMessage(ctx, viberPhone, viberID, viberMsg)
	   if err != nil {
	       log.Printf("Error sending Viber message: %v\n", err)
	   } else {
	       fmt.Printf("Viber Sent Response: %v\n", vResp)
	   }
	*/

	_ = viberService // Avoid unused variable error if block is commented
}
