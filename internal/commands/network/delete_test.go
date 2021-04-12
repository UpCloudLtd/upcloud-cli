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

func TestDeleteCommand(t *testing.T) {
	targetMethod := "DeleteNetwork"

	n := upcloud.Network{UUID: "0a30b5ca-d0e3-4f7c-81d0-f77d42ea6366", Name: "test-network"}

	for _, test := range []struct {
		name  string
		arg   string
		flags []string
		error string
		req   request.DeleteNetworkRequest
	}{
		{
			name: "delete network with UUID",
			arg:  n.UUID,
			req:  request.DeleteNetworkRequest{UUID: n.UUID},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			cachedNetworks = nil
			mService := smock.Service{}
			mService.On(targetMethod, &test.req).Return(nil)
			mService.On("GetNetworks").Return(&upcloud.Networks{Networks: []upcloud.Network{n}}, nil)
			conf := config.New()
			c := commands.BuildCommand(DeleteCommand(), nil, conf)
			err := c.Cobra().Flags().Parse(test.flags)
			assert.NoError(t, err)

			_, err = c.(commands.Command).Execute(commands.NewExecutor(conf, &mService), test.arg)

			if test.error != "" {
				assert.Errorf(t, err, test.error)
			} else {
				mService.AssertNumberOfCalls(t, targetMethod, 1)
			}
		})
	}

}
