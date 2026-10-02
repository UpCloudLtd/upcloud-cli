package database

import (
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/apierror"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/completion"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/format"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/labels"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/output"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/resolver"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/ui"
	upcloudv9 "github.com/UpCloudLtd/upcloud-go-api/v9/pkg/upcloud"
)

// ShowCommand creates the "database show" command
func ShowCommand() commands.Command {
	return &showCommand{
		BaseCommand: commands.New(
			"show",
			"Show database details",
			"upctl database show 9a8effcb-80e6-4a63-a7e5-066a6d093c14",
			"upctl database show my-pg-database",
			"upctl database show my-mysql-database",
		),
	}
}

type showCommand struct {
	*commands.BaseCommand
	resolver.CachingDatabase
	completion.Database
}

// Execute implements commands.MultipleArgumentCommand
func (s *showCommand) Execute(exec commands.Executor, uuidStr string) (output.Output, error) {
	client, err := v9Client(exec)
	if err != nil {
		return nil, err
	}
	id, err := parseDatabaseUUID(uuidStr)
	if err != nil {
		return nil, err
	}

	res, err := client.GetDatabaseWithResponse(exec.Context(), id)
	if err != nil {
		return nil, err
	}
	if res.JSON200 == nil {
		return nil, apierror.FromResponse(res.StatusCode(), res.Body)
	}
	db := res.JSON200
	legacyOutput, err := legacyDatabaseOutput(db)
	if err != nil {
		return nil, err
	}
	dbType := deref(db.Type)

	nodeRows := []output.TableRow{}
	for _, node := range deref(db.NodeStates) {
		nodeRows = append(nodeRows, output.TableRow{
			deref(node.Name),
			deref(node.Role),
			deref(node.State),
		})
	}

	componentsRows := []output.TableRow{}
	for _, component := range deref(db.Components) {
		componentsRows = append(componentsRows, output.TableRow{
			deref(component.Component),
			deref(component.Host),
			deref(component.Port),
			deref(component.Route),
			deref(component.Usage),
		})
	}

	maintenance := deref(db.Maintenance)
	detailSections := []output.DetailSection{
		{
			Title: "Overview:",
			Rows: []output.DetailRow{
				{Title: "UUID:", Value: uuidString(db.Uuid), Colour: ui.DefaultUUUIDColours},
				{Title: "Title:", Value: deref(db.Title)},
				{Title: "Name:", Value: deref(db.Name)},
				{Title: "Type:", Value: prettyDatabaseType(dbType)},
				{Title: "Version:", Value: getVersion(db), Format: format.PossiblyUnknownString},
				{Title: "Plan:", Value: deref(db.Plan)},
				{Title: "Zone:", Value: deref(db.Zone)},
				{Title: "State:", Value: deref(db.State), Format: format.DatabaseState},
				{Title: "Termination protection:", Value: deref(db.TerminationProtection), Format: format.Boolean},
			},
		},
	}

	if section, ok := planComponentsSection(db.PlanComponents); ok {
		detailSections = append(detailSections, section)
	}

	detailSections = append(detailSections,
		output.DetailSection{
			Title: "Maintenance schedule:",
			Rows: []output.DetailRow{
				{Title: "Weekday:", Value: deref(maintenance.Dow)},
				{Title: "Time:", Value: deref(maintenance.Time)},
			},
		},
		output.DetailSection{
			Title: "Authentication:",
			Rows: []output.DetailRow{
				{Title: "Service URI:", Value: deref(db.ServiceUri)},
				{Title: "Database name:", Value: mapString(db.ServiceUriParams, "dbname")},
				{Title: "Host:", Value: mapString(db.ServiceUriParams, "host")},
				{Title: "Password:", Value: mapString(db.ServiceUriParams, "password")},
				{Title: "Port:", Value: mapString(db.ServiceUriParams, "port")},
				{Title: "SSL mode:", Value: mapString(db.ServiceUriParams, "ssl_mode")},
				{Title: "User:", Value: mapString(db.ServiceUriParams, "user")},
			},
		},
	)

	if dbType == string(upcloudv9.DatabaseServiceTypeOpensearch) {
		acl, err := client.GetDatabaseAccessControlWithResponse(exec.Context(), id)
		if err != nil {
			return nil, err
		}
		if acl.JSON200 == nil {
			return nil, apierror.FromResponse(acl.StatusCode(), acl.Body)
		}

		detailSections = append(detailSections, output.DetailSection{
			Title: "Access control settings:",
			Rows: []output.DetailRow{
				{Title: "Access control:", Value: deref(acl.JSON200.AccessControl), Format: format.Boolean},
				{Title: "Extended access control:", Value: deref(acl.JSON200.ExtendedAccessControl), Format: format.Boolean},
			},
		})
	}

	// For JSON and YAML output, passthrough API response
	return output.MarshaledWithHumanOutput{
		Value: legacyOutput,
		Output: output.Combined{
			output.CombinedSection{
				Contents: output.Details{
					Sections: detailSections,
				},
			},
			labels.LabelsSection(deref(db.Labels), func(l upcloudv9.DatabaseLabelInformationResponse) (string, string) {
				return deref(l.Key), deref(l.Value)
			}, "database"),
			output.CombinedSection{
				Title: "Nodes:",
				Contents: output.Table{
					Columns: []output.TableColumn{
						{Key: "name", Header: "Name"},
						{Key: "role", Header: "Type"},
						{Key: "state", Header: "State"},
					},
					Rows: nodeRows,
				},
			},
			output.CombinedSection{
				Title: "Components:",
				Contents: output.Table{
					Columns: []output.TableColumn{
						{Key: "component", Header: "Component"},
						{Key: "host", Header: "Host"},
						{Key: "port", Header: "Port"},
						{Key: "route", Header: "Route"},
						{Key: "usage", Header: "Usage"},
					},
					Rows: componentsRows,
				},
			},
		},
	}, nil
}

func planComponentsSection(pc *upcloudv9.DatabasePlanComponentsResponse) (output.DetailSection, bool) {
	if pc == nil || pc.Compute == nil || pc.Compute.Name == nil {
		return output.DetailSection{}, false
	}

	return output.DetailSection{
		Title: "Plan components:",
		Rows: []output.DetailRow{
			{Title: "Compute:", Value: *pc.Compute.Name},
			{Title: "Node count:", Value: deref(pc.Compute.NodeCount)},
			{Title: "Storage per node (GiB):", Value: deref(deref(pc.Storage).TotalGib)},
			{Title: "Backups:", Value: string(deref(deref(pc.Backups).Name))},
		},
	}, true
}

var versionMetadataKeys = map[string]string{
	"mysql":      "mysql_version",
	"opensearch": "opensearch_version",
	"pg":         "pg_version",
	"valkey":     "valkey_version",
}

func getVersion(db *upcloudv9.DatabaseServiceInformationResponse) string {
	if db == nil || db.Metadata == nil {
		return ""
	}

	key, ok := versionMetadataKeys[deref(db.Type)]
	if !ok {
		return ""
	}
	return mapString(db.Metadata, key)
}

func prettyDatabaseType(serviceType string) string {
	switch upcloudv9.DatabaseServiceType(serviceType) {
	case upcloudv9.DatabaseServiceTypeMysql:
		return "MySQL"
	case upcloudv9.DatabaseServiceTypeOpensearch:
		return "OpenSearch"
	case upcloudv9.DatabaseServiceTypePg:
		return "PostgreSQL"
	default:
		return serviceType
	}
}
