package storage

import (
	"fmt"
	"github.com/UpCloudLtd/cli/internal/commands"
	"github.com/UpCloudLtd/cli/internal/ui"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/request"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/service"
	"github.com/spf13/pflag"
	"time"
)

type modifyCommand struct {
	*commands.BaseCommand
	service service.Storage
	params  modifyParams
}

type modifyParams struct {
	request.ModifyStorageRequest
	backupTime string
	backupInterval string
	backupRetention int
}

var DefaultModifyParams = &modifyParams{
	ModifyStorageRequest: request.ModifyStorageRequest{},
}

var defaultBackupRuleParams = upcloud.BackupRule{
	Interval:  "daily",
	Retention: 7,
}

func ModifyCommand(service service.Storage) commands.Command {
	return &modifyCommand{
		BaseCommand: commands.New("modify", "Modify a storage"),
		service:     service,
	}
}

func (s *modifyCommand) InitCommand() {
	s.params = modifyParams{ModifyStorageRequest: request.ModifyStorageRequest{}}

	flagSet := &pflag.FlagSet{}
	flagSet.StringVar(&s.params.Title, "title", DefaultModifyParams.Title, "Storage title")
	flagSet.IntVar(&s.params.Size, "size", DefaultModifyParams.Size, "Size of the storage in GiB")
	flagSet.StringVar(&s.params.backupTime, "backup-time", s.params.backupTime, "The time when to create a backup in HH:MM. Empty value means no backups.")
	flagSet.StringVar(&s.params.backupInterval, "backup-interval", "", "The interval of the backup.\nAvailable: daily,mon,tue,wed,thu,fri,sat,sun")
	flagSet.IntVar(&s.params.backupRetention, "backup-retention", 0, "How long to store the backups in days. The accepted range is 1-1095")

	s.AddFlags(flagSet)
}

func setBackupFields(storage *upcloud.Storage, p modifyParams, service service.Storage, req *request.ModifyStorageRequest) error {

	details, err := service.GetStorageDetails(&request.GetStorageDetailsRequest{UUID: storage.UUID})
	if err != nil {return err}

	var tv time.Time
	if p.backupTime != "" {
		tv, err = time.Parse("15:04", p.backupTime)
		if err != nil {
			return fmt.Errorf("invalid backup time %q", p.backupTime)
		}
	}

	if details.BackupRule == nil || details.BackupRule.Time == "" {
		if p.backupTime == "" && (p.backupInterval != "" || p.backupRetention != 0) {
			return fmt.Errorf("backup-time must be provided")
		}

		req.BackupRule = &upcloud.BackupRule{
			Time:      tv.Format("1504"),
		}
		if p.backupInterval == "" {
			req.BackupRule.Interval = defaultBackupRuleParams.Interval
		} else {
			req.BackupRule.Interval = p.backupInterval
		}
		if p.backupRetention == 0 {
			req.BackupRule.Retention = defaultBackupRuleParams.Retention
		} else {
			req.BackupRule.Retention = p.backupRetention
		}
	} else {
		req.BackupRule = details.BackupRule
		if p.backupTime != "" {
			req.BackupRule.Time = tv.Format("1504")
		}
		if p.backupInterval != "" {
			req.BackupRule.Interval = p.backupInterval
		}
		if p.backupRetention != 0 {
			req.BackupRule.Retention = p.backupRetention
		}
	}
	return nil
}

func (s *modifyCommand) MakeExecuteCommand() func(args []string) (interface{}, error) {
	return func(args []string) (interface{}, error) {

		return Request{
			BuildRequest: func(storage *upcloud.Storage) (interface{}, error) {
				req := s.params.ModifyStorageRequest
				if err := setBackupFields(storage, s.params, s.service, &req); err != nil {return nil, err}
				req.UUID = storage.UUID
				return &req, nil
			},
			Service: s.service,
			HandleContext: ui.HandleContext{
				RequestId:  func(in interface{}) string { return in.(*request.ModifyStorageRequest).UUID },
				MaxActions: maxStorageActions,
				InteractiveUi: s.Config().InteractiveUI(),
				ActionMsg:  "Modifying storage",
				Action: func(req interface{}) (interface{}, error) {
					return s.service.ModifyStorage(req.(*request.ModifyStorageRequest))
				},
			},
		}.Send(args)
	}
}
