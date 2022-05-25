package completion

import (
	"github.com/UpCloudLtd/upcloud-cli/internal/service"
	"github.com/UpCloudLtd/upcloud-go-api/v4/upcloud/request"
	"github.com/spf13/cobra"
)

// Database implements argument completion for databases, by uuid or title.
type Database struct {
}

// make sure Database implements the interface
var _ Provider = Database{}

// CompleteArgument implements completion.Provider
func (s Database) CompleteArgument(svc service.AllServices, toComplete string) ([]string, cobra.ShellCompDirective) {
	databases, err := svc.GetManagedDatabases(&request.GetManagedDatabasesRequest{})
	if err != nil {
		return None(toComplete)
	}
	var vals []string
	for _, db := range databases {
		vals = append(vals, db.UUID, db.Title)
	}
	return MatchStringPrefix(vals, toComplete, true), cobra.ShellCompDirectiveNoFileComp
}
