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

// Service creates a new service instance
func Service(config *config.Config) *service.Service {
	hc := &http.Client{Transport: &transport{}}

	whc := client.NewWithHTTPClient(
		config.Top().GetString("username"),
		config.Top().GetString("password"),
		hc)
	whc.UserAgent = fmt.Sprintf("upctl/%s", globals.Version)

	svc := service.New(whc)
	hc.Timeout = config.ClientTimeout()
	return svc
}
