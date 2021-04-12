package serverstorage

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

type attachCommand struct {
	*commands.BaseCommand
	resolver.CachingServer
	completion.Server
	params attachParams
}

type attachParams struct {
	request.AttachStorageRequest
	bootable bool
}

// AttachCommand creates the "server storage attach" command
func AttachCommand() commands.Command {
	return &attachCommand{
		BaseCommand: commands.New("attach", "Attach a storage as a device to a server"),
	}
}

var defaultAttachParams = &attachParams{
	AttachStorageRequest: request.AttachStorageRequest{
		Type:     "disk",
		BootDisk: 0,
		Address:  "virtio",
	},
}

// InitCommand implements Command.InitCommand
func (s *attachCommand) InitCommand() {
	s.params = attachParams{AttachStorageRequest: request.AttachStorageRequest{}}

	flagSet := &pflag.FlagSet{}
	flagSet.StringVar(&s.params.Type, "type", defaultAttachParams.Type, "Type of the attached storage. Available: disk, cdrom")
	flagSet.StringVar(&s.params.Address, "address", defaultAttachParams.Address, "Address where the storage device is attached on the server. \nSpecify only the bus name (ide/scsi/virtio) to auto-select next available address from that bus.")
	flagSet.StringVar(&s.params.StorageUUID, "storage", defaultAttachParams.StorageUUID, "UUID of the storage to attach.")
	flagSet.BoolVar(&s.params.bootable, "boot-disk", false, "Set attached device as the server's boot disk.")

	s.AddFlags(flagSet)
}

// MaximumExecutions implements command.Command
func (s *attachCommand) MaximumExecutions() int {
	return maxServerStorageActions
}

// Execute implements command.Command
func (s *attachCommand) Execute(exec commands.Executor, uuid string) (output.Output, error) {
	storageSvc := exec.Storage()

	if s.params.StorageUUID == "" {
		return nil, fmt.Errorf("storage is required")
	}

	strg, err := storage.SearchSingleStorage(s.params.StorageUUID, storageSvc)
	if err != nil {
		return nil, err
	}

	s.params.StorageUUID = strg.UUID
	s.params.BootDisk = defaultAttachParams.BootDisk

	if s.params.bootable {
		s.params.BootDisk = 1
	}
	req := s.params.AttachStorageRequest
	req.ServerUUID = uuid

	msg := fmt.Sprintf("Attaching storage %q to server %q", req.StorageUUID, req.ServerUUID)
	logline := exec.NewLogEntry(msg)
	logline.StartedNow()

	res, err := storageSvc.AttachStorage(&req)

	if err != nil {
		logline.SetMessage(ui.LiveLogEntryErrorColours.Sprintf("%s: failed (%v)", msg, err.Error()))
		logline.SetDetails(err.Error(), "error: ")
		return nil, err
	}

	logline.SetMessage(fmt.Sprintf("%s: success", msg))
	logline.MarkDone()

	return output.Marshaled{Value: res}, nil
}
