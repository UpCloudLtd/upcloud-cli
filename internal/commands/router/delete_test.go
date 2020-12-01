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

func TestDeleteCommand(t *testing.T) {
	methodName := "DeleteRouter"

	router := upcloud.Router{
		UUID: "97fbd082-30b0-11eb-adc1-0242ac120002",
		Name: "test-router",
	}

	for _, test := range []struct {
		name  string
		args  []string
		error string
		req   request.DeleteRouterRequest
	}{
		{
			name: "delete with UUID",
			args: []string{router.UUID},
			req:  request.DeleteRouterRequest{UUID: router.UUID},
		},
		{
			name: "delete with name",
			args: []string{router.Name},
			req:  request.DeleteRouterRequest{UUID: router.UUID},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			mrs := MockNetworkService{}
			mrs.On(methodName, &test.req).Return(nil)
			mrs.On("GetRouters", mock.Anything).Return(&upcloud.Routers{Routers: []upcloud.Router{router}}, nil)

			c := commands.BuildCommand(DeleteCommand(&mrs), nil, config.New(viper.New()))

			_, err := c.MakeExecuteCommand()(test.args)

			if test.error != "" {
				assert.Errorf(t, err, test.error)
			} else {
				mrs.AssertNumberOfCalls(t, methodName, 1)
			}
		})
	}

}
