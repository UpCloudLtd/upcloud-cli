package kubernetes

import (
	"errors"
	"fmt"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/completion"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/output"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/resolver"

	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud/request"
	"github.com/spf13/pflag"
	"gopkg.in/yaml.v3"
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
		"Absolute path for writing output. If the file exists, the config will be merged.")
	c.AddFlags(flagSet)

	// Deprecating uks in favor of k8s
	// TODO: Remove this in the future
	commands.SetSubcommandDeprecationHelp(c, []string{"uks"})
}

// Execute implements commands.MultipleArgumentCommand
func (c *configCommand) Execute(exec commands.Executor, uuid string) (output.Output, error) {
	// Deprecating uks
	// TODO: Remove this in the future
	commands.SetSubcommandExecutionDeprecationMessage(c, []string{"uks"}, "k8s")

	svc := exec.All()

	msg := fmt.Sprintf("Getting kubeconfig for Kubernetes cluster %s", uuid)
	exec.PushProgressStarted(msg)

	resp, err := svc.GetKubernetesKubeconfig(exec.Context(), &request.GetKubernetesKubeconfigRequest{
		UUID: uuid,
	})
	if err != nil {
		return commands.HandleError(exec, msg, err)
	}

	config, err := clientcmd.Load([]byte(resp))
	if err != nil {
		return commands.HandleError(exec, msg, err)
	}

	exec.PushProgressSuccess(msg)

	if c.Cobra().Flag("write").Changed {
		return c.write(exec, uuid, config)
	}

	return c.output(exec, config, resp, msg)
}

func (c *configCommand) output(exec commands.Executor, config *api.Config, resp string, msg string) (output.Output, error) {
	clusters := make([]output.TableRow, 0)
	for _, v := range config.Clusters {
		clusters = append(clusters, output.TableRow{
			v.Server,
			v.InsecureSkipTLSVerify,
		})
	}

	contexts := make([]output.TableRow, 0)
	for _, v := range config.Contexts {
		contexts = append(contexts, output.TableRow{
			v.Cluster,
			v.AuthInfo,
		})
	}

	var value interface{}
	err := yaml.Unmarshal([]byte(resp), &value)
	if err != nil {
		return commands.HandleError(exec, msg, err)
	}

	return output.MarshaledWithHumanOutput{
		Value: value,
		Output: output.Combined{
			output.CombinedSection{
				Contents: output.Details{
					Sections: []output.DetailSection{
						{
							Title: "Overview:",
							Rows: []output.DetailRow{
								{Title: "Current context:", Value: config.CurrentContext},
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

func (c *configCommand) write(exec commands.Executor, uuid string, config *api.Config) (output.Output, error) {
	msg := fmt.Sprintf("Writing kubeconfig for Kubernetes cluster %s to destination %s", uuid, c.pathOptions.GetDefaultFilename())
	exec.PushProgressStarted(msg)

	if c.Cobra().Flag("write").Value.String() == c.Cobra().Flag("write").DefValue {
		return commands.HandleError(exec, msg, errors.New("invalid write path"))
	}

	startingConfig, err := c.pathOptions.GetStartingConfig()
	if err != nil {
		return commands.HandleError(exec, msg, err)
	}

	err = clientcmd.ModifyConfig(c.pathOptions, mergeConfig(startingConfig, config), false)
	if err != nil {
		return commands.HandleError(exec, msg, err)
	}

	exec.PushProgressSuccess(msg)

	return output.None{}, nil
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
