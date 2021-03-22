package upapi

import (
	"fmt"
	"net/http"

	"github.com/UpCloudLtd/upcloud-go-api/upcloud/client"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/service"
	"github.com/hashicorp/go-cleanhttp"

	"github.com/UpCloudLtd/cli/internal/config"
	"github.com/UpCloudLtd/cli/internal/globals"
)

type transport struct{}

// RoundTrip implements http.RoundTripper
func (t *transport) RoundTrip(req *http.Request) (*http.Response, error) {
	return cleanhttp.DefaultTransport().RoundTrip(req)
}

// SetupService creates a new service instance and puts in the conf struct
func SetupService(config *config.Config) error {
	username := config.Top().GetString("username")
	password := config.Top().GetString("password")

	if username == "" || password == "" {
		err := `
User credentials not found, these must be set in config file or via environment vars
`
		return fmt.Errorf(err)
	}

	hc := &http.Client{Transport: &transport{}}
	hc.Timeout = config.ClientTimeout()

	whc := client.NewWithHTTPClient(
		username,
		password,
		hc,
	)
	whc.UserAgent = fmt.Sprintf("upctl/%s", globals.Version)

	config.Service = service.New(whc)

	return nil
}
