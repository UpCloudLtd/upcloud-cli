package storage

import "github.com/UpCloudLtd/cli/internal/commands"

func BackupCommand() commands.Command {
  return &storageCommand{commands.New("backup", "Manage backups")}
}

type backupCommand struct {
  *commands.BaseCommand
}
