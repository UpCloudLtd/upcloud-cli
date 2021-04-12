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

	router := upcloud.Router{Name: "test-router", UUID: "123123"}
	modifiedRouter := upcloud.Router{Name: "test-router-b", UUID: "123123"}
	for _, test := range []struct {
		name    string
		flags   []string
		arg     string
		error   string
		returns *upcloud.Router
		req     request.ModifyRouterRequest
	}{
		{
			name:  "arg is missing",
			flags: []string{},
			error: "router is required",
		},
		{
			name:  "name is missing",
			arg:   router.UUID,
			flags: []string{},
			error: "name is required",
		},
		{
			name:    "name is passed",
			arg:     router.UUID,
			flags:   []string{"--name", "router-2-b"},
			returns: &modifiedRouter,
			req:     request.ModifyRouterRequest{Name: "router-2-b", UUID: router.UUID},
		},
	} {
		targetMethod := "ModifyRouter"
		t.Run(test.name, func(t *testing.T) {
			cachedRouters = nil
			mService := smock.Service{}
			mService.On(targetMethod, &test.req).Return(test.returns, nil)
			mService.On("GetRouters", mock.Anything).Return(&upcloud.Routers{Routers: []upcloud.Router{router}}, nil)

			conf := config.New()

			c := commands.BuildCommand(ModifyCommand(), nil, conf)
			if err := c.Cobra().Flags().Parse(test.flags); err != nil {
				t.Fatal(err)
			}
			_, err := c.(commands.Command).Execute(commands.NewExecutor(conf, &mService), test.arg)

			if test.error != "" {
				assert.Errorf(t, err, test.error)
			} else {
				mService.AssertNumberOfCalls(t, targetMethod, 1)
			}
		})
	}

}
