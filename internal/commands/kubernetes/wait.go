package kubernetes

import (
	"context"
	"fmt"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"

	"github.com/UpCloudLtd/progress/messages"
	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud/request"
	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud/service"
)

// waitForClusterState waits for cluster to reach given state and updates progress message with key matching given msg. Finally, progress message is updated back to given msg and either done state or timeout warning.
func WaitForClusterState(uuid string, state upcloud.KubernetesClusterState, exec commands.Executor, msg string) {
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

func allNodeGroupsRunning(groups []upcloud.KubernetesNodeGroup) bool {
	for _, group := range groups {
		if group.State != upcloud.KubernetesNodeGroupStateRunning {
			return false
		}
	}
	return true
}

func waitUntilNodeGroupsRunning(uuid string, exec commands.Executor) error {
	svc := exec.All()

	_, err := service.Retry(exec.Context(), func(i int, ctx context.Context) (*[]upcloud.KubernetesNodeGroup, error) {
		groups, err := svc.GetKubernetesNodeGroups(exec.Context(), &request.GetKubernetesNodeGroupsRequest{
			ClusterUUID: uuid,
		})
		if err != nil {
			return nil, err
		}
		if allNodeGroupsRunning(groups) {
			return &groups, nil
		}

		return nil, nil //nolint:nilnil // Continue retrying
	}, nil)
	return err
}

func waitUntilClusterAndNodeGroupsRunning(uuid string, exec commands.Executor, msg string) {
	exec.PushProgressUpdateMessage(msg, fmt.Sprintf("Waiting for cluster %s to be in running state", uuid))

	if _, err := exec.All().WaitForKubernetesClusterState(exec.Context(), &request.WaitForKubernetesClusterStateRequest{
		UUID:         uuid,
		DesiredState: upcloud.KubernetesClusterStateRunning,
	}); err != nil {
		exec.PushProgressUpdate(messages.Update{
			Key:     msg,
			Message: msg,
			Status:  messages.MessageStatusWarning,
			Details: "Error: " + err.Error(),
		})
		return
	}

	exec.PushProgressUpdateMessage(msg, fmt.Sprintf("Waiting for cluster %s node-groups to be in running state", uuid))
	if err := waitUntilNodeGroupsRunning(uuid, exec); err != nil {
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
