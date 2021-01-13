package server_storage

import (
	"fmt"
	"github.com/UpCloudLtd/cli/internal/commands"
	"github.com/UpCloudLtd/cli/internal/commands/server"
	"github.com/UpCloudLtd/cli/internal/ui"
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
		BaseCommand: commands.New("detach", "Detaches a storage resource from a server"),
		serverSvc:   serverSvc,
		storageSvc:  storageSvc,
	}
}

var DefaultDetachParams = &detachParams{
	DetachStorageRequest: request.DetachStorageRequest{},
}

func (s *detachCommand) InitCommand() {
	s.SetPositionalArgHelp(PositionalArgHelp)
	s.ArgCompletion(server.GetArgCompFn(s.serverSvc))
	s.params = detachParams{DetachStorageRequest: request.DetachStorageRequest{}}

	flagSet := &pflag.FlagSet{}
	flagSet.StringVar(&s.params.Address, "address", DefaultDetachParams.Address, "Detach the storage attached to this address.\n[Required]")

	s.AddFlags(flagSet)
}

func (s *detachCommand) MakeExecuteCommand() func(args []string) (interface{}, error) {
	return func(args []string) (interface{}, error) {

		if s.params.Address == "" {
			return nil, fmt.Errorf("address is required")
		}

		return server.Request{
			BuildRequest: func(uuid string) interface{} {
				req := s.params.DetachStorageRequest
				req.ServerUUID = uuid
				return &req
			},
			Service: s.serverSvc,
			Handler: ui.HandleContext{
				MessageFn: func(in interface{}) string {
					req := in.(*request.DetachStorageRequest)
					return fmt.Sprintf("Detaching address %q from server %q", req.Address, req.ServerUUID)
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
