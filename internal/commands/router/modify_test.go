package router

import (
	"testing"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/config"
	smock "github.com/UpCloudLtd/upcloud-cli/v3/internal/mock"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/mockexecute"

	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud/request"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestModifyCommand(t *testing.T) {
	router := upcloud.Router{Name: "test-router", UUID: "123123"}
	modifiedRouter := upcloud.Router{Name: "test-router-b", UUID: "123123"}

	for _, test := range []struct {
		name    string
		args    []string
		error   string
		returns *upcloud.Router
		req     request.ModifyRouterRequest
	}{
		{
			name:  "name is missing",
			args:  []string{router.UUID},
			error: `required flag(s) "name" not set`,
		},
		{
			name:    "name is passed",
			args:    []string{"--name", "New name", router.UUID},
			returns: &modifiedRouter,
			req:     request.ModifyRouterRequest{Name: "New name", UUID: router.UUID},
		},
	} {
		targetMethod := "ModifyRouter"
		t.Run(test.name, func(t *testing.T) {
			mService := smock.Service{}
			req := test.req
			mService.On(targetMethod, &req).Return(test.returns, nil)
			mService.On("GetRouters", mock.Anything).Return(&upcloud.Routers{Routers: []upcloud.Router{router}}, nil)

			conf := config.New()

			c := commands.BuildCommand(ModifyCommand(), nil, conf)

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
