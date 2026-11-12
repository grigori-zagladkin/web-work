package config

import (
	"errors"
	"log/slog"
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

const MIN_JWT_LENGTH = 32

type Config struct {
	Env      string     `env:"APP_ENV" envDefault:"dev"`
	HTTPAddr string     `env:"HTTP_ADDR" envDefault:"8080"`
	LogLevel slog.Level `env:"LOG_LEVEL" envDefault:"info"`

	DatabaseURL string `env:"DATABASE_URL,required,notEmpty"`

	JWTSecret       string        `env:"JWT_SECRET,required,notEmpty"`
	AccessTokenTTL  time.Duration `env:"ACCESS_TOKEN_TTL" envDefault:"15m"`
	RefreshTokenTTL time.Duration `env:"REFRESH_TOKEN_TTL" envDefault:"720h"`

	CookieSecure bool `env:"COOKIE_SECURE" envDefault:"true"`

	TrustedOrigins []string `env:"TRUSTED_ORIGINS"`
}

func Load() (Config, error) {
	_ = godotenv.Load()

	cfg, err := env.ParseAs[Config]()
	if err != nil {
		return Config{}, err
	}
	if len(cfg.JWTSecret) < MIN_JWT_LENGTH {
		return Config{}, errors.New("JWT_Secret must be at least 32 characters")
	}
	return cfg, nil
}
