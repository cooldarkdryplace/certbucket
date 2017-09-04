package certbucket

import (
	"context"
	"io/ioutil"
	"strings"

	"cloud.google.com/go/storage"
	"golang.org/x/crypto/acme/autocert"
)

// Cache is an adapter that is used to satisfy autocert Cache interface.
// Under the hood Storage Bucket is used to store certificate.
type Cache struct {
	bucket *storage.BucketHandle
}

// New creates new certificate cache storage in Google Cloud Storage bucket.
// If bucket is not there it will be created.
func New(projectID, bucketName string) (*Cache, error) {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, err
	}

	cache := &Cache{client.Bucket(bucketName)}

	if err := cache.bucket.Create(ctx, projectID, nil); err != nil {
		// For now parsing error text. We do not care if bucket already own by our user.
		if !strings.Contains(err.Error(), "Error 409") {
			return nil, err
		}
	}

	return cache, nil
}

// Get returns a certificate data for the specified key.
// If there's no such key, Get returns autocert.ErrCacheMiss.
func (c *Cache) Get(ctx context.Context, key string) ([]byte, error) {
	rc, err := c.bucket.Object(key).NewReader(ctx)
	if err != nil {
		if err == storage.ErrObjectNotExist {
			return nil, autocert.ErrCacheMiss
		}
		return nil, err
	}
	defer rc.Close()

	return ioutil.ReadAll(rc)
}

// Put stores the data in the bucket under the specified key.
func (c *Cache) Put(ctx context.Context, key string, data []byte) error {
	wc := c.bucket.Object(key).NewWriter(ctx)
	defer wc.Close()

	_, err := wc.Write(data)

	return err
}

// Delete removes a certificate data from the cache under the specified key.
// If there's no such key in the cache, Delete returns nil.
func (c *Cache) Delete(ctx context.Context, key string) error {
	if err := c.bucket.Object(key).Delete(ctx); err != nil {
		if err != storage.ErrObjectNotExist {
			return err
		}
	}

	return nil
}
