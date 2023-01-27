//nolint:nilnil // Here nil, nil returns are used in not-implemented methods required to satisfy an interface
package mock

import (
	"context"

	"github.com/UpCloudLtd/upcloud-go-api/v5/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/v5/upcloud/request"
	"github.com/UpCloudLtd/upcloud-go-api/v5/upcloud/service"
	"github.com/stretchr/testify/mock"
)

// Service represents a mock upcloud API service
type Service struct {
	mock.Mock
}

// GetAccount implements service.Account.GetAccount
func (m *Service) GetAccount(context.Context) (*upcloud.Account, error) {
	args := m.Called()
	if args[0] == nil {
		return nil, args.Error(1)
	}
	return args[0].(*upcloud.Account), args.Error(1)
}

// GetZones implements service.Zones.GetZones
func (m *Service) GetZones(context.Context) (*upcloud.Zones, error) {
	args := m.Called()
	if args[0] == nil {
		return nil, args.Error(1)
	}
	return args[0].(*upcloud.Zones), args.Error(1)
}

// GetPriceZones implements service.Zones.GetPriceZones
func (m *Service) GetPriceZones(context.Context) (*upcloud.PriceZones, error) {
	return nil, nil
}

// GetPriceZones implements service.Zones.GetPriceZones
func (m *Service) GetTimeZones(context.Context) (*upcloud.TimeZones, error) {
	return nil, nil
}

// GetPlans implements service.Plan.GetPlans
func (m *Service) GetPlans(context.Context) (*upcloud.Plans, error) {
	args := m.Called()
	if args[0] == nil {
		return nil, args.Error(1)
	}
	return args[0].(*upcloud.Plans), args.Error(1)
}

// make sure Service implements service interfaces
var (
	_ service.Server       = &Service{}
	_ service.Storage      = &Service{}
	_ service.Firewall     = &Service{}
	_ service.Network      = &Service{}
	_ service.IPAddress    = &Service{}
	_ service.Cloud        = &Service{}
	_ service.Account      = &Service{}
	_ service.LoadBalancer = &Service{}
	_ service.Kubernetes   = &Service{}
)

// GetServerConfigurations implements service.Server.GetServerConfigurations
func (m *Service) GetServerConfigurations(context.Context) (*upcloud.ServerConfigurations, error) {
	args := m.Called()
	if args[0] == nil {
		return nil, args.Error(1)
	}
	return args[0].(*upcloud.ServerConfigurations), args.Error(1)
}

// GetServers implements service.Server.GetServers
func (m *Service) GetServers(context.Context) (*upcloud.Servers, error) {
	args := m.Called()
	if args[0] == nil {
		return nil, args.Error(1)
	}
	return args[0].(*upcloud.Servers), args.Error(1)
}

// GetServerDetails implements service.Server.GetServerDetails
func (m *Service) GetServerDetails(_ context.Context, r *request.GetServerDetailsRequest) (*upcloud.ServerDetails, error) {
	args := m.Called(r)
	if args[0] == nil {
		return nil, args.Error(1)
	}
	return args[0].(*upcloud.ServerDetails), args.Error(1)
}

// CreateServer implements service.Server.CreateServer
func (m *Service) CreateServer(_ context.Context, r *request.CreateServerRequest) (*upcloud.ServerDetails, error) {
	args := m.Called(r)
	if args[0] == nil {
		return nil, args.Error(1)
	}
	return args[0].(*upcloud.ServerDetails), args.Error(1)
}

// WaitForServerState implements service.Server.WaitForServerState
func (m *Service) WaitForServerState(_ context.Context, r *request.WaitForServerStateRequest) (*upcloud.ServerDetails, error) {
	args := m.Called(r)
	if args[0] == nil {
		return nil, args.Error(1)
	}
	return args[0].(*upcloud.ServerDetails), args.Error(1)
}

// StartServer implements service.Server.StartServer
func (m *Service) StartServer(_ context.Context, r *request.StartServerRequest) (*upcloud.ServerDetails, error) {
	args := m.Called(r)
	if args[0] == nil {
		return nil, args.Error(1)
	}
	return args[0].(*upcloud.ServerDetails), args.Error(1)
}

// StopServer implements service.Server.StopServer
func (m *Service) StopServer(_ context.Context, r *request.StopServerRequest) (*upcloud.ServerDetails, error) {
	args := m.Called(r)
	if args[0] == nil {
		return nil, args.Error(1)
	}
	return args[0].(*upcloud.ServerDetails), args.Error(1)
}

// RestartServer implements service.Server.RestartServer
func (m *Service) RestartServer(_ context.Context, r *request.RestartServerRequest) (*upcloud.ServerDetails, error) {
	args := m.Called(r)
	if args[0] == nil {
		return nil, args.Error(1)
	}
	return args[0].(*upcloud.ServerDetails), args.Error(1)
}

