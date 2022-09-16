# Possible exit codes

As mentioned in the [../README.md](../README.md#exit-codes), upctl sets exit code based on number of failed tasks up to exit code 99. This example demonstrates executions with few different exit codes.

To keep track of resources created during this example, we will use common prefix in all resource names.

```env
prefix=example-upctl-exit-codes-
```

Exit code 100 is set, for example, when command argument validation fails.

```sh exit_code=100
upctl server create
# Error: required flag(s) "hostname", "zone" not set
```

Let's create two server and stop one of those to later see other failing exit codes. These command should succceed, and thus return zero exit code.

```sh
upctl server create --hostname ${prefix}vm-1 --zone pl-waw1 --ssh-keys ~/.ssh/*.pub --wait
upctl server create --hostname ${prefix}vm-2 --zone pl-waw1 --ssh-keys ~/.ssh/*.pub --wait

upctl server stop ${prefix}vm-1 --wait
```

Now let's try to stop both both of the created servers. Exit code will be one, as `${prefix}vm-1` is already stopped and thus cannot be stopped again. `${prefix}vm-2`, though, will be stopped as it was online. Thus one of the two operations failed.

```sh exit_code=1
upctl server stop ${prefix}vm-1 ${prefix}vm-2 --wait
```

If we now try to run above command again, exit code will be two as both of the servers are already stopped. Thus both stop operations failed.

```sh exit_code=2
upctl server stop ${prefix}vm-1 ${prefix}vm-2 --wait
```

Finally, we can cleanup the created resources.

```sh
upctl server delete ${prefix}vm-1 ${prefix}vm-2 --delete-storages
```
