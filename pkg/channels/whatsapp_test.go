package channels

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/sipeed/picoclaw/pkg/bus"
	"github.com/sipeed/picoclaw/pkg/config"
)

// TestWhatsAppChannelConnection tests WhatsApp channel connection
func TestWhatsAppChannelConnection(t *testing.T) {
	// Create un servidor WebSocket de prueba que simula el bridge de WhatsApp
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			t.Logf("Error upgrading connection: %v", err)
			return
		}
		defer conn.Close()

		// Handle mensajes del cliente
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				break
			}

			// Echo back to confirm reception
			var msg map[string]interface{}
			if err := json.Unmarshal(message, &msg); err != nil {
				continue
			}

			// Send confirmation
			response := map[string]interface{}{
				"type":   "status",
				"status": "delivered",
			}
			if err := conn.WriteJSON(response); err != nil {
				break
			}
		}
	}))
	defer server.Close()

	// Convertir https:// a wss://
	wsURL := strings.Replace(server.URL, "https://", "wss://", 1)

	// WhatsApp configuration
	cfg := config.WhatsAppConfig{
		Enabled:   true,
		BridgeURL: wsURL,
		AllowFrom: []string{},
	}

	msgBus := bus.NewMessageBus()
	channel, err := NewWhatsAppChannel(cfg, msgBus)
	if err != nil {
		t.Fatalf("Error creating WhatsApp channel: %v", err)
	}

	// Connection test
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = channel.Start(ctx)
	if err != nil {
		t.Fatalf("Error starting WhatsApp channel: %v", err)
	}

	// Verify channel is connected
	if !channel.IsRunning() {
		t.Error("WhatsApp channel should be running after successful start")
	}

	// Test message sending
	testMsg := bus.OutboundMessage{
		Channel: "whatsapp",
		ChatID:  "test-user",
		Content: "Hello from test!",
	}

	err = channel.Send(ctx, testMsg)
	if err != nil {
		t.Errorf("Error sending message: %v", err)
	}

	// Detener el canal
	err = channel.Stop(ctx)
	if err != nil {
		t.Errorf("Error stopping WhatsApp channel: %v", err)
	}

	if channel.IsRunning() {
		t.Error("WhatsApp channel should not be running after stop")
	}
}

// TestWhatsAppMessageFormat prueba el formato de mensajes de WhatsApp
func TestWhatsAppMessageFormat(t *testing.T) {
	// Create un servidor WebSocket simple
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	receivedMessages := make([]map[string]interface{}, 0)
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		defer conn.Close()

		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				break
			}

			var msg map[string]interface{}
			if err := json.Unmarshal(message, &msg); err != nil {
				continue
			}

			receivedMessages = append(receivedMessages, msg)
		}
	}))
	defer server.Close()

	wsURL := strings.Replace(server.URL, "https://", "wss://", 1)

	cfg := config.WhatsAppConfig{
		Enabled:   true,
		BridgeURL: wsURL,
		AllowFrom: []string{},
	}

	msgBus := bus.NewMessageBus()
	channel, err := NewWhatsAppChannel(cfg, msgBus)
	if err != nil {
		t.Fatalf("Error creating WhatsApp channel: %v", err)
	}

	ctx := context.Background()
	err = channel.Start(ctx)
	if err != nil {
		t.Fatalf("Error starting WhatsApp channel: %v", err)
	}

	// Send mensaje de prueba
	testMsg := bus.OutboundMessage{
		Channel: "whatsapp",
		ChatID:  "+1234567890",
		Content: "Test message content",
	}

	err = channel.Send(ctx, testMsg)
	if err != nil {
		t.Errorf("Error sending message: %v", err)
	}

	// Dar tiempo para que el mensaje se procese
	time.Sleep(100 * time.Millisecond)

	// Verify que el mensaje fue recibido con el formato correcto
	if len(receivedMessages) == 0 {
		t.Fatal("No messages received")
	}

	lastMsg := receivedMessages[len(receivedMessages)-1]
	if lastMsg["type"] != "message" {
		t.Errorf("Expected message type 'message', got %v", lastMsg["type"])
	}
	if lastMsg["to"] != "+1234567890" {
		t.Errorf("Expected to '+1234567890', got %v", lastMsg["to"])
	}
	if lastMsg["content"] != "Test message content" {
		t.Errorf("Expected content 'Test message content', got %v", lastMsg["content"])
	}

	channel.Stop(ctx)
}

// TestWhatsAppAllowFrom prueba la funcionalidad de allow_from
func TestWhatsAppAllowFrom(t *testing.T) {
	cfg := config.WhatsAppConfig{
		Enabled:   true,
		BridgeURL: "ws://localhost:3001",
		AllowFrom: []string{"allowed-user", "+1234567890"},
	}

	msgBus := bus.NewMessageBus()
	channel, err := NewWhatsAppChannel(cfg, msgBus)
	if err != nil {
		t.Fatalf("Error creating WhatsApp channel: %v", err)
	}

	// Test usuarios permitidos
	if !channel.IsAllowed("allowed-user") {
		t.Error("allowed-user should be allowed")
	}
	if !channel.IsAllowed("+1234567890") {
		t.Error("+1234567890 should be allowed")
	}

	// Test usuario no permitido
	if channel.IsAllowed("not-allowed-user") {
		t.Error("not-allowed-user should not be allowed")
	}

	// Test empty list (should allow all)
	emptyAllowChannel, err := NewWhatsAppChannel(config.WhatsAppConfig{
		Enabled:   true,
		BridgeURL: "ws://localhost:3001",
		AllowFrom: []string{},
	}, msgBus)
	if err != nil {
		t.Fatalf("Error creating WhatsApp channel with empty allow_from: %v", err)
	}

	if !emptyAllowChannel.IsAllowed("any-user") {
		t.Error("With empty allow_from, any user should be allowed")
	}
}

// TestWhatsAppConfigDefaultValues prueba los valores por defecto
func TestWhatsAppConfigDefaultValues(t *testing.T) {
	cfg := config.DefaultConfig()

	// Verify valores por defecto
	if cfg.Channels.WhatsApp.Enabled {
		t.Error("WhatsApp should be disabled by default")
	}
	if cfg.Channels.WhatsApp.BridgeURL != "ws://localhost:3001" {
		t.Errorf("Default bridge_url should be 'ws://localhost:3001', got %s", cfg.Channels.WhatsApp.BridgeURL)
	}
	if len(cfg.Channels.WhatsApp.AllowFrom) != 0 {
		t.Error("Default allow_from should be empty")
	}
}
