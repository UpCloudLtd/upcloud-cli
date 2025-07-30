package objectstorage

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/UpCloudLtd/progress/messages"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/config"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/labels"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/output"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/ui"

	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud/request"
)

// CreateCommand creates the 'objectstorage create' command
func CreateCommand() commands.Command {
	return &createCommand{
		BaseCommand: commands.New(
			"create",
			"Create a new managed object storage service",
			"upctl object-storage create --name my-service --region europe-1",
			"upctl object-storage create --name my-service --region europe-1 --configured-status started",
		),
	}
}

type createCommand struct {
	*commands.BaseCommand
	params           request.CreateManagedObjectStorageRequest
	wait             config.OptionalBoolean
	labels           []string
	networks         []string
	configuredStatus string
}

// InitCommand implements Command.InitCommand
func (s *createCommand) InitCommand() {
	s.Cobra().Long = commands.WrapLongDescription(`Create a new managed object storage service

Creates a new managed object storage service in the specified region. The service can be started or stopped based on the configured status.`)

	fs := s.Cobra().Flags()

	fs.StringVar(&s.params.Name, "name", "", "The name of the service. Must be unique within customer account.")
	fs.StringVar(&s.params.Region, "region", "", "Region in which the service will be hosted.")
	fs.StringVar(&s.configuredStatus, "configured-status", "started", "Service status managed by the customer. Valid values: started, stopped")
	fs.StringArrayVar(&s.labels, "label", nil, "Labels to describe the service in `key=value` format, multiple can be declared.\nUsage: --label env=dev\n\n--label owner=operations")
	fs.StringArrayVar(&s.networks, "network", nil, "Networks for the service. At least one network is needed. If not specified, a public network will be used by default.\nNote: Only one public network is allowed per service, but multiple private networks are supported.\nPublic network: --network type=public,name=my-public-network,family=IPv4\nPrivate network: --network type=private,name=my-private-network,family=IPv4,uuid=03fc6b80-9039-4bb7-ae43-5ccbe0ae35ce")
	config.AddToggleFlag(fs, &s.wait, "wait", false, "Wait for service to be in running state before returning.")

	commands.Must(s.Cobra().MarkFlagRequired("name"))
	commands.Must(s.Cobra().MarkFlagRequired("region"))
}

// ExecuteWithoutArguments implements commands.NoArgumentCommand
func (s *createCommand) ExecuteWithoutArguments(exec commands.Executor) (output.Output, error) {
	if s.configuredStatus != "started" && s.configuredStatus != "stopped" {
		return nil, fmt.Errorf("configured-status must be either 'started' or 'stopped'")
	}
	s.params.ConfiguredStatus = upcloud.ManagedObjectStorageConfiguredStatus(s.configuredStatus)

	svc := exec.All()
	msg := fmt.Sprintf("Creating object storage service %v", s.params.Name)
	exec.PushProgressStarted(msg)

	// Process networks
	networks, err := s.processNetworks()
	if err != nil {
		return nil, err
	}

	// Validate network constraints
	err = s.validateCreateNetworks(networks)
	if err != nil {
		return nil, err
	}

	// If no networks provided, add a default public network
	if len(networks) == 0 {
		networks = []upcloud.ManagedObjectStorageNetwork{
			{
				Family: "IPv4",
				Name:   fmt.Sprintf("%s-public-network", s.params.Name),
				Type:   "public",
			},
		}
	}

	s.params.Networks = networks

	// Process labels
	labelSlice, err := labels.StringsToSliceOfLabels(s.labels)
	if err != nil {
		return nil, err
	}
	s.params.Labels = labelSlice

	res, err := svc.CreateManagedObjectStorage(exec.Context(), &s.params)
	if err != nil {
		return commands.HandleError(exec, msg, err)
	}

	if s.wait.Value() {
		waitForObjectStorageServiceState(res.UUID, upcloud.ManagedObjectStorageOperationalStateRunning, exec, msg)
	} else {
		exec.PushProgressSuccess(msg)
	}

	return output.MarshaledWithHumanDetails{Value: res, Details: []output.DetailRow{
		{Title: "UUID", Value: res.UUID, Colour: ui.DefaultUUUIDColours},
		{Title: "Name", Value: res.Name},
		{Title: "Region", Value: res.Region},
	}}, nil
}

func (s *createCommand) processNetworks() ([]upcloud.ManagedObjectStorageNetwork, error) {
	var networks []upcloud.ManagedObjectStorageNetwork
	for _, n := range s.networks {
		var network upcloud.ManagedObjectStorageNetwork
		parts := strings.Split(n, ",")
		for _, part := range parts {
			kv := strings.SplitN(part, "=", 2)
			if len(kv) != 2 {
				return nil, fmt.Errorf("invalid network format: %s", part)
			}
			key, value := kv[0], kv[1]
			switch key {
			case "type":
				network.Type = value
			case "name":
				network.Name = value
			case "family":
				network.Family = value
			case "uuid":
				network.UUID = &value
			}
		}
		networks = append(networks, network)
	}
	return networks, nil
}

func (s *createCommand) validateCreateNetworks(networks []upcloud.ManagedObjectStorageNetwork) error {
	var publicNetworkCount int
	for _, n := range networks {
		if n.Type == "public" {
			publicNetworkCount++
		}
	}
	if publicNetworkCount > 1 {
		return fmt.Errorf("only one public network is allowed per service")
	}

	return nil
}

// waitForObjectStorageServiceState waits for object storage service to reach given state and updates progress message with key matching given msg. Finally, progress message is updated back to given msg and either done state or timeout warning.
func waitForObjectStorageServiceState(uuid string, state upcloud.ManagedObjectStorageOperationalState, exec commands.Executor, msg string) {
	exec.PushProgressUpdateMessage(msg, fmt.Sprintf("Waiting for object storage service %s to be in %s state", uuid, state))

	ctx, cancel := context.WithTimeout(exec.Context(), 15*time.Minute)
	defer cancel()

	if _, err := exec.All().WaitForManagedObjectStorageOperationalState(ctx, &request.WaitForManagedObjectStorageOperationalStateRequest{
		UUID:         uuid,
		DesiredState: state,
	}); err != nil {
		exec.PushProgressUpdate(messages.Update{
			Key:     msg,
			Message: msg,
			Status:  messages.MessageStatusWarning,
			Details: "Error: " + err.Error(),
		})
		return
	}

	exec.PushProgressUpdateMessage(msg, msg)
	exec.PushProgressSuccess(msg)
}
