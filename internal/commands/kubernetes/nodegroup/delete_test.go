package nodegroup

import (
	"testing"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/config"
	smock "github.com/UpCloudLtd/upcloud-cli/v3/internal/mock"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/mockexecute"

	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud/request"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDeleteKubernetesNodeGroup(t *testing.T) {
	clusterUUID := "898c4cf0-524c-4fc1-9c47-8cc697ed2d52"

	for _, test := range []struct {
		name     string
		args     []string
		expected request.DeleteKubernetesNodeGroupRequest
		errorMsg string
	}{
		{
			name:     "no args",
			args:     []string{clusterUUID},
			errorMsg: `required flag(s) "name" not set`,
		},
		{
			name: "delete success",
			args: []string{
				clusterUUID,
				"--name", "my-node-group",
			},
			expected: request.DeleteKubernetesNodeGroupRequest{
				ClusterUUID: clusterUUID,
				Name:        "my-node-group",
			},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			conf := config.New()
			testCmd := DeleteCommand()
			mService := new(smock.Service)

			expected := test.expected
			mService.On("DeleteKubernetesNodeGroup", &expected).Return(nil)

			c := commands.BuildCommand(testCmd, nil, conf)

			c.Cobra().SetArgs(test.args)
			_, err := mockexecute.MockExecute(c, mService, conf)

			if test.errorMsg != "" {
				assert.EqualError(t, err, test.errorMsg)
			} else {
				require.NoError(t, err)
				mService.AssertNumberOfCalls(t, "DeleteKubernetesNodeGroup", 1)
			}
		})
	}
}
