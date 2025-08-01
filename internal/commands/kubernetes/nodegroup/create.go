package nodegroup

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/completion"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/config"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/labels"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/namedargs"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/output"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/resolver"

	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud/request"
	"github.com/spf13/pflag"
)

type CreateNodeGroupParams struct {
	Count                int
	Name                 string
	Plan                 string
	SSHKeys              []string
	Storage              string
	KubeletArgs          []string
	Labels               []string
	Taints               []string
	UtilityNetworkAccess config.OptionalBoolean
}

func GetCreateNodeGroupFlagSet(p *CreateNodeGroupParams) *pflag.FlagSet {
	fs := &pflag.FlagSet{}

	fs.IntVar(&p.Count, "count", 0, "Number of nodes in the node group")
	fs.StringArrayVar(&p.KubeletArgs, "kubelet-arg", []string{}, "Arguments to use when executing kubelet in `argument=value` format")
	fs.StringArrayVar(&p.Labels, "label", []string{}, "Labels to describe the nodes in `key=value` format. Use multiple times to define multiple labels. Labels are forwarded to the kubernetes nodes.")
	fs.StringVar(&p.Name, "name", "", "Node group name")
	fs.StringVar(&p.Plan, "plan", "", "Server plan to use for nodes in the node group. Run `upctl server plans` to list all available plans.")
	fs.StringArrayVar(&p.SSHKeys, "ssh-key", []string{}, "SSH keys to be configured as authorized keys to the nodes.")
	fs.StringVar(&p.Storage, "storage", "", "Storage template to use when creating the nodes. Defaults to `UpCloud K8s` public template.")
	fs.StringArrayVar(&p.Taints, "taint", []string{}, "Taints to be configured to the nodes in `key=value:effect` format")
	config.AddEnableOrDisableFlag(fs, &p.UtilityNetworkAccess, true, "utility-network-access", "utility network access. If disabled, nodes in this group will not have access to utility network")

	commands.Must(fs.SetAnnotation("count", commands.FlagAnnotationNoFileCompletions, nil))
	commands.Must(fs.SetAnnotation("kubelet-arg", commands.FlagAnnotationNoFileCompletions, nil))
	commands.Must(fs.SetAnnotation("label", commands.FlagAnnotationNoFileCompletions, nil))
	commands.Must(fs.SetAnnotation("name", commands.FlagAnnotationNoFileCompletions, nil))
	commands.Must(fs.SetAnnotation("storage", commands.FlagAnnotationNoFileCompletions, nil))
	commands.Must(fs.SetAnnotation("taint", commands.FlagAnnotationNoFileCompletions, nil))

	return fs
}

func processKubeletArg(in string) (upcloud.KubernetesKubeletArg, error) {
	split := strings.SplitN(in, "=", 2)
	if len(split) < 2 {
		return upcloud.KubernetesKubeletArg{}, fmt.Errorf("invalid kubelet-arg: %s", in)
	}

	return upcloud.KubernetesKubeletArg{
		Key:   split[0],
		Value: split[1],
	}, nil
}

func processTaint(in string) (upcloud.KubernetesTaint, error) {
	r := regexp.MustCompile(`^(.+)=(.+):(.+)`)
	s := r.FindStringSubmatch(in)
	if len(s) < 4 {
		return upcloud.KubernetesTaint{}, fmt.Errorf("invalid taint: %s", in)
	}

	return upcloud.KubernetesTaint{
		Effect: upcloud.KubernetesClusterTaintEffect(s[3]),
		Key:    s[1],
		Value:  s[2],
	}, nil
}

func ProcessNodeGroupParams(p CreateNodeGroupParams) (request.KubernetesNodeGroup, error) {
	ng := request.KubernetesNodeGroup{}

	kubeletArgs := make([]upcloud.KubernetesKubeletArg, 0)
	for _, v := range p.KubeletArgs {
		ka, err := processKubeletArg(v)
		if err != nil {
			return ng, err
		}

		kubeletArgs = append(kubeletArgs, ka)
	}

	labelSlice, err := labels.StringsToSliceOfLabels(p.Labels)
	if err != nil {
		return ng, err
	}

	sshKeys, err := commands.ParseSSHKeys(p.SSHKeys)
	if err != nil {
		return ng, err
	}

	taints := make([]upcloud.KubernetesTaint, 0)
	for _, v := range p.Taints {
		t, err := processTaint(v)
		if err != nil {
			return ng, err
		}

		taints = append(taints, t)
	}

	ng = request.KubernetesNodeGroup{
		Count:                p.Count,
		Labels:               labelSlice,
		Name:                 p.Name,
		Plan:                 p.Plan,
		SSHKeys:              sshKeys,
		Storage:              p.Storage,
		KubeletArgs:          kubeletArgs,
		Taints:               taints,
		UtilityNetworkAccess: upcloud.BoolPtr(p.UtilityNetworkAccess.Value()),
	}

	return ng, nil
}

type createCommand struct {
	*commands.BaseCommand
	p CreateNodeGroupParams
	completion.Kubernetes
	resolver.CachingKubernetes
}

// CreateCommand creates the "kubernetes nodegroup create" command
func CreateCommand() commands.Command {
	return &createCommand{
		BaseCommand: commands.New(
			"create",
			"Create a new node group into the specified cluster.",
			"upctl kubernetes nodegroup create 55199a44-4751-4e27-9394-7c7661910be3 --name secondary-node-group --count 3 --plan 2xCPU-4GB",
		),
	}
}

// InitCommand implements Command.InitCommand
func (s *createCommand) InitCommand() {
	fs := GetCreateNodeGroupFlagSet(&s.p)
	s.AddFlags(fs)

	commands.Must(s.Cobra().MarkFlagRequired("name"))
	commands.Must(s.Cobra().MarkFlagRequired("count"))
	commands.Must(s.Cobra().MarkFlagRequired("plan"))
}

func (s *createCommand) InitCommandWithConfig(cfg *config.Config) {
	commands.Must(s.Cobra().RegisterFlagCompletionFunc("plan", namedargs.CompletionFunc(completion.KubernetesPlan{}, cfg)))
}

// ExecuteSingleArgument implements commands.SingleArgumentCommand
func (s *createCommand) ExecuteSingleArgument(exec commands.Executor, arg string) (output.Output, error) {
	msg := fmt.Sprintf("Creating node group %s into cluster %v", s.p.Name, arg)
	exec.PushProgressStarted(msg)

	ng, err := ProcessNodeGroupParams(s.p)
	if err != nil {
		return commands.HandleError(exec, msg, err)
	}

	res, err := exec.All().CreateKubernetesNodeGroup(exec.Context(), &request.CreateKubernetesNodeGroupRequest{ClusterUUID: arg, NodeGroup: ng})
	if err != nil {
		return commands.HandleError(exec, msg, err)
	}

	exec.PushProgressSuccess(msg)

	return output.OnlyMarshaled{Value: res}, nil
}
