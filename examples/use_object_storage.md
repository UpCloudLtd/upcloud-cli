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

# Create access key and save credentials
access_key_output=$(upctl object-storage access-key create ${prefix}service --username ${prefix}user -o json)
access_key_id=$(echo "$access_key_output" | jq -r '.access_key_id')
secret_access_key=$(echo "$access_key_output" | jq -r '.secret_access_key')
```

Save the access key ID and secret access key from the output - you'll need these for S3 access.

Attach a policy to grant the user access to buckets:

```sh
upctl object-storage user policy attach ${prefix}service --username ${prefix}user --policy ECSS3FullAccess
```

Verify S3 access with AWS CLI:

```sh when="false"
# Get the service endpoint
service_endpoint=$(upctl object-storage show ${prefix}service -o json | jq -r '.endpoints[0].domain_name')

# Configure AWS CLI with your credentials and test access
AWS_ACCESS_KEY_ID=${access_key_id} \
AWS_SECRET_ACCESS_KEY=${secret_access_key} \
aws s3 ls --endpoint-url https://${service_endpoint}
```

Delete the managed object storage service along with all its sub-resources such as buckets and users:

```sh
upctl object-storage delete ${prefix}service --force
```