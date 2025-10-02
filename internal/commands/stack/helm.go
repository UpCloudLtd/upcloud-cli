package stack

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"dario.cat/mergo"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/storage/driver"
	"k8s.io/client-go/kubernetes"
	"k8s.io/helm/pkg/chartutil"
)

func CreateHelmLogFile(chartPath string) (*os.File, error) {
	timestamp := time.Now().Format("20060102-150405")
	logDir := filepath.Join(chartPath, "logs-"+timestamp)
	logFileName := "helm-deploy.log"
	logFilePath := filepath.Join(logDir, logFileName)

	if err := os.MkdirAll(logDir, 0o755); err != nil {
		return nil, fmt.Errorf("failed to create log directory %q: %w", logDir, err)
	}

	logFile, err := os.Create(logFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to create log file %q: %w", logFilePath, err)
	}
	return logFile, nil
}

// initHelmActionConfig initializes the Helm action configuration with the provided chart path and release name.
func InitHelmActionConfig(releaseName string, logFile *os.File) (*action.Configuration, error) {
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

// loadHelmCharts loads a Helm chart from the specified path and filters out test templates.
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
	err := WaitForAPIServer(kubeClient)
	if err != nil {
		return fmt.Errorf("waiting for API server: %w", err)
	}

	// Ensure the target namespace exists
	err = CreateNamespace(kubeClient, releaseName)
	if err != nil {
		return fmt.Errorf("ensuring namespace %q exists: %w", releaseName, err)
	}

	// Configure Helm logging to file
	logFile, err := CreateHelmLogFile(chartPath)
	if err != nil {
		return fmt.Errorf("creating Helm log file: %w", err)
	}
	defer logFile.Close()

	// Bootstrap Helm action config
	actionConfig, err := InitHelmActionConfig(releaseName, logFile)
	if err != nil {
		return fmt.Errorf("initializing Helm action config: %w", err)
	}

	// Load the Helm chart from the specified path
	ch, err := loadHelmCharts(chartPath)
	if err != nil {
		return fmt.Errorf("loading Helm chart from %q: %w", chartPath, err)
	}

	// Merge values from the specified files
	mergedVals, err := MergeHelmValueFiles(valuesFiles)
	if err != nil {
		return fmt.Errorf("merging values files: %w", err)
	}

	// Check for existing release in the correct namespace
	statusClient := action.NewStatus(actionConfig)
	_, err = statusClient.Run(releaseName)
	if err != nil && !errors.Is(err, driver.ErrReleaseNotFound) {
		return fmt.Errorf("error checking release status: %w", err)
	}

	// Install if not found
	if errors.Is(err, driver.ErrReleaseNotFound) && !upgrade {
		installClient := action.NewInstall(actionConfig)
		installClient.ReleaseName = releaseName
		installClient.Namespace = releaseName
		installClient.CreateNamespace = false
		installClient.Wait = true
		installClient.Timeout = 15 * time.Minute

		if _, err := installClient.Run(ch, mergedVals); err != nil {
			return fmt.Errorf("helm install failed: %w", err)
		}

		return nil
	}

	// Upgrade if it exists
	if upgrade {
		upgradeClient := action.NewUpgrade(actionConfig)
		upgradeClient.Namespace = releaseName
		upgradeClient.Atomic = true
		upgradeClient.Wait = true
		upgradeClient.Timeout = 10 * time.Minute

		if _, err := upgradeClient.Run(releaseName, ch, mergedVals); err != nil {
			return fmt.Errorf("helm upgrade failed: %w", err)
		}
	}
	return nil
}

// loadHelmChartsFromRepo loads a Helm chart from a specified repository URL and chart name.
func loadHelmChartsFromRepo(actionConfig *action.Configuration, repoURL, chartName, version string) (*chart.Chart, error) {
	settings := cli.New()
	installClient := action.NewInstall(actionConfig)
	installClient.RepoURL = repoURL
	installClient.Version = version

	chartPath, err := installClient.LocateChart(chartName, settings)
	if err != nil {
		return nil, fmt.Errorf("locating chart %q in repo %q: %w", chartName, repoURL, err)
	}

	ch, err := loader.Load(chartPath)
	if err != nil {
		return nil, fmt.Errorf("loading chart from %q: %w", chartPath, err)
	}

	return ch, nil
}

