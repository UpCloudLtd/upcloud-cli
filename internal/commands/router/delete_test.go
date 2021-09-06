package router_test

import (
	"testing"

	"github.com/UpCloudLtd/upcloud-cli/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/internal/commands/router"
	"github.com/UpCloudLtd/upcloud-cli/internal/config"
	smock "github.com/UpCloudLtd/upcloud-cli/internal/mock"

	"github.com/UpCloudLtd/upcloud-go-api/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/request"
	"github.com/gemalto/flume"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestDeleteCommand(t *testing.T) {
	t.Parallel()
	targetMethod := "DeleteRouter"

	testRouter := upcloud.Router{
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
			arg:  testRouter.UUID,
			req:  request.DeleteRouterRequest{UUID: testRouter.UUID},
		},
	} {
		// grab a local reference for parallel tests
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			mService := smock.Service{}
			mService.On(targetMethod, &test.req).Return(nil)
			mService.On("GetRouters", mock.Anything).Return(&upcloud.Routers{Routers: []upcloud.Router{testRouter}}, nil)

			conf := config.New()

			c := commands.BuildCommand(router.DeleteCommand(), nil, conf)

			_, err := c.(commands.MultipleArgumentCommand).Execute(commands.NewExecutor(conf, &mService, flume.New("test")), test.arg)

			if test.error != "" {
				assert.Errorf(t, err, test.error)
			} else {
				mService.AssertNumberOfCalls(t, targetMethod, 1)
			}
		})
	}
}
