package databaseindex

import (
	"fmt"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/completion"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/format"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/output"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/resolver"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/ui"

	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud/request"
)

// ListCommand creates the "database index list" command
func ListCommand() commands.Command {
	return &listCommand{
		BaseCommand: commands.New("list", "List current indices of the specified databases", "upctl database index list 55199a44-4751-4e27-9394-7c7661910be3"),
	}
}

type listCommand struct {
	*commands.BaseCommand
	resolver.CachingDatabase
	completion.Database
}

// Execute implements commands.MultipleArgumentCommand
func (s *listCommand) Execute(exec commands.Executor, uuid string) (output.Output, error) {
	svc := exec.All()
	indices, err := svc.GetManagedDatabaseIndices(exec.Context(), &request.GetManagedDatabaseIndicesRequest{ServiceUUID: uuid})
	if err != nil {
		return nil, err
	}

	rows := []output.TableRow{}
	for _, index := range indices {
		rows = append(rows, output.TableRow{
			index.IndexName,
			index.CreateTime,
			index.Health,
			index.Status,
			index.Docs,
			index.NumberOfReplicas,
			index.NumberOfShards,
			index.ReadOnlyAllowDelete,
			fmt.Sprintf("%d bytes", index.Size),
		})
	}

	return output.MarshaledWithHumanOutput{
		Value: indices,
		Output: output.Table{
			Columns: []output.TableColumn{
				{Key: "index_name", Header: "Name", Colour: ui.DefaultUUUIDColours},
				{Key: "create_time", Header: "Created"},
				{Key: "health", Header: "Health", Format: format.DatabaseIndexHealth},
				{Key: "status", Header: "Status", Format: format.DatabaseIndexState},
				{Key: "docs", Header: "Documents"},
				{Key: "number_of_replicas", Header: "Replicas"},
				{Key: "number_of_shards", Header: "Shards"},
				{Key: "read_only_allow_delete", Header: "Read-only & allow delete", Format: format.Boolean},
				{Key: "size", Header: "Size"},
			},
			Rows: rows,
		},
	}, nil
}
