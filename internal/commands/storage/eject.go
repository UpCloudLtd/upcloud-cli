package storage

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

type ejectCommand struct {
	*commands.BaseCommand
	service service.Storage
	params  ejectParams
	flagSet *pflag.FlagSet
}

type ejectParams struct {
	request.EjectCDROMRequest
}

func EjectCommand(service service.Storage) commands.Command {
	return &ejectCommand{
		BaseCommand: commands.New("eject", "Eject a CD-ROM"),
		service:     service,
	}
}

func (s *ejectCommand) MakeExecuteCommand() func(args []string) (interface{}, error) {
	return func(args []string) (interface{}, error) {

		if len(args) < 1 {
			return nil, fmt.Errorf("server uuid is required")
		}

		var ejectStorageRequests []request.EjectCDROMRequest
		for _, v := range args {
			s.params.EjectCDROMRequest.ServerUUID = v
			ejectStorageRequests = append(ejectStorageRequests, s.params.EjectCDROMRequest)
		}

		var (
			mu            sync.Mutex
			numOk         int
			serverDetails []*upcloud.ServerDetails
		)
		handler := func(idx int, e *ui.LogEntry) {
			req := ejectStorageRequests[idx]
			msg := fmt.Sprintf("Ejecting CD-ROM of server %q", req.ServerUUID)
			e.SetMessage(msg)
			e.Start()
			details, err := s.service.EjectCDROM(&req)
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
			NumTasks:           len(ejectStorageRequests),
			MaxConcurrentTasks: 5,
			EnableUI:           s.Config().InteractiveUI(),
		}, handler)
		if numOk != len(ejectStorageRequests) {
			return nil, fmt.Errorf("number of CD-ROM ejections that failed: %d", len(ejectStorageRequests)-numOk)
		}

		return serverDetails, nil
	}
}