// ModifyServer implements service.Server.ModifyServer
func (m *Service) ModifyServer(_ context.Context, r *request.ModifyServerRequest) (*upcloud.ServerDetails, error) {
	args := m.Called(r)
	if args[0] == nil {
		return nil, args.Error(1)
	}
	return args[0].(*upcloud.ServerDetails), args.Error(1)
}

// DeleteServer implements service.Server.DeleteServer
func (m *Service) DeleteServer(_ context.Context, r *request.DeleteServerRequest) error {
	args := m.Called(r)
	return args.Error(0)
}

// DeleteServerAndStorages implements service.Server.DeleteServerAndStorages
func (m *Service) DeleteServerAndStorages(_ context.Context, r *request.DeleteServerAndStoragesRequest) error {
	args := m.Called(r)
	return args.Error(0)
}

// GetStorages implements service.Storage.GetStorages
func (m *Service) GetStorages(_ context.Context, r *request.GetStoragesRequest) (*upcloud.Storages, error) {
	args := m.Called(r)
	if args[0] == nil {
		return nil, args.Error(1)
	}
	return args[0].(*upcloud.Storages), args.Error(1)
}

// GetStorageDetails implements service.Storage.GetStorageDetails
func (m *Service) GetStorageDetails(_ context.Context, r *request.GetStorageDetailsRequest) (*upcloud.StorageDetails, error) {
	args := m.Called(r)
	if args[0] == nil {
		return nil, args.Error(1)
	}
	return args[0].(*upcloud.StorageDetails), args.Error(1)
}

// CreateStorage implements service.Storage.CreateStorage
func (m *Service) CreateStorage(_ context.Context, r *request.CreateStorageRequest) (*upcloud.StorageDetails, error) {
	args := m.Called(r)
	if args[0] == nil {
		return nil, args.Error(1)
	}
	return args[0].(*upcloud.StorageDetails), args.Error(1)
}

// ModifyStorage implements service.Storage.ModifyStorage
func (m *Service) ModifyStorage(_ context.Context, r *request.ModifyStorageRequest) (*upcloud.StorageDetails, error) {
	args := m.Called(r)
	if args[0] == nil {
		return nil, args.Error(1)
	}
	return args[0].(*upcloud.StorageDetails), args.Error(1)
}

// AttachStorage implements service.Storage.AttachStorage
func (m *Service) AttachStorage(_ context.Context, r *request.AttachStorageRequest) (*upcloud.ServerDetails, error) {
	args := m.Called(r)
	if args[0] == nil {
		return nil, args.Error(1)
	}
	return args[0].(*upcloud.ServerDetails), args.Error(1)
}

// DetachStorage implements service.Storage.DetachStorage
func (m *Service) DetachStorage(_ context.Context, r *request.DetachStorageRequest) (*upcloud.ServerDetails, error) {
	args := m.Called(r)
	if args[0] == nil {
		return nil, args.Error(1)
	}
	return args[0].(*upcloud.ServerDetails), args.Error(1)
}

// CloneStorage implements service.Storage.CloneStorage
func (m *Service) CloneStorage(_ context.Context, r *request.CloneStorageRequest) (*upcloud.StorageDetails, error) {
	args := m.Called(r)
	if args[0] == nil {
		return nil, args.Error(1)
	}
	return args[0].(*upcloud.StorageDetails), args.Error(1)
}

// TemplatizeStorage implements service.Storage.TemplatizeStorage
func (m *Service) TemplatizeStorage(_ context.Context, r *request.TemplatizeStorageRequest) (*upcloud.StorageDetails, error) {
	args := m.Called(r)
	if args[0] == nil {
		return nil, args.Error(1)
	}
	return args[0].(*upcloud.StorageDetails), args.Error(1)
}

// WaitForStorageState implements service.Storage.WaitForStorageState
func (m *Service) WaitForStorageState(_ context.Context, r *request.WaitForStorageStateRequest) (*upcloud.StorageDetails, error) {
	args := m.Called(r)
	if args[0] == nil {
		return nil, args.Error(1)
	}
	return args[0].(*upcloud.StorageDetails), args.Error(1)
}

// LoadCDROM implements service.Storage.LoadCDDROM
func (m *Service) LoadCDROM(_ context.Context, r *request.LoadCDROMRequest) (*upcloud.ServerDetails, error) {
	args := m.Called(r)
	if args[0] == nil {
		return nil, args.Error(1)
	}
	return args[0].(*upcloud.ServerDetails), args.Error(1)
}

// EjectCDROM implements service.Storage.EjectCDROM
func (m *Service) EjectCDROM(_ context.Context, r *request.EjectCDROMRequest) (*upcloud.ServerDetails, error) {
	args := m.Called(r)
	if args[0] == nil {
		return nil, args.Error(1)
	}
	return args[0].(*upcloud.ServerDetails), args.Error(1)
}

