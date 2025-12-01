# Create a Kubernetes cluster and expose a service

This example demonstrates how to create a Kubernetes cluster with `upctl`, create a deployment, and expose it to the internet using a service.

To keep track of resources created during this example, we will use common prefix in all resource names. The variable definitions also include configuration for the cluster, node-group, and service. By default, we will use `NodePort` service type (as it is faster to deploy), but you can change `svc_type` value to `LoadBalancer` if you want to use UpCloud Load Balancer to expose the service.

```env
prefix=example-upctl-kubernetes-
zone=pl-waw1

# Network
cidr=172.30.100.0/24

# Cluster
plan=dev-md

# Node-group
ng_count=1
ng_name=default
ng_plan=2xCPU-4GB
test_label=upctl-example

# Service
svc_type=NodePort

KUBECONFIG=./kubeconfig.yaml
```

First, we will need a private network for the cluster.

```sh
upctl network create \
    --name ${prefix}net \
    --zone $zone \
    --ip-network address=$cidr,dhcp=true;
```

Next, we can create the Kubernetes cluster.

```sh
upctl kubernetes create \
    --name ${prefix}cluster \
    --network ${prefix}net \
    --plan $plan \
    --zone $zone \
    --kubernetes-api-allow-ip "0.0.0.0/0" \
    --node-group count=$ng_count,name=$ng_name,plan=$ng_plan,label="test=$test_label" \
    --wait=all;
```

Once the cluster is created, we can get the kubeconfig file to interact with the cluster using `kubectl`.

```sh
upctl kubernetes config ${prefix}cluster \
    --write $KUBECONFIG;
```

Now we can create a deployment and expose it using a service.

```sh
kubectl create deployment --image=ghcr.io/upcloudltd/hello hello-uks
kubectl expose deployment hello-uks --port=80 --target-port=80 --type=$svc_type
```

If using `NodePort` service type, we need to get the node IP and service port to access the service.

```sh
# Skip this code block if not using NodePort
test "$svc_type" = "NodePort" || exit 0;

# Get node IP
node_ip=$(kubectl get node -o json | jq -r '.items[0].status.addresses.[] | select(.type == "ExternalIP
    ").address')
# Get service port
svc_port=$(kubectl get service hello-uks -o json | jq -r '.spec.ports[0].nodePort')

# Wait until the service is reachable
until curl -sSf $node_ip:$svc_port; do
    sleep 15;
done;
```

If using `LoadBalancer` service type, we need to wait until the load balancer has been created and assigned a public hostname. Once the hostname is available, we can try to access the service.

```sh
# Skip this code block if not using LoadBalancer
test "$svc_type" = "LoadBalancer" || exit 0;

# Wait for hostname to be available in service status
until kubectl get service hello-uks -o json | jq -re .status.loadBalancer.ingress[0].hostname; do
    sleep 15;
done;

# Wait until the service is reachable
hostname=$(kubectl get service hello-uks -o json | jq -re .status.loadBalancer.ingress[0].hostname)
until curl -sSf $hostname; do
    sleep 15;
done;
```

Finally, we can clean up the created resources.

```sh cleanup
kubectl delete service hello-uks

upctl all purge --include "*${prefix}*";
```
