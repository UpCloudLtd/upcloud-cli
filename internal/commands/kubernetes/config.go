package kubernetes

import (
	"fmt"

	"github.com/UpCloudLtd/upcloud-cli/v2/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v2/internal/completion"
	"github.com/UpCloudLtd/upcloud-cli/v2/internal/output"
	"github.com/UpCloudLtd/upcloud-cli/v2/internal/resolver"
	"github.com/UpCloudLtd/upcloud-go-api/v4/upcloud/request"

	"github.com/spf13/pflag"
	"k8s.io/client-go/tools/clientcmd"
)

type configCommand struct {
	*commands.BaseCommand
	resolver.CachingKubernetes
	completion.Kubernetes
	write string
}

// ConfigCommand creates the "connection config" command
func ConfigCommand() commands.Command {
	return &configCommand{
		BaseCommand: commands.New(
			"config",
			"Output Kubernetes cluster kubeconfig",
			`upctl kubernetes config 0fa980c4-0e4f-460b-9869-11b7bd62b831 --output human`,
			`upctl kubernetes config 0fa980c4-0e4f-460b-9869-11b7bd62b831 --output yaml --write $KUBECONFIG`,
			`upctl kubernetes config 0fa980c4-0e4f-460b-9869-11b7bd62b831 --output yaml --write ./my_kubeconfig.yaml`,
		),
	}
}

// InitCommand implements Command.InitCommand
func (s *configCommand) InitCommand() {
	flagSet := &pflag.FlagSet{}
	flagSet.StringVar(&s.write, "write", "", "Target file path where to write config output. Default value \"\" (empty string) implies no writing will be done.")
	s.AddFlags(flagSet)
}

// Execute implements commands.MultipleArgumentCommand
func (s *configCommand) Execute(exec commands.Executor, uuid string) (output.Output, error) {
	svc := exec.All()

	msg := fmt.Sprintf("Getting kubeconfig for Kubernetes cluster %s", uuid)
	exec.PushProgressStarted(msg)

	kubeconfig, err := svc.GetKubernetesKubeconfig(&request.GetKubernetesKubeconfigRequest{
		UUID: uuid,
	})
	if err != nil {
		return commands.HandleError(exec, msg, err)
	}

	uksConfig, err := clientcmd.Load([]byte(kubeconfig))
	if err != nil {
		return commands.HandleError(exec, msg, err)
	}

	clusters := make([]output.TableRow, 0)
	for _, v := range uksConfig.Clusters {
		clusters = append(clusters, output.TableRow{
			v.Server,
			v.InsecureSkipTLSVerify,
		})
	}

	contexts := make([]output.TableRow, 0)
	for _, v := range uksConfig.Contexts {
		contexts = append(contexts, output.TableRow{
			v.Cluster,
			v.AuthInfo,
		})
	}

	exec.PushProgressSuccess(msg)

	if s.Cobra().Flag("write").Changed {
		msg := fmt.Sprintf("Writing kubeconfig for Kubernetes cluster %s to destination %s", uuid, s.write)
		exec.PushProgressStarted(msg)
		exec.PushProgressSuccess(msg)
	}

	return output.MarshaledWithHumanOutput{
		Value: kubeconfig,
		Output: output.Combined{
			output.CombinedSection{
				Contents: output.Details{
					Sections: []output.DetailSection{
						{
							Title: "Overview:",
							Rows: []output.DetailRow{
								{Title: "Current context:", Value: uksConfig.CurrentContext},
							},
						},
					},
				},
			},
			output.CombinedSection{
				Title: "Clusters",
				Contents: output.Table{
					Columns: []output.TableColumn{
						{Key: "server", Header: "Server"},
						{Key: "insecure_skip_tls_verify", Header: "Insecure skip TLS verify"},
					},
					Rows:       clusters,
					HideHeader: false,
				},
			},
			output.CombinedSection{
				Title: "Contexts",
				Contents: output.Table{
					Columns: []output.TableColumn{
						{Key: "cluster", Header: "Cluster"},
						{Key: "authinfo", Header: "Authinfo"},
					},
					Rows:       contexts,
					HideHeader: false,
				},
			},
		},
	}, nil
}
