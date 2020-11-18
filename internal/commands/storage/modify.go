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

type modifyCommand struct {
	*commands.BaseCommand
	service service.Storage
	params  modifyParams
	flagSet *pflag.FlagSet
}

type modifyParams struct {
	request.ModifyStorageRequest
}

var DefaultModifyParams = &modifyParams{
	ModifyStorageRequest: request.ModifyStorageRequest{BackupRule: &upcloud.BackupRule{}},
}

func newModifyParams() modifyParams {
	return modifyParams{ModifyStorageRequest: request.ModifyStorageRequest{BackupRule: &upcloud.BackupRule{}}}
}

func ModifyCommand(service service.Storage) commands.Command {
	return &modifyCommand{
		BaseCommand: commands.New("modify", "Modify a storage"),
		service:     service,
	}
}

func (s *modifyCommand) InitCommand() {
	s.params = newModifyParams()

	s.flagSet = &pflag.FlagSet{}
	s.flagSet.StringVar(&s.params.Title, "title", DefaultModifyParams.Title, "Storage title")
	s.flagSet.IntVar(&s.params.Size, "size", DefaultModifyParams.Size, "Size of the storage in GiB")
	s.flagSet.StringVar(&s.params.BackupRule.Time, "backup-time", DefaultModifyParams.BackupRule.Time, "The time when to create a backup in HH:MM. Empty value means no backups.")
	s.flagSet.StringVar(&s.params.BackupRule.Interval, "backup-interval", DefaultModifyParams.BackupRule.Interval, "The interval of the backup.\nAvailable: daily,mon,tue,wed,thu,fri,sat,sun")
	s.flagSet.IntVar(&s.params.BackupRule.Retention, "backup-retention", DefaultModifyParams.BackupRule.Retention, "How long to store the backups in days. The accepted range is 1-1095")

	s.AddFlags(s.flagSet)
}

func (s *modifyCommand) MakeExecuteCommand() func(args []string) (interface{}, error) {
	return func(args []string) (interface{}, error) {

		if len(args) < 1 {
			return nil, fmt.Errorf("storage uuid is required")
		}

		s.params.UUID = args[0]
		if s.params.BackupRule.Retention == DefaultModifyParams.BackupRule.Retention ||
			s.params.BackupRule.Time == DefaultModifyParams.BackupRule.Time ||
			s.params.BackupRule.Interval == DefaultModifyParams.BackupRule.Interval {
			s.params.BackupRule = nil
		}

		var modifyStorageRequests []request.ModifyStorageRequest
		for _, v := range args {
			s.params.ModifyStorageRequest.UUID = v
			modifyStorageRequests = append(modifyStorageRequests, s.params.ModifyStorageRequest)
		}

		var (
			mu             sync.Mutex
			numOk          int
			storageDetails []*upcloud.StorageDetails
		)
		handler := func(idx int, e *ui.LogEntry) {
			storageRequest := modifyStorageRequests[idx]
			msg := fmt.Sprintf("Modifying storage %q", storageRequest.UUID)
			e.SetMessage(msg)
			e.Start()
			details, err := s.service.ModifyStorage(&storageRequest)
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
			NumTasks:           len(modifyStorageRequests),
			MaxConcurrentTasks: 5,
			EnableUI:           s.Config().InteractiveUI(),
		}, handler)
		if numOk != len(modifyStorageRequests) {
			return nil, fmt.Errorf("number of storages that failed: %d", len(modifyStorageRequests)-numOk)
		}
		return storageDetails, nil
	}
}
