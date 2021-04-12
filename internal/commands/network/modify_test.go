package network

import (
	"testing"

	"github.com/UpCloudLtd/cli/internal/commands"
	"github.com/UpCloudLtd/cli/internal/config"
	smock "github.com/UpCloudLtd/cli/internal/mock"

	"github.com/UpCloudLtd/upcloud-go-api/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/request"
	"github.com/stretchr/testify/assert"
)

func TestModifyCommand(t *testing.T) {
	targetMethod := "ModifyNetwork"

	n := upcloud.Network{
		UUID:   "9abccbe8-8d47-40dd-a5af-c6598f38b11b",
		Name:   "test-network",
		Zone:   "fi-hel1",
		Router: "",
	}

	for _, test := range []struct {
		name     string
		flags    []string
		error    string
		expected request.ModifyNetworkRequest
	}{
		{
			name: "family is missing",
			flags: []string{
				"--name", n.Name,
				"--ip-network", "gateway=gw,dhcp=true",
			},
			error: "family is required",
		},
		{
			name: "with single network",
			flags: []string{
				"--name", n.Name,
				"--ip-network", "family=IPv4,\"dhcp-dns=one,two,three\",gateway=gw,dhcp=true",
			},
			expected: request.ModifyNetworkRequest{
				UUID: n.UUID,
				Name: n.Name,
				IPNetworks: []upcloud.IPNetwork{
					{
						Family:  upcloud.IPAddressFamilyIPv4,
						DHCP:    upcloud.FromBool(true),
						DHCPDns: []string{"one", "two", "three"},
						Gateway: "gw",
					},
				},
			},
		},
		{
			name: "with multiple network",
			flags: []string{
				"--name", n.Name,
				"--ip-network", "\"dhcp-dns=one,two,three\",gateway=gw,dhcp=false,family=IPv4", "--ip-network", "family=IPv6,dhcp-dns=four",
			},
			expected: request.ModifyNetworkRequest{
				UUID: n.UUID,
				Name: n.Name,
				IPNetworks: []upcloud.IPNetwork{
					{
						Family:  upcloud.IPAddressFamilyIPv4,
						DHCP:    upcloud.FromBool(false),
						DHCPDns: []string{"one", "two", "three"},
						Gateway: "gw",
					},
					{
						Family:  upcloud.IPAddressFamilyIPv6,
						DHCPDns: []string{"four"},
					},
				},
			},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			cachedNetworks = nil
			mService := smock.Service{}
			mService.On(targetMethod, &test.expected).Return(&upcloud.Network{}, nil)
			mService.On("GetNetworks").Return(&upcloud.Networks{Networks: []upcloud.Network{n}}, nil)
			conf := config.New()
			c := commands.BuildCommand(ModifyCommand(), nil, conf)
			err := c.Cobra().Flags().Parse(test.flags)
			assert.NoError(t, err)

			_, err = c.(commands.Command).Execute(
				commands.NewExecutor(conf, &mService),
				n.UUID,
			)

			if err != nil {
				assert.EqualError(t, err, test.error)
			} else {
				mService.AssertNumberOfCalls(t, targetMethod, 1)
			}
		})
	}
}
