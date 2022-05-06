package serverstorage

import (
	"fmt"

	"github.com/UpCloudLtd/upcloud-cli/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/internal/completion"
	"github.com/UpCloudLtd/upcloud-cli/internal/output"
	"github.com/UpCloudLtd/upcloud-cli/internal/resolver"

	"github.com/UpCloudLtd/upcloud-go-api/v4/upcloud/request"
	"github.com/spf13/pflag"
)

type detachCommand struct {
	*commands.BaseCommand
	resolver.CachingServer
	completion.Server
	params detachParams
}

type detachParams struct {
	request.DetachStorageRequest
}

var defaultDetachParams = &detachParams{
	DetachStorageRequest: request.DetachStorageRequest{},
}

// DetachCommand creates the "server storage detach" command
func DetachCommand() commands.Command {
	return &detachCommand{
		BaseCommand: commands.New(
			"detach",
			"Detaches a storage resource from a server",
			"upctl server storage detach 00038afc-d526-4148-af0e-d2f1eeaded9b --address virtio:1",
			"upctl server storage detach my_server1 --address virtio:2",
		),
	}
}

// InitCommand implements Command.InitCommand
func (s *detachCommand) InitCommand() {
	s.params = detachParams{DetachStorageRequest: request.DetachStorageRequest{}}

	flagSet := &pflag.FlagSet{}
	flagSet.StringVar(&s.params.Address, "address", defaultDetachParams.Address, "Detach the storage attached to this address.")

	s.AddFlags(flagSet)
}

// MaximumExecutions implements command.Command
func (s *detachCommand) MaximumExecutions() int {
	return maxServerStorageActions
}

// ExecuteSingleArgument implements command.SingleArgumentCommand
func (s *detachCommand) ExecuteSingleArgument(exec commands.Executor, uuid string) (output.Output, error) {
	storageSvc := exec.Storage()

	if s.params.Address == "" {
		return nil, fmt.Errorf("address is required")
	}

	req := s.params.DetachStorageRequest
	req.ServerUUID = uuid

	msg := fmt.Sprintf("Detaching storage address %q from server %q", req.Address, req.ServerUUID)
	logline := exec.NewLogEntry(msg)
	logline.StartedNow()

	res, err := storageSvc.DetachStorage(&req)

	if err != nil {
		return commands.HandleError(logline, fmt.Sprintf("%s: failed", msg), err)
	}

	logline.SetMessage(fmt.Sprintf("%s: success", msg))
	logline.MarkDone()

	return output.OnlyMarshaled{Value: res}, nil
}
