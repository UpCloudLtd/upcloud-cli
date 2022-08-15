package format

import (
	"fmt"

	"github.com/UpCloudLtd/upcloud-go-api/v4/upcloud"
	"github.com/jedib0t/go-pretty/v6/text"
)

// serverStateColour is a helper mapping server states to colours
func serverStateColour(state string) text.Colors {
	switch state {
	case upcloud.ServerStateStarted:
		return text.Colors{text.FgGreen}
	case upcloud.ServerStateError:
		return text.Colors{text.FgHiRed, text.Bold}
	case upcloud.ServerStateMaintenance:
		return text.Colors{text.FgYellow}
	default:
		return text.Colors{text.FgHiBlack}
	}
}

// ServerState implements Format function for server states
func ServerState(val interface{}) (text.Colors, string, error) {
	state, ok := val.(string)
	if !ok {
		return nil, "", fmt.Errorf("cannot parse server state from %T, expected string", val)
	}

	if state == "" {
		return PossiblyUnknownString(state)
	}

	return nil, serverStateColour(state).Sprint(state), nil
}
