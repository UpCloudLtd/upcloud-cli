package server

import (
	"fmt"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands/storage"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/completion"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/config"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/namedargs"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/output"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/resolver"

	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud/request"
	"github.com/spf13/pflag"
)

type loadCommand struct {
	*commands.BaseCommand
	resolver.CachingServer
	completion.Server
	params loadParams
}

type loadParams struct {
	request.LoadCDROMRequest
}

// LoadCommand creates the "server load" command
func LoadCommand() commands.Command {
	return &loadCommand{
		BaseCommand: commands.New(
			"load",
			"Load a CD-ROM into the server",
			"upctl server load my_server4 --storage 01000000-0000-4000-8000-000080030101",
		),
	}
}

var defaultLoadParams = &loadParams{
	LoadCDROMRequest: request.LoadCDROMRequest{},
}

// InitCommand implements Command.InitCommand
func (s *loadCommand) InitCommand() {
	s.params = loadParams{LoadCDROMRequest: request.LoadCDROMRequest{}}

	flagSet := &pflag.FlagSet{}
	flagSet.StringVar(&s.params.StorageUUID, "storage", defaultLoadParams.StorageUUID, "The UUID of the storage to be loaded in the CD-ROM device.")

	s.AddFlags(flagSet)

	commands.Must(s.Cobra().MarkFlagRequired("storage"))
}

func (s *loadCommand) InitCommandWithConfig(cfg *config.Config) {
	commands.Must(s.Cobra().RegisterFlagCompletionFunc("storage", namedargs.CompletionFunc(completion.StorageCDROMUUID{}, cfg)))
}

// Execute implements commands.MultipleArgumentCommand
func (s *loadCommand) Execute(exec commands.Executor, uuid string) (output.Output, error) {
	svc := exec.Storage()

	strg, err := storage.SearchSingleStorage(s.params.StorageUUID, exec)
	if err != nil {
		return nil, err
	}
	s.params.StorageUUID = strg.UUID

	req := s.params.LoadCDROMRequest
	req.ServerUUID = uuid

	msg := fmt.Sprintf("Loading %q as a CD-ROM of server %q", req.StorageUUID, req.ServerUUID)
	exec.PushProgressStarted(msg)

	res, err := svc.LoadCDROM(exec.Context(), &req)
	if err != nil {
		return commands.HandleError(exec, msg, err)
	}

	exec.PushProgressSuccess(msg)

	return output.OnlyMarshaled{Value: res}, nil
}