// CreateBackup implements service.Storage.CreateBackup
func (m *Service) CreateBackup(_ context.Context, r *request.CreateBackupRequest) (*upcloud.StorageDetails, error) {
	args := m.Called(r)
	if args[0] == nil {
		return nil, args.Error(1)
	}
	return args[0].(*upcloud.StorageDetails), args.Error(1)
}

// RestoreBackup implements service.Storage.RestoreBackup
func (m *Service) RestoreBackup(_ context.Context, r *request.RestoreBackupRequest) error {
	return m.Called(r).Error(0)
}

// CreateStorageImport implements service.Storage.CreateStorageImport
func (m *Service) CreateStorageImport(_ context.Context, r *request.CreateStorageImportRequest) (*upcloud.StorageImportDetails, error) {
	args := m.Called(r)
	if args[0] == nil {
		return nil, args.Error(1)
	}
	return args[0].(*upcloud.StorageImportDetails), args.Error(1)
}

// GetStorageImportDetails implements service.Storage.GetStorageImportDetails
func (m *Service) GetStorageImportDetails(_ context.Context, r *request.GetStorageImportDetailsRequest) (*upcloud.StorageImportDetails, error) {
	args := m.Called(r)
	if args[0] == nil {
		return nil, args.Error(1)
	}
	return args[0].(*upcloud.StorageImportDetails), args.Error(1)
}

// WaitForStorageImportCompletion implements service.Storage.WaitForStorageImportCompletion
func (m *Service) WaitForStorageImportCompletion(_ context.Context, r *request.WaitForStorageImportCompletionRequest) (*upcloud.StorageImportDetails, error) {
	args := m.Called(r)
	if args[0] == nil {
		return nil, args.Error(1)
	}
	return args[0].(*upcloud.StorageImportDetails), args.Error(1)
}

// DeleteStorage implements service.Storage.DeleteStorage
func (m *Service) DeleteStorage(_ context.Context, r *request.DeleteStorageRequest) error {
	return m.Called(r).Error(0)
}

// GetFirewallRules implements service.Firewall.GetFirewallRules
func (m *Service) GetFirewallRules(_ context.Context, r *request.GetFirewallRulesRequest) (*upcloud.FirewallRules, error) {
	args := m.Called(r)
	if args[0] == nil {
		return nil, args.Error(1)
	}
	return args[0].(*upcloud.FirewallRules), args.Error(1)
}

// GetFirewallRuleDetails implements service.Firewall.GetFirewallRuleDetails
func (m *Service) GetFirewallRuleDetails(_ context.Context, r *request.GetFirewallRuleDetailsRequest) (*upcloud.FirewallRule, error) {
	args := m.Called(r)
	if args[0] == nil {
		return nil, args.Error(1)
	}
	return args[0].(*upcloud.FirewallRule), args.Error(1)
}

// CreateFirewallRule implements service.Firewall.CreateFirewallRule
func (m *Service) CreateFirewallRule(_ context.Context, r *request.CreateFirewallRuleRequest) (*upcloud.FirewallRule, error) {
	args := m.Called(r)
	if args[0] == nil {
		return nil, args.Error(1)
	}
	return args[0].(*upcloud.FirewallRule), args.Error(1)
}

// CreateFirewallRules implements service.Firewall.CreateFirewallRules
func (m *Service) CreateFirewallRules(_ context.Context, r *request.CreateFirewallRulesRequest) error {
	args := m.Called(r)
	return args.Error(0)
}

// DeleteFirewallRule implements service.Firewall.DeleteFirewallRule
func (m *Service) DeleteFirewallRule(_ context.Context, r *request.DeleteFirewallRuleRequest) error {
	args := m.Called(r)
	return args.Error(0)
}

// GetNetworks implements service.Network.GetNetworks
func (m *Service) GetNetworks(context.Context, ...request.QueryFilter) (*upcloud.Networks, error) {
	args := m.Called()
	if args[0] == nil {
		return nil, args.Error(1)
	}
	return args[0].(*upcloud.Networks), args.Error(1)
}

// GetNetworksInZone implements service.Network.GetNetworksInZone
func (m *Service) GetNetworksInZone(_ context.Context, r *request.GetNetworksInZoneRequest) (*upcloud.Networks, error) {
	args := m.Called(r)
	if args[0] == nil {
		return nil, args.Error(1)
	}
	return args[0].(*upcloud.Networks), args.Error(1)
}

// CreateNetwork implements service.Network.CreateNetwork
func (m *Service) CreateNetwork(_ context.Context, r *request.CreateNetworkRequest) (*upcloud.Network, error) {
	args := m.Called(r)
	if args[0] == nil {
		return nil, args.Error(1)
	}
	return args[0].(*upcloud.Network), args.Error(1)
}

// GetNetworkDetails implements service.Network.GetNetworkDetails
func (m *Service) GetNetworkDetails(_ context.Context, r *request.GetNetworkDetailsRequest) (*upcloud.Network, error) {
	args := m.Called(r)
	if args[0] == nil {
		return nil, args.Error(1)
	}
	return args[0].(*upcloud.Network), args.Error(1)
}

