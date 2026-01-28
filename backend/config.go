package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	OpenAIAPIKey string
	OpenAIAPIURL string
	Port         string
}

func LoadConfig() (*Config, error) {
	// .envファイルを読み込み（存在しない場合は無視）
	if err := godotenv.Load(); err != nil {
		// .envファイルが存在しない場合は警告を出力するが、処理は続行
		fmt.Printf("Warning: .env file not found: %v\n", err)
	}

	config := &Config{
		OpenAIAPIKey: getEnv("OPENAI_API_KEY", ""),
		OpenAIAPIURL: getEnv("OPENAI_API_URL", "https://api.openai.com/v1/chat/completions"),
		Port:         getEnv("PORT", "8080"),
	}

	return config, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
