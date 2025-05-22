package all

import (
	"testing"

	"github.com/jedib0t/go-pretty/v6/text"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/config"
	smock "github.com/UpCloudLtd/upcloud-cli/v3/internal/mock"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/mockexecute"
	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func mockListResponses(mService *smock.Service) {
	mService.On("GetManagedDatabases", mock.Anything).Return(nil, nil)
	mService.On("GetManagedObjectStorages", mock.Anything).Return(objectStorages, nil)
	mService.On("GetNetworks").Return(networks, nil)
	mService.On("GetRouters").Return(&upcloud.Routers{}, nil)
	mService.On("GetServers").Return(&upcloud.Servers{}, nil)
	mService.On("GetServerGroups", mock.Anything).Return(nil, nil)
	mService.On("GetStorages", mock.Anything).Return(&upcloud.Storages{}, nil)
	mService.On("GetTags").Return(&upcloud.Tags{}, nil)
	mService.On("GetKubernetesClusters", mock.Anything).Return(nil, nil)
	mService.On("GetLoadBalancers", mock.Anything).Return(nil, nil)
}

var networks = &upcloud.Networks{Networks: []upcloud.Network{
	{
		Name: "tf-acc-test-network",
		UUID: "28e15cf5-8817-42ab-b017-970666be96ec",
		Type: upcloud.NetworkTypePrivate,
		Zone: "pl-waw1",
	},
	{
		Name: "uks-e2e-test-network",
		UUID: "f9f5ad16-a63a-4670-8449-c01d1e97281e",
		Type: upcloud.NetworkTypePrivate,
		Zone: "fi-hel1",
	},
}}

var objectStorages = []upcloud.ManagedObjectStorage{
	{
		Name:   "tf-acc-test-objsto",
		UUID:   "28e15cf5-8817-42ab-b017-970666be96ec",
		Region: "europe-1",
	},
	{
		Name:   "persistent-tf-acc-test-objsto",
		UUID:   "f9f5ad16-a63a-4670-8449-c01d1e97281e",
		Region: "apac-1",
	},
}

func TestListCommand(t *testing.T) {
	for _, test := range []struct {
		name     string
		args     []string
		expected string
	}{
		{
			name: "list all",
			expected: `
 Type             UUID                                   Name                          
──────────────── ────────────────────────────────────── ───────────────────────────────
 network          28e15cf5-8817-42ab-b017-970666be96ec   tf-acc-test-network           
 network          f9f5ad16-a63a-4670-8449-c01d1e97281e   uks-e2e-test-network          
 object-storage   f9f5ad16-a63a-4670-8449-c01d1e97281e   persistent-tf-acc-test-objsto 
 object-storage   28e15cf5-8817-42ab-b017-970666be96ec   tf-acc-test-objsto            

`,
		},
		{
			name: "list non-persistent tf-acc-test resources",
			args: []string{"--include", "*tf-acc-test*", "--exclude", "*persistent*"},
			expected: `
 Type             UUID                                   Name                
──────────────── ────────────────────────────────────── ─────────────────────
 network          28e15cf5-8817-42ab-b017-970666be96ec   tf-acc-test-network 
 object-storage   28e15cf5-8817-42ab-b017-970666be96ec   tf-acc-test-objsto  

`,
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			text.DisableColors()

			mService := smock.Service{}
			mockListResponses(&mService)

			conf := config.New()
			command := commands.BuildCommand(ListCommand(), nil, conf)

			command.Cobra().SetArgs(test.args)
			output, err := mockexecute.MockExecute(command, &mService, conf)

			assert.NoError(t, err)
			assert.Equal(t, test.expected, output)
		})
	}
}
