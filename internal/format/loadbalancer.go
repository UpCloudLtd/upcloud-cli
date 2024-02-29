package format

import (
	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud"
	"github.com/jedib0t/go-pretty/v6/text"
)

// loadBalancerOperationalStateColour maps load balancer states to colours
func loadBalancerOperationalStateColour(state upcloud.LoadBalancerOperationalState) text.Colors {
	switch state {
	case upcloud.LoadBalancerOperationalStateRunning:
		return text.Colors{text.FgGreen}
	case upcloud.LoadBalancerOperationalStateCheckup, upcloud.LoadBalancerOperationalStatePending, upcloud.LoadBalancerOperationalStateSetupAgent, upcloud.LoadBalancerOperationalStateSetupDNS, upcloud.LoadBalancerOperationalStateSetupLB, upcloud.LoadBalancerOperationalStateSetupNetwork, upcloud.LoadBalancerOperationalStateSetupServer:
		return text.Colors{text.FgYellow}
	default:
		return text.Colors{text.FgHiBlack}
	}
}

// LoadBalancerState implements Format function for load-balancer states
func LoadBalancerState(val interface{}) (text.Colors, string, error) {
	return usingColorFunction(loadBalancerOperationalStateColour, val)
}
