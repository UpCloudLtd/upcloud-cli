package server

import (
	"github.com/UpCloudLtd/upcloud-go-api/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/request"
	"github.com/stretchr/testify/mock"
)

// MockServerService implements a mock server service
// TODO: this should live somewhere else than in the server command util.go..
type MockServerService struct {
	mock.Mock
}

// GetServerConfigurations implements service.Server.GetServerConfigurations
func (m *MockServerService) GetServerConfigurations() (*upcloud.ServerConfigurations, error) {
	args := m.Called()
	return args[0].(*upcloud.ServerConfigurations), args.Error(1)
}

// GetServers implements service.Server.GetServers
func (m *MockServerService) GetServers() (*upcloud.Servers, error) {
	args := m.Called()
	return args[0].(*upcloud.Servers), args.Error(1)
}

// GetServerDetails implements service.Server.GetServerDetails
func (m *MockServerService) GetServerDetails(r *request.GetServerDetailsRequest) (*upcloud.ServerDetails, error) {
	args := m.Called(r)
	return args[0].(*upcloud.ServerDetails), args.Error(1)
}

// CreateServer implements service.Server.CreateServer
func (m *MockServerService) CreateServer(r *request.CreateServerRequest) (*upcloud.ServerDetails, error) {
	args := m.Called(r)
	return args[0].(*upcloud.ServerDetails), args.Error(1)
}

// WaitForServerState implements service.Server.WaitForServerState
func (m *MockServerService) WaitForServerState(r *request.WaitForServerStateRequest) (*upcloud.ServerDetails, error) {
	args := m.Called(r)
	return args[0].(*upcloud.ServerDetails), args.Error(1)
}

// StartServer implements service.Server.StartServer
func (m *MockServerService) StartServer(r *request.StartServerRequest) (*upcloud.ServerDetails, error) {
	args := m.Called(r)
	return args[0].(*upcloud.ServerDetails), args.Error(1)
}

// StopServer implements service.Server.StopServer
func (m *MockServerService) StopServer(r *request.StopServerRequest) (*upcloud.ServerDetails, error) {
	args := m.Called(r)
	return args[0].(*upcloud.ServerDetails), args.Error(1)
}

// RestartServer implements service.Server.RestartServer
func (m *MockServerService) RestartServer(r *request.RestartServerRequest) (*upcloud.ServerDetails, error) {
	args := m.Called(r)
	return args[0].(*upcloud.ServerDetails), args.Error(1)
}

// ModifyServer implements service.Server.ModifyServer
func (m *MockServerService) ModifyServer(r *request.ModifyServerRequest) (*upcloud.ServerDetails, error) {
	args := m.Called(r)
	return args[0].(*upcloud.ServerDetails), args.Error(1)
}

// DeleteServer implements service.Server.DeleteServer
func (m *MockServerService) DeleteServer(r *request.DeleteServerRequest) error {
	args := m.Called(r)
	return args.Error(0)
}

// DeleteServerAndStorages implements service.Server.DeleteServerAndStorages
func (m *MockServerService) DeleteServerAndStorages(r *request.DeleteServerAndStoragesRequest) error {
	args := m.Called(r)
	return args.Error(0)
}
