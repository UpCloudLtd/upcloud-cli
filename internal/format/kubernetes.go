package format

import (
	"github.com/UpCloudLtd/upcloud-go-api/v6/upcloud"
	"github.com/jedib0t/go-pretty/v6/text"
)

// kubernetesOperationalStateColour maps load balancer states to colours
func kubernetesOperationalStateColour(state upcloud.KubernetesClusterState) text.Colors {
	switch state {
	case upcloud.KubernetesClusterStateRunning:
		return text.Colors{text.FgGreen}
	case upcloud.KubernetesClusterStatePending:
		return text.Colors{text.FgYellow}
	default:
		return text.Colors{text.FgHiBlack}
	}
}

// KubernetesState implements Format function for Kubernetes states
func KubernetesState(val interface{}) (text.Colors, string, error) {
	return usingColorFunction(kubernetesOperationalStateColour, val)
}
