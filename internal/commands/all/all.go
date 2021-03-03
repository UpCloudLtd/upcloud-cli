package all

import (
	"github.com/UpCloudLtd/cli/internal/commands"
	"github.com/UpCloudLtd/cli/internal/commands/account"
	"github.com/UpCloudLtd/cli/internal/commands/ipaddress"
	"github.com/UpCloudLtd/cli/internal/commands/network"
	"github.com/UpCloudLtd/cli/internal/commands/networkinterface"
	"github.com/UpCloudLtd/cli/internal/commands/router"
	"github.com/UpCloudLtd/cli/internal/commands/server"
	"github.com/UpCloudLtd/cli/internal/commands/serverstorage"
	"github.com/UpCloudLtd/cli/internal/commands/storage"
	"github.com/UpCloudLtd/cli/internal/commands/serverfirewall"
	"github.com/UpCloudLtd/cli/internal/config"
	"github.com/UpCloudLtd/cli/internal/upapi"
)

// BuildCommands is the main function that sets up the commands provided by upctl.
func BuildCommands(mainCommand commands.Command, mainConfig *config.Config) {
	cfgFn := func() *config.Config { return config.New(mainConfig.Viper()) }
	svc := upapi.Service(cfgFn())

	// Servers
	serverCommand := commands.BuildCommand(server.BaseServerCommand(), mainCommand, cfgFn())
	commands.BuildCommand(server.ListCommand(svc), serverCommand, cfgFn())
	commands.BuildCommand(server.PlanListCommand(), serverCommand, cfgFn())
	commands.BuildCommand(server.ShowCommand(svc, svc), serverCommand, cfgFn())
	commands.BuildCommand(server.StartCommand(svc), serverCommand, cfgFn())
	commands.BuildCommand(server.RestartCommand(svc), serverCommand, cfgFn())
	commands.BuildCommand(server.StopCommand(svc), serverCommand, cfgFn())
	commands.BuildCommand(server.CreateCommand(svc, svc), serverCommand, cfgFn())
	commands.BuildCommand(server.ModifyCommand(svc), serverCommand, cfgFn())
	commands.BuildCommand(server.LoadCommand(svc, svc), serverCommand, cfgFn())
	commands.BuildCommand(server.EjectCommand(svc, svc), serverCommand, cfgFn())
	commands.BuildCommand(server.DeleteCommand(svc), serverCommand, cfgFn())

	// Server storage operations
	serverStorageCommand := commands.BuildCommand(serverstorage.BaseServerStorageCommand(), serverCommand, cfgFn())
	commands.BuildCommand(serverstorage.AttachCommand(svc, svc), serverStorageCommand, cfgFn())
	commands.BuildCommand(serverstorage.DetachCommand(svc, svc), serverStorageCommand, cfgFn())

	// Server firewall operations
	serverFirewallCommand := commands.BuildCommand(serverfirewall.BaseServerFirewallCommand(), serverCommand, cfgFn())
	commands.BuildCommand(serverfirewall.CreateCommand(svc, svc), serverFirewallCommand, cfgFn())
	commands.BuildCommand(serverfirewall.DeleteCommand(svc, svc), serverFirewallCommand, cfgFn())

	// Storages
	storageCommand := commands.BuildCommand(storage.BaseStorageCommand(), mainCommand, cfgFn())
	commands.BuildCommand(storage.ListCommand(svc), storageCommand, cfgFn())
	commands.BuildCommand(storage.CreateCommand(svc), storageCommand, cfgFn())
	commands.BuildCommand(storage.ModifyCommand(svc), storageCommand, cfgFn())
	commands.BuildCommand(storage.CloneCommand(svc), storageCommand, cfgFn())
	commands.BuildCommand(storage.TemplatizeCommand(svc), storageCommand, cfgFn())
	commands.BuildCommand(storage.DeleteCommand(svc), storageCommand, cfgFn())
	commands.BuildCommand(storage.ImportCommand(svc), storageCommand, cfgFn())
	commands.BuildCommand(storage.ShowCommand(svc, svc), storageCommand, cfgFn())

	backupCommand := commands.BuildCommand(storage.BackupCommand(), storageCommand, cfgFn())
	commands.BuildCommand(storage.CreateBackupCommand(svc), backupCommand, cfgFn())
	commands.BuildCommand(storage.RestoreBackupCommand(svc), backupCommand, cfgFn())

	// IP Addresses
	ipAddressCommand := commands.BuildCommand(ipaddress.BaseIPAddressCommand(), mainCommand, cfgFn())
	commands.BuildCommand(ipaddress.ListCommand(svc), ipAddressCommand, cfgFn())
	commands.BuildCommand(ipaddress.ShowCommand(svc), ipAddressCommand, cfgFn())
	commands.BuildCommand(ipaddress.ModifyCommand(svc), ipAddressCommand, cfgFn())
	commands.BuildCommand(ipaddress.AssignCommand(svc, svc), ipAddressCommand, cfgFn())
	commands.BuildCommand(ipaddress.RemoveCommand(svc), ipAddressCommand, cfgFn())

	// Networks
	networkCommand := commands.BuildCommand(network.BaseNetworkCommand(), mainCommand, cfgFn())
	commands.BuildCommand(network.CreateCommand(svc), networkCommand, cfgFn())
	commands.BuildCommand(network.ListCommand(svc), networkCommand, cfgFn())
	commands.BuildCommand(network.ShowCommand(svc, svc), networkCommand, cfgFn())
	commands.BuildCommand(network.ModifyCommand(svc), networkCommand, cfgFn())
	commands.BuildCommand(network.DeleteCommand(svc), networkCommand, cfgFn())

	// Network Interfaces
	networkInterfaceCommand := commands.BuildCommand(networkinterface.BaseNetworkInterfaceCommand(), serverCommand, cfgFn())
	commands.BuildCommand(networkinterface.CreateCommand(svc, svc), networkInterfaceCommand, cfgFn())
	commands.BuildCommand(networkinterface.ModifyCommand(svc, svc), networkInterfaceCommand, cfgFn())
	commands.BuildCommand(networkinterface.DeleteCommand(svc, svc), networkInterfaceCommand, cfgFn())

	// Routers
	routerCommand := commands.BuildCommand(router.BaseRouterCommand(), mainCommand, cfgFn())
	commands.BuildCommand(router.CreateCommand(svc), routerCommand, cfgFn())
	commands.BuildCommand(router.ListCommand(svc), routerCommand, cfgFn())
	commands.BuildCommand(router.ShowCommand(svc), routerCommand, cfgFn())
	commands.BuildCommand(router.ModifyCommand(svc), routerCommand, cfgFn())
	commands.BuildCommand(router.DeleteCommand(svc), routerCommand, cfgFn())

	// Account
	accountCommand := commands.BuildCommand(account.BaseAccountCommand(), mainCommand, cfgFn())
	commands.BuildCommand(account.ShowCommand(svc), accountCommand, cfgFn())

}
