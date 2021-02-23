package storage

import "github.com/UpCloudLtd/cli/internal/commands"

// BackupCommand creates the "storage backup" command
func BackupCommand() commands.Command {
	return &storageCommand{commands.New("backup", "Manage backups")}
}
