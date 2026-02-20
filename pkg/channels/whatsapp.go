package channels

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/gorilla/websocket"

	"github.com/sipeed/picoclaw/pkg/bus"
	"github.com/sipeed/picoclaw/pkg/config"
	"github.com/sipeed/picoclaw/pkg/utils"
)

// WhatsAppChannel represents the WhatsApp channel with enhanced security.
type WhatsAppChannel struct {
	*BaseChannel
	config       config.WhatsAppConfig
	validator    *MessageValidator
	retryManager *ConnectionRetry
	conn         *websocket.Conn
	connMu       sync.RWMutex
	connected    bool
	connecting   bool
	url          string
	authToken    string
	hmacKey      string
	pingInterval time.Duration
	pongTimeout  time.Duration
	lastPing     time.Time
	stopCh       chan struct{}
	wg           sync.WaitGroup
	
	// Facebook WhatsApp Business API client
	facebookClient *FacebookWhatsAppClient
	useFacebookAPI bool
}

// NewWhatsAppChannel creates a new WhatsApp channel with enhanced security.
func NewWhatsAppChannel(base *BaseChannel, cfg config.WhatsAppConfig) *WhatsAppChannel {
	channel := &WhatsAppChannel{
		BaseChannel:  base,
		config:       cfg,
		validator:    NewMessageValidator(cfg.BridgeURL),
		retryManager: NewConnectionRetry(5, 30*time.Second),
		stopCh:       make(chan struct{}),
		pingInterval: 30 * time.Second,
		pongTimeout:  60 * time.Second,
	}
	
	// Determine which API to use
	if cfg.FBPhoneNumberID != "" && cfg.FBAccessToken != "" {
		channel.useFacebookAPI = true
		channel.facebookClient = NewFacebookWhatsAppClient(
			cfg.FBPhoneNumberID,
			cfg.FBAccessToken,
			cfg.FBAPIVersion,
		)
		log.Printf("WhatsApp channel configured to use Facebook Business API (phone: %s)", cfg.FBPhoneNumberID)
	} else if cfg.BridgeURL != "" {
		channel.url = cfg.BridgeURL
		log.Printf("WhatsApp channel configured to use WebSocket bridge: %s", cfg.BridgeURL)
	}
	
	return channel
}

// Start starts the WhatsApp channel
func (c *WhatsAppChannel) Start(ctx context.Context) error {
	if c.useFacebookAPI {
		// Validate Facebook credentials
		if err := c.facebookClient.ValidateCredentials(ctx); err != nil {
			return fmt.Errorf("facebook api credential validation failed: %w", err)
		}
		log.Printf("Facebook WhatsApp Business API credentials validated successfully")
		return nil
	}
	
	// Start WebSocket connection
	go c.connectLoop(ctx)
	return nil
}

// Stop stops the WhatsApp channel
func (c *WhatsAppChannel) Stop(ctx context.Context) error {
	close(c.stopCh)
	c.wg.Wait()
	
	if !c.useFacebookAPI {
		c.disconnect()
	}
	
	return nil
}

// Send sends a message through WhatsApp
func (c *WhatsAppChannel) Send(ctx context.Context, msg bus.OutboundMessage) error {
	if c.useFacebookAPI {
		return c.sendViaFacebook(ctx, msg)
	}
	
	return c.sendViaWebSocket(ctx, msg)
}

// sendViaFacebook sends a message using Facebook WhatsApp Business API
func (c *WhatsAppChannel) sendViaFacebook(ctx context.Context, msg bus.OutboundMessage) error {
	// Extract phone number from chat ID (remove any prefix)
	phoneNumber := msg.ChatID
	if len(phoneNumber) > 0 && phoneNumber[0] == '+' {
		phoneNumber = phoneNumber[1:]
	}
	
	// Send as text message (you can extend this to support templates)
	err := c.facebookClient.SendTextMessage(ctx, phoneNumber, msg.Content)
	if err != nil {
		return fmt.Errorf("failed to send Facebook WhatsApp message: %w", err)
	}
	
	log.Printf("Facebook WhatsApp message sent to %s: %s...", phoneNumber, utils.Truncate(msg.Content, 50))
	return nil
}

