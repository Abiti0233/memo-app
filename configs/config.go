package configs

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost		 string
	DBPort		 string
	DBUser		 string
	DBPassword	 string
	DBName		 string
	ServerPort	 string
	OAuthClientID    string
	OAuthClientSecret string
	OAuthRedirectURL string
	JWTSecret         string
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
			log.Println("No .env file found")
	}

	return &Config{
			DBHost:     getEnv("DB_HOST", "localhost"),
			DBPort:     getEnv("DB_PORT", "5432"),
			DBUser:     getEnv("DB_USER", "postgres"),
			DBPassword: getEnv("DB_PASSWORD", "password"),
			DBName:     getEnv("DB_NAME", "memo-app-db"),
			ServerPort: getEnv("SERVER_PORT", "8081"),
			OAuthClientID:    getEnv("OAUTH_CLIENT_ID", ""),
      OAuthClientSecret: getEnv("OAUTH_CLIENT_SECRET", ""),
      OAuthRedirectURL: getEnv("OAUTH_REDIRECT_URL", "http://localhost:8081/auth/callback"),
			JWTSecret:         getEnv("JWT_SECRET", ""),
	}
}

func getEnv(key, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
			return value
	}
	return defaultVal
}