package server

import (
	"fmt"

	"github.com/UpCloudLtd/cli/internal/commands"
	"github.com/UpCloudLtd/cli/internal/commands/storage"
	"github.com/UpCloudLtd/cli/internal/completion"
	"github.com/UpCloudLtd/cli/internal/output"
	"github.com/UpCloudLtd/cli/internal/resolver"
	"github.com/UpCloudLtd/cli/internal/ui"

	"github.com/UpCloudLtd/upcloud-go-api/upcloud/request"
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
		BaseCommand: commands.New("load", "Load a CD-ROM into the server"),
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
}

func (s *loadCommand) Execute(exec commands.Executor, uuid string) (output.Output, error) {
	svc := exec.Storage()

	if s.params.StorageUUID == "" {
		return nil, fmt.Errorf("storage is required")
	}

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
		logline.SetMessage(ui.LiveLogEntryErrorColours.Sprintf("%s: failed (%v)", msg, err.Error()))
		logline.SetDetails(err.Error(), "error: ")
		return nil, err
	}

	logline.SetMessage(fmt.Sprintf("%s: request sent", msg))
	logline.MarkDone()

	return output.Marshaled{Value: res}, nil
}