// ModifyNetwork implements service.Network.ModifyNetwork
func (m *Service) ModifyNetwork(_ context.Context, r *request.ModifyNetworkRequest) (*upcloud.Network, error) {
	args := m.Called(r)
	if args[0] == nil {
		return nil, args.Error(1)
	}
	return args[0].(*upcloud.Network), args.Error(1)
}

// AttachNetworkRouter implements service.Network.AttachNetworkRouter
func (m *Service) AttachNetworkRouter(_ context.Context, r *request.AttachNetworkRouterRequest) error {
	args := m.Called(r)
	return args.Error(0)
}

// DetachNetworkRouter implements service.Network.DetachNetworkRouter
func (m *Service) DetachNetworkRouter(_ context.Context, r *request.DetachNetworkRouterRequest) error {
	args := m.Called(r)
	return args.Error(0)
}

// GetServerNetworks implements service.Network.GetServerNetworks
func (m *Service) GetServerNetworks(_ context.Context, r *request.GetServerNetworksRequest) (*upcloud.Networking, error) {
	args := m.Called(r)
	if args[0] == nil {
		return nil, args.Error(1)
	}
	return args[0].(*upcloud.Networking), args.Error(1)
}

// CreateNetworkInterface implements service.Network.CreateNetworkInterface
func (m *Service) CreateNetworkInterface(_ context.Context, r *request.CreateNetworkInterfaceRequest) (*upcloud.Interface, error) {
	args := m.Called(r)
	if args[0] == nil {
		return nil, args.Error(1)
	}
	return args[0].(*upcloud.Interface), args.Error(1)
}

// ModifyNetworkInterface implements service.Network.ModifyNetworkInterface
func (m *Service) ModifyNetworkInterface(_ context.Context, r *request.ModifyNetworkInterfaceRequest) (*upcloud.Interface, error) {
	args := m.Called(r)
	if args[0] == nil {
		return nil, args.Error(1)
	}
	return args[0].(*upcloud.Interface), args.Error(1)
}

// DeleteNetwork implements service.Network.DeleteNetwork
func (m *Service) DeleteNetwork(_ context.Context, r *request.DeleteNetworkRequest) error {
	args := m.Called(r)
	return args.Error(0)
}

// DeleteNetworkInterface implements service.Network.DeleteNetworkInterface
func (m *Service) DeleteNetworkInterface(_ context.Context, r *request.DeleteNetworkInterfaceRequest) error {
	args := m.Called(r)
	return args.Error(0)
}

// GetRouters implements service.Network.GetRouters
func (m *Service) GetRouters(context.Context, ...request.QueryFilter) (*upcloud.Routers, error) {
	args := m.Called()
	if args[0] == nil {
		return nil, args.Error(1)
	}
	return args[0].(*upcloud.Routers), args.Error(1)
}

// GetRouterDetails implements service.Network.GetRouterDetails
func (m *Service) GetRouterDetails(_ context.Context, r *request.GetRouterDetailsRequest) (*upcloud.Router, error) {
	args := m.Called(r)
	if args[0] == nil {
		return nil, args.Error(1)
	}
	return args[0].(*upcloud.Router), args.Error(1)
}

// CreateRouter implements service.Network.CreateRouter
func (m *Service) CreateRouter(_ context.Context, r *request.CreateRouterRequest) (*upcloud.Router, error) {
	args := m.Called(r)
	if args[0] == nil {
		return nil, args.Error(1)
	}
	return args[0].(*upcloud.Router), args.Error(1)
}

// ModifyRouter implements service.Network.ModifyRouter
func (m *Service) ModifyRouter(_ context.Context, r *request.ModifyRouterRequest) (*upcloud.Router, error) {
	args := m.Called(r)
	if args[0] == nil {
		return nil, args.Error(1)
	}
	return args[0].(*upcloud.Router), args.Error(1)
}

// DeleteRouter implements service.Network.DeleteRouter
func (m *Service) DeleteRouter(_ context.Context, r *request.DeleteRouterRequest) error {
	args := m.Called(r)
	return args.Error(0)
}

// GetIPAddresses implements service.Network.GetIPAddresses
func (m *Service) GetIPAddresses(context.Context) (*upcloud.IPAddresses, error) {
	args := m.Called()
	if args[0] == nil {
		return nil, args.Error(1)
	}
	return args[0].(*upcloud.IPAddresses), args.Error(1)
}

// GetIPAddressDetails implements service.Network.GetIPAddressDetails
func (m *Service) GetIPAddressDetails(_ context.Context, r *request.GetIPAddressDetailsRequest) (*upcloud.IPAddress, error) {
	args := m.Called(r)
	if args[0] == nil {
		return nil, args.Error(1)
	}
	return args[0].(*upcloud.IPAddress), args.Error(1)
}

