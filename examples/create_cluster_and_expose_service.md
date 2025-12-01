# Create a Kubernetes cluster and expose a service

This example demonstrates how to create a Kubernetes cluster with `upctl`, create a deployment, and expose it to the internet using a service of type `LoadBalancer`.

To keep track of resources created during this example, we will use common prefix in all resource names.

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

Now we can create a deployment and expose it using a `LoadBalancer` type service.

```sh
kubectl create deployment --image=ghcr.io/upcloudltd/hello hello-uks
kubectl expose deployment hello-uks --port=80 --target-port=80 --type=LoadBalancer
```

After creating the service, we need to wait until the load balancer has been created and assigned a public hostname.

```sh
until kubectl get service hello-uks -o json | jq -re .status.loadBalancer.ingress[0].hostname; do
    sleep 15;
done;
```

Once the hostname is available, we can try to access the service.

```sh
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
