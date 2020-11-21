package storage

import (
	"github.com/UpCloudLtd/cli/internal/commands"
	"github.com/UpCloudLtd/cli/internal/ui"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/request"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/service"
	"github.com/spf13/pflag"
)

type modifyCommand struct {
	*commands.BaseCommand
	service service.Storage
	params  modifyParams
}

type modifyParams struct {
	request.ModifyStorageRequest
}

var DefaultModifyParams = &modifyParams{
	ModifyStorageRequest: request.ModifyStorageRequest{BackupRule: &upcloud.BackupRule{}},
}

func ModifyCommand(service service.Storage) commands.Command {
	return &modifyCommand{
		BaseCommand: commands.New("modify", "Modify a storage"),
		service:     service,
	}
}

func (s *modifyCommand) InitCommand() {
	s.params = modifyParams{ModifyStorageRequest: request.ModifyStorageRequest{BackupRule: &upcloud.BackupRule{}}}

	flagSet := &pflag.FlagSet{}
	flagSet.StringVar(&s.params.Title, "title", DefaultModifyParams.Title, "Storage title")
	flagSet.IntVar(&s.params.Size, "size", DefaultModifyParams.Size, "Size of the storage in GiB")
	flagSet.StringVar(&s.params.BackupRule.Time, "backup-time", DefaultModifyParams.BackupRule.Time, "The time when to create a backup in HH:MM. Empty value means no backups.")
	flagSet.StringVar(&s.params.BackupRule.Interval, "backup-interval", DefaultModifyParams.BackupRule.Interval, "The interval of the backup.\nAvailable: daily,mon,tue,wed,thu,fri,sat,sun")
	flagSet.IntVar(&s.params.BackupRule.Retention, "backup-retention", DefaultModifyParams.BackupRule.Retention, "How long to store the backups in days. The accepted range is 1-1095")

	s.AddFlags(flagSet)
}

func (s *modifyCommand) MakeExecuteCommand() func(args []string) (interface{}, error) {
	return func(args []string) (interface{}, error) {

		if s.params.BackupRule.Retention == DefaultModifyParams.BackupRule.Retention ||
			s.params.BackupRule.Time == DefaultModifyParams.BackupRule.Time ||
			s.params.BackupRule.Interval == DefaultModifyParams.BackupRule.Interval {
			s.params.BackupRule = nil
		}

		return Request{
			BuildRequest: func(storage *upcloud.Storage) interface{} {
				req := s.params.ModifyStorageRequest
				req.UUID = storage.UUID
				return &req
			},
			Service: s.service,
			HandleContext: ui.HandleContext{
				RequestId:  func(in interface{}) string { return in.(*request.ModifyStorageRequest).UUID },
				MaxActions: maxStorageActions,
				ActionMsg:  "Modifying storage",
				Action: func(req interface{}) (interface{}, error) {
					return s.service.ModifyStorage(req.(*request.ModifyStorageRequest))
				},
			},
		}.Send(args)
	}
}
