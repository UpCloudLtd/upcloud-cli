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

type attachCommand struct {
	*commands.BaseCommand
	serverSvc service.Server
	storageSvc service.Storage
	params  attachParams
}

type attachParams struct {
	request.AttachStorageRequest
}

func AttachCommand(serverSvc service.Server, storageSvc service.Storage) commands.Command {
	return &attachCommand{
		BaseCommand: commands.New("attach-storage", "Attaches a storage as a device to a server"),
		serverSvc: serverSvc,
		storageSvc: storageSvc,
	}
}

var DefaultAttachParams = &attachParams{
	AttachStorageRequest: request.AttachStorageRequest{
		Type:     "disk",
		BootDisk: 0,
		Address:  "virtio",
	},
}

func (s *attachCommand) InitCommand() {
	s.params = attachParams{AttachStorageRequest: request.AttachStorageRequest{}}

	flagSet := &pflag.FlagSet{}
	flagSet.StringVar(&s.params.Type, "type", DefaultAttachParams.Type, "The type of the attached storage.\nAvailable: disk, cdrom")
	flagSet.StringVar(&s.params.Address, "address", DefaultAttachParams.Address, "The address where the storage device is attached on the server. \nSpecify only the bus name (ide/scsi/virtio) to auto-select next available address from that bus.")
	flagSet.StringVar(&s.params.StorageUUID, "storage", DefaultAttachParams.StorageUUID, "The UUID of the storage to attach.")
	flagSet.IntVar(&s.params.BootDisk, "boot-disk", DefaultAttachParams.BootDisk, "If the value is 1 the storage device will be used as a boot disk, unless overridden with the server boot_order attribute.")

	s.AddFlags(flagSet)
}

func (s *attachCommand) MakeExecuteCommand() func(args []string) (interface{}, error) {
	return func(args []string) (interface{}, error) {
		return Request{
			BuildRequest: func(server *upcloud.Server) interface{} {
				req := s.params.AttachStorageRequest
				req.ServerUUID = server.UUID
				return &req
			},
			Service:    s.serverSvc,
			ExactlyOne: true,
			HandleContext: ui.HandleContext{
				InteractiveUI: s.Config().InteractiveUI(),
				MaxActions:    maxServerActions,
				MessageFn: func(in interface{}) string {
					req := in.(*request.AttachStorageRequest)
					return fmt.Sprintf("Attaching storage %q to server %q", req.StorageUUID, req.ServerUUID)
				},
				Action: func(req interface{}) (interface{}, error) {
					return s.storageSvc.AttachStorage(req.(*request.AttachStorageRequest))
				},
			},
		}.Send(args)
	}
}
