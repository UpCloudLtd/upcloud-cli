package kubernetes

import (
	"testing"

	"github.com/UpCloudLtd/upcloud-cli/v2/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v2/internal/config"
	smock "github.com/UpCloudLtd/upcloud-cli/v2/internal/mock"
	"github.com/UpCloudLtd/upcloud-cli/v2/internal/mockexecute"

	"github.com/UpCloudLtd/upcloud-go-api/v6/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/v6/upcloud/request"
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
			expected: request.ModifyKubernetesClusterRequest{ClusterUUID: clusterUUID, Cluster: request.ModifyKubernetesCluster{ControlPlaneIPFilter: []string{}}},
		},
		{
			name: "one IP",
			args: []string{
				clusterUUID,
				"--kubernetes-api-allow-ip", "10.144.1.100",
			},
			expected: request.ModifyKubernetesClusterRequest{
				ClusterUUID: clusterUUID,
				Cluster:     request.ModifyKubernetesCluster{ControlPlaneIPFilter: []string{"10.144.1.100"}},
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
				Cluster:     request.ModifyKubernetesCluster{ControlPlaneIPFilter: []string{"10.144.1.100", "10.144.2.0/24"}},
			},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			conf := config.New()
			testCmd := ModifyCommand()
			mService := new(smock.Service)

			mService.On("ModifyKubernetesCluster", &test.expected).Return(&upcloud.KubernetesCluster{}, nil)

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
