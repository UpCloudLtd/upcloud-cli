package core

import (
	"embed"
	"fmt"
	"io/fs"
	"math/rand"
	"os"
	"path/filepath"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud/request"
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

// createNetwork creates a network with a random 10.0.X.0/24 subnet
// It will try 10 times to create a network with a random subnet
// If it fails to create a network after 10 attempts, it returns an error
func CreateNetwork(exec commands.Executor, networkName, location string) (*upcloud.Network, error) {
	var networkCreated = false
	var network *upcloud.Network

	for range 10 {
		// Generate random 10.0.X.0/24 subnet
		x := rand.Intn(254) + 1
		cidr := fmt.Sprintf("10.0.%d.0/24", x)
		var err error
		network, err = exec.Network().CreateNetwork(exec.Context(), &request.CreateNetworkRequest{
			Name:   networkName,
			Zone:   location,
			Router: "",
			Labels: []upcloud.Label{
				{Key: "stacks.upcloud.com/stack", Value: "supabase"},
				{Key: "stacks.upcloud.com/created-by", Value: "upctl"},
				{Key: "stacks.upcloud.com/chart-version", Value: "0.1.3"},
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
			return os.MkdirAll(targetPath, 0755)
		}

		// Make sure the parent directory exists before writing the file
		if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
			return err
		}

		data, err := fsys.ReadFile(path)
		if err != nil {
			return err
		}

		return os.WriteFile(targetPath, data, 0644)
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
