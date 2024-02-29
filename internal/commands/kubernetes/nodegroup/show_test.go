package nodegroup

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

func TestShowCommand(t *testing.T) {
	text.DisableColors()
	smallNodeGroup := upcloud.KubernetesNodeGroupDetails{
		KubernetesNodeGroup: upcloud.KubernetesNodeGroup{
			AntiAffinity: true,
			Count:        4,
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
			Name:  "small",
			Plan:  "1xCPU-1GB",
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
					Key:    "taintkey1",
					Value:  "taintvalue1",
				},
				{
					Effect: "NoSchedule",
					Key:    "taintkey2",
					Value:  "taintvalue2",
				},
			},
			Storage:              "storage-uuid",
			SSHKeys:              []string{"somekey"},
			UtilityNetworkAccess: true,
		},
		Nodes: []upcloud.KubernetesNode{
			{
				UUID:  "00b1fc81-d471-4ae0-90ea-bfef9dd326bf",
				Name:  "small-9zhrq",
				State: upcloud.KubernetesNodeStateRunning,
			},
			{
				UUID:  "009fc510-f65a-4c3f-9917-d7cf0bb3f15f",
				Name:  "small-nbqvp",
				State: upcloud.KubernetesNodeStateRunning,
			},
			{
				UUID:  "0039e13a-76a6-4253-90d9-54be10539016",
				Name:  "small-wzdxz",
				State: upcloud.KubernetesNodeStateRunning,
			},
			{
				UUID:  "00141c93-afee-4c43-be28-79792675eec2",
				Name:  "small-xnzvc",
				State: upcloud.KubernetesNodeStateRunning,
			},
		},
	}

	smallExpected := `  
  Overview
    Name:                   small        
    Count:                  4            
    Plan:                   1xCPU-1GB    
    State:                  running      
    Storage UUID:           storage-uuid 
    Storage name:           Test storage 
    Anti-affinity:          yes          
    Utility network access: yes          

  Labels:

     Key         Value                    
    ─────────── ──────────────────────────
     managedBy   upcloud-go-sdk-unit-test 
     another     label-thing              
    
  Taints:

     Key         Value         Effect     
    ─────────── ───────────── ────────────
     taintkey1   taintvalue1   NoExecute  
     taintkey2   taintvalue2   NoSchedule 
    
  Kubelet arguments:

     Key              Value            
    ──────────────── ──────────────────
     somekubeletkey   somekubeletvalue 
    
  Nodes:

     UUID                                   Name          State   
    ────────────────────────────────────── ───────────── ─────────
     00b1fc81-d471-4ae0-90ea-bfef9dd326bf   small-9zhrq   running 
     009fc510-f65a-4c3f-9917-d7cf0bb3f15f   small-nbqvp   running 
     0039e13a-76a6-4253-90d9-54be10539016   small-wzdxz   running 
     00141c93-afee-4c43-be28-79792675eec2   small-xnzvc   running 
    
`

	emptyNodeGroup := upcloud.KubernetesNodeGroupDetails{
		KubernetesNodeGroup: upcloud.KubernetesNodeGroup{
			AntiAffinity:         false,
			Count:                0,
			Name:                 "empty",
			Plan:                 "1xCPU-1GB",
			State:                upcloud.KubernetesNodeGroupStateRunning,
			Storage:              "storage-uuid",
			UtilityNetworkAccess: true,
		},
	}

	emptyExpected := `  
  Overview
    Name:                   empty        
    Count:                  0            
    Plan:                   1xCPU-1GB    
    State:                  running      
    Storage UUID:           storage-uuid 
    Storage name:           Test storage 
    Anti-affinity:          no           
    Utility network access: yes          

  Labels:

    No labels defined for this node group.
    
  Taints:

    No taints defined for this node group.
    
  Kubelet arguments:

    No kubelet arguments defined for this node group.
    
  Nodes:

    No nodes found for this node group.
    
`

	for _, test := range []struct {
		name      string
		nodeGroup *upcloud.KubernetesNodeGroupDetails
		expected  string
	}{
		{
			name:      "data in all tables",
			nodeGroup: &smallNodeGroup,
			expected:  smallExpected,
		},
		{
			name:      "empty states",
			nodeGroup: &emptyNodeGroup,
			expected:  emptyExpected,
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			mService := smock.Service{}
			mService.On("GetKubernetesClusters", mock.Anything).Return([]upcloud.KubernetesCluster{{Name: "mock-cluster", UUID: "fake-uuid"}}, nil)
			mService.On("GetKubernetesNodeGroup", mock.Anything).Return(test.nodeGroup, nil)
			mService.On("GetStorageDetails", mock.Anything).Return(&upcloud.StorageDetails{Storage: upcloud.Storage{Title: "Test storage"}}, nil)

			conf := config.New()
			command := commands.BuildCommand(ShowCommand(), nil, conf)

			// get resolver to initialize command cache
			_, err := command.(*showCommand).Get(context.TODO(), &mService)
			if err != nil {
				t.Fatal(err)
			}

			command.Cobra().SetArgs([]string{"cluster-name", "--name", test.nodeGroup.Name})
			output, err := mockexecute.MockExecute(command, &mService, conf)

			assert.NoError(t, err)
			assert.Equal(t, test.expected, output)
		})
	}
}
