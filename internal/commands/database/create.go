package database

import (
	"encoding/json"
	"fmt"
	"slices"
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
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

type createCommand struct {
	*commands.BaseCommand
	params createParams
	wait   config.OptionalBoolean
}

func CreateCommand() commands.Command {
	return &createCommand{
		BaseCommand: commands.New(
			"create",
			"Create a new database",
			`upctl database create \
				--title mydb \
				--zone fi-hel1 \
				--hostname-prefix mydb`,
			`upctl database create \
				--title mydb \
				--zone fi-hel1 \
				--type pg \
				--hostname-prefix mydb \
				--termination-protection \
				--label env=dev \
				--property max_connections=200`,
		),
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
	labels                []string
	networks              []string
	terminationProtection config.OptionalBoolean
	dbType                string
	properties            []string
}

func (s *createParams) processParams(t *upcloud.ManagedDatabaseType) error {
	if len(s.labels) > 0 {
		labelSlice, err := labels.StringsToSliceOfLabels(s.labels)
		if err != nil {
			return err
		}
		s.Labels = labelSlice
	}

	if s.dbType != "" {
		s.Type = upcloud.ManagedDatabaseServiceType(s.dbType)
	}

	if len(s.properties) > 0 {
		props, err := processProperties(s.properties, t)
		if err != nil {
			return fmt.Errorf("invalid properties: %w", err)
		}
		s.Properties = props
	}

	if len(s.networks) > 0 {
		networks, err := processNetworks(s.networks)
		if err != nil {
			return fmt.Errorf("invalid networks: %w", err)
		}
		s.Networks = networks
	}

	if s.terminationProtection.IsSet() {
		terminationProtection := s.terminationProtection.Value()
		s.TerminationProtection = &terminationProtection
	}
	return nil
}

func isStringProperty(key upcloud.ManagedDatabasePropertyKey, t *upcloud.ManagedDatabaseType) bool {
	if propType, ok := t.Properties[string(key)].Type.(string); ok {
		return propType == "string"
	}

	if propType, ok := t.Properties[string(key)].Type.([]string); ok {
		return slices.Contains(propType, "string")
	}

	return false
}

func processProperties(in []string, t *upcloud.ManagedDatabaseType) (request.ManagedDatabasePropertiesRequest, error) {
	resp := request.ManagedDatabasePropertiesRequest{}
	for _, prop := range in {
		parts := strings.SplitN(prop, "=", 2)
		if len(parts) != 2 {
			return resp, fmt.Errorf("invalid property format: %s, expected key=value", prop)
		}

		key := upcloud.ManagedDatabasePropertyKey(parts[0])
		value := parts[1]

		// Remove quotes from the start and end of the value if they exist
		if (strings.HasPrefix(value, "\"") && strings.HasSuffix(value, "\"")) ||
			(strings.HasPrefix(value, "'") && strings.HasSuffix(value, "'")) {
			value = value[1 : len(value)-1]
		}

		// Handle numerical string values, e.g. Postgres version
		if isStringProperty(key, t) {
			resp.Set(key, value)
			continue
		}

		var parsedValue interface{}
		if err := json.Unmarshal([]byte(value), &parsedValue); err != nil {
			resp.Set(key, value) // Set as plain string if parsing fails
		} else {
			resp.Set(key, parsedValue)
		}
	}
	return resp, nil
}

func processNetworks(in []string) ([]upcloud.ManagedDatabaseNetwork, error) {
	var networks []upcloud.ManagedDatabaseNetwork
	for _, netStr := range in {
		network := upcloud.ManagedDatabaseNetwork{}
		pairs := strings.Split(netStr, ",")

		for _, pair := range pairs {
			parts := strings.SplitN(pair, "=", 2)
			if len(parts) != 2 {
				return nil, fmt.Errorf("invalid network format: %s, expected key=value", pair)
			}

			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])

			switch key {
			case "family":
				network.Family = value
			case "name":
				network.Name = value
			case "type":
				network.Type = value
			case "uuid":
				network.UUID = &value
			default:
				return nil, fmt.Errorf("unknown network parameter: %s", key)
			}
		}
		networks = append(networks, network)
	}
	return networks, nil
}

