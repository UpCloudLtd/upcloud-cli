package stack

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands/stack/stackops"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands/stack/supabase"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands/stack/supabase/supabasechart"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands/stack/supabaseconfig"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/output"
	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud/request"
	"github.com/spf13/pflag"
)

func DeploySupabaseCommand() commands.Command {
	return &deploySupabaseCommand{
		BaseCommand: commands.New(
			"supabase",
			"Deploy a Supabase stack",
			"upctl stack deploy supabase <project-name>",
			"upctl stack deploy supabase my-new-project",
		),
	}
}

type deploySupabaseCommand struct {
	*commands.BaseCommand
	zone       string
	name       string
	configPath string
}

func (s *deploySupabaseCommand) InitCommand() {
	fs := &pflag.FlagSet{}
	fs.StringVar(&s.zone, "zone", s.zone, "Zone for the stack deployment")
	fs.StringVar(&s.name, "name", s.name, "Supabase stack name")
	fs.StringVar(&s.configPath, "configPath", s.configPath, "Optional path to a configuration file for the Supabase stack")
	s.AddFlags(fs)

	commands.Must(s.Cobra().MarkFlagRequired("zone"))
	commands.Must(s.Cobra().MarkFlagRequired("name"))

	// configPath is optional, but if provided, it should be a valid path
	if s.configPath != "" {
		if _, err := os.Stat(s.configPath); os.IsNotExist(err) {
			commands.Must(s.Cobra().MarkFlagRequired("configPath"))
		}
	}
}

func (s *deploySupabaseCommand) ExecuteWithoutArguments(exec commands.Executor) (output.Output, error) {
	msg := fmt.Sprintf("Creating supabase stack %v", s.name)
	//exec.PushProgressStarted(msg)

	// Command implementation for deploying a Supabase stack
	config, err := deploy(s.zone, s.name, s.configPath, exec)
	if err != nil {
		fmt.Printf("Error deploying Supabase stack: %+v\n", err)
		return commands.HandleError(exec, msg, err)
	}

	exec.PushProgressSuccess("Supabase stack created successfully")

	return output.Raw(summaryOutput(config)), nil
}

