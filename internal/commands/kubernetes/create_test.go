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
	network := upcloud.Network{
		UUID:       "aa39e313-d908-418a-a959-459699bdc83a",
		Name:       "test-network",
		IPNetworks: []upcloud.IPNetwork{{Address: "172.16.1.0/24"}},
	}
	networks := upcloud.Networks{Networks: []upcloud.Network{network}}

	nodeGroupArgs := func(network string) []string {
		return []string{
			"--name", "my-cluster",
			"--network", network,
			"--node-group", "count=2,kubelet-arg=log-flush-frequency=5s,label=owner=devteam,label=env=dev,name=my-node-group,plan=2xCPU-4GB,ssh-key=ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIMWq/xsiYPgA/HLsaWHcjAGnwU+pJy9BUmvIlMBpkdn2 admin@user.com,storage=01000000-0000-4000-8000-000160010100,taint=env=dev:NoSchedule,taint=env=dev2:NoSchedule",
			"--node-group", "count=1,name=my-node-group2,plan=2xCPU-4GB,ssh-key=ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIMWq/xsiYPgA/HLsaWHcjAGnwU+pJy9BUmvIlMBpkdn2 admin@user.com,disable-utility-network-access",
			"--zone", "de-fra1",
		}
	}
	nodeGroupRequest := request.CreateKubernetesClusterRequest{
		Name:        "my-cluster",
		Network:     "aa39e313-d908-418a-a959-459699bdc83a",
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
				UtilityNetworkAccess: upcloud.BoolPtr(true),
			},
			{
				Count:       1,
				KubeletArgs: []upcloud.KubernetesKubeletArg{},
				Labels:      []upcloud.Label{},
				Name:        "my-node-group2",
				Plan:        "2xCPU-4GB",
				SSHKeys: []string{
					"ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIMWq/xsiYPgA/HLsaWHcjAGnwU+pJy9BUmvIlMBpkdn2 admin@user.com",
				},
				Taints:               []upcloud.KubernetesTaint{},
				UtilityNetworkAccess: upcloud.BoolPtr(false),
			},
		},
		Plan: "development",
		Zone: "de-fra1",
	}

	prodArg := []string{"--plan", "production-small"}
	prodPlanRequest := nodeGroupRequest
	prodPlanRequest.Plan = "production-small"

	privateNodeGroupsArg := []string{"--private-node-groups"}
	privateNodeGroupsRequest := nodeGroupRequest
	privateNodeGroupsRequest.PrivateNodeGroups = true

	for _, test := range []struct {
		name    string
		args    []string
		request request.CreateKubernetesClusterRequest
		wantErr bool
	}{
		{
			name:    "2 node groups",
			args:    nodeGroupArgs(network.UUID),
			request: nodeGroupRequest,
			wantErr: false,
		},
		{
			name:    "resolve network from name",
			args:    nodeGroupArgs(network.Name),
			request: nodeGroupRequest,
			wantErr: false,
		},
		{
			name:    "use productions-small plan",
			args:    append(nodeGroupArgs(network.Name), prodArg...),
			request: prodPlanRequest,
			wantErr: false,
		},
		{
			name:    "with private node groups",
			args:    append(nodeGroupArgs(network.Name), privateNodeGroupsArg...),
			request: privateNodeGroupsRequest,
			wantErr: false,
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			conf := config.New()
			testCmd := CreateCommand()
			mService := new(smock.Service)

			mService.On("CreateKubernetesCluster", &test.request).Return(&upcloud.KubernetesCluster{}, nil)
			mService.On("GetNetworks").Return(&networks, nil)
			mService.On("GetNetworkDetails", mock.Anything).Return(&network, nil)

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
