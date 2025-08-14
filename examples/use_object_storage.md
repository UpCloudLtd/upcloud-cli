# Using Object Storage with upctl

This example demonstrates how to create and manage an object storage service, including managing buckets and configuring
S3-compatible access via access keys.

Set environment variables for convenience:

```env
prefix=example-upctl-
region=europe-1
```

Create a managed object storage service:

```sh
upctl object-storage create --name ${prefix}service --network type=public,name=${prefix}network,family=IPv4 --region ${region}
```

List all buckets in your service (will be empty at this point):

```sh
upctl object-storage bucket list ${prefix}service
```

Create a new bucket:

```sh
upctl object-storage bucket create ${prefix}service --name ${prefix}bucket
```

Create a user and an access key for S3-compatible access:

```sh
upctl object-storage user create ${prefix}service --username ${prefix}user
upctl object-storage access-key create ${prefix}service --username ${prefix}user
```

Once not needed anymore, delete the user:

```sh
upctl object-storage user delete ${prefix}service --username ${prefix}user
```

Delete also the managed object storage service along with its buckets:

```sh
upctl object-storage delete ${prefix}service --delete-buckets
```
