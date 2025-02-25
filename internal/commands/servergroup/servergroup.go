package servergroup

import (
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/namedargs"
	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud"
)

const (
	met           = "met"
	notApplicable = "-"
	unMet         = "unmet"
)

// BaseServergroupCommand creates the base "servergroup" command
func BaseServergroupCommand() commands.Command {
	baseCmd := commands.New("server-group", "Manage server groups")
	baseCmd.SetDeprecatedAliases([]string{"servergroup"})

	return &servergroupCommand{
		BaseCommand: baseCmd,
	}
}

type servergroupCommand struct {
	*commands.BaseCommand
}

// InitCommand implements Command.InitCommand
func (sg *servergroupCommand) InitCommand() {
	sg.Cobra().Aliases = []string{"sg", "servergroup"}
	// Deprecating servergroup in favour of server-group
	// TODO: Remove this in the future
	commands.SetDeprecationHelp(sg.Cobra(), sg.DeprecatedAliases())
}

func stringsToServerUUIDSlice(exec commands.Executor, servers []string) (upcloud.ServerUUIDSlice, error) {
	slice := make(upcloud.ServerUUIDSlice, 0)
	for _, v := range servers {
		if v != "" {
			serverUUID, err := namedargs.ResolveServer(exec, v)
			if err != nil {
				return nil, err
			}
			slice = append(slice, serverUUID)
		}
	}

	return slice, nil
}
