package mocks

import (
	"github.com/UpCloudLtd/upcloud-go-api/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/request"
	"github.com/stretchr/testify/mock"
)

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
