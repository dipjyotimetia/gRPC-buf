package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLoadMergesYAMLAndEnvOverrides(t *testing.T) {
	t.Setenv("CFG_SERVER_PORT", "9100")
	t.Setenv("SERVER_LOG_LEVEL", "warn")
	t.Setenv("SERVER_LOGIN_RPS", "12")
	t.Setenv("SERVER_LOGIN_BURST", "24")
	t.Setenv("CFG_DATABASE_URL", "postgres://override")
	t.Setenv("CFG_DATABASE_MAX_CONNS", "25")
	t.Setenv("SECURITY_AUTH_SKIP_SUFFIXES", "/Foo,/Bar")
	t.Setenv("CFG_SECURITY_JWT_SECRET", "override-secret")

	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")
	data := []byte(`environment: prod
server:
  port: 8001
  run_migrations: false
  log_level: error
  login_rps: 7
  login_burst: 14
database:
  url: postgres://yaml
  max_conns: 10
  min_conns: 2
  connect_timeout: 30s
security:
  jwt_secret: yaml-secret
  auth_skip_suffixes: ["/YAMLOnly"]
`)
	require.NoError(t, os.WriteFile(path, data, 0o600))

	cfg, err := Load(path)
	require.NoError(t, err)

	require.Equal(t, "prod", cfg.Environment)
	require.Equal(t, 9100, cfg.Server.Port)
	require.False(t, cfg.Server.RunMigrations)
	require.Equal(t, "warn", cfg.Server.LogLevel)
	require.Equal(t, 12, cfg.Server.LoginRPS)
	require.Equal(t, 24, cfg.Server.LoginBurst)
	require.Equal(t, "postgres://override", cfg.Database.URL)
	require.Equal(t, 25, cfg.Database.MaxConns)
	require.Equal(t, 2, cfg.Database.MinConns)
	require.Equal(t, "30s", cfg.Database.ConnectTimeout)
	require.Equal(t, "override-secret", cfg.Security.JWTSecret)
	require.ElementsMatch(t, []string{"/Foo", "/Bar"}, cfg.Security.AuthSkipSuffixes)
}

func TestLoadExpandsEnvPlaceholders(t *testing.T) {
	t.Setenv("DATABASE_BASE", "postgres://interpolated")
	t.Setenv("SECRET_VALUE", "fromenv")

	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")
	data := []byte(`environment: dev
database:
  url: ${DATABASE_BASE}/grpcbuf
security:
  jwt_secret: ${SECRET_VALUE}
`)
	require.NoError(t, os.WriteFile(path, data, 0o600))

	cfg, err := Load(path)
	require.NoError(t, err)

	require.Equal(t, "dev", cfg.Environment)
	require.Equal(t, 8080, cfg.Server.Port)
	require.Equal(t, "postgres://interpolated/grpcbuf", cfg.Database.URL)
	require.Equal(t, "fromenv", cfg.Security.JWTSecret)
	require.ElementsMatch(t, []string{"/RegisterUser", "/LoginUser"}, cfg.Security.AuthSkipSuffixes)
}
