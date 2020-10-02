package plan

import "github.com/UpCloudLtd/cli/internal/commands"

func PlanCommand() commands.Command {
	return &planCommand{
		Command: commands.New("plan", "Server plans"),
	}
}

type planCommand struct {
	commands.Command
}
