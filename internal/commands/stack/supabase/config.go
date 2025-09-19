package supabase

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"html/template"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type SupabaseConfig struct {
	LbHostname  string
	ClusterName string

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
		err = loadConfigFromFile(configPath, config)
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

func setIfNotEmpty(val string, setter func(string)) {
	if strings.TrimSpace(val) != "" {
		setter(val)
	}
}

func loadConfigFromFile(path string, config *SupabaseConfig) error {
	envMap, err := godotenv.Read(path)
	if err != nil {
		return err
	}

	// String fields
	setIfNotEmpty(envMap["DASHBOARD_USERNAME"], func(v string) { config.DashboardUsername = v })
	setIfNotEmpty(envMap["DASHBOARD_PASSWORD"], func(v string) { config.DashboardPassword = v })
	setIfNotEmpty(envMap["POSTGRES_PASSWORD"], func(v string) { config.PostgresPassword = v })
	setIfNotEmpty(envMap["S3_KEY_ID"], func(v string) { config.S3KeyID = v })
	setIfNotEmpty(envMap["S3_ACCESS_KEY"], func(v string) { config.S3AccessKey = v })
	setIfNotEmpty(envMap["S3_BUCKET"], func(v string) { config.S3BucketName = v })
	setIfNotEmpty(envMap["S3_ENDPOINT"], func(v string) { config.S3Endpoint = v })
	setIfNotEmpty(envMap["S3_REGION"], func(v string) { config.S3Region = v })
	setIfNotEmpty(envMap["SMTP_HOST"], func(v string) { config.SmtpHost = v })
	setIfNotEmpty(envMap["SMTP_PORT"], func(v string) { config.SmtpPort = v })
	setIfNotEmpty(envMap["SMTP_USER"], func(v string) { config.SmtpUsername = v })
	setIfNotEmpty(envMap["SMTP_PASS"], func(v string) { config.SmtpPassword = v })
	setIfNotEmpty(envMap["SMTP_SENDER_NAME"], func(v string) { config.SmtpSenderName = v })

	// Boolean flags — convert only if the value is not empty
	if v := strings.TrimSpace(envMap["ENABLE_S3"]); v != "" {
		config.S3Enabled = strings.ToLower(v) == "true"
	}
	if v := strings.TrimSpace(envMap["ENABLE_SMTP"]); v != "" {
		config.SmtpEnabled = strings.ToLower(v) == "true"
	}

	return nil
}

const valuesTemplate = `# values.secure.yaml — overrides for rotating Supabase secrets
secret:
  db:
    password: "{{.PostgresPassword}}"

  jwt:
    secret:      "{{.JWTSecret}}"
    anonKey:     "{{.AnonKey}}"
    serviceKey:  "{{.ServiceRoleKey}}"
    secretRef:   ""
    secretRefKey:
      anonKey:    anonKey
      serviceKey: serviceKey
      secret:     secret

  dashboard:
    username: "{{.DashboardUsername}}"
    password: "{{.DashboardPassword}}"

  s3:
    keyId: "{{.S3KeyID}}"
    accessKey: "{{.S3AccessKey}}"
    secretRef: ""
    secretRefKey:
      keyId: keyId
      accessKey: accessKey

    # SMTP configuration (if your chart references these)
    smtp:
        host:   "{{.SmtpHost}}"
        port:   "{{.SmtpPort}}"
        user:   "{{.SmtpUsername}}"
        pass:   "{{.SmtpPassword}}"
        sender: "{{.SmtpSenderName}}"

# POOLER (Supavisor) tenant ID
pooler:
  tenantId: "{{.PoolerTenantID}}"

storage:
  enabled: "{{.S3Enabled}}"               
  environment:
    STORAGE_BACKEND:            "s3"
    GLOBAL_S3_BUCKET:           "{{.S3BucketName}}"
    TENANT_ID:                  "supabase"
    GLOBAL_S3_ENDPOINT:         "{{.S3Endpoint}}"  
    GLOBAL_S3_PROTOCOL:         "https"  
    GLOBAL_S3_FORCE_PATH_STYLE: "true"  
    AWS_DEFAULT_REGION:         "{{.S3Region}}" 
`

// WriteSecureValues writes the secure values to a specified file using the provided SupabaseSecrets.
func WriteSecureValues(filePath string, secrets *SupabaseConfig) error {
	f, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer f.Close()

	tmpl, err := template.New("secure").Parse(valuesTemplate)
	if err != nil {
		return err
	}

	return tmpl.Execute(f, secrets)
}
