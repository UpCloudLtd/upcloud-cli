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
	clusterName := fmt.Sprintf("%s-%s-%s", stack.SupabaseResourceRootNameCluster, s.name, s.zone)
	networkName := fmt.Sprintf("%s-%s-%s", stack.SupabaseResourceRootNameNetwork, s.name, s.zone)
	var network *upcloud.Network

	msg := "Generating configuration files for Supabase stack deployment"
	exec.PushProgressStarted(msg)

	config, err := GenerateDefaultConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to generate default configuration for this deployment: %w", err)
	}
	config.Name = s.name
	config.Zone = s.zone

	// Validate early input config file if it exists
	if s.configPath != "" {
		err := loadConfigFromFile(s.configPath, config)
		if err != nil {
			return nil, fmt.Errorf("failed to load configuration from file: %w", err)
		}

		err = validateConfig(config)
		if err != nil {
			return nil, fmt.Errorf("failed to validate configuration file: %w", err)
		}
	}

	exec.PushProgressSuccess(msg)

	clusters, err := exec.All().GetKubernetesClusters(exec.Context(), &request.GetKubernetesClustersRequest{})
	if err != nil {
		return nil, fmt.Errorf("failed to get Kubernetes clusters: %w", err)
	}

	// Code will not update an existing cluster, it will create a new one
	if stack.ClusterExists(clusterName, clusters) {
		return nil, fmt.Errorf("a cluster with the name '%s' already exists", clusterName)
	}

	msg = fmt.Sprintf("Creating Kubernetes cluster %s in zone %s", clusterName, s.zone)
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

	msg = "Getting Kubernetes cluster details"
	exec.PushProgressStarted(msg)

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

	exec.PushProgressSuccess(msg)

	if shouldCreateObjectStorage(config) {
		objStorageMsg := "Generating object storage for Supabase stack deployment"
		exec.PushProgressStarted(objStorageMsg)
		objStorageName := fmt.Sprintf("%s-%s-%s", stack.SupabaseResourceRootNameObjStorage, s.name, s.zone)
		objStorageRegion, err := stack.GetObjectStorageRegionFromZone(exec, s.zone)
		if err != nil {
			return nil, fmt.Errorf("failed to get object storage region from zone: %w", err)
		}

		objStorage, err := exec.All().CreateManagedObjectStorage(exec.Context(), &request.CreateManagedObjectStorageRequest{
			ConfiguredStatus: upcloud.ManagedObjectStorageConfiguredStatusStarted,
			Region:           objStorageRegion,
			Name:             objStorageName,
			Labels: []upcloud.Label{
				{Key: "stacks.upcloud.com/stack", Value: string(stack.StackTypeSupabase)},
				{Key: "stacks.upcloud.com/created-by", Value: "upctl"},
				{Key: "stacks.upcloud.com/version", Value: string(stack.VersionV0_1_0_0)},
				{Key: "stacks.upcloud.com/name", Value: objStorageName},
			},
			Networks: []upcloud.ManagedObjectStorageNetwork{
				{
					UUID:   &network.UUID,
					Type:   upcloud.NetworkTypePrivate,
					Name:   network.Name,
					Family: upcloud.IPAddressFamilyIPv4,
				},
			},
		})
		if err != nil {
			return nil, fmt.Errorf("failed to create object storage for this deployment: %w", err)
		}

		objStorage, err = exec.All().WaitForManagedObjectStorageOperationalState(exec.Context(), &request.WaitForManagedObjectStorageOperationalStateRequest{
			DesiredState: upcloud.ManagedObjectStorageOperationalStateRunning,
			UUID:         objStorage.UUID,
		})
		if err != nil {
			return nil, fmt.Errorf("error while waiting for the object storage to become online: %w", err)
		}

		// Create user for the object storage
		user, err := exec.All().CreateManagedObjectStorageUser(exec.Context(), &request.CreateManagedObjectStorageUserRequest{
			Username:    "supabase-user",
			ServiceUUID: objStorage.UUID,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to create user for the object storage: %w", err)
		}

		err = exec.All().AttachManagedObjectStorageUserPolicy(exec.Context(), &request.AttachManagedObjectStorageUserPolicyRequest{
			ServiceUUID: objStorage.UUID,
			Username:    user.Username,
			Name:        "ECSS3FullAccess",
		})
		if err != nil {
			return nil, fmt.Errorf("failed to attach user policy to object storage: %w", err)
		}

		// Create an access key for the object storage
		userAccessKey, err := exec.All().CreateManagedObjectStorageUserAccessKey(exec.Context(), &request.CreateManagedObjectStorageUserAccessKeyRequest{
			Username:    user.Username,
			ServiceUUID: objStorage.UUID,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to create access key for the object storage: %w", err)
		}

		// Create a bucket for the object storage
		bucketName := "supabase-storage"
		_, err = exec.All().CreateManagedObjectStorageBucket(exec.Context(), &request.CreateManagedObjectStorageBucketRequest{
			Name:        bucketName,
			ServiceUUID: objStorage.UUID,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to create bucket for the object storage: %w", err)
		}

		// Update the config with the object storage details
		config.S3Enabled = true
		config.S3KeyID = userAccessKey.AccessKeyID
		config.S3AccessKey = *userAccessKey.SecretAccessKey
		config.S3BucketName = bucketName
		config.S3Region = objStorage.Region
		if len(objStorage.Endpoints) > 0 {
			config.S3Endpoint = fmt.Sprintf("https://%s", objStorage.Endpoints[0].DomainName)
		} else {
			return nil, fmt.Errorf("object storage has no endpoints")
		}

		err = validateConfig(config)
		if err != nil {
			return nil, fmt.Errorf("invalid configuration: %w", err)
		}

		exec.PushProgressSuccess(objStorageMsg)
	}

	valuesPath := filepath.Join(chartDir, "charts/supabase/values.example.yaml")
	securePath := filepath.Join(chartDir, fmt.Sprintf("supabase-%s-values.secure.yaml", s.name))
	updatedValuesFile := filepath.Join(chartDir, fmt.Sprintf("supabase-%s-values.updated.yaml", s.name))
	err = WriteSecureValues(securePath, config)
	if err != nil {
		return nil, fmt.Errorf("failed to write secrets file: %w", err)
	}

	msg = fmt.Sprintf("Deploying Supabase stack in cluster %s in zone %s", clusterName, s.zone)
	exec.PushProgressStarted(msg)

	// Deploy the Helm release
	err = stack.DeployHelmRelease(kubeClient, clusterName, filepath.Join(chartDir, "charts/supabase"), []string{valuesPath, securePath}, false)
	if err != nil {
		return nil, fmt.Errorf("failed to deploy Supabase: %w", err)
	}

	// Wait for the Supabase stack to be ready
	lbHostname, err := stack.WaitForLoadBalancer(kubeClient, clusterName, clusterName+"-supabase-kong", 60, 20*time.Second)
	if err != nil {
		return nil, fmt.Errorf("timeout: Supabase Kong Load Balancer: %w", err)
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

	err = stack.DeployHelmRelease(kubeClient, clusterName, filepath.Join(chartDir, "charts/supabase"), files, true)
	if err != nil {
		return nil, fmt.Errorf("failed to update endpoints in Supabase: %w", err)
	}

	exec.PushProgressSuccess(msg)

	return config, nil
}

func summaryOutput(config *SupabaseConfig) string {
	builder := &strings.Builder{}

	fmt.Fprintln(builder, "Supabase deployed successfully!")
	fmt.Fprintf(builder, "Project Name:            %s\n", config.Name)
	fmt.Fprintf(builder, "Project Zone:            %s\n", config.Zone)
	fmt.Fprintf(builder, "Public endpoint:      http://%s:8000\n", config.LbHostname)
	fmt.Fprintf(builder, "Namespace:            %s\n", config.ClusterName)
	fmt.Fprintf(builder, "ANON_KEY:             %s\n", config.AnonKey)
	fmt.Fprintf(builder, "SERVICE_ROLE_KEY:     %s\n", config.ServiceRoleKey)
	fmt.Fprintf(builder, "POSTGRES_PASSWORD:    %s\n", orNotSet(config.PostgresPassword))
	fmt.Fprintf(builder, "POOLER_TENANT_ID:     %s\n", config.PoolerTenantID)
	fmt.Fprintf(builder, "DASHBOARD_USERNAME:   %s\n", config.DashboardUsername)
	fmt.Fprintf(builder, "DASHBOARD_PASSWORD:   %s\n", config.DashboardPassword)

	// S3 section
	if config.S3Enabled {
		fmt.Fprintf(builder, "S3 ENABLED:           true\n")
		fmt.Fprintf(builder, "S3_KEY_ID:            %s\n", orNotSet(config.S3KeyID))
		fmt.Fprintf(builder, "S3_ACCESS_KEY:        %s\n", orNotSet(config.S3AccessKey))
		fmt.Fprintf(builder, "S3_BUCKET:            %s\n", orNotSet(config.S3BucketName))
		fmt.Fprintf(builder, "S3_ENDPOINT:          %s\n", orNotSet(config.S3Endpoint))
		fmt.Fprintf(builder, "S3_REGION:            %s\n", orNotSet(config.S3Region))
	} else {
		fmt.Fprintf(builder, "S3 ENABLED:           false\n")
		fmt.Fprintf(builder, "S3 CONFIG:            Not available (S3 is disabled)\n")
	}

	// SMTP section
	if config.SMTPEnabled {
		fmt.Fprintf(builder, "SMTP ENABLED:         true\n")
		fmt.Fprintf(builder, "SMTP_HOST:            %s\n", orNotSet(config.SMTPHost))
		fmt.Fprintf(builder, "SMTP_PORT:            %s\n", orNotSet(config.SMTPPort))
		fmt.Fprintf(builder, "SMTP_USER:            %s\n", orNotSet(config.SMTPUsername))
		fmt.Fprintf(builder, "SMTP_SENDER_NAME:     %s\n", orNotSet(config.SMTPSenderName))
		fmt.Fprintf(builder, "SMTP_SENDER_EMAIL:    %s\n", orNotSet(config.SMTPAdminEmail))
	} else {
		fmt.Fprintf(builder, "SMTP ENABLED:         false\n")
		fmt.Fprintf(builder, "SMTP CONFIG:          Not available (SMTP is disabled)\n")
	}

	return builder.String()
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

	if err := os.WriteFile(updatedValuesFile, output, 0o600); err != nil {
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
			newValue := strings.ReplaceAll(originalValue, "REPLACEME.upcloudlb.com", dnsPrefix+".upcloudlb.com")
			v.Value = newValue
			return true
		}
	}

	return false
}
