package kubernetes

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

func TestModifyKubernetesCluster(t *testing.T) {
	clusterUUID := "28c80353-98fd-4221-85e0-82d7603756ba"

	for _, test := range []struct {
		name     string
		args     []string
		expected request.ModifyKubernetesClusterRequest
		errorMsg string
	}{
		{
			name:     "no args",
			args:     []string{clusterUUID},
			expected: request.ModifyKubernetesClusterRequest{ClusterUUID: clusterUUID, Cluster: request.ModifyKubernetesCluster{}},
		},
		{
			name: "one IP",
			args: []string{
				clusterUUID,
				"--kubernetes-api-allow-ip", "10.144.1.100",
			},
			expected: request.ModifyKubernetesClusterRequest{
				ClusterUUID: clusterUUID,
				Cluster: request.ModifyKubernetesCluster{
					ControlPlaneIPFilter: &[]string{"10.144.1.100"},
					Labels:               nil,
				},
			},
		},
		{
			name: "IP and CIDR",
			args: []string{
				clusterUUID,
				"--kubernetes-api-allow-ip", "10.144.1.100",
				"--kubernetes-api-allow-ip", "10.144.2.0/24",
			},
			expected: request.ModifyKubernetesClusterRequest{
				ClusterUUID: clusterUUID,
				Cluster: request.ModifyKubernetesCluster{
					ControlPlaneIPFilter: &[]string{"10.144.1.100", "10.144.2.0/24"},
					Labels:               nil,
				},
			},
		},
		{
			name: "labels",
			args: []string{
				clusterUUID,
				"--label", "tool=upctl",
				"--label", "test=unittest",
			},
			expected: request.ModifyKubernetesClusterRequest{
				ClusterUUID: clusterUUID,
				Cluster: request.ModifyKubernetesCluster{
					ControlPlaneIPFilter: nil,
					Labels: &[]upcloud.Label{
						{Key: "tool", Value: "upctl"},
						{Key: "test", Value: "unittest"},
					},
				},
			},
		},
		{
			name: "clear-labels",
			args: []string{
				clusterUUID,
				"--clear-labels",
			},
			expected: request.ModifyKubernetesClusterRequest{
				ClusterUUID: clusterUUID,
				Cluster: request.ModifyKubernetesCluster{
					ControlPlaneIPFilter: nil,
					Labels:               &[]upcloud.Label{},
				},
			},
		},
		{
			name: "labels and clear-labels",
			args: []string{
				clusterUUID,
				"--label", "tool=upctl",
				"--label", "test=unittest",
				"--clear-labels",
			},
			errorMsg: "if any flags in the group [label clear-labels] are set none of the others can be; [clear-labels label] were all set",
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			conf := config.New()
			testCmd := ModifyCommand()
			mService := new(smock.Service)

			expected := test.expected
			mService.On("ModifyKubernetesCluster", &expected).Return(&upcloud.KubernetesCluster{}, nil)

			c := commands.BuildCommand(testCmd, nil, conf)

			c.Cobra().SetArgs(test.args)
			_, err := mockexecute.MockExecute(c, mService, conf)

			if test.errorMsg != "" {
				assert.EqualError(t, err, test.errorMsg)
			} else {
				require.NoError(t, err)
				mService.AssertNumberOfCalls(t, "ModifyKubernetesCluster", 1)
			}
		})
	}
}
