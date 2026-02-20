package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/caarlos0/env/v11"
)

// FlexibleStringSlice is a []string that also accepts JSON numbers,
// so allow_from can contain both "123" and 123.
type FlexibleStringSlice []string

func (f *FlexibleStringSlice) UnmarshalJSON(data []byte) error {
	// Try []string first
	var ss []string
	if err := json.Unmarshal(data, &ss); err == nil {
		*f = ss
		return nil
	}

	// Try []interface{} to handle mixed types
	var raw []interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	*f = make([]string, len(raw))
	for i, v := range raw {
		switch val := v.(type) {
		case string:
			(*f)[i] = val
		case float64:
			(*f)[i] = fmt.Sprintf("%.0f", val)
		default:
			return fmt.Errorf("invalid type in allow_from: %T", v)
		}
	}
	return nil
}

// Config represents the main configuration structure
type Config struct {
	mu sync.RWMutex

	// Global settings
	Debug       bool   `json:"debug" env:"PICOCLAW_DEBUG"`
	LogLevel    string `json:"log_level" env:"PICOCLAW_LOG_LEVEL"`
	BindAddress string `json:"bind_address" env:"PICOCLAW_BIND_ADDRESS"`
	
	// Security settings
	EnableAuth bool   `json:"enable_auth" env:"PICOCLAW_ENABLE_AUTH"`
	SecretKey  string `json:"secret_key" env:"PICOCLAW_SECRET_KEY"`
	
	// AI settings
	AI AIConfig `json:"ai"`
	
	// Channel configurations
	Channels ChannelsConfig `json:"channels"`
	
	// Raw JSON for unknown fields
	Raw json.RawMessage `json:"-"`
}

// AIConfig represents AI provider configuration
type AIConfig struct {
	DefaultProvider string            `json:"default_provider" env:"PICOCLAW_AI_DEFAULT_PROVIDER"`
	Providers       []ProviderConfig  `json:"providers"`
}

// ProviderConfig represents a single AI provider configuration
type ProviderConfig struct {
	Name        string            `json:"name"`
	Type        string            `json:"type"`
	Endpoint    string            `json:"endpoint"`
	APIKey      string            `json:"api_key"`
	Model       string            `json:"model"`
	MaxTokens   int               `json:"max_tokens"`
	Temperature float64           `json:"temperature"`
	Headers     map[string]string `json:"headers,omitempty"`
}

// ChannelsConfig represents all channel configurations
type ChannelsConfig struct {
	WhatsApp WhatsAppConfig `json:"whatsapp"`
	Telegram TelegramConfig `json:"telegram"`
	LINE     LINEConfig     `json:"line"`
	OneBot   OneBotConfig   `json:"onebot"`
}

// WhatsAppConfig represents WhatsApp channel configuration
type WhatsAppConfig struct {
	Enabled   bool                `json:"enabled" env:"PICOCLAW_CHANNELS_WHATSAPP_ENABLED"`
	BridgeURL string              `json:"bridge_url" env:"PICOCLAW_CHANNELS_WHATSAPP_BRIDGE_URL"`
	AllowFrom FlexibleStringSlice `json:"allow_from" env:"PICOCLAW_CHANNELS_WHATSAPP_ALLOW_FROM"`
	
	// Facebook WhatsApp Business API configuration
	FBPhoneNumberID string `json:"fb_phone_number_id" env:"PICOCLAW_CHANNELS_WHATSAPP_FB_PHONE_NUMBER_ID"`
	FBAccessToken   string `json:"fb_access_token" env:"PICOCLAW_CHANNELS_WHATSAPP_FB_ACCESS_TOKEN"`
	FBAPIVersion    string `json:"fb_api_version" env:"PICOCLAW_CHANNELS_WHATSAPP_FB_API_VERSION"`
}

// TelegramConfig represents Telegram channel configuration
type TelegramConfig struct {
	Enabled   bool                `json:"enabled" env:"PICOCLAW_CHANNELS_TELEGRAM_ENABLED"`
	Token     string              `json:"token" env:"PICOCLAW_CHANNELS_TELEGRAM_TOKEN"`
	Proxy     string              `json:"proxy" env:"PICOCLAW_CHANNELS_TELEGRAM_PROXY"`
	AllowFrom FlexibleStringSlice `json:"allow_from" env:"PICOCLAW_CHANNELS_TELEGRAM_ALLOW_FROM"`
}

