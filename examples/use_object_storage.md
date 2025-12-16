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

Save the access key ID and secret access key from the output - you'll need these for S3 access.

Attach a policy to grant the user access to buckets:

1. **Using the UpCloud API:**

   First, get your service UUID:
   ```sh
   upctl object-storage list
   ```

   Then attach the policy (replace `<service-uuid>` with your UUID):
   ```sh
   curl -X POST "https://api.upcloud.com/1.3/object-storage-2/<service-uuid>/users/${prefix}user/policies" \
     -H "Authorization: Bearer ${UPCLOUD_TOKEN}" \
     -H "Content-Type: application/json" \
     -d '{"name": "ECSS3FullAccess"}'
   ```

   A successful response returns HTTP status 204.

2. **Using the UpCloud Control Panel:**
   Navigate to Object Storage → Users → Select user → Attach Policy → ECSS3FullAccess

**Note:** Without attaching a policy, the user won't have permission to access buckets via AWS CLI or S3-compatible tools.

Verify S3 access with AWS CLI:

```sh
# Get the service endpoint
upctl object-storage show ${prefix}service

# Configure AWS CLI with your credentials and test access
AWS_ACCESS_KEY_ID=<your-access-key-id> \
AWS_SECRET_ACCESS_KEY=<your-secret-access-key> \
aws s3 ls --endpoint-url https://<service-endpoint>
```

Once not needed anymore, delete the user:

```sh
upctl object-storage user delete ${prefix}service --username ${prefix}user
```

Delete also the managed object storage service along with its buckets:

```sh
upctl object-storage delete ${prefix}service --delete-buckets
```
