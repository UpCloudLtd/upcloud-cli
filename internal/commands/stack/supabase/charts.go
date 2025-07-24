package supabase

import (
	"context"
	"embed"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"

	"dario.cat/mergo"
	"gopkg.in/yaml.v3"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/chartutil"
	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/storage/driver"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// ExtractChart extracts the embedded Supabase chart files to the specified target directory.
func ExtractChart(fsys embed.FS, targetDir string) error {
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

func waitForAPIServer(kubeClient *kubernetes.Clientset) error {
	for i := 1; i <= 30; i++ {
		if _, err := kubeClient.Discovery().ServerVersion(); err == nil {
			fmt.Printf("API ready after %d attempts\n", i)
			return nil
		} else {
			fmt.Printf("⏳ [%2d/30] still waiting for API… %v\n", i, err)
			time.Sleep(5 * time.Second)
		}
	}
	return fmt.Errorf("timed out waiting for API server")
}

func createNamespace(kubeClient *kubernetes.Clientset, namespace string) error {
	fmt.Printf("Ensuring namespace %q exists...\n", namespace)
	_, err := kubeClient.CoreV1().Namespaces().Get(context.Background(), namespace, v1.GetOptions{})
	if err != nil {
		// Check for StatusError and StatusReasonNotFound
		if apierrors.IsNotFound(err) {
			fmt.Printf("Namespace %q not found, creating...\n", namespace)
			_, createErr := kubeClient.CoreV1().Namespaces().Create(context.Background(), &corev1.Namespace{
				ObjectMeta: v1.ObjectMeta{
					Name: namespace,
				},
			}, v1.CreateOptions{})
			if createErr != nil {
				return fmt.Errorf("failed to create namespace %q: %w", namespace, createErr)
			}
			fmt.Printf("Namespace %q created successfully.\n", namespace)
		} else {
			return fmt.Errorf("error getting namespace %q: %w", namespace, err)
		}
	} else {
		fmt.Printf("Namespace %q already exists.\n", namespace)
	}
	return nil
}

func createHelmLogFile(chartPath string) (*os.File, error) {
	timestamp := time.Now().Format("20060102-150405")
	logDir := filepath.Join(chartPath, "logs-"+timestamp)
	logFileName := "deploy.log"
	logFilePath := filepath.Join(logDir, logFileName)

	if err := os.MkdirAll(logDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create log directory %q: %w", logDir, err)
	}

	logFile, err := os.Create(logFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to create log file %q: %w", logFilePath, err)
	}
	return logFile, nil
}

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

func mergeValueFiles(valuesFiles []string) (map[string]interface{}, error) {
	// Merge values files with override
	mergedVals := map[string]interface{}{}
	for _, vf := range valuesFiles {
		vals, err := chartutil.ReadValuesFile(vf)
		if err != nil {
			return nil, fmt.Errorf("reading values file %q: %w", vf, err)
		}

		valsMap := vals.AsMap()

		if err := mergo.Merge(&mergedVals, valsMap, mergo.WithOverride, mergo.WithTypeCheck); err != nil {
			return nil, fmt.Errorf("merging values from file %q: %w", vf, err)
		}
	}
	return mergedVals, nil
}

// initHelmActionConfig initializes the Helm action configuration with the provided chart path and release name.
func initHelmActionConfig(chartPath, releaseName string, logFile *os.File) (*action.Configuration, error) {
	settings := cli.New()
	actionConfig := new(action.Configuration)
	driverName := os.Getenv("HELM_DRIVER")
	if err := actionConfig.Init(
		settings.RESTClientGetter(),
		releaseName,
		driverName,
		func(format string, v ...interface{}) {
			msg := fmt.Sprintf(format, v...)

			// Write all messages to the log file
			if _, writeErr := logFile.WriteString(msg + "\n"); writeErr != nil {
				// Handle error writing to log file, but don't fail the main process
				fmt.Fprintf(os.Stderr, "Error writing to log file: %v\n", writeErr)
			}
		},
	); err != nil {
		return nil, fmt.Errorf("initializing helm action config: %w", err)
	}

	return actionConfig, nil
}

// deployHelmRelease will check if a Helm release named releaseName exists in namespace=releaseName.
// - If it doesn't exist and upgrade = false, it will install the chart at chartPath using valuesFiles.
// - If it does exist and upgrade = true, it will perform a Helm upgrade with the same chartPath and values.
func DeployHelmRelease(kubeClient *kubernetes.Clientset, releaseName string, chartPath string, valuesFiles []string, upgrade bool) error {
	// Set HELM_NAMESPACE environment variable to ensure resources land in the correct namespace
	// Bug: https://github.com/helm/helm/issues/9171 and https://github.com/helm/helm/issues/8780
	os.Setenv("HELM_NAMESPACE", releaseName)

	// Wait for the Kubernetes API server to be ready
	err := waitForAPIServer(kubeClient)
	if err != nil {
		return fmt.Errorf("waiting for API server: %w", err)
	}

	// Ensure the target namespace exists
	err = createNamespace(kubeClient, releaseName)
	if err != nil {
		return fmt.Errorf("ensuring namespace %q exists: %w", releaseName, err)
	}

	// Configure Helm logging to file
	logFile, err := createHelmLogFile(chartPath)
	if err != nil {
		return fmt.Errorf("creating Helm log file: %w", err)
	}
	defer logFile.Close()

	// Bootstrap Helm action config
	actionConfig, err := initHelmActionConfig(chartPath, releaseName, logFile)
	if err != nil {
		return fmt.Errorf("initializing Helm action config: %w", err)
	}

	// Load the Helm chart from the specified path
	ch, err := loadHelmCharts(chartPath)
	if err != nil {
		return fmt.Errorf("loading Helm chart from %q: %w", chartPath, err)
	}

	// Merge values from the specified files
	mergedVals, err := mergeValueFiles(valuesFiles)
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

func WaitForLoadBalancer(kubeClient *kubernetes.Clientset, namespace, serviceName string, maxRetries int, sleepInterval time.Duration) (string, error) {
	fmt.Println("Waiting for Kong LoadBalancer external hostname...")

	for i := 1; i <= maxRetries; i++ {
		svc, err := kubeClient.CoreV1().Services(namespace).Get(context.Background(), serviceName, v1.GetOptions{})
		if err != nil {
			fmt.Printf("Error getting service %s/%s: %v\n", namespace, serviceName, err)
		} else if len(svc.Status.LoadBalancer.Ingress) > 0 {
			hostname := svc.Status.LoadBalancer.Ingress[0].Hostname
			if hostname != "" {
				fmt.Println("Found LoadBalancer Hostname:", hostname)
				return hostname, nil
			}
		}

		fmt.Printf("⏳ Waiting for LoadBalancer... (%d/%d)\n", i, maxRetries)
		time.Sleep(sleepInterval)
	}

	return "", fmt.Errorf("timed out waiting for LoadBalancer hostname after %d attempts", maxRetries)
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
