package completion

import (
	"github.com/UpCloudLtd/upcloud-cli/internal/service"

	"github.com/UpCloudLtd/upcloud-go-api/upcloud/request"
	"github.com/spf13/cobra"
)

// Storage implements argument completion for routers, by uuid, name or hostname.
type Storage struct {
}

// make sure Storage implements the interface
var _ Provider = Storage{}

// CompleteArgument implements completion.Provider
func (s Storage) CompleteArgument(svc service.AllServices, toComplete string) ([]string, cobra.ShellCompDirective) {
	storages, err := svc.GetStorages(&request.GetStoragesRequest{})
	if err != nil {
		return None(toComplete)
	}
	var vals []string
	for _, v := range storages.Storages {
		vals = append(vals, v.UUID, v.Title)
	}
	return MatchStringPrefix(vals, toComplete, true), cobra.ShellCompDirectiveNoFileComp
}
