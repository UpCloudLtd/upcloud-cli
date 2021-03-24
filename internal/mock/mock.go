package mock

import (
	"github.com/UpCloudLtd/upcloud-go-api/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/request"
	"github.com/stretchr/testify/mock"
)

type MockService struct {
	mock.Mock
}

func (m *MockService) GetServerConfigurations() (*upcloud.ServerConfigurations, error) {
	args := m.Called()
	return args[0].(*upcloud.ServerConfigurations), args.Error(1)
}

func (m *MockService) GetServers() (*upcloud.Servers, error) {
	args := m.Called()
	return args[0].(*upcloud.Servers), args.Error(1)
}

func (m *MockService) GetServerDetails(r *request.GetServerDetailsRequest) (*upcloud.ServerDetails, error) {
	args := m.Called(r)
	return args[0].(*upcloud.ServerDetails), args.Error(1)
}

func (m *MockService) CreateServer(r *request.CreateServerRequest) (*upcloud.ServerDetails, error) {
	args := m.Called(r)
	return args[0].(*upcloud.ServerDetails), args.Error(1)
}

func (m *MockService) WaitForServerState(r *request.WaitForServerStateRequest) (*upcloud.ServerDetails, error) {
	args := m.Called(r)
	return args[0].(*upcloud.ServerDetails), args.Error(1)
}

func (m *MockService) StartServer(r *request.StartServerRequest) (*upcloud.ServerDetails, error) {
	args := m.Called(r)
	return args[0].(*upcloud.ServerDetails), args.Error(1)
}

func (m *MockService) StopServer(r *request.StopServerRequest) (*upcloud.ServerDetails, error) {
	args := m.Called(r)
	return args[0].(*upcloud.ServerDetails), args.Error(1)
}

func (m *MockService) RestartServer(r *request.RestartServerRequest) (*upcloud.ServerDetails, error) {
	args := m.Called(r)
	return args[0].(*upcloud.ServerDetails), args.Error(1)
}

func (m *MockService) ModifyServer(r *request.ModifyServerRequest) (*upcloud.ServerDetails, error) {
	args := m.Called(r)
	return args[0].(*upcloud.ServerDetails), args.Error(1)
}

func (m *MockService) DeleteServer(r *request.DeleteServerRequest) error {
	args := m.Called(r)
	return args.Error(0)
}

func (m *MockService) DeleteServerAndStorages(r *request.DeleteServerAndStoragesRequest) error {
	args := m.Called(r)
	return args.Error(0)
}

func (m *MockService) GetStorages(r *request.GetStoragesRequest) (*upcloud.Storages, error) {
	args := m.Called(r)
	return args[0].(*upcloud.Storages), args.Error(1)
}

func (m *MockService) GetStorageDetails(r *request.GetStorageDetailsRequest) (*upcloud.StorageDetails, error) {
	args := m.Called(r)
	return args[0].(*upcloud.StorageDetails), args.Error(1)
}

func (m *MockService) CreateStorage(r *request.CreateStorageRequest) (*upcloud.StorageDetails, error) {
	args := m.Called(r)
	return args[0].(*upcloud.StorageDetails), args.Error(1)
}

func (m *MockService) ModifyStorage(r *request.ModifyStorageRequest) (*upcloud.StorageDetails, error) {
	args := m.Called(r)
	return args[0].(*upcloud.StorageDetails), args.Error(1)
}

func (m *MockService) AttachStorage(r *request.AttachStorageRequest) (*upcloud.ServerDetails, error) {
	args := m.Called(r)
	return args[0].(*upcloud.ServerDetails), args.Error(1)
}

func (m *MockService) DetachStorage(r *request.DetachStorageRequest) (*upcloud.ServerDetails, error) {
	args := m.Called(r)
	return args[0].(*upcloud.ServerDetails), args.Error(1)
}

func (m *MockService) CloneStorage(r *request.CloneStorageRequest) (*upcloud.StorageDetails, error) {
	args := m.Called(r)
	return args[0].(*upcloud.StorageDetails), args.Error(1)
}

func (m *MockService) TemplatizeStorage(r *request.TemplatizeStorageRequest) (*upcloud.StorageDetails, error) {
	args := m.Called(r)
	return args[0].(*upcloud.StorageDetails), args.Error(1)
}

func (m *MockService) WaitForStorageState(r *request.WaitForStorageStateRequest) (*upcloud.StorageDetails, error) {
	args := m.Called(r)
	return args[0].(*upcloud.StorageDetails), args.Error(1)
}

func (m *MockService) LoadCDROM(r *request.LoadCDROMRequest) (*upcloud.ServerDetails, error) {
	args := m.Called(r)
	return args[0].(*upcloud.ServerDetails), args.Error(1)
}

func (m *MockService) EjectCDROM(r *request.EjectCDROMRequest) (*upcloud.ServerDetails, error) {
	args := m.Called(r)
	return args[0].(*upcloud.ServerDetails), args.Error(1)
}

func (m *MockService) CreateBackup(r *request.CreateBackupRequest) (*upcloud.StorageDetails, error) {
	args := m.Called(r)
	return args[0].(*upcloud.StorageDetails), args.Error(1)
}

func (m *MockService) RestoreBackup(r *request.RestoreBackupRequest) error {
	return m.Called(r).Error(0)
}

func (m *MockService) CreateStorageImport(r *request.CreateStorageImportRequest) (*upcloud.StorageImportDetails, error) {
	args := m.Called(r)
	return args[0].(*upcloud.StorageImportDetails), args.Error(1)
}

