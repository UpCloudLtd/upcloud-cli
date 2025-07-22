# Backup a server and re-create it using the backup

This example demonstrates how to backup a server with `upctl` and use the created backup to re-create the server.

To keep track of resources created during this example, we will use common prefix in all resource names.

```env
prefix=example-upctl-backup-
zone=pl-waw1
```

We will first create a ssh-key into the current working directory for configuring an nginx server via SSH connection.

```sh
ssh-keygen -t ed25519 -q -f "./id_ed25519" -N "" -C "upctl example"
```

We will then create a server with a single network interface and default template settings.

```sh
upctl server create \
    --hostname ${prefix}source-server \
    --zone ${zone} \
    --ssh-keys ./id_ed25519.pub \
    --network type=public \
    --wait
```

To have something to backup, we will install a nginx server and configure a non-default HTML content to serve.

```sh filename=configure-nginx.sh
#!/bin/sh -xe

apt-get update
apt-get install nginx -y
echo "Hello from $(hostname)"'!' > /var/www/html/index.html
```

To configure the server, we will parse the public IP of the server and run the above script using SSH connection. We can then use `curl` to ensure that the HTTP server serves the content we defined.

```sh
# Parse public IP of the server with jq
ip=$(upctl server show ${prefix}source-server -o json | jq -r '.networking.interfaces[] | select(.type == "public") | .ip_addresses[0].address')

# Wait for a moment for the ssh server to become available
sleep 30

# Run the script defined above
ssh -i id_ed25519 -o StrictHostKeyChecking=accept-new root@$ip "sh" < configure-nginx.sh

# Validate HTTP server response
test "$(curl -s $ip)" = 'Hello from example-upctl-backup-source-server!'
```

We will then backup the OS disk of the created server.

```sh
upctl storage backup create ${prefix}source-server-OS --title ${prefix}source-server-backup
```

After creating the backup, we can delete the source server and its storages.

```sh
upctl server stop --type hard --wait ${prefix}source-server
upctl server delete ${prefix}source-server --delete-storages
```

We can then create a new server based on the backup of the source servers disk.

```sh
upctl server create \
    --hostname ${prefix}restored-server \
    --zone ${zone} \
    --ssh-keys ./id_ed25519.pub \
    --network type=public \
    --storage action=clone,storage=${prefix}source-server-backup \
    --wait
```

To validate that the server was re-created successfully, we will parse the public IP of the server and use curl to see that the HTTP server is running.

```sh
# Parse public IP of the server with jq
ip=$(upctl server show ${prefix}restored-server -o json | jq -r '.networking.interfaces[] | select(.type == "public") | .ip_addresses[0].address')

# Wait until server returns expected response
for i in $(seq 1 9); do
    test "$(curl -s $ip)" = 'Hello from example-upctl-backup-source-server!' && break || true;
    sleep 15;
done;
```

Finally, we can cleanup the created resources.

```sh
# Delete the restored server and its storages
upctl server stop --type hard --wait ${prefix}restored-server
upctl server delete ${prefix}restored-server --delete-storages

# Delete the backup
upctl storage delete ${prefix}source-server-backup
```
