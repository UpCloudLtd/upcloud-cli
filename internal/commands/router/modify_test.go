package router

import (
	"testing"

	"github.com/UpCloudLtd/cli/internal/commands"
	"github.com/UpCloudLtd/cli/internal/config"
	smock "github.com/UpCloudLtd/cli/internal/mock"

	"github.com/UpCloudLtd/upcloud-go-api/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/request"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestModifyCommand(t *testing.T) {
	targetMethod := "ModifyRouter"

	router := upcloud.Router{Name: "test-router"}

	for _, test := range []struct {
		name  string
		args  []string
		error string
		req   request.ModifyRouterRequest
	}{
		{
			name:  "name is missing",
			args:  []string{},
			error: "name is required",
		},
		{
			name: "name is passed",
			args: []string{"--name", "router-2"},
			req:  request.ModifyRouterRequest{Name: "router-2"},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			cachedRouters = nil
			mService := smock.MockService{}
			mService.On(targetMethod, &test.req).Return(&router, nil)
			mService.On("GetRouters", mock.Anything).Return(&upcloud.Routers{Routers: []upcloud.Router{router}}, nil)

			c := commands.BuildCommand(ModifyCommand(&mService), nil, config.New())
			err := c.SetFlags(test.args)
			assert.NoError(t, err)

			_, err = c.MakeExecuteCommand()([]string{router.Name})

			if test.error != "" {
				assert.Errorf(t, err, test.error)
			} else {
				mService.AssertNumberOfCalls(t, targetMethod, 1)
			}
		})
	}

}
