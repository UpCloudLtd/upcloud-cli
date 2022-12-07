package all

import (
	"github.com/UpCloudLtd/upcloud-cli/v2/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v2/internal/commands/account"
	"github.com/UpCloudLtd/upcloud-cli/v2/internal/commands/database"
	databaseconnection "github.com/UpCloudLtd/upcloud-cli/v2/internal/commands/database/connection"
	"github.com/UpCloudLtd/upcloud-cli/v2/internal/commands/ipaddress"
	"github.com/UpCloudLtd/upcloud-cli/v2/internal/commands/kubernetes"
	"github.com/UpCloudLtd/upcloud-cli/v2/internal/commands/loadbalancer"
	"github.com/UpCloudLtd/upcloud-cli/v2/internal/commands/network"
	"github.com/UpCloudLtd/upcloud-cli/v2/internal/commands/root"
	"github.com/UpCloudLtd/upcloud-cli/v2/internal/commands/router"
	"github.com/UpCloudLtd/upcloud-cli/v2/internal/commands/server"
	serverfirewall "github.com/UpCloudLtd/upcloud-cli/v2/internal/commands/server/firewall"
	"github.com/UpCloudLtd/upcloud-cli/v2/internal/commands/server/networkinterface"
	serverstorage "github.com/UpCloudLtd/upcloud-cli/v2/internal/commands/server/storage"
	"github.com/UpCloudLtd/upcloud-cli/v2/internal/commands/storage"
	storagebackup "github.com/UpCloudLtd/upcloud-cli/v2/internal/commands/storage/backup"
	"github.com/UpCloudLtd/upcloud-cli/v2/internal/commands/zone"
	"github.com/UpCloudLtd/upcloud-cli/v2/internal/config"

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

	// Routers
	routerCommand := commands.BuildCommand(router.BaseRouterCommand(), rootCmd, conf)
	commands.BuildCommand(router.CreateCommand(), routerCommand.Cobra(), conf)
	commands.BuildCommand(router.ListCommand(), routerCommand.Cobra(), conf)
	commands.BuildCommand(router.ShowCommand(), routerCommand.Cobra(), conf)
	commands.BuildCommand(router.ModifyCommand(), routerCommand.Cobra(), conf)
	commands.BuildCommand(router.DeleteCommand(), routerCommand.Cobra(), conf)

	// Account
	accountCommand := commands.BuildCommand(account.BaseAccountCommand(), rootCmd, conf)
	commands.BuildCommand(account.ShowCommand(), accountCommand.Cobra(), conf)

	// Zone
	zoneCommand := commands.BuildCommand(zone.BaseZoneCommand(), rootCmd, conf)
	commands.BuildCommand(zone.ListCommand(), zoneCommand.Cobra(), conf)

	// Databases
	databaseCommand := commands.BuildCommand(database.BaseDatabaseCommand(), rootCmd, conf)
	commands.BuildCommand(database.ListCommand(), databaseCommand.Cobra(), conf)
	commands.BuildCommand(database.ShowCommand(), databaseCommand.Cobra(), conf)
	commands.BuildCommand(database.TypesCommand(), databaseCommand.Cobra(), conf)
	commands.BuildCommand(database.PlansCommand(), databaseCommand.Cobra(), conf)
	commands.BuildCommand(database.StartCommand(), databaseCommand.Cobra(), conf)
	commands.BuildCommand(database.StopCommand(), databaseCommand.Cobra(), conf)
	commands.BuildCommand(database.DeleteCommand(), databaseCommand.Cobra(), conf)

	// Database connections
	connectionsCommand := commands.BuildCommand(databaseconnection.BaseConnectionCommand(), databaseCommand.Cobra(), conf)
	commands.BuildCommand(databaseconnection.CancelCommand(), connectionsCommand.Cobra(), conf)
	commands.BuildCommand(databaseconnection.ListCommand(), connectionsCommand.Cobra(), conf)

	// LoadBalancers
	loadbalancerCommand := commands.BuildCommand(loadbalancer.BaseLoadBalancerCommand(), rootCmd, conf)
	commands.BuildCommand(loadbalancer.ListCommand(), loadbalancerCommand.Cobra(), conf)
	commands.BuildCommand(loadbalancer.ShowCommand(), loadbalancerCommand.Cobra(), conf)
	commands.BuildCommand(loadbalancer.DeleteCommand(), loadbalancerCommand.Cobra(), conf)
	commands.BuildCommand(loadbalancer.PlansCommand(), loadbalancerCommand.Cobra(), conf)

	// Kubernetes
	kubernetesCommand := commands.BuildCommand(kubernetes.BaseKubernetesCommand(), rootCmd, conf)
	commands.BuildCommand(kubernetes.CreateCommand(), kubernetesCommand.Cobra(), conf)
	commands.BuildCommand(kubernetes.DeleteCommand(), kubernetesCommand.Cobra(), conf)
	commands.BuildCommand(kubernetes.ListCommand(), kubernetesCommand.Cobra(), conf)
	commands.BuildCommand(kubernetes.ShowCommand(), kubernetesCommand.Cobra(), conf)
	commands.BuildCommand(kubernetes.VersionsCommand(), kubernetesCommand.Cobra(), conf)

	// Misc
	commands.BuildCommand(
		&root.CompletionCommand{
			BaseCommand: commands.New(
				"completion",
				"Generates shell completion",
				"upctl completion bash",
			),
		}, rootCmd, conf,
	)
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
