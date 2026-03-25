package config

import (
	"fmt"
	"log"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	// Database
	DBHost     string
	DBPort     string
	DBName     string
	DBUser     string
	DBPassword string
	DBSSLMode  string

	// Server
	ServerPort string

	// JWT
	JWTSecret         string
	JWTExpiryMinutes  int
	RefreshExpiryDays int
}

var AppConfig *Config

func LoadConfig() error {
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()

	// Set defaults
	viper.SetDefault("DB_HOST", "localhost")
	viper.SetDefault("DB_PORT", "5432")
	viper.SetDefault("DB_NAME", "meet_sushruta")
	viper.SetDefault("DB_USER", "postgres")
	viper.SetDefault("DB_PASSWORD", "password")
	viper.SetDefault("DB_SSLMODE", "disable")
	viper.SetDefault("SERVER_PORT", "8080")
	viper.SetDefault("JWT_SECRET", "your_jwt_secret_key_change_in_production")
	viper.SetDefault("JWT_EXPIRY_MINUTES", 60)
	viper.SetDefault("REFRESH_EXPIRY_DAYS", 7)

	// Try to read config file, but don't fail if it doesn't exist
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			log.Printf("Warning: Error reading config file: %v", err)
		}
	}

	AppConfig = &Config{
		DBHost:            viper.GetString("DB_HOST"),
		DBPort:            viper.GetString("DB_PORT"),
		DBName:            viper.GetString("DB_NAME"),
		DBUser:            viper.GetString("DB_USER"),
		DBPassword:        viper.GetString("DB_PASSWORD"),
		DBSSLMode:         viper.GetString("DB_SSLMODE"),
		ServerPort:        viper.GetString("SERVER_PORT"),
		JWTSecret:         viper.GetString("JWT_SECRET"),
		JWTExpiryMinutes:  viper.GetInt("JWT_EXPIRY_MINUTES"),
		RefreshExpiryDays: viper.GetInt("REFRESH_EXPIRY_DAYS"),
	}

	log.Println("Configuration loaded successfully")
	return nil
}

func GetConfig() *Config {
	if AppConfig == nil {
		log.Fatal("Config not loaded. Call LoadConfig() first")
	}
	return AppConfig
}

// GetDSN returns the PostgreSQL connection string
func (c *Config) GetDSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.DBHost,
		c.DBPort,
		c.DBUser,
		c.DBPassword,
		c.DBName,
		c.DBSSLMode,
	)
}

// GetJWTExpiryDuration returns JWT expiry as duration
func (c *Config) GetJWTExpiryDuration() time.Duration {
	return time.Duration(c.JWTExpiryMinutes) * time.Minute
}

// GetRefreshExpiryDuration returns refresh token expiry as duration
func (c *Config) GetRefreshExpiryDuration() time.Duration {
	return time.Duration(c.RefreshExpiryDays) * 24 * time.Hour
}
