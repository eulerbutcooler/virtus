package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Env    string `mapstructure:"APP_ENV"`
	Server ServerConfig
	DB     DBConfig
	Redis  RedisConfig
	JWT    JWTConfig
	Log    LogConfig
	Stripe StripeConfig
}

type StripeConfig struct {
	SecretKey     string `mapstructure:"STRIPE_SECRET_KEY"`
	WebhookSecret string `mapstructure:"STRIPE_WEBHOOK_SECRET"`
}

type ServerConfig struct {
	Host         string        `mapstructure:"SERVER_HOST"`
	Port         int           `mapstructure:"SERVER_PORT"`
	ReadTimeout  time.Duration `mapstructure:"SERVER_READ_TIMEOUT"`
	WriteTimeout time.Duration `mapstructure:"SERVER_WRITE_TIMEOUT"`
	IdleTimeout  time.Duration `mapstructure:"SERVER_IDLE_TIMEOUT"`
}

type DBConfig struct {
	Host     string `mapstructure:"DB_HOST"`
	Port     int    `mapstructure:"DB_PORT"`
	User     string `mapstructure:"DB_USER"`
	Password string `mapstructure:"DB_PASSWORD"`
	Name     string `mapstructure:"DB_NAME"`
	SSLMode  string `mapstructure:"DB_SSLMODE"`
}

type RedisConfig struct {
	Addr     string `mapstructure:"REDIS_ADDR"`
	Password string `mapstructure:"REDIS_PASSWORD"`
	DB       int    `mapstructure:"REDIS_DB"`
}

type JWTConfig struct {
	Secret          string        `mapstructure:"JWT_SECRET"`
	AccessTokenTTL  time.Duration `mapstructure:"JWT_ACCESS_TTL"`
	RefreshTokenTTL time.Duration `mapstructure:"JWT_REFRESH_TTL"`
}

type LogConfig struct {
	Level  string `mapstructure:"LOG_LEVEL"`
	Format string `mapstructure:"LOG_FORMAT"`
}

func Load() (*Config, error) {
	v := viper.New()
	v.SetConfigFile(".env")
	v.SetConfigType("env")
	_ = v.ReadInConfig()

	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	setDefaults(v)

	cfg := &Config{
		Env: v.GetString("APP_ENV"),
		Server: ServerConfig{
			Host:         v.GetString("SERVER_HOST"),
			Port:         v.GetInt("SERVER_PORT"),
			ReadTimeout:  v.GetDuration("SERVER_READ_TIMEOUT"),
			WriteTimeout: v.GetDuration("SERVER_WRITE_TIMEOUT"),
			IdleTimeout:  v.GetDuration("SERVER_IDLE_TIMEOUT"),
		},
		DB: DBConfig{
			Host:     v.GetString("DB_HOST"),
			Port:     v.GetInt("DB_PORT"),
			User:     v.GetString("DB_USER"),
			Password: v.GetString("DB_PASSWORD"),
			Name:     v.GetString("DB_NAME"),
			SSLMode:  v.GetString("DB_SSLMODE"),
		},
		Redis: RedisConfig{
			Addr:     v.GetString("REDIS_ADDR"),
			Password: v.GetString("REDIS_PASSWORD"),
			DB:       v.GetInt("REDIS_DB"),
		},
		JWT: JWTConfig{
			Secret:          v.GetString("JWT_SECRET"),
			AccessTokenTTL:  v.GetDuration("JWT_ACCESS_TTL"),
			RefreshTokenTTL: v.GetDuration("JWT_REFRESH_TTL"),
		},
		Log: LogConfig{
			Level:  v.GetString("LOG_LEVEL"),
			Format: v.GetString("LOG_FORMAT"),
		},
		Stripe: StripeConfig{
			SecretKey:     v.GetString("STRIPE_SECRET_KEY"),
			WebhookSecret: v.GetString("STRIPE_WEBHOOK_SECRET"),
		},
	}

	if err := validate(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

func setDefaults(v *viper.Viper) {
	v.SetDefault("APP_ENV", "development")

	v.SetDefault("SERVER_HOST", "0.0.0.0")
	v.SetDefault("SERVER_PORT", 8080)
	v.SetDefault("SERVER_READ_TIMEOUT", "10s")
	v.SetDefault("SERVER_WRITE_TIMEOUT", "30s")
	v.SetDefault("SERVER_IDLE_TIMEOUT", "60s")

	v.SetDefault("DB_HOST", "localhost")
	v.SetDefault("DB_PORT", 5432)
	v.SetDefault("DB_USER", "virtus")
	v.SetDefault("DB_NAME", "virtus")
	v.SetDefault("DB_SSLMODE", "disable")

	v.SetDefault("REDIS_ADDR", "localhost:6379")
	v.SetDefault("REDIS_PASSWORD", "")
	v.SetDefault("REDIS_DB", 0)

	v.SetDefault("JWT_ACCESS_TTL", "15m")
	v.SetDefault("JWT_REFRESH_TTL", "168h")

	v.SetDefault("LOG_LEVEL", "info")
	v.SetDefault("LOG_FORMAT", "json")
}

func validate(cfg *Config) error {
	if strings.TrimSpace(cfg.DB.Password) == "" {
		return fmt.Errorf("required env var %q is not set", "DB_PASSWORD")
	}
	if strings.TrimSpace(cfg.JWT.Secret) == "" {
		return fmt.Errorf("required env var %q is not set", "JWT_SECRET")
	}
	return nil
}
