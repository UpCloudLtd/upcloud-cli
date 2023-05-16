package nodegroup

import (
	"testing"

	"github.com/UpCloudLtd/upcloud-cli/v2/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v2/internal/config"
	smock "github.com/UpCloudLtd/upcloud-cli/v2/internal/mock"
	"github.com/UpCloudLtd/upcloud-cli/v2/internal/mockexecute"

	"github.com/UpCloudLtd/upcloud-go-api/v6/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/v6/upcloud/request"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestCreateKubernetesNodeGroup(t *testing.T) {
	clusterUUID := "9a8f4905-76d2-4a99-8e54-ab928fa42f66"

	for _, test := range []struct {
		name     string
		args     []string
		expected request.CreateKubernetesNodeGroupRequest
		errorMsg string
	}{
		{
			name:     "no args",
			args:     []string{clusterUUID},
			errorMsg: `required flag(s) "count", "name", "plan" not set`,
		},
		{
			name:     "no name",
			args:     []string{clusterUUID, "--name", "my-node-group"},
			errorMsg: `required flag(s) "count", "plan" not set`,
		},
		{
			name: "simple nodegroup",
			args: []string{clusterUUID, "--count", "2", "--name", "my-node-group", "--plan", "2xCPU-4GB"},
			expected: request.CreateKubernetesNodeGroupRequest{
				ClusterUUID: clusterUUID,
				NodeGroup: request.KubernetesNodeGroup{
					Count:       2,
					Name:        "my-node-group",
					Plan:        "2xCPU-4GB",
					Labels:      []upcloud.Label{},
					KubeletArgs: []upcloud.KubernetesKubeletArg{},
					Taints:      []upcloud.KubernetesTaint{},
				},
			},
		},
		{
			name: "simple nodegroup with labels",
			args: []string{clusterUUID, "--count", "2", "--name", "my-node-group", "--plan", "2xCPU-4GB", "--label", "key=value", "--label", "key-without-value"},
			expected: request.CreateKubernetesNodeGroupRequest{
				ClusterUUID: clusterUUID,
				NodeGroup: request.KubernetesNodeGroup{
					Count: 2,
					Name:  "my-node-group",
					Plan:  "2xCPU-4GB",
					Labels: []upcloud.Label{
						{Key: "key", Value: "value"},
						{Key: "key-without-value", Value: ""},
					},
					KubeletArgs: []upcloud.KubernetesKubeletArg{},
					Taints:      []upcloud.KubernetesTaint{},
				},
			},
		},
		{
			name: "complex nodegroup",
			args: []string{
				clusterUUID,
				"--count", "2",
				"--kubelet-arg=log-flush-frequency=5s",
				"--label=owner=devteam",
				"--label=env=dev",
				"--name=my-node-group",
				"--plan=2xCPU-4GB",
				"--ssh-key=ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIMWq/xsiYPgA/HLsaWHcjAGnwU+pJy9BUmvIlMBpkdn2 admin@user.com",
				"--storage=01000000-0000-4000-8000-000160010100",
				"--taint=env=dev:NoSchedule",
				"--taint=env=dev2:NoSchedule",
			},
			expected: request.CreateKubernetesNodeGroupRequest{
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
					Plan: "2xCPU-4GB",
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
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			conf := config.New()
			testCmd := CreateCommand()
			mService := new(smock.Service)

			mService.On("CreateKubernetesNodeGroup", &test.expected).Return(&upcloud.KubernetesNodeGroup{}, nil)
			mService.On("GetNetworkDetails", mock.Anything).Return(&upcloud.Network{IPNetworks: []upcloud.IPNetwork{{Address: "172.16.1.0/24"}}}, nil)

			c := commands.BuildCommand(testCmd, nil, conf)

			c.Cobra().SetArgs(test.args)
			_, err := mockexecute.MockExecute(c, mService, conf)

			if test.errorMsg != "" {
				assert.EqualError(t, err, test.errorMsg)
			} else {
				require.NoError(t, err)
				mService.AssertNumberOfCalls(t, "CreateKubernetesNodeGroup", 1)
			}
		})
	}
}
