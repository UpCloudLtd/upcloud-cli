package mock

import (
	"github.com/UpCloudLtd/upcloud-go-api/v4/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/v4/upcloud/request"
	"github.com/UpCloudLtd/upcloud-go-api/v4/upcloud/service"
	"github.com/stretchr/testify/mock"
)

// Service represents a mock upcloud API service
type Service struct {
	mock.Mock
}

// GetAccount implements service.Account.GetAccount
func (m *Service) GetAccount() (*upcloud.Account, error) {
	args := m.Called()
	if args[0] == nil {
		return nil, args.Error(1)
	}
	return args[0].(*upcloud.Account), args.Error(1)
}

// GetPlans implements service.Plan.GetPlans
func (m *Service) GetPlans() (*upcloud.Plans, error) {
	args := m.Called()
	if args[0] == nil {
		return nil, args.Error(1)
	}
	return args[0].(*upcloud.Plans), args.Error(1)
}

// make sure Service implements service interfaces
var _ service.Server = &Service{}
var _ service.Storage = &Service{}
var _ service.Firewall = &Service{}
var _ service.Network = &Service{}
var _ service.Plans = &Service{}
var _ service.Account = &Service{}

// GetServerConfigurations implements service.Server.GetServerConfigurations
func (m *Service) GetServerConfigurations() (*upcloud.ServerConfigurations, error) {
	args := m.Called()
	if args[0] == nil {
		return nil, args.Error(1)
	}
	return args[0].(*upcloud.ServerConfigurations), args.Error(1)
}

// GetServers implements service.Server.GetServers
func (m *Service) GetServers() (*upcloud.Servers, error) {
	args := m.Called()
	if args[0] == nil {
		return nil, args.Error(1)
	}
	return args[0].(*upcloud.Servers), args.Error(1)
}

// GetServerDetails implements service.Server.GetServerDetails
func (m *Service) GetServerDetails(r *request.GetServerDetailsRequest) (*upcloud.ServerDetails, error) {
	args := m.Called(r)
	if args[0] == nil {
		return nil, args.Error(1)
	}
	return args[0].(*upcloud.ServerDetails), args.Error(1)
}

// CreateServer implements service.Server.CreateServer
func (m *Service) CreateServer(r *request.CreateServerRequest) (*upcloud.ServerDetails, error) {
	args := m.Called(r)
	if args[0] == nil {
		return nil, args.Error(1)
	}
	return args[0].(*upcloud.ServerDetails), args.Error(1)
}

// WaitForServerState implements service.Server.WaitForServerState
func (m *Service) WaitForServerState(r *request.WaitForServerStateRequest) (*upcloud.ServerDetails, error) {
	args := m.Called(r)
	if args[0] == nil {
		return nil, args.Error(1)
	}
	return args[0].(*upcloud.ServerDetails), args.Error(1)
}

// StartServer implements service.Server.StartServer
func (m *Service) StartServer(r *request.StartServerRequest) (*upcloud.ServerDetails, error) {
	args := m.Called(r)
	if args[0] == nil {
		return nil, args.Error(1)
	}
	return args[0].(*upcloud.ServerDetails), args.Error(1)
}

// StopServer implements service.Server.StopServer
func (m *Service) StopServer(r *request.StopServerRequest) (*upcloud.ServerDetails, error) {
	args := m.Called(r)
	if args[0] == nil {
		return nil, args.Error(1)
	}
	return args[0].(*upcloud.ServerDetails), args.Error(1)
}

// RestartServer implements service.Server.RestartServer
func (m *Service) RestartServer(r *request.RestartServerRequest) (*upcloud.ServerDetails, error) {
	args := m.Called(r)
	if args[0] == nil {
		return nil, args.Error(1)
	}
	return args[0].(*upcloud.ServerDetails), args.Error(1)
}

// ModifyServer implements service.Server.ModifyServer
func (m *Service) ModifyServer(r *request.ModifyServerRequest) (*upcloud.ServerDetails, error) {
	args := m.Called(r)
	if args[0] == nil {
		return nil, args.Error(1)
	}
	return args[0].(*upcloud.ServerDetails), args.Error(1)
}

