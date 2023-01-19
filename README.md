# GCP Storage Proxy function

This Google Cloud Function acts as a proxy to list and download Google Storage Bucket objects.

## Usage

List objects in a bucket

```bash
curl  -H "Authorization: Bearer $(gcloud auth print-identity-token)" \
https://functionid.a.run.app/info?bucket=<bucket_name>
```

Download a file from a bucket

```bash
curl  -H "Authorization: Bearer $(gcloud auth print-identity-token)" \
https://functionid.a.run.app/download?bucket=<bucket_name>&filename=<file_name>
```

To call the function using custom authentication, fill the `apikey` query parameter with your API key stored in Firebase.

```bash
curl https://functionid.a.run.app/download?bucket=<bucket_name>&filename=<file_name>&apikey=<api_key>
```

## Deploy

Google Cloud CLI (gcloud) is required for deployment with the following command. 
Otherwise the function can be deployed from Google Cloud Console as well (requires manual source code copying).

```bash
gcloud functions deploy gcp-storage-proxy --trigger-http --gen2 --runtime go119 \
--region=europe-west1 --source . --entry-point Handler
```

Custom authentication via Firebase/Firestore is also supported. To use custom auth instead of native GCP authentication, deploy the function as follows

```bash
gcloud functions deploy gcp-storage-proxy --set-env-vars AUTH_ENABLED=true,PROJECT_ID=<your_project_id> --allow-unauthenticated --trigger-http --gen2 --runtime go119 --region=europe-west1 --source . --entry-point Handler
```

## Local testing

**WARNING: This doesn't seem to work on Windows systems**

For interacting with storage buckets locally, you need to be authenticated to GCP - either via gcloud CLI or a credentials file (found in `$HOME/.config/gcloud/application_default_credentials.json`).

```bash
git clone https://github.com/Loupeznik/gcp-storage-proxy
cd gcp-storage-proxy

export FUNCTION_TARGET=Handler
go run cmd/main.go

curl http://localhost:8080/info?bucket=<bucket_name>
```
