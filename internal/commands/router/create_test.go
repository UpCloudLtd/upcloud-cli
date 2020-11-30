package router

import (
	"github.com/UpCloudLtd/cli/internal/commands"
	"github.com/UpCloudLtd/cli/internal/config"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/request"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreateCommand(t *testing.T) {
	methodName := "CreateRouter"

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
			mrs := MockRouterService{}
			mrs.On(methodName, &test.req).Return(&router, nil)

			c := commands.BuildCommand(CreateCommand(&mrs), nil, config.New(viper.New()))
			c.SetFlags(test.args)

			_, err := c.MakeExecuteCommand()([]string{})

			if test.error != "" {
				assert.Errorf(t, err, test.error)
			} else {
				mrs.AssertNumberOfCalls(t, methodName, 1)
			}
		})
	}

}
