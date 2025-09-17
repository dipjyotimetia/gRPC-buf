package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/kelseyhightower/envconfig"
	"gopkg.in/yaml.v3"
)

type ServerConfig struct {
	Port               int      `yaml:"port" envconfig:"PORT"`
	CORSAllowedOrigins []string `yaml:"cors_allowed_origins" envconfig:"CORS_ALLOWED_ORIGINS"`
	RunMigrations      bool     `yaml:"run_migrations" envconfig:"RUN_MIGRATIONS"`
	LogLevel           string   `yaml:"log_level" envconfig:"LOG_LEVEL"`
	LoginRPS           int      `yaml:"login_rps" envconfig:"LOGIN_RPS"`
	LoginBurst         int      `yaml:"login_burst" envconfig:"LOGIN_BURST"`
}

type DatabaseConfig struct {
	URL            string `yaml:"url" envconfig:"URL"`
	MaxConns       int    `yaml:"max_conns" envconfig:"MAX_CONNS"`
	MinConns       int    `yaml:"min_conns" envconfig:"MIN_CONNS"`
	ConnectTimeout string `yaml:"connect_timeout" envconfig:"CONNECT_TIMEOUT"`
}

type SecurityConfig struct {
	JWTSecret        string   `yaml:"jwt_secret" envconfig:"JWT_SECRET"`
	JWTIssuer        string   `yaml:"jwt_issuer" envconfig:"JWT_ISSUER"`
	JWTAudience      string   `yaml:"jwt_audience" envconfig:"JWT_AUDIENCE"`
	AuthSkipSuffixes []string `yaml:"auth_skip_suffixes" envconfig:"AUTH_SKIP_SUFFIXES"`
}

type Config struct {
	Environment string         `yaml:"environment" envconfig:"ENVIRONMENT"`
	Server      ServerConfig   `yaml:"server" envconfig:"SERVER"`
	Database    DatabaseConfig `yaml:"database" envconfig:"DATABASE"`
	Security    SecurityConfig `yaml:"security" envconfig:"SECURITY"`
}

// Load hydrates configuration from an optional YAML file and environment variables.
// Env overrides use either direct names (e.g. SERVER_PORT) or the CFG_ prefix
// (e.g. CFG_SERVER_PORT). Environment values take precedence over file contents.
func Load(path string) (*Config, error) {
	cfg := Config{
		Environment: "dev",
		Server: ServerConfig{
			Port:          8080,
			RunMigrations: true,
			LogLevel:      "info",
			LoginRPS:      5,
			LoginBurst:    10,
		},
		Database: DatabaseConfig{
			ConnectTimeout: "60s",
		},
		Security: SecurityConfig{
			AuthSkipSuffixes: []string{"/RegisterUser", "/LoginUser"},
		},
	}
	if strings.TrimSpace(path) != "" {
		data, err := os.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("read config file: %w", err)
		}
		if err := yaml.Unmarshal(data, &cfg); err != nil {
			return nil, fmt.Errorf("unmarshal config file: %w", err)
		}
	}
	// Apply defaults & environment overrides via envconfig; allow either CFG_* or direct names.
	if err := envconfig.Process("", &cfg); err != nil {
		return nil, fmt.Errorf("load env config: %w", err)
	}
	if err := envconfig.Process("CFG", &cfg); err != nil {
		return nil, fmt.Errorf("load cfg env overrides: %w", err)
	}
	cfg.Database.URL = os.ExpandEnv(cfg.Database.URL)
	cfg.Security.JWTSecret = os.ExpandEnv(cfg.Security.JWTSecret)
	return &cfg, nil
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
	if strings.TrimSpace(c.Security.JWTSecret) == "" {
		if s := strings.TrimSpace(os.Getenv("JWT_SECRET")); s != "" {
			c.Security.JWTSecret = s
		}
	}
	if strings.EqualFold(c.Environment, "production") {
		if strings.TrimSpace(c.Database.URL) == "" {
			return fmt.Errorf("database.url is required in production")
		}
		if strings.TrimSpace(c.Security.JWTSecret) == "" {
			return fmt.Errorf("security.jwt_secret is required in production")
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

// (env-only export helper was removed as unused)
