package server

import (
	"fmt"

	"github.com/UpCloudLtd/upcloud-cli/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/internal/completion"
	"github.com/UpCloudLtd/upcloud-cli/internal/config"
	"github.com/UpCloudLtd/upcloud-cli/internal/output"
	"github.com/UpCloudLtd/upcloud-cli/internal/resolver"

	"github.com/UpCloudLtd/upcloud-go-api/v4/upcloud/request"
	"github.com/spf13/pflag"
)

// ModifyCommand creates the "server modify" command
func ModifyCommand() commands.Command {
	return &modifyCommand{
		BaseCommand: commands.New(
			"modify",
			"Modifies the configuration of an existing server",
			"upctl server modify 00bb4617-c592-4b32-b869-35a60b323b18 --plan 1xCPU-1GB",
			"upctl server modify 00bb4617-c592-4b32-b869-35a60b323b18 0053a6f5-e6d1-4b0b-b9dc-b90d0894e8d0 --plan 1xCPU-1GB",
			"upctl server modify my_server1 --plan 1xCPU-2GB",
			"upctl server modify myapp --hostname superapp",
		),
	}
}

type modifyCommand struct {
	*commands.BaseCommand
	params       modifyParams
	setMetadata  config.OptionalBoolean
	remoteAccess config.OptionalBoolean
	firewall     config.OptionalBoolean
	resolver.CachingServer
	completion.Server
}

type modifyParams struct {
	request.ModifyServerRequest
}

var defaultModifyParams = modifyParams{
	ModifyServerRequest: request.ModifyServerRequest{},
}

// InitCommand implements Command.InitCommand
func (s *modifyCommand) InitCommand() {
	s.params = modifyParams{ModifyServerRequest: request.ModifyServerRequest{}}
	flags := &pflag.FlagSet{}
	flags.StringVar(&s.params.BootOrder, "boot-order", defaultModifyParams.BootOrder, "The boot device order, disk / cdrom / network or comma separated combination.")
	flags.IntVar(&s.params.CoreNumber, "cores", defaultModifyParams.CoreNumber, "Number of cores. Sets server plan to custom.")
	flags.StringVar(&s.params.Hostname, "hostname", defaultModifyParams.Hostname, "Hostname.")
	config.AddEnableDisableFlags(flags, &s.firewall, "firewall", "firewall")
	// flags.StringVar(&s.params.Firewall, "firewall", defaultModifyParams.Firewall, "Enables or disables firewall on the server. You can manage firewall rules with the firewall command.\nAvailable: true, false")
	flags.IntVar(&s.params.MemoryAmount, "memory", defaultModifyParams.MemoryAmount, "Memory amount in MiB. Sets server plan to custom.")
	config.AddEnableDisableFlags(flags, &s.setMetadata, "metadata", "metadata service")
	// flags.StringVar(&s.params.metadata, "metadata", defaultModifyParams.metadata, "Enable metadata service.")
	flags.StringVar(&s.params.Plan, "plan", defaultModifyParams.Plan, "Server plan to use.")
	flags.StringVar(&s.params.SimpleBackup, "simple-backup", defaultModifyParams.SimpleBackup, "Simple backup rule. Format (HHMM,{dailies,weeklies,monthlies}).\nExample: 2300,dailies")
	flags.StringVar(&s.params.Title, "title", defaultModifyParams.Title, "A short, informational description.")
	flags.StringVar(&s.params.TimeZone, "time-zone", defaultModifyParams.TimeZone, "Time zone to set the RTC to.")
	flags.StringVar(&s.params.VideoModel, "video-model", defaultModifyParams.VideoModel, "Video interface model of the server.\nAvailable: vga,cirrus")
	config.AddEnableDisableFlags(flags, &s.remoteAccess, "remote-access", "remote access")
	// flags.StringVar(&s.params.remoteAccessEnabled, "remote-access-enabled", defaultModifyParams.remoteAccessEnabled, "Enables or disables the remote access.\nAvailable: true, false")
	flags.StringVar(&s.params.RemoteAccessType, "remote-access-type", defaultModifyParams.RemoteAccessType, "The remote access type.")
	flags.StringVar(&s.params.RemoteAccessPassword, "remote-access-password", defaultModifyParams.RemoteAccessPassword, "The remote access password.")

	s.AddFlags(flags)
}

// Execute implements commands.MultipleArgumentCommand
func (s *modifyCommand) Execute(exec commands.Executor, uuid string) (output.Output, error) {
	svc := exec.Server()

	// TODO: refactor out when go-api actually supports not-set upcloud.Booleans in requests
	// ref: https://app.asana.com/0/1191419140326561/1200258914439524
	serverDetails, err := svc.GetServerDetails(&request.GetServerDetailsRequest{UUID: uuid})
	if err != nil {
		return nil, err
	}

	s.params.RemoteAccessEnabled = s.remoteAccess.OverrideNotSet(serverDetails.RemoteAccessEnabled.Bool()).AsUpcloudBoolean()
	s.params.Metadata = s.setMetadata.OverrideNotSet(serverDetails.Metadata.Bool()).AsUpcloudBoolean()

	// TODO: refactor when go-api parameter is refactored
	switch s.firewall {
	case config.True:
		s.params.Firewall = "on"
	case config.False:
		s.params.Firewall = "off"
		// nb. no handling for not set, just pass in an empty string in the request
	}

	if s.params.CoreNumber != 0 || s.params.MemoryAmount != 0 {
		s.params.Plan = "custom" // Valid for all custom plans.
	}

	msg := fmt.Sprintf("Modifying server %v", uuid)
	exec.PushProgressStarted(msg)

	req := s.params.ModifyServerRequest
	req.UUID = uuid
	res, err := svc.ModifyServer(&req)
	if err != nil {
		return commands.HandleError(exec, msg, err)
	}

	exec.PushProgressSuccess(msg)

	return output.OnlyMarshaled{Value: res}, nil
}