// MergeHelmValueFiles merges multiple Helm values files into a single map.
// This function is used to combine values from multiple files, allowing for overrides.
func MergeHelmValueFiles(valuesFiles []string) (map[string]interface{}, error) {
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

// DeployHelmReleaseFromRepo deploys a Helm release from a remote repository.
func DeployHelmReleaseFromRepo(
	kubeClient *kubernetes.Clientset,
	releaseName, repoURL, chartName, version string,
	valsFiles []string,
	upgrade bool,
) error {
	os.Setenv("HELM_NAMESPACE", releaseName)

	// Ensure the target namespace exists
	err := CreateNamespace(kubeClient, releaseName)
	if err != nil {
		return fmt.Errorf("ensuring namespace %q exists: %w", releaseName, err)
	}

	// Prepare a log file for Helm debug output
	logFile, err := CreateHelmLogFile(filepath.Join(os.TempDir(), fmt.Sprintf("%s-helm.log", releaseName)))
	if err != nil {
		return fmt.Errorf("creating Helm log file: %w", err)
	}
	defer logFile.Close()

	// Initialize Helm action.Configuration
	actionConfig, err := InitHelmActionConfig(releaseName, logFile)
	if err != nil {
		return fmt.Errorf("initializing Helm action config failed")
	}

	// Locate and load the chart in the remote repo
	ch, err := loadHelmChartsFromRepo(actionConfig, repoURL, chartName, version)
	if err != nil {
		return fmt.Errorf("loading Helm chart from repo %q: %w", repoURL, err)
	}

	// Merge values from the specified files
	mergedVals, err := MergeHelmValueFiles(valsFiles)
	if err != nil {
		return fmt.Errorf("merging values files: %w", err)
	}

	// Check for existing release
	statusClient := action.NewStatus(actionConfig)
	_, err = statusClient.Run(releaseName)
	if err != nil && !errors.Is(err, driver.ErrReleaseNotFound) {
		return fmt.Errorf("checking release %q status: %w", releaseName, err)
	}

	// Install if missing
	if errors.Is(err, driver.ErrReleaseNotFound) && !upgrade {
		installClient := action.NewInstall(actionConfig)
		installClient.ReleaseName = releaseName
		installClient.Namespace = releaseName
		installClient.CreateNamespace = false
		installClient.Wait = true
		installClient.Timeout = 10 * time.Minute

		fmt.Fprintf(logFile, "Installing %q into %q\n", releaseName, releaseName)
		if _, err := installClient.Run(ch, mergedVals); err != nil {
			return fmt.Errorf("helm install failed: %w", err)
		}
		fmt.Fprintf(logFile, "Installed release %q\n", releaseName)
		return nil
	}

	// Upgrade if requested
	if upgrade {
		upgradeClient := action.NewUpgrade(actionConfig)
		upgradeClient.Namespace = releaseName
		upgradeClient.Atomic = true
		upgradeClient.Wait = true
		upgradeClient.Timeout = 10 * time.Minute

		fmt.Fprintf(logFile, "Upgrading %q in %q\n", releaseName, releaseName)
		if _, err := upgradeClient.Run(releaseName, ch, mergedVals); err != nil {
			return fmt.Errorf("helm upgrade failed: %w", err)
		}
		fmt.Fprintf(logFile, "Upgraded release %q\n", releaseName)
	}

	return nil
}

// UninstallHelmRelease uninstalls a Helm release in namespace=releaseName.
func UninstallHelmRelease(releaseName, logDir string) error {
	// return error if KUBECONFIG is not set
	if os.Getenv("KUBECONFIG") == "" {
		return errors.New("KUBECONFIG environment variable is not set")
	}

	if err := os.Setenv("HELM_NAMESPACE", releaseName); err != nil {
		return fmt.Errorf("set HELM_NAMESPACE: %w", err)
	}

	// Ensure logs are written to the same chartPath dir
	logFile, err := CreateHelmLogFile(logDir)
	if err != nil {
		return fmt.Errorf("creating Helm log file: %w", err)
	}
	defer logFile.Close()

	// Initialize Helm action config
	actionConfig, err := InitHelmActionConfig(releaseName, logFile)
	if err != nil {
		return fmt.Errorf("initializing Helm action config: %w", err)
	}

	// Prepare uninstall client
	uninstall := action.NewUninstall(actionConfig)
	uninstall.Wait = true
	uninstall.Timeout = 15 * time.Minute

	// Run uninstall
	resp, err := uninstall.Run(releaseName)
	if err != nil {
		return fmt.Errorf("uninstalling release %q: %w", releaseName, err)
	}

	if resp != nil {
		fmt.Fprintf(logFile, "Uninstalled release %q: %s\n", releaseName, resp.Info)
	}

	return nil
}
