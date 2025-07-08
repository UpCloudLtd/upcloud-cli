//nolint:nilnil // Here nil, nil returns and unused parameters are used in not-implemented methods required to satisfy an interface and it does not make sense to rename them when copying functions from our SDK
package mock

import (
	"context"
	"io"

	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud/request"
	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud/service"
	"github.com/stretchr/testify/mock"
)

// Service represents a mock upcloud API service
type Service struct {
	mock.Mock
}

func (m *Service) CreateToken(_ context.Context, r *request.CreateTokenRequest) (*upcloud.Token, error) {
	args := m.Called(r)
	if args[0] == nil {
		return nil, args.Error(1)
	}
	return args[0].(*upcloud.Token), args.Error(1)
}

func (m *Service) GetTokenDetails(_ context.Context, r *request.GetTokenDetailsRequest) (*upcloud.Token, error) {
	args := m.Called(r)
	if args[0] == nil {
		return nil, args.Error(1)
	}
	return args[0].(*upcloud.Token), args.Error(1)
}

func (m *Service) DeleteToken(_ context.Context, r *request.DeleteTokenRequest) error {
	args := m.Called(r)
	return args.Error(0)
}

func (m *Service) GetTokens(context.Context, *request.GetTokensRequest) (*upcloud.Tokens, error) {
	args := m.Called()
	if args[0] == nil {
		return nil, args.Error(1)
	}
	return args[0].(*upcloud.Tokens), args.Error(1)
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

// GetTimeZones implements service.Zones.GetPriceZones
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
	_ service.Server               = &Service{}
	_ service.Storage              = &Service{}
	_ service.Firewall             = &Service{}
	_ service.Network              = &Service{}
	_ service.IPAddress            = &Service{}
	_ service.Cloud                = &Service{}
	_ service.Account              = &Service{}
	_ service.LoadBalancer         = &Service{}
	_ service.Kubernetes           = &Service{}
	_ service.ServerGroup          = &Service{}
	_ service.ManagedObjectStorage = &Service{}
	_ service.Gateway              = &Service{}
	_ service.Token                = &Service{}
	_ service.AuditLog             = &Service{}
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

// RelocateServer implements service.Server.RelocateServer
func (m *Service) RelocateServer(_ context.Context, r *request.RelocateServerRequest) (*upcloud.ServerDetails, error) {
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
	return args.Error(0)
}

func (m *Service) CancelManagedDatabaseSession(_ context.Context, r *request.CancelManagedDatabaseSession) error {
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
	args := m.Called(r)
	if args[0] == nil {
		return nil, args.Error(1)
	}
	return args[0].(*upcloud.ManagedDatabase), args.Error(1)
}

func (m *Service) DeleteManagedDatabaseIndex(_ context.Context, r *request.DeleteManagedDatabaseIndexRequest) error {
	args := m.Called(r)
	if args[0] != nil {
		return args.Error(0)
	}
	return nil
}

func (m *Service) GetManagedDatabase(_ context.Context, r *request.GetManagedDatabaseRequest) (*upcloud.ManagedDatabase, error) {
	args := m.Called(r)
	if args[0] == nil {
		return nil, args.Error(1)
	}
	return args[0].(*upcloud.ManagedDatabase), args.Error(1)
}

func (m *Service) GetManagedDatabases(_ context.Context, r *request.GetManagedDatabasesRequest) ([]upcloud.ManagedDatabase, error) {
	args := m.Called(r)
	if args[0] == nil {
		return nil, args.Error(1)
	}
	return args[0].([]upcloud.ManagedDatabase), args.Error(1)
}

func (m *Service) GetManagedDatabaseAccessControl(_ context.Context, r *request.GetManagedDatabaseAccessControlRequest) (*upcloud.ManagedDatabaseAccessControl, error) {
	return nil, nil
}

func (m *Service) GetManagedDatabaseSessions(_ context.Context, r *request.GetManagedDatabaseSessionsRequest) (upcloud.ManagedDatabaseSessions, error) {
	args := m.Called(r)
	if args[0] == nil {
		return upcloud.ManagedDatabaseSessions{}, args.Error(1)
	}
	return args[0].(upcloud.ManagedDatabaseSessions), args.Error(1)
}

func (m *Service) GetManagedDatabaseIndices(_ context.Context, r *request.GetManagedDatabaseIndicesRequest) ([]upcloud.ManagedDatabaseIndex, error) {
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

func (m *Service) ModifyManagedDatabaseAccessControl(_ context.Context, r *request.ModifyManagedDatabaseAccessControlRequest) (*upcloud.ManagedDatabaseAccessControl, error) {
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

func (m *Service) GetLoadBalancerBackendTLSConfigs(_ context.Context, r *request.GetLoadBalancerBackendTLSConfigsRequest) ([]upcloud.LoadBalancerBackendTLSConfig, error) {
	return nil, nil
}

func (m *Service) GetLoadBalancerBackendTLSConfig(_ context.Context, r *request.GetLoadBalancerBackendTLSConfigRequest) (*upcloud.LoadBalancerBackendTLSConfig, error) {
	return nil, nil
}

func (m *Service) CreateLoadBalancerBackendTLSConfig(_ context.Context, r *request.CreateLoadBalancerBackendTLSConfigRequest) (*upcloud.LoadBalancerBackendTLSConfig, error) {
	return nil, nil
}

func (m *Service) ModifyLoadBalancerBackendTLSConfig(_ context.Context, r *request.ModifyLoadBalancerBackendTLSConfigRequest) (*upcloud.LoadBalancerBackendTLSConfig, error) {
	return nil, nil
}

func (m *Service) DeleteLoadBalancerBackendTLSConfig(_ context.Context, r *request.DeleteLoadBalancerBackendTLSConfigRequest) error {
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

func (m *Service) GetLoadBalancerDNSChallengeDomain(ctx context.Context, r *request.GetLoadBalancerDNSChallengeDomainRequest) (*upcloud.LoadBalancerDNSChallengeDomain, error) {
	args := m.Called(r)
	if args[0] == nil {
		return nil, args.Error(1)
	}
	return args[0].(*upcloud.LoadBalancerDNSChallengeDomain), args.Error(1)
}

func (m *Service) CreateKubernetesCluster(_ context.Context, r *request.CreateKubernetesClusterRequest) (*upcloud.KubernetesCluster, error) {
	args := m.Called(r)
	if args[0] == nil {
		return nil, args.Error(1)
	}
	return args[0].(*upcloud.KubernetesCluster), args.Error(1)
}

func (m *Service) ModifyKubernetesCluster(ctx context.Context, r *request.ModifyKubernetesClusterRequest) (*upcloud.KubernetesCluster, error) {
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

func (m *Service) GetKubernetesVersions(_ context.Context, r *request.GetKubernetesVersionsRequest) ([]upcloud.KubernetesVersion, error) {
	args := m.Called(r)
	if args[0] == nil {
		return nil, args.Error(1)
	}
	return args[0].([]upcloud.KubernetesVersion), args.Error(1)
}

func (m *Service) WaitForKubernetesClusterState(context.Context, *request.WaitForKubernetesClusterStateRequest) (*upcloud.KubernetesCluster, error) {
	return nil, nil
}

func (m *Service) GetKubernetesNodeGroups(ctx context.Context, r *request.GetKubernetesNodeGroupsRequest) ([]upcloud.KubernetesNodeGroup, error) {
	return nil, nil
}

func (m *Service) GetKubernetesNodeGroup(ctx context.Context, r *request.GetKubernetesNodeGroupRequest) (*upcloud.KubernetesNodeGroupDetails, error) {
	args := m.Called(r)
	if args[0] == nil {
		return nil, args.Error(1)
	}
	return args[0].(*upcloud.KubernetesNodeGroupDetails), args.Error(1)
}

func (m *Service) CreateKubernetesNodeGroup(ctx context.Context, r *request.CreateKubernetesNodeGroupRequest) (*upcloud.KubernetesNodeGroup, error) {
	args := m.Called(r)
	if args[0] == nil {
		return nil, args.Error(1)
	}
	return args[0].(*upcloud.KubernetesNodeGroup), args.Error(1)
}

func (m *Service) ModifyKubernetesNodeGroup(ctx context.Context, r *request.ModifyKubernetesNodeGroupRequest) (*upcloud.KubernetesNodeGroup, error) {
	args := m.Called(r)
	if args[0] == nil {
		return nil, args.Error(1)
	}
	return args[0].(*upcloud.KubernetesNodeGroup), args.Error(1)
}

func (m *Service) WaitForKubernetesNodeGroupState(ctx context.Context, r *request.WaitForKubernetesNodeGroupStateRequest) (*upcloud.KubernetesNodeGroup, error) {
	args := m.Called(r)
	if args[0] == nil {
		return nil, args.Error(1)
	}
	return args[0].(*upcloud.KubernetesNodeGroup), args.Error(1)
}

func (m *Service) DeleteKubernetesNodeGroup(ctx context.Context, r *request.DeleteKubernetesNodeGroupRequest) error {
	args := m.Called(r)
	if args[0] == nil {
		return args.Error(0)
	}
	return nil
}

func (m *Service) DeleteKubernetesNodeGroupNode(ctx context.Context, r *request.DeleteKubernetesNodeGroupNodeRequest) error {
	args := m.Called(r)
	if args[0] == nil {
		return args.Error(0)
	}
	return nil
}

func (m *Service) GetKubernetesPlans(ctx context.Context, r *request.GetKubernetesPlansRequest) ([]upcloud.KubernetesPlan, error) {
	args := m.Called(r)
	if args[0] == nil {
		return nil, args.Error(1)
	}
	return args[0].([]upcloud.KubernetesPlan), args.Error(1)
}

func (m *Service) GetKubernetesClusterAvailableUpgrades(ctx context.Context, r *request.GetKubernetesClusterAvailableUpgradesRequest) (*upcloud.KubernetesClusterAvailableUpgrades, error) {
	args := m.Called(r)
	if args[0] == nil {
		return nil, args.Error(1)
	}
	return args[0].(*upcloud.KubernetesClusterAvailableUpgrades), args.Error(1)
}

func (m *Service) UpgradeKubernetesCluster(ctx context.Context, r *request.UpgradeKubernetesClusterRequest) (*upcloud.KubernetesClusterUpgrade, error) {
	args := m.Called(r)
	if args[0] == nil {
		return nil, args.Error(1)
	}
	return args[0].(*upcloud.KubernetesClusterUpgrade), args.Error(1)
}

func (m *Service) CreateServerGroup(ctx context.Context, r *request.CreateServerGroupRequest) (*upcloud.ServerGroup, error) {
	args := m.Called(r)
	if args[0] == nil {
		return nil, args.Error(1)
	}
	return args[0].(*upcloud.ServerGroup), args.Error(1)
}

func (m *Service) GetServerGroups(ctx context.Context, r *request.GetServerGroupsRequest) (upcloud.ServerGroups, error) {
	args := m.Called(r)
	if args[0] == nil {
		return nil, args.Error(1)
	}
	return args[0].(upcloud.ServerGroups), args.Error(1)
}

func (m *Service) GetServerGroup(ctx context.Context, r *request.GetServerGroupRequest) (*upcloud.ServerGroup, error) {
	args := m.Called(r)
	if args[0] == nil {
		return nil, args.Error(1)
	}
	return args[0].(*upcloud.ServerGroup), args.Error(1)
}

func (m *Service) ModifyServerGroup(ctx context.Context, r *request.ModifyServerGroupRequest) (*upcloud.ServerGroup, error) {
	args := m.Called(r)
	if args[0] == nil {
		return nil, args.Error(1)
	}
	return args[0].(*upcloud.ServerGroup), args.Error(1)
}

func (m *Service) DeleteServerGroup(ctx context.Context, r *request.DeleteServerGroupRequest) error {
	args := m.Called(r)
	if args[0] == nil {
		return args.Error(0)
	}
	return nil
}

func (m *Service) AddServerToServerGroup(ctx context.Context, r *request.AddServerToServerGroupRequest) error {
	return nil
}

func (m *Service) RemoveServerFromServerGroup(ctx context.Context, r *request.RemoveServerFromServerGroupRequest) error {
	return nil
}

func (m *Service) GetManagedObjectStorageRegions(ctx context.Context, r *request.GetManagedObjectStorageRegionsRequest) ([]upcloud.ManagedObjectStorageRegion, error) {
	return nil, nil
}

func (m *Service) GetManagedObjectStorageRegion(ctx context.Context, r *request.GetManagedObjectStorageRegionRequest) (*upcloud.ManagedObjectStorageRegion, error) {
	return nil, nil
}

func (m *Service) CreateManagedObjectStorage(ctx context.Context, r *request.CreateManagedObjectStorageRequest) (*upcloud.ManagedObjectStorage, error) {
	return nil, nil
}

func (m *Service) GetManagedObjectStorages(ctx context.Context, r *request.GetManagedObjectStoragesRequest) ([]upcloud.ManagedObjectStorage, error) {
	args := m.Called(r)
	if args[0] == nil {
		return nil, args.Error(1)
	}
	return args[0].([]upcloud.ManagedObjectStorage), args.Error(1)
}

func (m *Service) GetManagedObjectStorage(ctx context.Context, r *request.GetManagedObjectStorageRequest) (*upcloud.ManagedObjectStorage, error) {
	return nil, nil
}

func (m *Service) ReplaceManagedObjectStorage(ctx context.Context, r *request.ReplaceManagedObjectStorageRequest) (*upcloud.ManagedObjectStorage, error) {
	return nil, nil
}

func (m *Service) ModifyManagedObjectStorage(ctx context.Context, r *request.ModifyManagedObjectStorageRequest) (*upcloud.ManagedObjectStorage, error) {
	return nil, nil
}

func (m *Service) DeleteManagedObjectStorage(ctx context.Context, r *request.DeleteManagedObjectStorageRequest) error {
	args := m.Called(r)
	if args[0] == nil {
		return args.Error(0)
	}
	return nil
}

func (m *Service) GetManagedObjectStorageMetrics(ctx context.Context, r *request.GetManagedObjectStorageMetricsRequest) (*upcloud.ManagedObjectStorageMetrics, error) {
	return nil, nil
}

func (m *Service) CreateManagedObjectStorageBucket(ctx context.Context, r *request.CreateManagedObjectStorageBucketRequest) (upcloud.ManagedObjectStorageBucketMetrics, error) {
	args := m.Called(r)
	if args[0] == nil {
		return upcloud.ManagedObjectStorageBucketMetrics{}, args.Error(1)
	}
	return args[0].(upcloud.ManagedObjectStorageBucketMetrics), args.Error(1)
}

func (m *Service) DeleteManagedObjectStorageBucket(ctx context.Context, r *request.DeleteManagedObjectStorageBucketRequest) error {
	return m.Called(r).Error(0)
}

func (m *Service) GetManagedObjectStorageBucketMetrics(ctx context.Context, r *request.GetManagedObjectStorageBucketMetricsRequest) ([]upcloud.ManagedObjectStorageBucketMetrics, error) {
	return nil, nil
}

func (m *Service) CreateManagedObjectStorageNetwork(ctx context.Context, r *request.CreateManagedObjectStorageNetworkRequest) (*upcloud.ManagedObjectStorageNetwork, error) {
	return nil, nil
}

func (m *Service) GetManagedObjectStorageNetworks(ctx context.Context, r *request.GetManagedObjectStorageNetworksRequest) ([]upcloud.ManagedObjectStorageNetwork, error) {
	return nil, nil
}

func (m *Service) GetManagedObjectStorageNetwork(ctx context.Context, r *request.GetManagedObjectStorageNetworkRequest) (*upcloud.ManagedObjectStorageNetwork, error) {
	return nil, nil
}

func (m *Service) DeleteManagedObjectStorageNetwork(ctx context.Context, r *request.DeleteManagedObjectStorageNetworkRequest) error {
	return nil
}

func (m *Service) CreateManagedObjectStorageUser(ctx context.Context, r *request.CreateManagedObjectStorageUserRequest) (*upcloud.ManagedObjectStorageUser, error) {
	return nil, nil
}

func (m *Service) GetManagedObjectStorageUsers(ctx context.Context, r *request.GetManagedObjectStorageUsersRequest) ([]upcloud.ManagedObjectStorageUser, error) {
	return nil, nil
}

func (m *Service) GetManagedObjectStorageUser(ctx context.Context, r *request.GetManagedObjectStorageUserRequest) (*upcloud.ManagedObjectStorageUser, error) {
	return nil, nil
}

func (m *Service) DeleteManagedObjectStorageUser(ctx context.Context, r *request.DeleteManagedObjectStorageUserRequest) error {
	return nil
}

func (m *Service) CreateManagedObjectStorageUserAccessKey(ctx context.Context, r *request.CreateManagedObjectStorageUserAccessKeyRequest) (*upcloud.ManagedObjectStorageUserAccessKey, error) {
	return nil, nil
}

func (m *Service) GetManagedObjectStorageUserAccessKeys(ctx context.Context, r *request.GetManagedObjectStorageUserAccessKeysRequest) ([]upcloud.ManagedObjectStorageUserAccessKey, error) {
	return nil, nil
}

func (m *Service) GetManagedObjectStorageUserAccessKey(ctx context.Context, r *request.GetManagedObjectStorageUserAccessKeyRequest) (*upcloud.ManagedObjectStorageUserAccessKey, error) {
	return nil, nil
}

func (m *Service) ModifyManagedObjectStorageUserAccessKey(ctx context.Context, r *request.ModifyManagedObjectStorageUserAccessKeyRequest) (*upcloud.ManagedObjectStorageUserAccessKey, error) {
	return nil, nil
}

func (m *Service) DeleteManagedObjectStorageUserAccessKey(ctx context.Context, r *request.DeleteManagedObjectStorageUserAccessKeyRequest) error {
	return nil
}

func (m *Service) WaitForManagedObjectStorageOperationalState(ctx context.Context, r *request.WaitForManagedObjectStorageOperationalStateRequest) (*upcloud.ManagedObjectStorage, error) {
	return nil, nil
}

func (m *Service) WaitForManagedObjectStorageDeletion(ctx context.Context, r *request.WaitForManagedObjectStorageDeletionRequest) error {
	return nil
}

func (m *Service) WaitForManagedObjectStorageBucketDeletion(ctx context.Context, r *request.WaitForManagedObjectStorageBucketDeletionRequest) error {
	return m.Called(r).Error(0)
}

func (m *Service) CreateManagedObjectStoragePolicy(ctx context.Context, r *request.CreateManagedObjectStoragePolicyRequest) (*upcloud.ManagedObjectStoragePolicy, error) {
	return nil, nil
}

func (m *Service) GetManagedObjectStoragePolicies(ctx context.Context, r *request.GetManagedObjectStoragePoliciesRequest) ([]upcloud.ManagedObjectStoragePolicy, error) {
	return nil, nil
}

func (m *Service) GetManagedObjectStoragePolicy(ctx context.Context, r *request.GetManagedObjectStoragePolicyRequest) (*upcloud.ManagedObjectStoragePolicy, error) {
	return nil, nil
}

func (m *Service) DeleteManagedObjectStoragePolicy(ctx context.Context, r *request.DeleteManagedObjectStoragePolicyRequest) error {
	return nil
}

func (m *Service) AttachManagedObjectStorageUserPolicy(ctx context.Context, r *request.AttachManagedObjectStorageUserPolicyRequest) error {
	return nil
}

func (m *Service) GetManagedObjectStorageUserPolicies(ctx context.Context, r *request.GetManagedObjectStorageUserPoliciesRequest) ([]upcloud.ManagedObjectStorageUserPolicy, error) {
	return nil, nil
}

func (m *Service) DetachManagedObjectStorageUserPolicy(ctx context.Context, r *request.DetachManagedObjectStorageUserPolicyRequest) error {
	return nil
}

func (m *Service) CreateManagedObjectStorageCustomDomain(ctx context.Context, r *request.CreateManagedObjectStorageCustomDomainRequest) error {
	return m.Called(r).Error(0)
}

func (m *Service) GetManagedObjectStorageCustomDomains(ctx context.Context, r *request.GetManagedObjectStorageCustomDomainsRequest) ([]upcloud.ManagedObjectStorageCustomDomain, error) {
	args := m.Called(r)
	if args[0] == nil {
		return nil, args.Error(1)
	}
	return args[0].([]upcloud.ManagedObjectStorageCustomDomain), args.Error(1)
}

func (m *Service) GetManagedObjectStorageCustomDomain(ctx context.Context, r *request.GetManagedObjectStorageCustomDomainRequest) (*upcloud.ManagedObjectStorageCustomDomain, error) {
	args := m.Called(r)
	if args[0] == nil {
		return nil, args.Error(1)
	}
	return args[0].(*upcloud.ManagedObjectStorageCustomDomain), args.Error(1)
}

func (m *Service) ModifyManagedObjectStorageCustomDomain(ctx context.Context, r *request.ModifyManagedObjectStorageCustomDomainRequest) (*upcloud.ManagedObjectStorageCustomDomain, error) {
	args := m.Called(r)
	if args[0] == nil {
		return nil, args.Error(1)
	}
	return args[0].(*upcloud.ManagedObjectStorageCustomDomain), args.Error(1)
}

func (m *Service) DeleteManagedObjectStorageCustomDomain(ctx context.Context, r *request.DeleteManagedObjectStorageCustomDomainRequest) error {
	return m.Called(r).Error(0)
}

func (m *Service) GetPermissions(ctx context.Context, r *request.GetPermissionsRequest) (upcloud.Permissions, error) {
	args := m.Called(r)
	if args[0] == nil {
		return nil, args.Error(1)
	}
	return args[0].(upcloud.Permissions), args.Error(1)
}

func (m *Service) GrantPermission(ctx context.Context, r *request.GrantPermissionRequest) (*upcloud.Permission, error) {
	args := m.Called(r)
	if args[0] == nil {
		return nil, args.Error(1)
	}
	return args[0].(*upcloud.Permission), args.Error(1)
}

func (m *Service) RevokePermission(ctx context.Context, r *request.RevokePermissionRequest) error {
	args := m.Called(r)
	if args[0] == nil {
		return args.Error(0)
	}
	return nil
}

// GetGateways implements service.Gateway.GetGateways
func (m *Service) GetGateways(_ context.Context, f ...request.QueryFilter) ([]upcloud.Gateway, error) {
	args := m.Called(f)
	if args[0] == nil {
		return nil, args.Error(1)
	}
	return args[0].([]upcloud.Gateway), args.Error(1)
}

// DeleteGateway implements service.Gateway.DeleteGateway
func (m *Service) DeleteGateway(_ context.Context, r *request.DeleteGatewayRequest) error {
	args := m.Called(r)
	return args.Error(0)
}

// CreateGateway implements service.Gateway.CreateGateway
func (m *Service) CreateGateway(_ context.Context, r *request.CreateGatewayRequest) (*upcloud.Gateway, error) {
	args := m.Called(r)
	if args[0] == nil {
		return nil, args.Error(1)
	}
	return args[0].(*upcloud.Gateway), args.Error(1)
}

// GetGateway implements service.Gateway.GetGateway
func (m *Service) GetGateway(_ context.Context, r *request.GetGatewayRequest) (*upcloud.Gateway, error) {
	args := m.Called(r)
	if args[0] == nil {
		return nil, args.Error(1)
	}
	return args[0].(*upcloud.Gateway), args.Error(1)
}

// ModifyGateway implements service.Gateway.ModifyGateway
func (m *Service) ModifyGateway(_ context.Context, r *request.ModifyGatewayRequest) (*upcloud.Gateway, error) {
	args := m.Called(r)
	if args[0] == nil {
		return nil, args.Error(1)
	}
	return args[0].(*upcloud.Gateway), args.Error(1)
}

func (m *Service) GetGatewayPlans(ctx context.Context) ([]upcloud.GatewayPlan, error) {
	return nil, nil
}

func (m *Service) GetGatewayMetrics(ctx context.Context, r *request.GetGatewayMetricsRequest) (*upcloud.GatewayMetrics, error) {
	return nil, nil
}

func (m *Service) GetGatewayConnections(ctx context.Context, r *request.GetGatewayConnectionsRequest) ([]upcloud.GatewayConnection, error) {
	return nil, nil
}

func (m *Service) GetGatewayConnection(ctx context.Context, r *request.GetGatewayConnectionRequest) (*upcloud.GatewayConnection, error) {
	return nil, nil
}

func (m *Service) CreateGatewayConnection(ctx context.Context, r *request.CreateGatewayConnectionRequest) (*upcloud.GatewayConnection, error) {
	return nil, nil
}

func (m *Service) ModifyGatewayConnection(ctx context.Context, r *request.ModifyGatewayConnectionRequest) (*upcloud.GatewayConnection, error) {
	return nil, nil
}

func (m *Service) DeleteGatewayConnection(ctx context.Context, r *request.DeleteGatewayConnectionRequest) error {
	return nil
}

func (m *Service) GetGatewayConnectionTunnels(ctx context.Context, r *request.GetGatewayConnectionTunnelsRequest) ([]upcloud.GatewayTunnel, error) {
	return nil, nil
}

func (m *Service) GetGatewayConnectionTunnel(ctx context.Context, r *request.GetGatewayConnectionTunnelRequest) (*upcloud.GatewayTunnel, error) {
	return nil, nil
}

func (m *Service) CreateGatewayConnectionTunnel(ctx context.Context, r *request.CreateGatewayConnectionTunnelRequest) (*upcloud.GatewayTunnel, error) {
	return nil, nil
}

func (m *Service) DeleteGatewayConnectionTunnel(ctx context.Context, r *request.DeleteGatewayConnectionTunnelRequest) error {
	return nil
}

func (m *Service) GetNetworkPeerings(ctx context.Context, f ...request.QueryFilter) (upcloud.NetworkPeerings, error) {
	return nil, nil
}

func (m *Service) GetNetworkPeering(ctx context.Context, r *request.GetNetworkPeeringRequest) (*upcloud.NetworkPeering, error) {
	return nil, nil
}

func (m *Service) CreateNetworkPeering(ctx context.Context, r *request.CreateNetworkPeeringRequest) (*upcloud.NetworkPeering, error) {
	return nil, nil
}

func (m *Service) ModifyNetworkPeering(ctx context.Context, r *request.ModifyNetworkPeeringRequest) (*upcloud.NetworkPeering, error) {
	return nil, nil
}

func (m *Service) DeleteNetworkPeering(ctx context.Context, r *request.DeleteNetworkPeeringRequest) error {
	args := m.Called(r)
	return args.Error(0)
}

func (m *Service) WaitForNetworkPeeringState(ctx context.Context, r *request.WaitForNetworkPeeringStateRequest) (*upcloud.NetworkPeering, error) {
	return nil, nil
}

func (m *Service) GetHosts(ctx context.Context) (*upcloud.Hosts, error) {
	return nil, nil
}

func (m *Service) GetHostDetails(ctx context.Context, r *request.GetHostDetailsRequest) (*upcloud.Host, error) {
	return nil, nil
}

func (m *Service) ModifyHost(ctx context.Context, r *request.ModifyHostRequest) (*upcloud.Host, error) {
	return nil, nil
}

func (m *Service) CreatePartnerAccount(ctx context.Context, r *request.CreatePartnerAccountRequest) (*upcloud.PartnerAccount, error) {
	args := m.Called(r)
	if args[0] == nil {
		return nil, args.Error(1)
	}
	return args[0].(*upcloud.PartnerAccount), args.Error(1)
}

func (m *Service) GetPartnerAccounts(ctx context.Context) ([]upcloud.PartnerAccount, error) {
	args := m.Called()
	if args[0] == nil {
		return nil, args.Error(1)
	}
	return args[0].([]upcloud.PartnerAccount), args.Error(1)
}

func (m *Service) GetTags(ctx context.Context) (*upcloud.Tags, error) {
	args := m.Called()
	if args[0] == nil {
		return nil, args.Error(1)
	}
	return args[0].(*upcloud.Tags), args.Error(1)
}

func (m *Service) CreateTag(ctx context.Context, r *request.CreateTagRequest) (*upcloud.Tag, error) {
	return nil, nil
}

func (m *Service) ModifyTag(ctx context.Context, r *request.ModifyTagRequest) (*upcloud.Tag, error) {
	return nil, nil
}

func (m *Service) DeleteTag(ctx context.Context, r *request.DeleteTagRequest) error {
	return nil
}

func (m *Service) TagServer(ctx context.Context, r *request.TagServerRequest) (*upcloud.ServerDetails, error) {
	return nil, nil
}

func (m *Service) UntagServer(ctx context.Context, r *request.UntagServerRequest) (*upcloud.ServerDetails, error) {
	return nil, nil
}

func (m *Service) WaitForLoadBalancerOperationalState(_ context.Context, r *request.WaitForLoadBalancerOperationalStateRequest) (*upcloud.LoadBalancer, error) {
	args := m.Called(r)
	if args[0] == nil {
		return nil, args.Error(1)
	}
	return args[0].(*upcloud.LoadBalancer), args.Error(1)
}

func (m *Service) WaitForLoadBalancerDeletion(_ context.Context, r *request.WaitForLoadBalancerDeletionRequest) error {
	args := m.Called(r)
	if args[0] == nil {
		return args.Error(0)
	}
	return nil
}

func (m *Service) GetDevicesAvailability(ctx context.Context) (*upcloud.DevicesAvailability, error) {
	args := m.Called()
	if args[0] == nil {
		return nil, args.Error(1)
	}
	return args[0].(*upcloud.DevicesAvailability), args.Error(1)
}

// AssignIPAddressToNetworkInterface implements service.Network.AssignIPAddressToNetworkInterface
func (m *Service) AssignIPAddressToNetworkInterface(_ context.Context, r *request.AssignIPAddressToNetworkInterfaceRequest) (*upcloud.IPAddress, error) {
	args := m.Called(r)
	if args[0] == nil {
		return nil, args.Error(1)
	}
	return args[0].(*upcloud.IPAddress), args.Error(1)
}

// DeleteIPAddressFromNetworkInterface implements service.Network.DeleteIPAddressFromNetworkInterface
func (m *Service) DeleteIPAddressFromNetworkInterface(_ context.Context, r *request.DeleteIPAddressFromNetworkInterfaceRequest) error {
	args := m.Called(r)
	return args.Error(0)
}

// GetBillingSummary implements service.Account.GetBillingSummary
func (m *Service) GetBillingSummary(_ context.Context, r *request.GetBillingSummaryRequest) (*upcloud.BillingSummary, error) {
	args := m.Called(r)
	if args[0] == nil {
		return nil, args.Error(1)
	}
	return args[0].(*upcloud.BillingSummary), args.Error(1)
}

// ExportAuditLog implements service.AuditLog.ExportAuditLog
func (m *Service) ExportAuditLog(_ context.Context, r *request.ExportAuditLogRequest) (io.ReadCloser, error) {
	args := m.Called(r)
	if args[0] == nil {
		return nil, args.Error(1)
	}
	return args[0].(io.ReadCloser), args.Error(1)
}
