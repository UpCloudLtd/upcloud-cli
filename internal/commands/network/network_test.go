package network

import (
	"github.com/UpCloudLtd/upcloud-go-api/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/request"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

type MockNetworkService struct {
	mock.Mock
}

func (m *MockNetworkService) GetNetworks() (*upcloud.Networks, error) {
	args := m.Called()
	return args[0].(*upcloud.Networks), args.Error(1)
}
func (m *MockNetworkService) GetNetworksInZone(r *request.GetNetworksInZoneRequest) (*upcloud.Networks, error) {
	args := m.Called(r)
	return args[0].(*upcloud.Networks), args.Error(1)
}
func (m *MockNetworkService) CreateNetwork(r *request.CreateNetworkRequest) (*upcloud.Network, error) {
	args := m.Called(r)
	return args[0].(*upcloud.Network), args.Error(1)
}
func (m *MockNetworkService) GetNetworkDetails(r *request.GetNetworkDetailsRequest) (*upcloud.Network, error) {
	args := m.Called(r)
	return args[0].(*upcloud.Network), args.Error(1)
}
func (m *MockNetworkService) ModifyNetwork(r *request.ModifyNetworkRequest) (*upcloud.Network, error) {
	args := m.Called(r)
	return args[0].(*upcloud.Network), args.Error(1)
}
func (m *MockNetworkService) GetServerNetworks(r *request.GetServerNetworksRequest) (*upcloud.Networking, error) {
	args := m.Called(r)
	return args[0].(*upcloud.Networking), args.Error(1)
}
func (m *MockNetworkService) CreateNetworkInterface(r *request.CreateNetworkInterfaceRequest) (*upcloud.Interface, error) {
	args := m.Called(r)
	return args[0].(*upcloud.Interface), args.Error(1)
}
func (m *MockNetworkService) ModifyNetworkInterface(r *request.ModifyNetworkInterfaceRequest) (*upcloud.Interface, error) {
	args := m.Called(r)
	return args[0].(*upcloud.Interface), args.Error(1)
}
func (m *MockNetworkService) DeleteNetwork(r *request.DeleteNetworkRequest) error {
	args := m.Called(r)
	return args.Error(0)
}
func (m *MockNetworkService) DeleteNetworkInterface(r *request.DeleteNetworkInterfaceRequest) error {
	args := m.Called(r)
	return args.Error(0)
}

type MockServerService struct {
	mock.Mock
}

func (m *MockServerService) GetServerConfigurations() (*upcloud.ServerConfigurations, error) {
	args := m.Called()
	return args[0].(*upcloud.ServerConfigurations), args.Error(1)
}
func (m *MockServerService) GetServers() (*upcloud.Servers, error) {
	args := m.Called()
	return args[0].(*upcloud.Servers), args.Error(1)
}
func (m *MockServerService) GetServerDetails(r *request.GetServerDetailsRequest) (*upcloud.ServerDetails, error) {
	args := m.Called(r)
	return args[0].(*upcloud.ServerDetails), args.Error(1)
}
func (m *MockServerService) CreateServer(r *request.CreateServerRequest) (*upcloud.ServerDetails, error) {
	args := m.Called(r)
	return args[0].(*upcloud.ServerDetails), args.Error(1)
}
func (m *MockServerService) WaitForServerState(r *request.WaitForServerStateRequest) (*upcloud.ServerDetails, error) {
	args := m.Called(r)
	return args[0].(*upcloud.ServerDetails), args.Error(1)
}
func (m *MockServerService) StartServer(r *request.StartServerRequest) (*upcloud.ServerDetails, error) {
	args := m.Called(r)
	return args[0].(*upcloud.ServerDetails), args.Error(1)
}
func (m *MockServerService) StopServer(r *request.StopServerRequest) (*upcloud.ServerDetails, error) {
	args := m.Called(r)
	return args[0].(*upcloud.ServerDetails), args.Error(1)
}
func (m *MockServerService) RestartServer(r *request.RestartServerRequest) (*upcloud.ServerDetails, error) {
	args := m.Called(r)
	return args[0].(*upcloud.ServerDetails), args.Error(1)
}
func (m *MockServerService) ModifyServer(r *request.ModifyServerRequest) (*upcloud.ServerDetails, error) {
	args := m.Called(r)
	return args[0].(*upcloud.ServerDetails), args.Error(1)
}
func (m *MockServerService) DeleteServer(r *request.DeleteServerRequest) error {
	args := m.Called(r)
	return args.Error(0)
}
func (m *MockServerService) DeleteServerAndStorages(r *request.DeleteServerAndStoragesRequest) error {
	args := m.Called(r)
	return args.Error(0)
}

func TestSearchAllNetworks(t *testing.T) {

	var Network1 = upcloud.Network{
		Name: "network-1",
		UUID: "28e15cf5-8817-42ab-b017-970666be96ec",
	}

	var Network2 = upcloud.Network{
		Name: "network-2",
		UUID: "f9f5ad16-a63a-4670-8449-c01d1e97281e",
	}

	var Network3 = upcloud.Network{
		Name: "network-3",
		UUID: "e157ce0a-eeb0-49fc-9f2c-a05c3ac57066",
	}

	var Network4 = upcloud.Network{
		Name: Network1.Name,
		UUID: "b3e49768-f13a-42c3-bea7-4e2471657f2f",
	}

	var networks = &upcloud.Networks{Networks: []upcloud.Network{Network1, Network2, Network3, Network4}}

	for _, testcase := range []struct {
		name         string
		args         []string
		expected     []string
		unique       bool
		additional   []upcloud.Storage
		backendCalls int
		errMsg       string
	}{
		{
			name:         "SingleUuid",
			args:         []string{Network2.UUID},
			expected:     []string{Network2.UUID},
			backendCalls: 0,
		},
		{
			name:         "MultipleUuidSearched",
			args:         []string{Network2.UUID, Network3.UUID},
			expected:     []string{Network2.UUID, Network3.UUID},
			backendCalls: 0,
		},
		{
			name:         "SingleName",
			args:         []string{Network2.Name},
			expected:     []string{Network2.UUID},
			backendCalls: 1,
		},
		{
			name:         "MultipleNamesSearched",
			args:         []string{Network2.Name, Network3.Name},
			expected:     []string{Network2.UUID, Network3.UUID},
			backendCalls: 1,
		},
		{
			name:         "MultipleNamesFound",
			args:         []string{Network1.Name},
			expected:     []string{Network1.UUID, Network4.UUID},
			backendCalls: 1,
		},
		{
			name:         "MultipleNamesFound_UniqueWanted",
			args:         []string{Network1.Name},
			expected:     []string{Network1.UUID, Network4.UUID},
			backendCalls: 1,
			unique:       true,
			errMsg:       "multiple networks matched to query \"" + Network1.Name + "\", use UUID to specify",
		},
	} {
		t.Run(testcase.name, func(t *testing.T) {
			cachedNetworks = nil
			mns := MockNetworkService{}
			mns.On("GetNetworks", mock.Anything).Return(networks, nil)

			result, err := searchAllNetworks(testcase.args, &mns, testcase.unique)

			if testcase.errMsg == "" {
				assert.Nil(t, err)
				assert.ElementsMatch(t, testcase.expected, result)
			} else {
				assert.Nil(t, result)
				assert.EqualError(t, err, testcase.errMsg)
			}
			mns.AssertNumberOfCalls(t, "GetNetworks", testcase.backendCalls)
		})
	}
}

func TestSearchSUniqueNetwork(t *testing.T) {

	var Network1 = upcloud.Network{
		Name: "network-1",
		UUID: "28e15cf5-8817-42ab-b017-970666be96ec",
	}

	var Network2 = upcloud.Network{
		Name: "network-2",
		UUID: "f9f5ad16-a63a-4670-8449-c01d1e97281e",
	}

	var Network3 = upcloud.Network{
		Name: "network-3",
		UUID: "e157ce0a-eeb0-49fc-9f2c-a05c3ac57066",
	}

	var Network4 = upcloud.Network{
		Name: Network1.Name,
		UUID: "b3e49768-f13a-42c3-bea7-4e2471657f2f",
	}

	var networks = &upcloud.Networks{Networks: []upcloud.Network{Network1, Network2, Network3, Network4}}

	for _, testcase := range []struct {
		name         string
		args         string
		expected     *upcloud.Network
		backendCalls int
		errMsg       string
	}{
		{
			name:         "SingleUuid",
			args:         Network2.UUID,
			expected:     &Network2,
			backendCalls: 1,
		},
		{
			name:         "SingleName",
			args:         Network2.Name,
			expected:     &Network2,
			backendCalls: 1,
		},
		{
			name:         "MultipleNamesFound_UniqueWanted",
			args:         Network1.Name,
			backendCalls: 1,
			errMsg:       "multiple networks matched to query \"" + Network1.Name + "\", use UUID to specify",
		},
	} {
		t.Run(testcase.name, func(t *testing.T) {
			cachedNetworks = nil
			mns := MockNetworkService{}
			mns.On("GetNetworks", mock.Anything).Return(networks, nil)

			result, err := SearchUniqueNetwork(testcase.args, &mns)

			if testcase.errMsg == "" {
				assert.Nil(t, err)
				assert.Equal(t, testcase.expected, result)
			} else {
				assert.Nil(t, result)
				assert.EqualError(t, err, testcase.errMsg)
			}
			mns.AssertNumberOfCalls(t, "GetNetworks", testcase.backendCalls)
		})
	}
}
