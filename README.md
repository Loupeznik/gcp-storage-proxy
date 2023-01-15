# GCP Storage Proxy function

This Google Cloud Function acts as a proxy to list and download Google Storage Bucket objects.

## Usage

List objects in a bucket

```bash
curl  -H "Authorization: Bearer $(gcloud auth print-identity-token)" \
https://functionid.a.run.app/info?bucket=<bucket_name>
```

List objects in a bucket

```bash
curl  -H "Authorization: Bearer $(gcloud auth print-identity-token)" \
https://functionid.a.run.app/download?bucket=<bucket_name>&filename=<file_name>
```

## Deploy

Google Cloud CLI (gcloud) is required for deployment with the following command. 
Otherwise the function can be deployed from Google Cloud Console as well (requires manual source code copying).

```bash
gcloud functions deploy gcp-storage-proxy --trigger-http --gen2 --runtime go119 \
--region=europe-west1 --source . --entry-point Handler
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
