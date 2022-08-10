package storage

import (
	"fmt"
	"time"

	"github.com/UpCloudLtd/upcloud-cli/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/internal/completion"
	"github.com/UpCloudLtd/upcloud-cli/internal/config"
	"github.com/UpCloudLtd/upcloud-cli/internal/output"
	"github.com/UpCloudLtd/upcloud-cli/internal/resolver"

	"github.com/UpCloudLtd/upcloud-go-api/v4/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/v4/upcloud/request"
	"github.com/UpCloudLtd/upcloud-go-api/v4/upcloud/service"
	"github.com/spf13/pflag"
)

type modifyCommand struct {
	*commands.BaseCommand
	completion.Storage
	resolver.CachingStorage
	params                        modifyParams
	autoresizePartitionFilesystem config.OptionalBoolean
}

type modifyParams struct {
	request.ModifyStorageRequest
	backupTime      string
	backupInterval  string
	backupRetention int
}

var defaultModifyParams = &modifyParams{
	ModifyStorageRequest: request.ModifyStorageRequest{},
}

var defaultBackupRuleParams = upcloud.BackupRule{
	Interval:  "daily",
	Retention: 7,
}

// ModifyCommand creates the "storage modify" command
func ModifyCommand() commands.Command {
	return &modifyCommand{
		BaseCommand: commands.New(
			"modify",
			"Modify a storage",
			`upctl storage modify 01271548-2e92-44bb-9774-d282508cc762 --title "My Storage" --size 20`,
			`upctl storage modify "My Storage" --size 25`,
		),
	}
}

// MaximumExecutions implements command.Command
func (s *modifyCommand) MaximumExecutions() int {
	return maxStorageActions
}

// InitCommand implements Command.InitCommand
func (s *modifyCommand) InitCommand() {
	s.params = modifyParams{ModifyStorageRequest: request.ModifyStorageRequest{}}

	flagSet := &pflag.FlagSet{}
	flagSet.StringVar(&s.params.Title, "title", defaultModifyParams.Title, "A short, informational description.")
	flagSet.IntVar(&s.params.Size, "size", defaultModifyParams.Size, "Size of the storage (GiB).")
	flagSet.StringVar(&s.params.backupTime, "backup-time", s.params.backupTime, "The time when to create a backup in HH:MM. Empty value means no backups.")
	flagSet.StringVar(&s.params.backupInterval, "backup-interval", "", "The interval of the backup.\nAvailable: daily,mon,tue,wed,thu,fri,sat,sun")
	flagSet.IntVar(&s.params.backupRetention, "backup-retention", 0, "How long to store the backups in days. The accepted range is 1-1095.")
	config.AddEnableOrDisableFlag(flagSet, &s.autoresizePartitionFilesystem, false, "filesystem-autoresize", "automatic resize of partition and filesystem when modifying storage size. Note that before the resize attempt is made, backup of the storage will be taken. If the resize attempt fails, the backup will be used to restore the storage and then deleted. If the resize attempt succeeds, backup will be kept. Taking and keeping backups incure costs.")

	s.AddFlags(flagSet)
}

func setBackupFields(storageUUID string, p modifyParams, service service.Storage, req *request.ModifyStorageRequest) error {
	details, err := service.GetStorageDetails(&request.GetStorageDetailsRequest{UUID: storageUUID})
	if err != nil {
		return err
	}

	var tv time.Time
	if p.backupTime != "" {
		tv, err = time.Parse("15:04", p.backupTime)
		if err != nil {
			return fmt.Errorf("invalid backup time %q", p.backupTime)
		}
		p.backupTime = tv.Format("1504")
	}

	var newBUR *upcloud.BackupRule
	if p.backupTime != "" || p.backupInterval != "" || p.backupRetention != 0 {
		newBUR = &upcloud.BackupRule{
			Interval:  p.backupInterval,
			Time:      p.backupTime,
			Retention: p.backupRetention,
		}
	}

	if details.BackupRule.Time == "" {
		if newBUR != nil {
			if newBUR.Time == "" {
				return fmt.Errorf("backup-time is required")
			}
			if newBUR.Interval == "" {
				newBUR.Interval = defaultBackupRuleParams.Interval
			}
			if newBUR.Retention == 0 {
				newBUR.Retention = defaultBackupRuleParams.Retention
			}
			req.BackupRule = newBUR
		} else {
			req.BackupRule = nil
		}
	} else {
		req.BackupRule = details.BackupRule
		if p.backupTime != "" {
			req.BackupRule.Time = p.backupTime
		}
		if p.backupInterval != "" {
			req.BackupRule.Interval = p.backupInterval
		}
		if p.backupRetention != 0 {
			req.BackupRule.Retention = p.backupRetention
		}
	}

	req.UUID = storageUUID

	return nil
}

// Execute implements commands.MultipleArgumentCommand
func (s *modifyCommand) Execute(exec commands.Executor, uuid string) (output.Output, error) {
	if s.autoresizePartitionFilesystem.Value() && s.params.Size == 0 {
		return nil, fmt.Errorf("filesystem autoresize is enabled, but new size is not specified")
	}

	svc := exec.Storage()
	msg := fmt.Sprintf("modifing storage %v", uuid)
	logline := exec.NewLogEntry(msg)

	logline.StartedNow()

	req := s.params.ModifyStorageRequest
	if err := setBackupFields(uuid, s.params, svc, &req); err != nil {
		return commands.HandleError(logline, fmt.Sprintf("%s: failed", msg), err)
	}

	res, err := svc.ModifyStorage(&req)
	if err != nil {
		return commands.HandleError(logline, fmt.Sprintf("%s: failed", msg), err)
	}

	// If autoresize is not enabled, then just consider the whole operation done and output the modify API call response
	if !s.autoresizePartitionFilesystem.Value() {
		logline.SetMessage(fmt.Sprintf("%s: done", msg))
		logline.MarkDone()

		return output.OnlyMarshaled{Value: res}, nil
	}

	logline.SetMessage(fmt.Sprintf("%s: resizing partition and filesystem", msg))
	backup, err := svc.ResizeStorageFilesystem(&request.ResizeStorageFilesystemRequest{UUID: uuid})
	// If there was an error during resize attempt, we consider the overall modify operation successful and just log warning about failed resize
	if err != nil {
		logline.SetMessage(fmt.Sprintf("%s: partially done", msg))
		logline.SetDetails(fmt.Sprintf("Partition and filesystem resize failed; storage was restored using backup taken right before resize attempt (%s)", err.Error()), "Error: ")
		logline.MarkWarning()

		return output.OnlyMarshaled{Value: res}, nil
	}

	logline.SetMessage(fmt.Sprintf("%s: done", msg))
	logline.MarkDone()

	out := struct {
		upcloud.StorageDetails
		LatestResizeBackup string `json:"latest_resize_backup,omitempty"`
	}{
		StorageDetails:     *res,
		LatestResizeBackup: backup.UUID,
	}

	return output.OnlyMarshaled{Value: out}, nil
}
