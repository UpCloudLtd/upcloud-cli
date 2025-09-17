package supabase

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands/stack"
	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud/request"
	"gopkg.in/yaml.v3"
)

type ChartVersion string

const (
	ChartVersionV0_1_3 ChartVersion = "v0.1.3"
)

func (s *deploySupabaseCommand) deploy(exec commands.Executor, chartDir string) (*SupabaseConfig, error) {
	clusterName := fmt.Sprintf("stack-supabase-cluster-%s-%s", s.name, s.zone)
	networkName := fmt.Sprintf("stack-supabase-net-%s-%s", s.name, s.zone)
	var network *upcloud.Network

	clusters, err := exec.All().GetKubernetesClusters(exec.Context(), &request.GetKubernetesClustersRequest{})

	if err != nil {
		return nil, fmt.Errorf("failed to get Kubernetes clusters: %w", err)
	}

	// Code will not update an existing cluster, it will create a new one
	if stack.ClusterExists(clusterName, clusters) {
		return nil, fmt.Errorf("a cluster with the name '%s' already exists", clusterName)
	}

	msg := fmt.Sprintf("Creating Kubernetes cluster %s in zone %s", clusterName, s.zone)
	exec.PushProgressStarted(msg)

	// Check if the network already exists
	networks, err := exec.Network().GetNetworks(exec.Context())
	if err != nil {
		return nil, fmt.Errorf("failed to get networks: %w", err)
	}

	network = stack.GetNetworkFromName(networkName, networks.Networks)

	// Create the network if it does not exist
	if network == nil {
		network, err = stack.CreateNetwork(exec, networkName, s.zone, stack.StackTypeSupabase)
		if err != nil {
			return nil, fmt.Errorf("failed to create network: %w for kubernetes deployment", err)
		}

		if network == nil {
			return nil, fmt.Errorf("created network %s is nil", networkName)
		}
		if len(network.IPNetworks) == 0 {
			return nil, fmt.Errorf("created network %s has no IP networks", networkName)
		}
	}

	cluster, err := exec.All().CreateKubernetesCluster(exec.Context(), &request.CreateKubernetesClusterRequest{
		Name:        clusterName,
		Network:     network.UUID,
		Zone:        s.zone,
		NetworkCIDR: network.IPNetworks[0].Address,
		Labels: []upcloud.Label{
			{Key: "stacks.upcloud.com/stack", Value: "supabase"},
			{Key: "stacks.upcloud.com/created-by", Value: "upctl"},
			{Key: "stacks.upcloud.com/chart-version", Value: string(ChartVersionV0_1_3)},
			{Key: "stacks.upcloud.com/version", Value: string(stack.VersionV0_1_0_0)},
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

	exec.PushProgressSuccess(msg)

	exec.PushProgressStarted("Setting up environment for Supabase stack deployment")
	// Get kubeconfig file for the cluster
	kubeconfigPath, err := stack.WriteKubeconfigToFile(exec, cluster.UUID, chartDir)
	if err != nil {
		return nil, fmt.Errorf("failed to write kubeconfig to file: %w", err)
	}

	// Create a Kubernetes client from the kubeconfig
	os.Setenv("KUBECONFIG", kubeconfigPath)
	kubeClient, err := stack.GetKubernetesClient(kubeconfigPath)
	if err != nil {
		return nil, fmt.Errorf("failed to create Kubernetes client: %w", err)
	}

	exec.PushProgressSuccess("Setting up environment for Supabase stack deployment")

	exec.PushProgressStarted("Generating configuration files for Supabase stack deployment")
	// Generate configuration for the Supabase stack from input config file
	config, err := Generate(s.configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to generate secrets: %w", err)
	}

	valuesPath := filepath.Join(chartDir, "charts/supabase/values.example.yaml")
	securePath := filepath.Join(chartDir, fmt.Sprintf("supabase-%s-values.secure.yaml", s.name))
	updatedValuesFile := filepath.Join(chartDir, fmt.Sprintf("supabase-%s-values.updated.yaml", s.name))
	err = WriteSecureValues(securePath, config)
	if err != nil {
		return nil, fmt.Errorf("failed to write secrets file: %w", err)
	}

	exec.PushProgressSuccess("Generating configuration files for Supabase stack deployment")

	msg = fmt.Sprintf("Deploying Supabase stack in cluster %s in zone %s", clusterName, s.zone)
	exec.PushProgressStarted(msg)

	// Deploy the Helm release
	err = stack.DeployHelmRelease(kubeClient, clusterName, filepath.Join(chartDir, "charts/supabase"), []string{valuesPath, securePath}, false)
	if err != nil {
		return nil, fmt.Errorf("failed to deploy Supabase Helm release: %w", err)
	}

	// Wait for the Supabase stack to be ready
	lbHostname, err := stack.WaitForLoadBalancer(kubeClient, clusterName, clusterName+"-supabase-kong", 60, 20*time.Second)
	if err != nil {
		return nil, fmt.Errorf("failed to wait for Supabase Kong Load Balancer: %w", err)
	}

	// Adding the cluster name and lbHostname to the config for reporting purposes
	config.LbHostname = lbHostname
	config.ClusterName = clusterName
	dnsPrefix := strings.TrimSuffix(lbHostname, ".upcloudlb.com")

	err = updateDNS(valuesPath, updatedValuesFile, dnsPrefix)
	if err != nil {
		return nil, fmt.Errorf("failed to update DNS entries in values file: %w", err)
	}

	files := []string{
		updatedValuesFile,
		securePath,
	}

	stack.DeployHelmRelease(kubeClient, clusterName, filepath.Join(chartDir, "charts/supabase"), files, true)

	exec.PushProgressSuccess(msg)

	return config, nil
}

func summaryOutput(config *SupabaseConfig) []byte {
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

// updateDNS updates the DNS entries in the values file with the provided dnsPrefix.
func updateDNS(valuesFile, updatedValuesFile, dnsPrefix string) error {
	input, err := os.ReadFile(valuesFile)
	if err != nil {
		return fmt.Errorf("failed to read values file: %w", err)
	}

	var root yaml.Node
	if err := yaml.Unmarshal(input, &root); err != nil {
		return fmt.Errorf("failed to parse YAML: %w", err)
	}

	// List of keys to update
	targetFields := [][][]string{
		{{"studio", "environment", "SUPABASE_PUBLIC_URL"}},
		{{"auth", "environment", "GOTRUE_SITE_URL"}},
		{{"auth", "environment", "API_EXTERNAL_URL"}},
		{{"rest", "environment", "POSTGREST_SITE_URL"}},
		{{"realtime", "environment", "PORTAL_URL"}},
		{{"storage", "environment", "FILE_STORAGE_BACKEND_URL"}},
		{{"kong", "environment", "SUPABASE_PUBLIC_URL"}},
		{{"kong", "environment", "SUPABASE_STUDIO_URL"}},
	}

	updated := false
	for _, path := range targetFields {
		if updateValueAtPath(&root, path[0], dnsPrefix) {
			updated = true
		}
	}

	if !updated {
		return fmt.Errorf("no DNS entries were updated")
	}

	output, err := yaml.Marshal(&root)
	if err != nil {
		return fmt.Errorf("failed to marshal YAML: %w", err)
	}

	if err := os.WriteFile(updatedValuesFile, output, 0644); err != nil {
		return fmt.Errorf("failed to write updated values file: %w", err)
	}

	return nil
}

// updateValueAtPath updates the value at the specified path in the YAML node with the given dnsPrefix.
// It returns true if the value was updated, false otherwise.
func updateValueAtPath(node *yaml.Node, path []string, dnsPrefix string) bool {
	if node.Kind != yaml.DocumentNode || len(node.Content) == 0 {
		return false
	}

	current := node.Content[0]
	for _, key := range path[:len(path)-1] {
		found := false
		for i := 0; i < len(current.Content); i += 2 {
			k := current.Content[i]
			v := current.Content[i+1]
			if k.Value == key {
				current = v
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}

	lastKey := path[len(path)-1]
	for i := 0; i < len(current.Content); i += 2 {
		k := current.Content[i]
		v := current.Content[i+1]
		if k.Value == lastKey {
			originalValue := v.Value
			newValue := strings.Replace(originalValue, "REPLACEME.upcloudlb.com", dnsPrefix+".upcloudlb.com", -1)
			v.Value = newValue
			return true
		}
	}

	return false
}
