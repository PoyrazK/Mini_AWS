package sdk

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestClient_CreateSecurityGroup(t *testing.T) {
	vpcID := "vpc-123"
	expectedSG := SecurityGroup{
		ID:          "sg-1",
		VPCID:       vpcID,
		Name:        "test-sg",
		Description: "test security group",
		CreatedAt:   time.Now(),
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/security-groups", r.URL.Path)
		assert.Equal(t, http.MethodPost, r.Method)

		var req map[string]string
		err := json.NewDecoder(r.Body).Decode(&req)
		assert.NoError(t, err)
		assert.Equal(t, vpcID, req["vpc_id"])
		assert.Equal(t, expectedSG.Name, req["name"])

		w.Header().Set("Content-Type", "application/json")
		resp := Response[SecurityGroup]{Data: expectedSG}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := NewClient(server.URL, "test-api-key")
	sg, err := client.CreateSecurityGroup(vpcID, "test-sg", "test security group")

	assert.NoError(t, err)
	assert.NotNil(t, sg)
	assert.Equal(t, expectedSG.ID, sg.ID)
}

func TestClient_ListSecurityGroups(t *testing.T) {
	vpcID := "vpc-123"
	expectedSGs := []SecurityGroup{
		{ID: "sg-1", Name: "sg-1", VPCID: vpcID},
		{ID: "sg-2", Name: "sg-2", VPCID: vpcID},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/security-groups", r.URL.Path)
		assert.Equal(t, "vpc_id="+vpcID, r.URL.RawQuery)
		assert.Equal(t, http.MethodGet, r.Method)

		w.Header().Set("Content-Type", "application/json")
		resp := Response[[]SecurityGroup]{Data: expectedSGs}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := NewClient(server.URL, "test-api-key")
	sgs, err := client.ListSecurityGroups(vpcID)

	assert.NoError(t, err)
	assert.Len(t, sgs, 2)
}

func TestClient_GetSecurityGroup(t *testing.T) {
	id := "sg-123"
	expectedSG := SecurityGroup{
		ID:   id,
		Name: "test-sg",
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/security-groups/"+id, r.URL.Path)
		assert.Equal(t, http.MethodGet, r.Method)

		w.Header().Set("Content-Type", "application/json")
		resp := Response[SecurityGroup]{Data: expectedSG}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := NewClient(server.URL, "test-api-key")
	sg, err := client.GetSecurityGroup(id)

	assert.NoError(t, err)
	assert.NotNil(t, sg)
	assert.Equal(t, expectedSG.ID, sg.ID)
}

func TestClient_DeleteSecurityGroup(t *testing.T) {
	id := "sg-123"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/security-groups/"+id, r.URL.Path)
		assert.Equal(t, http.MethodDelete, r.Method)

		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	client := NewClient(server.URL, "test-api-key")
	err := client.DeleteSecurityGroup(id)

	assert.NoError(t, err)
}

func TestClient_AddSecurityRule(t *testing.T) {
	groupID := "sg-123"
	rule := SecurityRule{
		Direction: "ingress",
		Protocol:  "tcp",
		PortMin:   80,
		PortMax:   80,
		CIDR:      "0.0.0.0/0",
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/security-groups/"+groupID+"/rules", r.URL.Path)
		assert.Equal(t, http.MethodPost, r.Method)

		var req SecurityRule
		err := json.NewDecoder(r.Body).Decode(&req)
		assert.NoError(t, err)
		assert.Equal(t, rule.Protocol, req.Protocol)

		w.Header().Set("Content-Type", "application/json")
		rule.ID = "rule-1"
		resp := Response[SecurityRule]{Data: rule}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := NewClient(server.URL, "test-api-key")
	result, err := client.AddSecurityRule(groupID, rule)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "rule-1", result.ID)
}

func TestClient_AttachSecurityGroup(t *testing.T) {
	instanceID := "inst-123"
	groupID := "sg-123"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/security-groups/attach", r.URL.Path)
		assert.Equal(t, http.MethodPost, r.Method)

		var req map[string]string
		err := json.NewDecoder(r.Body).Decode(&req)
		assert.NoError(t, err)
		assert.Equal(t, instanceID, req["instance_id"])
		assert.Equal(t, groupID, req["group_id"])

		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := NewClient(server.URL, "test-api-key")
	err := client.AttachSecurityGroup(instanceID, groupID)

	assert.NoError(t, err)
}
