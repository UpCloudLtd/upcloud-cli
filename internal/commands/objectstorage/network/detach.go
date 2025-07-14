package network

import (
	"fmt"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/completion"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/output"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/resolver"
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
	completion.ObjectStorage
	resolver.CachingObjectStorage
	params request.DeleteManagedObjectStorageNetworkRequest
}

// InitCommand implements Command.InitCommand
func (s *detachCommand) InitCommand() {
	fs := s.Cobra().Flags()
	fs.StringVar(&s.params.NetworkName, "name", "", "The name of the network to detach.")
	commands.Must(s.Cobra().MarkFlagRequired("name"))
}

// Execute implements Command.Execute
func (s *detachCommand) Execute(exec commands.Executor, serviceUUID string) (output.Output, error) {
	svc := exec.All()

	msg := fmt.Sprintf("Detaching network %s from service %s", s.params.NetworkName, serviceUUID)
	exec.PushProgressStarted(msg)

	// Use the dedicated network deletion endpoint
	s.params.ServiceUUID = serviceUUID

	err := svc.DeleteManagedObjectStorageNetwork(exec.Context(), &s.params)
	if err != nil {
		return commands.HandleError(exec, msg, err)
	}

	exec.PushProgressSuccess(msg)

	return output.OnlyMarshaled{Value: fmt.Sprintf("Network %s detached from service %s", s.params.NetworkName, serviceUUID)}, nil
}
