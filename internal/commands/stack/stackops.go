package stack

import (
	"context"
	"crypto/rand"
	"embed"
	"fmt"
	"io/fs"
	"math/big"
	"os"
	"path/filepath"
	"time"

	"github.com/UpCloudLtd/progress/messages"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands/all"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands/kubernetes"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands/network"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands/objectstorage"
	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud/request"
	"golang.org/x/sync/errgroup"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type StackType string

const (
	StackTypeDokku      StackType = "dokku"
	StackTypeSupabase   StackType = "supabase"
	StackTypeStarterKit StackType = "starter"
)

type Version string

const (
	VersionV0_1_0_0 Version = "v0.1.0.0"
	SupportEmail    string  = "support@upcloud.com"
)

type (
	SupabaseResourceRootName   string
	DokkuResourceRootName      string
	StarterKitResourceRootName string
)

var (
	SupabaseResourceRootNameCluster      SupabaseResourceRootName   = SupabaseResourceRootName(fmt.Sprintf("stack-%s-cluster", StackTypeSupabase))
	SupabaseResourceRootNameNetwork      SupabaseResourceRootName   = SupabaseResourceRootName(fmt.Sprintf("stack-%s-net", StackTypeSupabase))
	SupabaseResourceRootNameObjStorage   SupabaseResourceRootName   = SupabaseResourceRootName(fmt.Sprintf("stack-%s-os", StackTypeSupabase))
	DokkuResourceRootNameCluster         DokkuResourceRootName      = DokkuResourceRootName(fmt.Sprintf("stack-%s-cluster", StackTypeDokku))
	DokkuResourceRootNameNetwork         DokkuResourceRootName      = DokkuResourceRootName(fmt.Sprintf("stack-%s-net", StackTypeDokku))
	StarterKitResourceRootNameCluster    StarterKitResourceRootName = StarterKitResourceRootName(fmt.Sprintf("stack-%s-cluster", StackTypeStarterKit))
	StarterKitResourceRootNameNetwork    StarterKitResourceRootName = StarterKitResourceRootName(fmt.Sprintf("stack-%s-net", StackTypeStarterKit))
	StarterKitResourceRootNameObjStorage StarterKitResourceRootName = StarterKitResourceRootName(fmt.Sprintf("stack-%s-obj-sto", StackTypeStarterKit))
	StarterKitResourceRootNameDatabase   StarterKitResourceRootName = StarterKitResourceRootName(fmt.Sprintf("stack-%s-db", StackTypeStarterKit))
	StarterKitResourceRootNameRouter     StarterKitResourceRootName = StarterKitResourceRootName(fmt.Sprintf("stack-%s-router", StackTypeStarterKit))
)

// GetNetworkFromName retrieves the network from the net name, returns nil if it does not exist
func GetNetworkFromName(networkName string, networks []upcloud.Network) *upcloud.Network {
	for _, network := range networks {
		if network.Name == networkName {
			return &network
		}
	}
	return nil
}

// ClusterExists checks if a Kubernetes cluster with the given name exists in the provided list of clusters.
func ClusterExists(clusterName string, clusters []upcloud.KubernetesCluster) bool {
	for _, cluster := range clusters {
		if cluster.Name == clusterName {
			return true
		}
	}
	return false
}

// CreateNetwork creates a network with a random 10.0.X.0/24 subnet
// It will try 10 times to create a network with a random subnet
// If it fails to create a network after 10 attempts, it returns an error
func CreateNetwork(exec commands.Executor, networkName, location string, stackType StackType) (*upcloud.Network, error) {
	networkCreated := false
	var network *upcloud.Network

	for range 10 {
		// Generate random 10.0.X.0/24 subnet
		r, _ := rand.Int(rand.Reader, big.NewInt(254))
		x := r.Int64() + 1
		cidr := fmt.Sprintf("10.0.%d.0/24", x)
		var err error
		network, err = exec.Network().CreateNetwork(exec.Context(), &request.CreateNetworkRequest{
			Name:   networkName,
			Zone:   location,
			Router: "",
			Labels: []upcloud.Label{
				{Key: "stacks.upcloud.com/stack", Value: string(stackType)},
				{Key: "stacks.upcloud.com/created-by", Value: "upctl"},
				{Key: "stacks.upcloud.com/version", Value: string(VersionV0_1_0_0)},
				{Key: "stacks.upcloud.com/name", Value: networkName},
			},
			IPNetworks: []upcloud.IPNetwork{{Address: cidr, DHCP: upcloud.True, Family: upcloud.IPAddressFamilyIPv4}},
		})

		if err != nil {
			continue
		} else {
			networkCreated = true
			break
		}
	}
	if !networkCreated {
		return nil, fmt.Errorf("failed to create network after 10 attempts")
	}
	return network, nil
}

// ExtractFolder extracts the embedded chart files to the specified target directory.
func ExtractFolder(fsys embed.FS, targetDir string) error {
	return fs.WalkDir(fsys, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		targetPath := filepath.Join(targetDir, path)

		if d.IsDir() {
			return os.MkdirAll(targetPath, 0o755)
		}

		// Make sure the parent directory exists before writing the file
		if err := os.MkdirAll(filepath.Dir(targetPath), 0o755); err != nil {
			return err
		}

		data, err := fsys.ReadFile(path)
		if err != nil {
			return err
		}

		return os.WriteFile(targetPath, data, 0o600)
	})
}

