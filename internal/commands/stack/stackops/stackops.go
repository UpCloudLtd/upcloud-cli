package stackops

import (
	"bytes"
	"context"
	"embed"
	"fmt"
	"io"
	"io/fs"
	"math/rand"
	"os"
	"path/filepath"
	"time"

	"dario.cat/mergo"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud/request"
	"github.com/pkg/errors"
	"helm.sh/helm/pkg/chartutil"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/storage/driver"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/client-go/discovery"
	memorycached "k8s.io/client-go/discovery/cached/memory"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/restmapper"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/remotecommand"
	"k8s.io/kubectl/pkg/scheme"
	"sigs.k8s.io/kustomize/api/krusty"
	"sigs.k8s.io/kustomize/kyaml/filesys"
)

func WaitForLoadBalancer(kubeClient *kubernetes.Clientset, namespace, serviceName string, maxRetries int, sleepInterval time.Duration) (string, error) {
	fmt.Println("Waiting for LoadBalancer external hostname...")

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

// ExtractChart extracts the embedded chart files to the specified target directory.
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

func ClusterExists(clusterName string, clusters []upcloud.KubernetesCluster) bool {
	for _, cluster := range clusters {
		if cluster.Name == clusterName {
			return true
		}
	}
	return false
}

// getNetworkFromName retrieves the network from the net name, returns nil if it does not exist
func GetNetworkFromName(networkName string, networks []upcloud.Network) *upcloud.Network {
	for _, network := range networks {
		if network.Name == networkName {
			return &network
		}
	}
	return nil
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
		fmt.Println("Trying to create network with CIDR:", cidr)
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
			fmt.Printf("Failed to create network %s with CIDR %s: %v\n", networkName, cidr, err)
			continue
		} else {
			fmt.Println("Network created successfully:", network.Name, "with network.UUID", network.UUID)
			networkCreated = true
			break
		}
	}
	if !networkCreated {
		return nil, fmt.Errorf("failed to create network after 10 attempts")
	}
	return network, nil
}

func GetKubernetesClient(kubeconfigPath string) (*kubernetes.Clientset, error) {
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load kubeconfig: %w", err)
	}

	kubeClient, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create Kubernetes client: %w", err)
	}

	return kubeClient, nil
}

