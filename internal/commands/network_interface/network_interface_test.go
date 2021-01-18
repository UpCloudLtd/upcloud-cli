package network_interface

import (
	"github.com/UpCloudLtd/upcloud-go-api/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/request"
	"github.com/stretchr/testify/mock"
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

func (m *MockNetworkService) GetRouters() (*upcloud.Routers, error) {
	args := m.Called()
	return args[0].(*upcloud.Routers), args.Error(1)
}
func (m *MockNetworkService) GetRouterDetails(r *request.GetRouterDetailsRequest) (*upcloud.Router, error) {
	args := m.Called(r)
	return args[0].(*upcloud.Router), args.Error(1)
}
func (m *MockNetworkService) CreateRouter(r *request.CreateRouterRequest) (*upcloud.Router, error) {
	args := m.Called(r)
	return args[0].(*upcloud.Router), args.Error(1)
}
func (m *MockNetworkService) ModifyRouter(r *request.ModifyRouterRequest) (*upcloud.Router, error) {
	args := m.Called(r)
	return args[0].(*upcloud.Router), args.Error(1)
}
func (m *MockNetworkService) DeleteRouter(r *request.DeleteRouterRequest) error {
	args := m.Called(r)
	return args.Error(0)
}
