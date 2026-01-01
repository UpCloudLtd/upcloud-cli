package clierrors

import (
	"strings"
)

var _ ClientError = MissingCredentialsError{}

type MissingCredentialsError struct {
	ConfigFile   string
	ServiceName  string
	KeyringError error // New field to track keyring issues
}

func (err MissingCredentialsError) ErrorCode() int {
	return MissingCredentials
}

func (err MissingCredentialsError) Error() string {
	var msg strings.Builder

	msg.WriteString("Authentication credentials not found.\n\n")

	// Check if this is due to keyring issues
	if err.KeyringError != nil && isKeyringError(err.KeyringError) {
		msg.WriteString("System keyring is not accessible. ")
		msg.WriteString("You can:\n")
		msg.WriteString("  1. Save token to file: upctl account login --with-token --save-to-file\n")
		msg.WriteString("  2. Set environment variable: export UPCLOUD_TOKEN=your-token\n")
		msg.WriteString("  3. Add to config file: " + err.ConfigFile + "\n\n")
		msg.WriteString("Original error: " + err.KeyringError.Error())
	} else {
		msg.WriteString("Please configure authentication using one of these methods:\n")
		msg.WriteString("  1. Login with token: upctl account login --with-token\n")
		msg.WriteString("  2. Set environment variable: export UPCLOUD_TOKEN=your-token\n")
		msg.WriteString("  3. Add to config file: " + err.ConfigFile + "\n")

		if err.ServiceName != "" {
			msg.WriteString("\nNote: System keyring service '" + err.ServiceName + "' may not be available.\n")
			msg.WriteString("Use --save-to-file flag with login command if keyring fails.")
		}
	}

	return msg.String()
}

func isKeyringError(err error) bool {
	if err == nil {
		return false
	}
	errStr := strings.ToLower(err.Error())
	keyringIndicators := []string{
		"keyring", "secret", "dbus", "collection",
		"unlock", "gnome-keyring", "kwallet",
	}
	for _, indicator := range keyringIndicators {
		if strings.Contains(errStr, indicator) {
			return true
		}
	}
	return false
}
