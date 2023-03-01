package server

import (
	"testing"

	"github.com/UpCloudLtd/upcloud-cli/v2/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v2/internal/config"
	smock "github.com/UpCloudLtd/upcloud-cli/v2/internal/mock"
	"github.com/UpCloudLtd/upcloud-cli/v2/internal/mockexecute"

	"github.com/UpCloudLtd/upcloud-go-api/v6/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/v6/upcloud/request"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/stretchr/testify/assert"
)

func TestListServers(t *testing.T) {
	text.DisableColors()

	uuid := "server-list-test-server-uuid"
	servers := upcloud.Servers{
		Servers: []upcloud.Server{
			{
				Hostname: "server-list-test-server",
				UUID:     uuid,
				Plan:     "1xCPU-1GB",
				Zone:     "pl-waw1",
				State:    "started",
			},
		},
	}
	serverNetworks := upcloud.Networking{
		Interfaces: upcloud.ServerInterfaceSlice{
			{
				IPAddresses: upcloud.IPAddressSlice{
					{
						Access:   "utility",
						Address:  "10.0.100.1",
						Floating: upcloud.False,
					},
				},
				Type: "utility",
			},
			{
				IPAddresses: upcloud.IPAddressSlice{
					{
						Access:   "private",
						Address:  "10.0.99.2",
						Floating: upcloud.False,
					},
				},
				Type: "private",
			},
			{
				IPAddresses: upcloud.IPAddressSlice{
					{
						Access:   "public",
						Address:  "10.0.98.3",
						Floating: upcloud.False,
					},
					{
						Access:   "public",
						Address:  "10.0.97.4",
						Floating: upcloud.True,
					},
				},
				Type: "public",
			},
		},
	}

	ipaddressesTitle := "IP addresses"

	for _, test := range []struct {
		name              string
		args              []string
		outputContains    []string
		outputNotContains []string
	}{
		{
			name: "No args",
			args: []string{},
			outputNotContains: []string{
				ipaddressesTitle,
				"10.0.98.3",
			},
		},
		{
			name: "Show all IP Addresses",
			args: []string{"--show-ip-addresses"},
			outputContains: []string{
				ipaddressesTitle,
				"public: 10.0.97.4 (f)",
				"public: 10.0.98.3",
				"private: 10.0.99.2",
				"utility: 10.0.100.1",
			},
		},
		{
			name: "Show public Addresses",
			args: []string{"--show-ip-addresses=public"},
			outputContains: []string{
				ipaddressesTitle,
				"public: 10.0.97.4 (f)",
				"public: 10.0.98.3",
			},
			outputNotContains: []string{
				"private: 10.0.99.2",
				"utility: 10.0.100.1",
			},
		},
		{
			name: "Show private IP Addresses",
			args: []string{"--show-ip-addresses=private"},
			outputContains: []string{
				ipaddressesTitle,
				"private: 10.0.99.2",
			},
			outputNotContains: []string{
				"public: 10.0.97.4 (f)",
				"public: 10.0.98.3",
				"utility: 10.0.100.1",
			},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			conf := config.New()

			testCmd := ListCommand()
			mService := new(smock.Service)

			mService.On("GetServers").Return(&servers, nil)
			mService.On("GetServerNetworks", &request.GetServerNetworksRequest{ServerUUID: uuid}).Return(&serverNetworks, nil)

			c := commands.BuildCommand(testCmd, nil, conf)
			c.Cobra().SetArgs(test.args)

			output, err := mockexecute.MockExecute(c, mService, conf)
			assert.NoError(t, err)

			for _, contains := range test.outputContains {
				assert.Contains(t, output, contains)
			}

			for _, notContains := range test.outputNotContains {
				assert.NotContains(t, output, notContains)
			}
		})
	}
}
