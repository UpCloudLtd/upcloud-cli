package storage

import (
	"fmt"
	"time"

	"github.com/UpCloudLtd/cli/internal/commands"
	"github.com/UpCloudLtd/cli/internal/output"
	"github.com/UpCloudLtd/cli/internal/ui"

	"github.com/UpCloudLtd/upcloud-go-api/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/request"
	"github.com/spf13/pflag"
)

// CreateCommand creates the "storage create" command
func CreateCommand() commands.Command {
	return &createCommand{
		BaseCommand: commands.New("create", "Create a storage"),
	}
}

var defaultCreateParams = &createParams{
	CreateStorageRequest: request.CreateStorageRequest{
		Size: 10,
		Tier: upcloud.StorageTierMaxIOPS,
		BackupRule: &upcloud.BackupRule{
			Interval:  upcloud.BackupRuleIntervalDaily,
			Retention: 7,
		},
	},
}

func newCreateParams() createParams {
	return createParams{
		CreateStorageRequest: request.CreateStorageRequest{
			BackupRule: &upcloud.BackupRule{},
		},
	}
}

type createParams struct {
	request.CreateStorageRequest
	backupTime string
}

func (s *createParams) processParams() error {
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
	params  createParams
	flagSet *pflag.FlagSet
}

func createFlags(fs *pflag.FlagSet, dst, def *createParams) {
	fs.StringVar(&dst.Title, "title", def.Title, "Storage title.")
	fs.IntVar(&dst.Size, "size", def.Size, "Size of the storage in GiB.")
	fs.StringVar(&dst.Zone, "zone", def.Zone, "Physical location of the storage. See zone listing for valid zones.")
	fs.StringVar(&dst.Tier, "tier", def.Tier, "Storage tier.")
	fs.StringVar(&dst.backupTime, "backup-time", def.backupTime, "The time when to create a backup in HH:MM. Empty value means no backups.")
	fs.StringVar(&dst.BackupRule.Interval, "backup-interval", def.BackupRule.Interval, "The interval of the backup.\nAvailable: daily,mon,tue,wed,thu,fri,sat,sun")
	fs.IntVar(&dst.BackupRule.Retention, "backup-retention", def.BackupRule.Retention, "How long to store the backups in days. The accepted range is 1-1095")
}

// InitCommand implements Command.InitCommand
func (s *createCommand) InitCommand() {
	s.flagSet = &pflag.FlagSet{}
	s.params = newCreateParams()
	createFlags(s.flagSet, &s.params, defaultCreateParams)
	s.AddFlags(s.flagSet)
}

// Execute implements command.Command
func (s *createCommand) Execute(exec commands.Executor, _ string) (output.Output, error) {
	svc := exec.Storage()

	if s.params.Size == 0 || s.params.Zone == "" || s.params.Title == "" {
		return nil, fmt.Errorf("size, title and zone are required")
	}

	if err := s.params.processParams(); err != nil {
		return nil, err
	}

	msg := "creating storage"
	logline := exec.NewLogEntry(msg)
	logline.StartedNow()

	res, err := svc.CreateStorage(&s.params.CreateStorageRequest)
	if err != nil {
		logline.SetMessage(ui.LiveLogEntryErrorColours.Sprintf("%s: failed (%v)", msg, err.Error()))
		logline.SetDetails(err.Error(), "error: ")
		return nil, err
	}

	logline.SetMessage(fmt.Sprintf("%s: success", msg))
	logline.MarkDone()

	return output.Marshaled{Value: res}, nil
}
