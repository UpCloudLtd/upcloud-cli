package base

import (
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands/account"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands/account/permissions"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands/account/token"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands/all"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands/auditlog"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands/database"
	databaseindex "github.com/UpCloudLtd/upcloud-cli/v3/internal/commands/database/index"
	databaseproperties "github.com/UpCloudLtd/upcloud-cli/v3/internal/commands/database/properties"
	databasesession "github.com/UpCloudLtd/upcloud-cli/v3/internal/commands/database/session"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands/gateway"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands/host"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands/ipaddress"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands/kubernetes"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands/kubernetes/nodegroup"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands/loadbalancer"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands/network"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands/networkpeering"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands/objectstorage"
	objectstorageAccesskey "github.com/UpCloudLtd/upcloud-cli/v3/internal/commands/objectstorage/accesskey"
	objectstoragebucket "github.com/UpCloudLtd/upcloud-cli/v3/internal/commands/objectstorage/bucket"
	objectstoragelabel "github.com/UpCloudLtd/upcloud-cli/v3/internal/commands/objectstorage/label"
	objectstoragenetwork "github.com/UpCloudLtd/upcloud-cli/v3/internal/commands/objectstorage/network"
	objectstorageuser "github.com/UpCloudLtd/upcloud-cli/v3/internal/commands/objectstorage/user"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands/partner"
	partneraccount "github.com/UpCloudLtd/upcloud-cli/v3/internal/commands/partner/account"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands/root"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands/router"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands/server"
	serverfirewall "github.com/UpCloudLtd/upcloud-cli/v3/internal/commands/server/firewall"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands/server/networkinterface"
	serverstorage "github.com/UpCloudLtd/upcloud-cli/v3/internal/commands/server/storage"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands/servergroup"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands/stack"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands/stack/dokku"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands/stack/starterkit"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands/stack/supabase"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands/storage"
	storagebackup "github.com/UpCloudLtd/upcloud-cli/v3/internal/commands/storage/backup"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands/zone"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands/zone/devices"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/config"

	"github.com/spf13/cobra"
)

