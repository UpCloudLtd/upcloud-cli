package format

import (
	"github.com/UpCloudLtd/upcloud-go-api/v6/upcloud"
	"github.com/jedib0t/go-pretty/v6/text"
)

// databaseStateColour maps database states to colours
func databaseStateColour(state upcloud.ManagedDatabaseState) text.Colors {
	switch state {
	case upcloud.ManagedDatabaseStateRunning:
		return text.Colors{text.FgGreen}
	case "rebuilding", "rebalancing":
		return text.Colors{text.FgYellow}
	default:
		return text.Colors{text.FgHiBlack}
	}
}

// DatabaseState implements Format function for database states
func DatabaseState(val interface{}) (text.Colors, string, error) {
	return usingColorFunction(databaseStateColour, val)
}
