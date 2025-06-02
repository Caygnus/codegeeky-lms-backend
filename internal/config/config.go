package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/joho/godotenv"
	"github.com/omkar273/police/internal/types"
	"github.com/omkar273/police/internal/validator"
	"github.com/spf13/viper"
)

type Configuration struct {
	Server   ServerConfig   `validate:"required"`
	Logging  LoggingConfig  `validate:"required"`
	Postgres PostgresConfig `validate:"required"`
	Supabase SupabaseConfig `validate:"required"`
	Secrets  SecretsConfig  `validate:"required"`
}

type LoggingConfig struct {
	Level types.LogLevel `mapstructure:"level" validate:"required"`
}

type Env string

const (
	EnvLocal Env = "local"
	EnvDev   Env = "dev"
	EnvProd  Env = "prod"
)

type ServerConfig struct {
	Env     Env    `mapstructure:"env" validate:"required"`
	Address string `mapstructure:"address" validate:"required"`
}

type PostgresConfig struct {
	Host                   string `mapstructure:"host" validate:"required"`
	Port                   int    `mapstructure:"port" validate:"required"`
	User                   string `mapstructure:"user" validate:"required"`
	Password               string `mapstructure:"password" validate:"required"`
	DBName                 string `mapstructure:"dbname" validate:"required"`
	SSLMode                string `mapstructure:"sslmode" validate:"required"`
	MaxOpenConns           int    `mapstructure:"max_open_conns" default:"10"`
	MaxIdleConns           int    `mapstructure:"max_idle_conns" default:"5"`
	ConnMaxLifetimeMinutes int    `mapstructure:"conn_max_lifetime_minutes" default:"60"`
	AutoMigrate            bool   `mapstructure:"auto_migrate" default:"false"`
}

type SecretsConfig struct {
	PrivateKey string `mapstructure:"private_key" validate:"required"`
	PublicKey  string `mapstructure:"public_key" validate:"required"`
}

type SupabaseConfig struct {
	URL        string `mapstructure:"url" validate:"required"`
	Key        string `mapstructure:"key" validate:"required"`
	JWTSecret  string `mapstructure:"jwt_secret" validate:"required"`
	ServiceKey string `mapstructure:"service_key" validate:"required"`
}

func NewConfig() (*Configuration, error) {
	v := viper.New()

	// Step 1: Load `.env` if it exists
	envLoaded := godotenv.Load() == nil

	// Step 2: Initialize Viper
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath("./internal/config")
	v.AddConfigPath("./config")

	// Step 3: Set up environment variables support
	v.SetEnvPrefix("CAYGNUS")
	v.AutomaticEnv()

	// Step 4: Environment variable key mapping (e.g., CAYGNUS_SUPABASE_URL)
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Step 5: Read the YAML file
	configFileFound := true
	if err := v.ReadInConfig(); err != nil {
		if errors.As(err, &viper.ConfigFileNotFoundError{}) {
			configFileFound = false
			fmt.Printf("Warning: No config file found\n")
		} else {
			return nil, fmt.Errorf("error reading config file: %v", err)
		}
	} else {
		fmt.Printf("Using config file: %s\n", v.ConfigFileUsed())
	}

	// Check if we have any configuration source
	if !configFileFound && !envLoaded {
		fmt.Printf("Warning: Neither config.yaml nor .env file found. Checking environment variables...\n")
	}

	var cfg Configuration
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unable to decode into config struct, %v", err)
	}

	// Step 6: Load keys from files if not provided via config/env
	if err := cfg.loadKeysFromFiles(); err != nil {
		return nil, fmt.Errorf("error loading keys from files: %v", err)
	}

	// Step 7: Validate the configuration
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("configuration validation failed: %v\n\nPlease ensure you have either:\n1. A valid config.yaml file in ./internal/config/ or ./config/\n2. A .env file with required variables\n3. Environment variables with CAYGNUS_ prefix\n4. Key files in ./keys/ folder (private_key.pem, public_key.pem)\n\nRequired fields: server.address, logging.level, postgres.host, postgres.port, postgres.user, postgres.password, postgres.dbname, postgres.sslmode, supabase.url, supabase.key, supabase.jwt_secret, secrets.private_key, secrets.public_key", err)
	}

	return &cfg, nil
}

func (c Configuration) Validate() error {
	return validator.ValidateRequest(c)
}

// GetDefaultConfig returns a default configuration for local development
// This is useful for running scripts or other non-web applications
func GetDefaultConfig() *Configuration {
	return &Configuration{
		Server: ServerConfig{
			Address: ":8080",
		},
		Logging: LoggingConfig{
			Level: types.LogLevelDebug,
		},
	}
}

func (p PostgresConfig) GetDSN() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		p.Host,
		p.Port,
		p.User,
		p.Password,
		p.DBName,
		p.SSLMode,
	)
}

// loadKeysFromFiles attempts to load private and public keys from the keys folder
// if they are not already provided via configuration or environment variables
func (c *Configuration) loadKeysFromFiles() error {
	keysDir := "keys"

	// Load private key if not already set
	if c.Secrets.PrivateKey == "" {
		privateKeyPath := filepath.Join(keysDir, "private_key.pem")
		if content, err := os.ReadFile(privateKeyPath); err == nil {
			c.Secrets.PrivateKey = string(content)
			fmt.Printf("Loaded private key from: %s\n", privateKeyPath)
		} else {
			// Try alternative naming conventions
			altPaths := []string{
				filepath.Join(keysDir, "private.pem"),
				filepath.Join(keysDir, "id_rsa"),
				filepath.Join(keysDir, "private_key"),
			}

			for _, altPath := range altPaths {
				if content, err := os.ReadFile(altPath); err == nil {
					c.Secrets.PrivateKey = string(content)
					fmt.Printf("Loaded private key from: %s\n", altPath)
					break
				}
			}
		}
	}

	// Load public key if not already set
	if c.Secrets.PublicKey == "" {
		publicKeyPath := filepath.Join(keysDir, "public_key.pem")
		if content, err := os.ReadFile(publicKeyPath); err == nil {
			c.Secrets.PublicKey = string(content)
			fmt.Printf("Loaded public key from: %s\n", publicKeyPath)
		} else {
			// Try alternative naming conventions
			altPaths := []string{
				filepath.Join(keysDir, "public.pem"),
				filepath.Join(keysDir, "id_rsa.pub"),
				filepath.Join(keysDir, "public_key"),
			}

			for _, altPath := range altPaths {
				if content, err := os.ReadFile(altPath); err == nil {
					c.Secrets.PublicKey = string(content)
					fmt.Printf("Loaded public key from: %s\n", altPath)
					break
				}
			}
		}
	}

	return nil
}
