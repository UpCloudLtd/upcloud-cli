package objectstorage

import (
	"fmt"
	"time"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/completion"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/config"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/output"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/resolver"

	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud/request"
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
	resolver.CachingObjectStorage
	completion.ObjectStorage

	deleteUsers    config.OptionalBoolean
	deletePolicies config.OptionalBoolean
	deleteBuckets  config.OptionalBoolean
	wait           config.OptionalBoolean
}

// InitCommand implements Command.InitCommand
func (c *deleteCommand) InitCommand() {
	flags := &pflag.FlagSet{}
	config.AddToggleFlag(flags, &c.deleteUsers, "delete-users", false, "Delete all users from the service before deleting the object storage instance.")
	config.AddToggleFlag(flags, &c.deletePolicies, "delete-policies", false, "Delete all policies from the service before deleting the object storage instance.")
	config.AddToggleFlag(flags, &c.deleteBuckets, "delete-buckets", false, "Delete all buckets from the service before deleting the object storage instance.")
	config.AddToggleFlag(flags, &c.wait, "wait", false, "Wait until the object storage instance has been deleted before returning.")
	c.AddFlags(flags)

	// Deprecating objectstorage and objsto in favour of object-storage
	// TODO: Remove this in the future
	commands.SetSubcommandDeprecationHelp(c, []string{"objectstorage", "objsto"})
}

func Delete(exec commands.Executor, uuid string, deleteUsers, deletePolicies, deleteBuckets, wait bool) (output.Output, error) {
	svc := exec.All()
	msg := fmt.Sprintf("Deleting object storage service %v", uuid)
	exec.PushProgressStarted(msg)

	if deleteUsers {
		exec.PushProgressUpdateMessage(msg, fmt.Sprintf("Deleting users from the service %s", uuid))

		users, err := svc.GetManagedObjectStorageUsers(exec.Context(), &request.GetManagedObjectStorageUsersRequest{ServiceUUID: uuid})
		if err != nil {
			return commands.HandleError(exec, msg, err)
		}

		for _, user := range users {
			err = svc.DeleteManagedObjectStorageUser(exec.Context(), &request.DeleteManagedObjectStorageUserRequest{
				ServiceUUID: uuid,
				Username:    user.Username,
			})
			if err != nil {
				return commands.HandleError(exec, msg, err)
			}
		}
	}

	if deletePolicies {
		exec.PushProgressUpdateMessage(msg, fmt.Sprintf("Deleting policies from the service %s", uuid))

		policies, err := svc.GetManagedObjectStoragePolicies(exec.Context(), &request.GetManagedObjectStoragePoliciesRequest{ServiceUUID: uuid})
		if err != nil {
			return commands.HandleError(exec, msg, err)
		}

		for _, policy := range policies {
			if policy.System {
				continue
			}

			err = svc.DeleteManagedObjectStoragePolicy(exec.Context(), &request.DeleteManagedObjectStoragePolicyRequest{
				ServiceUUID: uuid,
				Name:        policy.Name,
			})
			if err != nil {
				return commands.HandleError(exec, msg, err)
			}
		}
	}

	if deleteBuckets {
		exec.PushProgressUpdateMessage(msg, fmt.Sprintf("Deleting buckets from the service %s", uuid))

		buckets, err := svc.GetManagedObjectStorageBucketMetrics(exec.Context(), &request.GetManagedObjectStorageBucketMetricsRequest{ServiceUUID: uuid})
		if err != nil {
			return commands.HandleError(exec, msg, err)
		}

		for _, bucket := range buckets {
			err = svc.DeleteManagedObjectStorageBucket(exec.Context(), &request.DeleteManagedObjectStorageBucketRequest{
				ServiceUUID: uuid,
				Name:        bucket.Name,
			})
			if err != nil {
				return commands.HandleError(exec, msg, err)
			}
		}

		exec.PushProgressUpdateMessage(msg, fmt.Sprintf("Waiting buckets to be deleted from the service %s", uuid))
		err = waitUntilBucketsHaveBeenDeleted(exec, uuid)
		if err != nil {
			return commands.HandleError(exec, msg, err)
		}
	}

	exec.PushProgressUpdateMessage(msg, msg)
	err := svc.DeleteManagedObjectStorage(exec.Context(), &request.DeleteManagedObjectStorageRequest{
		UUID: uuid,
	})
	if err != nil {
		return commands.HandleError(exec, msg, err)
	}

	if wait {
		exec.PushProgressUpdateMessage(msg, fmt.Sprintf("Waiting for object storage service %s to be deleted", uuid))
		err = svc.WaitForManagedObjectStorageDeletion(exec.Context(), &request.WaitForManagedObjectStorageDeletionRequest{UUID: uuid})
		if err != nil {
			return commands.HandleError(exec, msg, err)
		}
		exec.PushProgressUpdateMessage(msg, msg)
	}

	exec.PushProgressSuccess(msg)

	return output.None{}, nil
}

// Execute implements commands.MultipleArgumentCommand
func (c *deleteCommand) Execute(exec commands.Executor, arg string) (output.Output, error) {
	// Deprecating objectstorage and objsto in favour of object-storage
	// TODO: Remove this in the future
	commands.SetSubcommandExecutionDeprecationMessage(c, []string{"objectstorage", "objsto"}, "object-storage")

	return Delete(exec, arg, c.deleteUsers.Value(), c.deletePolicies.Value(), c.deleteBuckets.Value(), c.wait.Value())
}

func waitUntilBucketsHaveBeenDeleted(exec commands.Executor, serviceUUID string) error {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	ctx := exec.Context()
	svc := exec.All()

	for i := 0; ; i++ {
		select {
		case <-ticker.C:
			buckets, err := svc.GetManagedObjectStorageBucketMetrics(exec.Context(), &request.GetManagedObjectStorageBucketMetricsRequest{
				ServiceUUID: serviceUUID,
			})
			if err != nil {
				return err
			}

			if len(buckets) == 0 {
				return nil
			}
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}
