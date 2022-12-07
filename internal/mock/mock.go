//nolint:nilnil // Here nil, nil returns are used in not-implemented methods required to satisfy an interface
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

// GetZones implements service.Zones.GetZones
func (m *Service) GetZones() (*upcloud.Zones, error) {
	args := m.Called()
	if args[0] == nil {
		return nil, args.Error(1)
	}
	return args[0].(*upcloud.Zones), args.Error(1)
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

var (
	_ service.Storage  = &Service{}
	_ service.Firewall = &Service{}
	_ service.Network  = &Service{}
	_ service.Plans    = &Service{}
	_ service.Account  = &Service{}
)

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

func (m *Service) CancelManagedDatabaseConnection(r *request.CancelManagedDatabaseConnection) error {
	args := m.Called(r)
	if args[0] != nil {
		return args.Error(0)
	}
	return nil
}

func (m *Service) CloneManagedDatabase(r *request.CloneManagedDatabaseRequest) (*upcloud.ManagedDatabase, error) {
	return nil, nil
}

func (m *Service) CreateManagedDatabase(r *request.CreateManagedDatabaseRequest) (*upcloud.ManagedDatabase, error) {
	return nil, nil
}

func (m *Service) GetManagedDatabase(r *request.GetManagedDatabaseRequest) (*upcloud.ManagedDatabase, error) {
	return nil, nil
}

func (m *Service) GetManagedDatabases(r *request.GetManagedDatabasesRequest) ([]upcloud.ManagedDatabase, error) {
	args := m.Called(r)
	if args[0] == nil {
		return nil, args.Error(1)
	}
	return args[0].([]upcloud.ManagedDatabase), args.Error(1)
}

func (m *Service) GetManagedDatabaseConnections(r *request.GetManagedDatabaseConnectionsRequest) ([]upcloud.ManagedDatabaseConnection, error) {
	return nil, nil
}

func (m *Service) GetManagedDatabaseMetrics(r *request.GetManagedDatabaseMetricsRequest) (*upcloud.ManagedDatabaseMetrics, error) {
	return nil, nil
}

func (m *Service) GetManagedDatabaseLogs(r *request.GetManagedDatabaseLogsRequest) (*upcloud.ManagedDatabaseLogs, error) {
	return nil, nil
}

func (m *Service) GetManagedDatabaseQueryStatisticsMySQL(r *request.GetManagedDatabaseQueryStatisticsRequest) ([]upcloud.ManagedDatabaseQueryStatisticsMySQL, error) {
	return nil, nil
}

func (m *Service) GetManagedDatabaseQueryStatisticsPostgreSQL(r *request.GetManagedDatabaseQueryStatisticsRequest) ([]upcloud.ManagedDatabaseQueryStatisticsPostgreSQL, error) {
	return nil, nil
}

func (m *Service) GetManagedDatabaseServiceType(r *request.GetManagedDatabaseServiceTypeRequest) (*upcloud.ManagedDatabaseType, error) {
	args := m.Called(r)
	if args[0] == nil {
		return nil, args.Error(1)
	}
	return args[0].(*upcloud.ManagedDatabaseType), args.Error(1)
}

func (m *Service) GetManagedDatabaseServiceTypes(r *request.GetManagedDatabaseServiceTypesRequest) (map[string]upcloud.ManagedDatabaseType, error) {
	args := m.Called(r)
	if args[0] == nil {
		return nil, args.Error(1)
	}
	return args[0].(map[string]upcloud.ManagedDatabaseType), args.Error(1)
}

func (m *Service) DeleteManagedDatabase(r *request.DeleteManagedDatabaseRequest) error {
	args := m.Called(r)
	return args.Error(0)
}

func (m *Service) ModifyManagedDatabase(r *request.ModifyManagedDatabaseRequest) (*upcloud.ManagedDatabase, error) {
	return nil, nil
}

func (m *Service) UpgradeManagedDatabaseVersion(r *request.UpgradeManagedDatabaseVersionRequest) (*upcloud.ManagedDatabase, error) {
	return nil, nil
}

func (m *Service) GetManagedDatabaseVersions(r *request.GetManagedDatabaseVersionsRequest) ([]string, error) {
	return nil, nil
}

func (m *Service) StartManagedDatabase(r *request.StartManagedDatabaseRequest) (*upcloud.ManagedDatabase, error) {
	args := m.Called(r)
	if args[0] == nil {
		return nil, args.Error(1)
	}
	return args[0].(*upcloud.ManagedDatabase), args.Error(1)
}

func (m *Service) ShutdownManagedDatabase(r *request.ShutdownManagedDatabaseRequest) (*upcloud.ManagedDatabase, error) {
	args := m.Called(r)
	if args[0] == nil {
		return nil, args.Error(1)
	}
	return args[0].(*upcloud.ManagedDatabase), args.Error(1)
}

func (m *Service) WaitForManagedDatabaseState(r *request.WaitForManagedDatabaseStateRequest) (*upcloud.ManagedDatabase, error) {
	return nil, nil
}

func (m *Service) GetLoadBalancers(r *request.GetLoadBalancersRequest) ([]upcloud.LoadBalancer, error) {
	args := m.Called(r)
	if args[0] == nil {
		return nil, args.Error(1)
	}
	return args[0].([]upcloud.LoadBalancer), args.Error(1)
}

func (m *Service) GetLoadBalancer(r *request.GetLoadBalancerRequest) (*upcloud.LoadBalancer, error) {
	return nil, nil
}

func (m *Service) CreateLoadBalancer(r *request.CreateLoadBalancerRequest) (*upcloud.LoadBalancer, error) {
	return nil, nil
}

func (m *Service) ModifyLoadBalancer(r *request.ModifyLoadBalancerRequest) (*upcloud.LoadBalancer, error) {
	return nil, nil
}

func (m *Service) DeleteLoadBalancer(r *request.DeleteLoadBalancerRequest) error {
	args := m.Called(r)
	return args.Error(0)
}

func (m *Service) GetLoadBalancerBackends(r *request.GetLoadBalancerBackendsRequest) ([]upcloud.LoadBalancerBackend, error) {
	return nil, nil
}

func (m *Service) GetLoadBalancerBackend(r *request.GetLoadBalancerBackendRequest) (*upcloud.LoadBalancerBackend, error) {
	return nil, nil
}

func (m *Service) CreateLoadBalancerBackend(r *request.CreateLoadBalancerBackendRequest) (*upcloud.LoadBalancerBackend, error) {
	return nil, nil
}

func (m *Service) ModifyLoadBalancerBackend(r *request.ModifyLoadBalancerBackendRequest) (*upcloud.LoadBalancerBackend, error) {
	return nil, nil
}

func (m *Service) DeleteLoadBalancerBackend(r *request.DeleteLoadBalancerBackendRequest) error {
	return nil
}

func (m *Service) GetLoadBalancerBackendMembers(r *request.GetLoadBalancerBackendMembersRequest) ([]upcloud.LoadBalancerBackendMember, error) {
	return nil, nil
}

func (m *Service) GetLoadBalancerBackendMember(r *request.GetLoadBalancerBackendMemberRequest) (*upcloud.LoadBalancerBackendMember, error) {
	return nil, nil
}

func (m *Service) CreateLoadBalancerBackendMember(r *request.CreateLoadBalancerBackendMemberRequest) (*upcloud.LoadBalancerBackendMember, error) {
	return nil, nil
}

func (m *Service) ModifyLoadBalancerBackendMember(r *request.ModifyLoadBalancerBackendMemberRequest) (*upcloud.LoadBalancerBackendMember, error) {
	return nil, nil
}

func (m *Service) DeleteLoadBalancerBackendMember(r *request.DeleteLoadBalancerBackendMemberRequest) error {
	return nil
}

func (m *Service) GetLoadBalancerResolvers(r *request.GetLoadBalancerResolversRequest) ([]upcloud.LoadBalancerResolver, error) {
	return nil, nil
}

func (m *Service) CreateLoadBalancerResolver(r *request.CreateLoadBalancerResolverRequest) (*upcloud.LoadBalancerResolver, error) {
	return nil, nil
}

func (m *Service) GetLoadBalancerResolver(r *request.GetLoadBalancerResolverRequest) (*upcloud.LoadBalancerResolver, error) {
	return nil, nil
}

func (m *Service) ModifyLoadBalancerResolver(r *request.ModifyLoadBalancerResolverRequest) (*upcloud.LoadBalancerResolver, error) {
	return nil, nil
}

func (m *Service) DeleteLoadBalancerResolver(r *request.DeleteLoadBalancerResolverRequest) error {
	return nil
}

func (m *Service) GetLoadBalancerPlans(r *request.GetLoadBalancerPlansRequest) ([]upcloud.LoadBalancerPlan, error) {
	return nil, nil
}

func (m *Service) GetLoadBalancerFrontends(r *request.GetLoadBalancerFrontendsRequest) ([]upcloud.LoadBalancerFrontend, error) {
	return nil, nil
}

func (m *Service) GetLoadBalancerFrontend(r *request.GetLoadBalancerFrontendRequest) (*upcloud.LoadBalancerFrontend, error) {
	return nil, nil
}

func (m *Service) CreateLoadBalancerFrontend(r *request.CreateLoadBalancerFrontendRequest) (*upcloud.LoadBalancerFrontend, error) {
	return nil, nil
}

func (m *Service) ModifyLoadBalancerFrontend(r *request.ModifyLoadBalancerFrontendRequest) (*upcloud.LoadBalancerFrontend, error) {
	return nil, nil
}

func (m *Service) DeleteLoadBalancerFrontend(r *request.DeleteLoadBalancerFrontendRequest) error {
	return nil
}

func (m *Service) GetLoadBalancerFrontendRules(r *request.GetLoadBalancerFrontendRulesRequest) ([]upcloud.LoadBalancerFrontendRule, error) {
	return nil, nil
}

func (m *Service) GetLoadBalancerFrontendRule(r *request.GetLoadBalancerFrontendRuleRequest) (*upcloud.LoadBalancerFrontendRule, error) {
	return nil, nil
}

func (m *Service) CreateLoadBalancerFrontendRule(r *request.CreateLoadBalancerFrontendRuleRequest) (*upcloud.LoadBalancerFrontendRule, error) {
	return nil, nil
}

func (m *Service) ModifyLoadBalancerFrontendRule(r *request.ModifyLoadBalancerFrontendRuleRequest) (*upcloud.LoadBalancerFrontendRule, error) {
	return nil, nil
}

func (m *Service) ReplaceLoadBalancerFrontendRule(r *request.ReplaceLoadBalancerFrontendRuleRequest) (*upcloud.LoadBalancerFrontendRule, error) {
	return nil, nil
}

func (m *Service) DeleteLoadBalancerFrontendRule(r *request.DeleteLoadBalancerFrontendRuleRequest) error {
	return nil
}

func (m *Service) GetLoadBalancerFrontendTLSConfigs(r *request.GetLoadBalancerFrontendTLSConfigsRequest) ([]upcloud.LoadBalancerFrontendTLSConfig, error) {
	return nil, nil
}

func (m *Service) GetLoadBalancerFrontendTLSConfig(r *request.GetLoadBalancerFrontendTLSConfigRequest) (*upcloud.LoadBalancerFrontendTLSConfig, error) {
	return nil, nil
}

func (m *Service) CreateLoadBalancerFrontendTLSConfig(r *request.CreateLoadBalancerFrontendTLSConfigRequest) (*upcloud.LoadBalancerFrontendTLSConfig, error) {
	return nil, nil
}

func (m *Service) ModifyLoadBalancerFrontendTLSConfig(r *request.ModifyLoadBalancerFrontendTLSConfigRequest) (*upcloud.LoadBalancerFrontendTLSConfig, error) {
	return nil, nil
}

func (m *Service) DeleteLoadBalancerFrontendTLSConfig(r *request.DeleteLoadBalancerFrontendTLSConfigRequest) error {
	return nil
}

func (m *Service) GetLoadBalancerCertificateBundles(r *request.GetLoadBalancerCertificateBundlesRequest) ([]upcloud.LoadBalancerCertificateBundle, error) {
	return nil, nil
}

func (m *Service) GetLoadBalancerCertificateBundle(r *request.GetLoadBalancerCertificateBundleRequest) (*upcloud.LoadBalancerCertificateBundle, error) {
	return nil, nil
}

func (m *Service) CreateLoadBalancerCertificateBundle(r *request.CreateLoadBalancerCertificateBundleRequest) (*upcloud.LoadBalancerCertificateBundle, error) {
	return nil, nil
}

func (m *Service) ModifyLoadBalancerCertificateBundle(r *request.ModifyLoadBalancerCertificateBundleRequest) (*upcloud.LoadBalancerCertificateBundle, error) {
	return nil, nil
}

func (m *Service) DeleteLoadBalancerCertificateBundle(r *request.DeleteLoadBalancerCertificateBundleRequest) error {
	return nil
}

func (m *Service) CreateKubernetesCluster(r *request.CreateKubernetesClusterRequest) (*upcloud.KubernetesCluster, error) {
	args := m.Called(r)
	if args[0] == nil {
		return nil, args.Error(1)
	}
	return args[0].(*upcloud.KubernetesCluster), args.Error(1)
}

func (m *Service) DeleteKubernetesCluster(r *request.DeleteKubernetesClusterRequest) error {
	args := m.Called(r)
	return args.Error(0)
}

func (m *Service) GetKubernetesCluster(r *request.GetKubernetesClusterRequest) (*upcloud.KubernetesCluster, error) {
	args := m.Called(r)
	if args[0] == nil {
		return nil, args.Error(1)
	}
	return args[0].(*upcloud.KubernetesCluster), args.Error(1)
}

func (m *Service) GetKubernetesClusters(r *request.GetKubernetesClustersRequest) ([]upcloud.KubernetesCluster, error) {
	args := m.Called(r)
	if args[0] == nil {
		return nil, args.Error(1)
	}
	return args[0].([]upcloud.KubernetesCluster), args.Error(1)
}

func (m *Service) GetKubernetesKubeconfig(r *request.GetKubernetesKubeconfigRequest) (string, error) {
	args := m.Called(r)
	if args[0] == nil {
		return "", args.Error(1)
	}
	return args[0].(string), args.Error(1)
}

func (m *Service) GetKubernetesPlans(r *request.GetKubernetesPlansRequest) ([]upcloud.KubernetesPlan, error) {
	return nil, nil
}

func (m *Service) GetKubernetesVersions(r *request.GetKubernetesVersionsRequest) ([]string, error) {
	return nil, nil
}

func (m *Service) WaitForKubernetesClusterState(r *request.WaitForKubernetesClusterStateRequest) (*upcloud.KubernetesCluster, error) {
	return nil, nil
}
