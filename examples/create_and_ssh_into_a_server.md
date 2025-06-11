# Create and ssh into a server

This example demonstrates how to create a server with `upctl` and connect to the created server via ssh connection.

To keep track of resources created during this example, we will use common prefix in all resource names.

```env
prefix=example-upctl-ssh-server-
zone=pl-waw1
```

In order to be able to connect to the server we are going to create, we will need an ssh-key. If you already have a ssh-key available, you can skip this step. The example creates the ssh-key into the current working directory, if you want to use this key for other authentication purposes as well, create the key into `~/.ssh` directory instead.

```sh
# Create ssh-key into current working directory
ssh-keygen -t ed25519 -q -f "./id_ed25519" -N "" -C "upctl example"
```

Create a server using the above created ssh-key as login method.

```sh
upctl server create \
    --hostname ${prefix}server \
    --zone ${zone} \
    --ssh-keys ./id_ed25519.pub \
    --network type=public \
    --network type=utility \
    --wait
```

Find the IP address of the created server from the JSON output of `upctl server show` and execute hostname command via ssh connection on the created server.

```sh
# Parse public IP of the server with jq
ip=$(upctl server show ${prefix}server -o json | jq -r '.networking.interfaces[] | select(.type == "public") | .ip_addresses[0].address')

# Wait for a moment for the ssh server to become available
sleep 30

ssh -i id_ed25519 -o StrictHostKeyChecking=accept-new root@$ip "hostname"
```

Finally, we can cleanup the created resources.

```sh
upctl server stop --type hard --wait ${prefix}server
upctl server delete ${prefix}server --delete-storages
```
