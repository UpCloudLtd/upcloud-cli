package server

import (
	"fmt"
	"github.com/UpCloudLtd/cli/internal/commands"
	"github.com/UpCloudLtd/cli/internal/ui"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/request"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/service"
	"github.com/spf13/pflag"
	"sync"
)

func ModifyCommand(service service.Server) commands.Command {
	return &modifyCommand{
		BaseCommand: commands.New("modify", "Modifies the configuration of an existing server"),
		service:     service,
	}
}

type modifyCommand struct {
	*commands.BaseCommand
	service service.Server
	params  modifyParams
	flagSet *pflag.FlagSet
}

type modifyParams struct {
	request.ModifyServerRequest
	remoteAccessEnabled string
	metadata            string
}

var defModify = modifyParams{
	ModifyServerRequest: request.ModifyServerRequest{},
}

func (s *modifyCommand) InitCommand() {
	s.params = modifyParams{ModifyServerRequest: request.ModifyServerRequest{}}
	flags := &pflag.FlagSet{}
	flags.StringVar(&s.params.BootOrder, "boot-order", defModify.BootOrder, "The boot device order.")
	flags.IntVar(&s.params.CoreNumber, "cores", defModify.CoreNumber, "Number of cores")
	flags.StringVar(&s.params.Hostname, "hostname", defModify.Hostname, "Hostname")
	flags.StringVar(&s.params.Firewall, "firewall", defModify.Firewall, "Sets the firewall on or off. You can manage firewall rules with the firewall command\nAvailable: on, off")
	flags.IntVar(&s.params.MemoryAmount, "memory", defModify.MemoryAmount, "Memory amount in MiB")
	flags.StringVar(&s.params.metadata, "metadata", defModify.metadata, "Enable metadata service")
	flags.StringVar(&s.params.Plan, "plan", defModify.Plan, "Server plan to use. Set this to custom to use custom core/memory amounts.")
	flags.StringVar(&s.params.SimpleBackup, "simple-backup", defModify.SimpleBackup, "Simple backup rule. Format (HHMM,{dailies,weeklies,monthlies}).\nExample: 2300,dailies")
	flags.StringVar(&s.params.Title, "title", defModify.Title, "Visible name")
	flags.StringVar(&s.params.TimeZone, "time-zone", defModify.TimeZone, "Time zone to set the RTC to")
	flags.StringVar(&s.params.VideoModel, "video-model", defModify.VideoModel, "Video interface model of the server.\nAvailable: vga,cirrus")
	flags.StringVar(&s.params.remoteAccessEnabled, "remote-access-enabled", defModify.remoteAccessEnabled, "Enables or disables the remote access\nAvailable: true, false")
	flags.StringVar(&s.params.RemoteAccessType, "remote-access-type", defModify.RemoteAccessType, "The remote access type")
	flags.StringVar(&s.params.RemoteAccessPassword, "remote-access-password", defModify.RemoteAccessPassword, "The remote access password")

	s.AddFlags(flags)
}

func (s *modifyCommand) MakeExecuteCommand() func(args []string) (interface{}, error) {
	return func(args []string) (interface{}, error) {

		if len(args) < 1 {
			return nil, fmt.Errorf("server uuid is required")
		}

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

		var modifyServerRequests []request.ModifyServerRequest
		for _, v := range args {
			s.params.UUID = v
			modifyServerRequests = append(modifyServerRequests, s.params.ModifyServerRequest)
		}

		var (
			mu            sync.Mutex
			numOk         int
			serverDetails []*upcloud.ServerDetails
		)
		handler := func(idx int, e *ui.LogEntry) {
			storageRequest := modifyServerRequests[idx]
			msg := fmt.Sprintf("Modifying server %q", storageRequest.UUID)
			e.SetMessage(msg)
			e.Start()
			details, err := s.service.ModifyServer(&storageRequest)
			if err != nil {
				e.SetMessage(ui.LiveLogEntryErrorColours.Sprintf("%s: failed", msg))
				e.SetDetails(err.Error(), "error: ")
			} else {
				e.SetMessage(fmt.Sprintf("%s: done", msg))
				e.SetDetails(details.UUID, "UUID: ")
				mu.Lock()
				numOk++
				serverDetails = append(serverDetails, details)
				mu.Unlock()
			}
		}
		ui.StartWorkQueue(ui.WorkQueueConfig{
			NumTasks:           len(modifyServerRequests),
			MaxConcurrentTasks: 5,
			EnableUI:           s.Config().InteractiveUI(),
		}, handler)
		if numOk != len(modifyServerRequests) {
			return nil, fmt.Errorf("number of servers that failed: %d", len(modifyServerRequests)-numOk)
		}

		return serverDetails, nil

	}
}
