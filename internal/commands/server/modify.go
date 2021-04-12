package server

import (
	"fmt"

	"github.com/UpCloudLtd/cli/internal/commands"
	"github.com/UpCloudLtd/cli/internal/completion"
	"github.com/UpCloudLtd/cli/internal/output"
	"github.com/UpCloudLtd/cli/internal/resolver"
	"github.com/UpCloudLtd/cli/internal/ui"

	"github.com/UpCloudLtd/upcloud-go-api/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/request"
	"github.com/spf13/pflag"
)

// ModifyCommand creates the "server modify" command
func ModifyCommand() commands.Command {
	return &modifyCommand{
		BaseCommand: commands.New("modify", "Modifies the configuration of an existing server"),
	}
}

type modifyCommand struct {
	*commands.BaseCommand
	params modifyParams
	resolver.CachingServer
	completion.Server
}

type modifyParams struct {
	request.ModifyServerRequest
	remoteAccessEnabled string
	metadata            string
}

var defaultModifyParams = modifyParams{
	ModifyServerRequest: request.ModifyServerRequest{},
}

// InitCommand implements Command.InitCommand
func (s *modifyCommand) InitCommand() {
	s.params = modifyParams{ModifyServerRequest: request.ModifyServerRequest{}}
	flags := &pflag.FlagSet{}
	flags.StringVar(&s.params.BootOrder, "boot-order", defaultModifyParams.BootOrder, "The boot device order.")
	flags.IntVar(&s.params.CoreNumber, "cores", defaultModifyParams.CoreNumber, "Number of cores. Sets server plan to custom.")
	flags.StringVar(&s.params.Hostname, "hostname", defaultModifyParams.Hostname, "Hostname.")
	flags.StringVar(&s.params.Firewall, "firewall", defaultModifyParams.Firewall, "Enables or disables firewall on the server. You can manage firewall rules with the firewall command.\nAvailable: true, false")
	flags.IntVar(&s.params.MemoryAmount, "memory", defaultModifyParams.MemoryAmount, "Memory amount in MiB. Sets server plan to custom.")
	flags.StringVar(&s.params.metadata, "metadata", defaultModifyParams.metadata, "Enable metadata service.")
	flags.StringVar(&s.params.Plan, "plan", defaultModifyParams.Plan, "Server plan to use.")
	flags.StringVar(&s.params.SimpleBackup, "simple-backup", defaultModifyParams.SimpleBackup, "Simple backup rule. Format (HHMM,{dailies,weeklies,monthlies}).\nExample: 2300,dailies")
	flags.StringVar(&s.params.Title, "title", defaultModifyParams.Title, "Visible name.")
	flags.StringVar(&s.params.TimeZone, "time-zone", defaultModifyParams.TimeZone, "Time zone to set the RTC to.")
	flags.StringVar(&s.params.VideoModel, "video-model", defaultModifyParams.VideoModel, "Video interface model of the server.\nAvailable: vga,cirrus")
	flags.StringVar(&s.params.remoteAccessEnabled, "remote-access-enabled", defaultModifyParams.remoteAccessEnabled, "Enables or disables the remote access.\nAvailable: true, false")
	flags.StringVar(&s.params.RemoteAccessType, "remote-access-type", defaultModifyParams.RemoteAccessType, "The remote access type.")
	flags.StringVar(&s.params.RemoteAccessPassword, "remote-access-password", defaultModifyParams.RemoteAccessPassword, "The remote access password.")

	s.AddFlags(flags)
}

// Execute implements command.Command
func (s *modifyCommand) Execute(exec commands.Executor, uuid string) (output.Output, error) {

	remoteAccess := new(upcloud.Boolean)
	if err := remoteAccess.UnmarshalJSON([]byte(s.params.remoteAccessEnabled)); err != nil {
		return nil, err
	}
	s.params.RemoteAccessEnabled = *remoteAccess

	metadata := new(upcloud.Boolean)
	if err := metadata.UnmarshalJSON([]byte(s.params.metadata)); err != nil {
		return nil, err
	}
	s.params.Metadata = *metadata

	svc := exec.Server()

	//XXX: should fix the SDK with the correct type
	switch s.params.Firewall {
	case "true":
		s.params.Firewall = "on"
	case "false":
		s.params.Firewall = "off"
	}

	if s.params.CoreNumber != 0 || s.params.MemoryAmount != 0 {
		s.params.Plan = "custom" // Valid for all custom plans.
	}

	msg := fmt.Sprintf("modifing server %v", uuid)
	logline := exec.NewLogEntry(msg)

	logline.StartedNow()

	req := s.params.ModifyServerRequest
	req.UUID = uuid
	res, err := svc.ModifyServer(&req)
	if err != nil {
		logline.SetMessage(ui.LiveLogEntryErrorColours.Sprintf("%s: failed (%v)", msg, err.Error()))
		logline.SetDetails(err.Error(), "error: ")
		return nil, err
	}

	logline.SetMessage(fmt.Sprintf("%s: done", msg))
	logline.MarkDone()

	return output.Marshaled{Value: res}, nil
}
