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

type loadCommand struct {
	*commands.BaseCommand
	service service.Storage
	params  loadParams
	flagSet *pflag.FlagSet
}

type loadParams struct {
	request.LoadCDROMRequest
}

func newLoadParams() loadParams {
	return loadParams{LoadCDROMRequest: request.LoadCDROMRequest{}}
}

func LoadCommand(service service.Storage) commands.Command {
	return &loadCommand{
		BaseCommand: commands.New("load", "Load a CD-ROM"),
		service:     service,
	}
}

var DefaultLoadParams = &loadParams{
	LoadCDROMRequest: request.LoadCDROMRequest{},
}

func (s *loadCommand) InitCommand() {
	s.params = newLoadParams()

	s.flagSet = &pflag.FlagSet{}
	s.flagSet.StringVar(&s.params.StorageUUID, "storage", DefaultLoadParams.StorageUUID, "The UUID of the storage to be loaded in the CD-ROM device.")

	s.AddFlags(s.flagSet)
}

func (s *loadCommand) MakeExecuteCommand() func(args []string) (interface{}, error) {
	return func(args []string) (interface{}, error) {

		if len(args) != 1 {
			return nil, fmt.Errorf("server uuid is required")
		}

		s.params.ServerUUID = args[0]
		var loadStorageRequests []request.LoadCDROMRequest
		loadStorageRequests = append(loadStorageRequests, s.params.LoadCDROMRequest)

		var (
			mu            sync.Mutex
			numOk         int
			serverDetails []*upcloud.ServerDetails
		)
		handler := func(idx int, e *ui.LogEntry) {
			req := loadStorageRequests[idx]
			msg := fmt.Sprintf("Loading %q as a CD-ROM of server %q", req.StorageUUID, req.ServerUUID)
			e.SetMessage(msg)
			e.Start()
			details, err := s.service.LoadCDROM(&req)
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
			NumTasks:           len(loadStorageRequests),
			MaxConcurrentTasks: 5,
			EnableUI:           s.Config().InteractiveUI(),
		}, handler)
		if numOk != len(loadStorageRequests) {
			return nil, fmt.Errorf("number of CD-ROM loadings that failed: %d", len(loadStorageRequests)-numOk)
		}

		return serverDetails, nil
	}
}
