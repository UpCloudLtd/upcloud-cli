package server

import (
	"github.com/UpCloudLtd/cli/internal/commands"
	"github.com/UpCloudLtd/cli/internal/ui"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/request"
)

type ejectCommand struct {
	*commands.BaseCommand
	params ejectParams
}

type ejectParams struct {
	request.EjectCDROMRequest
}

func (s *ejectCommand) InitCommand() {
	s.SetPositionalArgHelp(PositionalArgHelp)
	s.ArgCompletion(GetServerArgumentCompletionFunction(s.Config()))
}

// EjectCommand creates the "server eject" command
func EjectCommand() commands.Command {
	return &ejectCommand{
		BaseCommand: commands.New("eject", "Eject a CD-ROM from the server"),
	}
}

// MakeExecuteCommand implements Command.MakeExecuteCommand
func (s *ejectCommand) MakeExecuteCommand() func(args []string) (interface{}, error) {
	return func(args []string) (interface{}, error) {
		serverSvc := s.Config().Service.Server()
		storageSvc := s.Config().Service.Storage()

		return Request{
			BuildRequest: func(uuid string) interface{} {
				req := s.params.EjectCDROMRequest
				req.ServerUUID = uuid
				return &req
			},
			Service: serverSvc,
			Handler: ui.HandleContext{
				RequestID:  func(in interface{}) string { return in.(*request.EjectCDROMRequest).ServerUUID },
				MaxActions: maxServerActions,
				ActionMsg:  "Ejecting CD-ROM of server",
				Action: func(req interface{}) (interface{}, error) {
					return storageSvc.EjectCDROM(req.(*request.EjectCDROMRequest))
				},
			},
		}.Send(args)
	}
}
