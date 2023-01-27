package nodegroup

import (
	"testing"

	"github.com/UpCloudLtd/upcloud-cli/v2/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v2/internal/config"
	smock "github.com/UpCloudLtd/upcloud-cli/v2/internal/mock"
	"github.com/UpCloudLtd/upcloud-cli/v2/internal/mockexecute"

	"github.com/UpCloudLtd/upcloud-go-api/v5/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/v5/upcloud/request"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestCreateKubernetesNodeGroup(t *testing.T) {
	clusterUUID := "9a8f4905-76d2-4a99-8e54-ab928fa42f66"

	for _, test := range []struct {
		name    string
		args    []string
		r       request.CreateKubernetesNodeGroupRequest
		wantErr bool
	}{
		{
			name: "1 nodegroup",
			args: []string{
				clusterUUID,
				"--count", "2",
				"--kubelet-arg=log-flush-frequency=5s",
				"--label=owner=devteam",
				"--label=env=dev",
				"--name=my-node-group",
				"--plan=K8S-2xCPU-4GB",
				"--ssh-key=ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIMWq/xsiYPgA/HLsaWHcjAGnwU+pJy9BUmvIlMBpkdn2 admin@user.com",
				"--storage=01000000-0000-4000-8000-000160010100",
				"--taint=env=dev:NoSchedule",
				"--taint=env=dev2:NoSchedule",
			},
			r: request.CreateKubernetesNodeGroupRequest{
				ClusterUUID: clusterUUID,
				NodeGroup: request.KubernetesNodeGroup{
					Count: 2,
					Labels: []upcloud.Label{
						{
							Key:   "owner",
							Value: "devteam",
						},
						{
							Key:   "env",
							Value: "dev",
						},
					},
					Name: "my-node-group",
					Plan: "K8S-2xCPU-4GB",
					SSHKeys: []string{
						"ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIMWq/xsiYPgA/HLsaWHcjAGnwU+pJy9BUmvIlMBpkdn2 admin@user.com",
					},
					Storage: "01000000-0000-4000-8000-000160010100",
					KubeletArgs: []upcloud.KubernetesKubeletArg{
						{
							Key:   "log-flush-frequency",
							Value: "5s",
						},
					},
					Taints: []upcloud.KubernetesTaint{
						{
							Effect: "NoSchedule",
							Key:    "env",
							Value:  "dev",
						},
						{
							Effect: "NoSchedule",
							Key:    "env",
							Value:  "dev2",
						},
					},
				},
			},
			wantErr: false,
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			conf := config.New()
			testCmd := CreateCommand()
			mService := new(smock.Service)

			mService.On("CreateKubernetesNodeGroup", &test.r).Return(&upcloud.KubernetesNodeGroup{}, nil)
			mService.On("GetNetworkDetails", mock.Anything).Return(&upcloud.Network{IPNetworks: []upcloud.IPNetwork{{Address: "172.16.1.0/24"}}}, nil)

			c := commands.BuildCommand(testCmd, nil, conf)

			c.Cobra().SetArgs(test.args)
			_, err := mockexecute.MockExecute(c, mService, conf)

			if test.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				mService.AssertNumberOfCalls(t, "CreateKubernetesNodeGroup", 1)
			}
		})
	}
}
