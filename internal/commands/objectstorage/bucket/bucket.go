package bucket

import (
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
)

// BaseBucketCommand creates the base "object-storage bucket" command
func BaseBucketCommand() commands.Command {
	return &bucketCommand{
		BaseCommand: commands.New("bucket", "Manage buckets in managed object storage services"),
	}
}

type bucketCommand struct {
	*commands.BaseCommand
}
