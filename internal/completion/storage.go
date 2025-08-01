package completion

import (
	"context"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/service"

	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud/request"
	"github.com/spf13/cobra"
)

// Storage implements argument completion for storages, by uuid or title.
type Storage struct{}

// make sure Storage implements the interface
var _ Provider = Storage{}

// CompleteArgument implements completion.Provider
func (s Storage) CompleteArgument(ctx context.Context, svc service.AllServices, toComplete string) ([]string, cobra.ShellCompDirective) {
	vals, err := storageCompletions(ctx, svc, &request.GetStoragesRequest{}, true)
	if err != nil {
		return None(toComplete)
	}
	return MatchStringPrefix(vals, toComplete, true), cobra.ShellCompDirectiveNoFileComp
}

// StorageUUID implements argument completion for storage UUIDs.
type StorageUUID struct{}

// make sure StorageUUID implements the interface
var _ Provider = StorageUUID{}

// CompleteArgument implements completion.Provider
func (s StorageUUID) CompleteArgument(ctx context.Context, svc service.AllServices, toComplete string) ([]string, cobra.ShellCompDirective) {
	vals, err := storageCompletions(ctx, svc, &request.GetStoragesRequest{}, false)
	if err != nil {
		return None(toComplete)
	}
	return MatchStringPrefix(vals, toComplete, true), cobra.ShellCompDirectiveNoFileComp
}

// StorageCDROMUUID implements argument completion for cd-rom storage UUIDs.
type StorageCDROMUUID struct{}

// make sure StorageCDROMUUID implements the interface
var _ Provider = StorageCDROMUUID{}

// CompleteArgument implements completion.Provider
func (s StorageCDROMUUID) CompleteArgument(ctx context.Context, svc service.AllServices, toComplete string) ([]string, cobra.ShellCompDirective) {
	vals, err := storageCompletions(ctx, svc, &request.GetStoragesRequest{Type: "cdrom"}, false)
	if err != nil {
		return None(toComplete)
	}
	return MatchStringPrefix(vals, toComplete, true), cobra.ShellCompDirectiveNoFileComp
}

func storageCompletions(ctx context.Context, svc service.AllServices, req *request.GetStoragesRequest, withTitles bool) ([]string, error) {
	storages, err := svc.GetStorages(ctx, req)
	if err != nil {
		return nil, err
	}
	var vals []string
	for _, v := range storages.Storages {
		vals = append(vals, v.UUID)
		if withTitles {
			vals = append(vals, v.Title)
		}
	}
	return vals, nil
}