// AssignIPAddress implements service.Network.AssignIPAddress
func (m *Service) AssignIPAddress(_ context.Context, r *request.AssignIPAddressRequest) (*upcloud.IPAddress, error) {
	args := m.Called(r)
	if args[0] == nil {
		return nil, args.Error(1)
	}
	return args[0].(*upcloud.IPAddress), args.Error(1)
}

// ModifyIPAddress implements service.Network.ModifyIPAddress
func (m *Service) ModifyIPAddress(_ context.Context, r *request.ModifyIPAddressRequest) (*upcloud.IPAddress, error) {
	args := m.Called(r)
	if args[0] == nil {
		return nil, args.Error(1)
	}
	return args[0].(*upcloud.IPAddress), args.Error(1)
}

// ReleaseIPAddress implements service.Network.ReleaseIPAddress
func (m *Service) ReleaseIPAddress(_ context.Context, r *request.ReleaseIPAddressRequest) error {
	args := m.Called(r)
	return args.Error(0)
}

// ResizeStorageFilesystem implements service.Storage.ResizeStorageFilesystem
func (m *Service) ResizeStorageFilesystem(_ context.Context, r *request.ResizeStorageFilesystemRequest) (*upcloud.ResizeStorageFilesystemBackup, error) {
	args := m.Called(r)
	if args[0] == nil {
		return nil, args.Error(1)
	}
	return args[0].(*upcloud.ResizeStorageFilesystemBackup), args.Error(1)
}

// CreateSubaccount implements service.Account.CreateSubaccount
func (m *Service) CreateSubaccount(_ context.Context, r *request.CreateSubaccountRequest) (*upcloud.AccountDetails, error) {
	args := m.Called(r)
	if args[0] == nil {
		return nil, args.Error(1)
	}
	return args[0].(*upcloud.AccountDetails), args.Error(1)
}

// GetAccountList implements service.Account.GetAccountList
func (m *Service) GetAccountList(context.Context) (upcloud.AccountList, error) {
	args := m.Called()
	if args[0] == nil {
		return nil, args.Error(1)
	}
	return args[0].(upcloud.AccountList), args.Error(1)
}

// GetAccountDetails implements service.Account.GetAccountDetails
func (m *Service) GetAccountDetails(_ context.Context, r *request.GetAccountDetailsRequest) (*upcloud.AccountDetails, error) {
	args := m.Called(r)
	if args[0] == nil {
		return nil, args.Error(1)
	}
	return args[0].(*upcloud.AccountDetails), args.Error(1)
}

// ModifySubaccount implements service.Account.ModifySubaccount
func (m *Service) ModifySubaccount(_ context.Context, r *request.ModifySubaccountRequest) (*upcloud.AccountDetails, error) {
	args := m.Called(r)
	if args[0] == nil {
		return nil, args.Error(1)
	}
	return args[0].(*upcloud.AccountDetails), args.Error(1)
}

// DeleteSubaccount implements service.Account.DeleteSubaccount
func (m *Service) DeleteSubaccount(_ context.Context, r *request.DeleteSubaccountRequest) error {
	args := m.Called(r)
	if args[0] == nil {
		return args.Error(1)
	}
	return nil
}

func (m *Service) CancelManagedDatabaseConnection(_ context.Context, r *request.CancelManagedDatabaseConnection) error {
	args := m.Called(r)
	if args[0] != nil {
		return args.Error(0)
	}
	return nil
}

func (m *Service) CloneManagedDatabase(_ context.Context, r *request.CloneManagedDatabaseRequest) (*upcloud.ManagedDatabase, error) {
	return nil, nil
}

func (m *Service) CreateManagedDatabase(_ context.Context, r *request.CreateManagedDatabaseRequest) (*upcloud.ManagedDatabase, error) {
	return nil, nil
}

func (m *Service) GetManagedDatabase(_ context.Context, r *request.GetManagedDatabaseRequest) (*upcloud.ManagedDatabase, error) {
	return nil, nil
}

func (m *Service) GetManagedDatabases(_ context.Context, r *request.GetManagedDatabasesRequest) ([]upcloud.ManagedDatabase, error) {
	args := m.Called(r)
	if args[0] == nil {
		return nil, args.Error(1)
	}
	return args[0].([]upcloud.ManagedDatabase), args.Error(1)
}

func (m *Service) GetManagedDatabaseConnections(_ context.Context, r *request.GetManagedDatabaseConnectionsRequest) ([]upcloud.ManagedDatabaseConnection, error) {
	return nil, nil
}

func (m *Service) GetManagedDatabaseMetrics(_ context.Context, r *request.GetManagedDatabaseMetricsRequest) (*upcloud.ManagedDatabaseMetrics, error) {
	return nil, nil
}

func (m *Service) GetManagedDatabaseLogs(_ context.Context, r *request.GetManagedDatabaseLogsRequest) (*upcloud.ManagedDatabaseLogs, error) {
	return nil, nil
}

