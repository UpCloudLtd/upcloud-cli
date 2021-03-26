package networkinterface

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
	targetMethod := "ModifyNetworkInterface"

	s := upcloud.Server{UUID: "c4cb35bc-3fb5-4cce-9951-79cab2225417"}
	servers := upcloud.Servers{Servers: []upcloud.Server{s}}
	network := upcloud.Network{UUID: "aa39e313-d908-418a-a959-459699bdc83a", Name: "test-network"}
	networks := upcloud.Networks{Networks: []upcloud.Network{network}}

	for _, test := range []struct {
		name  string
		args  []string
		error string
		req   request.ModifyNetworkInterfaceRequest
	}{
		{
			name:  "index is missing",
			args:  []string{},
			error: "index is required",
		},
		{
			name: "index is present, using default values",
			args: []string{
				"--index", "4",
			},
			req: request.ModifyNetworkInterfaceRequest{CurrentIndex: 4, ServerUUID: s.UUID},
		},
		{
			name: "index is present, all values modified",
			args: []string{
				"--index", "4",
				"--new-index", "5",
				"--bootable", "false",
				"--source-ip-filtering", "true",
				"--ip-addresses", "127.0.0.2,127.0.0.3,127.0.0.4",
			},
			req: request.ModifyNetworkInterfaceRequest{
				ServerUUID:        s.UUID,
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
			mService := smock.MockService{}
			mService.On("GetNetworks").Return(&networks, nil)
			mService.On(targetMethod, &test.req).Return(&upcloud.Interface{}, nil)

			mService.On("GetServers").Return(&servers, nil)
			c := commands.BuildCommand(ModifyCommand(&mService, &mService), nil, config.New())
			err := c.SetFlags(test.args)
			assert.NoError(t, err)

			_, err = c.MakeExecuteCommand()([]string{s.UUID})

			if test.error != "" {
				assert.Errorf(t, err, test.error)
			} else {
				mService.AssertNumberOfCalls(t, targetMethod, 1)
			}
		})
	}

}
