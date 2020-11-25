package server

import (
	"fmt"
	"github.com/UpCloudLtd/cli/internal/commands"
	"github.com/UpCloudLtd/cli/internal/ui"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/request"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/service"
	"github.com/spf13/pflag"
)

type loadCommand struct {
	*commands.BaseCommand
	serverSvc  service.Server
	storageSvc service.Storage
	params     loadParams
}

type loadParams struct {
	request.LoadCDROMRequest
}

func LoadCommand(serverSvc service.Server, storageSvc service.Storage) commands.Command {
	return &loadCommand{
		BaseCommand: commands.New("load", "Load a CD-ROM"),
		serverSvc:   serverSvc,
		storageSvc:  storageSvc,
	}
}

var DefaultLoadParams = &loadParams{
	LoadCDROMRequest: request.LoadCDROMRequest{},
}

func (s *loadCommand) InitCommand() {
	s.SetPositionalArgHelp(PositionalArgHelp)
	s.ArgCompletion(GetArgCompFn(s.serverSvc))
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
			Service:    s.serverSvc,
			ExactlyOne: true,
			HandleContext: ui.HandleContext{
				MessageFn: func(in interface{}) string {
					req := in.(*request.LoadCDROMRequest)
					return fmt.Sprintf("Loading %q as a CD-ROM of server %q", req.StorageUUID, req.ServerUUID)
				},
				MaxActions: maxServerActions,
				Action: func(req interface{}) (interface{}, error) {
					return s.storageSvc.LoadCDROM(req.(*request.LoadCDROMRequest))
				},
			},
		}.Send(args)
	}
}
