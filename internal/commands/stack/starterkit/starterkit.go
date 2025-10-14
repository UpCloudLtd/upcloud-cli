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
	clusterName := fmt.Sprintf("%s-%s-%s", stack.StarterKitResourceRootNameCluster, projectName, zone)
	networkName := fmt.Sprintf("%s-%s-%s", stack.StarterKitResourceRootNameNetwork, projectName, zone)
	objectStorageName := fmt.Sprintf("%s-%s-%s", stack.StarterKitResourceRootNameObjStorage, projectName, zone)
	dbName := fmt.Sprintf("%s-%s-%s", stack.StarterKitResourceRootNameDatabase, projectName, zone)
	routerName := fmt.Sprintf("%s-%s-%s", stack.StarterKitResourceRootNameRouter, projectName, zone)

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
			"a cluster with the name %q already exists. Names are generated as 'stack-starter-cluster-<project>-<zone>' (got project=%q, zone=%q)",
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
				"a database with the name %q already exists. Names are generated as 'stack-starter-db-<project>-<zone>'",
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
				"an object storage with the name %q already exists. Names are generated as 'stack-starter-obj-sto-<project>-<zone>'",
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
	objAcc *upcloud.ManagedObjectStorageUserAccessKey,
	objBucket string,
) string {
	var b strings.Builder

	// Header
	b.WriteString("Starter Kit deployed successfully!\n\n")

	// Kubernetes
	b.WriteString("KUBERNETES CLUSTER\n")
	if cluster != nil {
		b.WriteString(fmt.Sprintf("  Name:        %s\n", cluster.Name))
		b.WriteString(fmt.Sprintf("  UUID:        %s\n", cluster.UUID))
		b.WriteString(fmt.Sprintf("  Zone:        %s\n", cluster.Zone))
		b.WriteString(fmt.Sprintf("  Network:     %s\n", cluster.Network))
		if kubeconfigPath != "" {
			b.WriteString(fmt.Sprintf("  Kubeconfig:  %s\n", kubeconfigPath))
			b.WriteString(fmt.Sprintf("  Set env:     export KUBECONFIG=%s\n", kubeconfigPath))
			b.WriteString("  Test:        kubectl get nodes\n")
			b.WriteString("  Ingress LB:  kubectl -n ingress-nginx get svc ingress-nginx-controller -o jsonpath='{.status.loadBalancer.ingress[0].hostname}{\"\\n\"}'\n")
		}
	}
	b.WriteString("\n")

	// Network & Router
	b.WriteString("NETWORKING\n")
	if network != nil {
		b.WriteString(fmt.Sprintf("  Network:     %s (UUID: %s)\n", network.Name, network.UUID))
		if len(network.IPNetworks) > 0 {
			b.WriteString(fmt.Sprintf("  CIDR:        %s\n", network.IPNetworks[0].Address))
			if network.IPNetworks[0].DHCP == upcloud.True {
				b.WriteString("  DHCP:        enabled\n")
			} else {
				b.WriteString("  DHCP:        disabled\n")
			}
		}
	}
	if router != nil {
		b.WriteString(fmt.Sprintf("  Router:      %s (UUID: %s)\n", router.Name, router.UUID))
	}
	b.WriteString("\n")

	// Managed Database
	b.WriteString("MANAGED DATABASE\n")
	if db != nil {
		b.WriteString(fmt.Sprintf("  Name:        %s (UUID: %s)\n", db.Title, db.UUID))
		b.WriteString(fmt.Sprintf("  Type/Plan:   %s / %s\n", db.Type, db.Plan))
		b.WriteString(fmt.Sprintf("  State:       %s\n", db.State))
		b.WriteString(fmt.Sprintf("  ServiceURI:  %s\n", db.ServiceURI))
	} else {
		b.WriteString("  (not created)\n")
	}
	b.WriteString("\n")

	// Managed Object Storage
	b.WriteString("OBJECT STORAGE\n")
	if obj != nil {
		b.WriteString(fmt.Sprintf("  Name:        %s (UUID: %s)\n", obj.Name, obj.UUID))
		b.WriteString(fmt.Sprintf("  Region:      %s\n", obj.Region))
		b.WriteString(fmt.Sprintf("  State:       %s\n", obj.OperationalState))

		// If API provides endpoint(s)
		if len(obj.Endpoints) > 0 {
			b.WriteString(fmt.Sprintf("  DomainName:  %s\n", obj.Endpoints[0].DomainName))
			b.WriteString(fmt.Sprintf("  Type:        %s\n", obj.Endpoints[0].Type))
			b.WriteString(fmt.Sprintf("  IAMURL:      %s\n", obj.Endpoints[0].IAMURL))
			b.WriteString(fmt.Sprintf("  STSURL:      %s\n", obj.Endpoints[0].STSURL))
		}
		// If bucket was created
		if objBucket != "" {
			b.WriteString(fmt.Sprintf("  Bucket:      %s\n", objBucket))
		}
		// If access key was created
		if objAcc != nil {
			b.WriteString(fmt.Sprintf("  AccessKey:   %s\n", objAcc.AccessKeyID))
			b.WriteString(fmt.Sprintf("  SecretKey:   %s\n", *objAcc.SecretAccessKey))
		}
	} else {
		b.WriteString("  (not created)\n")
	}
	b.WriteString("\n")

	// ACCESS without any Load Balancer
	b.WriteString("ACCESS (no load balancer created)\n")
	b.WriteString("  1) Local dev via port-forward (recommended for quick testing):\n")
	b.WriteString("     kubectl -n <namespace> port-forward svc/<your-service> 8080:80\n")
	b.WriteString("     # then open http://localhost:8080\n\n")

	b.WriteString("  2) NodePort (reachable inside the private network):\n")
	b.WriteString("     # switch your Service to NodePort\n")
	b.WriteString("     kubectl -n <namespace> patch svc <your-service> -p '{\"spec\":{\"type\":\"NodePort\"}}'\n")
	b.WriteString("     # find the assigned nodePort and a node's private IP\n")
	b.WriteString("     kubectl -n <namespace> get svc <your-service> -o jsonpath='{.spec.ports[0].nodePort}{\"\\n\"}'\n")
	b.WriteString("     kubectl get nodes -o wide\n")
	b.WriteString("     # then browse: http://<node-private-ip>:<nodePort>\n")
	b.WriteString("     (for external access: use a VPN/bastion into this private network)\n\n")

	b.WriteString("  3) Private-only options:\n")
	b.WriteString("     - VPN/peering to the VPC and use ClusterIP/NodePort directly\n")
	b.WriteString("     - Bastion host with SSH tunnels (e.g., ssh -L 8080:<node-ip>:<nodePort> ...)\n\n")

	// Final tips
	b.WriteString("NEXT STEPS\n")
	if kubeconfigPath != "" {
		b.WriteString(fmt.Sprintf("  export KUBECONFIG=%s\n", kubeconfigPath))
	}
	b.WriteString("  Deploy ingress-nginx and your app, then point DNS (CNAME) to the LB hostname shown above.\n")

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

func createObjectStorage(ctx context.Context, exec commands.Executor, config *StarterKitConfig, network *upcloud.Network) (*upcloud.ManagedObjectStorage, *upcloud.ManagedObjectStorageUserAccessKey, string, error) {
	exec.PushProgressStarted("Deploying Object Storage")
	region, err := stack.GetObjectStorageRegionFromZone(exec, config.Zone)
	if err != nil {
		return nil, nil, "", fmt.Errorf("failed to validate object storage region from zone %q: %w, contact support because you might be using a new zone not supported by the deploy command", config.Zone, err)
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
		return nil, nil, "", fmt.Errorf("failed to create object storage: %w", err)
	}
	exec.PushProgressSuccess("Deploying Object Storage")

	objStorage, err = exec.All().WaitForManagedObjectStorageOperationalState(exec.Context(), &request.WaitForManagedObjectStorageOperationalStateRequest{
		DesiredState: upcloud.ManagedObjectStorageOperationalStateRunning,
		UUID:         objStorage.UUID,
	})
	if err != nil {
		return nil, nil, "", fmt.Errorf("error while waiting for the object storage to become online: %w", err)
	}

	// Create user for the object storage
	user, err := exec.All().CreateManagedObjectStorageUser(exec.Context(), &request.CreateManagedObjectStorageUserRequest{
		Username:    "starter-kit-user",
		ServiceUUID: objStorage.UUID,
	})
	if err != nil {
		return nil, nil, "", fmt.Errorf("failed to create user for the object storage: %w", err)
	}

	err = exec.All().AttachManagedObjectStorageUserPolicy(exec.Context(), &request.AttachManagedObjectStorageUserPolicyRequest{
		ServiceUUID: objStorage.UUID,
		Username:    user.Username,
		Name:        "ECSS3FullAccess",
	})
	if err != nil {
		return nil, nil, "", fmt.Errorf("failed to attach user policy to object storage: %w", err)
	}

	// Create an access key for the object storage
	userAccessKey, err := exec.All().CreateManagedObjectStorageUserAccessKey(exec.Context(), &request.CreateManagedObjectStorageUserAccessKeyRequest{
		Username:    user.Username,
		ServiceUUID: objStorage.UUID,
	})
	if err != nil {
		return nil, nil, "", fmt.Errorf("failed to create access key for the object storage: %w", err)
	}

	// Create a bucket for the object storage
	bucketName := "starter-kit-storage"
	_, err = exec.All().CreateManagedObjectStorageBucket(exec.Context(), &request.CreateManagedObjectStorageBucketRequest{
		Name:        bucketName,
		ServiceUUID: objStorage.UUID,
	})
	if err != nil {
		return nil, nil, "", fmt.Errorf("failed to create bucket for the object storage: %w", err)
	}

	return objStorage, userAccessKey, bucketName, nil
}
