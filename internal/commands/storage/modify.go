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
	backupTime      string
	backupInterval  string
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
	s.SetPositionalArgHelp(positionalArgHelp)
	s.ArgCompletion(GetArgCompFn(s.service))
	s.params = modifyParams{ModifyStorageRequest: request.ModifyStorageRequest{}}

	flagSet := &pflag.FlagSet{}
	flagSet.StringVar(&s.params.Title, "title", DefaultModifyParams.Title, "Storage title")
	flagSet.IntVar(&s.params.Size, "size", DefaultModifyParams.Size, "Size of the storage in GiB")
	flagSet.StringVar(&s.params.backupTime, "backup-time", s.params.backupTime, "The time when to create a backup in HH:MM. Empty value means no backups.")
	flagSet.StringVar(&s.params.backupInterval, "backup-interval", "", "The interval of the backup.\nAvailable: daily,mon,tue,wed,thu,fri,sat,sun")
	flagSet.IntVar(&s.params.backupRetention, "backup-retention", 0, "How long to store the backups in days. The accepted range is 1-1095")

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
			} else {
				if newBUR.Interval == "" {
					newBUR.Interval = defaultBackupRuleParams.Interval
				}
				if newBUR.Retention == 0 {
					newBUR.Retention = defaultBackupRuleParams.Retention
				}
				req.BackupRule = newBUR
			}
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

	return nil
}

func (s *modifyCommand) MakeExecuteCommand() func(args []string) (interface{}, error) {
	return func(args []string) (interface{}, error) {

		return Request{
			BuildRequest: func(uuid string) (interface{}, error) {
				req := s.params.ModifyStorageRequest
				if err := setBackupFields(uuid, s.params, s.service, &req); err != nil {
					return nil, err
				}
				req.UUID = uuid
				return &req, nil
			},
			Service: s.service,
			Handler: ui.HandleContext{
				RequestID:     func(in interface{}) string { return in.(*request.ModifyStorageRequest).UUID },
				MaxActions:    maxStorageActions,
				InteractiveUI: s.Config().InteractiveUI(),
				ActionMsg:     "Modifying storage",
				Action: func(req interface{}) (interface{}, error) {
					return s.service.ModifyStorage(req.(*request.ModifyStorageRequest))
				},
			},
		}.Send(args)
	}
}
