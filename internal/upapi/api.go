package upapi

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/UpCloudLtd/upcloud-go-api/upcloud/client"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/service"
	"github.com/hashicorp/go-cleanhttp"
	"github.com/spf13/viper"

	"github.com/UpCloudLtd/cli/internal/config"
)

var APIClient *service.Service

type transport struct{}

// RoundTrip implements http.RoundTripper
func (t *transport) RoundTrip(req *http.Request) (*http.Response, error) {
	return cleanhttp.DefaultTransport().RoundTrip(req)
}

// NewServiceClient creates a new service instance
func NewServiceClient(conf *config.Config, v *viper.Viper) error {
	var username, password string

	hc := &http.Client{Transport: &transport{}}
	hc.Timeout = conf.ClientTimeout()

	if username = v.GetString("USERNAME"); username == "" {
		return errors.New("Missing \"username\" for authentification,, it should be set in config file or via UPCLOUD_USERNAME environment variable\n")
	}

	if password = v.GetString("PASSWORD"); password == "" {
		return errors.New("Missing \"password\" for authentification, it should be set in config file or via UPCLOUD_PASSWORD environment variable\n")
	}

	fmt.Printf("2user: %v\n", username)
	fmt.Printf("2pass: %v\n", password)
	whc := client.NewWithHTTPClient(
		username,
		password,
		hc)
	whc.UserAgent = fmt.Sprintf("upctl/%s", config.Version)

	conf.Service = service.New(whc)

	return nil
}
