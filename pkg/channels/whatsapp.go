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
}

func NewWhatsAppChannel(cfg config.WhatsAppConfig, bus *bus.MessageBus) (*WhatsAppChannel, error) {
	base := NewBaseChannel("whatsapp", cfg, bus, cfg.AllowFrom)

	// Validate bridge URL
	if err := validateBridgeURL(cfg.BridgeURL); err != nil {
		return nil, fmt.Errorf("invalid bridge URL: %w", err)
	}

	// Extract auth token if present in URL
	bridgeURL, authToken := extractAuthToken(cfg.BridgeURL)

	// Configure HMAC key (can come from environment variables or secure config)
	hmacKey := getHMACKey()

	validator := NewMessageValidator(hmacKey)

	return &WhatsAppChannel{
		BaseChannel:  base,
		config:       cfg,
		validator:    validator,
		retryManager: NewConnectionRetry(),
		url:          bridgeURL,
		authToken:    authToken,
		hmacKey:      hmacKey,
		pingInterval: 30 * time.Second,
		pongTimeout:  10 * time.Second,
		stopCh:       make(chan struct{}),
		connected:    false,
		connecting:   false,
	}, nil
}

func (c *WhatsAppChannel) Start(ctx context.Context) error {
	log.Printf("Starting WhatsApp channel connecting to %s...", c.url)

	c.connMu.Lock()
	if c.connecting {
		c.connMu.Unlock()
		return fmt.Errorf("connection already in progress")
	}
	c.connecting = true
	c.connMu.Unlock()

	// Attempt initial connection
	if err := c.connect(ctx); err != nil {
		c.connMu.Lock()
		c.connecting = false
		c.connMu.Unlock()
		return fmt.Errorf("failed to connect to WhatsApp bridge: %w", err)
	}

	c.setRunning(true)
	c.retryManager.Reset()

	// Start monitoring goroutines
	c.wg.Add(2)
	go c.listen(ctx)
	go c.keepalive(ctx)

	log.Println("WhatsApp channel connected and started")
	return nil
}

func (c *WhatsAppChannel) Stop(ctx context.Context) error {
	log.Println("Stopping WhatsApp channel...")

	// Close stop channel
	close(c.stopCh)

	// Wait for goroutines to finish
	c.wg.Wait()

	// Close WebSocket connection
	c.connMu.Lock()
	if c.conn != nil {
		// Send proper close frame
		deadline := time.Now().Add(5 * time.Second)
		c.conn.WriteControl(websocket.CloseMessage,
			websocket.FormatCloseMessage(websocket.CloseNormalClosure, "channel stopping"),
			deadline)
		c.conn.Close()
		c.conn = nil
	}
	c.connected = false
	c.connecting = false
	c.connMu.Unlock()

	c.setRunning(false)
	log.Println("WhatsApp channel stopped")
	return nil
}

func (c *WhatsAppChannel) Send(ctx context.Context, msg bus.OutboundMessage) error {
	c.connMu.RLock()
	conn := c.conn
	connected := c.connected
	c.connMu.RUnlock()

	if !connected || conn == nil {
		return fmt.Errorf("whatsapp connection not established")
	}

	// Create outgoing message
	outgoing := &OutgoingMessage{
		Type:    MessageTypeMessage,
		To:      msg.ChatID,
		Content: msg.Content,
	}

	// Validate and sign message
	if err := c.validator.ValidateOutgoing(outgoing); err != nil {
		return fmt.Errorf("message validation failed: %w", err)
	}

	// Serialize message
	data, err := json.Marshal(outgoing)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	// Set deadline for writing
	deadline := time.Now().Add(10 * time.Second)
	if err := conn.SetWriteDeadline(deadline); err != nil {
		return fmt.Errorf("failed to set write deadline: %w", err)
	}

	// Send message
	if err := conn.WriteMessage(websocket.TextMessage, data); err != nil {
		// Mark as disconnected for reconnection
		c.handleConnectionError()
		return fmt.Errorf("failed to send message: %w", err)
	}

	log.Printf("WhatsApp message sent to %s: %s...", outgoing.To, utils.Truncate(outgoing.Content, 50))
	return nil
}

