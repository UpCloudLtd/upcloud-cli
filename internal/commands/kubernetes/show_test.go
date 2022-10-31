package kubernetes

import (
	"testing"

	"github.com/UpCloudLtd/upcloud-cli/v2/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v2/internal/config"
	smock "github.com/UpCloudLtd/upcloud-cli/v2/internal/mock"
	"github.com/UpCloudLtd/upcloud-cli/v2/internal/mockexecute"

	"github.com/UpCloudLtd/upcloud-go-api/v4/upcloud"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestShowCommand(t *testing.T) {
	text.DisableColors()
	cluster1 := upcloud.KubernetesCluster{
		Name:        "upcloud-go-sdk-unit-test",
		Network:     "03a98be3-7daa-443f-bb25-4bc6854b396c",
		NetworkCIDR: "172.16.0.0/24",
		NodeGroups: []upcloud.KubernetesNodeGroup{
			{
				Count: 1,
				Labels: []upcloud.Label{
					{
						Key:   "managedBy",
						Value: "upcloud-go-sdk-unit-test",
					},
					{
						Key:   "another",
						Value: "label-thing",
					},
				},
				Name: "upcloud-go-sdk-unit-test",
				Plan: "K8S-2xCPU-4GB",
				KubeletArgs: []upcloud.KubernetesKubeletArg{
					{
						Key:   "somekubeletkey",
						Value: "somekubeletvalue",
					},
				},
				Taints: []upcloud.KubernetesTaint{
					{
						Effect: "NoExecute",
						Key:    "sometaintkey",
						Value:  "sometaintvalue",
					},
					{
						Effect: "NoExecute",
						Key:    "sometaintkey",
						Value:  "sometaintvalue",
					},
					{
						Effect: "NoExecute",
						Key:    "sometaintkey",
						Value:  "sometaintvalue",
					},
				},
				SSHKeys: []string{"somekey"},
			},
		},
		State: upcloud.KubernetesClusterStateRunning,
		UUID:  "0ddab8f4-97c0-4222-91ba-85a4fff7499b",
		Zone:  "de-fra1",
	}

	expected := `
    
`

	mService := smock.Service{}
	mService.On("GetKubernetesClusters", mock.Anything).Return([]upcloud.KubernetesCluster{cluster1}, nil)
	mService.On("GetKubernetesCluster", mock.Anything).Return(&cluster1, nil)
	mService.On("GetNetworkDetails", mock.Anything).Return(&upcloud.Network{Name: "Test network"}, nil)

	conf := config.New()
	command := commands.BuildCommand(ShowCommand(), nil, conf)

	// get resolver to initialize command cache
	_, err := command.(*showCommand).Get(&mService)
	if err != nil {
		t.Fatal(err)
	}

	command.Cobra().SetArgs([]string{cluster1.UUID})
	output, err := mockexecute.MockExecute(command, &mService, conf)

	assert.NoError(t, err)
	assert.Equal(t, expected, output)
}