// CheckSSHKeys ensures both private and public SSH key files exist.
func CheckSSHKeys(privateKeyPath, publicKeyPath string) error {
	for _, p := range []string{privateKeyPath, publicKeyPath} {
		if stat, err := os.Stat(p); err != nil {
			return fmt.Errorf("SSH key not found at %s: %w", p, err)
		} else if stat.IsDir() {
			return fmt.Errorf("SSH key path %s is a directory", p)
		}
	}
	return nil
}

// GetObjectStorageRegionFromZone retrieves the object storage region for the selected zone. First occurrence is returned
func GetObjectStorageRegionFromZone(exec commands.Executor, zone string) (string, error) {
	regions, err := exec.All().GetManagedObjectStorageRegions(exec.Context(), &request.GetManagedObjectStorageRegionsRequest{})
	if err != nil {
		return "", fmt.Errorf("failed to get object storage regions: %w", err)
	}

	// find the region that contains the zone
	for i := range regions {
		// check primary zone
		if regions[i].PrimaryZone == zone {
			return regions[i].Name, nil
		}
		// check secondary zones
		for _, z := range regions[i].Zones {
			if z.Name == zone {
				return regions[i].Name, nil
			}
		}
	}
	return "", fmt.Errorf("zone %s not found in any object storage region", zone)
}

// waitForManagedObjectState waits for database to reach given state and updates progress message with key matching given msg. Finally, progress message is updated back to given msg and either done state or timeout warning.
func WaitForManagedObjectStorageState(uuid string, state upcloud.ManagedObjectStorageOperationalState, exec commands.Executor, msg string) {
	exec.PushProgressUpdateMessage(msg, fmt.Sprintf("Waiting for object storage %s to be in %s state", uuid, state))

	ctx, cancel := context.WithTimeout(exec.Context(), 15*time.Minute)
	defer cancel()

	if _, err := exec.All().WaitForManagedObjectStorageOperationalState(ctx, &request.WaitForManagedObjectStorageOperationalStateRequest{
		UUID:         uuid,
		DesiredState: state,
	}); err != nil {
		exec.PushProgressUpdate(messages.Update{
			Key:     msg,
			Message: msg,
			Status:  messages.MessageStatusWarning,
			Details: "Error: " + err.Error(),
		})
		return
	}

	exec.PushProgressUpdateMessage(msg, msg)
	exec.PushProgressSuccess(msg)
}

