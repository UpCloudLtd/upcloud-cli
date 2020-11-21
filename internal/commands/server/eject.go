package server

import (
	"github.com/UpCloudLtd/cli/internal/commands"
	"github.com/UpCloudLtd/cli/internal/interfaces"
	"github.com/UpCloudLtd/cli/internal/ui"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/request"
	"github.com/spf13/pflag"
)

type ejectCommand struct {
	*commands.BaseCommand
	service interfaces.ServerAndStorage
	params  ejectParams
	flagSet *pflag.FlagSet
}

type ejectParams struct {
	request.EjectCDROMRequest
}

func EjectCommand(service interfaces.ServerAndStorage) commands.Command {
	return &ejectCommand{
		BaseCommand: commands.New("eject", "Eject a CD-ROM"),
		service:     service,
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
			Service: s.service,
			HandleContext: ui.HandleContext{
				RequestId:  func(in interface{}) string { return in.(*request.EjectCDROMRequest).ServerUUID },
				MaxActions: maxServerActions,
				ActionMsg:  "Ejecting CD-ROM of server",
				Action: func(req interface{}) (interface{}, error) {
					return s.service.EjectCDROM(req.(*request.EjectCDROMRequest))
				},
			},
		}.Send(args)
	}
}
