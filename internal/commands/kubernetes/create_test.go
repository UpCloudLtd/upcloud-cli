package kubernetes

import (
	"testing"

	"github.com/UpCloudLtd/upcloud-cli/v2/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v2/internal/config"
	smock "github.com/UpCloudLtd/upcloud-cli/v2/internal/mock"
	"github.com/UpCloudLtd/upcloud-cli/v2/internal/mockexecute"

	"github.com/UpCloudLtd/upcloud-go-api/v6/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/v6/upcloud/request"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestCreateKubernetes(t *testing.T) {
	for _, test := range []struct {
		name    string
		args    []string
		r       request.CreateKubernetesClusterRequest
		wantErr bool
	}{
		{
			name: "1 node group",
			args: []string{
				"--name", "my-cluster",
				"--network", "03e5ca07-f36c-4957-a676-e001e40441eb",
				"--node-group", "count=2,kubelet-arg=log-flush-frequency=5s,label=owner=devteam,label=env=dev,name=my-node-group,plan=K8S-2xCPU-4GB,ssh-key=ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIMWq/xsiYPgA/HLsaWHcjAGnwU+pJy9BUmvIlMBpkdn2 admin@user.com,storage=01000000-0000-4000-8000-000160010100,taint=env=dev:NoSchedule,taint=env=dev2:NoSchedule",
				"--zone", "de-fra1",
			},
			r: request.CreateKubernetesClusterRequest{
				Name:        "my-cluster",
				Network:     "03e5ca07-f36c-4957-a676-e001e40441eb",
				NetworkCIDR: "172.16.1.0/24",
				NodeGroups: []request.KubernetesNodeGroup{
					{
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
				Zone: "de-fra1",
			},
			wantErr: false,
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			conf := config.New()
			testCmd := CreateCommand()
			mService := new(smock.Service)

			mService.On("CreateKubernetesCluster", &test.r).Return(&upcloud.KubernetesCluster{}, nil)
			mService.On("GetNetworkDetails", mock.Anything).Return(&upcloud.Network{IPNetworks: []upcloud.IPNetwork{{Address: "172.16.1.0/24"}}}, nil)

			c := commands.BuildCommand(testCmd, nil, conf)

			c.Cobra().SetArgs(test.args)
			_, err := mockexecute.MockExecute(c, mService, conf)

			if test.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				mService.AssertNumberOfCalls(t, "CreateKubernetesCluster", 1)
			}
		})
	}
}
