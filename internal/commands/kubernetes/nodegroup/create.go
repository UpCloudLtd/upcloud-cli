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

const (
	CloudNativePlanPrefix = "CLOUDNATIVE-"
	GPUPlanPrefix         = "GPU-"
)

var validStorageTiers = []string{upcloud.StorageTierMaxIOPS, upcloud.StorageTierStandard, upcloud.StorageTierHDD}

type CreateNodeGroupParams struct {
	Count                int
	Name                 string
	Plan                 string
	SSHKeys              []string
	Storage              string
	StorageSize          int
	StorageTier          string
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
	fs.IntVar(&p.StorageSize, "storage-size", 0, fmt.Sprintf("Custom storage size in GiB. Only applicable for Cloud Native (%s*) and GPU (%s*) plans. If not specified, uses plan default.", CloudNativePlanPrefix, GPUPlanPrefix))
	fs.StringVar(&p.StorageTier, "storage-tier", "", fmt.Sprintf("Storage tier (maxiops, standard, hdd). Only applicable for Cloud Native (%s*) and GPU (%s*) plans. If not specified, uses plan default.", CloudNativePlanPrefix, GPUPlanPrefix))
	fs.StringArrayVar(&p.Taints, "taint", []string{}, "Taints to be configured to the nodes in `key=value:effect` format")
	config.AddEnableOrDisableFlag(fs, &p.UtilityNetworkAccess, true, "utility-network-access", "utility network access. If disabled, nodes in this group will not have access to utility network")

	commands.Must(fs.SetAnnotation("count", commands.FlagAnnotationNoFileCompletions, nil))
	commands.Must(fs.SetAnnotation("kubelet-arg", commands.FlagAnnotationNoFileCompletions, nil))
	commands.Must(fs.SetAnnotation("label", commands.FlagAnnotationNoFileCompletions, nil))
	commands.Must(fs.SetAnnotation("name", commands.FlagAnnotationNoFileCompletions, nil))
	commands.Must(fs.SetAnnotation("ssh-key", commands.FlagAnnotationNoFileCompletions, nil))
	commands.Must(fs.SetAnnotation("storage", commands.FlagAnnotationNoFileCompletions, nil))
	commands.Must(fs.SetAnnotation("storage-size", commands.FlagAnnotationNoFileCompletions, nil))
	commands.Must(fs.SetAnnotation("storage-tier", commands.FlagAnnotationFixedCompletions, []string{upcloud.StorageTierMaxIOPS, upcloud.StorageTierStandard, upcloud.StorageTierHDD}))
	commands.Must(fs.SetAnnotation("taint", commands.FlagAnnotationNoFileCompletions, nil))

	return fs
}

// supportStorageCustomization checks if a plan supports storage customization
func supportStorageCustomization(planName string) bool {
	return strings.HasPrefix(planName, CloudNativePlanPrefix) ||
		strings.HasPrefix(planName, GPUPlanPrefix)
}

// validateStorageTier checks if the storage tier is valid
func validateStorageTier(tier string) error {
	if tier == "" {
		return nil // Empty is valid (uses plan default)
	}

	for _, validTier := range validStorageTiers {
		if tier == validTier {
			return nil
		}
	}

	return fmt.Errorf("invalid storage tier %q, must be one of: %s", tier, strings.Join(validStorageTiers, ", "))
}

// validateStorageSize checks if the storage size is valid
func validateStorageSize(size int) error {
	if size == 0 {
		return nil // Zero is valid (uses plan default)
	}

	if size < 25 {
		return fmt.Errorf("storage size must be at least 25 GiB, got %d", size)
	}

	if size > 4096 {
		return fmt.Errorf("storage size cannot exceed 4096 GiB, got %d", size)
	}

	return nil
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

	hasStorageCustomization := p.StorageSize > 0 || p.StorageTier != ""
	if hasStorageCustomization && !supportStorageCustomization(p.Plan) {
		return ng, fmt.Errorf("storage customization (--storage-size, --storage-tier) is only supported for Cloud Native (%s*) and GPU (%s*) plans, got plan: %s", CloudNativePlanPrefix, GPUPlanPrefix, p.Plan)
	}

	if err := validateStorageTier(p.StorageTier); err != nil {
		return ng, err
	}

	if err := validateStorageSize(p.StorageSize); err != nil {
		return ng, err
	}

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

	// Set storage customization for supported plans
	if hasStorageCustomization {
		if strings.HasPrefix(p.Plan, CloudNativePlanPrefix) {
			ng.CloudNativePlan = &upcloud.KubernetesNodeGroupCloudNativePlan{}
			if p.StorageSize > 0 {
				ng.CloudNativePlan.StorageSize = p.StorageSize
			}
			if p.StorageTier != "" {
				ng.CloudNativePlan.StorageTier = upcloud.StorageTier(p.StorageTier)
			}
		} else if strings.HasPrefix(p.Plan, GPUPlanPrefix) {
			ng.GPUPlan = &upcloud.KubernetesNodeGroupGPUPlan{}
			if p.StorageSize > 0 {
				ng.GPUPlan.StorageSize = p.StorageSize
			}
			if p.StorageTier != "" {
				ng.GPUPlan.StorageTier = upcloud.StorageTier(p.StorageTier)
			}
		}
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
			"upctl kubernetes nodegroup create 55199a44-4751-4e27-9394-7c7661910be3 --name gpu-nodes --count 2 --plan GPU-8xCPU-64GB-1xL40S --storage-size 1024 --storage-tier maxiops --label gpu=NVIDIA-L40S",
			"upctl kubernetes nodegroup create 55199a44-4751-4e27-9394-7c7661910be3 --name cloud-native-nodes --count 4 --plan CLOUDNATIVE-4xCPU-8GB --storage-size 50 --storage-tier standard",
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
