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

type templatizeCommand struct {
	*commands.BaseCommand
	service service.Storage
	params  templatizeParams
	flagSet *pflag.FlagSet
}

type templatizeParams struct {
	request.TemplatizeStorageRequest
}

func newTemplatizeParams() templatizeParams {
	return templatizeParams{TemplatizeStorageRequest: request.TemplatizeStorageRequest{}}
}

func TemplatizeCommand(service service.Storage) commands.Command {
	return &templatizeCommand{
		BaseCommand: commands.New("templatize", "Templatize a storage"),
		service:     service,
	}
}

var DefaultTemplatizeParams = &templatizeParams{
	TemplatizeStorageRequest: request.TemplatizeStorageRequest{},
}

func (s *templatizeCommand) InitCommand() {
	s.params = newTemplatizeParams()

	s.flagSet = &pflag.FlagSet{}
	s.flagSet.StringVar(&s.params.Title, "title", DefaultTemplatizeParams.Title, "A short, informational description.")

	s.AddFlags(s.flagSet)
}

func (s *templatizeCommand) MakeExecuteCommand() func(args []string) (interface{}, error) {
	return func(args []string) (interface{}, error) {

		if len(args) < 1 {
			return nil, fmt.Errorf("storage uuid is required")
		}

		var templatizeStorageRequests []request.TemplatizeStorageRequest
		for _, v := range args {
			s.params.TemplatizeStorageRequest.UUID = v
			templatizeStorageRequests = append(templatizeStorageRequests, s.params.TemplatizeStorageRequest)
		}

		var (
			mu             sync.Mutex
			numOk          int
			storageDetails []*upcloud.StorageDetails
		)
		handler := func(idx int, e *ui.LogEntry) {
			req := templatizeStorageRequests[idx]
			msg := fmt.Sprintf("Templatizing storage %q", req.UUID)
			e.SetMessage(msg)
			e.Start()
			details, err := s.service.TemplatizeStorage(&req)
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
			NumTasks:           len(templatizeStorageRequests),
			MaxConcurrentTasks: 5,
			EnableUI:           s.Config().InteractiveUI(),
		}, handler)
		if numOk != len(templatizeStorageRequests) {
			return nil, fmt.Errorf("number of template creations that failed: %d", len(templatizeStorageRequests)-numOk)
		}

		return storageDetails, nil
	}
}
