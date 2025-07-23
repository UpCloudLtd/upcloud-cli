package stack

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands/stack/supabase"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands/stack/supabase/supabasechart"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands/stack/supabaseconfig"
	stacksecrets "github.com/UpCloudLtd/upcloud-cli/v3/internal/commands/stack/supabaseconfig"
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
	location   string
	name       string
	configPath string
}

func (s *deploySupabaseCommand) InitCommand() {
	fs := &pflag.FlagSet{}
	fs.StringVar(&s.location, "location", s.location, "Select the location (region) for the stack deployment")
	fs.StringVar(&s.name, "name", s.name, "Specify the name of the Supabase project")
	fs.StringVar(&s.configPath, "configPath", s.configPath, "Optional path to a configuration file for the Supabase stack")
	s.AddFlags(fs)

	commands.Must(s.Cobra().MarkFlagRequired("location"))
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
	error := deploy(s.location, s.name, s.configPath, exec)
	if error != nil {
		fmt.Printf("Error deploying Supabase stack: %+v\n", error)
		return commands.HandleError(exec, msg, error)
	}

	exec.PushProgressSuccess("Supabase stack created successfully")

	return output.Raw([]byte("Commamnd executed successfully")), nil
}

func deploy(location, name, configPath string, exec commands.Executor) error {
	fmt.Printf("Deploying Supabase stack '%s' in location '%s'\n", name, location)

	clusterName := fmt.Sprintf("supabase-%s-%s", name, location)
	networkName := fmt.Sprintf("supabase-net-%s-%s", name, location)
	var networkID string
	var networkCIDR string

	clusters, err := exec.All().GetKubernetesClusters(exec.Context(), &request.GetKubernetesClustersRequest{})

	if err != nil {
		return fmt.Errorf("failed to get Kubernetes clusters: %w", err)
	}

	for _, cluster := range clusters {
		if cluster.Name == clusterName {
			// TODO: We need to include the upgrade option when the cluster already exists
			return fmt.Errorf("a cluster with the name '%s' already exists", clusterName)
		}
	}

	fmt.Println("Creating Kubernetes cluster:", clusterName)
	// Check if the network already exists
	networks, err := exec.Network().GetNetworks(exec.Context())
	if err != nil {
		return fmt.Errorf("failed to get networks: %w", err)
	}

	networkExists := false
	for _, network := range networks.Networks {
		if network.Name == networkName {
			networkExists = true
			networkID = network.UUID
			break
		}
	}

	// Create the network if it does not exist
	if !networkExists {
		fmt.Println("Network does NOT exists, creating:", networkName)
		var networkCreated bool = false

		for i := 0; i < 10; i++ {
			// Generate random 10.0.X.0/24 subnet
			x := rand.Intn(254) + 1
			cidr := fmt.Sprintf("10.0.%d.0/24", x)
			fmt.Println("Trying to create network with CIDR:", cidr)

			network, err := exec.Network().CreateNetwork(exec.Context(), &request.CreateNetworkRequest{
				Name:       networkName,
				Zone:       location,
				Router:     "",
				IPNetworks: []upcloud.IPNetwork{{Address: cidr, DHCP: 1, Family: upcloud.IPAddressFamilyIPv4}},
			})

			if err != nil {
				fmt.Printf("Failed to create network %s with CIDR %s: %v\n", networkName, cidr, err)
				continue
			} else {
				fmt.Println("Network created successfully:", network.Name, "with network.UUID", network.UUID)
				networkCreated = true
				networkID = network.UUID
				networkCIDR = cidr
				break
			}
		}
		if !networkCreated {
			return fmt.Errorf("failed to create network after 10 attempts")
		}
	}

	fmt.Println(("Moving to creating Kubernetes cluster..."))
	cluster, err := exec.All().CreateKubernetesCluster(exec.Context(), &request.CreateKubernetesClusterRequest{
		Name:        clusterName,
		Network:     networkID,
		Zone:        location,
		NetworkCIDR: networkCIDR,
		Labels: []upcloud.Label{
			{Key: "stacks.upcloud.com/stack", Value: "supabase"},
			{Key: "stacks.upcloud.com/created-by", Value: "upctl"},
			{Key: "stacks.upcloud.com/chart-version", Value: "0.1.3"},
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
		return fmt.Errorf("failed to create Kubernetes cluster: %w", err)
	}

	fmt.Println("Kubernetes cluster created successfully:", cluster.Name, "with UUID", cluster.UUID)

	fmt.Println("Requesting kubeconfig for the cluster...")
	kubeconfig, err := exec.All().GetKubernetesKubeconfig(exec.Context(), &request.GetKubernetesKubeconfigRequest{
		UUID: cluster.UUID,
	})

	if err != nil {
		return fmt.Errorf("failed to get kubeconfig for cluster %s: %w", cluster.UUID, err)
	}

	fmt.Println("Kubeconfig received successfully for cluster:", cluster.Name)

	// Create a tmp dir for this deployment
	chartDir, err := os.MkdirTemp("", fmt.Sprintf("supabase-%s-%s", name, location))
	if err != nil {
		return fmt.Errorf("failed to make temp dir for deployment: %w", err)
	}

	// Save kubeconfig to temp file
	kubeconfigPath := filepath.Join(chartDir, "kubeconfig.yaml")
	if err := os.WriteFile(kubeconfigPath, []byte(kubeconfig), 0600); err != nil {
		return fmt.Errorf("failed to write kubeconfig: %w", err)
	}

	// Point helm (client-go) to kubeconfig
	os.Setenv("KUBECONFIG", kubeconfigPath)

	fmt.Println("Kubeconfig received successfully and saved to:", kubeconfigPath)

	// TODO: Read deploy_supabase.env file and use it to set up the Supabase stack

	// Generate configuration for the Supabase stack
	config, err := supabaseconfig.Generate(configPath)
	if err != nil {
		return fmt.Errorf("failed to generate secrets: %w", err)
	}

	// Print the secrets to the console
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		fmt.Println("Failed to marshal secrets:", err)
	} else {
		fmt.Println(string(data))
	}

	// clean up at the end
	//defer os.RemoveAll(chartDir)

	valuesPath := filepath.Join(chartDir, "charts/supabase/values.example.yaml")
	securePath := filepath.Join(chartDir, fmt.Sprintf("supabase-%s-values.secure.yaml", name))
	updatedValuesFile := filepath.Join(chartDir, fmt.Sprintf("supabase-%s-values.updated.yaml", name))
	err = stacksecrets.WriteSecureValues(securePath, config)
	if err != nil {
		return fmt.Errorf("failed to write secrets file: %w", err)
	}

	// unpack the supabase charts into that temp dir
	if err := supabase.ExtractChart(supabasechart.ChartFS, chartDir); err != nil {
		return fmt.Errorf("failed to extract supabase chart: %w", err)
	}

	// Deploy the Helm release
	// FOR INSTALL:
	// VALUES_FILE="${CHART_DIR}/values.example.yaml"
	// SECURE_VALUES_FILE="${CHART_DIR}/values.secure.yaml" -> This is the file in chartDir
	err = supabase.DeployHelmRelease(clusterName, filepath.Join(chartDir, "charts/supabase"), []string{valuesPath, securePath}, false)
	if err != nil {
		return fmt.Errorf("failed to deploy Supabase Helm release: %w", err)
	}

	fmt.Println("Supabase Helm release deployed successfully")

	fmt.Println("Waiting on Supabase Kong Load Balancer to be ready...")
	// Wait for the Supabase stack to be ready
	lbHostname, err := supabase.WaitForLoadBalancer(clusterName, clusterName+"-supabase-kong", 60, 20*time.Second)

	if err != nil {
		return fmt.Errorf("failed to wait for Supabase Kong Load Balancer: %w", err)
	}

	fmt.Println("Supabase Kong Load Balancer ", lbHostname, " is ready")

	// Update the DNS entries in the values file with the provided dnsPrefix
	fmt.Println("Updating deployment DNS entries with prefix:", lbHostname)

	dnsPrefix := strings.TrimSuffix(lbHostname, ".upcloudlb.com")

	err = supabase.UpdateDNS(valuesPath, updatedValuesFile, dnsPrefix)
	if err != nil {
		return fmt.Errorf("failed to update DNS entries in values file: %w", err)
	}

	files := []string{
		updatedValuesFile,
		securePath,
	}

	// If a configPath is provided, add it to the files to be used for the Helm release
	//if configPath != "" {
	//	files = append(files, configPath)
	//}

	supabase.DeployHelmRelease(clusterName, filepath.Join(chartDir, "charts/supabase"), files, true)

	fmt.Println("Supabase stack deployed successfully with updated DNS entries")

	printSummary(lbHostname, clusterName, config)

	return nil
}

func printSummary(lbHostname, namespace string, config *stacksecrets.SupabaseConfig) {
	fmt.Println("Supabase deployed successfully!")
	fmt.Printf("Public endpoint:      http://%s:8000\n", lbHostname)
	fmt.Printf("Namespace:            %s\n", namespace)
	fmt.Printf("ANON_KEY:             %s\n", config.AnonKey)
	fmt.Printf("SERVICE_ROLE_KEY:     %s\n", config.ServiceRoleKey)
	fmt.Printf("POSTGRES_PASSWORD:    %s\n", orNotSet(config.PostgresPassword))
	fmt.Printf("POOLER_TENANT_ID:     %s\n", config.PoolerTenantID)
	fmt.Printf("DASHBOARD_USERNAME:   %s\n", config.DashboardUsername)
	fmt.Printf("DASHBOARD_PASSWORD:   %s\n", config.DashboardPassword)
	fmt.Printf("S3 ENABLED:           %t\n", config.S3Enabled)
	fmt.Printf("S3_BUCKET:            %s\n", orNotSet(config.S3BucketName))
	fmt.Printf("S3_ENDPOINT:          %s\n", orNotSet(config.S3Endpoint))
	fmt.Printf("S3_REGION:            %s\n", orNotSet(config.S3Region))
	fmt.Printf("SMTP ENABLED:         %t\n", config.SmtpEnabled)
	fmt.Printf("SMTP_HOST:           %s\n", orNotSet(config.SmtpHost))
	fmt.Printf("SMTP_PORT:            %s\n", orNotSet(config.SmtpPort))
	fmt.Printf("SMTP_USER:            %s\n", orNotSet(config.SmtpUsername))
	fmt.Printf("SMTP_SENDER_NAME:     %s\n", orNotSet(config.SmtpSenderName))
}

func orNotSet(val string) string {
	if val == "" {
		return "not set"
	}
	return val
}
