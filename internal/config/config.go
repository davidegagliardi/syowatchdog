package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	// Image monitoring settings
	ImageURL      string        `yaml:"image_url" env:"IMAGE_URL"`
	CheckInterval time.Duration `yaml:"check_interval" env:"CHECK_INTERVAL"`
	StoragePath   string        `yaml:"storage_path" env:"STORAGE_PATH"`

	// Telegram settings
	TelegramBotToken string `yaml:"telegram_bot_token" env:"TELEGRAM_BOT_TOKEN"`
	TelegramChatID   int64  `yaml:"telegram_chat_id" env:"TELEGRAM_CHAT_ID"`

	// General settings
	LogLevel string `yaml:"log_level" env:"LOG_LEVEL"`
}

func LoadConfig(configPath string) (*Config, error) {
	// Default values
	config := &Config{
		CheckInterval: 5 * time.Minute,
		StoragePath:   "./data",
		LogLevel:      "info",
	}

	// Load from file if provided
	if configPath != "" {
		data, err := os.ReadFile(configPath)
		if err != nil {
			return nil, err
		}
		if err := yaml.Unmarshal(data, config); err != nil {
			return nil, err
		}
	}

	// Override with environment variables
	if url := os.Getenv("IMAGE_URL"); url != "" {
		config.ImageURL = url
	}
	if interval := os.Getenv("CHECK_INTERVAL"); interval != "" {
		if d, err := time.ParseDuration(interval); err == nil {
			config.CheckInterval = d
		}
	}
	if token := os.Getenv("TELEGRAM_BOT_TOKEN"); token != "" {
		config.TelegramBotToken = token
	}
	if chatID := os.Getenv("TELEGRAM_CHAT_ID"); chatID != "" {
		if id, err := strconv.ParseInt(chatID, 10, 64); err == nil {
			config.TelegramChatID = id
		}
	}
	if path := os.Getenv("STORAGE_PATH"); path != "" {
		config.StoragePath = path
	}

	return config, nil
}

func (c *Config) Validate() error {
	if c.ImageURL == "" {
		return fmt.Errorf("image_url is required")
	}
	if c.TelegramBotToken == "" {
		return fmt.Errorf("telegram_bot_token is required")
	}
	if c.TelegramChatID == 0 {
		return fmt.Errorf("telegram_chat_id is required")
	}
	return nil
}
