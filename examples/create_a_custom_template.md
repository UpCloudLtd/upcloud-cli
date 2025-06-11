# Create a custom template

This example demonstrates how to create a custom template with `upctl`.

To keep track of resources created during this example, we will use common prefix in all resource names.

```env
prefix=example-upctl-custom-template-
zone=pl-waw1
```

First, we will create server which disk will be used as a source for the custom template.

```sh
# Create ssh-key into current working directory
ssh-keygen -t ed25519 -q -f "./id_ed25519" -N ""

upctl server create \
    --hostname ${prefix}source-server \
    --zone ${zone} \
    --ssh-keys ./id_ed25519.pub \
    --network type=public \
    --network type=utility \
    --wait
```

After the server has started, you can connect to it and prepare the disk to be templatized. Then, to be able to templatize the storage disk, we will stop the server.

```sh
upctl server stop --type hard --wait ${prefix}source-server
```

The default name for the OS storage of servers created with `upctl` is `${server-title}-OS`, in this case `${prefix}source-server-OS`. We can use either that or the UUID of the storage, when creating the template. UUID of the storage can be printed, for example, by processing `json` output with jq.

```sh
upctl server show ${prefix}source-server -o json \
    | jq -r ".storage[0].uuid"
```

Now we are ready for creating the template.

```sh
upctl storage templatise ${prefix}source-server-OS \
    --title ${prefix}template \
    --wait
```

Once the template is created, we can delete the source server 

```sh
upctl server delete ${prefix}source-server --delete-storages
```

To test that the template creation succeeded, create a new server from the just created template.

```sh
upctl server create \
    --hostname ${prefix}server \
    --zone ${zone} \
    --network type=public \
    --network type=utility \
    --os ${prefix}template \
    --wait
```

Finally, we can cleanup the created resources.

```sh
upctl server stop --type hard --wait ${prefix}server
upctl server delete ${prefix}server --delete-storages
upctl storage delete ${prefix}template
```
