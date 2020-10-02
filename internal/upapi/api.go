package upapi

import (
	"fmt"
	"net/http"

	"github.com/UpCloudLtd/upcloud-go-api/upcloud/client"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/service"
	"github.com/hashicorp/go-cleanhttp"
	"github.com/spf13/viper"

	"github.com/UpCloudLtd/cli/internal/globals"
)

type transport struct{}

func (t *transport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Add("User-Agent", fmt.Sprintf("upctl/%s", globals.Version))
	return cleanhttp.DefaultTransport().RoundTrip(req)
}

func Service(config *viper.Viper) *service.Service {
	return service.New(client.NewWithHTTPClient(
		config.GetString("username"),
		config.GetString("password"),
		&http.Client{Transport: &transport{}},
	))
}
