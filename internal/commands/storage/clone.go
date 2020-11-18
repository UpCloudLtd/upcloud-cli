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

type cloneCommand struct {
	*commands.BaseCommand
	service service.Storage
	params  cloneParams
	flagSet *pflag.FlagSet
}

type cloneParams struct {
	request.CloneStorageRequest
}

func newCloneParams() cloneParams {
	return cloneParams{CloneStorageRequest: request.CloneStorageRequest{}}
}

func CloneCommand(service service.Storage) commands.Command {
	return &cloneCommand{
		BaseCommand: commands.New("clone", "Clone a storage"),
		service:     service,
	}
}

var DefaultCloneParams = &cloneParams{
	CloneStorageRequest: request.CloneStorageRequest{
		Tier: "hdd",
	},
}

func (s *cloneCommand) InitCommand() {
	s.params = newCloneParams()

	s.flagSet = &pflag.FlagSet{}
	s.flagSet.StringVar(&s.params.Tier, "tier", DefaultCloneParams.Tier, "The storage tier to use.")
	s.flagSet.StringVar(&s.params.Title, "title", DefaultCloneParams.Title, "A short, informational description.")
	s.flagSet.StringVar(&s.params.Zone, "zone", DefaultCloneParams.Zone, "The zone in which the storage will be created, e.g. fi-hel1.")

	s.AddFlags(s.flagSet)
}

func (s *cloneCommand) MakeExecuteCommand() func(args []string) (interface{}, error) {
	return func(args []string) (interface{}, error) {

		if len(args) < 1 {
			return nil, fmt.Errorf("storage uuid is required")
		}

		var cloneStorageRequests []request.CloneStorageRequest
		for _, v := range args {
			s.params.CloneStorageRequest.UUID = v
			cloneStorageRequests = append(cloneStorageRequests, s.params.CloneStorageRequest)
		}

		var (
			mu             sync.Mutex
			numOk          int
			storageDetails []*upcloud.StorageDetails
		)
		handler := func(idx int, e *ui.LogEntry) {
			req := cloneStorageRequests[idx]
			msg := fmt.Sprintf("Cloneing storage %q", req.UUID)
			e.SetMessage(msg)
			e.Start()
			details, err := s.service.CloneStorage(&req)
			if err != nil {
				e.SetMessage(ui.LiveLogEntryErrorColours.Sprintf("%s: failed", msg))
				e.SetDetails(err.Error(), "error: ")
			} else {
				e.SetMessage(fmt.Sprintf("%s: done", msg))
				e.SetDetails(details.UUID, "UUID: ")
				mu.Lock()
				numOk++
				storageDetails = append(storageDetails, details)
				mu.Unlock()
			}
		}
		ui.StartWorkQueue(ui.WorkQueueConfig{
			NumTasks:           len(cloneStorageRequests),
			MaxConcurrentTasks: 5,
			EnableUI:           s.Config().InteractiveUI(),
		}, handler)
		if numOk != len(cloneStorageRequests) {
			return nil, fmt.Errorf("number of clone operations that failed: %d", len(cloneStorageRequests)-numOk)
		}

		return storageDetails, nil
	}
}