// LINEConfig represents LINE channel configuration
type LINEConfig struct {
	Enabled           bool                `json:"enabled" env:"PICOCLAW_CHANNELS_LINE_ENABLED"`
	ChannelSecret     string              `json:"channel_secret" env:"PICOCLAW_CHANNELS_LINE_CHANNEL_SECRET"`
	ChannelAccessToken string             `json:"channel_access_token" env:"PICOCLAW_CHANNELS_LINE_CHANNEL_ACCESS_TOKEN"`
	AllowFrom         FlexibleStringSlice `json:"allow_from" env:"PICOCLAW_CHANNELS_LINE_ALLOW_FROM"`
}

// OneBotConfig represents OneBot channel configuration
type OneBotConfig struct {
	Enabled   bool                `json:"enabled" env:"PICOCLAW_CHANNELS_ONEBOT_ENABLED"`
	Endpoint  string              `json:"endpoint" env:"PICOCLAW_CHANNELS_ONEBOT_ENDPOINT"`
	AccessToken string            `json:"access_token" env:"PICOCLAW_CHANNELS_ONEBOT_ACCESS_TOKEN"`
	AllowFrom FlexibleStringSlice `json:"allow_from" env:"PICOCLAW_CHANNELS_ONEBOT_ALLOW_FROM"`
}

// Load loads configuration from file and environment
func Load(configPath string) (*Config, error) {
	cfg := &Config{}
	
	// Load from file if exists
	if configPath != "" {
		if err := loadFromFile(configPath, cfg); err != nil {
			return nil, fmt.Errorf("failed to load config from file: %w", err)
		}
	}
	
	// Override with environment variables
	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("failed to parse environment: %w", err)
	}
	
	// Apply defaults
	cfg.applyDefaults()
	
	// Validate configuration
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}
	
	return cfg, nil
}

// loadFromFile loads configuration from JSON file
func loadFromFile(configPath string, cfg *Config) error {
	// Expand home directory
	configPath = expandPath(configPath)
	
	// Read file
	data, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil // File doesn't exist, use defaults
		}
		return fmt.Errorf("failed to read config file: %w", err)
	}
	
	// Parse JSON
	if err := json.Unmarshal(data, cfg); err != nil {
		return fmt.Errorf("failed to parse config JSON: %w", err)
	}
	
	return nil
}

// applyDefaults applies default values
func (c *Config) applyDefaults() {
	if c.BindAddress == "" {
		c.BindAddress = "0.0.0.0:8080"
	}
	if c.LogLevel == "" {
		c.LogLevel = "info"
	}
	if c.AI.DefaultProvider == "" {
		c.AI.DefaultProvider = "openai"
	}
	
	// Set default Facebook API version
	if c.Channels.WhatsApp.FBAPIVersion == "" {
		c.Channels.WhatsApp.FBAPIVersion = "v22.0"
	}
}

// Validate validates the configuration
func (c *Config) Validate() error {
	// Validate channels
	if c.Channels.WhatsApp.Enabled {
		// Check if either bridge URL or Facebook API credentials are provided
		hasBridge := c.Channels.WhatsApp.BridgeURL != ""
		hasFBAPI := c.Channels.WhatsApp.FBPhoneNumberID != "" && c.Channels.WhatsApp.FBAccessToken != ""
		
		if !hasBridge && !hasFBAPI {
			return fmt.Errorf("whatsapp: either bridge_url or facebook api credentials (fb_phone_number_id and fb_access_token) must be provided")
		}
		
		if hasBridge && hasFBAPI {
			return fmt.Errorf("whatsapp: cannot use both bridge_url and facebook api simultaneously")
		}
	}
	
	return nil
}

// GetProvider returns a provider configuration by name
func (c *Config) GetProvider(name string) (*ProviderConfig, error) {
	for _, provider := range c.AI.Providers {
		if provider.Name == name {
			return &provider, nil
		}
	}
	return nil, fmt.Errorf("provider %s not found", name)
}

// IsProviderEnabled checks if a provider is enabled
func (c *Config) IsProviderEnabled(name string) bool {
	provider, err := c.GetProvider(name)
	return err == nil && provider.APIKey != ""
}

// Lock locks the configuration for reading
func (c *Config) RLock() {
	c.mu.RLock()
}

// RUnlock unlocks the configuration after reading
func (c *Config) RUnlock() {
	c.mu.RUnlock()
}

// Lock locks the configuration for writing
func (c *Config) Lock() {
	c.mu.Lock()
}

// Unlock unlocks the configuration after writing
func (c *Config) Unlock() {
	c.mu.Unlock()
}

// expandPath expands ~ to home directory
func expandPath(path string) string {
	if path == "" {
		return path
	}
	
	if path[0] == '~' {
		home, err := os.UserHomeDir()
		if err == nil {
			path = filepath.Join(home, path[1:])
		}
	}
	
	return path
}