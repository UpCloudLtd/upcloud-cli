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

func (t *transport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Add("User-Agent", fmt.Sprintf("upctl/%s", globals.Version))
	return cleanhttp.DefaultTransport().RoundTrip(req)
}

func Service(config *config.Config) *service.Service {
	hc := &http.Client{Transport: &transport{}}

	svc := service.New(client.NewWithHTTPClient(
		config.Top().GetString("username"),
		config.Top().GetString("password"),
		hc))
	hc.Timeout = config.ClientTimeout()
	return svc
}
