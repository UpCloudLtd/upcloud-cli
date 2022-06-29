package router

import (
	"testing"

	"github.com/UpCloudLtd/upcloud-cli/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/internal/config"
	smock "github.com/UpCloudLtd/upcloud-cli/internal/mock"
	"github.com/UpCloudLtd/upcloud-cli/internal/mockexecute"

	"github.com/UpCloudLtd/upcloud-go-api/v4/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/v4/upcloud/request"
	"github.com/stretchr/testify/assert"
)

func TestCreateCommand(t *testing.T) {
	targetMethod := "CreateRouter"

	router := upcloud.Router{Name: "test-router"}

	for _, test := range []struct {
		name  string
		flags []string
		error string
		req   request.CreateRouterRequest
	}{
		{
			name:  "name is missing",
			flags: []string{},
			error: `required flag(s) "name" not set`,
		},
		{
			name:  "name is passed",
			flags: []string{"--name", router.Name},
			req:   request.CreateRouterRequest{Name: router.Name},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			mService := smock.Service{}
			mService.On(targetMethod, &test.req).Return(&router, nil)

			conf := config.New()

			c := commands.BuildCommand(CreateCommand(), nil, conf)

			c.Cobra().SetArgs(test.flags)
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
