package network

import (
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/completion"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/output"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/resolver"

	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud/request"
)

// AttachCommand creates the 'objectstorage network attach' command
func AttachCommand() commands.Command {
	return &attachCommand{
		BaseCommand: commands.New(
			"attach",
			"Attach a network to a managed object storage service",
			"upctl object-storage network attach <service-uuid> --type public --name my-public-net --family IPv4",
			"upctl object-storage network attach <service-uuid> --type private --name my-private-net --family IPv4 --uuid 03fc6b80-9039-4bb7-ae43-5ccbe0ae35ce",
		),
	}
}

type attachCommand struct {
	*commands.BaseCommand
	completion.ObjectStorage
	resolver.CachingObjectStorage
	params        request.CreateManagedObjectStorageNetworkRequest
	networkUUID   string
	networkType   string
	networkName   string
	networkFamily string
}

// InitCommand implements Command.InitCommand
func (s *attachCommand) InitCommand() {
	fs := s.Cobra().Flags()
	fs.StringVar(&s.networkUUID, "network", "", "The UUID of the network to attach.")
	fs.StringVar(&s.networkType, "type", "", "The type of network (public or private).")
	fs.StringVar(&s.networkName, "name", "", "The name for the network.")
	fs.StringVar(&s.networkFamily, "family", "", "The IP family for the network (IPv4 or IPv6).")
	commands.Must(s.Cobra().MarkFlagRequired("network"))
	commands.Must(s.Cobra().MarkFlagRequired("type"))
	commands.Must(s.Cobra().MarkFlagRequired("name"))
	commands.Must(s.Cobra().MarkFlagRequired("family"))
}

// Execute implements Command.Execute
func (s *attachCommand) Execute(exec commands.Executor, serviceUUID string) (output.Output, error) {
	svc := exec.All()

	msg := "Attaching network " + s.networkUUID + " to service " + serviceUUID
	exec.PushProgressStarted(msg)

	// Use the dedicated network creation endpoint
	s.params.ServiceUUID = serviceUUID
	s.params.UUID = s.networkUUID
	s.params.Type = s.networkType
	s.params.Name = s.networkName
	s.params.Family = s.networkFamily

	res, err := svc.CreateManagedObjectStorageNetwork(exec.Context(), &s.params)
	if err != nil {
		return commands.HandleError(exec, msg, err)
	}

	exec.PushProgressSuccess(msg)

	return output.MarshaledWithHumanDetails{Value: res, Details: []output.DetailRow{
		{Title: "Network UUID", Value: s.networkUUID},
		{Title: "Service UUID", Value: serviceUUID},
	}}, nil
}
