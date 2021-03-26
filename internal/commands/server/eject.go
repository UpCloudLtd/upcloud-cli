package server

import (
	"github.com/UpCloudLtd/cli/internal/commands"
	"github.com/UpCloudLtd/cli/internal/ui"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/request"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/service"
)

type ejectCommand struct {
	*commands.BaseCommand
	serverSvc  service.Server
	storageSvc service.Storage
	params     ejectParams
}

type ejectParams struct {
	request.EjectCDROMRequest
}

func (s *ejectCommand) InitCommand() {
	s.SetPositionalArgHelp(PositionalArgHelp)
	s.ArgCompletion(GetServerArgumentCompletionFunction(s.serverSvc))
}

// EjectCommand creates the "server eject" command
func EjectCommand(serverSvc service.Server, storageSvc service.Storage) commands.Command {
	return &ejectCommand{
		BaseCommand: commands.New("eject", "Eject a CD-ROM from the server"),
		serverSvc:   serverSvc,
		storageSvc:  storageSvc,
	}
}

// MakeExecuteCommand implements Command.MakeExecuteCommand
func (s *ejectCommand) MakeExecuteCommand() func(args []string) (interface{}, error) {
	return func(args []string) (interface{}, error) {
		return Request{
			BuildRequest: func(uuid string) interface{} {
				req := s.params.EjectCDROMRequest
				req.ServerUUID = uuid
				return &req
			},
			Service: s.serverSvc,
			Handler: ui.HandleContext{
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
