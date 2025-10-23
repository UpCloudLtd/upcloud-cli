package utils

import (
	"errors"
	"net/http"

	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud"
)

func IsNotFoundError(err error) bool {
	var ucErr *upcloud.Problem
	if errors.As(err, &ucErr) && ucErr.Status == http.StatusNotFound {
		return true
	}

	return false
}
