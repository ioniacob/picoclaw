package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/sipeed/picoclaw/pkg/channels"
)

func main() {
	// Example usage of Facebook WhatsApp Business API integration
	
	// IMPORTANT: Replace these with your actual Facebook credentials
	phoneNumberID := "YOUR_PHONE_NUMBER_ID" // Your phone number ID from Facebook
	accessToken := "YOUR_ACCESS_TOKEN"       // Your access token from Facebook
	apiVersion := "v22.0"                    // API version (optional, defaults to v22.0)
	
	// Create Facebook WhatsApp client
	client := channels.NewFacebookWhatsAppClient(phoneNumberID, accessToken, apiVersion)
	
	// Validate credentials
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	
	if err := client.ValidateCredentials(ctx); err != nil {
		log.Fatalf("Failed to validate credentials: %v", err)
	}
	fmt.Println("✅ Facebook WhatsApp Business API credentials validated successfully")
	
	// Example 1: Send the exact template message from your curl command
	recipient := "RECIPIENT_PHONE_NUMBER" // Phone number from your curl command
	templateName := "hello_world"
	languageCode := "en_US"
	
	fmt.Printf("📤 Sending template message '%s' to %s...\n", templateName, recipient)
	err := client.SendTemplateMessage(ctx, recipient, templateName, languageCode, nil)
	if err != nil {
		log.Fatalf("Failed to send template message: %v", err)
	}
	fmt.Println("✅ Template message sent successfully")
	
	// Example 2: Send a text message
	fmt.Printf("📤 Sending text message to %s...\n", recipient)
	err = client.SendTextMessage(ctx, recipient, "Hello from PicoClaw! This is a test message.")
	if err != nil {
		log.Fatalf("Failed to send text message: %v", err)
	}
	fmt.Println("✅ Text message sent successfully")
	
	// Example 3: Send a template with parameters
	fmt.Printf("📤 Sending template message with parameters to %s...\n", recipient)
	
	// Create template components with parameters
	components := []channels.TemplateComponent{
		{
			Type: "body",
			Parameters: []channels.TemplateParameter{
				{Type: "text", Text: "John Doe"},
				{Type: "text", Text: "PicoClaw"},
			},
		},
	}
	
	err = client.SendTemplateMessage(ctx, recipient, "welcome_message", "en_US", components)
	if err != nil {
		log.Printf("Failed to send parameterized template: %v", err)
	} else {
		fmt.Println("✅ Parameterized template message sent successfully")
	}
	
	fmt.Println("\n🎉 All Facebook WhatsApp Business API operations completed successfully!")
}

// Example configuration for PicoClaw
func exampleConfig() {
	fmt.Println("\n📋 Example PicoClaw configuration for Facebook WhatsApp Business API:")
	fmt.Println(`{
  "channels": {
    "whatsapp": {
      "enabled": true,
      "fb_phone_number_id": "YOUR_PHONE_NUMBER_ID",
      "fb_access_token": "YOUR_ACCESS_TOKEN",
      "fb_api_version": "v22.0",
      "allow_from": ["ALLOWED_PHONE_NUMBERS"]
    }
  }
}`)
	
	fmt.Println("\n🌍 Environment variables alternative:")
	fmt.Println("PICOCLAW_CHANNELS_WHATSAPP_ENABLED=true")
	fmt.Println("PICOCLAW_CHANNELS_WHATSAPP_FB_PHONE_NUMBER_ID=YOUR_PHONE_NUMBER_ID")
	fmt.Println("PICOCLAW_CHANNELS_WHATSAPP_FB_ACCESS_TOKEN=YOUR_ACCESS_TOKEN")
	fmt.Println("PICOCLAW_CHANNELS_WHATSAPP_FB_API_VERSION=v22.0")
	fmt.Println("PICOCLAW_CHANNELS_WHATSAPP_ALLOW_FROM=[\"ALLOWED_PHONE_NUMBERS\"]")
}

// Example curl command equivalent
func exampleCurlEquivalent() {
	fmt.Println("\n🔧 Your curl command equivalent in Go:")
	fmt.Println(`// This is what the FacebookWhatsAppClient.SendTemplateMessage() does internally:
ctx := context.Background()
client := channels.NewFacebookWhatsAppClient("YOUR_PHONE_NUMBER_ID", "YOUR_ACCESS_TOKEN", "v22.0")

message := channels.FacebookMessageRequest{
    MessagingProduct: "whatsapp",
    To:               "RECIPIENT_PHONE_NUMBER",
    Type:             "template",
    Template: &channels.FacebookTemplate{
        Name:     "hello_world",
        Language: channels.FacebookLanguage{Code: "en_US"},
    },
}

// The client handles the HTTP POST to https://graph.facebook.com/v22.0/YOUR_PHONE_NUMBER_ID/messages
// with Authorization: Bearer YOUR_ACCESS_TOKEN
// and Content-Type: application/json`)
}