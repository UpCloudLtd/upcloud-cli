package kubernetes

import (
	"fmt"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"

	"github.com/UpCloudLtd/progress/messages"
	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud/request"
)

// waitForClusterState waits for cluster to reach given state and updates progress message with key matching given msg. Finally, progress message is updated back to given msg and either done state or timeout warning.
func waitForClusterState(uuid string, state upcloud.KubernetesClusterState, exec commands.Executor, msg string) {
	exec.PushProgressUpdateMessage(msg, fmt.Sprintf("Waiting for cluster %s to be in %s state", uuid, state))

	if _, err := exec.All().WaitForKubernetesClusterState(exec.Context(), &request.WaitForKubernetesClusterStateRequest{
		UUID:         uuid,
		DesiredState: state,
	}); err != nil {
		exec.PushProgressUpdate(messages.Update{
			Key:     msg,
			Message: msg,
			Status:  messages.MessageStatusWarning,
			Details: "Error: " + err.Error(),
		})
		return
	}

	exec.PushProgressUpdateMessage(msg, msg)
	exec.PushProgressSuccess(msg)
}
