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

func TestModifyCommand(t *testing.T) {
	t.Parallel()
	testRouter := upcloud.Router{Name: "test-router", UUID: "123123"}
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
			arg:   testRouter.UUID,
			flags: []string{},
			error: "name is required",
		},
		{
			name:    "name is passed",
			arg:     testRouter.UUID,
			flags:   []string{"--name", "router-2-b"},
			returns: &modifiedRouter,
			req:     request.ModifyRouterRequest{Name: "router-2-b", UUID: testRouter.UUID},
		},
	} {
		// grab a local reference for parallel tests
		test := test
		targetMethod := "ModifyRouter"
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			mService := smock.Service{}
			mService.On(targetMethod, &test.req).Return(test.returns, nil)
			mService.On("GetRouters", mock.Anything).Return(&upcloud.Routers{Routers: []upcloud.Router{testRouter}}, nil)

			conf := config.New()

			c := commands.BuildCommand(router.ModifyCommand(), nil, conf)
			if err := c.Cobra().Flags().Parse(test.flags); err != nil {
				t.Fatal(err)
			}
			_, err := c.(commands.SingleArgumentCommand).ExecuteSingleArgument(commands.NewExecutor(conf, &mService, flume.New("test")), test.arg)

			if test.error != "" {
				assert.Errorf(t, err, test.error)
			} else {
				mService.AssertNumberOfCalls(t, targetMethod, 1)
			}
		})
	}
}
