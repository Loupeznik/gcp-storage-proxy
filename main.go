package function

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"cloud.google.com/go/firestore"
	"cloud.google.com/go/storage"
	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"google.golang.org/api/iterator"
)

func init() {
	functions.HTTP("GetBucket", getBucket)
	functions.HTTP("DownloadFile", downloadFile)
	functions.HTTP("Handler", handler)
}

func isAuthEnabled() bool {
	isEnabled, err := strconv.ParseBool(os.Getenv("AUTH_ENABLED"))

	if err != nil {
		isEnabled = false
	}

	return isEnabled
}

func handler(w http.ResponseWriter, r *http.Request) {
	if isAuthEnabled() {
		ctx := context.Background()

		firestoreClient := setupFirestore(ctx)

		_, err := firestoreClient.Collection("users").Where("API-KEY", "==", r.URL.Query().Get("apikey")).Documents(ctx).Next()

		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("401 - Unauthorized"))
			return
		}
	}

	if strings.HasPrefix(r.RequestURI, "/info") {
		getBucket(w, r)
		return
	}

	if strings.HasPrefix(r.RequestURI, "/download") {
		downloadFile(w, r)
		return
	}

	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("404 - Endpoint not found"))
}

func getBucket(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	client, err := storage.NewClient(ctx)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	bucketName := r.URL.Query().Get("bucket")

	if bucketName == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("400 - Bucket name was not specified"))
		return
	}

	bucket := client.Bucket(bucketName)

	query := &storage.Query{}

	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	objects := bucket.Objects(ctx, query)

	var result []*storage.ObjectAttrs

	for {
		obj, err := objects.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("404 - Bucket not found or is inaccessible"))
			return
		}

		result = append(result, obj)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func downloadFile(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	client, err := storage.NewClient(ctx)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	bucketName := r.URL.Query().Get("bucket")
	fileName := r.URL.Query().Get("filename")

	if bucketName == "" || fileName == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("400 - Bucket or file name were not specified"))
		return
	}

	bucket := client.Bucket(bucketName)

	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	file, err := bucket.Object(fileName).NewReader(ctx)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("400 - Cannot open requested file"))
		return
	}

	defer file.Close()

	bytes, err := ioutil.ReadAll(file)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("500 - Unable to read data from bucket"))
		return
	}

	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fileName))
	w.Write(bytes)
}

func setupFirestore(ctx context.Context) *firestore.Client {
	projectID := os.Getenv("PROJECT_ID")

	client, err := firestore.NewClient(ctx, projectID)

	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	return client
}
