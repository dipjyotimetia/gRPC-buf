package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	kfn "github.com/knadh/koanf/v2"
)

type ServerConfig struct {
    Port               int      `yaml:"port"`
    CORSAllowedOrigins []string `yaml:"cors_allowed_origins"`
    RunMigrations      bool     `yaml:"run_migrations"`
    LogLevel           string   `yaml:"log_level"`
    LoginRPS           int      `yaml:"login_rps"`
    LoginBurst         int      `yaml:"login_burst"`
}

type DatabaseConfig struct {
    URL      string `yaml:"url"`
    MaxConns int    `yaml:"max_conns"`
    MinConns int    `yaml:"min_conns"`
    ConnectTimeout string `yaml:"connect_timeout"`
}

type SecurityConfig struct {
    JWTSecret string `yaml:"jwt_secret"`
    JWTIssuer string `yaml:"jwt_issuer"`
    JWTAudience string `yaml:"jwt_audience"`
    AuthSkipSuffixes []string `yaml:"auth_skip_suffixes"`
}

type Config struct {
    Environment string         `yaml:"environment"`
    Server      ServerConfig   `yaml:"server"`
    Database    DatabaseConfig `yaml:"database"`
    Security    SecurityConfig `yaml:"security"`
}

// Load reads the YAML config file at path and applies env overrides via koanf.
// Env overrides use the prefix CFG_ and double underscore (__) to denote nesting.
// Example: CFG_SERVER__PORT=9090 overrides server.port.
func Load(path string) (*Config, error) {
    k := kfn.New(".")
    if err := k.Load(file.Provider(path), yaml.Parser()); err != nil {
        return nil, fmt.Errorf("load config file: %w", err)
    }
	// Env overrides with prefix CFG_. Replace __ with . for nested keys.
	if err := k.Load(env.Provider("CFG_", ".", func(s string) string {
		s = strings.TrimPrefix(s, "CFG_")
		s = strings.ReplaceAll(s, "__", ".")
		return strings.ToLower(s)
	}), nil); err != nil {
		return nil, fmt.Errorf("load env overrides: %w", err)
	}

    var c Config
    if err := k.Unmarshal("", &c); err != nil {
        return nil, fmt.Errorf("unmarshal config: %w", err)
    }
    // Best-effort env interpolation for secret fields
    c.Database.URL = os.ExpandEnv(c.Database.URL)
    c.Security.JWTSecret = os.ExpandEnv(c.Security.JWTSecret)
    return &c, nil
}

// Validate checks the configuration and returns an error if required fields are missing
// or invalid for the current environment.
func (c *Config) Validate() error {
	if c == nil {
		return errors.New("nil config")
	}
	env := strings.ToLower(strings.TrimSpace(c.Environment))
	if env == "" {
		env = "dev"
		c.Environment = env
	}
	if c.Server.Port <= 0 || c.Server.Port > 65535 {
		return fmt.Errorf("invalid server.port: %d", c.Server.Port)
	}
	if env == "prod" || env == "production" {
		if strings.TrimSpace(c.Database.URL) == "" {
			return errors.New("database.url is required in production")
		}
		if strings.TrimSpace(c.Security.JWTSecret) == "" {
			return errors.New("security.jwt_secret is required in production")
		}
	}
	return nil
}

// ResolvePath determines a config file path based on ENVIRONMENT if CONFIG_PATH is not provided.
func ResolvePath() (string, error) {
	if cp := os.Getenv("CONFIG_PATH"); strings.TrimSpace(cp) != "" {
		return cp, nil
	}
	env := strings.ToLower(strings.TrimSpace(os.Getenv("ENVIRONMENT")))
	switch env {
	case "dev", "development", "local", "":
		return filepath.Join("config", "local.yaml"), nil
	case "prod", "production":
		return filepath.Join("config", "production.yaml"), nil
	default:
		return "", fmt.Errorf("unknown ENVIRONMENT %q; set CONFIG_PATH explicitly", env)
	}
}

// ApplyEnv exports configuration as environment variables, so existing code can consume it.
func (c *Config) ApplyEnv() error {
	if c == nil {
		return errors.New("nil config")
	}
	set := func(k, v string) { _ = os.Setenv(k, v) }

	set("ENVIRONMENT", c.Environment)
	if c.Server.Port != 0 {
		set("PORT", fmt.Sprintf("%d", c.Server.Port))
	}
	if len(c.Server.CORSAllowedOrigins) > 0 {
		set("CORS_ALLOWED_ORIGINS", strings.Join(c.Server.CORSAllowedOrigins, ","))
	}
	if c.Server.RunMigrations {
		set("RUN_MIGRATIONS", "true")
	} else {
		set("RUN_MIGRATIONS", "false")
	}

	if c.Database.URL != "" {
		set("DATABASE_URL", c.Database.URL)
	}
	if c.Database.MaxConns > 0 {
		set("DB_MAX_CONNS", fmt.Sprintf("%d", c.Database.MaxConns))
	}
	if c.Database.MinConns >= 0 {
		set("DB_MIN_CONNS", fmt.Sprintf("%d", c.Database.MinConns))
	}

    if c.Security.JWTSecret != "" {
        set("JWT_SECRET", c.Security.JWTSecret)
    }
    if c.Security.JWTIssuer != "" {
        set("JWT_ISSUER", c.Security.JWTIssuer)
    }
    if c.Security.JWTAudience != "" {
        set("JWT_AUDIENCE", c.Security.JWTAudience)
    }
    if len(c.Security.AuthSkipSuffixes) > 0 {
        set("AUTH_SKIP_SUFFIXES", strings.Join(c.Security.AuthSkipSuffixes, ","))
    }
    if c.Server.LogLevel != "" {
        set("LOG_LEVEL", strings.ToLower(c.Server.LogLevel))
    }
    return nil
}
