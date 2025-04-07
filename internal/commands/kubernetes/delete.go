package kubernetes

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/completion"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/config"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/output"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/resolver"
	"github.com/spf13/pflag"

	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud/request"
)

// DeleteCommand creates the "kubernetes delete" command
func DeleteCommand() commands.Command {
	return &deleteCommand{
		BaseCommand: commands.New(
			"delete",
			"Delete a Kubernetes cluster",
			"upctl kubernetes delete 55199a44-4751-4e27-9394-7c7661910be3",
			"upctl kubernetes delete my-kubernetes-cluster",
		),
	}
}

func (s *deleteCommand) InitCommand() {
	// Deprecating uks
	// TODO: Remove this in the future
	commands.SetSubcommandDeprecationHelp(s, []string{"uks"})

	flags := &pflag.FlagSet{}
	config.AddToggleFlag(flags, &s.wait, "wait", false, "Wait until the Kubernetes cluster has been deleted before returning.")
	s.AddFlags(flags)
}

type deleteCommand struct {
	*commands.BaseCommand
	resolver.CachingKubernetes
	completion.Kubernetes

	wait config.OptionalBoolean
}

func Delete(exec commands.Executor, uuid string, wait bool) (output.Output, error) {
	svc := exec.All()
	msg := fmt.Sprintf("Deleting Kubernetes cluster %v", uuid)
	exec.PushProgressStarted(msg)

	err := svc.DeleteKubernetesCluster(exec.Context(), &request.DeleteKubernetesClusterRequest{
		UUID: uuid,
	})
	if err != nil {
		return commands.HandleError(exec, msg, err)
	}

	if wait {
		exec.PushProgressUpdateMessage(msg, fmt.Sprintf("Waiting for Kubernetes cluster %s to be deleted", uuid))
		err = waitUntilClusterDeleted(exec, uuid)
		if err != nil {
			return commands.HandleError(exec, msg, err)
		}
		exec.PushProgressUpdateMessage(msg, msg)
	}

	exec.PushProgressSuccess(msg)

	return output.None{}, nil
}

// Execute implements commands.MultipleArgumentCommand
func (s *deleteCommand) Execute(exec commands.Executor, arg string) (output.Output, error) {
	// Deprecating uks
	// TODO: Remove this in the future
	commands.SetSubcommandExecutionDeprecationMessage(s, []string{"uks"}, "k8s")

	return Delete(exec, arg, s.wait.Value())
}

func waitUntilClusterDeleted(exec commands.Executor, uuid string) error {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	ctx := exec.Context()
	svc := exec.All()

	for i := 0; ; i++ {
		select {
		case <-ticker.C:
			_, err := svc.GetKubernetesCluster(exec.Context(), &request.GetKubernetesClusterRequest{
				UUID: uuid,
			})
			if err != nil {
				var ucErr *upcloud.Problem
				if errors.As(err, &ucErr) && ucErr.Status == http.StatusNotFound {
					return nil
				}

				return err
			}
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}
