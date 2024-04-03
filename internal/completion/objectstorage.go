package completion

import (
	"context"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/service"
	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud/request"
	"github.com/spf13/cobra"
)

// ObjectStorage implements argument completion for gateways, by uuid or name.
type ObjectStorage struct{}

// make sure ObjectStorage implements the interface
var _ Provider = ObjectStorage{}

// CompleteArgument implements completion.Provider
func (s ObjectStorage) CompleteArgument(ctx context.Context, svc service.AllServices, toComplete string) ([]string, cobra.ShellCompDirective) {
	objectstorages, err := svc.GetManagedObjectStorages(ctx, &request.GetManagedObjectStoragesRequest{})
	if err != nil {
		return None(toComplete)
	}
	var vals []string
	for _, objsto := range objectstorages {
		vals = append(vals, objsto.UUID, objsto.Name)
	}

	return MatchStringPrefix(vals, toComplete, true), cobra.ShellCompDirectiveNoFileComp
}
