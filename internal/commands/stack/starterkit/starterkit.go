package starterkit

import (
	"context"
	"fmt"
	"strings"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands/stack"
	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud/request"
)

type Version string

const (
	VersionV1 Version = "v1.0.0"
)

type StarterKitConfig struct {
	ProjectName       string
	Zone              string
	NetworkName       string
	ClusterName       string
	ObjectStorageName string
	DBName            string
	RouterName        string
}

func GenerateStarterKitConfig(projectName, zone string) *StarterKitConfig {
	clusterName := fmt.Sprintf("starter-kit-cluster-%s-%s", projectName, zone)
	networkName := fmt.Sprintf("starter-kit-network-%s-%s", projectName, zone)
	objectStorageName := fmt.Sprintf("starter-kit-object-storage-%s-%s", projectName, zone)
	dbName := fmt.Sprintf("starter-kit-db-%s-%s", projectName, zone)
	routerName := fmt.Sprintf("starter-kit-router-%s-%s", projectName, zone)

	return &StarterKitConfig{
		ProjectName:       projectName,
		Zone:              zone,
		NetworkName:       networkName,
		ClusterName:       clusterName,
		ObjectStorageName: objectStorageName,
		DBName:            dbName,
		RouterName:        routerName,
	}
}

func (c *StarterKitConfig) Validate(exec commands.Executor) error {
	// Check if a network with the same NetworkName already exists
	networks, err := exec.Network().GetNetworks(exec.Context())
	if err != nil {
		return fmt.Errorf("failed to get networks: %w", err)
	}

	network := stack.GetNetworkFromName(c.NetworkName, networks.Networks)
	if network != nil {
		return fmt.Errorf(
			"a network with the name %q already exists. Names are generated as 'starter-kit-network-<project>-<zone>'",
			c.NetworkName,
		)
	}

	// Check if a router with the same RouterName already exists
	routers, err := exec.All().GetRouters(exec.Context())
	if err != nil {
		return fmt.Errorf("failed to get routers: %w", err)
	}

	for _, router := range routers.Routers {
		if router.Name == c.RouterName {
			return fmt.Errorf(
				"a router with the name %q already exists. Names are generated as 'starter-kit-router-<project>-<zone>'",
				c.RouterName,
			)
		}
	}

	// Check if a kubernetes cluster with the name clusterName already exists
	clusters, err := exec.All().GetKubernetesClusters(exec.Context(), &request.GetKubernetesClustersRequest{})

	if err != nil {
		return fmt.Errorf("failed to get Kubernetes clusters: %w", err)
	}

	if stack.ClusterExists(c.ClusterName, clusters) {
		return fmt.Errorf(
			"a cluster with the name %q already exists. Names are generated as 'starter-kit-cluster-<project>-<zone>' (got project=%q, zone=%q)",
			c.ClusterName, c.ProjectName, c.Zone,
		)
	}

	// Check if a database with the same DBName already exists
	dbs, err := exec.All().GetAllManagedDatabases(exec.Context())
	if err != nil {
		return fmt.Errorf("failed to get databases: %w", err)
	}

	for _, db := range dbs {
		if db.Name == c.DBName {
			return fmt.Errorf(
				"a database with the name %q already exists. Names are generated as 'starter-kit-db-<project>-<zone>'",
				c.DBName,
			)
		}
	}

	// Check if an object storage with the same ObjectStorageName already exists
	objectstorages, err := exec.All().GetManagedObjectStorages(exec.Context(), &request.GetManagedObjectStoragesRequest{})
	if err != nil {
		return fmt.Errorf("failed to get object storages: %w", err)
	}
	for _, objsto := range objectstorages {
		if objsto.Name == c.ObjectStorageName {
			return fmt.Errorf(
				"an object storage with the name %q already exists. Names are generated as 'starter-kit-object-storage-<project>-<zone>'",
				c.ObjectStorageName,
			)
		}
	}

	// Check we have a mapping from the selected zone to a valid object storage region
	_, err = stack.GetObjectStorageRegionFromZone(exec, c.Zone)
	if err != nil {
		return fmt.Errorf("failed to validate object storage region from zone %q: %w, contact support because you might be using a new zone not supported by the deploy command", c.Zone, err)
	}

	return nil
}

