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

type detachCommand struct {
	*commands.BaseCommand
	serverSvc  service.Server
	storageSvc service.Storage
	params     detachParams
}

type detachParams struct {
	request.DetachStorageRequest
}

func DetachCommand(serverSvc service.Server, storageSvc service.Storage) commands.Command {
	return &detachCommand{
		BaseCommand: commands.New("detach-storage", "Detaches a storage resource from a server"),
		serverSvc:   serverSvc,
		storageSvc:  storageSvc,
	}
}

var DefaultDetachParams = &detachParams{
	DetachStorageRequest: request.DetachStorageRequest{},
}

func (s *detachCommand) InitCommand() {
	s.params = detachParams{DetachStorageRequest: request.DetachStorageRequest{}}

	flagSet := &pflag.FlagSet{}
	flagSet.StringVar(&s.params.Address, "address", DefaultDetachParams.Address, "Detach the storage attached to this address.")

	s.AddFlags(flagSet)
}

func (s *detachCommand) MakeExecuteCommand() func(args []string) (interface{}, error) {
	return func(args []string) (interface{}, error) {
		return Request{
			BuildRequest: func(server *upcloud.Server) interface{} {
				req := s.params.DetachStorageRequest
				req.ServerUUID = server.UUID
				return &req
			},
			Service: s.serverSvc,
			HandleContext: ui.HandleContext{
				MessageFn: func(in interface{}) string {
					req := in.(*request.DetachStorageRequest)
					return fmt.Sprintf("Detaching address %q to server %q", req.Address, req.ServerUUID)
				},
				InteractiveUI: s.Config().InteractiveUI(),
				MaxActions:    maxServerActions,
				Action: func(req interface{}) (interface{}, error) {
					return s.storageSvc.DetachStorage(req.(*request.DetachStorageRequest))
				},
			},
		}.Send(args)
	}
}
