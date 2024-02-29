package objectstorage

import (
	"fmt"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/config"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/output"

	"github.com/UpCloudLtd/upcloud-go-api/v7/upcloud/request"
	"github.com/spf13/pflag"
)

// DeleteCommand creates the "objectstorage delete" command
func DeleteCommand() commands.Command {
	return &deleteCommand{
		BaseCommand: commands.New(
			"delete",
			"Delete a Managed object storage service",
			"upctl objectstorage delete 55199a44-4751-4e27-9394-7c7661910be8",
		),
	}
}

type deleteCommand struct {
	*commands.BaseCommand
	deleteUsers    config.OptionalBoolean
	deletePolicies config.OptionalBoolean
}

// InitCommand implements Command.InitCommand
func (c *deleteCommand) InitCommand() {
	flags := &pflag.FlagSet{}
	config.AddToggleFlag(flags, &c.deleteUsers, "delete-users", false, "Delete all users from the service before deleting the object storage instance.")
	config.AddToggleFlag(flags, &c.deletePolicies, "delete-policies", false, "Delete all policies from the service before deleting the object storage instance.")
	c.AddFlags(flags)
}

// Execute implements commands.MultipleArgumentCommand
func (c *deleteCommand) Execute(exec commands.Executor, arg string) (output.Output, error) {
	svc := exec.All()
	msg := fmt.Sprintf("Deleting object storage service %v", arg)
	exec.PushProgressStarted(msg)

	if c.deleteUsers.Value() {
		exec.PushProgressUpdateMessage(msg, fmt.Sprintf("Deleting users from the service %s", arg))

		users, err := svc.GetManagedObjectStorageUsers(exec.Context(), &request.GetManagedObjectStorageUsersRequest{ServiceUUID: arg})
		if err != nil {
			return commands.HandleError(exec, msg, err)
		}

		for _, user := range users {
			err = svc.DeleteManagedObjectStorageUser(exec.Context(), &request.DeleteManagedObjectStorageUserRequest{
				ServiceUUID: arg,
				Username:    user.Username,
			})
			if err != nil {
				return commands.HandleError(exec, msg, err)
			}
		}
	}

	if c.deletePolicies.Value() {
		exec.PushProgressUpdateMessage(msg, fmt.Sprintf("Deleting policies from the service %s", arg))

		policies, err := svc.GetManagedObjectStoragePolicies(exec.Context(), &request.GetManagedObjectStoragePoliciesRequest{ServiceUUID: arg})
		if err != nil {
			return commands.HandleError(exec, msg, err)
		}

		for _, policy := range policies {
			err = svc.DeleteManagedObjectStoragePolicy(exec.Context(), &request.DeleteManagedObjectStoragePolicyRequest{
				ServiceUUID: arg,
				Name:        policy.Name,
			})
			if err != nil {
				return commands.HandleError(exec, msg, err)
			}
		}
	}

	err := svc.DeleteManagedObjectStorage(exec.Context(), &request.DeleteManagedObjectStorageRequest{
		UUID: arg,
	})
	if err != nil {
		return commands.HandleError(exec, msg, err)
	}

	exec.PushProgressSuccess(msg)

	return output.None{}, nil
}
