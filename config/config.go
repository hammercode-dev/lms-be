package config

import (
	"strings"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var (
	c *Config
)

type (
	Config struct {
		APP_ENV         string
		APP_NAME        string
		APP_PORT        string
		DB_POSTGRES_DSN string
		JWT_SECRET_KEY  string

		SMTP_HOST     string
		SMTP_PORT     string
		SMTP_EMAIL    string
		SMTP_PASSWORD string

		CORS_ALLOWED_ORIGINS []string
		CORS_ALLOWED_METHODS []string
		CORS_ALLOWED_HEADERS []string

		BaseURL     string
		BASE_URL_FE string
	}
)

func LoadConfig() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		logrus.Fatal("Error loading .env file", err)
	}

	// Enable VIPER to read Environment Variables
	viper.AutomaticEnv()
}

func GetConfig() Config {

	if c == nil {
		// default cors
		origins := []string{"*"}
		methods := []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}
		headers := []string{"Accept", "Authorization", "Content-Type"}

		corsOrigins := viper.GetString("CORS_ALLOWED_ORIGINS")
		if corsOrigins != "" {
			origins = strings.Split(corsOrigins, ",")
		}

		corsHeaders := viper.GetString("CORS_ALLOWED_HEADERS")
		if corsHeaders != "" {
			headers = strings.Split(corsHeaders, ",")
		}

		corsMethods := viper.GetString("CORS_ALLOWED_METHODS")
		if corsHeaders != "" {
			methods = strings.Split(corsMethods, ",")
		}

		c = &Config{
			APP_ENV:              viper.GetString("APP_ENV"),
			APP_NAME:             viper.GetString("APP_NAME"),
			APP_PORT:             viper.GetString("APP_PORT"),
			DB_POSTGRES_DSN:      viper.GetString("DB_POSTGRES_DSN"),
			JWT_SECRET_KEY:       viper.GetString("JWT_SECRET_KEY"),
			SMTP_HOST:            viper.GetString("SMTP_HOST"),
			SMTP_PORT:            viper.GetString("SMTP_PORT"),
			SMTP_EMAIL:           viper.GetString("SMTP_EMAIL"),
			SMTP_PASSWORD:        viper.GetString("SMTP_PASSWORD"),
			CORS_ALLOWED_ORIGINS: origins,
			CORS_ALLOWED_METHODS: methods,
			CORS_ALLOWED_HEADERS: headers,
			BaseURL:              viper.GetString("BASE_URL"),
			BASE_URL_FE:          viper.GetString("BASE_URL_FE"),
		}
	}

	return *c
}
