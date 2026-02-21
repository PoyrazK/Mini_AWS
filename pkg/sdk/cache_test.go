package sdk_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/poyrazk/thecloud/pkg/sdk"
	"github.com/stretchr/testify/assert"
)

const (
	cacheTestAPIKey      = "test-api-key"
	cacheTestName        = "my-cache"
	cacheTestID          = "cache-1"
	cacheTestKey         = "foo"
	cacheTestVal         = "bar"
	cacheTestBasePath    = "/api/v1/caches"
	cacheTestContentType = "Content-Type"
	cacheTestAppJSON     = "application/json"
)

func newCacheTestServer(t *testing.T) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set(cacheTestContentType, cacheTestAppJSON)

		switch {
		case r.URL.Path == cacheTestBasePath && r.Method == http.MethodPost:
			var body map[string]interface{}
			_ = json.NewDecoder(r.Body).Decode(&body)
			if body["name"] == cacheTestName {
				w.WriteHeader(http.StatusCreated)
				_ = json.NewEncoder(w).Encode(map[string]interface{}{
					"data": map[string]interface{}{
						"id":   cacheTestID,
						"name": cacheTestName,
					},
				})
			}
		case r.URL.Path == cacheTestBasePath && r.Method == http.MethodGet:
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"data": []map[string]interface{}{
					{"id": cacheTestID, "name": cacheTestName},
				},
			})
		case r.URL.Path == cacheTestBasePath+"/"+cacheTestID && r.Method == http.MethodDelete:
			w.WriteHeader(http.StatusNoContent)
		case r.URL.Path == fmt.Sprintf("%s/%s/items/%s", cacheTestBasePath, cacheTestID, cacheTestKey):
			if r.Method == http.MethodPut {
				w.WriteHeader(http.StatusNoContent)
			} else if r.Method == http.MethodGet {
				w.WriteHeader(http.StatusOK)
				_ = json.NewEncoder(w).Encode(map[string]interface{}{
					"data": cacheTestVal,
				})
			} else if r.Method == http.MethodDelete {
				w.WriteHeader(http.StatusNoContent)
			}
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
}

func TestCacheSDK(t *testing.T) {
	ts := newCacheTestServer(t)
	defer ts.Close()

	client := sdk.NewClient(ts.URL+"/api/v1", cacheTestAPIKey)

	t.Run("CreateCache", func(t *testing.T) {
		c, err := client.CreateCache(cacheTestName, "redis", 100, nil)
		assert.NoError(t, err)
		if c != nil {
			assert.Equal(t, cacheTestName, c.Name)
		}
	})

	t.Run("ListCaches", func(t *testing.T) {
		cs, err := client.ListCaches()
		assert.NoError(t, err)
		assert.Len(t, cs, 1)
	})

	t.Run("SetCacheItem", func(t *testing.T) {
		err := client.SetCacheItem(cacheTestID, cacheTestKey, cacheTestVal, 0)
		assert.NoError(t, err)
	})

	t.Run("GetCacheItem", func(t *testing.T) {
		val, err := client.GetCacheItem(cacheTestID, cacheTestKey)
		assert.NoError(t, err)
		assert.Equal(t, cacheTestVal, val)
	})

	t.Run("DeleteCacheItem", func(t *testing.T) {
		err := client.DeleteCacheItem(cacheTestID, cacheTestKey)
		assert.NoError(t, err)
	})

	t.Run("DeleteCache", func(t *testing.T) {
		err := client.DeleteCache(cacheTestID)
		assert.NoError(t, err)
	})
}

func TestCacheSDKErrors(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	client := sdk.NewClient(server.URL+"/api/v1", cacheTestAPIKey)

	_, err := client.CreateCache("c", "redis", 10, nil)
	assert.Error(t, err)

	_, err = client.ListCaches()
	assert.Error(t, err)

	err = client.SetCacheItem(cacheTestID, "k", "v", 0)
	assert.Error(t, err)

	_, err = client.GetCacheItem(cacheTestID, "k")
	assert.Error(t, err)

	err = client.DeleteCacheItem(cacheTestID, "k")
	assert.Error(t, err)

	err = client.DeleteCache(cacheTestID)
	assert.Error(t, err)
}