// DeleteServer implements service.Server.DeleteServer
func (m *Service) DeleteServer(r *request.DeleteServerRequest) error {
	args := m.Called(r)
	return args.Error(0)
}

// DeleteServerAndStorages implements service.Server.DeleteServerAndStorages
func (m *Service) DeleteServerAndStorages(r *request.DeleteServerAndStoragesRequest) error {
	args := m.Called(r)
	return args.Error(0)
}

// GetStorages implements service.Storage.GetStorages
func (m *Service) GetStorages(r *request.GetStoragesRequest) (*upcloud.Storages, error) {
	args := m.Called(r)
	if args[0] == nil {
		return nil, args.Error(1)
	}
	return args[0].(*upcloud.Storages), args.Error(1)
}

// GetStorageDetails implements service.Storage.GetStorageDetails
func (m *Service) GetStorageDetails(r *request.GetStorageDetailsRequest) (*upcloud.StorageDetails, error) {
	args := m.Called(r)
	if args[0] == nil {
		return nil, args.Error(1)
	}
	return args[0].(*upcloud.StorageDetails), args.Error(1)
}

// CreateStorage implements service.Storage.CreateStorage
func (m *Service) CreateStorage(r *request.CreateStorageRequest) (*upcloud.StorageDetails, error) {
	args := m.Called(r)
	if args[0] == nil {
		return nil, args.Error(1)
	}
	return args[0].(*upcloud.StorageDetails), args.Error(1)
}

// ModifyStorage implements service.Storage.ModifyStorage
func (m *Service) ModifyStorage(r *request.ModifyStorageRequest) (*upcloud.StorageDetails, error) {
	args := m.Called(r)
	if args[0] == nil {
		return nil, args.Error(1)
	}
	return args[0].(*upcloud.StorageDetails), args.Error(1)
}

// AttachStorage implements service.Storage.AttachStorage
func (m *Service) AttachStorage(r *request.AttachStorageRequest) (*upcloud.ServerDetails, error) {
	args := m.Called(r)
	if args[0] == nil {
		return nil, args.Error(1)
	}
	return args[0].(*upcloud.ServerDetails), args.Error(1)
}

// DetachStorage implements service.Storage.DetachStorage
func (m *Service) DetachStorage(r *request.DetachStorageRequest) (*upcloud.ServerDetails, error) {
	args := m.Called(r)
	if args[0] == nil {
		return nil, args.Error(1)
	}
	return args[0].(*upcloud.ServerDetails), args.Error(1)
}

// CloneStorage implements service.Storage.CloneStorage
func (m *Service) CloneStorage(r *request.CloneStorageRequest) (*upcloud.StorageDetails, error) {
	args := m.Called(r)
	if args[0] == nil {
		return nil, args.Error(1)
	}
	return args[0].(*upcloud.StorageDetails), args.Error(1)
}

// TemplatizeStorage implements service.Storage.TemplatizeStorage
func (m *Service) TemplatizeStorage(r *request.TemplatizeStorageRequest) (*upcloud.StorageDetails, error) {
	args := m.Called(r)
	if args[0] == nil {
		return nil, args.Error(1)
	}
	return args[0].(*upcloud.StorageDetails), args.Error(1)
}

// WaitForStorageState implements service.Storage.WaitForStorageState
func (m *Service) WaitForStorageState(r *request.WaitForStorageStateRequest) (*upcloud.StorageDetails, error) {
	args := m.Called(r)
	if args[0] == nil {
		return nil, args.Error(1)
	}
	return args[0].(*upcloud.StorageDetails), args.Error(1)
}

// LoadCDROM implements service.Storage.LoadCDDROM
func (m *Service) LoadCDROM(r *request.LoadCDROMRequest) (*upcloud.ServerDetails, error) {
	args := m.Called(r)
	if args[0] == nil {
		return nil, args.Error(1)
	}
	return args[0].(*upcloud.ServerDetails), args.Error(1)
}

// EjectCDROM implements service.Storage.EjectCDROM
func (m *Service) EjectCDROM(r *request.EjectCDROMRequest) (*upcloud.ServerDetails, error) {
	args := m.Called(r)
	if args[0] == nil {
		return nil, args.Error(1)
	}
	return args[0].(*upcloud.ServerDetails), args.Error(1)
}

