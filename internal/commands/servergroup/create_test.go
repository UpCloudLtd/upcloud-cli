package servergroup

import (
	"testing"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/config"
	smock "github.com/UpCloudLtd/upcloud-cli/v3/internal/mock"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/mockexecute"

	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud/request"
	"github.com/stretchr/testify/assert"
)

func TestCreateServerGroup(t *testing.T) {
	serverGroupDefault := upcloud.ServerGroup{
		Title:              "test",
		AntiAffinityPolicy: upcloud.ServerGroupAntiAffinityPolicyBestEffort,
	}

	for _, test := range []struct {
		name  string
		args  []string
		req   request.CreateServerGroupRequest
		error string
	}{
		{
			name: "use default values",
			args: []string{
				"--title", "test",
			},
			req: request.CreateServerGroupRequest{
				Title:              "test",
				AntiAffinityPolicy: upcloud.ServerGroupAntiAffinityPolicyBestEffort,
			},
		},
		{
			name: "missing title",
			args: []string{},
			req: request.CreateServerGroupRequest{
				AntiAffinityPolicy: upcloud.ServerGroupAntiAffinityPolicyBestEffort,
			},
			error: `required flag(s) "title" not set`,
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			conf := config.New()
			testCmd := CreateCommand()
			mService := new(smock.Service)

			req := test.req
			mService.On("CreateServerGroup", &req).Return(&serverGroupDefault, nil)

			c := commands.BuildCommand(testCmd, nil, conf)

			c.Cobra().SetArgs(test.args)
			_, err := mockexecute.MockExecute(c, mService, conf)

			if test.error != "" {
				if err == nil {
					t.Errorf("expected error '%v', got nil", test.error)
				} else {
					assert.Equal(t, test.error, err.Error())
				}
			} else {
				mService.AssertNumberOfCalls(t, "CreateServerGroup", 1)
			}
		})
	}
}
