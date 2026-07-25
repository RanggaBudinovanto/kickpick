package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AppEnv            string
	AppPort           string
	AppURL            string
	DatabaseURL       string
	RedisURL          string
	JWTAccessSecret   string
	JWTRefreshSecret  string
	CORSAllowedOrigin string
	ResendAPIKey      string
	EmailFrom         string
	// ShopeePartnerID/ShopeePartnerKey are unset until a Shopee Affiliate
	// Program/Open Platform account is approved (PENDING.md) — the shopee
	// adapter (internal/scraper/shopee) checks Configured() and stays
	// inactive rather than error on every scheduled run until then.
	ShopeePartnerID  string
	ShopeePartnerKey string
}

func Load() *Config {
	_ = godotenv.Load()

	return &Config{
		AppEnv:            getEnv("APP_ENV", "development"),
		AppPort:           getEnv("APP_PORT", "8080"),
		AppURL:            getEnv("APP_URL", "http://localhost:3000"),
		DatabaseURL:       getEnv("DATABASE_URL", ""),
		RedisURL:          getEnv("REDIS_URL", ""),
		JWTAccessSecret:   getEnv("JWT_ACCESS_SECRET", ""),
		JWTRefreshSecret:  getEnv("JWT_REFRESH_SECRET", ""),
		CORSAllowedOrigin: getEnv("CORS_ALLOWED_ORIGIN", "http://localhost:3000"),
		ResendAPIKey:      getEnv("RESEND_API_KEY", ""),
		EmailFrom:         getEnv("EMAIL_FROM", "KickPick <noreply@kickpick.id>"),
		ShopeePartnerID:   getEnv("SHOPEE_PARTNER_ID", ""),
		ShopeePartnerKey:  getEnv("SHOPEE_PARTNER_KEY", ""),
	}
}

func getEnv(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		return v
	}
	return fallback
}
