package storage

import (
	"fmt"
	"github.com/UpCloudLtd/cli/internal/commands"
	"github.com/UpCloudLtd/cli/internal/ui"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/request"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/service"
	"github.com/spf13/pflag"
	"sync"
)

type restoreBackupCommand struct {
	*commands.BaseCommand
	service service.Storage
	params  restoreBackupParams
	flagSet *pflag.FlagSet
}

type restoreBackupParams struct {
	request.RestoreBackupRequest
}

func RestoreBackupCommand(service service.Storage) commands.Command {
	return &restoreBackupCommand{
		BaseCommand: commands.New("restore-backup", "Restore backup of a storage"),
		service:     service,
	}
}

func (s *restoreBackupCommand) MakeExecuteCommand() func(args []string) (interface{}, error) {
	return func(args []string) (interface{}, error) {

		if len(args) < 1 {
			return nil, fmt.Errorf("storage uuid is required")
		}

		var restoreBackupStorageRequests []request.RestoreBackupRequest
		for _, v := range args {
			s.params.RestoreBackupRequest.UUID = v
			restoreBackupStorageRequests = append(restoreBackupStorageRequests, s.params.RestoreBackupRequest)
		}

		var (
			mu    sync.Mutex
			numOk int
		)
		handler := func(idx int, e *ui.LogEntry) {
			req := restoreBackupStorageRequests[idx]
			msg := fmt.Sprintf("Creating backup of storage %q", req.UUID)
			e.SetMessage(msg)
			e.Start()
			err := s.service.RestoreBackup(&req)
			if err != nil {
				e.SetMessage(ui.LiveLogEntryErrorColours.Sprintf("%s: failed", msg))
				e.SetDetails(err.Error(), "error: ")
			} else {
				e.SetMessage(fmt.Sprintf("%s: done", msg))
				e.SetDetails(req.UUID, "UUID: ")
				mu.Lock()
				numOk++
				mu.Unlock()
			}
		}
		ui.StartWorkQueue(ui.WorkQueueConfig{
			NumTasks:           len(restoreBackupStorageRequests),
			MaxConcurrentTasks: 5,
			EnableUI:           s.Config().InteractiveUI(),
		}, handler)
		if numOk != len(restoreBackupStorageRequests) {
			return nil, fmt.Errorf("number of backup creations that failed: %d", len(restoreBackupStorageRequests)-numOk)
		}

		return nil, nil
	}
}
