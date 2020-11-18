package storage

import (
	"fmt"
	"github.com/UpCloudLtd/cli/internal/commands"
	"github.com/UpCloudLtd/cli/internal/ui"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/request"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/service"
	"github.com/spf13/pflag"
	"sync"
)

type attachCommand struct {
	*commands.BaseCommand
	service service.Storage
	params  attachParams
	flagSet *pflag.FlagSet
}

type attachParams struct {
	request.AttachStorageRequest
}

func newAttachParams() attachParams {
	return attachParams{AttachStorageRequest: request.AttachStorageRequest{}}
}

func AttachCommand(service service.Storage) commands.Command {
	return &attachCommand{
		BaseCommand: commands.New("attach", "Attach a storage"),
		service:     service,
	}
}

var DefaultAttachParams = &attachParams{
	AttachStorageRequest: request.AttachStorageRequest{
		Type:     "disk",
		BootDisk: 0,
	},
}

func (s *attachCommand) InitCommand() {
	s.params = newAttachParams()

	s.flagSet = &pflag.FlagSet{}
	s.flagSet.StringVar(&s.params.Type, "type", DefaultAttachParams.Type, "The type of the attached storage.\nAvailable: disk, cdrom")
	s.flagSet.StringVar(&s.params.Address, "address", DefaultAttachParams.Address, "The address where the storage device is attached on the server. \nSpecify only the bus name (ide/scsi/virtio) to auto-select next available address from that bus.")
	s.flagSet.StringVar(&s.params.StorageUUID, "storage", DefaultAttachParams.StorageUUID, "The UUID of the storage to attach.")
	s.flagSet.IntVar(&s.params.BootDisk, "boot-disk", DefaultAttachParams.BootDisk, "If the value is 1 the storage device will be used as a boot disk, unless overridden with the server boot_order attribute.")

	s.AddFlags(s.flagSet)
}

func (s *attachCommand) MakeExecuteCommand() func(args []string) (interface{}, error) {
	return func(args []string) (interface{}, error) {

		if len(args) != 1 {
			return nil, fmt.Errorf("server uuid is required")
		}

		s.params.ServerUUID = args[0]
		var attachStorageRequests []request.AttachStorageRequest
		attachStorageRequests = append(attachStorageRequests, s.params.AttachStorageRequest)

		var (
			mu            sync.Mutex
			numOk         int
			serverDetails *upcloud.ServerDetails
		)
		handler := func(idx int, e *ui.LogEntry) {
			req := attachStorageRequests[idx]
			msg := fmt.Sprintf("Attaching storage %q to server %q", req.StorageUUID, req.ServerUUID)
			e.SetMessage(msg)
			e.Start()
			details, err := s.service.AttachStorage(&req)
			if err != nil {
				e.SetMessage(ui.LiveLogEntryErrorColours.Sprintf("%s: failed", msg))
				e.SetDetails(err.Error(), "error: ")
			} else {
				e.SetMessage(fmt.Sprintf("%s: done", msg))
				e.SetDetails(details.UUID, "UUID: ")
				mu.Lock()
				numOk++
				serverDetails = details
				mu.Unlock()
			}
		}
		ui.StartWorkQueue(ui.WorkQueueConfig{
			NumTasks:           len(attachStorageRequests),
			MaxConcurrentTasks: 5,
			EnableUI:           s.Config().InteractiveUI(),
		}, handler)
		if numOk != len(attachStorageRequests) {
			return nil, fmt.Errorf("number of operations that failed: %d", len(attachStorageRequests)-numOk)
		}

		return serverDetails, nil
	}
}
