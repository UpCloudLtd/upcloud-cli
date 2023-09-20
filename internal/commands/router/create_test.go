package router

import (
	"fmt"
	"testing"

	"github.com/UpCloudLtd/upcloud-cli/v2/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v2/internal/config"
	smock "github.com/UpCloudLtd/upcloud-cli/v2/internal/mock"
	"github.com/UpCloudLtd/upcloud-cli/v2/internal/mockexecute"

	"github.com/UpCloudLtd/upcloud-go-api/v6/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/v6/upcloud/request"
	"github.com/stretchr/testify/assert"
)

func TestCreateCommand(t *testing.T) {
	targetMethod := "CreateRouter"

	router := upcloud.Router{
		Name: "test-router",
		StaticRoutes: []upcloud.StaticRoute{
			{
				Name:    "test-static-route",
				Route:   "0.0.0.0/0",
				Nexthop: "10.0.0.100",
			},
		},
	}

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
		{
			name: "name and a static route are passed",
			flags: []string{
				"--name", router.Name,
				"--static-route", fmt.Sprintf("name=%s,nexthop=%s,route=%s",
					router.StaticRoutes[0].Name,
					router.StaticRoutes[0].Nexthop,
					router.StaticRoutes[0].Route,
				),
			},
			req: request.CreateRouterRequest{
				Name: router.Name,
				StaticRoutes: []upcloud.StaticRoute{
					{
						Name:    "test-static-route",
						Route:   "0.0.0.0/0",
						Nexthop: "10.0.0.100",
					},
				},
			},
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
