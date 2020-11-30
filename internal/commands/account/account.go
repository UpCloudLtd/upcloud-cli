package account

import (
	"github.com/UpCloudLtd/cli/internal/commands"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud"
)

type Account interface {
	GetAccount() (*upcloud.Account, error)
}

func AccountCommand() commands.Command {
	return &accountCommand{commands.New("account", "Manage account")}
}

type accountCommand struct {
	*commands.BaseCommand
}
