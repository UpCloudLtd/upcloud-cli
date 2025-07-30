package supabase

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands/stack/stackops"
	"gopkg.in/yaml.v3"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/storage/driver"
	"k8s.io/client-go/kubernetes"
)

func loadHelmCharts(chartPath string) (*chart.Chart, error) {
	// Load helm charts from disk
	ch, err := loader.Load(chartPath)
	if err != nil {
		return nil, fmt.Errorf("loading chart from %q: %w", chartPath, err)
	}

	// Remove test templates from the chart
	filteredTemplates := []*chart.File{}
	for _, f := range ch.Templates {
		if !strings.HasPrefix(f.Name, "templates/test/") {
			filteredTemplates = append(filteredTemplates, f)
		}
	}
	ch.Templates = filteredTemplates

	return ch, nil
}

// deployHelmRelease will check if a Helm release named releaseName exists in namespace=releaseName.
// - If it doesn't exist and upgrade = false, it will install the chart at chartPath using valuesFiles.
// - If it does exist and upgrade = true, it will perform a Helm upgrade with the same chartPath and values.
func DeployHelmRelease(kubeClient *kubernetes.Clientset, releaseName string, chartPath string, valuesFiles []string, upgrade bool) error {
	// Set HELM_NAMESPACE environment variable to ensure resources land in the correct namespace
	// Bug: https://github.com/helm/helm/issues/9171 and https://github.com/helm/helm/issues/8780
	os.Setenv("HELM_NAMESPACE", releaseName)

	// Wait for the Kubernetes API server to be ready
	err := stackops.WaitForAPIServer(kubeClient)
	if err != nil {
		return fmt.Errorf("waiting for API server: %w", err)
	}

	// Ensure the target namespace exists
	err = stackops.CreateNamespace(kubeClient, releaseName)
	if err != nil {
		return fmt.Errorf("ensuring namespace %q exists: %w", releaseName, err)
	}

	// Configure Helm logging to file
	logFile, err := stackops.CreateHelmLogFile(chartPath)
	if err != nil {
		return fmt.Errorf("creating Helm log file: %w", err)
	}
	defer logFile.Close()

	// Bootstrap Helm action config
	actionConfig, err := stackops.InitHelmActionConfig(releaseName, logFile)
	if err != nil {
		return fmt.Errorf("initializing Helm action config: %w", err)
	}

	// Load the Helm chart from the specified path
	ch, err := loadHelmCharts(chartPath)
	if err != nil {
		return fmt.Errorf("loading Helm chart from %q: %w", chartPath, err)
	}

	// Merge values from the specified files
	mergedVals, err := stackops.MergeValueFiles(valuesFiles)
	if err != nil {
		return fmt.Errorf("merging values files: %w", err)
	}

	// Check for existing release in the *correct* namespace
	statusClient := action.NewStatus(actionConfig)
	_, err = statusClient.Run(releaseName)
	if err != nil && !errors.Is(err, driver.ErrReleaseNotFound) {
		return fmt.Errorf("error checking release status: %w", err)
	}

	// Install if not found
	if errors.Is(err, driver.ErrReleaseNotFound) && !upgrade {
		fmt.Printf("Release %q not found: installing into namespace %q\n", releaseName, releaseName)
		installClient := action.NewInstall(actionConfig)
		installClient.ReleaseName = releaseName
		installClient.Namespace = releaseName
		installClient.CreateNamespace = false
		installClient.Wait = true
		installClient.Timeout = 10 * time.Minute

		if _, err := installClient.Run(ch, mergedVals); err != nil {
			return fmt.Errorf("helm install failed: %w", err)
		}
		fmt.Printf("Helm release %q successfully installed\n", releaseName)
		return nil
	}

	// Upgrade if it exists
	if upgrade {
		fmt.Printf("Release %q already exists: upgrading in namespace %q\n", releaseName, releaseName)
		upgradeClient := action.NewUpgrade(actionConfig)
		upgradeClient.Namespace = releaseName
		upgradeClient.Atomic = true
		upgradeClient.Wait = true
		upgradeClient.Timeout = 10 * time.Minute

		if _, err := upgradeClient.Run(releaseName, ch, mergedVals); err != nil {
			return fmt.Errorf("helm upgrade failed: %w", err)
		}
		fmt.Printf("Helm release %q successfully upgraded\n", releaseName)
	}
	return nil
}

// UpdateDNS updates the DNS entries in the values file with the provided dnsPrefix.
func UpdateDNS(valuesFile, updatedValuesFile, dnsPrefix string) error {
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