// BuildCommands is the main function that sets up the commands provided by upctl.
func BuildCommands(rootCmd *cobra.Command, conf *config.Config) {
	// Servers
	serverCommand := commands.BuildCommand(server.BaseServerCommand(), rootCmd, conf)
	commands.BuildCommand(server.ListCommand(), serverCommand.Cobra(), conf)
	commands.BuildCommand(server.PlanListCommand(), serverCommand.Cobra(), conf)
	commands.BuildCommand(server.ShowCommand(), serverCommand.Cobra(), conf)
	commands.BuildCommand(server.StartCommand(), serverCommand.Cobra(), conf)
	commands.BuildCommand(server.RestartCommand(), serverCommand.Cobra(), conf)
	commands.BuildCommand(server.StopCommand(), serverCommand.Cobra(), conf)
	commands.BuildCommand(server.CreateCommand(), serverCommand.Cobra(), conf)
	commands.BuildCommand(server.ModifyCommand(), serverCommand.Cobra(), conf)
	commands.BuildCommand(server.LoadCommand(), serverCommand.Cobra(), conf)
	commands.BuildCommand(server.EjectCommand(), serverCommand.Cobra(), conf)
	commands.BuildCommand(server.DeleteCommand(), serverCommand.Cobra(), conf)
	commands.BuildCommand(server.RelocateCommand(), serverCommand.Cobra(), conf)

	// Server Network Interfaces
	networkInterfaceCommand := commands.BuildCommand(networkinterface.BaseNetworkInterfaceCommand(), serverCommand.Cobra(), conf)
	commands.BuildCommand(networkinterface.CreateCommand(), networkInterfaceCommand.Cobra(), conf)
	commands.BuildCommand(networkinterface.ModifyCommand(), networkInterfaceCommand.Cobra(), conf)
	commands.BuildCommand(networkinterface.DeleteCommand(), networkInterfaceCommand.Cobra(), conf)

	// Server storage operations
	serverStorageCommand := commands.BuildCommand(serverstorage.BaseServerStorageCommand(), serverCommand.Cobra(), conf)
	commands.BuildCommand(serverstorage.AttachCommand(), serverStorageCommand.Cobra(), conf)
	commands.BuildCommand(serverstorage.DetachCommand(), serverStorageCommand.Cobra(), conf)

	// Server firewall operations
	serverFirewallCommand := commands.BuildCommand(serverfirewall.BaseServerFirewallCommand(), serverCommand.Cobra(), conf)
	commands.BuildCommand(serverfirewall.CreateCommand(), serverFirewallCommand.Cobra(), conf)
	commands.BuildCommand(serverfirewall.DeleteCommand(), serverFirewallCommand.Cobra(), conf)
	commands.BuildCommand(serverfirewall.ShowCommand(), serverFirewallCommand.Cobra(), conf)

	// Storages
	storageCommand := commands.BuildCommand(storage.BaseStorageCommand(), rootCmd, conf)
	commands.BuildCommand(storage.ListCommand(), storageCommand.Cobra(), conf)
	commands.BuildCommand(storage.CreateCommand(), storageCommand.Cobra(), conf)
	commands.BuildCommand(storage.ModifyCommand(), storageCommand.Cobra(), conf)
	commands.BuildCommand(storage.CloneCommand(), storageCommand.Cobra(), conf)
	commands.BuildCommand(storage.TemplatizeCommand(), storageCommand.Cobra(), conf)
	commands.BuildCommand(storage.DeleteCommand(), storageCommand.Cobra(), conf)
	commands.BuildCommand(storage.ImportCommand(), storageCommand.Cobra(), conf)
	commands.BuildCommand(storage.ShowCommand(), storageCommand.Cobra(), conf)

	backupCommand := commands.BuildCommand(storagebackup.BackupCommand(), storageCommand.Cobra(), conf)
	commands.BuildCommand(storagebackup.CreateBackupCommand(), backupCommand.Cobra(), conf)
	commands.BuildCommand(storagebackup.RestoreBackupCommand(), backupCommand.Cobra(), conf)

	// IP Addresses
	ipAddressCommand := commands.BuildCommand(ipaddress.BaseIPAddressCommand(), rootCmd, conf)
	commands.BuildCommand(ipaddress.ListCommand(), ipAddressCommand.Cobra(), conf)
	commands.BuildCommand(ipaddress.ShowCommand(), ipAddressCommand.Cobra(), conf)
	commands.BuildCommand(ipaddress.ModifyCommand(), ipAddressCommand.Cobra(), conf)
	commands.BuildCommand(ipaddress.AssignCommand(), ipAddressCommand.Cobra(), conf)
	commands.BuildCommand(ipaddress.RemoveCommand(), ipAddressCommand.Cobra(), conf)

	// Networks
	networkCommand := commands.BuildCommand(network.BaseNetworkCommand(), rootCmd, conf)
	commands.BuildCommand(network.CreateCommand(), networkCommand.Cobra(), conf)
	commands.BuildCommand(network.ListCommand(), networkCommand.Cobra(), conf)
	commands.BuildCommand(network.ShowCommand(), networkCommand.Cobra(), conf)
	commands.BuildCommand(network.ModifyCommand(), networkCommand.Cobra(), conf)
	commands.BuildCommand(network.DeleteCommand(), networkCommand.Cobra(), conf)

	// Network peerings
	networkPeeringCommand := commands.BuildCommand(networkpeering.BaseNetworkPeeringCommand(), rootCmd, conf)
	commands.BuildCommand(networkpeering.ListCommand(), networkPeeringCommand.Cobra(), conf)
	commands.BuildCommand(networkpeering.DeleteCommand(), networkPeeringCommand.Cobra(), conf)
	commands.BuildCommand(networkpeering.DisableCommand(), networkPeeringCommand.Cobra(), conf)

	// Routers
	routerCommand := commands.BuildCommand(router.BaseRouterCommand(), rootCmd, conf)
	commands.BuildCommand(router.CreateCommand(), routerCommand.Cobra(), conf)
	commands.BuildCommand(router.ListCommand(), routerCommand.Cobra(), conf)
	commands.BuildCommand(router.ShowCommand(), routerCommand.Cobra(), conf)
	commands.BuildCommand(router.ModifyCommand(), routerCommand.Cobra(), conf)
	commands.BuildCommand(router.DeleteCommand(), routerCommand.Cobra(), conf)

	// Account
	accountCommand := commands.BuildCommand(account.BaseAccountCommand(), rootCmd, conf)
	commands.BuildCommand(account.LoginCommand(), accountCommand.Cobra(), conf)
	commands.BuildCommand(account.ShowCommand(), accountCommand.Cobra(), conf)
	commands.BuildCommand(account.ListCommand(), accountCommand.Cobra(), conf)
	commands.BuildCommand(account.DeleteCommand(), accountCommand.Cobra(), conf)

	// Account permissions
	permissionsCommand := commands.BuildCommand(permissions.BasePermissionsCommand(), accountCommand.Cobra(), conf)
	commands.BuildCommand(permissions.ListCommand(), permissionsCommand.Cobra(), conf)

	// Account token
	tokenCommand := commands.BuildCommand(token.BaseTokenCommand(), accountCommand.Cobra(), conf)
	commands.BuildCommand(token.CreateCommand(), tokenCommand.Cobra(), conf)
	commands.BuildCommand(token.ListCommand(), tokenCommand.Cobra(), conf)
	commands.BuildCommand(token.ShowCommand(), tokenCommand.Cobra(), conf)
	commands.BuildCommand(token.DeleteCommand(), tokenCommand.Cobra(), conf)

	// Zone
	zoneCommand := commands.BuildCommand(zone.BaseZoneCommand(), rootCmd, conf)
	commands.BuildCommand(zone.ListCommand(), zoneCommand.Cobra(), conf)
	devicesCommand := commands.BuildCommand(devices.DevicesCommand(), zoneCommand.Cobra(), conf)
	commands.BuildCommand(devices.ShowCommand(), devicesCommand.Cobra(), conf)

	// Databases
	databaseCommand := commands.BuildCommand(database.BaseDatabaseCommand(), rootCmd, conf)
	commands.BuildCommand(database.CreateCommand(), databaseCommand.Cobra(), conf)
	commands.BuildCommand(database.ListCommand(), databaseCommand.Cobra(), conf)
	commands.BuildCommand(database.ShowCommand(), databaseCommand.Cobra(), conf)
	commands.BuildCommand(database.TypesCommand(), databaseCommand.Cobra(), conf)
	commands.BuildCommand(database.PlansCommand(), databaseCommand.Cobra(), conf)
	commands.BuildCommand(database.StartCommand(), databaseCommand.Cobra(), conf)
	commands.BuildCommand(database.StopCommand(), databaseCommand.Cobra(), conf)
	commands.BuildCommand(database.DeleteCommand(), databaseCommand.Cobra(), conf)

	// Database sessions
	sessionsCommand := commands.BuildCommand(databasesession.BaseSessionCommand(), databaseCommand.Cobra(), conf)
	commands.BuildCommand(databasesession.CancelCommand(), sessionsCommand.Cobra(), conf)
	commands.BuildCommand(databasesession.ListCommand(), sessionsCommand.Cobra(), conf)

	// Database sessions
	propertiesCommand := commands.BuildCommand(databaseproperties.PropertiesCommand(), databaseCommand.Cobra(), conf)
	for _, i := range []struct {
		serviceName string
		serviceType string
	}{
		{serviceName: "MySQL", serviceType: "mysql"},
		{serviceName: "OpenSearch", serviceType: "opensearch"},
		{serviceName: "PostgreSQL", serviceType: "pg"},
		{serviceName: "Redis", serviceType: "redis"},
		{serviceName: "Valkey", serviceType: "valkey"},
	} {
		typeCommand := commands.BuildCommand(databaseproperties.DBTypeCommand(i.serviceType, i.serviceName), propertiesCommand.Cobra(), conf)
		commands.BuildCommand(databaseproperties.ShowCommand(i.serviceType, i.serviceName), typeCommand.Cobra(), conf)
	}

	// Database indices
	indexCommand := commands.BuildCommand(databaseindex.BaseIndexCommand(), databaseCommand.Cobra(), conf)
	commands.BuildCommand(databaseindex.DeleteCommand(), indexCommand.Cobra(), conf)
	commands.BuildCommand(databaseindex.ListCommand(), indexCommand.Cobra(), conf)

	// LoadBalancers
	loadbalancerCommand := commands.BuildCommand(loadbalancer.BaseLoadBalancerCommand(), rootCmd, conf)
	commands.BuildCommand(loadbalancer.ListCommand(), loadbalancerCommand.Cobra(), conf)
	commands.BuildCommand(loadbalancer.ShowCommand(), loadbalancerCommand.Cobra(), conf)
	commands.BuildCommand(loadbalancer.DeleteCommand(), loadbalancerCommand.Cobra(), conf)
	commands.BuildCommand(loadbalancer.PlansCommand(), loadbalancerCommand.Cobra(), conf)

	// Kubernetes
	kubernetesCommand := commands.BuildCommand(kubernetes.BaseKubernetesCommand(), rootCmd, conf)
	commands.BuildCommand(kubernetes.CreateCommand(), kubernetesCommand.Cobra(), conf)
	commands.BuildCommand(kubernetes.ModifyCommand(), kubernetesCommand.Cobra(), conf)
	commands.BuildCommand(kubernetes.ConfigCommand(), kubernetesCommand.Cobra(), conf)
	commands.BuildCommand(kubernetes.DeleteCommand(), kubernetesCommand.Cobra(), conf)
	commands.BuildCommand(kubernetes.ListCommand(), kubernetesCommand.Cobra(), conf)
	commands.BuildCommand(kubernetes.ShowCommand(), kubernetesCommand.Cobra(), conf)
	commands.BuildCommand(kubernetes.VersionsCommand(), kubernetesCommand.Cobra(), conf)
	commands.BuildCommand(kubernetes.PlansCommand(), kubernetesCommand.Cobra(), conf)

	// Kubernetes nodegroup operations
	nodeGroupCommand := commands.BuildCommand(nodegroup.BaseNodeGroupCommand(), kubernetesCommand.Cobra(), conf)
	commands.BuildCommand(nodegroup.CreateCommand(), nodeGroupCommand.Cobra(), conf)
	commands.BuildCommand(nodegroup.ScaleCommand(), nodeGroupCommand.Cobra(), conf)
	commands.BuildCommand(nodegroup.ShowCommand(), nodeGroupCommand.Cobra(), conf)
	commands.BuildCommand(nodegroup.DeleteCommand(), nodeGroupCommand.Cobra(), conf)

	// Server group operations
	serverGroupCommand := commands.BuildCommand(servergroup.BaseServergroupCommand(), rootCmd, conf)
	commands.BuildCommand(servergroup.CreateCommand(), serverGroupCommand.Cobra(), conf)
	commands.BuildCommand(servergroup.DeleteCommand(), serverGroupCommand.Cobra(), conf)
	commands.BuildCommand(servergroup.ListCommand(), serverGroupCommand.Cobra(), conf)
	commands.BuildCommand(servergroup.ModifyCommand(), serverGroupCommand.Cobra(), conf)
	commands.BuildCommand(servergroup.ShowCommand(), serverGroupCommand.Cobra(), conf)

	// Managed object storage operations
	objectStorageCommand := commands.BuildCommand(objectstorage.BaseobjectstorageCommand(), rootCmd, conf)
	commands.BuildCommand(objectstorage.CreateCommand(), objectStorageCommand.Cobra(), conf)
	commands.BuildCommand(objectstorage.DeleteCommand(), objectStorageCommand.Cobra(), conf)
	commands.BuildCommand(objectstorage.ListCommand(), objectStorageCommand.Cobra(), conf)
	commands.BuildCommand(objectstorage.ShowCommand(), objectStorageCommand.Cobra(), conf)
	commands.BuildCommand(objectstorage.RegionsCommand(), objectStorageCommand.Cobra(), conf)

	// Object storage user management
	userCommand := commands.BuildCommand(objectstorageuser.BaseUserCommand(), objectStorageCommand.Cobra(), conf)
	commands.BuildCommand(objectstorageuser.CreateCommand(), userCommand.Cobra(), conf)
	commands.BuildCommand(objectstorageuser.DeleteCommand(), userCommand.Cobra(), conf)
	commands.BuildCommand(objectstorageuser.ListCommand(), userCommand.Cobra(), conf)

	// Object storage access key management
	accessKeyCommand := commands.BuildCommand(objectstorageAccesskey.BaseAccessKeyCommand(), objectStorageCommand.Cobra(), conf)
	commands.BuildCommand(objectstorageAccesskey.CreateCommand(), accessKeyCommand.Cobra(), conf)
	commands.BuildCommand(objectstorageAccesskey.DeleteCommand(), accessKeyCommand.Cobra(), conf)
	commands.BuildCommand(objectstorageAccesskey.ListCommand(), accessKeyCommand.Cobra(), conf)

	// Object storage network management
	objectStorageNetworkCommand := commands.BuildCommand(objectstoragenetwork.BaseNetworkCommand(), objectStorageCommand.Cobra(), conf)
	commands.BuildCommand(objectstoragenetwork.AttachCommand(), objectStorageNetworkCommand.Cobra(), conf)
	commands.BuildCommand(objectstoragenetwork.DetachCommand(), objectStorageNetworkCommand.Cobra(), conf)
	commands.BuildCommand(objectstoragenetwork.ListCommand(), objectStorageNetworkCommand.Cobra(), conf)

	// Object storage bucket management
	bucketCommand := commands.BuildCommand(objectstoragebucket.BaseBucketCommand(), objectStorageCommand.Cobra(), conf)
	commands.BuildCommand(objectstoragebucket.CreateCommand(), bucketCommand.Cobra(), conf)
	commands.BuildCommand(objectstoragebucket.DeleteCommand(), bucketCommand.Cobra(), conf)
	commands.BuildCommand(objectstoragebucket.ListCommand(), bucketCommand.Cobra(), conf)

	// Object storage label management
	labelCommand := commands.BuildCommand(objectstoragelabel.BaseLabelCommand(), objectStorageCommand.Cobra(), conf)
	commands.BuildCommand(objectstoragelabel.AddCommand(), labelCommand.Cobra(), conf)
	commands.BuildCommand(objectstoragelabel.RemoveCommand(), labelCommand.Cobra(), conf)
	commands.BuildCommand(objectstoragelabel.ListCommand(), labelCommand.Cobra(), conf)

	// Network Gateway operations
	gatewayCommand := commands.BuildCommand(gateway.BaseGatewayCommand(), rootCmd, conf)
	commands.BuildCommand(gateway.DeleteCommand(), gatewayCommand.Cobra(), conf)
	commands.BuildCommand(gateway.ListCommand(), gatewayCommand.Cobra(), conf)
	commands.BuildCommand(gateway.PlansCommand(), gatewayCommand.Cobra(), conf)

	// Host operations
	hostCommand := commands.BuildCommand(host.BaseHostCommand(), rootCmd, conf)
	commands.BuildCommand(host.ListCommand(), hostCommand.Cobra(), conf)

	// Partner API
	partnerCommand := commands.BuildCommand(partner.BasePartnerCommand(), rootCmd, conf)
	partnerAccountCommand := commands.BuildCommand(partneraccount.BaseAccountCommand(), partnerCommand.Cobra(), conf)
	commands.BuildCommand(partneraccount.CreateCommand(), partnerAccountCommand.Cobra(), conf)
	commands.BuildCommand(partneraccount.ListCommand(), partnerAccountCommand.Cobra(), conf)

	// Audit log operations
	auditlogCommand := commands.BuildCommand(auditlog.BaseAuditLogCommand(), rootCmd, conf)
	commands.BuildCommand(auditlog.ExportCommand(), auditlogCommand.Cobra(), conf)

	// Operations for managing all resources at once
	allCommand := commands.BuildCommand(all.BaseAllCommand(), rootCmd, conf)
	commands.BuildCommand(all.PurgeCommand(), allCommand.Cobra(), conf)
	commands.BuildCommand(all.ListCommand(), allCommand.Cobra(), conf)

	// Stack operations
	stackCommand := commands.BuildCommand(stack.BaseStackCommand(), rootCmd, conf)
	stackDeployCommand := commands.BuildCommand(stack.DeployCommand(), stackCommand.Cobra(), conf)
	stackDestroyCommand := commands.BuildCommand(stack.DestroyCommand(), stackCommand.Cobra(), conf)
	commands.BuildCommand(supabase.DeploySupabaseCommand(), stackDeployCommand.Cobra(), conf)
	commands.BuildCommand(dokku.DeployDokkuCommand(), stackDeployCommand.Cobra(), conf)
	commands.BuildCommand(starterkit.DeployStarterKitCommand(), stackDeployCommand.Cobra(), conf)
	commands.BuildCommand(supabase.DestroySupabaseCommand(), stackDestroyCommand.Cobra(), conf)
	commands.BuildCommand(dokku.DestroyDokkuCommand(), stackDestroyCommand.Cobra(), conf)
	commands.BuildCommand(starterkit.DestroyStarterKitCommand(), stackDestroyCommand.Cobra(), conf)

	// Misc
	commands.BuildCommand(
		&root.VersionCommand{
			BaseCommand: commands.New(
				"version",
				"Display software information",
				"upctl version",
			),
		}, rootCmd, conf,
	)
}
