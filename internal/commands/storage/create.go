package storage

import (
	"fmt"
	"time"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/completion"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/config"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/namedargs"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/output"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/ui"

	"github.com/UpCloudLtd/upcloud-go-api/v6/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/v6/upcloud/request"
	"github.com/spf13/pflag"
)

// CreateCommand creates the "storage create" command
func CreateCommand() commands.Command {
	return &createCommand{
		BaseCommand: commands.New(
			"create",
			"Create a storage",
			`upctl storage create --zone pl-waw1 --title "My Storage"`,
			"upctl storage create --zone pl-waw1 --title my_storage --size 20 --backup-interval wed --backup-retention 4",
		),
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

func applyCreateFlags(fs *pflag.FlagSet, dst, def *createParams) {
	fs.StringVar(&dst.Title, "title", def.Title, "A short, informational description.")
	fs.IntVar(&dst.Size, "size", def.Size, "Size of the storage in GiB.")
	fs.StringVar(&dst.Zone, "zone", def.Zone, namedargs.ZoneDescription("storage"))
	fs.StringVar(&dst.Tier, "tier", def.Tier, "Storage tier.")
	fs.StringVar(&dst.backupTime, "backup-time", def.backupTime, "The time when to create a backup in HH:MM. Empty value means no backups.")
	fs.StringVar(&dst.BackupRule.Interval, "backup-interval", def.BackupRule.Interval, "The interval of the backup.\nAvailable: daily,mon,tue,wed,thu,fri,sat,sun")
	fs.IntVar(&dst.BackupRule.Retention, "backup-retention", def.BackupRule.Retention, "How long to store the backups in days. The accepted range is 1-1095")
}

// InitCommand implements Command.InitCommand
func (s *createCommand) InitCommand() {
	s.flagSet = &pflag.FlagSet{}
	s.params = newCreateParams()
	applyCreateFlags(s.flagSet, &s.params, defaultCreateParams)

	s.AddFlags(s.flagSet)
	_ = s.Cobra().MarkFlagRequired("title")
	_ = s.Cobra().MarkFlagRequired("zone")
}

func (s *createCommand) InitCommandWithConfig(cfg *config.Config) {
	_ = s.Cobra().RegisterFlagCompletionFunc("zone", namedargs.CompletionFunc(completion.Zone{}, cfg))
}

// ExecuteWithoutArguments implements commands.NoArgumentCommand
func (s *createCommand) ExecuteWithoutArguments(exec commands.Executor) (output.Output, error) {
	svc := exec.Storage()

	if err := s.params.processParams(); err != nil {
		return nil, err
	}

	msg := fmt.Sprintf("Creating storage %s", s.params.Title)
	exec.PushProgressStarted(msg)

	res, err := svc.CreateStorage(exec.Context(), &s.params.CreateStorageRequest)
	if err != nil {
		return commands.HandleError(exec, msg, err)
	}

	exec.PushProgressSuccess(msg)

	return output.MarshaledWithHumanDetails{Value: res, Details: []output.DetailRow{
		{Title: "UUID", Value: res.UUID, Colour: ui.DefaultUUUIDColours},
	}}, nil
}
