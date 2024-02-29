package completion

import (
	"context"
	"fmt"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/service"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/utils"
	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud/request"
	"github.com/spf13/cobra"
)

// Database implements argument completion for databases, by uuid or title.
type Database struct{}

// make sure Database implements the interface
var _ Provider = Database{}

// CompleteArgument implements completion.Provider
func (s Database) CompleteArgument(ctx context.Context, svc service.AllServices, toComplete string) ([]string, cobra.ShellCompDirective) {
	databases, err := svc.GetManagedDatabases(ctx, &request.GetManagedDatabasesRequest{})
	if err != nil {
		return None(toComplete)
	}
	var vals []string
	for _, db := range databases {
		vals = append(vals, db.UUID, db.Title)
	}
	return MatchStringPrefix(vals, toComplete, true), cobra.ShellCompDirectiveNoFileComp
}

// DatabaseType implements argument completion for database types.
type DatabaseType struct{}

// make sure DatabaseType implements the interface
var _ Provider = DatabaseType{}

// CompleteArgument implements completion.Provider
func (s DatabaseType) CompleteArgument(ctx context.Context, svc service.AllServices, toComplete string) ([]string, cobra.ShellCompDirective) {
	dbTypes, err := svc.GetManagedDatabaseServiceTypes(ctx, &request.GetManagedDatabaseServiceTypesRequest{})
	if err != nil {
		return None(toComplete)
	}
	var vals []string
	for _, t := range dbTypes {
		vals = append(vals, t.Name)
	}
	return MatchStringPrefix(vals, toComplete, true), cobra.ShellCompDirectiveNoFileComp
}

// DatabaseProperty implements argument completion for database properties.
type DatabaseProperty struct {
	ServiceType string
}

// make sure DatabaseType implements the interface
var _ Provider = DatabaseProperty{}

// CompleteArgument implements completion.Provider
func (s DatabaseProperty) CompleteArgument(ctx context.Context, svc service.AllServices, toComplete string) ([]string, cobra.ShellCompDirective) {
	dbType, err := svc.GetManagedDatabaseServiceType(ctx, &request.GetManagedDatabaseServiceTypeRequest{Type: s.ServiceType})
	if err != nil {
		return None(toComplete)
	}

	properties := utils.GetFlatDatabaseProperties(dbType.Properties)
	var vals []string
	for key, details := range properties {
		description := details.Title
		if description == key && details.Description != "" {
			description = details.Description
		}

		vals = append(vals, fmt.Sprintf("%s\t%s", key, description))
	}
	return MatchStringPrefix(vals, toComplete, true), cobra.ShellCompDirectiveNoFileComp
}
