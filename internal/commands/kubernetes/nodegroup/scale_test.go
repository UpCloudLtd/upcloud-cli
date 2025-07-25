package nodegroup

import (
	"testing"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/config"
	smock "github.com/UpCloudLtd/upcloud-cli/v3/internal/mock"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/mockexecute"

	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud/request"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestScaleKubernetesNodeGroup(t *testing.T) {
	clusterUUID := "28c80353-98fd-4221-85e0-82d7603756ba"

	for _, test := range []struct {
		name     string
		args     []string
		expected request.ModifyKubernetesNodeGroupRequest
		errorMsg string
	}{
		{
			name:     "no args",
			args:     []string{clusterUUID},
			errorMsg: `required flag(s) "name", "count" not set`,
		},
		{
			name: "delete success",
			args: []string{
				clusterUUID,
				"--name", "my-node-group",
				"--count", "6",
			},
			expected: request.ModifyKubernetesNodeGroupRequest{
				ClusterUUID: clusterUUID,
				Name:        "my-node-group",
				NodeGroup:   request.ModifyKubernetesNodeGroup{Count: 6},
			},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			conf := config.New()
			testCmd := ScaleCommand()
			mService := new(smock.Service)

			expected := test.expected
			mService.On("ModifyKubernetesNodeGroup", &expected).Return(&upcloud.KubernetesNodeGroup{}, nil)

			c := commands.BuildCommand(testCmd, nil, conf)

			c.Cobra().SetArgs(test.args)
			_, err := mockexecute.MockExecute(c, mService, conf)

			if test.errorMsg != "" {
				assert.EqualError(t, err, test.errorMsg)
			} else {
				require.NoError(t, err)
				mService.AssertNumberOfCalls(t, "ModifyKubernetesNodeGroup", 1)
			}
		})
	}
}
