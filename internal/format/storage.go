package format

import (
	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud"
	"github.com/jedib0t/go-pretty/v6/text"
)

func storageStateColour(state string) text.Colors {
	switch state {
	case upcloud.StorageStateOnline, upcloud.StorageStateSyncing:
		return text.Colors{text.FgGreen}
	case upcloud.StorageStateError:
		return text.Colors{text.FgHiRed, text.Bold}
	case upcloud.StorageStateMaintenance:
		return text.Colors{text.FgYellow}
	case upcloud.StorageStateCloning, upcloud.StorageStateBackuping:
		return text.Colors{text.FgHiMagenta, text.Bold}
	default:
		return text.Colors{text.FgHiBlack}
	}
}

// StorageState implements Format function for storage states
func StorageState(val interface{}) (text.Colors, string, error) {
	return usingColorFunction(storageStateColour, val)
}
