package format

import (
	"github.com/UpCloudLtd/upcloud-go-api/v6/upcloud"
	"github.com/jedib0t/go-pretty/v6/text"
)

// kubernetesClusterStateColour maps kubernetes cluster states to colours
func kubernetesClusterStateColour(state upcloud.KubernetesClusterState) text.Colors {
	switch state {
	case upcloud.KubernetesClusterStateRunning:
		return text.Colors{text.FgGreen}
	case upcloud.KubernetesClusterStatePending:
		return text.Colors{text.FgYellow}
	default:
		return text.Colors{text.FgHiBlack}
	}
}

// kubernetesNodeGroupStateColour maps kubernetes node-group states to colours
func kubernetesNodeGroupStateColour(state upcloud.KubernetesNodeGroupState) text.Colors {
	switch state {
	case upcloud.KubernetesNodeGroupStateRunning:
		return text.Colors{text.FgGreen}
	case upcloud.KubernetesNodeGroupStatePending:
		return text.Colors{text.FgYellow}
	default:
		return text.Colors{text.FgHiBlack}
	}
}

// KubernetesClusterState implements Format function for Kubernetes cluster states
func KubernetesClusterState(val interface{}) (text.Colors, string, error) {
	return usingColorFunction(kubernetesClusterStateColour, val)
}

// KubernetesNodeGroupState implements Format function for Kubernetes node-group states
func KubernetesNodeGroupState(val interface{}) (text.Colors, string, error) {
	return usingColorFunction(kubernetesNodeGroupStateColour, val)
}
