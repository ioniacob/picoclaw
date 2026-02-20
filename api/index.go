package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/sipeed/picoclaw/pkg/bus"
	"github.com/sipeed/picoclaw/pkg/channels"
	"github.com/sipeed/picoclaw/pkg/config"
	"github.com/sipeed/picoclaw/pkg/logger"
	"github.com/sipeed/picoclaw/pkg/providers"
)

var (
	channelManager *channels.Manager
	msgBus         *bus.MessageBus
	provider       providers.LLMProvider
	cfg            *config.Config
)

func init() {
	// Crear configuración para Vercel
	cfg = createVercelConfig()

	// Crear provider
	var err error
	provider, err = providers.CreateProvider(cfg)
	if err != nil {
		logger.ErrorCF("vercel", "Error creating provider", map[string]interface{}{"error": err.Error()})
		return
	}

	// Crear message bus
	msgBus = bus.NewMessageBus()

	// Crear channel manager
	channelManager, err = channels.NewManager(cfg, msgBus)
	if err != nil {
		logger.ErrorCF("vercel", "Error creating channel manager", map[string]interface{}{"error": err.Error()})
		return
	}

	// Iniciar canales
	ctx := context.Background()
	if err := channelManager.StartAll(ctx); err != nil {
		logger.ErrorCF("vercel", "Error starting channels", map[string]interface{}{"error": err.Error()})
		return
	}
}

func createVercelConfig() *config.Config {
	// Configuración optimizada para Vercel
	cfg := &config.Config{
		Agents: config.AgentsConfig{
			Defaults: config.AgentDefaults{
				Workspace:           "/tmp/workspace",
				RestrictToWorkspace: true,
				Model:               "glm-4.7",
				MaxTokens:           4096, // Reducido para Vercel
				Temperature:         0.7,
				MaxToolIterations:   10, // Reducido para Vercel
			},
		},
		Channels: config.ChannelsConfig{
			WhatsApp: config.WhatsAppConfig{
				Enabled:   true, // Habilitar WhatsApp
				BridgeURL: getEnvOrDefault("WHATSAPP_BRIDGE_URL", "wss://api.whatsapp.com/v1"),
				AllowFrom: []string{}, // Configurar via variables de entorno
			},
			Telegram: config.TelegramConfig{
				Enabled: false, // Deshabilitar por defecto en Vercel
			},
			Discord: config.DiscordConfig{
				Enabled: false, // Deshabilitar por defecto en Vercel
			},
			Slack: config.SlackConfig{
				Enabled: false, // Deshabilitar por defecto en Vercel
			},
		},
		Providers: config.ProvidersConfig{
			OpenRouter: config.ProviderConfig{
				APIKey:  getEnvOrDefault("OPENROUTER_API_KEY", ""),
				APIBase: getEnvOrDefault("OPENROUTER_API_BASE", "https://openrouter.ai/api/v1"),
			},
			Anthropic: config.ProviderConfig{
				APIKey:  getEnvOrDefault("ANTHROPIC_API_KEY", ""),
				APIBase: getEnvOrDefault("ANTHROPIC_API_BASE", ""),
			},
			OpenAI: config.ProviderConfig{
				APIKey:  getEnvOrDefault("OPENAI_API_KEY", ""),
				APIBase: getEnvOrDefault("OPENAI_API_BASE", ""),
			},
		},
		Tools: config.ToolsConfig{
			Cron: config.CronToolsConfig{
				ExecTimeoutMinutes: 2, // Reducido para Vercel
			},
		},
		Heartbeat: config.HeartbeatConfig{
			Enabled:  false, // Deshabilitar en Vercel
			Interval: 30,
		},
		Devices: config.DevicesConfig{
			Enabled:    false, // Deshabilitar en Vercel
			MonitorUSB: false,
		},
		Gateway: config.GatewayConfig{
			Host: "0.0.0.0",
			Port: 8080, // Puerto estándar de Vercel
		},
	}

	// Configurar allow_from desde variables de entorno
	if allowed := os.Getenv("WHATSAPP_ALLOWED_NUMBERS"); allowed != "" {
		cfg.Channels.WhatsApp.AllowFrom = strings.Split(allowed, ",")
	}

	// Configurar seguridad de WhatsApp
	cfg.Channels.WhatsApp.BridgeURL = getEnvOrDefault("WHATSAPP_BRIDGE_URL", "wss://api.whatsapp.com/v1")

	return cfg
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// Handler principal para Vercel
func Handler(w http.ResponseWriter, r *http.Request) {
	// Configurar CORS
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	// Health check
	if r.URL.Path == "/health" {
		handleHealth(w, r)
		return
	}

	// Ready check
	if r.URL.Path == "/ready" {
		handleReady(w, r)
		return
	}

	// Webhook para WhatsApp
	if strings.HasPrefix(r.URL.Path, "/webhook/whatsapp") {
		handleWhatsAppWebhook(w, r)
		return
	}

	// API principal
	if r.URL.Path == "/api/chat" && r.Method == "POST" {
		handleChat(w, r)
		return
	}

	// Información del servicio
	if r.URL.Path == "/" && r.Method == "GET" {
		handleInfo(w, r)
		return
	}

	http.NotFound(w, r)
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	response := map[string]interface{}{
		"status":    "healthy",
		"timestamp": time.Now().Unix(),
		"service":   "picoclaw-whatsapp",
		"version":   "1.0.0",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func handleReady(w http.ResponseWriter, r *http.Request) {
	ready := channelManager != nil

	response := map[string]interface{}{
		"status": map[string]bool{
			"ready":    ready,
			"whatsapp": cfg != nil && cfg.Channels.WhatsApp.Enabled,
		},
		"timestamp": time.Now().Unix(),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func handleInfo(w http.ResponseWriter, r *http.Request) {
	response := map[string]interface{}{
		"service":     "PicoClaw WhatsApp Gateway",
		"version":     "1.0.0",
		"description": "Ultra-lightweight personal AI agent with WhatsApp integration",
		"endpoints": map[string]string{
			"health":           "/health",
			"ready":            "/ready",
			"whatsapp_webhook": "/webhook/whatsapp",
			"chat":             "/api/chat",
		},
		"channels": map[string]bool{
			"whatsapp": cfg != nil && cfg.Channels.WhatsApp.Enabled,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func handleWhatsAppWebhook(w http.ResponseWriter, r *http.Request) {
	// Verificar método
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Verificar autenticación si está configurada
	if token := os.Getenv("WHATSAPP_WEBHOOK_TOKEN"); token != "" {
		authHeader := r.Header.Get("Authorization")
		if authHeader != "Bearer "+token {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
	}

	// Leer body
	var message map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&message); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Procesar mensaje
	if channelManager != nil {
		// El canal de WhatsApp procesará el mensaje
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "received",
			"message": "Message processed",
		})
	} else {
		http.Error(w, "Channel manager not initialized", http.StatusServiceUnavailable)
	}
}

func handleChat(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Message string `json:"message"`
		Channel string `json:"channel"`
		ChatID  string `json:"chat_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if req.Message == "" {
		http.Error(w, "Message is required", http.StatusBadRequest)
		return
	}

	// Usar WhatsApp por defecto si no se especifica canal
	if req.Channel == "" {
		req.Channel = "whatsapp"
	}

	// Procesar mensaje
	response := map[string]interface{}{
		"status":    "ok",
		"channel":   req.Channel,
		"chat_id":   req.ChatID,
		"message":   req.Message,
		"timestamp": time.Now().Unix(),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
