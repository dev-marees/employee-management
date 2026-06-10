package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/viper"
)

// Config holds all configuration for the application, populated from
// environment variables (and an optional .env file).
type Config struct {
	App      AppConfig
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
	Log      LogConfig
}

type AppConfig struct {
	Name string `mapstructure:"APP_NAME"`
	Env  string `mapstructure:"APP_ENV"`
}

type ServerConfig struct {
	Port            string        `mapstructure:"SERVER_PORT"`
	ReadTimeout     time.Duration `mapstructure:"SERVER_READ_TIMEOUT"`
	WriteTimeout    time.Duration `mapstructure:"SERVER_WRITE_TIMEOUT"`
	ShutdownTimeout time.Duration `mapstructure:"SERVER_SHUTDOWN_TIMEOUT"`
}

type DatabaseConfig struct {
	Host     string `mapstructure:"DB_HOST"`
	Port     string `mapstructure:"DB_PORT"`
	User     string `mapstructure:"DB_USER"`
	Password string `mapstructure:"DB_PASSWORD"`
	Name     string `mapstructure:"DB_NAME"`
	SSLMode  string `mapstructure:"DB_SSLMODE"`
	TimeZone string `mapstructure:"DB_TIMEZONE"`

	MaxOpenConns    int           `mapstructure:"DB_MAX_OPEN_CONNS"`
	MaxIdleConns    int           `mapstructure:"DB_MAX_IDLE_CONNS"`
	ConnMaxLifetime time.Duration `mapstructure:"DB_CONN_MAX_LIFETIME"`
}

type JWTConfig struct {
	AccessSecret  string        `mapstructure:"JWT_ACCESS_SECRET"`
	RefreshSecret string        `mapstructure:"JWT_REFRESH_SECRET"`
	AccessTTL     time.Duration `mapstructure:"JWT_ACCESS_TTL"`
	RefreshTTL    time.Duration `mapstructure:"JWT_REFRESH_TTL"`
	Issuer        string        `mapstructure:"JWT_ISSUER"`
}

type LogConfig struct {
	Level string `mapstructure:"LOG_LEVEL"`
}

// DSN returns a PostgreSQL connection string suitable for GORM/pgx.
func (d DatabaseConfig) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s TimeZone=%s",
		d.Host, d.Port, d.User, d.Password, d.Name, d.SSLMode, d.TimeZone,
	)
}

// Load reads configuration from environment variables. An optional .env file
// is loaded automatically when present. Defaults are applied for every key so
// the service can boot in a minimal environment.
func Load() (*Config, error) {
	v := viper.New()

	v.SetConfigFile(".env")
	v.SetConfigType("env")
	// Ignore a missing .env file — real environments inject vars directly.
	_ = v.ReadInConfig()

	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	setDefaults(v)

	cfg := &Config{
		App: AppConfig{
			Name: v.GetString("APP_NAME"),
			Env:  v.GetString("APP_ENV"),
		},
		Server: ServerConfig{
			Port:            v.GetString("SERVER_PORT"),
			ReadTimeout:     v.GetDuration("SERVER_READ_TIMEOUT"),
			WriteTimeout:    v.GetDuration("SERVER_WRITE_TIMEOUT"),
			ShutdownTimeout: v.GetDuration("SERVER_SHUTDOWN_TIMEOUT"),
		},
		Database: DatabaseConfig{
			Host:            v.GetString("DB_HOST"),
			Port:            v.GetString("DB_PORT"),
			User:            v.GetString("DB_USER"),
			Password:        v.GetString("DB_PASSWORD"),
			Name:            v.GetString("DB_NAME"),
			SSLMode:         v.GetString("DB_SSLMODE"),
			TimeZone:        v.GetString("DB_TIMEZONE"),
			MaxOpenConns:    v.GetInt("DB_MAX_OPEN_CONNS"),
			MaxIdleConns:    v.GetInt("DB_MAX_IDLE_CONNS"),
			ConnMaxLifetime: v.GetDuration("DB_CONN_MAX_LIFETIME"),
		},
		JWT: JWTConfig{
			AccessSecret:  v.GetString("JWT_ACCESS_SECRET"),
			RefreshSecret: v.GetString("JWT_REFRESH_SECRET"),
			AccessTTL:     v.GetDuration("JWT_ACCESS_TTL"),
			RefreshTTL:    v.GetDuration("JWT_REFRESH_TTL"),
			Issuer:        v.GetString("JWT_ISSUER"),
		},
		Log: LogConfig{
			Level: v.GetString("LOG_LEVEL"),
		},
	}

	if err := cfg.validate(); err != nil {
		return nil, err
	}
	return cfg, nil
}

func (c *Config) validate() error {
	if c.JWT.AccessSecret == "" || c.JWT.RefreshSecret == "" {
		return fmt.Errorf("config: JWT secrets must be set")
	}
	if c.Database.Host == "" || c.Database.Name == "" {
		return fmt.Errorf("config: database host and name must be set")
	}
	return nil
}

func (c *Config) IsProduction() bool {
	return strings.EqualFold(c.App.Env, "production")
}

func setDefaults(v *viper.Viper) {
	v.SetDefault("APP_NAME", "employee-management-api")
	v.SetDefault("APP_ENV", "development")

	v.SetDefault("SERVER_PORT", "8080")
	v.SetDefault("SERVER_READ_TIMEOUT", "15s")
	v.SetDefault("SERVER_WRITE_TIMEOUT", "15s")
	v.SetDefault("SERVER_SHUTDOWN_TIMEOUT", "10s")

	v.SetDefault("DB_HOST", "localhost")
	v.SetDefault("DB_PORT", "5432")
	v.SetDefault("DB_USER", "postgres")
	v.SetDefault("DB_PASSWORD", "postgres")
	v.SetDefault("DB_NAME", "ems")
	v.SetDefault("DB_SSLMODE", "disable")
	v.SetDefault("DB_TIMEZONE", "UTC")
	v.SetDefault("DB_MAX_OPEN_CONNS", 25)
	v.SetDefault("DB_MAX_IDLE_CONNS", 25)
	v.SetDefault("DB_CONN_MAX_LIFETIME", "5m")

	v.SetDefault("JWT_ACCESS_TTL", "15m")
	v.SetDefault("JWT_REFRESH_TTL", "168h") // 7 days
	v.SetDefault("JWT_ISSUER", "employee-management-api")

	v.SetDefault("LOG_LEVEL", "info")
}
