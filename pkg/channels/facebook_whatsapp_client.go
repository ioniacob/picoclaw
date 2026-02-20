package channels

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// FacebookWhatsAppClient handles WhatsApp Business API through Facebook Graph API
type FacebookWhatsAppClient struct {
	phoneNumberID string
	accessToken   string
	apiVersion    string
	httpClient    *http.Client
	baseURL       string
}

// FacebookMessageRequest represents the message structure for Facebook WhatsApp API
type FacebookMessageRequest struct {
	MessagingProduct string                 `json:"messaging_product"`
	To               string                 `json:"to"`
	Type             string                 `json:"type"`
	Template         *FacebookTemplate      `json:"template,omitempty"`
	Text             *FacebookTextMessage   `json:"text,omitempty"`
	Image            *FacebookMediaMessage  `json:"image,omitempty"`
	Audio            *FacebookMediaMessage  `json:"audio,omitempty"`
	Video            *FacebookMediaMessage  `json:"video,omitempty"`
	Document         *FacebookMediaMessage  `json:"document,omitempty"`
}

// FacebookTemplate represents a template message
type FacebookTemplate struct {
	Name     string            `json:"name"`
	Language FacebookLanguage  `json:"language"`
	Components []TemplateComponent `json:"components,omitempty"`
}

// FacebookLanguage represents the language configuration
type FacebookLanguage struct {
	Code string `json:"code"`
}

// TemplateComponent represents template components
type TemplateComponent struct {
	Type       string                 `json:"type"`
	Parameters []TemplateParameter    `json:"parameters,omitempty"`
	Text       string                 `json:"text,omitempty"`
}

// TemplateParameter represents template parameters
type TemplateParameter struct {
	Type string `json:"type"`
	Text string `json:"text,omitempty"`
}

// FacebookTextMessage represents a text message
type FacebookTextMessage struct {
	Body string `json:"body"`
}

// FacebookMediaMessage represents a media message
type FacebookMediaMessage struct {
	ID   string `json:"id,omitempty"`
	Link string `json:"link,omitempty"`
	Caption string `json:"caption,omitempty"`
}

// FacebookMessageResponse represents the API response
type FacebookMessageResponse struct {
	MessagingProduct string   `json:"messaging_product"`
	Contacts          []Contact `json:"contacts"`
	Messages          []Message `json:"messages"`
}

// Contact represents contact information
type Contact struct {
	Input string `json:"input"`
	WaID  string `json:"wa_id"`
}

// Message represents message information
type Message struct {
	ID string `json:"id"`
}

// FacebookErrorResponse represents error responses
type FacebookErrorResponse struct {
	Error FacebookError `json:"error"`
}

// FacebookError represents an API error
type FacebookError struct {
	Message      string `json:"message"`
	Type         string `json:"type"`
	Code         int    `json:"code"`
	ErrorSubcode int    `json:"error_subcode,omitempty"`
	FBTraceID    string `json:"fbtrace_id"`
}

// NewFacebookWhatsAppClient creates a new Facebook WhatsApp client
func NewFacebookWhatsAppClient(phoneNumberID, accessToken, apiVersion string) *FacebookWhatsAppClient {
	if apiVersion == "" {
		apiVersion = "v22.0"
	}
	
	return &FacebookWhatsAppClient{
		phoneNumberID: phoneNumberID,
		accessToken:   accessToken,
		apiVersion:    apiVersion,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		baseURL: "https://graph.facebook.com",
	}
}

// SendTemplateMessage sends a template message
func (c *FacebookWhatsAppClient) SendTemplateMessage(ctx context.Context, to, templateName, languageCode string, components []TemplateComponent) error {
	message := FacebookMessageRequest{
		MessagingProduct: "whatsapp",
		To:               to,
		Type:             "template",
		Template: &FacebookTemplate{
			Name:     templateName,
			Language: FacebookLanguage{Code: languageCode},
			Components: components,
		},
	}
	
	return c.sendMessage(ctx, message)
}

// SendTextMessage sends a text message
func (c *FacebookWhatsAppClient) SendTextMessage(ctx context.Context, to, text string) error {
	message := FacebookMessageRequest{
		MessagingProduct: "whatsapp",
		To:               to,
		Type:             "text",
		Text: &FacebookTextMessage{
			Body: text,
		},
	}
	
	return c.sendMessage(ctx, message)
}

// sendMessage sends the actual message to Facebook API
func (c *FacebookWhatsAppClient) sendMessage(ctx context.Context, message FacebookMessageRequest) error {
	url := fmt.Sprintf("%s/%s/%s/messages", c.baseURL, c.apiVersion, c.phoneNumberID)
	
	jsonData, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}
	
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	
	req.Header.Set("Authorization", "Bearer "+c.accessToken)
	req.Header.Set("Content-Type", "application/json")
	
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}
	
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		var errorResp FacebookErrorResponse
		if err := json.Unmarshal(body, &errorResp); err != nil {
			return fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
		}
		return fmt.Errorf("Facebook API error: %s (type: %s, code: %d)", 
			errorResp.Error.Message, errorResp.Error.Type, errorResp.Error.Code)
	}
	
	var successResp FacebookMessageResponse
	if err := json.Unmarshal(body, &successResp); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}
	
	return nil
}

// ValidateCredentials validates the Facebook credentials
func (c *FacebookWhatsAppClient) ValidateCredentials(ctx context.Context) error {
	url := fmt.Sprintf("%s/%s/%s", c.baseURL, c.apiVersion, c.phoneNumberID)
	
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	
	req.Header.Set("Authorization", "Bearer "+c.accessToken)
	
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("credential validation failed (status %d): %s", resp.StatusCode, string(body))
	}
	
	return nil
}