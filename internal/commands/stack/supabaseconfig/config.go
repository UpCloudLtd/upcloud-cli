package supabaseconfig

import (
	"bufio"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type SupabaseConfig struct {
	JWTSecret         string
	AnonKey           string
	ServiceRoleKey    string
	PostgresPassword  string
	PoolerTenantID    string
	DashboardUsername string
	DashboardPassword string

	S3Enabled    bool
	S3KeyID      string
	S3AccessKey  string
	S3BucketName string
	S3Region     string
	S3Endpoint   string

	SmtpEnabled    bool
	SmtpHost       string
	SmtpPort       string
	SmtpUsername   string
	SmtpPassword   string
	SmtpSenderName string
}

func Generate(configPath string) (*SupabaseConfig, error) {
	// Generate a random JWT secret and API keys
	jwtSecret, err := generateJWTSecret()
	if err != nil {
		return nil, err
	}

	// Use fixed timestamps for comparison. You can replace with time.Now().Unix() in production.
	iat := int64(1752786000)
	exp := int64(1910552400)

	anonKey, err := signJWT("anon", jwtSecret, iat, exp)
	if err != nil {
		return nil, err
	}

	serviceKey, err := signJWT("service_role", jwtSecret, iat, exp)
	if err != nil {
		return nil, err
	}

	// Read configuration from the provided configPath
	config := &SupabaseConfig{
		JWTSecret:         jwtSecret,
		AnonKey:           anonKey,
		ServiceRoleKey:    serviceKey,
		DashboardUsername: "supabase",
		DashboardPassword: generateRandomString(20),
		PostgresPassword:  generateRandomString(20),
		PoolerTenantID:    generateRandomString(20),
		S3Enabled:         false,
		SmtpEnabled:       false,
	}

	if configPath != "" {
		config, err = loadConfigFromFile(configPath)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to load configuration from file: %w", err)
	}

	return config, nil
}

func generateJWTSecret() (string, error) {
	bytes := make([]byte, 20) // 20 bytes = 40 hex characters
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func generateRandomString(length int) string {
	bytes := make([]byte, length)
	_, _ = rand.Read(bytes)
	return base64.RawURLEncoding.EncodeToString(bytes)[:length]
}

func signJWT(role string, jwtSecret string, iat int64, exp int64) (string, error) {
	type jwtHeader struct {
		Alg string `json:"alg"`
		Typ string `json:"typ"`
	}

	type jwtPayload struct {
		Role string `json:"role"`
		Iss  string `json:"iss"`
		Iat  int64  `json:"iat"`
		Exp  int64  `json:"exp"`
	}

	header := jwtHeader{
		Alg: "HS256",
		Typ: "JWT",
	}

	payload := jwtPayload{
		Role: role,
		Iss:  "supabase",
		Iat:  iat,
		Exp:  exp,
	}

	hb, err := json.Marshal(header)
	if err != nil {
		return "", err
	}
	pb, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	h64 := base64.RawURLEncoding.EncodeToString(hb)
	p64 := base64.RawURLEncoding.EncodeToString(pb)
	toSign := fmt.Sprintf("%s.%s", h64, p64)

	h := hmac.New(sha256.New, []byte(jwtSecret))
	h.Write([]byte(toSign))
	sig := base64.RawURLEncoding.EncodeToString(h.Sum(nil))

	return fmt.Sprintf("%s.%s", toSign, sig), nil
}

func loadConfigFromFile(path string) (*SupabaseConfig, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	config := &SupabaseConfig{}
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()

		// Skip empty lines and comments
		if strings.TrimSpace(line) == "" || strings.HasPrefix(line, "#") {
			continue
		}

		kv := strings.SplitN(line, "=", 2)
		if len(kv) != 2 {
			continue // Or return error if format must be strict
		}

		key := strings.TrimSpace(kv[0])
		val := strings.TrimSpace(kv[1])

		switch key {
		case "DASHBOARD_USERNAME":
			config.DashboardUsername = val
		case "DASHBOARD_PASSWORD":
			config.DashboardPassword = val
		case "POSTGRES_PASSWORD":
			config.PostgresPassword = val
		case "ENABLE_S3":
			config.S3Enabled = strings.ToLower(val) == "true"
		case "S3_KEY_ID":
			config.S3KeyID = val
		case "S3_ACCESS_KEY":
			config.S3AccessKey = val
		case "S3_BUCKET":
			config.S3BucketName = val
		case "S3_ENDPOINT":
			config.S3Endpoint = val
		case "S3_REGION":
			config.S3Region = val
		case "ENABLE_SMTP":
			config.SmtpEnabled = strings.ToLower(val) == "true"
		case "SMTP_HOST":
			config.SmtpHost = val
		case "SMTP_PORT":
			config.SmtpPort = val
		case "SMTP_USER":
			config.SmtpUsername = val
		case "SMTP_PASS":
			config.SmtpPassword = val
		case "SMTP_SENDER_NAME":
			config.SmtpSenderName = val
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return config, nil
}
