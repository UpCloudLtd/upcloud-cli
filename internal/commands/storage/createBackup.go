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

type createBackupCommand struct {
	*commands.BaseCommand
	service service.Storage
	params  createBackupParams
	flagSet *pflag.FlagSet
}

type createBackupParams struct {
	request.CreateBackupRequest
}

func newCreateBackupParams() createBackupParams {
	return createBackupParams{CreateBackupRequest: request.CreateBackupRequest{}}
}

func CreateBackupCommand(service service.Storage) commands.Command {
	return &createBackupCommand{
		BaseCommand: commands.New("create-backup", "Create backup of a storage"),
		service:     service,
	}
}

var DefaultCreateBackupParams = &createBackupParams{
	CreateBackupRequest: request.CreateBackupRequest{},
}

func (s *createBackupCommand) InitCommand() {
	s.params = newCreateBackupParams()

	s.flagSet = &pflag.FlagSet{}
	s.flagSet.StringVar(&s.params.Title, "title", DefaultCreateBackupParams.Title, "A short, informational description.")

	s.AddFlags(s.flagSet)
}

func (s *createBackupCommand) MakeExecuteCommand() func(args []string) (interface{}, error) {
	return func(args []string) (interface{}, error) {

		if len(args) < 1 {
			return nil, fmt.Errorf("storage uuid is required")
		}

		var createBackupStorageRequests []request.CreateBackupRequest
		for _, v := range args {
			s.params.CreateBackupRequest.UUID = v
			createBackupStorageRequests = append(createBackupStorageRequests, s.params.CreateBackupRequest)
		}

		var (
			mu             sync.Mutex
			numOk          int
			storageDetails []*upcloud.StorageDetails
		)
		handler := func(idx int, e *ui.LogEntry) {
			req := createBackupStorageRequests[idx]
			msg := fmt.Sprintf("Creating backup of storage %q", req.UUID)
			e.SetMessage(msg)
			e.Start()
			details, err := s.service.CreateBackup(&req)
			if err != nil {
				e.SetMessage(ui.LiveLogEntryErrorColours.Sprintf("%s: failed", msg))
				e.SetDetails(err.Error(), "error: ")
			} else {
				e.SetMessage(fmt.Sprintf("%s: done", msg))
				e.SetDetails(details.UUID, "UUID: ")
				mu.Lock()
				numOk++
				storageDetails = append(storageDetails, details)
				mu.Unlock()
			}
		}
		ui.StartWorkQueue(ui.WorkQueueConfig{
			NumTasks:           len(createBackupStorageRequests),
			MaxConcurrentTasks: 5,
			EnableUI:           s.Config().InteractiveUI(),
		}, handler)
		if numOk != len(createBackupStorageRequests) {
			return nil, fmt.Errorf("number of backup creations that failed: %d", len(createBackupStorageRequests)-numOk)
		}

		return storageDetails, nil
	}
}
