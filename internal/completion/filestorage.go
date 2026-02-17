package completion

import (
	"context"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/service"
	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud/request"
	"github.com/spf13/cobra"
)

// FileStorage implements argument completion for filestorages, by uuid or name.
type FileStorage struct{}

// make sure FileStorage implements the interface
var _ Provider = FileStorage{}

// CompleteArgument implements completion.Provider
func (s FileStorage) CompleteArgument(ctx context.Context, svc service.AllServices, toComplete string) ([]string, cobra.ShellCompDirective) {
	filestorages, err := svc.GetFileStorages(ctx, &request.GetFileStoragesRequest{})
	if err != nil {
		return None(toComplete)
	}
	var vals []string
	for _, filesto := range filestorages {
		vals = append(vals, filesto.UUID, filesto.Name)
	}

	return MatchStringPrefix(vals, toComplete, true), cobra.ShellCompDirectiveNoFileComp
}
