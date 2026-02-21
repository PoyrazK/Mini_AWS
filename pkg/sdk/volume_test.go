package sdk

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

const (
	volumeAPIKey          = "test-key"
	volumeName            = "test-vol"
	volumeNewName         = "new-name"
	volumeContentType     = "Content-Type"
	volumeApplicationJSON = "application/json"
	volumePath            = "/api/v1/volumes"
)

var volumeID = uuid.New()

func TestClientVolume(t *testing.T) {
	mockVolume := Volume{
		ID:     volumeID,
		Name:   volumeName,
		SizeGB: 10,
		Status: "available",
	}

	t.Run("ListVolumes", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, volumePath, r.URL.Path)
			assert.Equal(t, "GET", r.Method)
			w.Header().Set(volumeContentType, volumeApplicationJSON)
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(Response[[]Volume]{Data: []Volume{mockVolume}})
		}))
		defer server.Close()

		client := NewClient(server.URL+"/api/v1", volumeAPIKey)
		vols, err := client.ListVolumes()
		assert.NoError(t, err)
		assert.Len(t, vols, 1)
		assert.Equal(t, volumeID, vols[0].ID)
	})

	t.Run("GetVolume", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, volumePath+"/"+volumeID.String(), r.URL.Path)
			w.Header().Set(volumeContentType, volumeApplicationJSON)
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(Response[Volume]{Data: mockVolume})
		}))
		defer server.Close()

		client := NewClient(server.URL+"/api/v1", volumeAPIKey)
		vol, err := client.GetVolume(volumeID.String())
		assert.NoError(t, err)
		assert.Equal(t, volumeID, vol.ID)
	})

	t.Run("CreateVolume", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, volumePath, r.URL.Path)
			assert.Equal(t, "POST", r.Method)

			var body map[string]interface{}
			_ = json.NewDecoder(r.Body).Decode(&body)
			assert.Equal(t, volumeNewName, body["name"])
			assert.InDelta(t, float64(20), body["size_gb"], 0.01)

			w.Header().Set(volumeContentType, volumeApplicationJSON)
			w.WriteHeader(http.StatusCreated)
			_ = json.NewEncoder(w).Encode(Response[Volume]{Data: mockVolume})
		}))
		defer server.Close()

		client := NewClient(server.URL+"/api/v1", volumeAPIKey)
		vol, err := client.CreateVolume(volumeNewName, 20)
		assert.NoError(t, err)
		assert.NotNil(t, vol)
	})

	t.Run("DeleteVolume", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, volumePath+"/"+volumeID.String(), r.URL.Path)
			assert.Equal(t, "DELETE", r.Method)
			w.WriteHeader(http.StatusNoContent)
		}))
		defer server.Close()

		client := NewClient(server.URL+"/api/v1", volumeAPIKey)
		err := client.DeleteVolume(volumeID.String())
		assert.NoError(t, err)
	})
}

func TestClientVolumeErrors(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("boom"))
	}))
	defer server.Close()

	client := NewClient(server.URL+"/api/v1", volumeAPIKey)

	_, err := client.ListVolumes()
	assert.Error(t, err)

	_, err = client.GetVolume(volumeID.String())
	assert.Error(t, err)

	_, err = client.CreateVolume("v", 10)
	assert.Error(t, err)

	err = client.DeleteVolume(volumeID.String())
	assert.Error(t, err)
}
