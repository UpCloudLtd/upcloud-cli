package network

import (
	"fmt"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/output"

	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud/request"
)

// DetachCommand creates the 'objectstorage network detach' command
func DetachCommand() commands.Command {
	return &detachCommand{
		BaseCommand: commands.New(
			"detach",
			"Detach a network from a managed object storage service",
			"upctl object-storage network detach <service-uuid> --name my-network",
		),
	}
}

type detachCommand struct {
	*commands.BaseCommand
	params      request.DeleteManagedObjectStorageNetworkRequest
	networkName string
}

// InitCommand implements Command.InitCommand
func (s *detachCommand) InitCommand() {
	fs := s.Cobra().Flags()
	fs.StringVar(&s.networkName, "name", "", "The name of the network to detach.")
	commands.Must(s.Cobra().MarkFlagRequired("name"))
}

// Execute implements Command.Execute
func (s *detachCommand) Execute(exec commands.Executor, uuid string) (output.Output, error) {
	if s.networkName == "" {
		return nil, fmt.Errorf("network name is required")
	}

	svc := exec.All()

	msg := fmt.Sprintf("Detaching network %s from service %s", s.networkName, uuid)
	exec.PushProgressStarted(msg)

	// Use the dedicated network deletion endpoint
	s.params.ServiceUUID = uuid
	s.params.NetworkName = s.networkName

	err := svc.DeleteManagedObjectStorageNetwork(exec.Context(), &s.params)
	if err != nil {
		return commands.HandleError(exec, msg, err)
	}

	exec.PushProgressSuccess(msg)

	return output.OnlyMarshaled{Value: fmt.Sprintf("Network %s detached from service %s", s.networkName, uuid)}, nil
}
