package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	TwitchClientID     string
	TwitchClientSecret string
	TwitchRedirectURI  string
	Port               string
	Env                string
}

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system env vars")
	}
	return &Config{
		TwitchClientID:     mustGet("TWITCH_CLIENT_ID"),
		TwitchClientSecret: mustGet("TWITCH_CLIENT_SECRET"),
		TwitchRedirectURI:  mustGet("TWITCH_REDIRECT_URI"),
		Port:               getOrDefault("PORT", "3000"),
		Env:                getOrDefault("ENV", "development"),
	}
}

func mustGet(key string) string {
	val := os.Getenv(key)
	if val == "" {
		log.Fatalf("required env var %s is not set", key)
	}
	return val
}

func getOrDefault(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}