// InitCommand implements commands.InitializeCommand
func (s *createCommand) InitCommand() {
	flags := &pflag.FlagSet{}
	s.params = createParams{CreateManagedDatabaseRequest: request.CreateManagedDatabaseRequest{}}
	def := defaultCreateParams
	flags.StringVar(&s.params.HostNamePrefix, "hostname-prefix", def.HostNamePrefix, "A host name prefix for the database")
	flags.StringVar(&s.params.Title, "title", def.Title, "A short, informational description.")
	flags.StringVar(&s.params.Plan, "plan", def.Plan, "Plan to use for the database. Run `upctl database plans [database type]` to list all available plans.")
	flags.StringVar(&s.params.Zone, "zone", def.Zone, namedargs.ZoneDescription("database"))
	flags.StringVar(&s.params.dbType, "type", string(def.Type), "Type of the database")
	flags.StringVar(&s.params.Maintenance.DayOfWeek, "maintenance-dow", def.Maintenance.DayOfWeek, "Full name of weekday in English, lower case(sunday) for automatic maintenance day of the week. Set randomly if not provided.")
	flags.StringVar(&s.params.Maintenance.Time, "maintenance-time", def.Maintenance.Time, "Database time in UTC of automatic maintenance HH:MM:SS. Set randomly if not provided.")
	flags.StringSliceVar(&s.params.labels, "label", def.labels, "Labels to describe the database in `key=value` format, multiple can be declared.\nUsage: --label env=dev\n\n--label owner=operations")
	flags.StringArrayVar(&s.params.networks, "network", def.networks, "A network interface for the database, multiple can be declared.\nUsage: --network name=network-name,family=IPv4,type=private,uuid=030e83d2-d413-4d19-b1c9-af05cdb60c1f")
	config.AddEnableOrDisableFlag(flags, &s.params.terminationProtection, def.terminationProtection.Value(), "termination-protection", "termination protection to prevent the database instance from being powered off or deleted")

	flags.StringArrayVar(&s.params.properties, "property", nil, "Properties for the database in `key=value` format. Can be specified multiple times.")
	config.AddToggleFlag(flags, &s.wait, "wait", false, "Wait for database to be in running state before returning.")

	s.AddFlags(flags)

	commands.Must(s.Cobra().MarkFlagRequired("title"))
	commands.Must(s.Cobra().MarkFlagRequired("zone"))
	commands.Must(s.Cobra().MarkFlagRequired("hostname-prefix"))
	for _, flag := range []string{"hostname-prefix", "title", "plan", "maintenance-dow", "maintenance-time", "label", "network", "property"} {
		commands.Must(s.Cobra().RegisterFlagCompletionFunc(flag, cobra.NoFileCompletions))
	}
}

func (s *createCommand) InitCommandWithConfig(cfg *config.Config) {
	commands.Must(s.Cobra().RegisterFlagCompletionFunc("type", namedargs.CompletionFunc(completion.DatabaseType{}, cfg)))
	commands.Must(s.Cobra().RegisterFlagCompletionFunc("zone", namedargs.CompletionFunc(completion.Zone{}, cfg)))
}

// ExecuteWithoutArguments implements commands.NoArgumentCommand
func (s *createCommand) ExecuteWithoutArguments(exec commands.Executor) (output.Output, error) {
	svc := exec.All()
	msg := fmt.Sprintf("Creating database %v", s.params.Title)
	exec.PushProgressStarted(msg)

	t, err := svc.GetManagedDatabaseServiceType(exec.Context(), &request.GetManagedDatabaseServiceTypeRequest{
		Type: s.params.dbType,
	})
	if err != nil {
		return commands.HandleError(exec, msg, err)
	}

	if err := s.params.processParams(t); err != nil {
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
