package config

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	AppEnv                 string
	Port                   string
	DatabaseURL            string
	CORSAllowOrigins       string
	BodyLimitMB            int
	RateLimitMax           int
	JWTSecret              string
	JWTExpiryHours         int
	CookieSecure           bool
	CookieSameSite         string
	FaceServiceURL         string
	InternalAPIKey         string
	FaceServiceTimeoutSec  int
	BiometricEncryptionKey string
}

// Load memuat konfigurasi dari file .env (jika ada) dan variabel lingkungan sistem
func Load() (*Config, error) {
	loadDotEnv(".env")

	appEnv := getEnv("APP_ENV", "development")
	cookieSecure := appEnv == "production"
	if cs := os.Getenv("COOKIE_SECURE"); cs != "" {
		cookieSecure = cs == "true" || cs == "1"
	}

	cfg := &Config{
		AppEnv:                 appEnv,
		Port:                   getEnv("PORT", getEnv("APP_PORT", "8000")),
		DatabaseURL:            os.Getenv("DATABASE_URL"),
		CORSAllowOrigins:       getEnv("CORS_ALLOW_ORIGINS", "http://localhost:8443,http://127.0.0.1:8443"),
		BodyLimitMB:            getEnvInt("BODY_LIMIT_MB", 10),
		RateLimitMax:           getEnvInt("RATE_LIMIT_MAX", 60),
		JWTSecret:              getEnv("JWT_SECRET", "labsync-dev-secret-key-change-in-production"),
		JWTExpiryHours:         getEnvInt("JWT_EXPIRY_HOURS", 24),
		CookieSecure:           cookieSecure,
		CookieSameSite:         getEnv("COOKIE_SAMESITE", "Lax"),
		FaceServiceURL:         getEnv("FACE_SERVICE_URL", "http://localhost:8001"),
		InternalAPIKey:         getEnv("INTERNAL_API_KEY", "labsync-internal-secret-key"),
		FaceServiceTimeoutSec:  getEnvInt("FACE_SERVICE_TIMEOUT_SEC", 10),
		BiometricEncryptionKey: getEnv("BIOMETRIC_ENCRYPTION_KEY", "01234567890123456789012345678901"),
	}

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

// Validate memeriksa validitas konfigurasi
func (c *Config) Validate() error {
	if c.Port == "" {
		return errors.New("PORT tidak boleh kosong")
	}
	if _, err := strconv.Atoi(c.Port); err != nil {
		return fmt.Errorf("PORT tidak valid: %s", c.Port)
	}
	if c.AppEnv == "production" && c.DatabaseURL == "" {
		return errors.New("DATABASE_URL wajib diisi pada environment production")
	}
	if c.AppEnv == "production" && (c.JWTSecret == "" || c.JWTSecret == "labsync-dev-secret-key-change-in-production") {
		return errors.New("JWT_SECRET wajib dikonfigurasi secara aman pada environment production")
	}
	if c.AppEnv == "production" && (len(c.BiometricEncryptionKey) < 32 || c.BiometricEncryptionKey == "01234567890123456789012345678901") {
		return errors.New("BIOMETRIC_ENCRYPTION_KEY minimal 32 karakter dan wajib diubah pada environment production")
	}
	if len(c.BiometricEncryptionKey) < 32 {
		return errors.New("BIOMETRIC_ENCRYPTION_KEY minimal 32 karakter (256-bit)")
	}
	if c.BodyLimitMB <= 0 {
		return errors.New("BODY_LIMIT_MB harus lebih besar dari 0")
	}
	if c.RateLimitMax <= 0 {
		return errors.New("RATE_LIMIT_MAX harus lebih besar dari 0")
	}
	if c.JWTExpiryHours <= 0 {
		return errors.New("JWT_EXPIRY_HOURS harus lebih besar dari 0")
	}
	if c.FaceServiceTimeoutSec <= 0 {
		return errors.New("FACE_SERVICE_TIMEOUT_SEC harus lebih besar dari 0")
	}
	return nil
}

// BodyLimitBytes mengembalikan batas ukuran body dalam satuan bytes
func (c *Config) BodyLimitBytes() int {
	return c.BodyLimitMB * 1024 * 1024
}

// ponytail: parser .env sederhana berbasis stdlib cukup sampai butuh multiline syntax
func loadDotEnv(filepath string) {
	file, err := os.Open(filepath)
	if err != nil {
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		val := strings.TrimSpace(parts[1])
		val = strings.Trim(val, `"'`)
		if _, exists := os.LookupEnv(key); !exists {
			_ = os.Setenv(key, val)
		}
	}
}

func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}

func getEnvInt(key string, defaultVal int) int {
	if val := os.Getenv(key); val != "" {
		if parsed, err := strconv.Atoi(val); err == nil {
			return parsed
		}
	}
	return defaultVal
}
