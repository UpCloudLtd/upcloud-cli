package stack

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands/stack/stackops"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/output"
	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud/request"
	"github.com/spf13/pflag"
)

//go:embed dokku/config/**
var dokkuChartFS embed.FS

func DeployDokkuCommand() commands.Command {
	return &deployDokkuCommand{
		BaseCommand: commands.New(
			"dokku",
			"Deploy a Dokku Builder stack",
			"upctl stack deploy dokku <project-name>",
			"upctl stack deploy dokku my-new-project",
		),
	}
}

type deployDokkuCommand struct {
	*commands.BaseCommand
	zone             string
	name             string
	githubPAT        string
	githubUser       string
	certManagerEmail string
	globalDomain     string
	numNodes         int
	sshPath          string
	sshPubPath       string
	githubPackageUrl string
}

func getHomeDir() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		// Handle error as you prefer, e.g., log it or return a default value
		return ""
	}
	return homeDir
}

func (s *deployDokkuCommand) InitCommand() {
	defaultSSH := filepath.Join(getHomeDir(), ".ssh", "id_rsa")
	defaultSSHPub := filepath.Join(getHomeDir(), ".ssh", "id_rsa.pub")

	fs := &pflag.FlagSet{}
	fs.StringVar(&s.zone, "zone", s.zone, "Zone (location) for the stack deployment")
	fs.StringVar(&s.name, "name", s.name, "Dokku project name")
	fs.StringVar(&s.githubPAT, "github-pat", s.githubPAT, "GitHub Personal Access Token. Used to allow Dokku to push your app images to your GitHub Container Registry. Make sure it has write:packages and read:packages permissions")
	fs.StringVar(&s.githubUser, "github-user", s.githubUser, "Used to allow Dokku to push your app images to your GitHub Container Registry")
	fs.StringVar(&s.certManagerEmail, "cert-manager-email", "ops@example.com", "Email for TLS cert registration (default: ops@example.com)")
	fs.StringVar(&s.globalDomain, "global-domain", s.globalDomain, "Example: example.com. If you do not have a domain name leave this empty and it will get the value of the ingress nginx load balancer automatically. Example: lb-0a39e6584…")
	fs.IntVar(&s.numNodes, "num-nodes", 3, "Number of nodes in the Dokku cluster (default: 3)")
	fs.StringVar(&s.sshPath, "ssh-path", defaultSSH, "Path to your private SSH key (default: ~/.ssh/id_rsa). Needed to be able to ‘git push dokku@<host>:<app>’ when deploying apps with git push")
	fs.StringVar(&s.sshPubPath, "ssh-path-pub", defaultSSHPub, "Path to your public SSH key (default: ~/.ssh/id_rsa.pub)")
	fs.StringVar(&s.githubPackageUrl, "github-package-url", "ghcr.io", "Container registry hostname (default: ghcr.io)")
	s.AddFlags(fs)

	commands.Must(s.Cobra().MarkFlagRequired("zone"))
	commands.Must(s.Cobra().MarkFlagRequired("name"))
	commands.Must(s.Cobra().MarkFlagRequired("github-pat"))
	commands.Must(s.Cobra().MarkFlagRequired("github-user"))
}

func (s *deployDokkuCommand) ExecuteWithoutArguments(exec commands.Executor) (output.Output, error) {
	msg := fmt.Sprintf("Creating dokku stack %v", s.name)
	exec.PushProgressStarted(msg)

	// Create a tmp dir for this deployment
	configDir, err := os.MkdirTemp("", fmt.Sprintf("dokku-%s-%s", s.name, s.zone))
	if err != nil {
		return nil, fmt.Errorf("failed to make temp dir for deployment: %w", err)
	}

	//defer os.RemoveAll(configDir)

	// unpack the dokku charts and config files into that temp dir
	if err := stackops.ExtractChart(dokkuChartFS, configDir); err != nil {
		return nil, fmt.Errorf("failed to extract dokku charts and configuration files: %w", err)
	}

	if err = s.deploy(exec, configDir); err != nil {
		return nil, fmt.Errorf("failed to deploy dokku stack: %w", err)
	}

	exec.PushProgressSuccess(msg)

	return output.Raw([]byte("Command executed successfully")), nil
}

