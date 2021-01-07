package server

import "github.com/UpCloudLtd/cli/internal/commands"

func PlanCommand() commands.Command {
	return &planCommand{
		BaseCommand: commands.New("plan", "Server plans"),
	}
}

type planCommand struct {
	*commands.BaseCommand
}
