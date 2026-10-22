package config

import (
	"os"
	"testing"
)

func TestConfigValidate(t *testing.T) {
	tests := []struct {
		name    string
		cfg     Config
		wantErr bool
	}{
		{
			name: "Valid default development config",
			cfg: Config{
				AppEnv:                 "development",
				Port:                   "8000",
				CORSAllowOrigins:       "http://localhost:8443",
				BodyLimitMB:            10,
				RateLimitMax:           60,
				JWTExpiryHours:         24,
				JWTSecret:              "dev-secret",
				BiometricEncryptionKey: "01234567890123456789012345678901",
				FaceServiceTimeoutSec:  10,
			},
			wantErr: false,
		},
		{
			name: "Invalid port string",
			cfg: Config{
				AppEnv:         "development",
				Port:           "abc",
				BodyLimitMB:    10,
				RateLimitMax:   60,
				JWTExpiryHours: 24,
			},
			wantErr: true,
		},
		{
			name: "Missing DatabaseURL in production",
			cfg: Config{
				AppEnv:         "production",
				Port:           "8000",
				DatabaseURL:    "",
				BodyLimitMB:    10,
				RateLimitMax:   60,
				JWTExpiryHours: 24,
				JWTSecret:      "production-secret-1234567890",
			},
			wantErr: true,
		},
		{
			name: "Default/weak JWTSecret in production",
			cfg: Config{
				AppEnv:         "production",
				Port:           "8000",
				DatabaseURL:    "postgresql://localhost/db",
				BodyLimitMB:    10,
				RateLimitMax:   60,
				JWTExpiryHours: 24,
				JWTSecret:      "labsync-dev-secret-key-change-in-production",
			},
			wantErr: true,
		},
		{
			name: "Invalid BodyLimitMB",
			cfg: Config{
				AppEnv:         "development",
				Port:           "8000",
				BodyLimitMB:    0,
				RateLimitMax:   60,
				JWTExpiryHours: 24,
			},
			wantErr: true,
		},
		{
			name: "Invalid RateLimitMax",
			cfg: Config{
				AppEnv:         "development",
				Port:           "8000",
				BodyLimitMB:    10,
				RateLimitMax:   0,
				JWTExpiryHours: 24,
			},
			wantErr: true,
		},
		{
			name: "Invalid JWTExpiryHours",
			cfg: Config{
				AppEnv:         "development",
				Port:           "8000",
				BodyLimitMB:    10,
				RateLimitMax:   60,
				JWTExpiryHours: 0,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.cfg.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestConfigLoadWithEnv(t *testing.T) {
	_ = os.Setenv("PORT", "9090")
	_ = os.Setenv("BODY_LIMIT_MB", "20")
	defer func() {
		_ = os.Unsetenv("PORT")
		_ = os.Unsetenv("BODY_LIMIT_MB")
	}()

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() unexpected error: %v", err)
	}

	if cfg.Port != "9090" {
		t.Errorf("expected Port 9090, got %s", cfg.Port)
	}
	if cfg.BodyLimitMB != 20 {
		t.Errorf("expected BodyLimitMB 20, got %d", cfg.BodyLimitMB)
	}
	if cfg.BodyLimitBytes() != 20*1024*1024 {
		t.Errorf("expected BodyLimitBytes 20MB, got %d", cfg.BodyLimitBytes())
	}
}
