package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	AppEnv               string
	AppAddr              string
	AppBaseURL           string
	SessionSecret        string
	AzureTenantID        string
	AzureClientID        string
	AzureClientSecret    string
	AzureRedirectURL     string
	AzureAllowedGroupIDs []string
	DuckDBPath           string
}

func Load() (Config, error) {
	_ = godotenv.Load()

	cfg := Config{
		AppEnv:            env("APP_ENV", "development"),
		AppAddr:           env("APP_ADDR", "127.0.0.1:8080"),
		AppBaseURL:        env("APP_BASE_URL", "http://127.0.0.1:8080"),
		SessionSecret:     os.Getenv("SESSION_SECRET"),
		AzureTenantID:     os.Getenv("AZURE_TENANT_ID"),
		AzureClientID:     os.Getenv("AZURE_CLIENT_ID"),
		AzureClientSecret: os.Getenv("AZURE_CLIENT_SECRET"),
		AzureRedirectURL:  os.Getenv("AZURE_REDIRECT_URL"),
		DuckDBPath:        env("DUCKDB_PATH", "./mabata.duckdb"),
	}

	groups := strings.TrimSpace(os.Getenv("AZURE_ALLOWED_GROUP_IDS"))
	if groups != "" {
		for _, g := range strings.Split(groups, ",") {
			g = strings.TrimSpace(g)
			if g != "" {
				cfg.AzureAllowedGroupIDs = append(cfg.AzureAllowedGroupIDs, g)
			}
		}
	}

	missing := []string{}
	for k, v := range map[string]string{
		"SESSION_SECRET":      cfg.SessionSecret,
		"AZURE_TENANT_ID":     cfg.AzureTenantID,
		"AZURE_CLIENT_ID":     cfg.AzureClientID,
		"AZURE_CLIENT_SECRET": cfg.AzureClientSecret,
		"AZURE_REDIRECT_URL":  cfg.AzureRedirectURL,
	} {
		if strings.TrimSpace(v) == "" {
			missing = append(missing, k)
		}
	}
	if len(missing) > 0 {
		return Config{}, fmt.Errorf("missing required env vars: %s", strings.Join(missing, ", "))
	}

	return cfg, nil
}

func env(key, fallback string) string {
	if v := strings.TrimSpace(os.Getenv(key)); v != "" {
		return v
	}
	return fallback
}