func (m *MockService) GetStorageImportDetails(r *request.GetStorageImportDetailsRequest) (*upcloud.StorageImportDetails, error) {
	args := m.Called(r)
	return args[0].(*upcloud.StorageImportDetails), args.Error(1)
}

func (m *MockService) WaitForStorageImportCompletion(r *request.WaitForStorageImportCompletionRequest) (*upcloud.StorageImportDetails, error) {
	args := m.Called(r)
	return args[0].(*upcloud.StorageImportDetails), args.Error(1)
}

func (m *MockService) DeleteStorage(r *request.DeleteStorageRequest) error {
	return m.Called(r).Error(0)
}

func (m *MockService) GetFirewallRules(r *request.GetFirewallRulesRequest) (*upcloud.FirewallRules, error) {
	args := m.Called(r)
	return args[0].(*upcloud.FirewallRules), args.Error(1)
}

func (m *MockService) GetFirewallRuleDetails(r *request.GetFirewallRuleDetailsRequest) (*upcloud.FirewallRule, error) {
	args := m.Called(r)
	return args[0].(*upcloud.FirewallRule), args.Error(1)
}

func (m *MockService) CreateFirewallRule(r *request.CreateFirewallRuleRequest) (*upcloud.FirewallRule, error) {
	args := m.Called(r)
	return args[0].(*upcloud.FirewallRule), args.Error(1)
}

func (m *MockService) CreateFirewallRules(r *request.CreateFirewallRulesRequest) error {
	args := m.Called(r)
	return args.Error(0)
}

func (m *MockService) DeleteFirewallRule(r *request.DeleteFirewallRuleRequest) error {
	args := m.Called(r)
	return args.Error(0)
}

func (m *MockService) GetNetworks() (*upcloud.Networks, error) {
	args := m.Called()
	return args[0].(*upcloud.Networks), args.Error(1)
}

func (m *MockService) GetNetworksInZone(r *request.GetNetworksInZoneRequest) (*upcloud.Networks, error) {
	args := m.Called(r)
	return args[0].(*upcloud.Networks), args.Error(1)
}

func (m *MockService) CreateNetwork(r *request.CreateNetworkRequest) (*upcloud.Network, error) {
	args := m.Called(r)
	return args[0].(*upcloud.Network), args.Error(1)
}

func (m *MockService) GetNetworkDetails(r *request.GetNetworkDetailsRequest) (*upcloud.Network, error) {
	args := m.Called(r)
	return args[0].(*upcloud.Network), args.Error(1)
}

func (m *MockService) ModifyNetwork(r *request.ModifyNetworkRequest) (*upcloud.Network, error) {
	args := m.Called(r)
	return args[0].(*upcloud.Network), args.Error(1)
}

func (m *MockService) GetServerNetworks(r *request.GetServerNetworksRequest) (*upcloud.Networking, error) {
	args := m.Called(r)
	return args[0].(*upcloud.Networking), args.Error(1)
}

func (m *MockService) CreateNetworkInterface(r *request.CreateNetworkInterfaceRequest) (*upcloud.Interface, error) {
	args := m.Called(r)
	return args[0].(*upcloud.Interface), args.Error(1)
}

func (m *MockService) ModifyNetworkInterface(r *request.ModifyNetworkInterfaceRequest) (*upcloud.Interface, error) {
	args := m.Called(r)
	return args[0].(*upcloud.Interface), args.Error(1)
}

func (m *MockService) DeleteNetwork(r *request.DeleteNetworkRequest) error {
	args := m.Called(r)
	return args.Error(0)
}

func (m *MockService) DeleteNetworkInterface(r *request.DeleteNetworkInterfaceRequest) error {
	args := m.Called(r)
	return args.Error(0)
}

func (m *MockService) GetRouters() (*upcloud.Routers, error) {
	args := m.Called()
	return args[0].(*upcloud.Routers), args.Error(1)
}

func (m *MockService) GetRouterDetails(r *request.GetRouterDetailsRequest) (*upcloud.Router, error) {
	args := m.Called(r)
	return args[0].(*upcloud.Router), args.Error(1)
}

func (m *MockService) CreateRouter(r *request.CreateRouterRequest) (*upcloud.Router, error) {
	args := m.Called(r)
	return args[0].(*upcloud.Router), args.Error(1)
}

func (m *MockService) ModifyRouter(r *request.ModifyRouterRequest) (*upcloud.Router, error) {
	args := m.Called(r)
	return args[0].(*upcloud.Router), args.Error(1)
}

func (m *MockService) DeleteRouter(r *request.DeleteRouterRequest) error {
	args := m.Called(r)
	return args.Error(0)
}

func (m *MockService) GetIPAddresses() (*upcloud.IPAddresses, error) {
	args := m.Called()
	return args[0].(*upcloud.IPAddresses), args.Error(1)
}

func (m *MockService) GetIPAddressDetails(r *request.GetIPAddressDetailsRequest) (*upcloud.IPAddress, error) {
	args := m.Called(r)
	return args[0].(*upcloud.IPAddress), args.Error(1)
}

func (m *MockService) AssignIPAddress(r *request.AssignIPAddressRequest) (*upcloud.IPAddress, error) {
	args := m.Called(r)
	return args[0].(*upcloud.IPAddress), args.Error(1)
}

func (m *MockService) ModifyIPAddress(r *request.ModifyIPAddressRequest) (*upcloud.IPAddress, error) {
	args := m.Called(r)
	return args[0].(*upcloud.IPAddress), args.Error(1)
}

func (m *MockService) ReleaseIPAddress(r *request.ReleaseIPAddressRequest) error {
	args := m.Called(r)
	return args.Error(0)
}