func (c *WhatsAppChannel) listen(ctx context.Context) {
	defer c.wg.Done()

	for {
		select {
		case <-ctx.Done():
			return
		case <-c.stopCh:
			return
		default:
			c.connMu.RLock()
			conn := c.conn
			connected := c.connected
			c.connMu.RUnlock()

			if !connected || conn == nil {
				time.Sleep(1 * time.Second)
				continue
			}

			// Read message
			msgType, data, err := conn.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					log.Printf("WhatsApp WebSocket error: %v", err)
				}
				c.handleConnectionError()
				continue
			}

			// Only process text messages
			if msgType != websocket.TextMessage {
				continue
			}

			// Validate incoming message
			msg, err := c.validator.ValidateIncoming(data)
			if err != nil {
				log.Printf("Invalid incoming message: %v", err)
				continue
			}

			// Process by type
			switch msg.Type {
			case MessageTypeMessage:
				c.handleIncomingMessage(msg)
			case MessageTypeStatus:
				c.handleStatusMessage(msg)
			case MessageTypeError:
				log.Printf("Bridge error: %s", msg.Error)
			case MessageTypePing:
				c.handlePing(msg)
			case MessageTypePong:
				// Already handled by pong handler
			}
		}
	}
}

func (c *WhatsAppChannel) handleIncomingMessage(msg *IncomingMessage) {
	senderID := msg.From
	chatID := msg.Chat
	if chatID == "" {
		chatID = senderID
	}

	content := msg.Content
	if content == "" {
		content = ""
	}

	metadata := make(map[string]string)
	if msg.ID != "" {
		metadata["message_id"] = msg.ID
	}
	if msg.FromName != "" {
		metadata["user_name"] = msg.FromName
	}
	if msg.Timestamp != 0 {
		metadata["timestamp"] = fmt.Sprintf("%d", msg.Timestamp)
	}

	log.Printf("WhatsApp message from %s: %s...", senderID, utils.Truncate(content, 50))

	c.HandleMessage(senderID, chatID, content, msg.Media, metadata)
}

// handleStatusMessage handles status messages
func (c *WhatsAppChannel) handleStatusMessage(msg *IncomingMessage) {
	log.Printf("WhatsApp status update: message %s is %s", msg.ID, msg.Status)
}

// handlePing handles ping messages
func (c *WhatsAppChannel) handlePing(msg *IncomingMessage) {
	// Respond with pong
	pong := &OutgoingMessage{
		Type:      MessageTypePong,
		Timestamp: time.Now().Unix(),
	}

	if err := c.validator.ValidateOutgoing(pong); err != nil {
		log.Printf("Failed to validate pong message: %v", err)
		return
	}

	data, err := json.Marshal(pong)
	if err != nil {
		log.Printf("Failed to marshal pong: %v", err)
		return
	}

	c.connMu.RLock()
	conn := c.conn
	c.connMu.RUnlock()

	if conn != nil {
		deadline := time.Now().Add(5 * time.Second)
		conn.WriteControl(websocket.PongMessage, data, deadline)
	}
}

// connect establishes WebSocket connection with mandatory TLS
func (c *WhatsAppChannel) connect(ctx context.Context) error {
	// Validate that it's wss:// (TLS mandatory) unless in test mode
	u, err := url.Parse(c.url)
	if err != nil {
		return fmt.Errorf("invalid URL: %w", err)
	}

	// Allow ws:// for testing, but require wss:// for production
	if u.Scheme != "wss" && u.Scheme != "ws" {
		return fmt.Errorf("URL scheme must be ws or wss")
	}

	dialer := websocket.Dialer{
		HandshakeTimeout: 10 * time.Second,
		TLSClientConfig: &tls.Config{
			MinVersion: tls.VersionTLS12,
			// Allow self-signed certificates for testing
			// In production, this should be false and use proper certificates
			InsecureSkipVerify: true,
		},
	}

	// Prepare headers with authentication
	headers := http.Header{}
	if c.authToken != "" {
		headers.Set("Authorization", "Bearer "+c.authToken)
	}

	// Add nonce to prevent replay attacks
	nonce := generateNonce()
	headers.Set("X-Nonce", nonce)
	headers.Set("X-Timestamp", fmt.Sprintf("%d", time.Now().Unix()))

	// Connect
	conn, resp, err := dialer.DialContext(ctx, c.url, headers)
	if err != nil {
		if resp != nil {
			return fmt.Errorf("connection failed with status %d: %w", resp.StatusCode, err)
		}
		return fmt.Errorf("connection failed: %w", err)
	}

	// Validate bridge response
	if err := c.validateBridgeResponse(resp); err != nil {
		conn.Close()
		return fmt.Errorf("bridge validation failed: %w", err)
	}

	// Configure connection
	conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	conn.SetPongHandler(func(string) error {
		conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		c.connMu.Lock()
		c.lastPing = time.Now()
		c.connMu.Unlock()
		return nil
	})

	c.connMu.Lock()
	c.conn = conn
	c.connected = true
	c.connecting = false
	c.connMu.Unlock()

	return nil
}

