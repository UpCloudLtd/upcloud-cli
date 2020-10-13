package storage

import (
	"fmt"
	"sync"
	"time"

	"github.com/UpCloudLtd/upcloud-go-api/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/request"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/service"
	"github.com/spf13/pflag"

	"github.com/UpCloudLtd/cli/internal/commands"
	"github.com/UpCloudLtd/cli/internal/ui"
	"github.com/UpCloudLtd/cli/internal/upapi"
)

func CreateCommand() commands.Command {
	return &createCommand{
		BaseCommand: commands.New("create", "Create a storage"),
	}
}

var DefaultCreateParams = &createParams{
	CreateStorageRequest: request.CreateStorageRequest{
		Size: 10,
		Tier: "maxiops",
		BackupRule: &upcloud.BackupRule{
			Interval:  "daily",
			Retention: 7,
		},
	},
}

func newCreateParams() createParams {
	return createParams{CreateStorageRequest: request.CreateStorageRequest{BackupRule: &upcloud.BackupRule{}}}
}

type createParams struct {
	request.CreateStorageRequest
	backupTime string
}

func (s *createParams) processParams(srv *service.Service) error {
	if s.backupTime != "" {
		tv, err := time.Parse("15:04", s.backupTime)
		if err != nil {
			return fmt.Errorf("invalid backup time %q", s.backupTime)
		}
		s.BackupRule.Time = tv.Format("1504")
	} else {
		s.BackupRule = nil
	}
	return nil
}

type createCommand struct {
	*commands.BaseCommand
	service            *service.Service
	firstCreateStorage createParams
	flagSet            *pflag.FlagSet
}

func (s *createCommand) initService() {
	if s.service == nil {
		s.service = upapi.Service(s.Config())
	}
}

func createFlags(fs *pflag.FlagSet, dst, def *createParams) {
	fs.StringVar(&dst.Title, "title", def.Title, "Storage title")
	fs.IntVar(&dst.Size, "size", def.Size, "Size of the storage in GiB")
	fs.StringVar(&dst.Zone, "zone", def.Zone, "The zone to create the storage on")
	fs.StringVar(&dst.Tier, "tier", def.Tier, "Storage tier")
	fs.StringVar(&dst.backupTime, "backup-time", def.backupTime,
		"The time when to create a backup in HH:MM. Empty value means no backups.")
	fs.StringVar(&dst.BackupRule.Interval, "backup-interval", def.BackupRule.Interval,
		"The interval of the backup.\nAvailable: daily,mon,tue,wed,thu,fri,sat,sun")
	fs.IntVar(&dst.BackupRule.Retention, "backup-retention", def.BackupRule.Retention,
		"How long to store the backups in days. The accepted range is 1-1095")
}

func (s *createCommand) InitCommand() {
	s.flagSet = &pflag.FlagSet{}
	s.firstCreateStorage = newCreateParams()
	createFlags(s.flagSet, &s.firstCreateStorage, DefaultCreateParams)
	s.AddFlags(s.flagSet)
}

func (s *createCommand) MakeExecuteCommand() func(args []string) error {
	return func(args []string) error {
		s.initService()
		var createStorages []request.CreateStorageRequest
		if err := s.firstCreateStorage.processParams(s.service); err != nil {
			return err
		}
		createStorages = append(createStorages, s.firstCreateStorage.CreateStorageRequest)

		// Process additional storage create args
		var additionalCreateArgs = make([]string, 0, len(args))
		for i, arg := range args {
			if arg == "--" || i == len(args)-1 {
				if i == len(args)-1 && arg != "--" {
					additionalCreateArgs = append(additionalCreateArgs, arg)
				}
				if len(additionalCreateArgs) > 0 {
					fs := &pflag.FlagSet{}
					dst := newCreateParams()
					createFlags(fs, &dst, &s.firstCreateStorage)
					if err := fs.Parse(additionalCreateArgs); err != nil {
						return err
					}
					if err := dst.processParams(s.service); err != nil {
						return err
					}
					createStorages = append(createStorages, dst.CreateStorageRequest)
				}
				additionalCreateArgs = additionalCreateArgs[:0]
				continue
			}
			additionalCreateArgs = append(additionalCreateArgs, arg)
		}

		var (
			mu              sync.Mutex
			numOk           int
			createdStorages []*upcloud.StorageDetails
		)
		handler := func(idx int, e *ui.LogEntry) {
			storage := createStorages[idx]
			msg := fmt.Sprintf("Creating storage %q", storage.Title)
			e.SetMessage(msg)
			e.Start()
			details, err := s.service.CreateStorage(&storage)
			if err != nil {
				e.SetMessage(ui.LiveLogEntryErrorColours.Sprintf("%s: failed", msg))
				e.SetDetails(err.Error(), "error: ")
			} else {
				e.SetMessage(fmt.Sprintf("%s: done", msg))
				e.SetDetails(details.UUID, "UUID: ")
				mu.Lock()
				numOk++
				createdStorages = append(createdStorages, details)
				mu.Unlock()
			}
		}
		ui.StartWorkQueue(ui.WorkQueueConfig{
			NumTasks:           len(createStorages),
			MaxConcurrentTasks: 5,
			EnableUI:           s.Config().InteractiveUI(),
		}, handler)
		if numOk != len(createStorages) {
			return fmt.Errorf("number of storages that failed: %d", len(createStorages)-numOk)
		}
		return s.HandleOutput(createdStorages)
	}
}

func (s *createCommand) HandleOutput(out interface{}) error {
	results := out.([]*upcloud.StorageDetails)
	var uuids []string
	for _, res := range results {
		uuids = append(uuids, res.UUID)
	}

	if !s.Config().OutputHuman() {
		return s.BaseCommand.HandleOutput(uuids)
	}
	return nil
}
