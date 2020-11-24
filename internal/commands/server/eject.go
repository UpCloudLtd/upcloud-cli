package server

import (
	"github.com/UpCloudLtd/cli/internal/commands"
	"github.com/UpCloudLtd/cli/internal/ui"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/request"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/service"
	"github.com/spf13/pflag"
)

type ejectCommand struct {
	*commands.BaseCommand
	serverSvc service.Server
	storageSvc service.Storage
	params  ejectParams
	flagSet *pflag.FlagSet
}

type ejectParams struct {
	request.EjectCDROMRequest
}

func EjectCommand(serverSvc service.Server, storageSvc service.Storage) commands.Command {
	return &ejectCommand{
		BaseCommand: commands.New("eject", "Eject a CD-ROM"),
		serverSvc: serverSvc,
		storageSvc: storageSvc,
	}
}

func (s *ejectCommand) MakeExecuteCommand() func(args []string) (interface{}, error) {
	return func(args []string) (interface{}, error) {
		return Request{
			BuildRequest: func(server *upcloud.Server) interface{} {
				req := s.params.EjectCDROMRequest
				req.ServerUUID = server.UUID
				return &req
			},
			Service: s.serverSvc,
			HandleContext: ui.HandleContext{
				RequestID:  func(in interface{}) string { return in.(*request.EjectCDROMRequest).ServerUUID },
				MaxActions: maxServerActions,
				ActionMsg:  "Ejecting CD-ROM of server",
				Action: func(req interface{}) (interface{}, error) {
					return s.storageSvc.EjectCDROM(req.(*request.EjectCDROMRequest))
				},
			},
		}.Send(args)
	}
}