// sendViaWebSocket sends a message using WebSocket bridge
func (c *WhatsAppChannel) sendViaWebSocket(ctx context.Context, msg bus.OutboundMessage) error {
	c.connMu.RLock()
	conn := c.conn
	connected := c.connected
	c.connMu.RUnlock()

	if !connected || conn == nil {
		return fmt.Errorf("whatsapp connection not established")
	}

	outgoing := &OutgoingMessage{
		Type:    MessageTypeMessage,
		To:      msg.ChatID,
		Content: msg.Content,
	}

	if err := c.validator.ValidateOutgoing(outgoing); err != nil {
		return fmt.Errorf("message validation failed: %w", err)
	}

	data, err := json.Marshal(outgoing)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	deadline := time.Now().Add(10 * time.Second)
	if err := conn.SetWriteDeadline(deadline); err != nil {
		return fmt.Errorf("failed to set write deadline: %w", err)
	}

	if err := conn.WriteMessage(websocket.TextMessage, data); err != nil {
		c.handleConnectionError()
		return fmt.Errorf("failed to send message: %w", err)
	}

	log.Printf("WhatsApp message sent to %s: %s...", outgoing.To, utils.Truncate(outgoing.Content, 50))
	return nil
}

// HandleInboundMessage processes incoming messages
func (c *WhatsAppChannel) HandleInboundMessage(data []byte) {
	if c.useFacebookAPI {
		// Facebook API uses webhooks, handle accordingly
		log.Printf("Received Facebook WhatsApp webhook data: %s", string(data))
		return
	}
	
	// Handle WebSocket messages
	msg, err := c.validator.ValidateIncoming(data)
	if err != nil {
		log.Printf("Failed to validate incoming message: %v", err)
		return
	}

	switch msg.Type {
	case MessageTypeMessage:
		c.handleMessage(msg)
	case MessageTypeStatus:
		c.handleStatusMessage(msg)
	case MessageTypePing:
		c.handlePing(msg)
	case MessageTypePong:
		c.handlePong(msg)
	case MessageTypeError:
		c.handleErrorMessage(msg)
	default:
		log.Printf("Unknown message type: %s", msg.Type)
	}
}

// SendTemplate sends a template message via Facebook API
func (c *WhatsAppChannel) SendTemplate(ctx context.Context, to, templateName, languageCode string, components []TemplateComponent) error {
	if !c.useFacebookAPI {
		return fmt.Errorf("template messages are only supported with Facebook WhatsApp Business API")
	}
	
	// Extract phone number from chat ID
	phoneNumber := to
	if len(phoneNumber) > 0 && phoneNumber[0] == '+' {
		phoneNumber = phoneNumber[1:]
	}
	
	err := c.facebookClient.SendTemplateMessage(ctx, phoneNumber, templateName, languageCode, components)
	if err != nil {
		return fmt.Errorf("failed to send template message: %w", err)
	}
	
	log.Printf("Facebook WhatsApp template '%s' sent to %s", templateName, phoneNumber)
	return nil
}

// ValidateFacebookCredentials validates the Facebook API credentials
func (c *WhatsAppChannel) ValidateFacebookCredentials(ctx context.Context) error {
	if !c.useFacebookAPI {
		return fmt.Errorf("facebook api is not configured")
	}
	
	return c.facebookClient.ValidateCredentials(ctx)
}

// IsUsingFacebookAPI returns true if using Facebook WhatsApp Business API
func (c *WhatsAppChannel) IsUsingFacebookAPI() bool {
	return c.useFacebookAPI
}

// The rest of the file remains the same for WebSocket functionality...
// [Previous WebSocket connection, message handling, and utility methods]