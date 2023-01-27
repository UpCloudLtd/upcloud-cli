package nodegroup

import (
	"testing"

	"github.com/UpCloudLtd/upcloud-cli/v2/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v2/internal/config"
	smock "github.com/UpCloudLtd/upcloud-cli/v2/internal/mock"
	"github.com/UpCloudLtd/upcloud-cli/v2/internal/mockexecute"

	"github.com/UpCloudLtd/upcloud-go-api/v5/upcloud/request"
	"github.com/stretchr/testify/require"
)

func TestDeleteKubernetesNodeGroup(t *testing.T) {
	clusterUUID := "898c4cf0-524c-4fc1-9c47-8cc697ed2d52"

	for _, test := range []struct {
		name    string
		args    []string
		r       request.DeleteKubernetesNodeGroupRequest
		wantErr bool
	}{
		{
			name: "delete success",
			args: []string{
				clusterUUID,
				"--name", "my-node-group",
			},
			r: request.DeleteKubernetesNodeGroupRequest{
				ClusterUUID: clusterUUID,
				Name:        "my-node-group",
			},
			wantErr: false,
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			conf := config.New()
			testCmd := DeleteCommand()
			mService := new(smock.Service)

			mService.On("DeleteKubernetesNodeGroup", &test.r).Return(nil)

			c := commands.BuildCommand(testCmd, nil, conf)

			c.Cobra().SetArgs(test.args)
			_, err := mockexecute.MockExecute(c, mService, conf)

			if test.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				mService.AssertNumberOfCalls(t, "DeleteKubernetesNodeGroup", 1)
			}
		})
	}
}
