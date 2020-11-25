package network

import (
	"github.com/UpCloudLtd/cli/internal/commands"
	"github.com/UpCloudLtd/cli/internal/config"
	"github.com/UpCloudLtd/cli/internal/mocks"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/request"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreateCommand(t *testing.T) {

	for _, test := range []struct {
		name     string
		args     []string
		expected request.CreateNetworkRequest
	}{
		{
			name: "without network",
			args: []string{"--name", "asdf", "--zone", "fi-hel1"},
			expected: request.CreateNetworkRequest{
				Name: "asdf",
				Zone: "fi-hel1",
			},
		},
		{
			name: "with single network",
			args: []string{"--ip-network", "family=IPv4,\"dns=one,two,three\",gateway=gw,dhcp=true"},
			expected: request.CreateNetworkRequest{
				IPNetworks: []upcloud.IPNetwork{
					{
						Family:           "IPv4",
						DHCP:             upcloud.FromBool(true),
						DHCPDefaultRoute: upcloud.FromBool(false),
						DHCPDns:          []string{"one", "two", "three"},
						Gateway:          "gw",
					},
				},
			},
		},
		{
			name: "with multiple network",
			args: []string{"--ip-network", "family=IPv4,\"dns=one,two,three\",gateway=gw,dhcp", "--ip-network", "family=IPv6,\"dns=four\""},
			expected: request.CreateNetworkRequest{
				IPNetworks: []upcloud.IPNetwork{
					{
						Family:           "IPv4",
						DHCP:             upcloud.FromBool(true),
						DHCPDefaultRoute: upcloud.FromBool(false),
						DHCPDns:          []string{"one", "two", "three"},
						Gateway:          "gw",
					},
					{
						Family:           "IPv6",
						DHCP:             upcloud.FromBool(false),
						DHCPDefaultRoute: upcloud.FromBool(false),
						DHCPDns:          []string{"four"},
					},
				},
			},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			mns := MockNetworkService()
			cc := createCommand{
				BaseCommand: commands.New("create", "Create a network"),
				service:     mns,
				params:      createParams{},
			}
			comm := commands.BuildCommand(&cc, nil, config.New(viper.New()))
			mocks.SetFlags(comm, test.args)

			res, err := cc.BuildRequest()

			assert.Nil(t, err)
			assert.Equal(t, res, &test.expected)
		})
	}
}
