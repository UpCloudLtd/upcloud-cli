package all

import (
	"github.com/UpCloudLtd/cli/internal/commands"
	"github.com/UpCloudLtd/cli/internal/commands/account"
	ip_address "github.com/UpCloudLtd/cli/internal/commands/ip_address"
	"github.com/UpCloudLtd/cli/internal/commands/network"
	"github.com/UpCloudLtd/cli/internal/commands/network_interface"
	"github.com/UpCloudLtd/cli/internal/commands/plan"
	"github.com/UpCloudLtd/cli/internal/commands/router"
	"github.com/UpCloudLtd/cli/internal/commands/server"
	"github.com/UpCloudLtd/cli/internal/commands/storage"
	"github.com/UpCloudLtd/cli/internal/config"
	"github.com/UpCloudLtd/cli/internal/upapi"
)

func BuildCommands(mainCommand commands.Command, mainConfig *config.Config) {
	cfgFn := func() *config.Config { return config.New(mainConfig.Viper()) }
	svc := upapi.Service(cfgFn())

	// Plans
	planCommand := commands.BuildCommand(plan.PlanCommand(), mainCommand, cfgFn())
	commands.BuildCommand(plan.ListCommand(), planCommand, cfgFn())

	// Servers
	serverCommand := commands.BuildCommand(server.ServerCommand(), mainCommand, cfgFn())
	commands.BuildCommand(server.ListCommand(svc), serverCommand, cfgFn())
	commands.BuildCommand(server.ConfigurationsCommand(svc), serverCommand, cfgFn())
	commands.BuildCommand(server.ShowCommand(svc), serverCommand, cfgFn())
	commands.BuildCommand(server.StartCommand(svc), serverCommand, cfgFn())
	commands.BuildCommand(server.RestartCommand(svc), serverCommand, cfgFn())
	commands.BuildCommand(server.StopCommand(svc), serverCommand, cfgFn())
	commands.BuildCommand(server.CreateCommand(svc, svc), serverCommand, cfgFn())
	commands.BuildCommand(server.ModifyCommand(svc), serverCommand, cfgFn())
	commands.BuildCommand(server.AttachCommand(svc, svc), serverCommand, cfgFn())
	commands.BuildCommand(server.LoadCommand(svc, svc), serverCommand, cfgFn())
	commands.BuildCommand(server.DetachCommand(svc, svc), serverCommand, cfgFn())
	commands.BuildCommand(server.EjectCommand(svc, svc), serverCommand, cfgFn())
	commands.BuildCommand(server.DeleteCommand(svc), serverCommand, cfgFn())

	// Network Interfaces
	networkInterfaceCommand := commands.BuildCommand(network_interface.NetworkInterfaceCommand(), serverCommand, cfgFn())
	commands.BuildCommand(network_interface.CreateCommand(svc), networkInterfaceCommand, cfgFn())
	commands.BuildCommand(network_interface.ModifyCommand(svc), networkInterfaceCommand, cfgFn())
	commands.BuildCommand(network_interface.DeleteCommand(svc), networkInterfaceCommand, cfgFn())

	// Storages
	storageCommand := commands.BuildCommand(storage.StorageCommand(), mainCommand, cfgFn())
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
	ipAddressCommand := commands.BuildCommand(ip_address.IpAddressCommand(), mainCommand, cfgFn())
	commands.BuildCommand(ip_address.ListCommand(svc), ipAddressCommand, cfgFn())
	commands.BuildCommand(ip_address.ShowCommand(svc), ipAddressCommand, cfgFn())
	commands.BuildCommand(ip_address.ModifyCommand(svc), ipAddressCommand, cfgFn())
	commands.BuildCommand(ip_address.AssignCommand(svc), ipAddressCommand, cfgFn())
	commands.BuildCommand(ip_address.ReleaseCommand(svc), ipAddressCommand, cfgFn())

	// Networks
	networkCommand := commands.BuildCommand(network.NetworkCommand(), mainCommand, cfgFn())
	commands.BuildCommand(network.CreateCommand(svc), networkCommand, cfgFn())
	commands.BuildCommand(network.ListCommand(svc), networkCommand, cfgFn())
	commands.BuildCommand(network.ShowCommand(svc, svc), networkCommand, cfgFn())
	commands.BuildCommand(network.ModifyCommand(svc), networkCommand, cfgFn())
	commands.BuildCommand(network.DeleteCommand(svc), networkCommand, cfgFn())

	// Routers
	routerCommand := commands.BuildCommand(router.RouterCommand(), mainCommand, cfgFn())
	commands.BuildCommand(router.CreateCommand(svc), routerCommand, cfgFn())
	commands.BuildCommand(router.ListCommand(svc), routerCommand, cfgFn())
	commands.BuildCommand(router.ShowCommand(svc), routerCommand, cfgFn())
	commands.BuildCommand(router.ModifyCommand(svc), routerCommand, cfgFn())
	commands.BuildCommand(router.DeleteCommand(svc), routerCommand, cfgFn())

	// Account
	accountCommand := commands.BuildCommand(account.AccountCommand(), mainCommand, cfgFn())
	commands.BuildCommand(account.ShowCommand(svc), accountCommand, cfgFn())

}
