package storagebackup

import "github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"

// BackupCommand creates the "storage backup" command
func BackupCommand() commands.Command {
	return &backupCommand{
		commands.New("backup", "Manage backups"),
	}
}

type backupCommand struct {
	*commands.BaseCommand
}
