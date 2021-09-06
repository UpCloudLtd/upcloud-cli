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
)

func TestCreateCommand(t *testing.T) {
	t.Parallel()
	targetMethod := "CreateRouter"

	testRouter := upcloud.Router{Name: "test-router"}

	for _, test := range []struct {
		name  string
		flags []string
		error string
		req   request.CreateRouterRequest
	}{
		{
			name:  "name is missing",
			flags: []string{},
			error: "name is required",
		},
		{
			name:  "name is passed",
			flags: []string{"--name", testRouter.Name},
			req:   request.CreateRouterRequest{Name: testRouter.Name},
		},
	} {
		// grab a local reference for parallel tests
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			mService := smock.Service{}
			mService.On(targetMethod, &test.req).Return(&testRouter, nil)

			conf := config.New()

			c := commands.BuildCommand(router.CreateCommand(), nil, conf)
			err := c.Cobra().Flags().Parse(test.flags)
			assert.NoError(t, err)

			_, err = c.(commands.NoArgumentCommand).ExecuteWithoutArguments(commands.NewExecutor(conf, &mService, flume.New("test")))

			if test.error != "" {
				assert.Errorf(t, err, test.error)
			} else {
				mService.AssertNumberOfCalls(t, targetMethod, 1)
			}
		})
	}
}