func (m *Service) GetManagedDatabaseQueryStatisticsMySQL(_ context.Context, r *request.GetManagedDatabaseQueryStatisticsRequest) ([]upcloud.ManagedDatabaseQueryStatisticsMySQL, error) {
	return nil, nil
}

func (m *Service) GetManagedDatabaseQueryStatisticsPostgreSQL(_ context.Context, r *request.GetManagedDatabaseQueryStatisticsRequest) ([]upcloud.ManagedDatabaseQueryStatisticsPostgreSQL, error) {
	return nil, nil
}

func (m *Service) GetManagedDatabaseServiceType(_ context.Context, r *request.GetManagedDatabaseServiceTypeRequest) (*upcloud.ManagedDatabaseType, error) {
	args := m.Called(r)
	if args[0] == nil {
		return nil, args.Error(1)
	}
	return args[0].(*upcloud.ManagedDatabaseType), args.Error(1)
}

func (m *Service) GetManagedDatabaseServiceTypes(_ context.Context, r *request.GetManagedDatabaseServiceTypesRequest) (map[string]upcloud.ManagedDatabaseType, error) {
	args := m.Called(r)
	if args[0] == nil {
		return nil, args.Error(1)
	}
	return args[0].(map[string]upcloud.ManagedDatabaseType), args.Error(1)
}

func (m *Service) DeleteManagedDatabase(_ context.Context, r *request.DeleteManagedDatabaseRequest) error {
	args := m.Called(r)
	return args.Error(0)
}

func (m *Service) ModifyManagedDatabase(_ context.Context, r *request.ModifyManagedDatabaseRequest) (*upcloud.ManagedDatabase, error) {
	return nil, nil
}

func (m *Service) UpgradeManagedDatabaseVersion(_ context.Context, r *request.UpgradeManagedDatabaseVersionRequest) (*upcloud.ManagedDatabase, error) {
	return nil, nil
}

func (m *Service) GetManagedDatabaseVersions(_ context.Context, r *request.GetManagedDatabaseVersionsRequest) ([]string, error) {
	return nil, nil
}

func (m *Service) StartManagedDatabase(_ context.Context, r *request.StartManagedDatabaseRequest) (*upcloud.ManagedDatabase, error) {
	args := m.Called(r)
	if args[0] == nil {
		return nil, args.Error(1)
	}
	return args[0].(*upcloud.ManagedDatabase), args.Error(1)
}

func (m *Service) ShutdownManagedDatabase(_ context.Context, r *request.ShutdownManagedDatabaseRequest) (*upcloud.ManagedDatabase, error) {
	args := m.Called(r)
	if args[0] == nil {
		return nil, args.Error(1)
	}
	return args[0].(*upcloud.ManagedDatabase), args.Error(1)
}

func (m *Service) WaitForManagedDatabaseState(_ context.Context, r *request.WaitForManagedDatabaseStateRequest) (*upcloud.ManagedDatabase, error) {
	return nil, nil
}

func (m *Service) GetLoadBalancers(_ context.Context, r *request.GetLoadBalancersRequest) ([]upcloud.LoadBalancer, error) {
	args := m.Called(r)
	if args[0] == nil {
		return nil, args.Error(1)
	}
	return args[0].([]upcloud.LoadBalancer), args.Error(1)
}

func (m *Service) GetLoadBalancer(_ context.Context, r *request.GetLoadBalancerRequest) (*upcloud.LoadBalancer, error) {
	return nil, nil
}

func (m *Service) CreateLoadBalancer(_ context.Context, r *request.CreateLoadBalancerRequest) (*upcloud.LoadBalancer, error) {
	return nil, nil
}

func (m *Service) ModifyLoadBalancer(_ context.Context, r *request.ModifyLoadBalancerRequest) (*upcloud.LoadBalancer, error) {
	return nil, nil
}

func (m *Service) DeleteLoadBalancer(_ context.Context, r *request.DeleteLoadBalancerRequest) error {
	args := m.Called(r)
	return args.Error(0)
}

func (m *Service) GetLoadBalancerBackends(_ context.Context, r *request.GetLoadBalancerBackendsRequest) ([]upcloud.LoadBalancerBackend, error) {
	return nil, nil
}

func (m *Service) GetLoadBalancerBackend(_ context.Context, r *request.GetLoadBalancerBackendRequest) (*upcloud.LoadBalancerBackend, error) {
	return nil, nil
}

func (m *Service) CreateLoadBalancerBackend(_ context.Context, r *request.CreateLoadBalancerBackendRequest) (*upcloud.LoadBalancerBackend, error) {
	return nil, nil
}

func (m *Service) ModifyLoadBalancerBackend(_ context.Context, r *request.ModifyLoadBalancerBackendRequest) (*upcloud.LoadBalancerBackend, error) {
	return nil, nil
}

func (m *Service) DeleteLoadBalancerBackend(_ context.Context, r *request.DeleteLoadBalancerBackendRequest) error {
	return nil
}

