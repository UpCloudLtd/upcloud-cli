package server

import (
	"fmt"
	"github.com/UpCloudLtd/cli/internal/commands"
	"github.com/UpCloudLtd/cli/internal/interfaces"
	"github.com/UpCloudLtd/cli/internal/ui"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/request"
	"github.com/spf13/pflag"
)

type loadCommand struct {
	*commands.BaseCommand
	service interfaces.ServerAndStorage
	params  loadParams
}

type loadParams struct {
	request.LoadCDROMRequest
}

func LoadCommand(service interfaces.ServerAndStorage) commands.Command {
	return &loadCommand{
		BaseCommand: commands.New("load", "Load a CD-ROM"),
		service:     service,
	}
}

var DefaultLoadParams = &loadParams{
	LoadCDROMRequest: request.LoadCDROMRequest{},
}

func (s *loadCommand) InitCommand() {
	s.params = loadParams{LoadCDROMRequest: request.LoadCDROMRequest{}}

	flagSet := &pflag.FlagSet{}
	flagSet.StringVar(&s.params.StorageUUID, "storage", DefaultLoadParams.StorageUUID, "The UUID of the storage to be loaded in the CD-ROM device.")

	s.AddFlags(flagSet)
}

func (s *loadCommand) MakeExecuteCommand() func(args []string) (interface{}, error) {
	return func(args []string) (interface{}, error) {
		return Request{
			BuildRequest: func(server *upcloud.Server) interface{} {
				req := s.params.LoadCDROMRequest
				req.ServerUUID = server.UUID
				return &req
			},
			Service: s.service,
			ExactlyOne: true,
			HandleContext: ui.HandleContext{
				MessageFn: func(in interface{}) string {
					req := in.(*request.LoadCDROMRequest)
					return fmt.Sprintf("Loading %q as a CD-ROM of server %q", req.StorageUUID, req.ServerUUID)
				},
				MaxActions: maxServerActions,
				Action: func(req interface{}) (interface{}, error) {
					return s.service.LoadCDROM(req.(*request.LoadCDROMRequest))
				},
			},
		}.Send(args)
	}
}
