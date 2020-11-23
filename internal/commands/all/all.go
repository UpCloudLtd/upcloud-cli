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
	commands.BuildCommand(server.ListCommand(), serverCommand, cfgFn())
	commands.BuildCommand(server.ShowCommand(), serverCommand, cfgFn())
	commands.BuildCommand(server.StartCommand(), serverCommand, cfgFn())
	commands.BuildCommand(server.StopCommand(), serverCommand, cfgFn())
	commands.BuildCommand(server.CreateCommand(), serverCommand, cfgFn())
	commands.BuildCommand(server.DeleteCommand(), serverCommand, cfgFn())

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
		storage.ShowCommand(svc),
	}
	for _, stgCmd := range stgCmds {
		commands.BuildCommand(stgCmd, storageCommand, cfgFn())
	}
	backupCommand := commands.BuildCommand(storage.BackupCommand(), storageCommand, cfgFn())
	commands.BuildCommand(storage.CreateBackupCommand(svc), backupCommand, cfgFn())
	commands.BuildCommand(storage.RestoreBackupCommand(svc), backupCommand, cfgFn())
}
