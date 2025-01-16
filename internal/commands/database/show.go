package database

import (
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/completion"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/format"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/labels"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/output"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/resolver"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/ui"
	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud/request"
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
func (s *showCommand) Execute(exec commands.Executor, uuid string) (output.Output, error) {
	svc := exec.All()
	db, err := svc.GetManagedDatabase(exec.Context(), &request.GetManagedDatabaseRequest{UUID: uuid})
	if err != nil {
		return nil, err
	}

	nodeRows := []output.TableRow{}
	for _, node := range db.NodeStates {
		nodeRows = append(nodeRows, output.TableRow{
			node.Name,
			node.Role,
			node.State,
		})
	}

	componentsRows := []output.TableRow{}
	for _, component := range db.Components {
		componentsRows = append(componentsRows, output.TableRow{
			component.Component,
			component.Host,
			component.Port,
			component.Route,
			component.Usage,
		})
	}

	detailSections := []output.DetailSection{
		{
			Title: "Overview:",
			Rows: []output.DetailRow{
				{Title: "UUID:", Value: db.UUID, Colour: ui.DefaultUUUIDColours},
				{Title: "Title:", Value: db.Title},
				{Title: "Name:", Value: db.Name},
				{Title: "Type:", Value: prettyDatabaseType(db.Type)},
				{Title: "Version:", Value: getVersion(db), Format: format.PossiblyUnknownString},
				{Title: "Plan:", Value: db.Plan},
				{Title: "Zone:", Value: db.Zone},
				{Title: "State:", Value: db.State, Format: format.DatabaseState},
				{Title: "Termination protection:", Value: db.TerminationProtection, Format: format.Boolean},
			},
		},
		{
			Title: "Maintenance schedule:",
			Rows: []output.DetailRow{
				{Title: "Weekday:", Value: db.Maintenance.DayOfWeek},
				{Title: "Time:", Value: db.Maintenance.Time},
			},
		},
		{
			Title: "Authentication:",
			Rows: []output.DetailRow{
				{Title: "Service URI:", Value: db.ServiceURI},
				{Title: "Database name:", Value: db.ServiceURIParams.DatabaseName},
				{Title: "Host:", Value: db.ServiceURIParams.Host},
				{Title: "Password:", Value: db.ServiceURIParams.Password},
				{Title: "Port:", Value: db.ServiceURIParams.Port},
				{Title: "SSL mode:", Value: db.ServiceURIParams.SSLMode},
				{Title: "User:", Value: db.ServiceURIParams.User},
			},
		},
	}

	if db.Type == upcloud.ManagedDatabaseServiceTypeOpenSearch {
		acl, err := svc.GetManagedDatabaseAccessControl(exec.Context(), &request.GetManagedDatabaseAccessControlRequest{ServiceUUID: uuid})
		if err != nil {
			return nil, err
		}

		detailSections = append(detailSections, output.DetailSection{
			Title: "Access control settings:",
			Rows: []output.DetailRow{
				{Title: "Access control:", Value: acl.ACLsEnabled, Format: format.Boolean},
				{Title: "Extended access control:", Value: acl.ExtendedACLsEnabled, Format: format.Boolean},
			},
		})
	}

	// For JSON and YAML output, passthrough API response
	return output.MarshaledWithHumanOutput{
		Value: db,
		Output: output.Combined{
			output.CombinedSection{
				Contents: output.Details{
					Sections: detailSections,
				},
			},
			labels.GetLabelsSectionWithResourceType(db.Labels, "database"),
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

func getVersion(db *upcloud.ManagedDatabase) string {
	if db == nil || db.Metadata == nil {
		return ""
	}

	switch db.Type {
	case "mysql":
		return db.Metadata.MySQLVersion
	case "opensearch":
		return db.Metadata.OpenSearchVersion
	case "pg":
		return db.Metadata.PGVersion
	case "redis":
		return db.Metadata.RedisVersion //nolint:staticcheck // To be removed when Redis support has been removed
	}
	return ""
}

func prettyDatabaseType(serviceType upcloud.ManagedDatabaseServiceType) string {
	switch serviceType {
	case upcloud.ManagedDatabaseServiceTypeMySQL:
		return "MySQL"
	case upcloud.ManagedDatabaseServiceTypeOpenSearch:
		return "OpenSearch"
	case upcloud.ManagedDatabaseServiceTypePostgreSQL:
		return "PostgreSQL"
	default:
		return string(serviceType)
	}
}
