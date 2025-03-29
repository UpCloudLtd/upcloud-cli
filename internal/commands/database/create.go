package database

import (
	"fmt"
	"strings"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/completion"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/config"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/labels"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/namedargs"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/output"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/ui"
	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud/request"
	"github.com/spf13/pflag"
)

type createCommand struct {
	*commands.BaseCommand
	params createParams
	wait   config.OptionalBoolean
}

func CreateCommand() commands.Command {
	return &createCommand{
		BaseCommand: commands.New("create", "Create a new database", "upctl database create"),
	}
}

var defaultCreateParams = createParams{
	CreateManagedDatabaseRequest: request.CreateManagedDatabaseRequest{
		Plan: "2x2xCPU-4GB-100GB",
		Type: upcloud.ManagedDatabaseServiceTypeMySQL,
	},
}

type createParams struct {
	request.CreateManagedDatabaseRequest
	labels              []string
	networks            []string
	terminateProtection bool
	dbType              string
	properties          []string
}

func (s *createParams) processParams() error {
	if len(s.labels) > 0 {
		labelSlice, err := labels.StringsToSliceOfLabels(s.labels)
		if err != nil {
			return err
		}
		s.Labels = labelSlice
	}

	if s.terminateProtection {
		s.TerminationProtection = &s.terminateProtection
	}

	if s.dbType != "" {
		s.Type = upcloud.ManagedDatabaseServiceType(s.dbType)
	}

	if len(s.properties) > 0 {
		var props request.ManagedDatabasePropertiesRequest
		for _, prop := range s.properties {
			parts := strings.SplitN(prop, "=", 2)
			if len(parts) != 2 {
				return fmt.Errorf("invalid property format: %s, expected key=value", prop)
			}
			key := upcloud.ManagedDatabasePropertyKey(parts[0])
			props.Set(key, parts[1])
		}

		s.Properties = props
	}

	return nil
}

// InitCommand implements commands.InitializeCommand
func (s *createCommand) InitCommand() {
	flags := &pflag.FlagSet{}
	s.params = createParams{CreateManagedDatabaseRequest: request.CreateManagedDatabaseRequest{}}
	def := defaultCreateParams
	flags.StringVar(&s.params.HostNamePrefix, "host-name-prefix", def.HostNamePrefix, "A host name prefix for the database")
	flags.StringVar(&s.params.Title, "title", def.Title, "A short, informational description.")
	flags.StringVar(&s.params.Plan, "plan", def.Plan, "Plan to use for the database. Run `upctl database plans [database type]` to list all available plans.")
	flags.StringVar(&s.params.Zone, "zone", def.Zone, namedargs.ZoneDescription("database"))
	flags.StringVar(&s.params.dbType, "type", string(def.Type), "Type of the database")
	flags.StringSliceVar(&s.params.labels, "labels", def.labels, "Labels to describe the database in `key=value` format, multiple can be declared.\nUsage: --label env=dev\n\n--label owner=operations")
	flags.StringSliceVar(&s.params.networks, "networks", def.networks, "A network interface for the database, multiple can be declared.\nUsage: --network family=IPv4,type=public\n\n--network type=private,network=037a530b-533e-4cef-b6ad-6af8094bb2bc,ip-address=10.0.0.1")
	flags.BoolVar(&s.params.terminateProtection, "terminate-protection", def.terminateProtection, "Prevents the database from being powered off, or deleted.")
	flags.StringSliceVar(&s.params.properties, "property", nil, "Properties for the database in the format key=value (can be specified multiple times)")
	config.AddToggleFlag(flags, &s.wait, "wait", false, "Wait for database to be in running state before returning.")

	s.AddFlags(flags)

	commands.Must(s.Cobra().MarkFlagRequired("title"))
	commands.Must(s.Cobra().MarkFlagRequired("zone"))
	commands.Must(s.Cobra().MarkFlagRequired("host-name-prefix"))
}

func (s *createCommand) InitCommandWithConfig(cfg *config.Config) {
	commands.Must(s.Cobra().RegisterFlagCompletionFunc("zone", namedargs.CompletionFunc(completion.Zone{}, cfg)))
}

// ExecuteWithoutArguments implements commands.NoArgumentCommand
func (s *createCommand) ExecuteWithoutArguments(exec commands.Executor) (output.Output, error) {
	svc := exec.All()
	msg := fmt.Sprintf("Creating database %v", s.params.Title)
	exec.PushProgressStarted(msg)

	if err := s.params.processParams(); err != nil {
		return nil, err
	}

	req := s.params.CreateManagedDatabaseRequest
	res, err := svc.CreateManagedDatabase(exec.Context(), &req)
	if err != nil {
		return commands.HandleError(exec, msg, err)
	}

	if s.wait.Value() {
		waitForManagedDatabaseState(res.UUID, upcloud.ManagedDatabaseStateRunning, exec, msg)
	} else {
		exec.PushProgressSuccess(msg)
	}

	return output.MarshaledWithHumanDetails{Value: res, Details: []output.DetailRow{
		{Title: "UUID", Value: res.UUID, Colour: ui.DefaultUUUIDColours},
	}}, nil
}
