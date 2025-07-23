package supabaseconfig

import (
	"os"
	"text/template"
)

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
