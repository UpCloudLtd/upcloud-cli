package clierrors

import (
	"errors"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/apierror"
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
	var errCode string
	if prob := (*apierror.Problem)(nil); errors.As(err, &prob) {
		errCode = prob.ErrorCode()
	} else if prob := (*upcloud.Problem)(nil); errors.As(err, &prob) {
		errCode = prob.ErrorCode()
	}

	return errCode == upcloud.ErrCodeAuthenticationFailed || errCode == "INVALID_CREDENTIALS"
}
