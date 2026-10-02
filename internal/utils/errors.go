package utils

import (
	"errors"
	"net/http"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/apierror"
	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud"
)

func IsNotFoundError(err error) bool {
	if apierror.IsNotFound(err) {
		return true
	}

	var ucErr *upcloud.Problem
	if errors.As(err, &ucErr) && ucErr.Status == http.StatusNotFound {
		return true
	}

	return false
}
