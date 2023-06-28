package format

import (
	"github.com/jedib0t/go-pretty/v6/text"
)

// serverGroupAntiAffinityStateColour is a helper mapping server group anti-affinity states to colours
func serverGroupAntiAffinityStateColour(state string) text.Colors {
	switch state {
	case "met":
		return text.Colors{text.FgGreen}
	case "unmet":
		return text.Colors{text.FgHiRed, text.Bold}
	default:
		return text.Colors{text.FgHiBlack}
	}
}

// ServerGroupAntiAffinityState implements Format function for server group anti-affinity states
func ServerGroupAntiAffinityState(val interface{}) (text.Colors, string, error) {
	return usingColorFunction(serverGroupAntiAffinityStateColour, val)
}
