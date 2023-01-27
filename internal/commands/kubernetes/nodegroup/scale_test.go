package nodegroup

import (
	"testing"

	"github.com/UpCloudLtd/upcloud-cli/v2/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v2/internal/config"
	smock "github.com/UpCloudLtd/upcloud-cli/v2/internal/mock"
	"github.com/UpCloudLtd/upcloud-cli/v2/internal/mockexecute"

	"github.com/UpCloudLtd/upcloud-go-api/v5/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/v5/upcloud/request"
	"github.com/stretchr/testify/require"
)

func TestScaleKubernetesNodeGroup(t *testing.T) {
	clusterUUID := "28c80353-98fd-4221-85e0-82d7603756ba"

	for _, test := range []struct {
		name    string
		args    []string
		r       request.ModifyKubernetesNodeGroupRequest
		wantErr bool
	}{
		{
			name: "delete success",
			args: []string{
				clusterUUID,
				"--name", "my-node-group",
				"--count=6",
			},
			r: request.ModifyKubernetesNodeGroupRequest{
				ClusterUUID: clusterUUID,
				Name:        "my-node-group",
				NodeGroup:   request.ModifyKubernetesNodeGroup{Count: 6},
			},
			wantErr: false,
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			conf := config.New()
			testCmd := ScaleCommand()
			mService := new(smock.Service)

			mService.On("ModifyKubernetesNodeGroup", &test.r).Return(&upcloud.KubernetesNodeGroup{}, nil)

			c := commands.BuildCommand(testCmd, nil, conf)

			c.Cobra().SetArgs(test.args)
			_, err := mockexecute.MockExecute(c, mService, conf)

			if test.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				mService.AssertNumberOfCalls(t, "ModifyKubernetesNodeGroup", 1)
			}
		})
	}
}
