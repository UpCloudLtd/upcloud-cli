# Use user-data to configure a server

This example demonstrates how to create a server with `upctl` and run a container inside the server using an user-data script.

To keep track of resources created during this example, we will use common prefix in all resource names.

```env
prefix=example-upctl-server-userdata-
zone=pl-waw1
```

First, create a file for the user-data script. The script installs [podman](https://podman.io/) and uses it to run a [hello container](https://github.com/UpCloudLtd/hello-container) that will be exposed in the port 80 of the server.

```sh filename=user-data.sh
#!/bin/sh
sudo apt-get update
sudo apt-get install -y podman
podman run -d -p 80:80 ghcr.io/upcloudltd/hello
```

Create a server and use the user-data script from `user-data.sh` to configure the server.

```sh
# Create ssh-key into current working directory
ssh-keygen -t ed25519 -q -f "./id_ed25519" -N ""

upctl server create \
    --hostname ${prefix}server \
    --zone ${zone} \
    --plan "DEV-1xCPU-1GB-10GB" \
    --ssh-keys ./id_ed25519.pub \
    --network type=public \
    --user-data "$(cat user-data.sh)" \
    --wait;
```

After the server has started, we can use curl to wait for the container to be available.

```sh
ip=$(upctl server show ${prefix}server -o json | jq -r '.networking.interfaces[] | select(.type == "public") | .ip_addresses[0].address')
until curl -sf http://$ip; do
  sleep 5;
done;
```

Finally, we can cleanup the created resources.

```sh
upctl all purge -i ${prefix}*
```
