package nodegroup

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/UpCloudLtd/upcloud-cli/v2/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v2/internal/completion"
	"github.com/UpCloudLtd/upcloud-cli/v2/internal/output"
	"github.com/UpCloudLtd/upcloud-cli/v2/internal/resolver"
	"github.com/UpCloudLtd/upcloud-go-api/v5/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/v5/upcloud/request"
	"github.com/spf13/pflag"
)

type CreateNodeGroupParams struct {
	Count       int
	Name        string
	Plan        string
	SSHKeys     []string
	Storage     string
	KubeletArgs []string
	Labels      []string
	Taints      []string
}

func GetCreateNodeGroupFlagSet(p *CreateNodeGroupParams) *pflag.FlagSet {
	fs := &pflag.FlagSet{}

	fs.IntVar(&p.Count, "count", 0, "")
	fs.StringArrayVar(&p.KubeletArgs, "kubelet-arg", []string{}, "")
	fs.StringArrayVar(&p.Labels, "label", []string{}, "")
	fs.StringVar(&p.Name, "name", "", "")
	fs.StringVar(&p.Plan, "plan", "", "")
	fs.StringArrayVar(&p.SSHKeys, "ssh-key", []string{}, "")
	fs.StringVar(&p.Storage, "storage", "", "")
	fs.StringArrayVar(&p.Taints, "taint", []string{}, "")

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

func processLabel(in string) (upcloud.Label, error) {
	split := strings.SplitN(in, "=", 2)
	if len(split) < 2 {
		return upcloud.Label{}, fmt.Errorf("invalid label: %s", in)
	}

	return upcloud.Label{
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

func ProcessNodeGroupParams[T upcloud.KubernetesNodeGroup | request.KubernetesNodeGroup](p CreateNodeGroupParams) (T, error) {
	ng := T{}

	kubeletArgs := make([]upcloud.KubernetesKubeletArg, 0)
	for _, v := range p.KubeletArgs {
		ka, err := processKubeletArg(v)
		if err != nil {
			return ng, err
		}

		kubeletArgs = append(kubeletArgs, ka)
	}

	labels := make([]upcloud.Label, 0)
	for _, v := range p.Labels {
		l, err := processLabel(v)
		if err != nil {
			return ng, err
		}

		labels = append(labels, l)
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

	ng = T{
		Count:       p.Count,
		Labels:      labels,
		Name:        p.Name,
		Plan:        p.Plan,
		SSHKeys:     sshKeys,
		Storage:     p.Storage,
		KubeletArgs: kubeletArgs,
		Taints:      taints,
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

	_ = s.Cobra().MarkFlagRequired("name")
	_ = s.Cobra().MarkFlagRequired("count")
	_ = s.Cobra().MarkFlagRequired("plan")
}

// ExecuteSingleArgument implements commands.SingleArgumentCommand
func (s *createCommand) ExecuteSingleArgument(exec commands.Executor, arg string) (output.Output, error) {
	msg := fmt.Sprintf("Creating node group %s into cluster %v", s.p.Name, arg)
	exec.PushProgressStarted(msg)

	ng, err := ProcessNodeGroupParams[request.KubernetesNodeGroup](s.p)
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
