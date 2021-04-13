package storage

import (
	"fmt"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/request"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/service"

	"github.com/UpCloudLtd/upcloud-cli/internal/commands"
)

var (
	maxStorageActions = 10
	// CachedStorages stores the cached list of storages in order to not hit the service more than once
	// TODO: refactor
	CachedStorages []upcloud.Storage
)

// BaseStorageCommand creates the base "storage" command
func BaseStorageCommand() commands.Command {
	return &storageCommand{
		commands.New("storage", "Manage storages", ""),
	}
}

type storageCommand struct {
	*commands.BaseCommand
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
	if len(CachedStorages) == 0 {
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

// SearchSingleStorage returns exactly one storage where title or uuid matches uuidOrTitle
// TODO: remove the cross-command dependencies
func SearchSingleStorage(uuidOrTitle string, service service.Storage) (*upcloud.Storage, error) {
	matchedResults, err := searchStorage(&CachedStorages, service, uuidOrTitle, true)
	if err != nil {
		return nil, err
	}
	return matchedResults[0], nil
}
