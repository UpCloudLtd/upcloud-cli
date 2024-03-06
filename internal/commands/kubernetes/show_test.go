package kubernetes

import (
	"context"
	"testing"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/config"
	smock "github.com/UpCloudLtd/upcloud-cli/v3/internal/mock"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/mockexecute"

	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var testCluster = upcloud.KubernetesCluster{
	ControlPlaneIPFilter: []string{"10.144.1.100", "10.144.2.0/24"},
	Labels: []upcloud.Label{
		{Key: "test", Value: "upctl-unittest"},
	},
	Name:        "upcloud-upctl-unit-test",
	Network:     "03a98be3-7daa-443f-bb25-4bc6854b396c",
	NetworkCIDR: "172.16.1.0/24",
	NodeGroups: []upcloud.KubernetesNodeGroup{
		{
			Count: 4,
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
			Name:  "upcloud-go-sdk-unit-test",
			Plan:  "2xCPU-4GB",
			State: upcloud.KubernetesNodeGroupStateRunning,
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
			Storage:              "storage-uuid",
			SSHKeys:              []string{"somekey"},
			UtilityNetworkAccess: true,
		}, {
			Count: 8,
			Labels: []upcloud.Label{
				{
					Key:   "managedBy",
					Value: "upcloud-go-sdk-unit-test-2",
				},
				{
					Key:   "another2",
					Value: "label-thing-2",
				},
			},
			Name:  "upcloud-go-sdk-unit-test-2",
			Plan:  "4xCPU-8GB",
			State: upcloud.KubernetesNodeGroupStatePending,
			KubeletArgs: []upcloud.KubernetesKubeletArg{
				{
					Key:   "somekubeletkey2",
					Value: "somekubeletvalue2",
				},
			},
			Taints: []upcloud.KubernetesTaint{
				{
					Effect: "NoSchedule",
					Key:    "sometaintkey2",
					Value:  "sometaintvalue2",
				},
			},
			Storage:              "storage-uuid-2",
			SSHKeys:              []string{"somekey2"},
			UtilityNetworkAccess: false,
		},
	},
	State:   upcloud.KubernetesClusterStateRunning,
	UUID:    "0ddab8f4-97c0-4222-91ba-85a4fff7499b",
	Version: "2.54",
	Zone:    "de-fra1",
}

func TestShowCommand(t *testing.T) {
	text.DisableColors()

	expected := `  
  Overview:
    UUID:                       0ddab8f4-97c0-4222-91ba-85a4fff7499b 
    Name:                       upcloud-upctl-unit-test              
    Version:                    2.54                                 
    Network UUID:               03a98be3-7daa-443f-bb25-4bc6854b396c 
    Network name:               Test network                         
    Network CIDR:               172.16.1.0/24                        
    Kubernetes API allowed IPs: 10.144.1.100,                        
                                10.144.2.0/24                        
    Private node groups:        no                                   
    Zone:                       de-fra1                              
    Operational state:          running                              

  Labels:

     Key    Value          
    ────── ────────────────
     test   upctl-unittest 
    
  Node groups:

     Name                         Count   Plan        Anti affinity   Utility network access   State   
    ──────────────────────────── ─────── ─────────── ─────────────── ──────────────────────── ─────────
     upcloud-go-sdk-unit-test         4   2xCPU-4GB   no              yes                      running 
     upcloud-go-sdk-unit-test-2       8   4xCPU-8GB   no              no                       pending 
    
`

	mService := smock.Service{}
	mService.On("GetKubernetesClusters", mock.Anything).Return([]upcloud.KubernetesCluster{testCluster}, nil)
	mService.On("GetKubernetesCluster", mock.Anything).Return(&testCluster, nil)
	mService.On("GetNetworkDetails", mock.Anything).Return(&upcloud.Network{Name: "Test network"}, nil)
	mService.On("GetStorageDetails", mock.Anything).Return(&upcloud.StorageDetails{Storage: upcloud.Storage{Title: "Test storage"}}, nil)

	conf := config.New()
	command := commands.BuildCommand(ShowCommand(), nil, conf)

	// get resolver to initialize command cache
	_, err := command.(*showCommand).Get(context.TODO(), &mService)
	if err != nil {
		t.Fatal(err)
	}

	command.Cobra().SetArgs([]string{testCluster.UUID})
	output, err := mockexecute.MockExecute(command, &mService, conf)

	assert.NoError(t, err)
	assert.Equal(t, expected, output)
}
