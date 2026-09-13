package config

import "os"

type Config struct {
	Port          string
	AIAPIURL      string
	AIAPIKey      string
	AIFallbackURL string
	AIFallbackKey string
	RedisURL      string
}

func Load() *Config {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	aiURL := os.Getenv("AI_API_URL")
	if aiURL == "" {
		aiURL = "https://api.groq.com/openai/v1/chat/completions"
	}
	return &Config{
		Port:          port,
		AIAPIURL:      aiURL,
		AIAPIKey:      os.Getenv("AI_API_KEY"),
		AIFallbackURL: os.Getenv("AI_FALLBACK_URL"),
		AIFallbackKey: os.Getenv("AI_FALLBACK_KEY"),
		RedisURL:      os.Getenv("REDIS_URL"),
	}
}
