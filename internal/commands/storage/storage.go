package storage

import (
	"errors"
	"fmt"
	"github.com/UpCloudLtd/cli/internal/ui"
	"github.com/spf13/cobra"
	"time"

	"github.com/UpCloudLtd/upcloud-go-api/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/request"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/service"
	"github.com/jedib0t/go-pretty/v6/text"

	"github.com/UpCloudLtd/cli/internal/commands"
)

var (
	maxStorageActions = 10
	cachedStorages    []upcloud.Storage
)

const minStorageSize = 10
const positionalArgHelp = "<UUID or Title>"

func StorageCommand() commands.Command {
	return &storageCommand{commands.New("storage", "Manage storages")}
}

type storageCommand struct {
	*commands.BaseCommand
}

func StateColour(state string) text.Colors {
	switch state {
	case upcloud.StorageStateOnline, upcloud.StorageStateSyncing:
		return text.Colors{text.FgGreen}
	case upcloud.StorageStateError:
		return text.Colors{text.FgHiRed, text.Bold}
	case upcloud.StorageStateMaintenance:
		return text.Colors{text.FgYellow}
	case upcloud.StorageStateCloning, upcloud.StorageStateBackuping:
		return text.Colors{text.FgHiMagenta, text.Bold}
	default:
		return text.Colors{text.FgHiBlack}
	}
}

func ImportStateColour(state string) text.Colors {
	switch state {
	case "completed":
		return text.Colors{text.FgGreen}
	case "failed":
		return text.Colors{text.FgHiRed, text.Bold}
	case "pending", "importing":
		return text.Colors{text.FgYellow}
	case "cancelling":
		return text.Colors{text.FgHiMagenta, text.Bold}
	default:
		return text.Colors{text.FgHiBlack}
	}
}

func matchStorages(storages []upcloud.Storage, searchVal string) []*upcloud.Storage {
	var r []*upcloud.Storage
	for _, storage := range storages {
		storage := storage
		if storage.Title == searchVal || storage.UUID == searchVal {
			r = append(r, &storage)
		}
	}
	return r
}

func searchStorage(storagesPtr *[]upcloud.Storage, service service.Storage, uuidOrTitle string, unique bool) ([]*upcloud.Storage, error) {
	if storagesPtr == nil || service == nil {
		return nil, fmt.Errorf("no storages or service passed")
	}
	storages := *storagesPtr
	if len(cachedStorages) == 0 {
		res, err := service.GetStorages(&request.GetStoragesRequest{})
		if err != nil {
			return nil, err
		}
		storages = res.Storages
		*storagesPtr = storages
	}
	matched := matchStorages(storages, uuidOrTitle)
	if len(matched) == 0 {
		return nil, fmt.Errorf("no storage with uuid, name or title %q was found", uuidOrTitle)
	}
	if len(matched) > 1 && unique {
		return nil, fmt.Errorf("multiple storages matched to query %q, use UUID to specify", uuidOrTitle)
	}
	return matched, nil
}

func SearchAllStorages(uuidOrTitle []string, service service.Storage, unique bool) ([]*upcloud.Storage, error) {
	var result []*upcloud.Storage
	for _, id := range uuidOrTitle {
		matchedResults, err := searchStorage(&cachedStorages, service, id, unique)
		if err != nil {
			return nil, err
		}
		result = append(result, matchedResults...)
	}
	return result, nil
}

func SearchSingleStorage(uuidOrTitle string, service service.Storage) (*upcloud.Storage, error) {
	matchedResults, err := searchStorage(&cachedStorages, service, uuidOrTitle, true)
	if err != nil {
		return nil, err
	}
	return matchedResults[0], nil
}

func WaitForImportState(service *service.Service, uuid, desiredState string, timeout time.Duration) (*upcloud.StorageImportDetails, error) {
	timer := time.After(timeout)
	for {
		time.Sleep(5 * time.Second)
		details, err := service.GetStorageImportDetails(&request.GetStorageImportDetailsRequest{UUID: uuid})
		if err != nil {
			return nil, err
		}
		switch details.State {
		case upcloud.StorageImportStateFailed:
			return nil, errors.New("import in failed state")
		case upcloud.StorageImportStateCancelled:
			return nil, errors.New("import in cancelled state")
		case desiredState:
			return details, nil
		}
		select {
		case <-timer:
			return nil, fmt.Errorf("timed out while waiting an import to transition into %q", desiredState)
		default:
		}
	}
}

var getStorageDetailsUuid = func(in interface{}) string { return in.(*upcloud.StorageDetails).UUID }

type Request struct {
	ExactlyOne   bool
	BuildRequest func(storage *upcloud.Storage) (interface{}, error)
	Service      service.Storage
	ui.HandleContext
}

func (s Request) Send(args []string) (interface{}, error) {
	if s.ExactlyOne && len(args) != 1 {
		return nil, fmt.Errorf("single storage uuid is required")
	}
	if len(args) < 1 {
		return nil, fmt.Errorf("at least one storage uuid is required")
	}

	storages, err := SearchAllStorages(args, s.Service, true)
	if err != nil {
		return nil, err
	}

	var requests []interface{}
	for _, storage := range storages {
		request, err := s.BuildRequest(storage)
		if err != nil {
			return nil, err
		}
		requests = append(requests, request)
	}

	return s.Handle(requests)
}

func GetArgCompFn(s service.Storage) func(toComplete string) ([]string, cobra.ShellCompDirective) {
	return func(toComplete string) ([]string, cobra.ShellCompDirective) {
		storages, err := s.GetStorages(&request.GetStoragesRequest{})
		if err != nil {
			return nil, cobra.ShellCompDirectiveDefault
		}
		var vals []string
		for _, v := range storages.Storages {
			vals = append(vals, v.UUID, v.Title)
		}
		return commands.MatchStringPrefix(vals, toComplete, false), cobra.ShellCompDirectiveNoSpace | cobra.ShellCompDirectiveNoFileComp
	}
}
