package dokku

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands/stack"
	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud/request"
	"helm.sh/helm/v3/pkg/cli"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
)

func (s *deployDokkuCommand) deploy(exec commands.Executor, configDir string) error {
	clusterName := fmt.Sprintf("%s-%s-%s", stack.DokkuResourceRootNameCluster, s.name, s.zone)
	networkName := fmt.Sprintf("%s-%s-%s", stack.DokkuResourceRootNameNetwork, s.name, s.zone)
	var network *upcloud.Network

	// Check if the cluster already exists
	clusters, err := exec.All().GetKubernetesClusters(exec.Context(), &request.GetKubernetesClustersRequest{})
	if err != nil {
		return fmt.Errorf("failed to get Kubernetes clusters: %w", err)
	}

	// Return if the cluster already exists
	if stack.ClusterExists(clusterName, clusters) {
		return fmt.Errorf("a cluster with the name '%s' already exists", clusterName)
	}

	msg := fmt.Sprintf("Creating Kubernetes cluster %s in zone %s", clusterName, s.zone)
	exec.PushProgressStarted(msg)

	// Check if the network already exists
	networks, err := exec.Network().GetNetworks(exec.Context())
	if err != nil {
		return fmt.Errorf("failed to get networks: %w", err)
	}

	network = stack.GetNetworkFromName(networkName, networks.Networks)

	// Create the network if it does not exist
	if network == nil {
		network, err = stack.CreateNetwork(exec, networkName, s.zone, stack.StackTypeDokku)
		if err != nil {
			return fmt.Errorf("failed to create network: %w", err)
		}

		if network == nil {
			return fmt.Errorf("created network %s is nil", networkName)
		}
		if len(network.IPNetworks) == 0 {
			return fmt.Errorf("created network %s has no IP networks", networkName)
		}
	}

	cluster, err := exec.All().CreateKubernetesCluster(exec.Context(), &request.CreateKubernetesClusterRequest{
		Name:        clusterName,
		Network:     network.UUID,
		Zone:        s.zone,
		NetworkCIDR: network.IPNetworks[0].Address,
		Labels: []upcloud.Label{
			{Key: "stacks.upcloud.com/stack", Value: "dokku"},
			{Key: "stacks.upcloud.com/created-by", Value: "upctl"},
			{Key: "stacks.upcloud.com/dokku-version", Value: "0.1.3"}, // Dokku chart version
			{Key: "stacks.upcloud.com/version", Value: string(stack.VersionV0_1_0_0)},
			{Key: "stacks.upcloud.com/name", Value: clusterName},
		},
		NodeGroups: []request.KubernetesNodeGroup{
			{
				Name:  "dokku-node-group",
				Count: s.numNodes,
				Plan:  "2xCPU-4GB",
			},
		},
	})
	if err != nil {
		return fmt.Errorf("failed to create Kubernetes cluster: %w", err)
	}

	exec.PushProgressSuccess(msg)
	exec.PushProgressStarted("Setting up environment for Dokku stack deployment")

	// Get kubeconfig file for the cluster
	kubeconfigPath, err := stack.WriteKubeconfigToFile(exec, cluster.UUID, configDir)
	if err != nil {
		return fmt.Errorf("failed to write kubeconfig to file: %w", err)
	}

	// Create a Kubernetes client from the kubeconfig
	os.Setenv("KUBECONFIG", kubeconfigPath)
	kubeClient, err := stack.GetKubernetesClient(kubeconfigPath)
	if err != nil {
		return fmt.Errorf("failed to create Kubernetes client: %w", err)
	}

	// Wait for the Kubernetes API server to be ready
	_, err = exec.All().WaitForKubernetesClusterState(exec.Context(), &request.WaitForKubernetesClusterStateRequest{
		UUID:         cluster.UUID,
		DesiredState: upcloud.KubernetesClusterStateRunning,
	})
	if err != nil {
		return err
	}

	exec.PushProgressSuccess("Setting up environment for Dokku stack deployment")
	exec.PushProgressStarted("Deploying Dokku stack")
	// Deploy nginx ingress controller
	ingressValuesFilePath := filepath.Join(configDir, "config/ingress/values.yaml")
	err = stack.DeployHelmReleaseFromRepo(
		kubeClient,
		"ingress-nginx",
		"https://kubernetes.github.io/ingress-nginx",
		"ingress-nginx",
		"4.13.0",
		[]string{ingressValuesFilePath},
		false)
	if err != nil {
		return fmt.Errorf("failed to deploy ingress-nginx: %w", err)
	}

	// Wait for the ingress controller to be ready
	lbHostname, err := stack.WaitForLoadBalancer(kubeClient, "ingress-nginx", "ingress-nginx-controller", 30, 10*time.Second)
	if err != nil {
		return fmt.Errorf("failed to wait for ingress-nginx load balancer: %w", err)
	}

	// If no global domain is provided, use the load balancer hostname
	if s.globalDomain == "" {
		s.globalDomain = lbHostname
	}

	// Deploy cert-manager
	override, err := os.CreateTemp("", "certmgr-override-*.yaml")
	if err != nil {
		return err
	}
	defer func(name string) {
		errRemove := os.Remove(name)
		if errRemove != nil {
			fmt.Printf("failed to remove temp file %s: %v\n", name, errRemove)
		}
	}(override.Name())

	_, err = override.WriteString("installCRDs: true\n")
	if err != nil {
		return err
	}
	err = override.Close()
	if err != nil {
		return err
	}

	err = stack.DeployHelmReleaseFromRepo(
		kubeClient,
		"cert-manager",
		"https://charts.jetstack.io",
		"cert-manager",
		"v1.11.0",
		[]string{override.Name()},
		false,
	)
	if err != nil {
		return fmt.Errorf("failed to deploy cert-manager: %w", err)
	}

	// Deploy Dokku
	dokkuConfigPath := filepath.Join(configDir, "config/dokku")
	err = stack.ApplyKustomize(dokkuConfigPath)
	if err != nil {
		return fmt.Errorf("failed to apply Dokku kustomize configuration: %w", err)
	}

	// Configure Dokku with the provided parameters
	lbHostName, nodeIP, err := configureDokku(
		kubeClient,
		s.sshPath,
		s.sshPubPath,
		s.githubPAT,
		s.githubUser,
		s.githubPackageURL,
		s.globalDomain,
		s.certManagerEmail)
	if err != nil {
		return fmt.Errorf("failed to configure Dokku: %w", err)
	}

	exec.PushProgressSuccess("Deploying Dokku stack")

	// Print final instructions for the user
	printFinalInstructions(kubeconfigPath, s.globalDomain, s.sshPath, lbHostName, nodeIP)

	return nil
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

func configureDokku(
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
	if err := stack.CheckSSHKeys(privateKeyPath, publicKeyPath); err != nil {
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
	if _, _, err := stack.ExecInPod(restConfig, kubeClient,
		"dokku", podName, "dokku",
		[]string{"dokku", "ssh-keys:add", "admin"},
		pub,
	); err != nil {
		return "", "", fmt.Errorf("dokku ssh-keys:add: %w", err)
	}

	// Registry login
	patReader := bytes.NewBufferString(githubPAT)
	if _, _, err := stack.ExecInPod(restConfig, kubeClient,
		"dokku", podName, "dokku",
		[]string{"dokku", "registry:login", registryURL, githubUser, "--password-stdin"},
		patReader,
	); err != nil {
		return "", "", fmt.Errorf("dokku registry:login: %w", err)
	}

	// Create registry-credential secret in dokku namespace
	cmd := []string{
		"sh", "-c",
		"kubectl create secret generic registry-credential " +
			"--from-file=.dockerconfigjson=/home/dokku/.docker/config.json " +
			"--type=kubernetes.io/dockerconfigjson --dry-run=client -o yaml | " +
			"kubectl apply -n dokku -f -",
	}
	if _, _, err := stack.ExecInPod(restConfig, kubeClient, "dokku", podName, "dokku", cmd, nil); err != nil {
		return "", "", fmt.Errorf("creating registry-credential: %w", err)
	}

	// Dokku config:set --global CERT_MANAGER_EMAIL=…
	if _, _, err := stack.ExecInPod(restConfig, kubeClient,
		"dokku", podName, "dokku",
		[]string{"dokku", "config:set", "--global", "CERT_MANAGER_EMAIL=" + certManagerEmail},
		nil,
	); err != nil {
		return "", "", fmt.Errorf("setting CERT_MANAGER_EMAIL: %w", err)
	}

	// dokku domains:set-global
	if _, _, err := stack.ExecInPod(restConfig, kubeClient,
		"dokku", podName, "dokku",
		[]string{"dokku", "domains:set-global", globalDomain},
		nil,
	); err != nil {
		return "", "", fmt.Errorf("setting global domain: %w", err)
	}

	// registry:set server & image-repo-template
	if _, _, err := stack.ExecInPod(restConfig, kubeClient,
		"dokku", podName, "dokku",
		[]string{"dokku", "registry:set", "--global", "server", registryURL},
		nil,
	); err != nil {
		return "", "", fmt.Errorf("registry:set server: %w", err)
	}

	imageTemplate := fmt.Sprintf("%s/{{ .AppName }}", githubUser)
	if _, _, err := stack.ExecInPod(restConfig, kubeClient,
		"dokku", podName, "dokku",
		[]string{"dokku", "registry:set", "--global", "image-repo-template", imageTemplate},
		nil,
	); err != nil {
		return "", "", fmt.Errorf("registry:set image-repo-template: %w", err)
	}

	// builder:set herokuish
	if _, _, err := stack.ExecInPod(restConfig, kubeClient,
		"dokku", podName, "dokku",
		[]string{"dokku", "builder:set", "--global", "selected", "herokuish"},
		nil,
	); err != nil {
		return "", "", fmt.Errorf("builder:set herokuish: %w", err)
	}

	// Wait for ingress load-balancer up to 10m
	lb, err := stack.WaitForLoadBalancer(kubeClient, "ingress-nginx", "ingress-nginx-controller", 60, 10*time.Second)
	if err != nil {
		return "", "", fmt.Errorf("waiting ingress loadbalancer: %w", err)
	}

	// Node external IP
	nip, err := stack.GetNodeExternalIP(kubeClient)
	if err != nil {
		return lb, "", fmt.Errorf("getting node external IP: %w", err)
	}

	return lb, nip, nil
}

func printFinalInstructions(kubeconfigPath, globalDomain, sshKeyPath, lbHostname, nodeIP string) {
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
