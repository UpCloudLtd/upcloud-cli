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
					Count:                2,
					Name:                 "my-node-group",
					Plan:                 "2xCPU-4GB",
					Labels:               []upcloud.Label{},
					KubeletArgs:          []upcloud.KubernetesKubeletArg{},
					Taints:               []upcloud.KubernetesTaint{},
					UtilityNetworkAccess: upcloud.BoolPtr(true),
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
					KubeletArgs:          []upcloud.KubernetesKubeletArg{},
					Taints:               []upcloud.KubernetesTaint{},
					UtilityNetworkAccess: upcloud.BoolPtr(true),
				},
			},
		},
		{
			name: "complex nodegroup",
			args: []string{
				clusterUUID,
				"--count", "2",
				"--kubelet-arg", "log-flush-frequency=5s",
				"--label", "owner=devteam",
				"--label", "env=dev",
				"--name", "my-node-group",
				"--plan", "2xCPU-4GB",
				"--ssh-key", "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIMWq/xsiYPgA/HLsaWHcjAGnwU+pJy9BUmvIlMBpkdn2 admin@user.com",
				"--storage", "01000000-0000-4000-8000-000160010100",
				"--taint", "env=dev:NoSchedule",
				"--taint", "env=dev2:NoSchedule",
				"--disable-utility-network-access",
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
					UtilityNetworkAccess: upcloud.BoolPtr(false),
				},
			},
		},
		{
			name: "GPU plan with storage customization",
			args: []string{clusterUUID, "--count", "2", "--name", "gpu-nodes", "--plan", "GPU-8xCPU-64GB-1xL40S", "--storage-size", "1024", "--storage-tier", "maxiops"},
			expected: request.CreateKubernetesNodeGroupRequest{
				ClusterUUID: clusterUUID,
				NodeGroup: request.KubernetesNodeGroup{
					Count:                2,
					Name:                 "gpu-nodes",
					Plan:                 "GPU-8xCPU-64GB-1xL40S",
					Labels:               []upcloud.Label{},
					KubeletArgs:          []upcloud.KubernetesKubeletArg{},
					Taints:               []upcloud.KubernetesTaint{},
					UtilityNetworkAccess: upcloud.BoolPtr(true),
					GPUPlan: &upcloud.KubernetesNodeGroupGPUPlan{
						StorageSize: 1024,
						StorageTier: upcloud.StorageTierMaxIOPS,
					},
				},
			},
		},
		{
			name: "Cloud Native plan with storage customization",
			args: []string{clusterUUID, "--count", "3", "--name", "cloud-native-nodes", "--plan", "CLOUDNATIVE-4xCPU-8GB", "--storage-size", "50", "--storage-tier", "standard"},
			expected: request.CreateKubernetesNodeGroupRequest{
				ClusterUUID: clusterUUID,
				NodeGroup: request.KubernetesNodeGroup{
					Count:                3,
					Name:                 "cloud-native-nodes",
					Plan:                 "CLOUDNATIVE-4xCPU-8GB",
					Labels:               []upcloud.Label{},
					KubeletArgs:          []upcloud.KubernetesKubeletArg{},
					Taints:               []upcloud.KubernetesTaint{},
					UtilityNetworkAccess: upcloud.BoolPtr(true),
					CloudNativePlan: &upcloud.KubernetesNodeGroupCloudNativePlan{
						StorageSize: 50,
						StorageTier: upcloud.StorageTierStandard,
					},
				},
			},
		},
		{
			name:     "storage customization with unsupported plan",
			args:     []string{clusterUUID, "--count", "2", "--name", "regular-nodes", "--plan", "2xCPU-4GB", "--storage-size", "100"},
			errorMsg: "storage customization (--storage-size, --storage-tier) is only supported for Cloud Native (CLOUDNATIVE-*) and GPU (GPU-*) plans, got plan: 2xCPU-4GB",
		},
		{
			name:     "invalid storage tier",
			args:     []string{clusterUUID, "--count", "2", "--name", "gpu-nodes", "--plan", "GPU-8xCPU-64GB-1xL40S", "--storage-tier", "invalid"},
			errorMsg: "invalid storage tier \"invalid\", must be one of: maxiops, standard, hdd",
		},
		{
			name:     "invalid storage size too small",
			args:     []string{clusterUUID, "--count", "2", "--name", "gpu-nodes", "--plan", "GPU-8xCPU-64GB-1xL40S", "--storage-size", "20"},
			errorMsg: "storage size must be at least 25 GiB, got 20",
		},
		{
			name:     "invalid storage size too large",
			args:     []string{clusterUUID, "--count", "2", "--name", "gpu-nodes", "--plan", "GPU-8xCPU-64GB-1xL40S", "--storage-size", "5000"},
			errorMsg: "storage size cannot exceed 4096 GiB, got 5000",
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			conf := config.New()
			testCmd := CreateCommand()
			mService := new(smock.Service)

			expected := test.expected
			mService.On("CreateKubernetesNodeGroup", &expected).Return(&upcloud.KubernetesNodeGroup{}, nil)
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

func TestSupportStorageCustomization(t *testing.T) {
	tests := []struct {
		planName string
		expected bool
	}{
		{"2xCPU-4GB", false},
		{"4xCPU-8GB", false},
		{"HICPU-8xCPU-16GB", false},
		{"HIMEM-4xCPU-32GB", false},
		{"DEV-1xCPU-1GB", false},
		{"CLOUDNATIVE-2xCPU-4GB", true},
		{"CLOUDNATIVE-4xCPU-8GB", true},
		{"GPU-8xCPU-64GB-1xL40S", true},
		{"", false},
	}

	for _, test := range tests {
		t.Run(test.planName, func(t *testing.T) {
			result := supportStorageCustomization(test.planName)
			assert.Equal(t, test.expected, result)
		})
	}
}

func TestValidateStorageTier(t *testing.T) {
	tests := []struct {
		tier     string
		hasError bool
	}{
		{"", false},
		{"maxiops", false},
		{"standard", false},
		{"hdd", false},
		{"archive", true},
		{"invalid", true},
		{"MAXIOPS", true},
		{"Standard", true},
	}

	for _, test := range tests {
		t.Run(test.tier, func(t *testing.T) {
			err := validateStorageTier(test.tier)
			if test.hasError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
