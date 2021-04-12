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

func TestDeleteCommand(t *testing.T) {
	targetMethod := "DeleteRouter"

	router := upcloud.Router{
		UUID: "97fbd082-30b0-11eb-adc1-0242ac120002",
		Name: "test-router",
	}

	for _, test := range []struct {
		name  string
		arg   string
		error string
		req   request.DeleteRouterRequest
	}{
		{
			name: "delete with UUID",
			arg:  router.UUID,
			req:  request.DeleteRouterRequest{UUID: router.UUID},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			mService := smock.Service{}
			mService.On(targetMethod, &test.req).Return(nil)
			mService.On("GetRouters", mock.Anything).Return(&upcloud.Routers{Routers: []upcloud.Router{router}}, nil)

			conf := config.New()

			c := commands.BuildCommand(DeleteCommand(), nil, conf)

			_, err := c.(commands.Command).Execute(commands.NewExecutor(conf, &mService), test.arg)

			if test.error != "" {
				assert.Errorf(t, err, test.error)
			} else {
				mService.AssertNumberOfCalls(t, targetMethod, 1)
			}
		})
	}

}
