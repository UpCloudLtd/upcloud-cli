package router

import (
	"github.com/UpCloudLtd/cli/internal/commands"
	"github.com/UpCloudLtd/cli/internal/config"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/request"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestModifyCommand(t *testing.T) {
	methodName := "ModifyRouter"

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
			mrs := MockRouterService{}
			mrs.On(methodName, &test.req).Return(&router, nil)
			mrs.On("GetRouters", mock.Anything).Return(&upcloud.Routers{Routers: []upcloud.Router{router}}, nil)

			c := commands.BuildCommand(ModifyCommand(&mrs), nil, config.New(viper.New()))
			c.SetFlags(test.args)

			_, err := c.MakeExecuteCommand()([]string{router.Name})

			if test.error != "" {
				assert.Errorf(t, err, test.error)
			} else {
				mrs.AssertNumberOfCalls(t, methodName, 1)
			}
		})
	}

}