// CreateBackup implements service.Storage.CreateBackup
func (m *Service) CreateBackup(r *request.CreateBackupRequest) (*upcloud.StorageDetails, error) {
	args := m.Called(r)
	if args[0] == nil {
		return nil, args.Error(1)
	}
	return args[0].(*upcloud.StorageDetails), args.Error(1)
}

// RestoreBackup implements service.Storage.RestoreBackup
func (m *Service) RestoreBackup(r *request.RestoreBackupRequest) error {
	return m.Called(r).Error(0)
}

// CreateStorageImport implements service.Storage.CreateStorageImport
func (m *Service) CreateStorageImport(r *request.CreateStorageImportRequest) (*upcloud.StorageImportDetails, error) {
	args := m.Called(r)
	if args[0] == nil {
		return nil, args.Error(1)
	}
	return args[0].(*upcloud.StorageImportDetails), args.Error(1)
}

// GetStorageImportDetails implements service.Storage.GetStorageImportDetails
func (m *Service) GetStorageImportDetails(r *request.GetStorageImportDetailsRequest) (*upcloud.StorageImportDetails, error) {
	args := m.Called(r)
	if args[0] == nil {
		return nil, args.Error(1)
	}
	return args[0].(*upcloud.StorageImportDetails), args.Error(1)
}

// WaitForStorageImportCompletion implements service.Storage.WaitForStorageImportCompletion
func (m *Service) WaitForStorageImportCompletion(r *request.WaitForStorageImportCompletionRequest) (*upcloud.StorageImportDetails, error) {
	args := m.Called(r)
	if args[0] == nil {
		return nil, args.Error(1)
	}
	return args[0].(*upcloud.StorageImportDetails), args.Error(1)
}

// DeleteStorage implements service.Storage.DeleteStorage
func (m *Service) DeleteStorage(r *request.DeleteStorageRequest) error {
	return m.Called(r).Error(0)
}

// GetFirewallRules implements service.Firewall.GetFirewallRules
func (m *Service) GetFirewallRules(r *request.GetFirewallRulesRequest) (*upcloud.FirewallRules, error) {
	args := m.Called(r)
	if args[0] == nil {
		return nil, args.Error(1)
	}
	return args[0].(*upcloud.FirewallRules), args.Error(1)
}

// GetFirewallRuleDetails implements service.Firewall.GetFirewallRuleDetails
func (m *Service) GetFirewallRuleDetails(r *request.GetFirewallRuleDetailsRequest) (*upcloud.FirewallRule, error) {
	args := m.Called(r)
	if args[0] == nil {
		return nil, args.Error(1)
	}
	return args[0].(*upcloud.FirewallRule), args.Error(1)
}

// CreateFirewallRule implements service.Firewall.CreateFirewallRule
func (m *Service) CreateFirewallRule(r *request.CreateFirewallRuleRequest) (*upcloud.FirewallRule, error) {
	args := m.Called(r)
	if args[0] == nil {
		return nil, args.Error(1)
	}
	return args[0].(*upcloud.FirewallRule), args.Error(1)
}

// CreateFirewallRules implements service.Firewall.CreateFirewallRules
func (m *Service) CreateFirewallRules(r *request.CreateFirewallRulesRequest) error {
	args := m.Called(r)
	return args.Error(0)
}

// DeleteFirewallRule implements service.Firewall.DeleteFirewallRule
func (m *Service) DeleteFirewallRule(r *request.DeleteFirewallRuleRequest) error {
	args := m.Called(r)
	return args.Error(0)
}

// GetNetworks implements service.Network.GetNetworks
func (m *Service) GetNetworks() (*upcloud.Networks, error) {
	args := m.Called()
	if args[0] == nil {
		return nil, args.Error(1)
	}
	return args[0].(*upcloud.Networks), args.Error(1)
}

// GetNetworksInZone implements service.Network.GetNetworksInZone
func (m *Service) GetNetworksInZone(r *request.GetNetworksInZoneRequest) (*upcloud.Networks, error) {
	args := m.Called(r)
	if args[0] == nil {
		return nil, args.Error(1)
	}
	return args[0].(*upcloud.Networks), args.Error(1)
}

