package server

import "github.com/UpCloudLtd/cli/internal/commands"

func ServerCommand() commands.Command {
	return &planCommand{
		Command: commands.New("server", "List, show & control servers"),
	}
}

type planCommand struct {
	commands.Command
}
