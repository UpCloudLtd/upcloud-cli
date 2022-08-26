package format

import (
	"fmt"

	"github.com/UpCloudLtd/upcloud-go-api/v4/upcloud"
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
	state, ok := val.(string)
	if !ok {
		return nil, "", fmt.Errorf("cannot parse storage state from %T, expected string", val)
	}

	if state == "" {
		return PossiblyUnknownString(state)
	}

	return nil, storageStateColour(state).Sprint(state), nil
}