// CreateNetwork implements service.Network.CreateNetwork
func (m *Service) CreateNetwork(r *request.CreateNetworkRequest) (*upcloud.Network, error) {
	args := m.Called(r)
	if args[0] == nil {
		return nil, args.Error(1)
	}
	return args[0].(*upcloud.Network), args.Error(1)
}

// GetNetworkDetails implements service.Network.GetNetworkDetails
func (m *Service) GetNetworkDetails(r *request.GetNetworkDetailsRequest) (*upcloud.Network, error) {
	args := m.Called(r)
	if args[0] == nil {
		return nil, args.Error(1)
	}
	return args[0].(*upcloud.Network), args.Error(1)
}

// ModifyNetwork implements service.Network.ModifyNetwork
func (m *Service) ModifyNetwork(r *request.ModifyNetworkRequest) (*upcloud.Network, error) {
	args := m.Called(r)
	if args[0] == nil {
		return nil, args.Error(1)
	}
	return args[0].(*upcloud.Network), args.Error(1)
}

// AttachNetworkRouter implements service.Network.AttachNetworkRouter
func (m *Service) AttachNetworkRouter(r *request.AttachNetworkRouterRequest) error {
	args := m.Called(r)
	return args.Error(0)
}

// DetachNetworkRouter implements service.Network.DetachNetworkRouter
func (m *Service) DetachNetworkRouter(r *request.DetachNetworkRouterRequest) error {
	args := m.Called(r)
	return args.Error(0)
}

// GetServerNetworks implements service.Network.GetServerNetworks
func (m *Service) GetServerNetworks(r *request.GetServerNetworksRequest) (*upcloud.Networking, error) {
	args := m.Called(r)
	if args[0] == nil {
		return nil, args.Error(1)
	}
	return args[0].(*upcloud.Networking), args.Error(1)
}

// CreateNetworkInterface implements service.Network.CreateNetworkInterface
func (m *Service) CreateNetworkInterface(r *request.CreateNetworkInterfaceRequest) (*upcloud.Interface, error) {
	args := m.Called(r)
	if args[0] == nil {
		return nil, args.Error(1)
	}
	return args[0].(*upcloud.Interface), args.Error(1)
}

// ModifyNetworkInterface implements service.Network.ModifyNetworkInterface
func (m *Service) ModifyNetworkInterface(r *request.ModifyNetworkInterfaceRequest) (*upcloud.Interface, error) {
	args := m.Called(r)
	if args[0] == nil {
		return nil, args.Error(1)
	}
	return args[0].(*upcloud.Interface), args.Error(1)
}

// DeleteNetwork implements service.Network.DeleteNetwork
func (m *Service) DeleteNetwork(r *request.DeleteNetworkRequest) error {
	args := m.Called(r)
	return args.Error(0)
}

// DeleteNetworkInterface implements service.Network.DeleteNetworkInterface
func (m *Service) DeleteNetworkInterface(r *request.DeleteNetworkInterfaceRequest) error {
	args := m.Called(r)
	return args.Error(0)
}

// GetRouters implements service.Network.GetRouters
func (m *Service) GetRouters() (*upcloud.Routers, error) {
	args := m.Called()
	if args[0] == nil {
		return nil, args.Error(1)
	}
	return args[0].(*upcloud.Routers), args.Error(1)
}

// GetRouterDetails implements service.Network.GetRouterDetails
func (m *Service) GetRouterDetails(r *request.GetRouterDetailsRequest) (*upcloud.Router, error) {
	args := m.Called(r)
	if args[0] == nil {
		return nil, args.Error(1)
	}
	return args[0].(*upcloud.Router), args.Error(1)
}

// CreateRouter implements service.Network.CreateRouter
func (m *Service) CreateRouter(r *request.CreateRouterRequest) (*upcloud.Router, error) {
	args := m.Called(r)
	if args[0] == nil {
		return nil, args.Error(1)
	}
	return args[0].(*upcloud.Router), args.Error(1)
}

