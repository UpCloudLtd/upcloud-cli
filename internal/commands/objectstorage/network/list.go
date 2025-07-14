package network

import (
	"fmt"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/completion"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/output"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/resolver"
	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud/request"
)

// ListCommand creates the "network list" command
func ListCommand() commands.Command {
	return &listNetworksCommand{
		BaseCommand: commands.New(
			"list",
			"List networks in a managed object storage service",
			"upctl object-storage network list <service-uuid>",
			"upctl object-storage network list my-service",
		),
	}
}

type listNetworksCommand struct {
	*commands.BaseCommand
	completion.ObjectStorage
	resolver.CachingObjectStorage
}

// InitCommand implements Command.InitCommand
func (s *listNetworksCommand) InitCommand() {
	s.Cobra().Long = commands.WrapLongDescription(`List networks in a managed object storage service

Lists all networks attached to the specified managed object storage service, showing their names, UUIDs, types, and families.`)
}

// MaximumExecutions implements commands.MultipleArgumentCommand
func (s *listNetworksCommand) MaximumExecutions() int {
	return 1
}

// Execute implements commands.MultipleArgumentCommand
func (s *listNetworksCommand) Execute(exec commands.Executor, serviceUUID string) (output.Output, error) {
	svc := exec.All()

	msg := fmt.Sprintf("Listing networks in service %s", serviceUUID)
	exec.PushProgressStarted(msg)

	// Get service details to access networks
	req := &request.GetManagedObjectStorageRequest{
		UUID: serviceUUID,
	}

	exec.PushProgressUpdateMessage(msg, msg)
	res, err := svc.GetManagedObjectStorage(exec.Context(), req)
	if err != nil {
		return commands.HandleError(exec, msg, err)
	}

	exec.PushProgressSuccess(msg)

	rows := []output.TableRow{}
	for _, network := range res.Networks {
		uuid := "unknown"
		if network.UUID != nil {
			uuid = *network.UUID
		}
		rows = append(rows, output.TableRow{
			network.Name,
			uuid,
			network.Type,
			network.Family,
		})
	}

	return output.MarshaledWithHumanOutput{
		Value: res.Networks,
		Output: output.Table{
			Columns: []output.TableColumn{
				{Key: "name", Header: "Name"},
				{Key: "uuid", Header: "UUID"},
				{Key: "type", Header: "Type"},
				{Key: "family", Header: "Family"},
			},
			Rows: rows,
		},
	}, nil
}
