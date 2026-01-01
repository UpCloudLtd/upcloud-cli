package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/adrg/xdg"
	"github.com/gemalto/flume"
	"gopkg.in/yaml.v3"
)

const (
	// CredentialsFileName is the standard credentials file name for UpCloud tools
	CredentialsFileName = "credentials"
	// CredentialsFileMode is the required permissions for credentials file (user read/write only)
	CredentialsFileMode = 0600
)

// GetCredentialsFilePath returns the standard path for UpCloud credentials file
func GetCredentialsFilePath() string {
	// Check XDG_CONFIG_HOME/upcloud/credentials first
	xdgPath := filepath.Join(xdg.ConfigHome, "upcloud", CredentialsFileName)

	// Fallback to ~/.upcloud/credentials for compatibility
	homePath := filepath.Join(os.Getenv("HOME"), ".upcloud", CredentialsFileName)

	// Use XDG path if directory exists or if neither exists (prefer XDG for new installs)
	if _, err := os.Stat(filepath.Dir(xdgPath)); err == nil {
		return xdgPath
	}
	if _, err := os.Stat(filepath.Dir(homePath)); err == nil {
		return homePath
	}

	// Default to XDG for new installations
	return xdgPath
}

// CredentialsFileData represents the structure of the credentials file
type CredentialsFileData struct {
	Token    string `yaml:"token,omitempty"`
	Username string `yaml:"username,omitempty"`
	Password string `yaml:"password,omitempty"`
}

// SaveTokenToCredentialsFile saves token to the standard credentials file with proper permissions
func SaveTokenToCredentialsFile(token string) (string, error) {
	credPath := GetCredentialsFilePath()
	credDir := filepath.Dir(credPath)

	// Create directory with secure permissions
	if err := os.MkdirAll(credDir, 0700); err != nil {
		return "", fmt.Errorf("failed to create credentials directory: %w", err)
	}

	// Prepare credentials data
	creds := CredentialsFileData{
		Token: token,
	}

	data, err := yaml.Marshal(&creds)
	if err != nil {
		return "", fmt.Errorf("failed to marshal credentials: %w", err)
	}

	// Write with secure permissions atomically
	tempFile := credPath + ".tmp"
	if err := os.WriteFile(tempFile, data, CredentialsFileMode); err != nil {
		return "", fmt.Errorf("failed to write credentials file: %w", err)
	}

	// Move atomically
	if err := os.Rename(tempFile, credPath); err != nil {
		os.Remove(tempFile) // Clean up
		return "", fmt.Errorf("failed to save credentials file: %w", err)
	}

	// Verify permissions (in case umask affected them)
	if err := os.Chmod(credPath, CredentialsFileMode); err != nil {
		// Log warning but don't fail (only if logger is available)
		if logger := flume.FromContext(nil); logger != nil {
			logger.Info("Could not set secure permissions on credentials file", "path", credPath)
		}
	}

	return credPath, nil
}

// LoadCredentialsFile loads credentials from the standard file if it exists
func LoadCredentialsFile() (CredentialsFileData, error) {
	credPath := GetCredentialsFilePath()

	// Check if file exists
	info, err := os.Stat(credPath)
	if err != nil {
		if os.IsNotExist(err) {
			return CredentialsFileData{}, nil // File doesn't exist, not an error
		}
		return CredentialsFileData{}, err
	}

	// Warn if permissions are too open (but only if logger is available)
	mode := info.Mode().Perm()
	if mode != CredentialsFileMode {
		// Try to get logger, but don't fail if it's not available
		if logger := flume.FromContext(nil); logger != nil {
			logger.Info("Credentials file has insecure permissions",
				"path", credPath,
				"current", fmt.Sprintf("%04o", mode),
				"expected", fmt.Sprintf("%04o", CredentialsFileMode))
		}
	}

	// Read file
	data, err := os.ReadFile(credPath)
	if err != nil {
		return CredentialsFileData{}, fmt.Errorf("failed to read credentials file: %w", err)
	}

	var creds CredentialsFileData
	if err := yaml.Unmarshal(data, &creds); err != nil {
		return CredentialsFileData{}, fmt.Errorf("failed to parse credentials file: %w", err)
	}

	return creds, nil
}

// IsKeyringError checks if an error is related to keyring access
func IsKeyringError(err error) bool {
	if err == nil {
		return false
	}
	errStr := err.Error()
	keyringIndicators := []string{
		"keyring", "secret", "dbus", "collection",
		"unlock", "gnome-keyring", "kwallet",
		"windows credential", "failed to unlock",
	}
	errLower := strings.ToLower(errStr)
	for _, indicator := range keyringIndicators {
		if strings.Contains(errLower, indicator) {
			return true
		}
	}
	return false
}