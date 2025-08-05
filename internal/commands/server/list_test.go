package server

import (
	"testing"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/config"
	smock "github.com/UpCloudLtd/upcloud-cli/v3/internal/mock"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/mockexecute"

	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud/request"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/stretchr/testify/assert"
)

const expectedJSONOutput = `
{
  "servers": [
    {
      "core_number": "0",
      "hostname": "server-list-test-server",
      "license": 0,
      "memory_amount": "0",
      "plan": "1xCPU-1GB",
      "progress": "0",
      "state": "started",
      "tags": null,
      "title": "",
      "uuid": "server-list-test-server-uuid",
      "zone": "pl-waw1",
      "ip_addresses": [
        {
          "access": "public",
          "address": "10.0.97.4",
		  "dhcp_provided": "no",
          "family": "",
          "part_of_plan": "no",
          "ptr_record": "",
          "server": "",
          "mac": "",
          "floating": "yes",
          "zone": ""
        },
        {
          "access": "public",
          "address": "10.0.98.3",
		  "dhcp_provided": "no",
          "family": "",
          "part_of_plan": "no",
          "ptr_record": "",
          "server": "",
          "mac": "",
          "floating": "no",
          "zone": ""
        },
        {
          "access": "private",
          "address": "10.0.99.2",
		  "dhcp_provided": "no",
          "family": "",
          "part_of_plan": "no",
          "ptr_record": "",
          "server": "",
          "mac": "",
          "floating": "no",
          "zone": ""
        },
        {
          "access": "utility",
          "address": "10.0.100.1",
		  "dhcp_provided": "no",
          "family": "",
          "part_of_plan": "no",
          "ptr_record": "",
          "server": "",
          "mac": "",
          "floating": "no",
          "zone": ""
        }
      ]
    }
  ]
}
`

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
		json              bool
		outputContains    []string
		outputNotContains []string
		outputJSONEquals  string
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
		{
			name:             "JSON output",
			args:             []string{"--show-ip-addresses"},
			json:             true,
			outputJSONEquals: expectedJSONOutput,
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			conf := config.New()
			if test.json {
				conf.Viper().Set(config.KeyOutput, config.ValueOutputJSON)
			}

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

			if len(test.outputJSONEquals) > 0 {
				assert.JSONEq(t, test.outputJSONEquals, output)
			}
		})
	}
}
