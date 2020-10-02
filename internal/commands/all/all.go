package all

import (
	"github.com/spf13/viper"

	"github.com/UpCloudLtd/cli/internal/commands"
	"github.com/UpCloudLtd/cli/internal/commands/plan"
	"github.com/UpCloudLtd/cli/internal/commands/server"
)

func BuildCommands(mainCommand commands.Command, mainConfig *viper.Viper) {
	// Plans
	planCommand := commands.BuildCommand(plan.PlanCommand(), mainCommand, mainConfig)
	commands.BuildCommand(plan.ListCommand(), planCommand, mainConfig)

	// Servers
	serverCommand := commands.BuildCommand(server.ServerCommand(), mainCommand, mainConfig)
	commands.BuildCommand(server.ListCommand(), serverCommand, mainConfig)
	commands.BuildCommand(server.ShowCommand(), serverCommand, mainConfig)
}
