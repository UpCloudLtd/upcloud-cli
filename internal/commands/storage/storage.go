package storage

import (
	"errors"
	"fmt"
	"time"

	"github.com/UpCloudLtd/upcloud-go-api/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/request"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/service"
	"github.com/jedib0t/go-pretty/v6/text"

	"github.com/UpCloudLtd/cli/internal/commands"
	"github.com/UpCloudLtd/cli/internal/validation"
)

var (
	maxStorageActions = 10
	cachedStorages    []upcloud.Storage
	cachedTemplates   []upcloud.Storage
	cachedServers     []upcloud.Server
)

const minStorageSize = 10

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

func searchStorage(storagesPtr *[]upcloud.Storage, service service.Storage, uuidOrTitle string, unique bool) (*upcloud.Storage, error) {
	if storagesPtr == nil || service == nil {
		return nil, fmt.Errorf("no storages or service passed")
	}
	storages := *storagesPtr
	if err := validation.Uuid4(uuidOrTitle); err != nil || storages == nil {
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
	return matched[0], nil
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
