package all

import (
	"github.com/UpCloudLtd/cli/internal/commands"
	ip_address "github.com/UpCloudLtd/cli/internal/commands/ip-address"
	"github.com/UpCloudLtd/cli/internal/commands/network"
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
	svrCmds := []commands.Command{
		server.ListCommand(svc),
		server.ConfigurationsCommand(svc),
		server.ShowCommand(svc),
		server.StartCommand(svc),
		server.RestartCommand(svc),
		server.StopCommand(svc),
		server.CreateCommand(svc, svc),
		server.ModifyCommand(svc),
		server.AttachCommand(svc, svc),
		server.LoadCommand(svc, svc),
		server.DetachCommand(svc, svc),
		server.EjectCommand(svc, svc),
		server.DeleteCommand(svc),
	}
	for _, svrCmds := range svrCmds {
		commands.BuildCommand(svrCmds, serverCommand, cfgFn())
	}

	// Storages
	storageCommand := commands.BuildCommand(storage.StorageCommand(), mainCommand, cfgFn())
	stgCmds := []commands.Command{
		storage.ListCommand(svc),
		storage.CreateCommand(svc),
		storage.ModifyCommand(svc),
		storage.CloneCommand(svc),
		storage.TemplatizeCommand(svc),
		storage.DeleteCommand(svc),
		storage.ImportCommand(svc),
		storage.ShowCommand(svc, svc),
	}
	for _, stgCmd := range stgCmds {
		commands.BuildCommand(stgCmd, storageCommand, cfgFn())
	}
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
	commands.BuildCommand(network.ShowCommand(svc), networkCommand, cfgFn())
	commands.BuildCommand(network.ModifyCommand(svc), networkCommand, cfgFn())
	commands.BuildCommand(network.DeleteCommand(svc), networkCommand, cfgFn())

	// Routers
	routerCommand := commands.BuildCommand(router.RouterCommand(), mainCommand, cfgFn())
	commands.BuildCommand(router.CreateCommand(svc), routerCommand, cfgFn())
	commands.BuildCommand(router.ListCommand(svc), routerCommand, cfgFn())
	commands.BuildCommand(router.ShowCommand(svc), routerCommand, cfgFn())
	commands.BuildCommand(router.ModifyCommand(svc), routerCommand, cfgFn())
	commands.BuildCommand(router.DeleteCommand(svc), routerCommand, cfgFn())
}
