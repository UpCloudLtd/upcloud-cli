package network

import (
	"github.com/UpCloudLtd/cli/internal/commands"
	"github.com/UpCloudLtd/cli/internal/config"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/request"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestModifyCommand(t *testing.T) {
	methodName := "ModifyNetwork"

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
			mns := MockNetworkService{}
			mns.On(methodName, &test.expected).Return(&upcloud.Network{}, nil)
			mns.On("GetNetworks").Return(&upcloud.Networks{Networks: []upcloud.Network{n}}, nil)
			c := commands.BuildCommand(ModifyCommand(&mns), nil, config.New(viper.New()))
			c.SetFlags(test.flags)

			_, err := c.MakeExecuteCommand()([]string{n.Name})

			if err != nil {
				assert.Equal(t, test.error, err.Error())
			} else {
				mns.AssertNumberOfCalls(t, methodName, 1)
			}
		})
	}
}