func deploy(location, name, configPath string, exec commands.Executor) (*supabaseconfig.SupabaseConfig, error) {
	fmt.Printf("Deploying Supabase stack '%s' in location '%s'\n", name, location)

	clusterName := fmt.Sprintf("stack-supabase-cluster-%s-%s", name, location)
	networkName := fmt.Sprintf("stack-supabase-net-%s-%s", name, location)
	var network *upcloud.Network

	clusters, err := exec.All().GetKubernetesClusters(exec.Context(), &request.GetKubernetesClustersRequest{})

	if err != nil {
		return nil, fmt.Errorf("failed to get Kubernetes clusters: %w", err)
	}

	// TODO: We need to include the upgrade option when the cluster already exists
	// For now we are not supporting upgrades or updates
	if stackops.ClusterExists(clusterName, clusters) {
		return nil, fmt.Errorf("a cluster with the name '%s' already exists", clusterName)
	}

	fmt.Println("Creating Kubernetes cluster:", clusterName)
	// Check if the network already exists
	networks, err := exec.Network().GetNetworks(exec.Context())
	if err != nil {
		return nil, fmt.Errorf("failed to get networks: %w", err)
	}

	network = stackops.GetNetworkFromName(networkName, networks.Networks)

	// Create the network if it does not exist
	if network == nil {
		fmt.Println("Network does NOT exist, creating:", networkName)
		network, err = stackops.CreateNetwork(exec, networkName, location)
		if err != nil {
			return nil, fmt.Errorf("failed to create network: %w", err)
		}

		if network == nil {
			return nil, fmt.Errorf("created network %s is nil", networkName)
		}
		if len(network.IPNetworks) == 0 {
			return nil, fmt.Errorf("created network %s has no IP networks", networkName)
		}
	}

	fmt.Println(("Creating Kubernetes cluster..."))
	cluster, err := exec.All().CreateKubernetesCluster(exec.Context(), &request.CreateKubernetesClusterRequest{
		Name:        clusterName,
		Network:     network.UUID,
		Zone:        location,
		NetworkCIDR: network.IPNetworks[0].Address,
		Labels: []upcloud.Label{
			{Key: "stacks.upcloud.com/stack", Value: "supabase"},
			{Key: "stacks.upcloud.com/created-by", Value: "upctl"},
			{Key: "stacks.upcloud.com/chart-version", Value: "0.1.3"},
			{Key: "stacks.upcloud.com/name", Value: clusterName},
		},
		NodeGroups: []request.KubernetesNodeGroup{
			{
				Name:  "supabase-node-group",
				Count: 1,
				Plan:  "2xCPU-4GB",
			},
		},
	})

	if err != nil {
		return nil, fmt.Errorf("failed to create Kubernetes cluster: %w", err)
	}

	fmt.Println("Kubernetes cluster created successfully:", cluster.Name, "with UUID", cluster.UUID)

	fmt.Println("Requesting kubeconfig for the cluster...")
	kubeconfig, err := exec.All().GetKubernetesKubeconfig(exec.Context(), &request.GetKubernetesKubeconfigRequest{
		UUID: cluster.UUID,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to get kubeconfig for cluster %s: %w", cluster.UUID, err)
	}

	fmt.Println("Kubeconfig received successfully for cluster:", cluster.Name)

	// Create a tmp dir for this deployment
	chartDir, err := os.MkdirTemp("", fmt.Sprintf("supabase-%s-%s", name, location))
	if err != nil {
		return nil, fmt.Errorf("failed to make temp dir for deployment: %w", err)
	}

	// clean up at the end
	//defer os.RemoveAll(chartDir)

	// Save kubeconfig to temp file
	kubeconfigPath := filepath.Join(chartDir, "kubeconfig.yaml")
	if err := os.WriteFile(kubeconfigPath, []byte(kubeconfig), 0600); err != nil {
		return nil, fmt.Errorf("failed to write kubeconfig: %w", err)
	}

	// Create a Kubernetes client from the kubeconfig
	os.Setenv("KUBECONFIG", kubeconfigPath)
	kubeClient, err := stackops.GetKubernetesClient(kubeconfigPath)
	if err != nil {
		return nil, fmt.Errorf("failed to create Kubernetes client: %w", err)
	}

	fmt.Println("Kubeconfig received successfully and saved to:", kubeconfigPath)

	// Generate configuration for the Supabase stack from input config file
	config, err := supabaseconfig.Generate(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to generate secrets: %w", err)
	}

	valuesPath := filepath.Join(chartDir, "charts/supabase/values.example.yaml")
	securePath := filepath.Join(chartDir, fmt.Sprintf("supabase-%s-values.secure.yaml", name))
	updatedValuesFile := filepath.Join(chartDir, fmt.Sprintf("supabase-%s-values.updated.yaml", name))
	err = supabaseconfig.WriteSecureValues(securePath, config)
	if err != nil {
		return nil, fmt.Errorf("failed to write secrets file: %w", err)
	}

	// unpack the supabase charts into that temp dir
	if err := stackops.ExtractChart(supabasechart.ChartFS, chartDir); err != nil {
		return nil, fmt.Errorf("failed to extract supabase chart: %w", err)
	}

	// Deploy the Helm release
	err = supabase.DeployHelmRelease(kubeClient, clusterName, filepath.Join(chartDir, "charts/supabase"), []string{valuesPath, securePath}, false)
	if err != nil {
		return nil, fmt.Errorf("failed to deploy Supabase Helm release: %w", err)
	}

	fmt.Println("Supabase Helm release deployed successfully")

	fmt.Println("Waiting on Supabase Kong Load Balancer to be ready...")
	// Wait for the Supabase stack to be ready
	lbHostname, err := stackops.WaitForLoadBalancer(kubeClient, clusterName, clusterName+"-supabase-kong", 60, 20*time.Second)

	// Adding the cluster name and lbHostname to the config for reporting purposes
	config.LbHostname = lbHostname
	config.ClusterName = clusterName

	if err != nil {
		return nil, fmt.Errorf("failed to wait for Supabase Kong Load Balancer: %w", err)
	}

	fmt.Println("Supabase Kong Load Balancer ", lbHostname, " is ready")

	// Update the DNS entries in the values file with the provided dnsPrefix
	fmt.Println("Updating deployment DNS entries with prefix:", lbHostname)

	dnsPrefix := strings.TrimSuffix(lbHostname, ".upcloudlb.com")

	err = supabase.UpdateDNS(valuesPath, updatedValuesFile, dnsPrefix)
	if err != nil {
		return nil, fmt.Errorf("failed to update DNS entries in values file: %w", err)
	}

	files := []string{
		updatedValuesFile,
		securePath,
	}

	supabase.DeployHelmRelease(kubeClient, clusterName, filepath.Join(chartDir, "charts/supabase"), files, true)

	fmt.Println("Supabase stack deployed successfully with updated DNS entries")

	return config, nil
}

func summaryOutput(config *supabaseconfig.SupabaseConfig) []byte {
	builder := &strings.Builder{}

	fmt.Fprintln(builder, "Supabase deployed successfully!")
	fmt.Fprintf(builder, "Public endpoint:      http://%s:8000\n", config.LbHostname)
	fmt.Fprintf(builder, "Namespace:            %s\n", config.ClusterName)
	fmt.Fprintf(builder, "ANON_KEY:             %s\n", config.AnonKey)
	fmt.Fprintf(builder, "SERVICE_ROLE_KEY:     %s\n", config.ServiceRoleKey)
	fmt.Fprintf(builder, "POSTGRES_PASSWORD:    %s\n", orNotSet(config.PostgresPassword))
	fmt.Fprintf(builder, "POOLER_TENANT_ID:     %s\n", config.PoolerTenantID)
	fmt.Fprintf(builder, "DASHBOARD_USERNAME:   %s\n", config.DashboardUsername)
	fmt.Fprintf(builder, "DASHBOARD_PASSWORD:   %s\n", config.DashboardPassword)
	fmt.Fprintf(builder, "S3 ENABLED:           %t\n", config.S3Enabled)
	fmt.Fprintf(builder, "S3_BUCKET:            %s\n", orNotSet(config.S3BucketName))
	fmt.Fprintf(builder, "S3_ENDPOINT:          %s\n", orNotSet(config.S3Endpoint))
	fmt.Fprintf(builder, "S3_REGION:            %s\n", orNotSet(config.S3Region))
	fmt.Fprintf(builder, "SMTP ENABLED:         %t\n", config.SmtpEnabled)
	fmt.Fprintf(builder, "SMTP_HOST:           %s\n", orNotSet(config.SmtpHost))
	fmt.Fprintf(builder, "SMTP_PORT:            %s\n", orNotSet(config.SmtpPort))
	fmt.Fprintf(builder, "SMTP_USER:            %s\n", orNotSet(config.SmtpUsername))
	fmt.Fprintf(builder, "SMTP_SENDER_NAME:     %s\n", orNotSet(config.SmtpSenderName))

	return []byte(builder.String())
}

func orNotSet(val string) string {
	if val == "" {
		return "not set"
	}
	return val
}
