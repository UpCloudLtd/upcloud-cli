package server

import (
	"fmt"

	"github.com/UpCloudLtd/upcloud-cli/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/internal/commands/storage"
	"github.com/UpCloudLtd/upcloud-cli/internal/completion"
	"github.com/UpCloudLtd/upcloud-cli/internal/output"
	"github.com/UpCloudLtd/upcloud-cli/internal/resolver"

	"github.com/UpCloudLtd/upcloud-go-api/v4/upcloud/request"
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

	_ = s.Cobra().MarkFlagRequired("storage")
}

// Execute implements commands.MultipleArgumentCommand
func (s *loadCommand) Execute(exec commands.Executor, uuid string) (output.Output, error) {
	svc := exec.Storage()

	strg, err := storage.SearchSingleStorage(s.params.StorageUUID, svc)
	if err != nil {
		return nil, err
	}
	s.params.StorageUUID = strg.UUID

	req := s.params.LoadCDROMRequest
	req.ServerUUID = uuid

	msg := fmt.Sprintf("Loading %q as a CD-ROM of server %q", req.StorageUUID, req.ServerUUID)
	logline := exec.NewLogEntry(msg)

	logline.StartedNow()
	logline.SetMessage(fmt.Sprintf("%s: sending request", msg))

	res, err := svc.LoadCDROM(&req)
	if err != nil {
		return commands.HandleError(logline, fmt.Sprintf("%s: failed", msg), err)
	}

	logline.SetMessage(fmt.Sprintf("%s: request sent", msg))
	logline.MarkDone()

	return output.OnlyMarshaled{Value: res}, nil
}