func DestroyStack(exec commands.Executor, name, zone string, deleteStorage, deleteObjectStorage bool, stackType StackType) error {
	switch stackType {
	case StackTypeSupabase:
		logDir := os.TempDir()
		clusterName := fmt.Sprintf("%s-%s-%s", SupabaseResourceRootNameCluster, name, zone)
		msg := fmt.Sprintf("Searching cluster %s in zone %s", clusterName, zone)
		exec.PushProgressStarted(msg)
		clusters, err := exec.All().GetKubernetesClusters(exec.Context(), &request.GetKubernetesClustersRequest{})
		if err != nil {
			return fmt.Errorf("failed to get kubernetes clusters: %w", err)
		}

		var cluster *upcloud.KubernetesCluster
		for _, cl := range clusters {
			if cl.Name == clusterName {
				cluster = &cl
				break
			}
		}

		if cluster == nil {
			return fmt.Errorf("a cluster with the name '%s' was not found", clusterName)
		}

		exec.PushProgressSuccess(msg)

		msg = fmt.Sprintln("Preparing to start deleting resources")
		exec.PushProgressStarted(msg)

		kubeconfigPath, err := WriteKubeconfigToFile(exec, cluster.UUID, logDir)
		if err != nil {
			return fmt.Errorf("failed to write kubeconfig to file: %w", err)
		}

		if err := os.Setenv("KUBECONFIG", kubeconfigPath); err != nil {
			return fmt.Errorf("set KUBECONFIG: %w", err)
		}

		exec.PushProgressSuccess(msg)

		kubeClient, err := GetKubernetesClient(kubeconfigPath)
		if err != nil {
			return fmt.Errorf("failed to create Kubernetes client: %w", err)
		}

		var uuids []string
		// Collect PVC volume UUIDs for deletion if deleteStorage is true
		// This is done before uninstalling the helm release to ensure we get all PVCs
		// associated with the Supabase stack
		// If deleteStorage is false, we skip this step and leave the PVCs intact
		if deleteStorage {
			msg = fmt.Sprintln("Collecting PVC volume UUIDs for deletion")
			exec.PushProgressStarted(msg)

			uuids, err = CollectPVCVolumeUUIDs(exec.Context(), exec, kubeClient, clusterName)
			if err != nil {
				return fmt.Errorf("failed to retrieve PVC volume UUIDs: %w", err)
			}

			exec.PushProgressSuccess(msg)
		}

		msg = fmt.Sprintf("Uninstalling Supabase helm release: %s", clusterName)
		exec.PushProgressStarted(msg)

		err = UninstallHelmRelease(clusterName, logDir)
		if err != nil {
			return fmt.Errorf("failed to uninstall helm release: %w", err)
		}

		exec.PushProgressSuccess(msg)

		if deleteStorage && len(uuids) > 0 {
			msg = fmt.Sprintln("Deleting PVC volumes")
			exec.PushProgressStarted(msg)

			err = DeletePVCVolumesByUUIDs(exec, uuids)
			if err != nil {
				return fmt.Errorf("failed to delete PVC volumes: %w", err)
			}
			exec.PushProgressSuccess(msg)
		}

		msg = fmt.Sprintf("Deleting namespace: %s", clusterName)
		exec.PushProgressStarted(msg)

		err = kubeClient.CoreV1().Namespaces().Delete(exec.Context(), clusterName, v1.DeleteOptions{})
		if err != nil {
			return fmt.Errorf("failed to delete namespace: %w", err)
		}

		// Wait until namespace is deleted
		err = WaitForNamespaceDeletion(exec.Context(), kubeClient, clusterName, 10*time.Minute)
		if err != nil {
			return fmt.Errorf("failed to wait for namespace deletion: %w", err)
		}

		exec.PushProgressSuccess(msg)

		g, _ := errgroup.WithContext(exec.Context())
		g.Go(func() error {
			msg = fmt.Sprintf("Deleting Kubernetes cluster %s in zone %s", clusterName, zone)
			exec.PushProgressStarted(msg)

			_, err = kubernetes.Delete(exec, cluster.UUID, true)
			if err != nil {
				return fmt.Errorf("failed to delete kubernetes cluster: %w", err)
			}

			exec.PushProgressSuccess(msg)
			return nil
		})

		// Delete the object storage if it exists
		g.Go(func() error {
			objectStorageName := fmt.Sprintf("%s-%s-%s", SupabaseResourceRootNameObjStorage, name, zone)
			storages, err := exec.All().GetManagedObjectStorages(exec.Context(), &request.GetManagedObjectStoragesRequest{})
			if err != nil {
				return fmt.Errorf("failed to get object storages: %w", err)
			}

			var storageUUID string
			for _, sto := range storages {
				if sto.Name == objectStorageName {
					storageUUID = sto.UUID
					break
				}
			}

			if storageUUID == "" {
				exec.PushProgressUpdateMessage(msg, fmt.Sprintf("Object storage %s not found, skipping deletion", objectStorageName))
				return nil
			}

			if deleteObjectStorage {
				_, err = objectstorage.Delete(exec, storageUUID, true, true, true, true)
				if err != nil {
					return fmt.Errorf("failed to delete object storage: %w", err)
				}

			} else {
				exec.PushProgressUpdateMessage(msg, fmt.Sprintf("Detaching Object Storage %s from network", objectStorageName))
				objStorageNetworks, err := exec.All().GetManagedObjectStorageNetworks(exec.Context(), &request.GetManagedObjectStorageNetworksRequest{
					ServiceUUID: storageUUID,
				})
				if err != nil {
					return fmt.Errorf("failed to get object storage networks: %w", err)
				}

				for _, n := range objStorageNetworks {
					if n.UUID != nil {
						exec.PushProgressUpdateMessage(msg, fmt.Sprintf("Detaching network %s with UUID: %s from object storage %s", n.Name, *n.UUID, objectStorageName))
						err = exec.All().DeleteManagedObjectStorageNetwork(exec.Context(), &request.DeleteManagedObjectStorageNetworkRequest{
							ServiceUUID: storageUUID,
							NetworkName: n.Name,
						})
						if err != nil {
							msg = "failed to delete object storage network"
							return fmt.Errorf("failed to delete object storage networks: %w", err)
						}

					}
				}
			}
			return nil
		})

		if err := g.Wait(); err != nil {
			return err
		}

		msg = fmt.Sprintf("Deleting network %s in zone %s", clusterName, zone)
		exec.PushProgressStarted(msg)

		networkName := fmt.Sprintf("%s-%s-%s", SupabaseResourceRootNameNetwork, name, zone)
		networks, err := exec.All().GetNetworks(exec.Context())
		if err != nil {
			return fmt.Errorf("failed to get networks: %w", err)
		}

		for _, net := range networks.Networks {
			if net.Name == networkName {
				msg := fmt.Sprintf("Deleting network %s in zone %s", networkName, zone)
				exec.PushProgressStarted(msg)
				_, err = network.Delete(exec, net.UUID)
				if err != nil {
					return fmt.Errorf("failed to delete network: %w", err)
				}
				exec.PushProgressSuccess(msg)
				break
			}
		}
		exec.PushProgressSuccess(msg)
	case StackTypeDokku:
		clusterName := fmt.Sprintf("%s-%s-%s", DokkuResourceRootNameCluster, name, zone)
		networkName := fmt.Sprintf("%s-%s-%s", DokkuResourceRootNameNetwork, name, zone)

		resources, err := all.ListResources(exec, []string{clusterName, networkName}, []string{})
		if err != nil {
			return err
		}

		err = all.DeleteResources(exec, resources, 16)
		if err != nil {
			return err
		}
	case StackTypeStarterKit:
		clusterName := fmt.Sprintf("%s-%s-%s", StarterKitResourceRootNameCluster, name, zone)
		networkName := fmt.Sprintf("%s-%s-%s", StarterKitResourceRootNameNetwork, name, zone)
		objectStorageName := fmt.Sprintf("%s-%s-%s", StarterKitResourceRootNameObjStorage, name, zone)
		dbName := fmt.Sprintf("%s-%s-%s", StarterKitResourceRootNameDatabase, name, zone)
		routerName := fmt.Sprintf("%s-%s-%s", StarterKitResourceRootNameRouter, name, zone)

		resources, err := all.ListResources(exec, []string{clusterName, objectStorageName, dbName, routerName, networkName}, []string{})
		if err != nil {
			return err
		}

		err = all.DeleteResources(exec, resources, 16)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("unsupported stack type: %s", stackType)
	}
	return nil
}

func waitForStorageToBeDetached(exec commands.Executor, storageUUID string, timeout time.Duration) error {
	start := time.Now()
	for {
		storage, err := exec.All().GetStorageDetails(exec.Context(), &request.GetStorageDetailsRequest{UUID: storageUUID})
		if err != nil {
			return fmt.Errorf("get storage details: %w", err)
		}

		// UpCloud storages can only be deleted when detached and in 'online' state
		if storage.State == upcloud.StorageStateOnline && len(storage.ServerUUIDs) == 0 {
			return nil
		}

		if time.Since(start) > timeout {
			return fmt.Errorf(
				"timeout waiting for storage %s to become deletable (state=%s, servers=%v)",
				storageUUID, storage.State, storage.ServerUUIDs,
			)
		}

		time.Sleep(5 * time.Second)
	}
}