// ModifyRouter implements service.Network.ModifyRouter
func (m *Service) ModifyRouter(r *request.ModifyRouterRequest) (*upcloud.Router, error) {
	args := m.Called(r)
	if args[0] == nil {
		return nil, args.Error(1)
	}
	return args[0].(*upcloud.Router), args.Error(1)
}

// DeleteRouter implements service.Network.DeleteRouter
func (m *Service) DeleteRouter(r *request.DeleteRouterRequest) error {
	args := m.Called(r)
	return args.Error(0)
}

// GetIPAddresses implements service.Network.GetIPAddresses
func (m *Service) GetIPAddresses() (*upcloud.IPAddresses, error) {
	args := m.Called()
	if args[0] == nil {
		return nil, args.Error(1)
	}
	return args[0].(*upcloud.IPAddresses), args.Error(1)
}

// GetIPAddressDetails implements service.Network.GetIPAddressDetails
func (m *Service) GetIPAddressDetails(r *request.GetIPAddressDetailsRequest) (*upcloud.IPAddress, error) {
	args := m.Called(r)
	if args[0] == nil {
		return nil, args.Error(1)
	}
	return args[0].(*upcloud.IPAddress), args.Error(1)
}

// AssignIPAddress implements service.Network.AssignIPAddress
func (m *Service) AssignIPAddress(r *request.AssignIPAddressRequest) (*upcloud.IPAddress, error) {
	args := m.Called(r)
	if args[0] == nil {
		return nil, args.Error(1)
	}
	return args[0].(*upcloud.IPAddress), args.Error(1)
}

// ModifyIPAddress implements service.Network.ModifyIPAddress
func (m *Service) ModifyIPAddress(r *request.ModifyIPAddressRequest) (*upcloud.IPAddress, error) {
	args := m.Called(r)
	if args[0] == nil {
		return nil, args.Error(1)
	}
	return args[0].(*upcloud.IPAddress), args.Error(1)
}

// ReleaseIPAddress implements service.Network.ReleaseIPAddress
func (m *Service) ReleaseIPAddress(r *request.ReleaseIPAddressRequest) error {
	args := m.Called(r)
	return args.Error(0)
}

// ResizeStorageFilesystem implements service.Storage.ResizeStorageFilesystem
func (m *Service) ResizeStorageFilesystem(r *request.ResizeStorageFilesystemRequest) (*upcloud.ResizeStorageFilesystemBackup, error) {
	args := m.Called(r)
	if args[0] == nil {
		return nil, args.Error(1)
	}
	return args[0].(*upcloud.ResizeStorageFilesystemBackup), args.Error(1)
}

// CreateSubaccount implements service.Account.CreateSubaccount
func (m *Service) CreateSubaccount(r *request.CreateSubaccountRequest) (*upcloud.AccountDetails, error) {
	args := m.Called(r)
	if args[0] == nil {
		return nil, args.Error(1)
	}
	return args[0].(*upcloud.AccountDetails), args.Error(1)
}

// GetAccountList implements service.Account.GetAccountList
func (m *Service) GetAccountList() (upcloud.AccountList, error) {
	args := m.Called()
	if args[0] == nil {
		return nil, args.Error(1)
	}
	return args[0].(upcloud.AccountList), args.Error(1)
}

// GetAccountDetails implements service.Account.GetAccountDetails
func (m *Service) GetAccountDetails(r *request.GetAccountDetailsRequest) (*upcloud.AccountDetails, error) {
	args := m.Called(r)
	if args[0] == nil {
		return nil, args.Error(1)
	}
	return args[0].(*upcloud.AccountDetails), args.Error(1)
}

// ModifySubaccount implements service.Account.ModifySubaccount
func (m *Service) ModifySubaccount(r *request.ModifySubaccountRequest) (*upcloud.AccountDetails, error) {
	args := m.Called(r)
	if args[0] == nil {
		return nil, args.Error(1)
	}
	return args[0].(*upcloud.AccountDetails), args.Error(1)
}

// DeleteSubaccount implements service.Account.DeleteSubaccount
func (m *Service) DeleteSubaccount(r *request.DeleteSubaccountRequest) error {
	args := m.Called(r)
	if args[0] == nil {
		return args.Error(1)
	}
	return nil
}
