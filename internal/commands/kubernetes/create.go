package kubernetes

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/UpCloudLtd/upcloud-cli/v2/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v2/internal/completion"
	"github.com/UpCloudLtd/upcloud-cli/v2/internal/output"
	"github.com/UpCloudLtd/upcloud-cli/v2/internal/ui"

	"github.com/UpCloudLtd/upcloud-go-api/v4/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/v4/upcloud/request"
	"github.com/spf13/pflag"
)

// CreateCommand creates the "kubernetes create" command
func CreateCommand() commands.Command {
	return &createCommand{
		BaseCommand: commands.New(
			"create",
			"Create a Kubernetes cluster",
			`upctl kubernetes create \
				--name my-cluster \
				--network 03e5ca07-f36c-4957-a676-e001e40441eb \
				--network-cidr "172.16.1.0/24" \
				--node-group count=2,name=my-minimal-node-group,plan=K8S-2xCPU-4GB, \
				--zone de-fra1`,
			`upctl kubernetes create \
				--name my-cluster \
				--network 03e5ca07-f36c-4957-a676-e001e40441eb \
				--network-cidr "172.16.1.0/24" \
				--node-group count=4,kubelet-arg="log-flush-frequency=5s",label="owner=devteam",label="env=dev",name=my-node-group,plan=K8S-4xCPU-8GB,ssh-key="ssh-ed25519 AAAAo admin@user.com",ssh-key="/path/to/your/public/ssh/key.pub",storage=01000000-0000-4000-8000-000160010100,taint="env=dev:NoSchedule",taint="env=dev2:NoSchedule" \
				--zone de-fra1`,
		),
	}
}

type createParams struct {
	request.CreateKubernetesClusterRequest
	nodeGroups []string
}

type createNodeGroupParams struct {
	count       int
	name        string
	plan        string
	sshKeys     []string
	storage     string
	kubeletArgs []string
	labels      []string
	taints      []string
}

func (p *createParams) processParams() error {
	ngs := make([]upcloud.KubernetesNodeGroup, 0)

	for _, v := range p.nodeGroups {
		ng, err := processNodeGroup(v)
		if err != nil {
			return err
		}
		ngs = append(ngs, ng)
	}
	p.NodeGroups = ngs

	return nil
}

func processNodeGroup(in string) (upcloud.KubernetesNodeGroup, error) {
	fs := &pflag.FlagSet{}
	ng := upcloud.KubernetesNodeGroup{}

	args, err := commands.ParseN(in, 2)
	if err != nil {
		return ng, err
	}

	p := createNodeGroupParams{}
	fs.IntVar(&p.count, "count", 0, "")
	fs.StringArrayVar(&p.kubeletArgs, "kubelet-arg", []string{}, "")
	fs.StringArrayVar(&p.labels, "label", []string{}, "")
	fs.StringVar(&p.name, "name", "", "")
	fs.StringVar(&p.plan, "plan", "", "")
	fs.StringArrayVar(&p.sshKeys, "ssh-key", []string{}, "")
	fs.StringVar(&p.storage, "storage", "", "")
	fs.StringArrayVar(&p.taints, "taint", []string{}, "")

	err = fs.Parse(args)

	if err != nil {
		return ng, err
	}

	kubeletArgs := make([]upcloud.KubernetesKubeletArg, 0)
	for _, v := range p.kubeletArgs {
		ka, err := processKubeletArg(v)
		if err != nil {
			return ng, err
		}

		kubeletArgs = append(kubeletArgs, ka)
	}

	labels := make([]upcloud.Label, 0)
	for _, v := range p.labels {
		l, err := processLabel(v)
		if err != nil {
			return ng, err
		}

		labels = append(labels, l)
	}

	sshKeys, err := commands.ParseSSHKeys(p.sshKeys)
	if err != nil {
		return ng, err
	}

	taints := make([]upcloud.KubernetesTaint, 0)
	for _, v := range p.taints {
		t, err := processTaint(v)
		if err != nil {
			return ng, err
		}

		taints = append(taints, t)
	}

	ng = upcloud.KubernetesNodeGroup{
		Count:       p.count,
		Labels:      labels,
		Name:        p.name,
		Plan:        p.plan,
		SSHKeys:     sshKeys,
		Storage:     p.storage,
		KubeletArgs: kubeletArgs,
		Taints:      taints,
	}

	return ng, nil
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

type createCommand struct {
	*commands.BaseCommand
	params createParams
	completion.Kubernetes
}

// InitCommand implements Command.InitCommand
func (c *createCommand) InitCommand() {
	fs := &pflag.FlagSet{}
	c.params = createParams{CreateKubernetesClusterRequest: request.CreateKubernetesClusterRequest{}}

	fs.StringVar(&c.params.Name, "name", "", "Kubernetes cluster name.")
	fs.StringVar(&c.params.Network, "network", "", "Network to use. The value should be UUID of a private network.")
	fs.StringVar(&c.params.NetworkCIDR, "network-cidr", "", "CIDR of the network being used.")
	fs.StringArrayVar(
		&c.params.nodeGroups,
		"node-group",
		[]string{},
		"Node group(s) for running workloads, multiple can be declared.\n"+
			"Usage: --node-group "+
			"count=8,"+
			"kubelet-arg=\"log-flush-frequency=5s\","+
			"label=\"owner=devteam\","+
			"label=\"env=dev\","+
			"name=my-node-group,"+
			"plan=K8S-2xCPU-4GB,"+
			"ssh-key=\"ssh-ed25519 AAAAo admin@user.com\","+
			"ssh-key=\"/path/to/your/public/ssh/key.pub\","+
			"storage=01000000-0000-4000-8000-000160010100,"+
			"taint=\"env=dev:NoSchedule\","+
			"taint=\"env=dev2:NoSchedule\"",
	)
	fs.StringVar(&c.params.Zone, "zone", "", "Zone where to create the server.")
	c.AddFlags(fs)

	_ = c.Cobra().MarkFlagRequired("name")
	_ = c.Cobra().MarkFlagRequired("network")
	_ = c.Cobra().MarkFlagRequired("network-cidr")
	_ = c.Cobra().MarkFlagRequired("node-group")
	_ = c.Cobra().MarkFlagRequired("zone")
}

// ExecuteWithoutArguments implements commands.NoArgumentCommand
func (c *createCommand) ExecuteWithoutArguments(exec commands.Executor) (output.Output, error) {
	svc := exec.All()
	msg := fmt.Sprintf("Creating cluster %s", c.params.Name)
	exec.PushProgressStarted(msg)
	if err := c.params.processParams(); err != nil {
		return nil, err
	}

	exec.Debug("sending r", "params", c.params)

	r := c.params.CreateKubernetesClusterRequest

	res, err := svc.CreateKubernetesCluster(&r)
	if err != nil {
		return commands.HandleError(exec, msg, err)
	}

	exec.PushProgressSuccess(msg)

	return output.MarshaledWithHumanDetails{Value: res, Details: []output.DetailRow{
		{Title: "UUID", Value: res.UUID, Colour: ui.DefaultUUUIDColours},
	}}, nil
}
