package clierrors

import (
	"errors"

	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud"
)

var _ ClientError = InvalidCredentialsError{}

type InvalidCredentialsError struct{}

func (err InvalidCredentialsError) ErrorCode() int {
	return InvalidCredentials
}

func (err InvalidCredentialsError) Error() string {
	return "invalid user credentials, authentication failed using the given username and password"
}

func CheckAuthenticationFailed(err error) bool {
	prob := &upcloud.Problem{}

	if errors.As(err, &prob) {
		errCode := prob.ErrorCode()
		if errCode == upcloud.ErrCodeAuthenticationFailed || errCode == "INVALID_CREDENTIALS" {
			return true
		}
	}

	return false
}
