# Possible exit codes

`upctl` sets exit code based on number of failed tasks up to exit code 99. This example demonstrates executions with few different exit codes.

To keep track of resources created during this example, we will use common prefix in all resource names.

```env
prefix=example-upctl-exit-codes-
zone=pl-waw1
```

Exit code 100 is set, for example, when command argument validation fails.

```sh exit_code=100
upctl server create
# Error: required flag(s) "hostname", "zone" not set
```

Let's create two servers and stop one of those to later see other failing exit codes. This example uses `--type hard` when stopping the servers as the OS might not be completely up and running when the server reaches running state. These command should succeed, and thus return zero exit code.

```sh
# Create ssh-key into current working directory
ssh-keygen -t ed25519 -q -f "./id_ed25519" -N ""

upctl server create --hostname ${prefix}vm-1 --zone ${zone} --ssh-keys ./id_ed25519.pub --wait
upctl server create --hostname ${prefix}vm-2 --zone ${zone} --ssh-keys ./id_ed25519.pub --wait

upctl server stop --type hard ${prefix}vm-1 --wait
```

Now let's try to stop both both of the created servers. Exit code will be one, as `${prefix}vm-1` is already stopped and thus cannot be stopped again. `${prefix}vm-2`, though, will be stopped as it was online. Thus one of the two operations failed.

```sh exit_code=1
upctl server stop --type hard ${prefix}vm-1 ${prefix}vm-2 --wait
```

If we now try to run above command again, exit code will be two as both of the servers are already stopped. Thus both stop operations failed.

```sh exit_code=2
upctl server stop --type hard ${prefix}vm-1 ${prefix}vm-2 --wait
```

Let's cleanup the created resources while we have working credentials.

```sh
upctl server delete ${prefix}vm-1 ${prefix}vm-2 --delete-storages
```

To test validation of credentials, we configure `UPCLOUD_TOKEN` environment variable with invalid value.

```env
UPCLOUD_TOKEN=invalid
```

When authentication fails, upctl sets exit code to 103.

```sh exit_code=103
upctl server list
# Error: invalid user credentials, authentication failed using the given username and password
```
