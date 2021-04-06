package all

import (
	"github.com/UpCloudLtd/cli/internal/commands"
	"github.com/UpCloudLtd/cli/internal/commands/account"
	"github.com/UpCloudLtd/cli/internal/commands/ipaddress"
	"github.com/UpCloudLtd/cli/internal/commands/network"

	// "github.com/UpCloudLtd/cli/internal/commands/ipaddress"
	// "github.com/UpCloudLtd/cli/internal/commands/network"
	// "github.com/UpCloudLtd/cli/internal/commands/networkinterface"
	// "github.com/UpCloudLtd/cli/internal/commands/router"
	"github.com/UpCloudLtd/cli/internal/commands/root"
	// "github.com/UpCloudLtd/cli/internal/commands/serverfirewall"
	// "github.com/UpCloudLtd/cli/internal/commands/serverstorage"
	// "github.com/UpCloudLtd/cli/internal/commands/storage"
	"github.com/UpCloudLtd/cli/internal/config"

	"github.com/spf13/cobra"
)

// BuildCommands is the main function that sets up the commands provided by upctl.
func BuildCommands(rootCmd *cobra.Command, conf *config.Config) {

	// Servers
	/*	serverCommand := commands.BuildCommand(server.BaseServerCommand(), rootCmd, conf)
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
	*/
	// // Server storage operations
	// serverStorageCommand := commands.BuildCommand(serverstorage.BaseServerStorageCommand(), serverCommand, conf)
	// commands.BuildCommand(serverstorage.AttachCommand(svc, svc), serverStorageCommand, conf)
	// commands.BuildCommand(serverstorage.DetachCommand(svc, svc), serverStorageCommand, conf)

	// // Server firewall operations
	// serverFirewallCommand := commands.BuildCommand(serverfirewall.BaseServerFirewallCommand(), serverCommand, conf)
	// commands.BuildCommand(serverfirewall.CreateCommand(svc, svc), serverFirewallCommand, conf)
	// commands.BuildCommand(serverfirewall.DeleteCommand(svc, svc), serverFirewallCommand, conf)
	// commands.BuildCommand(serverfirewall.ShowCommand(svc, svc), serverFirewallCommand, conf)

	// // Storages
	// storageCommand := commands.BuildCommand(storage.BaseStorageCommand(), mainCommand, conf)
	// commands.BuildCommand(storage.ListCommand(svc), storageCommand, conf)
	// commands.BuildCommand(storage.CreateCommand(svc), storageCommand, conf)
	// commands.BuildCommand(storage.ModifyCommand(svc), storageCommand, conf)
	// commands.BuildCommand(storage.CloneCommand(svc), storageCommand, conf)
	// commands.BuildCommand(storage.TemplatizeCommand(svc), storageCommand, conf)
	// commands.BuildCommand(storage.DeleteCommand(svc), storageCommand, conf)
	// commands.BuildCommand(storage.ImportCommand(svc), storageCommand, conf)
	// commands.BuildCommand(storage.ShowCommand(svc, svc), storageCommand, conf)

	// backupCommand := commands.BuildCommand(storage.BackupCommand(), storageCommand, conf)
	// commands.BuildCommand(storage.CreateBackupCommand(svc), backupCommand, conf)
	// commands.BuildCommand(storage.RestoreBackupCommand(svc), backupCommand, conf)

	// // IP Addresses
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
	// commands.BuildCommand(network.ModifyCommand(svc), networkCommand, conf)
	commands.BuildCommand(network.DeleteCommand(), networkCommand.Cobra(), conf)

	// // Network Interfaces
	// networkInterfaceCommand := commands.BuildCommand(networkinterface.BaseNetworkInterfaceCommand(), serverCommand, conf)
	// commands.BuildCommand(networkinterface.CreateCommand(svc, svc), networkInterfaceCommand, conf)
	// commands.BuildCommand(networkinterface.ModifyCommand(svc, svc), networkInterfaceCommand, conf)
	// commands.BuildCommand(networkinterface.DeleteCommand(svc, svc), networkInterfaceCommand, conf)

	// // Routers
	// routerCommand := commands.BuildCommand(router.BaseRouterCommand(), mainCommand, conf)
	// commands.BuildCommand(router.CreateCommand(svc), routerCommand, conf)
	// commands.BuildCommand(router.ListCommand(svc), routerCommand, conf)
	// commands.BuildCommand(router.ShowCommand(svc), routerCommand, conf)
	// commands.BuildCommand(router.ModifyCommand(svc), routerCommand, conf)
	// commands.BuildCommand(router.DeleteCommand(svc), routerCommand, conf)

	// Account
	accountCommand := commands.BuildCommand(account.BaseAccountCommand(), rootCmd, conf)
	commands.BuildCommand(account.ShowCommand(), accountCommand.Cobra(), conf)

	// Misc
	commands.BuildCommand(
		&root.CompletionCommand{
			BaseCommand: commands.New(
				"completion",
				"Generates shell completion",
			),
		}, rootCmd, conf,
	)
	commands.BuildCommand(
		&root.VersionCommand{
			BaseCommand: commands.New(
				"version",
				"Display software infomation",
			),
		}, rootCmd, conf,
	)
}
