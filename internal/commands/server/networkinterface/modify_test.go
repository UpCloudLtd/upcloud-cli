package networkinterface

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

func TestModifyCommand(t *testing.T) {
	targetMethod := "ModifyNetworkInterface"

	server := upcloud.Server{UUID: "c4cb35bc-3fb5-4cce-9951-79cab2225417"}
	servers := upcloud.Servers{Servers: []upcloud.Server{server}}
	network := upcloud.Network{UUID: "aa39e313-d908-418a-a959-459699bdc83a", Name: "test-network"}
	networks := upcloud.Networks{Networks: []upcloud.Network{network}}

	for _, test := range []struct {
		name  string
		flags []string
		error string
		req   request.ModifyNetworkInterfaceRequest
	}{
		{
			name:  "index is missing",
			flags: []string{},
			error: `required flag(s) "index" not set`,
		},
		{
			name: "index is present, using default values",
			flags: []string{
				"--index", "4",
			},
			req: request.ModifyNetworkInterfaceRequest{CurrentIndex: 4, ServerUUID: server.UUID},
		},
		{
			name: "index is present, all values modified",
			flags: []string{
				"--index", "4",
				"--new-index", "5",
				"--bootable", "false",
				"--source-ip-filtering", "true",
				"--ip-addresses", "127.0.0.2,127.0.0.3,127.0.0.4",
			},
			req: request.ModifyNetworkInterfaceRequest{
				ServerUUID:        server.UUID,
				CurrentIndex:      4,
				NewIndex:          5,
				Bootable:          upcloud.FromBool(false),
				SourceIPFiltering: upcloud.FromBool(true),
				IPAddresses: request.CreateNetworkInterfaceIPAddressSlice{
					{Address: "127.0.0.2", Family: upcloud.IPAddressFamilyIPv4},
					{Address: "127.0.0.3", Family: upcloud.IPAddressFamilyIPv4},
					{Address: "127.0.0.4", Family: upcloud.IPAddressFamilyIPv4},
				},
			},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			mService := smock.Service{}
			mService.On("GetNetworks").Return(&networks, nil)
			req := test.req
			mService.On(targetMethod, &req).Return(&upcloud.Interface{}, nil)

			mService.On("GetServers").Return(&servers, nil)
			conf := config.New()

			c := commands.BuildCommand(ModifyCommand(), nil, conf)

			c.Cobra().SetArgs(append(test.flags, server.UUID))
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
