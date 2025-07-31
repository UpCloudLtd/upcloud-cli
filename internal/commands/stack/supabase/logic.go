package supabase

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands/stack/core"
	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud/request"
	"gopkg.in/yaml.v3"
)

func (s *deploySupabaseCommand) deploy(exec commands.Executor, chartDir string) (*SupabaseConfig, error) {
	clusterName := fmt.Sprintf("stack-supabase-cluster-%s-%s", s.name, s.zone)
	networkName := fmt.Sprintf("stack-supabase-net-%s-%s", s.name, s.zone)
	var network *upcloud.Network

	clusters, err := exec.All().GetKubernetesClusters(exec.Context(), &request.GetKubernetesClustersRequest{})

	if err != nil {
		return nil, fmt.Errorf("failed to get Kubernetes clusters: %w", err)
	}

	// TODO: We need to include the upgrade option when the cluster already exists
	// For now we are not supporting upgrades or updates
	if core.ClusterExists(clusterName, clusters) {
		return nil, fmt.Errorf("a cluster with the name '%s' already exists", clusterName)
	}

	fmt.Println("Creating Kubernetes cluster:", clusterName)
	// Check if the network already exists
	networks, err := exec.Network().GetNetworks(exec.Context())
	if err != nil {
		return nil, fmt.Errorf("failed to get networks: %w", err)
	}

	network = core.GetNetworkFromName(networkName, networks.Networks)

	// Create the network if it does not exist
	if network == nil {
		fmt.Println("Network does NOT exist, creating:", networkName)
		network, err = core.CreateNetwork(exec, networkName, s.zone)
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

	fmt.Println(("Creating Kubernetes cluster..."))
	cluster, err := exec.All().CreateKubernetesCluster(exec.Context(), &request.CreateKubernetesClusterRequest{
		Name:        clusterName,
		Network:     network.UUID,
		Zone:        s.zone,
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

	// Save kubeconfig to temp file
	kubeconfigPath := filepath.Join(chartDir, "kubeconfig.yaml")
	if err := os.WriteFile(kubeconfigPath, []byte(kubeconfig), 0600); err != nil {
		return nil, fmt.Errorf("failed to write kubeconfig: %w", err)
	}

	// Create a Kubernetes client from the kubeconfig
	os.Setenv("KUBECONFIG", kubeconfigPath)
	kubeClient, err := core.GetKubernetesClient(kubeconfigPath)
	if err != nil {
		return nil, fmt.Errorf("failed to create Kubernetes client: %w", err)
	}

	fmt.Println("Kubeconfig received successfully and saved to:", kubeconfigPath)

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

	// Deploy the Helm release
	err = core.DeployHelmRelease(kubeClient, clusterName, filepath.Join(chartDir, "charts/supabase"), []string{valuesPath, securePath}, false)
	if err != nil {
		return nil, fmt.Errorf("failed to deploy Supabase Helm release: %w", err)
	}

	fmt.Println("Supabase Helm release deployed successfully")

	fmt.Println("Waiting on Supabase Kong Load Balancer to be ready...")
	// Wait for the Supabase stack to be ready
	lbHostname, err := core.WaitForLoadBalancer(kubeClient, clusterName, clusterName+"-supabase-kong", 60, 20*time.Second)

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

	err = updateDNS(valuesPath, updatedValuesFile, dnsPrefix)
	if err != nil {
		return nil, fmt.Errorf("failed to update DNS entries in values file: %w", err)
	}

	files := []string{
		updatedValuesFile,
		securePath,
	}

	core.DeployHelmRelease(kubeClient, clusterName, filepath.Join(chartDir, "charts/supabase"), files, true)

	fmt.Println("Supabase stack deployed successfully with updated DNS entries")

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
