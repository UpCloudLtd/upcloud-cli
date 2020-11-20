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

type detachCommand struct {
	*commands.BaseCommand
	service service.Storage
	params  detachParams
	flagSet *pflag.FlagSet
}

type detachParams struct {
	request.DetachStorageRequest
}

func newDetachParams() detachParams {
	return detachParams{DetachStorageRequest: request.DetachStorageRequest{}}
}

func DetachCommand(service service.Storage) commands.Command {
	return &detachCommand{
		BaseCommand: commands.New("detach", "Detaches a storage resource from a server"),
		service:     service,
	}
}

var DefaultDetachParams = &detachParams{
	DetachStorageRequest: request.DetachStorageRequest{},
}

func (s *detachCommand) InitCommand() {
	s.params = newDetachParams()

	s.flagSet = &pflag.FlagSet{}
	s.flagSet.StringVar(&s.params.Address, "address", DefaultDetachParams.Address, "Detach the storage attached to this address.")

	s.AddFlags(s.flagSet)
}

func (s *detachCommand) MakeExecuteCommand() func(args []string) (interface{}, error) {
	return func(args []string) (interface{}, error) {

		if len(args) < 1 {
			return nil, fmt.Errorf("server uuid is required")
		}

		var detachStorageRequests []request.DetachStorageRequest
		for _, v := range args {
			s.params.DetachStorageRequest.ServerUUID = v
			detachStorageRequests = append(detachStorageRequests, s.params.DetachStorageRequest)
		}

		var (
			mu            sync.Mutex
			numOk         int
			serverDetails []*upcloud.ServerDetails
		)
		handler := func(idx int, e *ui.LogEntry) {
			req := detachStorageRequests[idx]
			msg := fmt.Sprintf("Detaching address %q to server %q", req.Address, req.ServerUUID)
			e.SetMessage(msg)
			e.Start()
			details, err := s.service.DetachStorage(&req)
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
			NumTasks:           len(detachStorageRequests),
			MaxConcurrentTasks: 5,
			EnableUI:           s.Config().InteractiveUI(),
		}, handler)
		if numOk != len(detachStorageRequests) {
			return nil, fmt.Errorf("number of operations that failed: %d", len(detachStorageRequests)-numOk)
		}

		return serverDetails, nil
	}
}