func (m *Service) GetLoadBalancerBackendMembers(_ context.Context, r *request.GetLoadBalancerBackendMembersRequest) ([]upcloud.LoadBalancerBackendMember, error) {
	return nil, nil
}

func (m *Service) GetLoadBalancerBackendMember(_ context.Context, r *request.GetLoadBalancerBackendMemberRequest) (*upcloud.LoadBalancerBackendMember, error) {
	return nil, nil
}

func (m *Service) CreateLoadBalancerBackendMember(_ context.Context, r *request.CreateLoadBalancerBackendMemberRequest) (*upcloud.LoadBalancerBackendMember, error) {
	return nil, nil
}

func (m *Service) ModifyLoadBalancerBackendMember(_ context.Context, r *request.ModifyLoadBalancerBackendMemberRequest) (*upcloud.LoadBalancerBackendMember, error) {
	return nil, nil
}

func (m *Service) DeleteLoadBalancerBackendMember(_ context.Context, r *request.DeleteLoadBalancerBackendMemberRequest) error {
	return nil
}

func (m *Service) GetLoadBalancerResolvers(_ context.Context, r *request.GetLoadBalancerResolversRequest) ([]upcloud.LoadBalancerResolver, error) {
	return nil, nil
}

func (m *Service) CreateLoadBalancerResolver(_ context.Context, r *request.CreateLoadBalancerResolverRequest) (*upcloud.LoadBalancerResolver, error) {
	return nil, nil
}

func (m *Service) GetLoadBalancerResolver(_ context.Context, r *request.GetLoadBalancerResolverRequest) (*upcloud.LoadBalancerResolver, error) {
	return nil, nil
}

func (m *Service) ModifyLoadBalancerResolver(_ context.Context, r *request.ModifyLoadBalancerResolverRequest) (*upcloud.LoadBalancerResolver, error) {
	return nil, nil
}

func (m *Service) DeleteLoadBalancerResolver(_ context.Context, r *request.DeleteLoadBalancerResolverRequest) error {
	return nil
}

func (m *Service) GetLoadBalancerPlans(_ context.Context, r *request.GetLoadBalancerPlansRequest) ([]upcloud.LoadBalancerPlan, error) {
	return nil, nil
}

func (m *Service) GetLoadBalancerFrontends(_ context.Context, r *request.GetLoadBalancerFrontendsRequest) ([]upcloud.LoadBalancerFrontend, error) {
	return nil, nil
}

func (m *Service) GetLoadBalancerFrontend(_ context.Context, r *request.GetLoadBalancerFrontendRequest) (*upcloud.LoadBalancerFrontend, error) {
	return nil, nil
}

func (m *Service) CreateLoadBalancerFrontend(_ context.Context, r *request.CreateLoadBalancerFrontendRequest) (*upcloud.LoadBalancerFrontend, error) {
	return nil, nil
}

func (m *Service) ModifyLoadBalancerFrontend(_ context.Context, r *request.ModifyLoadBalancerFrontendRequest) (*upcloud.LoadBalancerFrontend, error) {
	return nil, nil
}

func (m *Service) DeleteLoadBalancerFrontend(_ context.Context, r *request.DeleteLoadBalancerFrontendRequest) error {
	return nil
}

func (m *Service) GetLoadBalancerFrontendRules(_ context.Context, r *request.GetLoadBalancerFrontendRulesRequest) ([]upcloud.LoadBalancerFrontendRule, error) {
	return nil, nil
}

func (m *Service) GetLoadBalancerFrontendRule(_ context.Context, r *request.GetLoadBalancerFrontendRuleRequest) (*upcloud.LoadBalancerFrontendRule, error) {
	return nil, nil
}

func (m *Service) CreateLoadBalancerFrontendRule(_ context.Context, r *request.CreateLoadBalancerFrontendRuleRequest) (*upcloud.LoadBalancerFrontendRule, error) {
	return nil, nil
}

func (m *Service) ModifyLoadBalancerFrontendRule(_ context.Context, r *request.ModifyLoadBalancerFrontendRuleRequest) (*upcloud.LoadBalancerFrontendRule, error) {
	return nil, nil
}

func (m *Service) ReplaceLoadBalancerFrontendRule(_ context.Context, r *request.ReplaceLoadBalancerFrontendRuleRequest) (*upcloud.LoadBalancerFrontendRule, error) {
	return nil, nil
}

func (m *Service) DeleteLoadBalancerFrontendRule(_ context.Context, r *request.DeleteLoadBalancerFrontendRuleRequest) error {
	return nil
}

func (m *Service) GetLoadBalancerFrontendTLSConfigs(_ context.Context, r *request.GetLoadBalancerFrontendTLSConfigsRequest) ([]upcloud.LoadBalancerFrontendTLSConfig, error) {
	return nil, nil
}

func (m *Service) GetLoadBalancerFrontendTLSConfig(_ context.Context, r *request.GetLoadBalancerFrontendTLSConfigRequest) (*upcloud.LoadBalancerFrontendTLSConfig, error) {
	return nil, nil
}

