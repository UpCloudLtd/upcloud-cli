package clierrors

import "fmt"

var _ ClientError = MissingCredentialsError{}

type MissingCredentialsError struct {
	ConfigFile string
}

func (err MissingCredentialsError) ErrorCode() int {
	return MissingCredentials
}

func (err MissingCredentialsError) Error() string {
	return fmt.Sprintf("user credentials not found, these must be set in config file (%s) or via environment variables", err.ConfigFile)
}
