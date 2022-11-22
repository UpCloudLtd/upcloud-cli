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
	"k8s.io/client-go/tools/clientcmd/api"
)

type configCommand struct {
	*commands.BaseCommand
	resolver.CachingKubernetes
	completion.Kubernetes
	pathOptions *clientcmd.PathOptions
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
		pathOptions: clientcmd.NewDefaultPathOptions(),
	}
}

// InitCommand implements Command.InitCommand
func (c *configCommand) InitCommand() {
	flagSet := &pflag.FlagSet{}
	flagSet.StringVar(
		&c.pathOptions.LoadingRules.ExplicitPath,
		"write",
		"",
		"Absolute path for writing output. If the file exists, results will be merged. Default value \"\" implies no writing will be done.")
	c.AddFlags(flagSet)
}

// Execute implements commands.MultipleArgumentCommand
func (c *configCommand) Execute(exec commands.Executor, uuid string) (output.Output, error) {
	svc := exec.All()

	msg := fmt.Sprintf("Getting kubeconfig for Kubernetes cluster %s", uuid)
	exec.PushProgressStarted(msg)

	resp, err := svc.GetKubernetesKubeconfig(&request.GetKubernetesKubeconfigRequest{
		UUID: uuid,
	})
	if err != nil {
		return commands.HandleError(exec, msg, err)
	}

	uksConfig, err := clientcmd.Load([]byte(resp))
	if err != nil {
		return commands.HandleError(exec, msg, err)
	}

	exec.PushProgressSuccess(msg)

	if c.Cobra().Flag("write").Changed {
		msg := fmt.Sprintf("Writing kubeconfig for Kubernetes cluster %s to destination %s", uuid, c.pathOptions.GetDefaultFilename())
		exec.PushProgressStarted(msg)

		startingConfig, err := c.pathOptions.GetStartingConfig()
		if err != nil {
			return commands.HandleError(exec, msg, err)
		}

		err = clientcmd.ModifyConfig(c.pathOptions, mergeConfig(startingConfig, uksConfig), false)
		if err != nil {
			return commands.HandleError(exec, msg, err)
		}

		exec.PushProgressSuccess(msg)

		return output.None{}, nil
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

	return output.MarshaledWithHumanOutput{
		Value: resp,
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

func mergeConfig(startingConfig, newConfig *api.Config) api.Config {
	startingConfig.CurrentContext = newConfig.CurrentContext

	for k, v := range newConfig.Clusters {
		startingConfig.Clusters[k] = v
	}
	for k, v := range newConfig.AuthInfos {
		startingConfig.AuthInfos[k] = v
	}
	for k, v := range newConfig.Contexts {
		startingConfig.Contexts[k] = v
	}

	return *startingConfig
}
