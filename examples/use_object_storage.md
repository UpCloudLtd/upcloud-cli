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

1. **Using the UpCloud API:**

   First, get your service UUID:
   ```sh
   service_uuid=$(upctl object-storage list -o json | jq -r ".[] | select(.name == \"${prefix}service\") | .uuid")
   ```

   Then attach the policy (requires UPCLOUD_TOKEN environment variable):
   ```sh when="false"
   # Note: This command requires a bearer token which can be created via the UpCloud Control Panel
   # The when=false flag skips this in automated tests since only username/password are available in CI
   curl -X POST "https://api.upcloud.com/1.3/object-storage-2/${service_uuid}/users/${prefix}user/policies" \
     -H "Authorization: Bearer ${UPCLOUD_TOKEN}" \
     -H "Content-Type: application/json" \
     -d '{"name": "ECSS3FullAccess"}'
   ```

   A successful response returns HTTP status 204.

2. **Using the UpCloud Control Panel:**
   Navigate to Object Storage → Users → Select user → Attach Policy → ECSS3FullAccess

**Note:** Without attaching a policy, the user won't have permission to access buckets via AWS CLI or S3-compatible tools.

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