func buildSummary(
	cluster *upcloud.KubernetesCluster,
	kubeconfigPath string,
	network *upcloud.Network,
	router *upcloud.Router,
	db *upcloud.ManagedDatabase,
	obj *upcloud.ManagedObjectStorage,
) string {
	var b strings.Builder

	// Header
	fmt.Fprintf(&b, "Starter Kit deployed successfully!\n\n")

	// Kubernetes
	fmt.Fprintf(&b, "KUBERNETES CLUSTER\n")
	if cluster != nil {
		fmt.Fprintf(&b, "  Name:        %s\n", cluster.Name)
		fmt.Fprintf(&b, "  UUID:        %s\n", cluster.UUID)
		fmt.Fprintf(&b, "  Zone:        %s\n", cluster.Zone)
		fmt.Fprintf(&b, "  Network:     %s\n", cluster.Network)
		if kubeconfigPath != "" {
			fmt.Fprintf(&b, "  Kubeconfig:  %s\n", kubeconfigPath)
			fmt.Fprintf(&b, "  Set env:     export KUBECONFIG=%s\n", kubeconfigPath)
			fmt.Fprintf(&b, "  Test:        kubectl get nodes\n")
			fmt.Fprintf(&b, "  Ingress LB:  kubectl -n ingress-nginx get svc ingress-nginx-controller -o jsonpath='{.status.loadBalancer.ingress[0].hostname}{\"\\n\"}'\n")
		}
	}
	b.WriteString("\n")

	// Network & Router
	fmt.Fprintf(&b, "NETWORKING\n")
	if network != nil {
		fmt.Fprintf(&b, "  Network:     %s (UUID: %s)\n", network.Name, network.UUID)
		if len(network.IPNetworks) > 0 {
			fmt.Fprintf(&b, "  CIDR:        %s\n", network.IPNetworks[0].Address)
			if network.IPNetworks[0].DHCP == upcloud.True {
				fmt.Fprintf(&b, "  DHCP:        enabled\n")
			} else {
				fmt.Fprintf(&b, "  DHCP:        disabled\n")
			}
		}
	}
	if router != nil {
		fmt.Fprintf(&b, "  Router:      %s (UUID: %s)\n", router.Name, router.UUID)
	}
	b.WriteString("\n")

	// Managed Database
	fmt.Fprintf(&b, "MANAGED DATABASE\n")
	if db != nil {
		fmt.Fprintf(&b, "  Name:        %s (UUID: %s)\n", db.Title, db.UUID)
		fmt.Fprintf(&b, "  Type/Plan:   %s / %s\n", db.Type, db.Plan)
		fmt.Fprintf(&b, "  State:       %s\n", db.State)
		// TODO: We are showing the password in clear text here. Is that ok?
		fmt.Fprintf(&b, "  ServiceURI :        %s\n", db.ServiceURI)
	} else {
		fmt.Fprintf(&b, "  (not created)\n")
	}
	b.WriteString("\n")

	// Managed Object Storage
	fmt.Fprintf(&b, "OBJECT STORAGE\n")
	if obj != nil {
		fmt.Fprintf(&b, "  Name:        %s (UUID: %s)\n", obj.Name, obj.UUID)
		fmt.Fprintf(&b, "  Region:      %s\n", obj.Region)
		fmt.Fprintf(&b, "  State:      %s\n", obj.OperationalState)

		// If API provides endpoint(s)
		if len(obj.Endpoints) > 0 {
			fmt.Fprintf(&b, "  DomainName:    %s\n", obj.Endpoints[0].DomainName)
			fmt.Fprintf(&b, "  Type:    %s\n", obj.Endpoints[0].Type)
			fmt.Fprintf(&b, "  IAMURL:    %s\n", obj.Endpoints[0].IAMURL)
			fmt.Fprintf(&b, "  STSURL:    %s\n", obj.Endpoints[0].STSURL)
		}

	} else {
		fmt.Fprintf(&b, "  (not created)\n")
	}
	b.WriteString("\n")

	// ACCESS without any Load Balancer
	fmt.Fprintf(&b, "ACCESS (no load balancer created)\n")
	fmt.Fprintf(&b, "  1) Local dev via port-forward (recommended for quick testing):\n")
	fmt.Fprintf(&b, "     kubectl -n <namespace> port-forward svc/<your-service> 8080:80\n")
	fmt.Fprintf(&b, "     # then open http://localhost:8080\n\n")

	fmt.Fprintf(&b, "  2) NodePort (reachable inside the private network):\n")
	fmt.Fprintf(&b, "     # switch your Service to NodePort\n")
	fmt.Fprintf(&b, "     kubectl -n <namespace> patch svc <your-service> -p '{\"spec\":{\"type\":\"NodePort\"}}'\n")
	fmt.Fprintf(&b, "     # find the assigned nodePort and a node's private IP\n")
	fmt.Fprintf(&b, "     kubectl -n <namespace> get svc <your-service> -o jsonpath='{.spec.ports[0].nodePort}{\"\\n\"}'\n")
	fmt.Fprintf(&b, "     kubectl get nodes -o wide\n")
	fmt.Fprintf(&b, "     # then browse: http://<node-private-ip>:<nodePort>\n")
	fmt.Fprintf(&b, "     (for external access: use a VPN/bastion into this private network)\n\n")

	fmt.Fprintf(&b, "  3) Private-only options:\n")
	fmt.Fprintf(&b, "     - VPN/peering to the VPC and use ClusterIP/NodePort directly\n")
	fmt.Fprintf(&b, "     - Bastion host with SSH tunnels (e.g., ssh -L 8080:<node-ip>:<nodePort> ...)\n\n")

	// Final tips
	fmt.Fprintf(&b, "NEXT STEPS\n")
	if kubeconfigPath != "" {
		fmt.Fprintf(&b, "  export KUBECONFIG=%s\n", kubeconfigPath)
	}
	fmt.Fprintf(&b, "  Deploy ingress-nginx and your app, then point DNS (CNAME) to the LB hostname shown above.\n")

	return b.String()
}