func (m *Service) CreateLoadBalancerFrontendTLSConfig(_ context.Context, r *request.CreateLoadBalancerFrontendTLSConfigRequest) (*upcloud.LoadBalancerFrontendTLSConfig, error) {
	return nil, nil
}

func (m *Service) ModifyLoadBalancerFrontendTLSConfig(_ context.Context, r *request.ModifyLoadBalancerFrontendTLSConfigRequest) (*upcloud.LoadBalancerFrontendTLSConfig, error) {
	return nil, nil
}

func (m *Service) DeleteLoadBalancerFrontendTLSConfig(_ context.Context, r *request.DeleteLoadBalancerFrontendTLSConfigRequest) error {
	return nil
}

func (m *Service) GetLoadBalancerCertificateBundles(_ context.Context, r *request.GetLoadBalancerCertificateBundlesRequest) ([]upcloud.LoadBalancerCertificateBundle, error) {
	return nil, nil
}

func (m *Service) GetLoadBalancerCertificateBundle(_ context.Context, r *request.GetLoadBalancerCertificateBundleRequest) (*upcloud.LoadBalancerCertificateBundle, error) {
	return nil, nil
}

func (m *Service) CreateLoadBalancerCertificateBundle(_ context.Context, r *request.CreateLoadBalancerCertificateBundleRequest) (*upcloud.LoadBalancerCertificateBundle, error) {
	return nil, nil
}

func (m *Service) ModifyLoadBalancerCertificateBundle(_ context.Context, r *request.ModifyLoadBalancerCertificateBundleRequest) (*upcloud.LoadBalancerCertificateBundle, error) {
	return nil, nil
}

func (m *Service) DeleteLoadBalancerCertificateBundle(_ context.Context, r *request.DeleteLoadBalancerCertificateBundleRequest) error {
	return nil
}

func (m *Service) ModifyLoadBalancerNetwork(ctx context.Context, r *request.ModifyLoadBalancerNetworkRequest) (*upcloud.LoadBalancerNetwork, error) {
	return nil, nil
}

func (m *Service) CreateKubernetesCluster(_ context.Context, r *request.CreateKubernetesClusterRequest) (*upcloud.KubernetesCluster, error) {
	args := m.Called(r)
	if args[0] == nil {
		return nil, args.Error(1)
	}
	return args[0].(*upcloud.KubernetesCluster), args.Error(1)
}

func (m *Service) DeleteKubernetesCluster(_ context.Context, r *request.DeleteKubernetesClusterRequest) error {
	args := m.Called(r)
	return args.Error(0)
}

func (m *Service) GetKubernetesCluster(_ context.Context, r *request.GetKubernetesClusterRequest) (*upcloud.KubernetesCluster, error) {
	args := m.Called(r)
	if args[0] == nil {
		return nil, args.Error(1)
	}
	return args[0].(*upcloud.KubernetesCluster), args.Error(1)
}

func (m *Service) GetKubernetesClusters(_ context.Context, r *request.GetKubernetesClustersRequest) ([]upcloud.KubernetesCluster, error) {
	args := m.Called(r)
	if args[0] == nil {
		return nil, args.Error(1)
	}
	return args[0].([]upcloud.KubernetesCluster), args.Error(1)
}

func (m *Service) GetKubernetesKubeconfig(_ context.Context, r *request.GetKubernetesKubeconfigRequest) (string, error) {
	args := m.Called(r)
	if args[0] == nil {
		return "", args.Error(1)
	}
	return args[0].(string), args.Error(1)
}

func (m *Service) GetKubernetesVersions(_ context.Context, r *request.GetKubernetesVersionsRequest) ([]string, error) {
	return nil, nil
}

func (m *Service) WaitForKubernetesClusterState(context.Context, *request.WaitForKubernetesClusterStateRequest) (*upcloud.KubernetesCluster, error) {
	return nil, nil
}

func (m *Service) GetKubernetesNodeGroups(ctx context.Context, r *request.GetKubernetesNodeGroupsRequest) ([]upcloud.KubernetesNodeGroup, error) {
	return nil, nil
}

func (m *Service) GetKubernetesNodeGroup(ctx context.Context, r *request.GetKubernetesNodeGroupRequest) (*upcloud.KubernetesNodeGroup, error) {
	return nil, nil
}

func (m *Service) CreateKubernetesNodeGroup(ctx context.Context, r *request.CreateKubernetesNodeGroupRequest) (*upcloud.KubernetesNodeGroup, error) {
	return nil, nil
}

func (m *Service) ModifyKubernetesNodeGroup(ctx context.Context, r *request.ModifyKubernetesNodeGroupRequest) (*upcloud.KubernetesNodeGroup, error) {
	return nil, nil
}

func (m *Service) DeleteKubernetesNodeGroup(ctx context.Context, r *request.DeleteKubernetesNodeGroupRequest) error {
	return nil
}
