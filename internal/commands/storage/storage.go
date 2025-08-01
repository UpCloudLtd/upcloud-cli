package storage

import (
	"context"
	"fmt"
	"time"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"

	"github.com/UpCloudLtd/progress/messages"
	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud/request"
)

var (
	maxStorageActions = 10
	tiers             = []string{
		upcloud.StorageTierMaxIOPS,
		upcloud.StorageTierStandard,
		upcloud.StorageTierHDD,
	}
	backupIntervals = []string{
		upcloud.BackupRuleIntervalDaily,
		upcloud.BackupRuleIntervalMonday,
		upcloud.BackupRuleIntervalTuesday,
		upcloud.BackupRuleIntervalWednesday,
		upcloud.BackupRuleIntervalThursday,
		upcloud.BackupRuleIntervalFriday,
		upcloud.BackupRuleIntervalSaturday,
		upcloud.BackupRuleIntervalSunday,
	}
	// CachedStorages stores the cached list of storages in order to not hit the service more than once
	// TODO: refactor
	CachedStorages []upcloud.Storage
)

// BaseStorageCommand creates the base "storage" command
func BaseStorageCommand() commands.Command {
	return &storageCommand{
		commands.New("storage", "Manage storages"),
	}
}

type storageCommand struct {
	*commands.BaseCommand
}

// InitCommand implements Command.InitCommand
func (st *storageCommand) InitCommand() {
	st.Cobra().Aliases = []string{"st"}
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

func searchStorage(storagesPtr *[]upcloud.Storage, exec commands.Executor, uuidOrTitle string, unique bool) ([]*upcloud.Storage, error) {
	if storagesPtr == nil || exec == nil {
		return nil, fmt.Errorf("no storages or executor passed")
	}
	storages := *storagesPtr
	if len(CachedStorages) == 0 {
		res, err := exec.All().GetStorages(exec.Context(), &request.GetStoragesRequest{})
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
func SearchSingleStorage(uuidOrTitle string, exec commands.Executor) (*upcloud.Storage, error) {
	matchedResults, err := searchStorage(&CachedStorages, exec, uuidOrTitle, true)
	if err != nil {
		return nil, err
	}
	return matchedResults[0], nil
}

// waitForStorageState waits for storage to reach given state and updates progress message with key matching given msg. Finally, progress message is updated back to given msg and either done state or timeout warning.
func waitForStorageState(uuid, state string, exec commands.Executor, msg string) {
	exec.PushProgressUpdateMessage(msg, fmt.Sprintf("Waiting for storage %s to be in %s state", uuid, state))

	ctx, cancel := context.WithTimeout(exec.Context(), 15*time.Minute)
	defer cancel()

	if _, err := exec.All().WaitForStorageState(ctx, &request.WaitForStorageStateRequest{
		UUID:         uuid,
		DesiredState: state,
	}); err != nil {
		exec.PushProgressUpdate(messages.Update{
			Key:     msg,
			Message: msg,
			Status:  messages.MessageStatusWarning,
			Details: "Error: " + err.Error(),
		})
		return
	}

	exec.PushProgressUpdateMessage(msg, msg)
	exec.PushProgressSuccess(msg)
}
