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
	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud/request"
)

type StackType string

const (
	StackTypeDokku      StackType = "dokku"
	StackTypeSupabase   StackType = "supabase"
	StackTypeStarterKit StackType = "starter-kit"
)

type Version string

const (
	VersionV0_1_0_0 Version = "v0.1.0.0"
	SupportEmail    string  = "support@upcloud.com"
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
