package network

import (
	"testing"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/config"
	smock "github.com/UpCloudLtd/upcloud-cli/v3/internal/mock"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/mockexecute"

	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud/request"
	"github.com/stretchr/testify/assert"
)

func TestCreateCommand(t *testing.T) {
	targetMethod := "CreateNetwork"

	n := upcloud.Network{
		UUID:   "9abccbe8-8d47-40dd-a5af-c6598f38b11b",
		Name:   "test-network",
		Zone:   "fi-hel1",
		Router: "",
	}

	for _, test := range []struct {
		name     string
		args     []string
		error    string
		expected request.CreateNetworkRequest
	}{
		{
			name:  "name is missing",
			args:  []string{"--zone", n.Zone},
			error: `required flag(s) "name", "ip-network" not set`,
		},
		{
			name:  "zone is missing",
			args:  []string{"--name", n.Name},
			error: `required flag(s) "zone", "ip-network" not set`,
		},
		{
			name:  "without network",
			args:  []string{"--name", n.Name, "--zone", n.Zone},
			error: `required flag(s) "ip-network" not set`,
		},
		{
			name: "with single network",
			args: []string{
				"--name", n.Name,
				"--zone", n.Zone,
				"--ip-network", "address=127.0.0.1,\"dhcp-dns=one,two,three\",gateway=gw,dhcp=true",
			},
			expected: request.CreateNetworkRequest{
				Name: n.Name,
				Zone: n.Zone,
				IPNetworks: []upcloud.IPNetwork{
					{
						Address: "127.0.0.1",
						Family:  upcloud.IPAddressFamilyIPv4,
						DHCP:    upcloud.FromBool(true),
						DHCPDns: []string{"one", "two", "three"},
						Gateway: "gw",
					},
				},
			},
		},
		{
			name: "with DHCP parameters",
			args: []string{
				"--name", n.Name,
				"--zone", n.Zone,
				"--ip-network", "address=127.0.0.1,dhcp=true,dhcp-default-route=true",
			},
			expected: request.CreateNetworkRequest{
				Name: n.Name,
				Zone: n.Zone,
				IPNetworks: []upcloud.IPNetwork{
					{
						Address:          "127.0.0.1",
						Family:           upcloud.IPAddressFamilyIPv4,
						DHCP:             upcloud.FromBool(true),
						DHCPDefaultRoute: upcloud.FromBool(true),
					},
				},
			},
		},
		{
			name: "with multiple network",
			args: []string{
				"--name", n.Name,
				"--zone", n.Zone,
				"--ip-network", "\"dhcp-dns=one,two,three\",gateway=gw,dhcp=false,address=127.0.0.1", "--ip-network", "address=2001:0db8:85a3:0000:0000:8a2e:0370:7334/32,dhcp-dns=four",
			},
			expected: request.CreateNetworkRequest{
				Name: n.Name,
				Zone: n.Zone,
				IPNetworks: []upcloud.IPNetwork{
					{
						Address: "127.0.0.1",
						Family:  upcloud.IPAddressFamilyIPv4,
						DHCP:    upcloud.FromBool(false),
						DHCPDns: []string{"one", "two", "three"},
						Gateway: "gw",
					},
					{
						Address: "2001:0db8:85a3:0000:0000:8a2e:0370:7334/32",
						Family:  upcloud.IPAddressFamilyIPv6,
						DHCPDns: []string{"four"},
					},
				},
			},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			mService := smock.Service{}
			expected := test.expected
			mService.On(targetMethod, &expected).Return(&upcloud.Network{}, nil)
			conf := config.New()

			c := commands.BuildCommand(CreateCommand(), nil, conf)

			c.Cobra().SetArgs(test.args)
			_, err := mockexecute.MockExecute(c, &mService, conf)

			if test.error != "" {
				assert.EqualError(t, err, test.error)
			} else {
				assert.NoError(t, err)
				mService.AssertNumberOfCalls(t, targetMethod, 1)
			}
		})
	}
}
