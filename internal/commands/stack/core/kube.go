package core

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud/request"
	"helm.sh/helm/v3/pkg/cli"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
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

// WaitForLoadBalancer waits for a Kubernetes LoadBalancer service to become available and returns its hostname.
func WaitForLoadBalancer(kubeClient *kubernetes.Clientset, namespace, serviceName string, maxRetries int, sleepInterval time.Duration) (string, error) {
	for i := 1; i <= maxRetries; i++ {
		svc, err := kubeClient.CoreV1().Services(namespace).Get(context.Background(), serviceName, v1.GetOptions{})
		if err == nil && len(svc.Status.LoadBalancer.Ingress) > 0 {
			hostname := svc.Status.LoadBalancer.Ingress[0].Hostname
			if hostname != "" {
				return hostname, nil
			}
		}

		time.Sleep(sleepInterval)
	}

	return "", fmt.Errorf("timed out waiting for LoadBalancer hostname after %d attempts", maxRetries)
}

// GetKubernetesClient creates a Kubernetes client from the provided kubeconfig path.
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

// WaitForAPIServer waits for the Kubernetes API server to be ready.
func WaitForAPIServer(kubeClient *kubernetes.Clientset) error {
	for i := 1; i <= 30; i++ {
		if _, err := kubeClient.Discovery().ServerVersion(); err == nil {
			return nil
		} else {
			time.Sleep(5 * time.Second)
		}
	}
	return fmt.Errorf("timed out waiting for API server")
}

// CreateNamespace creates a Kubernetes namespace if it does not already exist.
func CreateNamespace(kubeClient *kubernetes.Clientset, namespace string) error {
	_, err := kubeClient.CoreV1().Namespaces().Get(context.Background(), namespace, v1.GetOptions{})
	if err != nil {
		// Check for StatusError and StatusReasonNotFound
		if apierrors.IsNotFound(err) {
			_, createErr := kubeClient.CoreV1().Namespaces().Create(context.Background(), &corev1.Namespace{
				ObjectMeta: v1.ObjectMeta{
					Name: namespace,
				},
			}, v1.CreateOptions{})
			if createErr != nil {
				return fmt.Errorf("failed to create namespace %q: %w", namespace, createErr)
			}
		} else {
			return fmt.Errorf("error getting namespace %q: %w", namespace, err)
		}
	}
	return nil
}

// ApplyKustomize builds the kustomization at `dir` and applies
// all resulting resources server-side (using server-side apply).
func ApplyKustomize(dir string) error {
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
		_, err = dr.Apply(context.Background(), name, &obj, v1.ApplyOptions{
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
		pod, err := kubeClient.CoreV1().Pods(namespace).Get(context.Background(), podName, v1.GetOptions{})
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

// GetNodeExternalIP returns the first ExternalIP of any node in the cluster.
func GetNodeExternalIP(kubeClient *kubernetes.Clientset) (string, error) {
	nodes, err := kubeClient.CoreV1().
		Nodes().
		List(context.Background(), v1.ListOptions{})
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

// writeKubeconfigToFile retrieves the kubeconfig for the given cluster and writes it to a file
func WriteKubeconfigToFile(exec commands.Executor, clusterId string, configDir string) (string, error) {
	kubeconfig, err := exec.All().GetKubernetesKubeconfig(exec.Context(), &request.GetKubernetesKubeconfigRequest{
		UUID: clusterId,
	})

	if err != nil {
		return "", fmt.Errorf("failed to get kubeconfig for cluster %s: %w", clusterId, err)
	}

	kubeconfigPath := filepath.Join(configDir, "kubeconfig.yaml")
	if err := os.WriteFile(kubeconfigPath, []byte(kubeconfig), 0600); err != nil {
		return "", fmt.Errorf("failed to write kubeconfig: %w", err)
	}

	return kubeconfigPath, nil
}