func WaitForAPIServer(kubeClient *kubernetes.Clientset) error {
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

func CreateNamespace(kubeClient *kubernetes.Clientset, namespace string) error {
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

func CreateHelmLogFile(chartPath string) (*os.File, error) {
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

func loadHelmChartsFromRepo(actionConfig *action.Configuration, repoURL, chartName, version string) (*chart.Chart, error) {
	settings := cli.New()
	installClient := action.NewInstall(actionConfig)
	installClient.ChartPathOptions.RepoURL = repoURL
	installClient.ChartPathOptions.Version = version

	chartPath, err := installClient.ChartPathOptions.LocateChart(chartName, settings)
	if err != nil {
		return nil, fmt.Errorf("locating chart %q in repo %q: %w", chartName, repoURL, err)
	}

	ch, err := loader.Load(chartPath)
	if err != nil {
		return nil, fmt.Errorf("loading chart from %q: %w", chartPath, err)
	}

	return ch, nil
}

func MergeValueFiles(valuesFiles []string) (map[string]interface{}, error) {
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

	// Wait for the API server
	if err := WaitForAPIServer(kubeClient); err != nil {
		return fmt.Errorf("waiting for API server: %w", err)
	}

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
	mergedVals, err := MergeValueFiles(valsFiles)
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

// ApplyKustomize builds the kustomization at `dir` and applies
// all resulting resources server-side (using server-side apply).
func ApplyKustomize(dir string) error {
	// Build with kustomize
	fsys := filesys.MakeFsOnDisk()
	kustomizer := krusty.MakeKustomizer(krusty.MakeDefaultOptions())
	resMap, err := kustomizer.Run(fsys, dir)
	if err != nil {
		return fmt.Errorf("kustomize build %q: %w", dir, err)
	}

	// yamlBytes contains the YAML output of the kustomization, all resources should be in here
	yamlBytes, err := resMap.AsYaml()
	if err != nil {
		return fmt.Errorf("serializing kustomize output: %w", err)
	}

	// Set up dynamic client & RESTMapper
	settings := cli.New()
	restConfig, err := settings.RESTClientGetter().ToRESTConfig()
	if err != nil {
		return fmt.Errorf("getting REST config: %w", err)
	}
	dyn, err := dynamic.NewForConfig(restConfig)
	if err != nil {
		return fmt.Errorf("creating dynamic client: %w", err)
	}
	disc, err := discovery.NewDiscoveryClientForConfig(restConfig)
	if err != nil {
		return fmt.Errorf("creating discovery client: %w", err)
	}
	cacheClient := memorycached.NewMemCacheClient(disc)
	mapper := restmapper.NewDeferredDiscoveryRESTMapper(cacheClient)

	// Decode YAML documents one by one
	dec := yaml.NewYAMLOrJSONDecoder(bytes.NewReader(yamlBytes), 4096)
	for {
		var obj unstructured.Unstructured
		if err := dec.Decode(&obj.Object); err != nil {
			if err == io.EOF {
				break
			}
			return fmt.Errorf("decoding manifest: %w", err)
		}
		// Skip empty docs
		if len(obj.Object) == 0 {
			continue
		}

		gvk := obj.GroupVersionKind()
		mapping, err := mapper.RESTMapping(
			schema.GroupKind{Group: gvk.Group, Kind: gvk.Kind},
			gvk.Version,
		)
		if err != nil {
			return fmt.Errorf("finding REST mapping for %v: %w", gvk, err)
		}

		// Choose namespaced vs cluster-scoped
		var dr dynamic.ResourceInterface
		if mapping.Scope.Name() == meta.RESTScopeNameNamespace {
			ns := obj.GetNamespace()
			if ns == "" {
				ns = corev1.NamespaceDefault
			}
			dr = dyn.Resource(mapping.Resource).Namespace(ns)
		} else {
			dr = dyn.Resource(mapping.Resource)
		}

		// Apply to the cluster via server-side apply
		name := obj.GetName()
		// Use “upctl” as the field manager
		_, err = dr.Apply(context.Background(), name, &obj, metav1.ApplyOptions{
			FieldManager: "upctl",
			Force:        true,
		})
		if err != nil {
			return fmt.Errorf("applying %s/%s: %w", mapping.Resource.Resource, name, err)
		}
	}

	return nil
}

// ExecInPod runs a command in the given pod and container and returns stdout/stderr.
// If containerName is empty, it will auto-detect the container if there is only one in the pod.
// If there are multiple containers, it returns an error asking to specify the container name.
func ExecInPod(
	restConfig *rest.Config,
	kubeClient *kubernetes.Clientset,
	namespace, podName string,
	containerName string, // "" means autodetect if only 1 container
	cmd []string,
	stdin io.Reader,
) (stdout, stderr []byte, err error) {
	// Auto-detect container if not provided
	if containerName == "" {
		pod, err := kubeClient.CoreV1().Pods(namespace).Get(context.Background(), podName, metav1.GetOptions{})
		if err != nil {
			return nil, nil, fmt.Errorf("failed to get pod %s: %w", podName, err)
		}

		switch len(pod.Spec.Containers) {
		case 0:
			return nil, nil, fmt.Errorf("pod %s has no containers", podName)
		case 1:
			containerName = pod.Spec.Containers[0].Name
		default:
			names := make([]string, 0, len(pod.Spec.Containers))
			for _, c := range pod.Spec.Containers {
				names = append(names, c.Name)
			}

			return nil, nil, fmt.Errorf("a container name must be specified for pod %s, choose one of: %v", podName, names)
		}
	}

	req := kubeClient.CoreV1().RESTClient().
		Post().
		Resource("pods").
		Namespace(namespace).
		Name(podName).
		SubResource("exec").
		Param("container", containerName).
		VersionedParams(&corev1.PodExecOptions{
			Command: cmd,
			Stdin:   stdin != nil,
			Stdout:  true,
			Stderr:  true,
			TTY:     false,
		}, scheme.ParameterCodec)

	executor, err := remotecommand.NewSPDYExecutor(restConfig, "POST", req.URL())
	if err != nil {
		return nil, nil, fmt.Errorf("building executor: %w", err)
	}

	var outBuf, errBuf bytes.Buffer
	err = executor.StreamWithContext(context.Background(), remotecommand.StreamOptions{
		Stdin:  stdin,
		Stdout: &outBuf,
		Stderr: &errBuf,
	})

	return outBuf.Bytes(), errBuf.Bytes(), err
}

// CheckDokkuInstallation ensures the Dokku Deployment exists and its pods become Ready.
// It first does a single GET on the Deployment, then polls every 5s for up to 5m
// until at least one pod with label app=dokku reports PodReady.
func CheckDokkuInstallation(kubeClient *kubernetes.Clientset) error {
	if _, err := kubeClient.AppsV1().
		Deployments("dokku").
		Get(context.Background(), "dokku", metav1.GetOptions{}); err != nil {
		return fmt.Errorf("dokku deployment missing: %w", err)
	}

	err := wait.PollUntilContextTimeout(
		context.Background(),
		5*time.Second,
		5*time.Minute,
		true,
		func(ctx context.Context) (done bool, err error) {
			pods, err := kubeClient.CoreV1().
				Pods("dokku").
				List(ctx, metav1.ListOptions{LabelSelector: "app=dokku"})
			if err != nil {
				// transient error, try again
				return false, nil
			}
			if len(pods.Items) == 0 {
				// no pods yet
				return false, nil
			}
			// check if any pod is Ready
			for _, pod := range pods.Items {
				for _, cond := range pod.Status.Conditions {
					if cond.Type == corev1.PodReady && cond.Status == corev1.ConditionTrue {
						return true, nil
					}
				}
			}
			// none ready yet
			return false, nil
		},
	)
	if err != nil {
		return fmt.Errorf("timed out waiting for dokku pods to be ready: %w", err)
	}

	return nil
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

// GetDokkuPodName returns the first pod name in namespace "dokku" with label app=dokku.
func GetDokkuPodName(kubeClient *kubernetes.Clientset) (string, error) {
	pods, err := kubeClient.CoreV1().
		Pods("dokku").
		List(context.Background(), metav1.ListOptions{LabelSelector: "app=dokku"})
	if err != nil {
		return "", fmt.Errorf("listing dokku pods: %w", err)
	}
	if len(pods.Items) == 0 {
		return "", fmt.Errorf("no dokku pods found")
	}
	return pods.Items[0].Name, nil
}

// GetNodeExternalIP returns the first ExternalIP of any node in the cluster.
func GetNodeExternalIP(kubeClient *kubernetes.Clientset) (string, error) {
	nodes, err := kubeClient.CoreV1().
		Nodes().
		List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return "", fmt.Errorf("listing nodes: %w", err)
	}
	if len(nodes.Items) == 0 {
		return "", fmt.Errorf("no nodes found")
	}
	for _, addr := range nodes.Items[0].Status.Addresses {
		if addr.Type == corev1.NodeExternalIP {
			return addr.Address, nil
		}
	}
	return "", fmt.Errorf("no external IP found for node %s", nodes.Items[0].Name)
}

func ConfigureDokku(
	kubeClient *kubernetes.Clientset,
	privateKeyPath, publicKeyPath,
	githubPAT, githubUser, registryURL,
	globalDomain, certManagerEmail string,
) (lbHostname, nodeIP string, err error) {
	// Load REST config so we can ExecInPod
	settings := cli.New()
	restConfig, err := settings.RESTClientGetter().ToRESTConfig()
	if err != nil {
		return "", "", fmt.Errorf("loading kubeconfig: %w", err)
	}

	// Check Dokku is installed & ready
	if err := CheckDokkuInstallation(kubeClient); err != nil {
		return "", "", err
	}

	// Check SSH keys locally
	if err := CheckSSHKeys(privateKeyPath, publicKeyPath); err != nil {
		return "", "", err
	}

	// Find the Dokku pod
	podName, err := GetDokkuPodName(kubeClient)
	if err != nil {
		return "", "", err
	}

	// ssh-keys:add admin <public-key>
	pub, err := os.Open(publicKeyPath)
	if err != nil {
		return "", "", fmt.Errorf("opening public key: %w", err)
	}
	if _, _, err := ExecInPod(restConfig, kubeClient,
		"dokku", podName, "dokku",
		[]string{"dokku", "ssh-keys:add", "admin"},
		pub,
	); err != nil {
		return "", "", fmt.Errorf("dokku ssh-keys:add: %w", err)
	}

	// Registry login
	patReader := bytes.NewBufferString(githubPAT)
	if _, _, err := ExecInPod(restConfig, kubeClient,
		"dokku", podName, "dokku",
		[]string{"dokku", "registry:login", registryURL, githubUser, "--password-stdin"},
		patReader,
	); err != nil {
		return "", "", fmt.Errorf("dokku registry:login: %w", err)
	}

	// Create registry-credential secret in dokku namespace
	cmd := []string{"sh", "-c",
		"kubectl create secret generic registry-credential " +
			"--from-file=.dockerconfigjson=/home/dokku/.docker/config.json " +
			"--type=kubernetes.io/dockerconfigjson --dry-run=client -o yaml | " +
			"kubectl apply -n dokku -f -",
	}
	if _, _, err := ExecInPod(restConfig, kubeClient, "dokku", podName, "dokku", cmd, nil); err != nil {
		return "", "", fmt.Errorf("creating registry-credential: %w", err)
	}

	// Dokku config:set --global CERT_MANAGER_EMAIL=…
	if _, _, err := ExecInPod(restConfig, kubeClient,
		"dokku", podName, "dokku",
		[]string{"dokku", "config:set", "--global", "CERT_MANAGER_EMAIL=" + certManagerEmail},
		nil,
	); err != nil {
		return "", "", fmt.Errorf("setting CERT_MANAGER_EMAIL: %w", err)
	}

	// dokku domains:set-global
	if _, _, err := ExecInPod(restConfig, kubeClient,
		"dokku", podName, "dokku",
		[]string{"dokku", "domains:set-global", globalDomain},
		nil,
	); err != nil {
		return "", "", fmt.Errorf("setting global domain: %w", err)
	}

	// registry:set server & image-repo-template
	if _, _, err := ExecInPod(restConfig, kubeClient,
		"dokku", podName, "dokku",
		[]string{"dokku", "registry:set", "--global", "server", registryURL},
		nil,
	); err != nil {
		return "", "", fmt.Errorf("registry:set server: %w", err)
	}

	imageTemplate := fmt.Sprintf("%s/{{ .AppName }}", githubUser)
	if _, _, err := ExecInPod(restConfig, kubeClient,
		"dokku", podName, "dokku",
		[]string{"dokku", "registry:set", "--global", "image-repo-template", imageTemplate},
		nil,
	); err != nil {
		return "", "", fmt.Errorf("registry:set image-repo-template: %w", err)
	}

	// builder:set herokuish
	if _, _, err := ExecInPod(restConfig, kubeClient,
		"dokku", podName, "dokku",
		[]string{"dokku", "builder:set", "--global", "selected", "herokuish"},
		nil,
	); err != nil {
		return "", "", fmt.Errorf("builder:set herokuish: %w", err)
	}

	// Wait for ingress load-balancer up to 10m
	lb, err := WaitForLoadBalancer(kubeClient, "ingress-nginx", "ingress-nginx-controller", 60, 10*time.Second)
	if err != nil {
		return "", "", fmt.Errorf("waiting ingress loadbalancer: %w", err)
	}

	// Node external IP
	nip, err := GetNodeExternalIP(kubeClient)
	if err != nil {
		return lb, "", fmt.Errorf("getting node external IP: %w", err)
	}

	return lb, nip, nil
}

func PrintFinalInstructions(kubeconfigPath, globalDomain, sshKeyPath, lbHostname, nodeIP string) {
	fmt.Println()
	fmt.Println("✅ Dokku installed and configured successfully!")
	fmt.Println()
	fmt.Println("---------------------------------------------")
	fmt.Println("📌 Before deploying apps you will have to set your local environment:")
	fmt.Printf("export KUBECONFIG=%s\n", kubeconfigPath)
	fmt.Printf("export GLOBAL_DOMAIN=%s\n", globalDomain)
	fmt.Println()
	fmt.Println("🌍 If you have a DNS name add the following DNS record to your domain:")
	fmt.Printf("  CNAME *.%s %s\n", globalDomain, lbHostname)
	fmt.Println()
	fmt.Println("🧪 Otherwise, follow the local testing instructions below.")
	fmt.Println()
	fmt.Println("📁 Add ssh config to your ~/.ssh/config (create it if you don't have it already):")
	fmt.Println("Host dokku")
	fmt.Printf("  Hostname %s\n", nodeIP)
	fmt.Println("  Port 30022")
	fmt.Println("  User dokku")
	fmt.Printf("  IdentityFile %s\n", sshKeyPath)
	fmt.Println()
	fmt.Println("---------------------------------------------")
	fmt.Println("🚀 Deploy your first app:")
	fmt.Println()
	fmt.Println("1. Create Dokku app:")
	fmt.Println("   make create-app APP_NAME=demo-app")
	fmt.Println()
	fmt.Println("2. Clone a sample app (e.g. Heroku Node.js sample):")
	fmt.Println("   mkdir apps && cd apps")
	fmt.Println("   git clone https://github.com/heroku/node-js-sample.git demo-app")
	fmt.Println("   cd demo-app")
	fmt.Println()
	fmt.Println("3. Set the Git remote:")
	fmt.Println("   git remote add dokku dokku@dokku:demo-app")
	fmt.Println()
	fmt.Println("4. Push the app:")
	fmt.Println("   git push dokku master")
	fmt.Println()
	fmt.Println("🌐 Local testing (if you don't have a real DNS)")
	fmt.Println()
	fmt.Println("1. Get the external IP (<EXTERNAL-IP) of the load balancer:")
	fmt.Printf("   dig +short %s\n", lbHostname)
	fmt.Println("2. Edit your local /etc/hosts file:")
	fmt.Println("   sudo vim /etc/hosts")
	fmt.Println()
	fmt.Println("   Add a line like this:")
	fmt.Printf("     <EXTERNAL-IP> demo-app.%s\n", globalDomain)
	fmt.Println()
	fmt.Println("   Example:")
	fmt.Printf("     5.22.219.157 demo-app.%s\n", globalDomain)
	fmt.Println()
	fmt.Println("3. Open your browser and visit:")
	fmt.Printf("   https://demo-app.%s\n", globalDomain)
	fmt.Println()
	fmt.Println("---------------------------------------------")
	fmt.Println("📦 You can repeat this for more apps using:")
	fmt.Println("   make create-app APP_NAME=another-app")
	fmt.Println("   git remote add dokku dokku@dokku:another-app")
	fmt.Println("   git push dokku master")
	fmt.Println()
	fmt.Println("---------------------------------------------")
	fmt.Println("🛠️ We have several make targets to help you manage your Dokku apps:")
	fmt.Println("Run 'make help' to see the available commands.")
	fmt.Println()
	fmt.Println("✨ Enjoy 🚀")
}
