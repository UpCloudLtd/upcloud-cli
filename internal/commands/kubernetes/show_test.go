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
				Storage: "storage-uuid",
				SSHKeys: []string{"somekey"},
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
				Name: "upcloud-go-sdk-unit-test-2",
				Plan: "K8S-4xCPU-8GB",
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
				Storage: "storage-uuid-2",
				SSHKeys: []string{"somekey2"},
			},
		},
		State: upcloud.KubernetesClusterStateRunning,
		UUID:  "0ddab8f4-97c0-4222-91ba-85a4fff7499b",
		Zone:  "de-fra1",
	}

	expected := `  
  Overview:
    UUID:              0ddab8f4-97c0-4222-91ba-85a4fff7499b 
    Name:              upcloud-go-sdk-unit-test             
    Network UUID:      03a98be3-7daa-443f-bb25-4bc6854b396c 
    Network name       Test network                         
    Network CIDR:      172.16.0.0/24                        
    Zone               de-fra1                              
    Operational state: running                              

  Node groups:

     Name                         Count   Plan            Storage          Kubelet args               Labels                              Taints                 
    ──────────────────────────── ─────── ─────────────── ──────────────── ────────────────────────── ─────────────────────────────────── ────────────────────────
     upcloud-go-sdk-unit-test         4   K8S-2xCPU-4GB   storage-uuid     Key: somekubeletkey        Key: managedBy                      Key: sometaintkey      
                                                                           Value: somekubeletvalue    Value: upcloud-go-sdk-unit-test     Value: sometaintvalue  
                                                                                                                                          Effect: NoExecute      
                                                                                                      Key: another                                               
                                                                                                      Value: label-thing                  Key: sometaintkey      
                                                                                                                                          Value: sometaintvalue  
                                                                                                                                          Effect: NoExecute      
                                                                                                                                                                 
                                                                                                                                          Key: sometaintkey      
                                                                                                                                          Value: sometaintvalue  
                                                                                                                                          Effect: NoExecute      
     upcloud-go-sdk-unit-test-2       8   K8S-4xCPU-8GB   storage-uuid-2   Key: somekubeletkey2       Key: managedBy                      Key: sometaintkey2     
                                                                           Value: somekubeletvalue2   Value: upcloud-go-sdk-unit-test-2   Value: sometaintvalue2 
                                                                                                                                          Effect: NoSchedule     
                                                                                                      Key: another2                                              
                                                                                                      Value: label-thing-2                                       
    
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
