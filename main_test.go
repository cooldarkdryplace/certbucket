package certbucket

import (
	"context"
	"os"
	"testing"

	"cloud.google.com/go/storage"
	"golang.org/x/crypto/acme/autocert"
)

var (
	projectID   = os.Getenv("GOOGLE_PROJECT_ID")
	cacheBucket = projectID + "_autocert-cache-bucket"
)

func deleteBucket() {
	ctx := context.Background()

	// If this fails – test will also fail and error will be logged.
	client, _ := storage.NewClient(ctx)

	// Ignore error, if bucket is not there we don't care.
	client.Bucket(cacheBucket).Delete(ctx)
}

func TestMain(m *testing.M) {
	deleteBucket()
	defer deleteBucket()

	os.Exit(m.Run())
}

func TestNewCreatesBucket(t *testing.T) {
	cache, err := New(projectID, cacheBucket)
	if err != nil {
		t.Fatalf("Failed to create cache bucket: %s", err)
	}

	if err := cache.bucket.Delete(context.Background()); err != nil {
		t.Error("Failed to delete bucket:", err)
	}
}

func TestCache(t *testing.T) {
	cache, err := New(projectID, cacheBucket)
	if err != nil {
		t.Fatalf("Failed to create cache: %s", err)
	}

	ctx := context.Background()

	key := "your_cert_key"
	value := []byte("your_cert_data")

	if _, err := cache.Get(ctx, key); err != autocert.ErrCacheMiss {
		t.Errorf("Got err: %v, expected: %q", err, autocert.ErrCacheMiss)
	}

	if err := cache.Delete(ctx, key); err != nil {
		t.Errorf("Got err: %v, expected: nil", err)
	}

	if err := cache.Put(ctx, key, value); err != nil {
		t.Errorf("Unexpected error during cert Put: %s", err)
	}

	actual, err := cache.Get(ctx, key)
	if err != nil {
		t.Errorf("Unexpected error during cert Get: %s", err)
	}

	if string(actual) != string(value) {
		t.Errorf("Got cert value: %s, expected: %s", actual, value)
	}

	if err := cache.Delete(ctx, key); err != nil {
		t.Errorf("Got err: %v, expected: nil", err)
	}
}