func createKubernetes(ctx context.Context, exec commands.Executor, config *StarterKitConfig, network *upcloud.Network, projectDir string) (*upcloud.KubernetesCluster, string, error) {
	exec.PushProgressStarted("Deploying Kubernetes cluster")
	cluster, err := exec.All().CreateKubernetesCluster(ctx, &request.CreateKubernetesClusterRequest{
		Name:        config.ClusterName,
		Network:     network.UUID,
		Zone:        config.Zone,
		NetworkCIDR: network.IPNetworks[0].Address,
		Labels: []upcloud.Label{
			{Key: "stacks.upcloud.com/stack", Value: string(stack.StackTypeStarterKit)},
			{Key: "stacks.upcloud.com/created-by", Value: "upctl"},
			{Key: "stacks.upcloud.com/version", Value: string(stack.VersionV0_1_0_0)},
			{Key: "stacks.upcloud.com/name", Value: config.ClusterName},
		},
		NodeGroups: []request.KubernetesNodeGroup{
			{
				Name:  "starter-kit-node-group",
				Count: 1,
				Plan:  "2xCPU-4GB",
			},
		},
	})

	if err != nil {
		return nil, "", fmt.Errorf("failed to create Kubernetes cluster: %w", err)
	}

	// Write kubeconfig file for the cluster to disk
	kubeconfigPath, err := stack.WriteKubeconfigToFile(exec, cluster.UUID, projectDir)
	if err != nil {
		return nil, "", fmt.Errorf("failed to write kubeconfig to file: %w", err)
	}
	exec.PushProgressSuccess("Deploying Kubernetes cluster")
	return cluster, kubeconfigPath, nil
}

func createDatabase(ctx context.Context, exec commands.Executor, config *StarterKitConfig, network *upcloud.Network) (*upcloud.ManagedDatabase, error) {
	exec.PushProgressStarted("Deploying Database")
	db, err := exec.All().CreateManagedDatabase(ctx, &request.CreateManagedDatabaseRequest{
		HostNamePrefix: config.DBName,
		Title:          config.DBName,
		Zone:           config.Zone,
		Plan:           "1x1xCPU-2GB-25GB",
		Type:           "pg",
		Networks: []upcloud.ManagedDatabaseNetwork{
			{
				UUID:   &network.UUID,
				Type:   upcloud.NetworkTypePrivate,
				Name:   config.NetworkName,
				Family: upcloud.IPAddressFamilyIPv4,
			},
		},
		Labels: []upcloud.Label{
			{Key: "stacks.upcloud.com/stack", Value: string(stack.StackTypeStarterKit)},
			{Key: "stacks.upcloud.com/created-by", Value: "upctl"},
			{Key: "stacks.upcloud.com/version", Value: string(stack.VersionV0_1_0_0)},
			{Key: "stacks.upcloud.com/name", Value: config.DBName},
		},
	})

	if err != nil {
		return nil, fmt.Errorf("failed to create managed database: %w", err)
	}
	exec.PushProgressSuccess("Deploying Database")
	return db, nil
}

func createObjectStorage(ctx context.Context, exec commands.Executor, config *StarterKitConfig, network *upcloud.Network) (*upcloud.ManagedObjectStorage, error) {
	exec.PushProgressStarted("Deploying Object Storage")
	stack.GetObjectStorageRegionFromZone(exec, config.Zone)
	region, err := stack.GetObjectStorageRegionFromZone(exec, config.Zone)
	if err != nil {
		return nil, fmt.Errorf("failed to validate object storage region from zone %q: %w, contact support because you might be using a new zone not supported by the deploy command", config.Zone, err)
	}

	objStorage, err := exec.All().CreateManagedObjectStorage(ctx, &request.CreateManagedObjectStorageRequest{
		Name:             config.ObjectStorageName,
		Region:           region,
		ConfiguredStatus: upcloud.ManagedObjectStorageConfiguredStatusStarted,
		Networks: []upcloud.ManagedObjectStorageNetwork{
			{
				UUID:   &network.UUID,
				Type:   upcloud.NetworkTypePrivate,
				Name:   config.NetworkName,
				Family: upcloud.IPAddressFamilyIPv4,
			},
		},
		Labels: []upcloud.Label{
			{Key: "stacks.upcloud.com/stack", Value: string(stack.StackTypeStarterKit)},
			{Key: "stacks.upcloud.com/created-by", Value: "upctl"},
			{Key: "stacks.upcloud.com/version", Value: string(stack.VersionV0_1_0_0)},
			{Key: "stacks.upcloud.com/name", Value: config.ObjectStorageName},
		},
	})

	if err != nil {
		return nil, fmt.Errorf("failed to create object storage: %w", err)
	}
	exec.PushProgressSuccess("Deploying Object Storage")
	return objStorage, nil
}
