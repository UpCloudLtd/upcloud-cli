package format

import (
	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud"
	"github.com/jedib0t/go-pretty/v6/text"
)

// networkPeeringStateStateColour maps network peering states to colours
func networkPeeringStateStateColour(state upcloud.NetworkPeeringState) text.Colors {
	switch state {
	case upcloud.NetworkPeeringStateActive:
		return text.Colors{text.FgGreen}
	case upcloud.NetworkPeeringStateProvisioning:
		return text.Colors{text.FgYellow}
	case upcloud.NetworkPeeringStateDeletedPeerNetwork, upcloud.NetworkPeeringStateError, upcloud.NetworkPeeringStateMissingLocalRouter, upcloud.NetworkPeeringStateMissingPeerRouter:
		return text.Colors{text.FgRed}
	default:
		return text.Colors{text.FgHiBlack}
	}
}

// NetworkPeeringState implements Format function for network peering states
func NetworkPeeringState(val interface{}) (text.Colors, string, error) {
	return usingColorFunction(networkPeeringStateStateColour, val)
}