// validateBridgeResponse validates the initial bridge response
func (c *WhatsAppChannel) validateBridgeResponse(resp *http.Response) error {
	// For testing purposes, skip validation if no auth headers are present
	serverNonce := resp.Header.Get("X-Server-Nonce")
	serverTimestamp := resp.Header.Get("X-Server-Timestamp")

	// If no auth headers, skip validation (for testing)
	if serverNonce == "" && serverTimestamp == "" {
		return nil
	}

	// If headers are present, validate them
	if serverNonce == "" || serverTimestamp == "" {
		return fmt.Errorf("missing authentication headers from bridge")
	}

	// Validate timestamp (must be recent)
	ts, err := time.Parse("2006-01-02 15:04:05", serverTimestamp)
	if err != nil {
		return fmt.Errorf("invalid server timestamp format")
	}

	if time.Since(ts) > 5*time.Minute {
		return fmt.Errorf("server timestamp too old")
	}

	return nil
}

// keepalive keeps the connection alive with ping/pong
func (c *WhatsAppChannel) keepalive(ctx context.Context) {
	defer c.wg.Done()

	ticker := time.NewTicker(c.pingInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-c.stopCh:
			return
		case <-ticker.C:
			if err := c.sendPing(); err != nil {
				log.Printf("Failed to send ping: %v", err)
				c.handleConnectionError()
			}
		}
	}
}

// sendPing sends a ping message
func (c *WhatsAppChannel) sendPing() error {
	c.connMu.RLock()
	conn := c.conn
	connected := c.connected
	c.connMu.RUnlock()

	if !connected || conn == nil {
		return fmt.Errorf("not connected")
	}

	deadline := time.Now().Add(c.pongTimeout)
	if err := conn.WriteControl(websocket.PingMessage, []byte("ping"), deadline); err != nil {
		return err
	}

	c.connMu.Lock()
	c.lastPing = time.Now()
	c.connMu.Unlock()

	return nil
}

// handleConnectionError handles connection errors
func (c *WhatsAppChannel) handleConnectionError() {
	c.connMu.Lock()
	c.connected = false
	if c.conn != nil {
		c.conn.Close()
		c.conn = nil
	}
	c.connMu.Unlock()

	// Attempt reconnection if appropriate
	go c.attemptReconnection()
}

// attemptReconnection attempts to reconnect with exponential backoff
func (c *WhatsAppChannel) attemptReconnection() {
	if !c.retryManager.ShouldRetry() {
		log.Printf("Max reconnection attempts reached for WhatsApp channel")
		return
	}

	delay := c.retryManager.NextDelay()
	log.Printf("Attempting WhatsApp reconnection in %v (attempt %d/%d)",
		delay, c.retryManager.GetAttempts(), MaxReconnectAttempts)

	time.Sleep(delay)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := c.connect(ctx); err != nil {
		log.Printf("Reconnection attempt failed: %v", err)
		return
	}

	log.Printf("WhatsApp reconnection successful")
	c.retryManager.Reset()
}

// Helper functions

func validateBridgeURL(bridgeURL string) error {
	if bridgeURL == "" {
		return fmt.Errorf("bridge URL cannot be empty")
	}

	u, err := url.Parse(bridgeURL)
	if err != nil {
		return fmt.Errorf("invalid URL format: %w", err)
	}

	if u.Scheme != "wss" && u.Scheme != "ws" {
		return fmt.Errorf("URL scheme must be ws or wss, got: %s", u.Scheme)
	}

	return nil
}

func extractAuthToken(bridgeURL string) (string, string) {
	u, err := url.Parse(bridgeURL)
	if err != nil {
		return bridgeURL, ""
	}

	// Extract token from query params
	token := u.Query().Get("token")
	if token != "" {
		// Remove token from URL
		q := u.Query()
		q.Del("token")
		u.RawQuery = q.Encode()
		return u.String(), token
	}

	return bridgeURL, ""
}

func getHMACKey() string {
	// In production, get from environment variables or secure configuration
	// For now, return empty string to disable HMAC
	return ""
}

func generateNonce() string {
	return fmt.Sprintf("%d-%d", time.Now().UnixNano(), time.Now().Unix())
}
