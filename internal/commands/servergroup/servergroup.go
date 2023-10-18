package servergroup

import (
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/namedargs"
	"github.com/UpCloudLtd/upcloud-go-api/v6/upcloud"
)

const (
	met           = "met"
	notApplicable = "-"
	unMet         = "unmet"
)

// BaseServergroupCommand creates the base "servergroup" command
func BaseServergroupCommand() commands.Command {
	return &servergroupCommand{
		commands.New("servergroup", "Manage server groups"),
	}
}

type servergroupCommand struct {
	*commands.BaseCommand
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
