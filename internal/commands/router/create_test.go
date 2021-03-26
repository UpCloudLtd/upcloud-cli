package router

import (
	"testing"

	"github.com/UpCloudLtd/cli/internal/commands"
	"github.com/UpCloudLtd/cli/internal/config"
	smock "github.com/UpCloudLtd/cli/internal/mock"

	"github.com/UpCloudLtd/upcloud-go-api/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/request"
	"github.com/stretchr/testify/assert"
)

func TestCreateCommand(t *testing.T) {
	targetMethod := "CreateRouter"

	router := upcloud.Router{Name: "test-router"}

	for _, test := range []struct {
		name  string
		args  []string
		error string
		req   request.CreateRouterRequest
	}{
		{
			name:  "name is missing",
			args:  []string{},
			error: "name is required",
		},
		{
			name: "name is passed",
			args: []string{"--name", router.Name},
			req:  request.CreateRouterRequest{Name: router.Name},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			mService := smock.MockService{}
			mService.On(targetMethod, &test.req).Return(&router, nil)

			c := commands.BuildCommand(CreateCommand(&mService), nil, config.New())
			err := c.SetFlags(test.args)
			assert.NoError(t, err)

			_, err = c.MakeExecuteCommand()([]string{})

			if test.error != "" {
				assert.Errorf(t, err, test.error)
			} else {
				mService.AssertNumberOfCalls(t, targetMethod, 1)
			}
		})
	}

}
