package format

import (
	"fmt"
	"strings"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/ui"
	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud"
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

// databaseIndexHealthColour maps database index health to colours
func databaseIndexHealthColour(health string) text.Colors {
	switch health {
	case "green":
		return text.Colors{text.FgGreen}
	case "red":
		return text.Colors{text.FgRed}
	case "yellow":
		return text.Colors{text.FgYellow}
	default:
		return text.Colors{text.FgHiBlack}
	}
}

// DatabaseIndexHealth implements Format function for database index health
func DatabaseIndexHealth(val interface{}) (text.Colors, string, error) {
	return usingColorFunction(databaseIndexHealthColour, val)
}

// databaseIndexStatusColour maps database index status to colours
func databaseIndexStatusColour(status string) text.Colors {
	switch status {
	case "closed":
		return text.Colors{text.FgRed}
	case "open":
		return text.Colors{text.FgGreen}
	default:
		return text.Colors{text.FgHiBlack}
	}
}

// DatabaseIndexState implements Format function for database index states
func DatabaseIndexState(val interface{}) (text.Colors, string, error) {
	return usingColorFunction(databaseIndexStatusColour, val)
}

// databaseSessionStatusColour maps database session status to colours
func databaseSessionStatusColour(status string) text.Colors {
	switch status {
	case "active":
		return text.Colors{text.FgGreen}
	case "idle":
		return text.Colors{text.FgYellow}
	default:
		return text.Colors{text.FgHiBlack}
	}
}

// DatabaseSessionState implements Format function for database session states
func DatabaseSessionState(val interface{}) (text.Colors, string, error) {
	return usingColorFunction(databaseSessionStatusColour, val)
}

// DatabaseSessionPID implements Format function for database session PID
func DatabaseSessionPID(val interface{}) (text.Colors, string, error) {
	if s, ok := val.(string); ok {
		return ui.DefaultUUUIDColours, strings.TrimPrefix(s, "pid_"), nil
	}

	return nil, "", fmt.Errorf("cannot parse '%v' (%T) as string", val, val)
}
