package all

import (
	"github.com/UpCloudLtd/cli/internal/commands"
	"github.com/UpCloudLtd/cli/internal/commands/plan"
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
		server.CreateCommand(svc),
		server.ModifyCommand(svc),
		server.AttachCommand(svc),
		server.LoadCommand(svc),
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
		server.DetachCommand(svc),
		server.EjectCommand(svc),
		storage.CloneCommand(svc),
		storage.TemplatizeCommand(svc),
		storage.CreateBackupCommand(svc),
		storage.RestoreBackupCommand(svc),
		storage.DeleteCommand(svc),
		storage.ImportCommand(svc),
		storage.ShowCommand(svc),
	}
	for _, stgCmd := range stgCmds {
		commands.BuildCommand(stgCmd, storageCommand, cfgFn())
	}
}