func (s *deployDokkuCommand) deploy(exec commands.Executor, configDir string) error {
	clusterName := fmt.Sprintf("stack-dokku-cluster-%s-%s", s.name, s.zone)
	networkName := fmt.Sprintf("stack-dokku-net-%s-%s", s.name, s.zone)
	var network *upcloud.Network

	msg := fmt.Sprintf("Setting up kubernetes cluster:%s in location '%s'\n", clusterName, s.zone)
	exec.PushProgressStarted(msg)

	// Check if the cluster already exists
	clusters, err := exec.All().GetKubernetesClusters(exec.Context(), &request.GetKubernetesClustersRequest{})

	if err != nil {
		return fmt.Errorf("failed to get Kubernetes clusters: %w", err)
	}

	// Return if the cluster already exists
	if stackops.ClusterExists(clusterName, clusters) {
		return fmt.Errorf("a cluster with the name '%s' already exists", clusterName)
	}

	exec.PushProgressUpdateMessage(msg, "Creating Kubernetes cluster...")

	// Check if the network already exists
	networks, err := exec.Network().GetNetworks(exec.Context())
	if err != nil {
		return fmt.Errorf("failed to get networks: %w", err)
	}

	network = stackops.GetNetworkFromName(networkName, networks.Networks)

	// Create the network if it does not exist
	if network == nil {
		exec.PushProgressUpdateMessage(msg, fmt.Sprintf("Network does NOT exist, creating: %s", networkName))
		network, err = stackops.CreateNetwork(exec, networkName, s.zone)
		if err != nil {
			return fmt.Errorf("failed to create network: %w", err)
		}

		if network == nil {
			return fmt.Errorf("created network %s is nil", networkName)
		}
		if len(network.IPNetworks) == 0 {
			return fmt.Errorf("created network %s has no IP networks", networkName)
		}
	}

	cluster, err := exec.All().CreateKubernetesCluster(exec.Context(), &request.CreateKubernetesClusterRequest{
		Name:        clusterName,
		Network:     network.UUID,
		Zone:        s.zone,
		NetworkCIDR: network.IPNetworks[0].Address,
		Labels: []upcloud.Label{
			{Key: "stacks.upcloud.com/stack", Value: "dokku"},
			{Key: "stacks.upcloud.com/created-by", Value: "upctl"},
			{Key: "stacks.upcloud.com/dokku-version", Value: "0.1.3"},
			{Key: "stacks.upcloud.com/version", Value: "1.0.0"},
			{Key: "stacks.upcloud.com/name", Value: clusterName},
		},
		NodeGroups: []request.KubernetesNodeGroup{
			{
				Name:  "dokku-node-group",
				Count: s.numNodes,
				Plan:  "2xCPU-4GB",
			},
		},
	})

	if err != nil {
		return fmt.Errorf("failed to create Kubernetes cluster: %w", err)
	}

	// Get kubeconfig file for the cluster
	kubeconfigPath, err := writeKubeconfigToFile(exec, cluster.UUID, configDir)
	if err != nil {
		return fmt.Errorf("failed to write kubeconfig to file: %w", err)
	}

	// Create a Kubernetes client from the kubeconfig
	os.Setenv("KUBECONFIG", kubeconfigPath)
	kubeClient, err := stackops.GetKubernetesClient(kubeconfigPath)
	if err != nil {
		return fmt.Errorf("failed to create Kubernetes client: %w", err)
	}

	// Wait for the Kubernetes API server to be ready
	stackops.WaitForAPIServer(kubeClient)

	// Deploy nginx ingress controller
	ingressValuesFilePath := filepath.Join(configDir, "dokku/config/ingress/values.yaml")
	err = stackops.DeployHelmReleaseFromRepo(
		kubeClient,
		"ingress-nginx",
		"https://kubernetes.github.io/ingress-nginx",
		"ingress-nginx",
		"4.13.0",
		[]string{ingressValuesFilePath},
		false)

	if err != nil {
		return fmt.Errorf("failed to deploy ingress-nginx: %w", err)
	}

	// Wait for the ingress controller to be ready
	lbHostname, err := stackops.WaitForLoadBalancer(kubeClient, "ingress-nginx", "ingress-nginx-controller", 30, 10*time.Second)
	if err != nil {
		return fmt.Errorf("failed to wait for ingress-nginx load balancer: %w", err)
	}

	// If no global domain is provided, use the load balancer hostname
	if s.globalDomain == "" {
		s.globalDomain = lbHostname
	}

	// Deploy cert-manager
	override, err := os.CreateTemp("", "certmgr-override-*.yaml")
	if err != nil {
		return err
	}
	defer os.Remove(override.Name())

	override.WriteString("installCRDs: true\n")
	override.Close()

	err = stackops.DeployHelmReleaseFromRepo(
		kubeClient,
		"cert-manager",
		"https://charts.jetstack.io",
		"cert-manager",
		"v1.11.0",
		[]string{override.Name()},
		false,
	)

	if err != nil {
		return fmt.Errorf("failed to deploy cert-manager: %w", err)
	}

	// Deploy Dokku
	dokkuConfigPath := filepath.Join(configDir, "dokku/config/dokku")
	err = stackops.ApplyKustomize(dokkuConfigPath)
	if err != nil {
		return fmt.Errorf("failed to apply Dokku kustomize configuration: %w", err)
	}

	// Configure Dokku with the provided parameters
	lbHostName, nodeIp, err := stackops.ConfigureDokku(
		kubeClient,
		s.sshPath,
		s.sshPubPath,
		s.githubPAT,
		s.githubUser,
		s.githubPackageUrl,
		s.globalDomain,
		s.certManagerEmail)

	if err != nil {
		return fmt.Errorf("failed to configure Dokku: %w", err)
	}

	// Print final instructions for the user
	stackops.PrintFinalInstructions(kubeconfigPath, s.globalDomain, s.sshPath, lbHostName, nodeIp)

	//msg = fmt.Sprintln("Kubernetes cluster created successfully:", cluster.Name, "with UUID", cluster.UUID)
	exec.PushProgressSuccess(msg)

	return nil
}

// writeKubeconfigToFile retrieves the kubeconfig for the given cluster and writes it to a file
func writeKubeconfigToFile(exec commands.Executor, clusterId string, configDir string) (string, error) {
	kubeconfig, err := exec.All().GetKubernetesKubeconfig(exec.Context(), &request.GetKubernetesKubeconfigRequest{
		UUID: clusterId,
	})

	if err != nil {
		return "", fmt.Errorf("failed to get kubeconfig for cluster %s: %w", clusterId, err)
	}

	kubeconfigPath := filepath.Join(configDir, "kubeconfig.yaml")
	if err := os.WriteFile(kubeconfigPath, []byte(kubeconfig), 0600); err != nil {
		return "", fmt.Errorf("failed to write kubeconfig: %w", err)
	}

	return kubeconfigPath, nil
}
