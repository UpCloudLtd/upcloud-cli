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
	"k8s.io/client-go/tools/clientcmd"
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

func waitForAPIServer(kubeconfigPath string) error {
	cfg, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	if err != nil {
		return fmt.Errorf("loading kubeconfig: %w", err)
	}
	clientset, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		return fmt.Errorf("building clientset: %w", err)
	}

	for i := 1; i <= 30; i++ {
		if _, err := clientset.Discovery().ServerVersion(); err == nil {
			fmt.Printf("API ready after %d attempts\n", i)
			return nil
		} else {
			fmt.Printf("⏳ [%2d/30] still waiting for API… %v\n", i, err)
			time.Sleep(5 * time.Second)
		}
	}
	return fmt.Errorf("timed out waiting for API server")
}

// deployHelmRelease will check if a Helm release named releaseName exists in namespace=releaseName.
// - If it doesn't exist and upgrade = false, it will install the chart at chartPath using valuesFiles.
// - If it does exist and upgrade = true, it will perform a Helm upgrade with the same chartPath and values.
func DeployHelmRelease(releaseName string, chartPath string, valuesFiles []string, upgrade bool) error {
	kubeconfigPath := os.Getenv("KUBECONFIG")

	// Wait for the Kubernetes API server to be ready
	err := waitForAPIServer(kubeconfigPath) // Use the loaded kubeconfig path
	if err != nil {
		return fmt.Errorf("waiting for API server: %w", err)
	}

	// Load kubeconfig and create Kubernetes client
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	if err != nil {
		return fmt.Errorf("loading kubeconfig: %w", err)
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return fmt.Errorf("building clientset: %w", err)
	}

	// Ensure the target namespace exists
	fmt.Printf("Ensuring namespace %q exists...\n", releaseName)
	_, err = clientset.CoreV1().Namespaces().Get(context.Background(), releaseName, v1.GetOptions{})
	if err != nil {
		// Correctly check for StatusError and StatusReasonNotFound
		if apierrors.IsNotFound(err) {
			fmt.Printf("Namespace %q not found, creating...\n", releaseName)
			_, createErr := clientset.CoreV1().Namespaces().Create(context.Background(), &corev1.Namespace{
				ObjectMeta: v1.ObjectMeta{
					Name: releaseName,
				},
			}, v1.CreateOptions{})
			if createErr != nil {
				return fmt.Errorf("failed to create namespace %q: %w", releaseName, createErr)
			}
			fmt.Printf("Namespace %q created successfully.\n", releaseName)
		} else {
			return fmt.Errorf("error getting namespace %q: %w", releaseName, err)
		}
	} else {
		fmt.Printf("Namespace %q already exists.\n", releaseName)
	}

	// Set HELM_NAMESPACE environment variable to ensure resources land in the correct namespace
	// This is needed because otherwise the resources end up in the default namespace even if the actionConfig is set to the target namespace.
	// https://github.com/helm/helm/issues/9171 and https://github.com/helm/helm/issues/8780
	os.Setenv("HELM_NAMESPACE", releaseName)
	fmt.Printf("Set HELM_NAMESPACE to %q\n", releaseName)

	// Bootstrap Helm action config
	settings := cli.New()
	actionConfig := new(action.Configuration)
	driverName := os.Getenv("HELM_DRIVER")
	if err := actionConfig.Init(
		settings.RESTClientGetter(),
		releaseName,
		driverName,
		func(format string, v ...interface{}) {
			fmt.Printf(format, v...)
		},
	); err != nil {
		return fmt.Errorf("initializing helm action config: %w", err)
	}

	// Load the chart from disk
	ch, err := loader.Load(chartPath)
	if err != nil {
		return fmt.Errorf("loading chart from %q: %w", chartPath, err)
	}

	// Remove test templates from the chart
	filteredTemplates := []*chart.File{}
	for _, f := range ch.Templates {
		if !strings.HasPrefix(f.Name, "templates/test/") {
			filteredTemplates = append(filteredTemplates, f)
		}
	}
	ch.Templates = filteredTemplates

	// Merge values files with override
	mergedVals := map[string]interface{}{}
	for _, vf := range valuesFiles {
		vals, err := chartutil.ReadValuesFile(vf)
		if err != nil {
			return fmt.Errorf("reading values file %q: %w", vf, err)
		}

		valsMap := vals.AsMap()

		if err := mergo.Merge(&mergedVals, valsMap, mergo.WithOverride, mergo.WithTypeCheck); err != nil {
			fmt.Printf("error merging values from file %s: %v\n", vf, err)
			return fmt.Errorf("merging values from file %q: %w", vf, err)
		}
	}

	// --- START: Added code for dry-run debugging ---
	/*
		fmt.Println("\n--- Performing Helm Dry Run to inspect rendered manifests ---")
		dryRunClient := action.NewInstall(actionConfig)
		dryRunClient.ReleaseName = releaseName
		dryRunClient.Namespace = releaseName
		dryRunClient.DryRun = true       // Crucially, set dry run to true
		dryRunClient.ClientOnly = true   // Render client-side only (no API calls)
		dryRunClient.IncludeCRDs = false // Exclude CRDs for cleaner output if not relevant to namespace issue

		renderedRelease, err := dryRunClient.Run(ch, mergedVals)
		if err != nil {
			return fmt.Errorf("helm dry-run failed: %w", err)
		}
		fmt.Printf("--- START: Rendered Manifests for %s in namespace %s ---\n", releaseName, releaseName)
		fmt.Println(renderedRelease.Manifest)
		fmt.Printf("--- END: Rendered Manifests ---\n")
		fmt.Println("--- Dry Run Complete ---")
	*/
	// --- END: Added code for dry-run debugging ---

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

func WaitForLoadBalancer(namespace, serviceName string, maxRetries int, sleepInterval time.Duration) (string, error) {
	// Load kubeconfig from KUBECONFIG env var
	kubeconfigPath := os.Getenv("KUBECONFIG")
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	if err != nil {
		return "", fmt.Errorf("failed to load kubeconfig: %w", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return "", fmt.Errorf("failed to create Kubernetes client: %w", err)
	}

	fmt.Println("Waiting for Kong LoadBalancer external hostname...")

	for i := 1; i <= maxRetries; i++ {
		svc, err := clientset.CoreV1().Services(namespace).Get(context.Background(), serviceName, v1.GetOptions{})
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

	// Now we are at the map where the final key should be updated
	lastKey := path[len(path)-1]
	for i := 0; i < len(current.Content); i += 2 {
		k := current.Content[i]
		v := current.Content[i+1]
		if k.Value == lastKey {
			v.Value = fmt.Sprintf("http://%s.%s", strings.ToLower(lastKey), dnsPrefix+".upcloudlb.com")
			return true
		}
	}

	return false
}
