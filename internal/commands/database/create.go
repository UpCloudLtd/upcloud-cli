package database

import (
	"encoding/json"
	"errors"
	"fmt"
	"slices"
	"strconv"
	"strings"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/apierror"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/completion"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/config"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/labels"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/namedargs"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/output"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/ui"
	upcloudv9 "github.com/UpCloudLtd/upcloud-go-api/v9/pkg/upcloud"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

const (
	defaultPlanCompute    = "rdb.standard.2CPU-8GB"
	defaultPlanNodeCount  = 2
	defaultPlanStorageGiB = 100
	defaultPlanBackups    = "regular"
)

var planComponentFlags = []string{"plan-compute", "plan-node-count", "plan-storage-gib", "plan-backups"}

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
				--plan-compute rdb.standard.2CPU-8GB \
				--plan-node-count 2 \
				--plan-storage-gib 120 \
				--plan-backups regular`,
			`upctl database create \
				--title mydb \
				--zone fi-hel1 \
				--type pg \
				--hostname-prefix mydb \
				--plan 2x2xCPU-4GB-100GB \
				--enable-termination-protection \
				--label env=dev \
				--property max_connections=200`,
		),
	}
}

type createParams struct {
	title                 string
	zone                  string
	hostnamePrefix        string
	dbType                string
	plan                  string
	planCompute           string
	planNodeCountValue    string
	planNodeCount         int
	planStorageGiBValue   string
	planStorageGiB        int
	planBackups           string
	maintenanceDow        string
	maintenanceTime       string
	labels                []string
	networks              []string
	properties            []string
	terminationProtection config.OptionalBoolean
}

func supportsPlanComponents(dbType string) bool {
	return dbType == string(upcloudv9.DatabaseServiceTypePg) || dbType == string(upcloudv9.DatabaseServiceTypeMysql)
}

func (p *createParams) request(properties map[string]propertySchema) (upcloudv9.DatabaseServiceCreateOpenAPI, error) {
	req := upcloudv9.DatabaseServiceCreateOpenAPI{
		HostnamePrefix: p.hostnamePrefix,
		Title:          p.title,
		Type:           upcloudv9.DatabaseServiceType(p.dbType),
		Zone:           p.zone,
	}

	if p.maintenanceDow != "" || p.maintenanceTime != "" {
		req.Maintenance = &struct {
			Dow  upcloudv9.DatabaseMaintenanceDow  `json:"dow"`
			Time upcloudv9.DatabaseMaintenanceTime `json:"time"`
		}{
			Dow:  upcloudv9.DatabaseMaintenanceDow(p.maintenanceDow),
			Time: p.maintenanceTime,
		}
	}

	if len(p.labels) > 0 {
		v9Labels, err := labels.StringsToLabels(p.labels, func(key, value string) upcloudv9.DatabaseLabelCreate {
			return upcloudv9.DatabaseLabelCreate{Key: key, Value: value}
		})
		if err != nil {
			return req, err
		}
		req.Labels = &v9Labels
	}

	if len(p.properties) > 0 {
		props, err := processProperties(p.properties, properties)
		if err != nil {
			return req, fmt.Errorf("invalid properties: %w", err)
		}
		// Properties are a generated union of per-engine types; setting the raw JSON keeps untyped key=value input.
		raw, err := json.Marshal(props)
		if err != nil {
			return req, fmt.Errorf("invalid properties: %w", err)
		}
		var v9Props upcloudv9.DatabaseServiceCreateOpenAPI_Properties
		if err := v9Props.UnmarshalJSON(raw); err != nil {
			return req, fmt.Errorf("invalid properties: %w", err)
		}
		req.Properties = &v9Props
	}

	if len(p.networks) > 0 {
		networks, err := processNetworks(p.networks)
		if err != nil {
			return req, fmt.Errorf("invalid networks: %w", err)
		}
		req.Networks = &networks
	}

	if p.terminationProtection.IsSet() {
		terminationProtection := p.terminationProtection.Value()
		req.TerminationProtection = &terminationProtection
	}
	return req, nil
}

// validateZone checks that the database type is offered in the requested zone.
func (p *createParams) validateZone(catalog *upcloudv9.DatabasePlansResponse) error {
	for _, serviceType := range catalog.ServiceTypes {
		if string(serviceType.Type) != p.dbType {
			continue
		}
		if !slices.Contains(serviceType.Zones, p.zone) {
			return fmt.Errorf("--zone %q is not available for database type %q, valid zones are %s", p.zone, p.dbType, strings.Join(serviceType.Zones, ", "))
		}
		return nil
	}
	return fmt.Errorf("database type %q is not available in the plan catalog", p.dbType)
}

func (p *createParams) validateLegacyPlan(plans []upcloudv9.DatabaseServicePlanResponse) error {
	for _, plan := range plans {
		if plan.Plan == nil || *plan.Plan != p.plan {
			continue
		}
		if plan.Zones == nil || plan.Zones.Zone == nil {
			return fmt.Errorf("availability zones are not available for database plan %q", p.plan)
		}

		zones := make([]string, 0, len(*plan.Zones.Zone))
		for _, zone := range *plan.Zones.Zone {
			if zone.Name != nil {
				zones = append(zones, *zone.Name)
			}
		}
		if !slices.Contains(zones, p.zone) {
			return fmt.Errorf("--zone %q is not available for database plan %q, valid zones are %s", p.zone, p.plan, strings.Join(zones, ", "))
		}
		return nil
	}

	hint := fmt.Sprintf("run \"upctl database plans %s\" to list valid values", p.dbType)
	if supportsPlanComponents(p.dbType) {
		hint = fmt.Sprintf("run \"upctl database plans %s --show-legacy\" to list legacy plans", p.dbType)
	}
	return fmt.Errorf("--plan %q is not available for database type %q, %s", p.plan, p.dbType, hint)
}

// applyPlanComponents validates the plan components against the plan catalog and sets them on the request.
func (p *createParams) applyPlanComponents(catalog *upcloudv9.DatabasePlansResponse, changed func(string) bool, req *upcloudv9.DatabaseServiceCreateOpenAPI) error {
	hint := fmt.Sprintf("run \"upctl database plans %s\" to list valid values", p.dbType)
	for _, serviceType := range catalog.ServiceTypes {
		if string(serviceType.Type) != p.dbType {
			continue
		}
		for _, shape := range serviceType.ComputeShapes {
			if shape.Compute != p.planCompute {
				continue
			}

			nodeCounts := make([]string, 0, len(shape.NodeCounts))
			for _, n := range shape.NodeCounts {
				nodeCounts = append(nodeCounts, fmt.Sprint(n))
			}
			if !slices.Contains(nodeCounts, fmt.Sprint(p.planNodeCount)) {
				return fmt.Errorf("--plan-node-count %d is not available for compute shape %s, valid node counts are %s", p.planNodeCount, shape.Compute, strings.Join(nodeCounts, ", "))
			}
			req.PlanCompute = &p.planCompute
			req.PlanNodeCount = &p.planNodeCount

			if shape.Storage == nil {
				if changed("plan-storage-gib") {
					return fmt.Errorf("compute shape %s does not support --plan-storage-gib", shape.Compute)
				}
			} else {
				step := int(shape.Storage.StepGib)
				valid := false
				ranges := make([]string, 0, len(shape.Storage.Options))
				for _, option := range shape.Storage.Options {
					base, maxGiB := int(option.BaseGib), int(option.MaxGib)
					ranges = append(ranges, fmt.Sprintf("%d-%d", base, maxGiB))
					if p.planStorageGiB == base || (step > 0 && p.planStorageGiB > base && p.planStorageGiB <= maxGiB && (p.planStorageGiB-base)%step == 0) {
						valid = true
					}
				}
				if !valid {
					return fmt.Errorf("--plan-storage-gib %d is not available for compute shape %s, valid totals per node are %s GiB in steps of %d GiB", p.planStorageGiB, shape.Compute, strings.Join(ranges, ", "), step)
				}
				req.PlanStorageGib = &p.planStorageGiB
			}

			if shape.Backups == nil {
				if changed("plan-backups") {
					return fmt.Errorf("compute shape %s does not support --plan-backups", shape.Compute)
				}
			} else {
				tiers := make([]string, 0, len(*shape.Backups))
				for _, tier := range *shape.Backups {
					tiers = append(tiers, string(tier))
				}
				if !slices.Contains(tiers, p.planBackups) {
					return fmt.Errorf("--plan-backups %q is not available for compute shape %s, valid backup tiers are %s", p.planBackups, shape.Compute, strings.Join(tiers, ", "))
				}
				backups := upcloudv9.DatabaseServiceCreateOpenAPIPlanBackups(p.planBackups)
				req.PlanBackups = &backups
			}
			return nil
		}
		return fmt.Errorf("--plan-compute %q is not available for %s, %s", p.planCompute, p.dbType, hint)
	}
	return fmt.Errorf("database type %s has no component plans, use --plan", p.dbType)
}

func processProperties(in []string, schemas map[string]propertySchema) (map[string]any, error) {
	resp := map[string]any{}
	for _, prop := range in {
		parts := strings.SplitN(prop, "=", 2)
		if len(parts) != 2 {
			return resp, fmt.Errorf("invalid property format: %s, expected key=value", prop)
		}

		key := parts[0]
		value := parts[1]

		// Remove quotes from the start and end of the value if they exist
		if (strings.HasPrefix(value, "\"") && strings.HasSuffix(value, "\"")) ||
			(strings.HasPrefix(value, "'") && strings.HasSuffix(value, "'")) {
			value = value[1 : len(value)-1]
		}

		// Handle numerical string values, e.g. Postgres version
		if schemas[key].allowsString() {
			resp[key] = value
			continue
		}

		var parsedValue any
		if err := json.Unmarshal([]byte(value), &parsedValue); err != nil {
			resp[key] = value // Set as plain string if parsing fails
		} else {
			resp[key] = parsedValue
		}
	}
	return resp, nil
}

func processNetworks(in []string) ([]upcloudv9.DatabaseNetworkCreate, error) {
	var networks []upcloudv9.DatabaseNetworkCreate
	for _, netStr := range in {
		network := upcloudv9.DatabaseNetworkCreate{}
		pairs := strings.SplitSeq(netStr, ",")

		for pair := range pairs {
			parts := strings.SplitN(pair, "=", 2)
			if len(parts) != 2 {
				return nil, fmt.Errorf("invalid network format: %s, expected key=value", pair)
			}

			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])

			switch key {
			case "family":
				network.Family = upcloudv9.DatabaseNetworkFamily(value)
			case "name":
				network.Name = value
			case "type":
				network.Type = upcloudv9.DatabaseNetworkType(value)
			case "uuid":
				id, err := uuid.Parse(value)
				if err != nil {
					return nil, fmt.Errorf("invalid network uuid %q: %w", value, err)
				}
				network.Uuid = &id
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
	s.Cobra().Long = commands.WrapLongDescription(`Create a new database

For pg and mysql, select the plan by its components with --plan-compute, --plan-node-count, --plan-storage-gib, and --plan-backups. When any of these flags is given, the omitted ones use their defaults. Run "upctl database plans <type>" to list the valid compute shapes, node counts, storage ranges, and backup tiers. --plan-storage-gib is the total storage per node, not additional storage.

--plan is deprecated for pg and mysql, but still supported. Run "upctl database plans <type> --show-legacy" to list their legacy plan names. --plan is required for other database types; run "upctl database plans <type>" to list their plans. --plan cannot be combined with the --plan-* flags. Without any plan flag, pg and mysql use the default of every --plan-* flag.`)

	flags := &pflag.FlagSet{}
	s.params = createParams{}
	flags.StringVar(&s.params.hostnamePrefix, "hostname-prefix", "", "A host name prefix for the database")
	flags.StringVar(&s.params.title, "title", "", "A short, informational description.")
	flags.StringVar(&s.params.plan, "plan", "", "Deprecated for pg and mysql, use the --plan-* flags instead. Plan name, required for database types other than pg and mysql.")
	flags.StringVar(&s.params.planCompute, "plan-compute", defaultPlanCompute, "Compute shape for pg and mysql, see the Compute shape column of \"upctl database plans <type>\".")
	flags.StringVar(&s.params.planNodeCountValue, "plan-node-count", strconv.Itoa(defaultPlanNodeCount), "Number of nodes for pg and mysql.")
	flags.StringVar(&s.params.planStorageGiBValue, "plan-storage-gib", strconv.Itoa(defaultPlanStorageGiB), "Total storage per node in GiB for pg and mysql.")
	flags.StringVar(&s.params.planBackups, "plan-backups", defaultPlanBackups, "Backup tier for pg and mysql.")
	flags.StringVar(&s.params.zone, "zone", "", namedargs.ZoneDescription("database"))
	flags.StringVar(&s.params.dbType, "type", string(upcloudv9.DatabaseServiceTypeMysql), "Type of the database")
	flags.StringVar(&s.params.maintenanceDow, "maintenance-dow", "", "Full name of weekday in English, lower case(sunday) for automatic maintenance day of the week. Set randomly if not provided. Requires --maintenance-time.")
	flags.StringVar(&s.params.maintenanceTime, "maintenance-time", "", "Database time in UTC of automatic maintenance HH:MM:SS. Set randomly if not provided. Requires --maintenance-dow.")
	flags.StringSliceVar(&s.params.labels, "label", nil, "Labels to describe the database in `key=value` format, multiple can be declared.\nUsage: --label env=dev\n\n--label owner=operations")
	flags.StringArrayVar(&s.params.networks, "network", nil, "A network interface for the database, multiple can be declared.\nUsage: --network name=network-name,family=IPv4,type=private,uuid=030e83d2-d413-4d19-b1c9-af05cdb60c1f")
	config.AddEnableOrDisableFlag(flags, &s.params.terminationProtection, false, "termination-protection", "termination protection to prevent the database instance from being powered off or deleted")

	flags.StringArrayVar(&s.params.properties, "property", nil, "Properties for the database in `key=value` format. Can be specified multiple times.")
	config.AddToggleFlag(flags, &s.wait, "wait", false, "Wait for database to be in running state before returning.")

	s.AddFlags(flags)

	commands.Must(s.Cobra().MarkFlagRequired("title"))
	commands.Must(s.Cobra().MarkFlagRequired("zone"))
	commands.Must(s.Cobra().MarkFlagRequired("hostname-prefix"))
	s.Cobra().MarkFlagsRequiredTogether("maintenance-dow", "maintenance-time")
	for _, flag := range append([]string{"hostname-prefix", "title", "plan", "maintenance-dow", "maintenance-time", "label", "network", "property"}, planComponentFlags...) {
		commands.Must(s.Cobra().RegisterFlagCompletionFunc(flag, cobra.NoFileCompletions))
	}
}

func (s *createCommand) InitCommandWithConfig(cfg *config.Config) {
	commands.Must(s.Cobra().RegisterFlagCompletionFunc("type", namedargs.CompletionFunc(completion.DatabaseType{}, cfg)))
	commands.Must(s.Cobra().RegisterFlagCompletionFunc("zone", namedargs.CompletionFunc(completion.Zone{}, cfg)))
}

// usesPlanComponents reports whether the plan is selected by components, and rejects plan flags that cannot be used for the database type.
func (s *createCommand) usesPlanComponents() (bool, error) {
	flags := s.Cobra().Flags()
	componentFlags := slices.ContainsFunc(planComponentFlags, flags.Changed)
	switch {
	case flags.Changed("plan") && componentFlags:
		return false, errors.New("--plan cannot be combined with component plan flags (--plan-compute, --plan-node-count, --plan-storage-gib, or --plan-backups)")
	case componentFlags && !supportsPlanComponents(s.params.dbType):
		return false, fmt.Errorf("the --plan-* flags are only supported for pg and mysql, use --plan for database type %q", s.params.dbType)
	case flags.Changed("plan"):
		return false, nil
	case supportsPlanComponents(s.params.dbType):
		var err error
		s.params.planNodeCount, err = positivePlanInteger("plan-node-count", s.params.planNodeCountValue)
		if err != nil {
			return false, err
		}
		s.params.planStorageGiB, err = positivePlanInteger("plan-storage-gib", s.params.planStorageGiBValue)
		if err != nil {
			return false, err
		}
		return true, nil
	default:
		return false, fmt.Errorf("--plan is required for database type %q, run \"upctl database plans %s\" to list available plans", s.params.dbType, s.params.dbType)
	}
}

func positivePlanInteger(flag, value string) (int, error) {
	parsed, err := strconv.Atoi(value)
	if err != nil || parsed < 1 {
		return 0, fmt.Errorf("--%s must be a positive integer, got %q", flag, value)
	}
	return parsed, nil
}

// ExecuteWithoutArguments implements commands.NoArgumentCommand
func (s *createCommand) ExecuteWithoutArguments(exec commands.Executor) (output.Output, error) {
	useComponents, err := s.usesPlanComponents()
	if err != nil {
		return nil, err
	}
	if s.Cobra().Flags().Changed("plan") && supportsPlanComponents(s.params.dbType) {
		_, _ = fmt.Fprintln(s.Cobra().ErrOrStderr(), "Deprecation Warning: --plan is deprecated for pg and mysql; use --plan-compute, --plan-node-count, --plan-storage-gib, and --plan-backups instead.")
	}

	client, err := v9Client(exec)
	if err != nil {
		return nil, err
	}

	msg := fmt.Sprintf("Creating database %v", s.params.title)
	exec.PushProgressStarted(msg)

	properties, legacyPlans, err := getServiceTypeDetails(exec.Context(), client, s.params.dbType)
	if err != nil {
		return commands.HandleError(exec, msg, err)
	}

	req, err := s.params.request(properties)
	if err != nil {
		return nil, err
	}

	if useComponents {
		plans, err := client.ListDatabasePlansWithResponse(exec.Context())
		if err != nil {
			return commands.HandleError(exec, msg, err)
		}
		if plans.JSON200 == nil {
			return commands.HandleError(exec, msg, apierror.FromResponse(plans.StatusCode(), plans.Body))
		}
		if err := s.params.validateZone(plans.JSON200); err != nil {
			return commands.HandleError(exec, msg, err)
		}
		if err := s.params.applyPlanComponents(plans.JSON200, s.Cobra().Flags().Changed, &req); err != nil {
			return commands.HandleError(exec, msg, err)
		}
	} else {
		if err := s.params.validateLegacyPlan(legacyPlans); err != nil {
			return commands.HandleError(exec, msg, err)
		}
		req.Plan = &s.params.plan
	}

	res, err := client.CreateDatabaseWithResponse(exec.Context(), req)
	if err != nil {
		return commands.HandleError(exec, msg, err)
	}
	if res.JSON201 == nil {
		return commands.HandleError(exec, msg, apierror.FromResponse(res.StatusCode(), res.Body))
	}
	db := res.JSON201
	legacyOutput, err := legacyDatabaseOutput(db)
	if err != nil {
		return commands.HandleError(exec, msg, err)
	}
	if db.Uuid == nil {
		return commands.HandleError(exec, msg, errors.New("the API response did not include the database UUID"))
	}
	id := db.Uuid.String()

	if s.wait.Value() {
		WaitForManagedDatabaseState(id, databaseStateRunning, exec, msg)
	} else {
		exec.PushProgressSuccess(msg)
	}

	return output.MarshaledWithHumanDetails{Value: legacyOutput, Details: []output.DetailRow{
		{Title: "UUID", Value: id, Colour: ui.DefaultUUUIDColours},
	}}, nil
}
