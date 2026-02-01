package tests

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/poyrazk/thecloud/internal/core/domain"
	"github.com/poyrazk/thecloud/pkg/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const gatewayRoutesPath = "/gateway/routes"

func TestGatewayE2E(t *testing.T) {
	if err := waitForServer(); err != nil {
		t.Fatalf("Failing Gateway E2E test: %v", err)
	}

	client := &http.Client{Timeout: 15 * time.Second}
	token := registerAndLogin(t, client, "gateway-tester@thecloud.local", "Gateway Tester")

	// Use a unique suffix for route names to avoid collisions in E2E environment
	ts := time.Now().UnixNano() % 100000

	t.Run("CreateAndListPatternRoute", func(t *testing.T) {
		// 1. Create a pattern-based route
		// We'll use httpbin.org to verify the proxying works
		pattern := fmt.Sprintf("/httpbin-%d/{method}", ts)
		routeName := fmt.Sprintf("httpbin-pattern-%d", ts)
		targetURL := "https://httpbin.org"

		payload := map[string]interface{}{
			"name":         routeName,
			"path_prefix":  pattern, // API currently uses path_prefix field for the pattern
			"target_url":   targetURL,
			"strip_prefix": true,
			"rate_limit":   100,
		}

		resp := postRequest(t, client, testutil.TestBaseURL+gatewayRoutesPath, token, payload)
		require.Equal(t, http.StatusCreated, resp.StatusCode)
		defer resp.Body.Close()

		var res struct {
			Data domain.GatewayRoute `json:"data"`
		}
		require.NoError(t, json.NewDecoder(resp.Body).Decode(&res))
		assert.Equal(t, "pattern", res.Data.PatternType)
		assert.Equal(t, pattern, res.Data.PathPattern)

		// 2. List routes and verify it's there
		listResp := getRequest(t, client, testutil.TestBaseURL+gatewayRoutesPath, token)
		require.Equal(t, http.StatusOK, listResp.StatusCode)
		defer listResp.Body.Close()

		var listRes struct {
			Data []domain.GatewayRoute `json:"data"`
		}
		require.NoError(t, json.NewDecoder(listResp.Body).Decode(&listRes))

		found := false
		for _, r := range listRes.Data {
			if r.Name == routeName {
				found = true
				break
			}
		}
		assert.True(t, found)
	})

	t.Run("VerifyPatternProxying", func(t *testing.T) {
		// Give the gateway a moment to refresh routes
		time.Sleep(2 * time.Second)

		// Test GET request through gateway
		req, _ := http.NewRequest("GET", fmt.Sprintf("%s/gw/httpbin-%d/get", testutil.TestBaseURL, ts), nil)
		req.Header.Set(testutil.TestHeaderAPIKey, token)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var httpbinResp struct {
			URL string `json:"url"`
		}
		require.NoError(t, json.NewDecoder(resp.Body).Decode(&httpbinResp))
		assert.Contains(t, httpbinResp.URL, "/get")
	})

	t.Run("VerifyRegexPatternProxying", func(t *testing.T) {
		pattern := fmt.Sprintf("/status-%d/{code:[0-9]+}", ts)
		routeName := fmt.Sprintf("status-code-%d", ts)
		targetURL := "https://httpbin.org/status"

		payload := map[string]interface{}{
			"name":         routeName,
			"path_prefix":  pattern,
			"target_url":   targetURL,
			"strip_prefix": true,
			"rate_limit":   100,
		}

		resp := postRequest(t, client, testutil.TestBaseURL+gatewayRoutesPath, token, payload)
		require.Equal(t, http.StatusCreated, resp.StatusCode)
		resp.Body.Close()

		time.Sleep(2 * time.Second)

		// This should match
		req, _ := http.NewRequest("GET", fmt.Sprintf("%s/gw/status-%d/201", testutil.TestBaseURL, ts), nil)
		resp, err := client.Do(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusCreated, resp.StatusCode)
		resp.Body.Close()

		// This should NOT match (letters instead of numbers)
		req, _ = http.NewRequest("GET", fmt.Sprintf("%s/gw/status-%d/abc", testutil.TestBaseURL, ts), nil)
		resp, err = client.Do(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
		resp.Body.Close()
	})

	t.Run("VerifyWildcardProxying", func(t *testing.T) {
		pattern := fmt.Sprintf("/wild-%d/*", ts)
		routeName := fmt.Sprintf("wildcard-route-%d", ts)
		targetURL := "https://httpbin.org/anything"

		payload := map[string]interface{}{
			"name":         routeName,
			"path_prefix":  pattern,
			"target_url":   targetURL,
			"strip_prefix": true,
			"rate_limit":   100,
		}

		resp := postRequest(t, client, testutil.TestBaseURL+gatewayRoutesPath, token, payload)
		require.Equal(t, http.StatusCreated, resp.StatusCode)
		resp.Body.Close()

		time.Sleep(2 * time.Second)

		req, _ := http.NewRequest("GET", fmt.Sprintf("%s/gw/wild-%d/foo/bar", testutil.TestBaseURL, ts), nil)
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var httpbinResp struct {
			URL string `json:"url"`
		}
		require.NoError(t, json.NewDecoder(resp.Body).Decode(&httpbinResp))
		assert.Contains(t, httpbinResp.URL, "/anything/foo/bar")
	})
}
