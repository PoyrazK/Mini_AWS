package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/poyrazk/thecloud/pkg/sdk"
)

func setupAPIServer(t *testing.T) *httptest.Server {
	t.Helper()

	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/security-groups":
			resp := sdk.Response[[]sdk.SecurityGroup]{
				Data: []sdk.SecurityGroup{
					{
						ID:          "sg-1",
						VPCID:       "vpc-1",
						Name:        "default",
						Description: "default group",
						ARN:         "arn:sg:1",
						CreatedAt:   time.Now().UTC(),
					},
				},
			}
			_ = json.NewEncoder(w).Encode(resp)
		case r.Method == http.MethodPost && r.URL.Path == "/security-groups":
			resp := sdk.Response[sdk.SecurityGroup]{
				Data: sdk.SecurityGroup{
					ID:          "sg-2",
					VPCID:       "vpc-2",
					Name:        "web",
					Description: "web group",
					ARN:         "arn:sg:2",
					CreatedAt:   time.Now().UTC(),
				},
			}
			_ = json.NewEncoder(w).Encode(resp)
		case r.Method == http.MethodGet && r.URL.Path == "/vpcs":
			resp := sdk.Response[[]sdk.VPC]{
				Data: []sdk.VPC{
					{
						ID:        "vpc-1",
						Name:      "main",
						CIDRBlock: "10.0.0.0/16",
						NetworkID: "net-1",
						VXLANID:   1001,
						Status:    "available",
						CreatedAt: time.Now().UTC(),
					},
				},
			}
			_ = json.NewEncoder(w).Encode(resp)
		case r.Method == http.MethodPost && r.URL.Path == "/vpcs":
			resp := sdk.Response[sdk.VPC]{
				Data: sdk.VPC{
					ID:        "vpc-2",
					Name:      "demo",
					CIDRBlock: "10.1.0.0/16",
					NetworkID: "net-2",
					VXLANID:   1002,
					Status:    "available",
					CreatedAt: time.Now().UTC(),
				},
			}
			_ = json.NewEncoder(w).Encode(resp)
		case r.Method == http.MethodGet && r.URL.Path == "/vpcs/vpc-1/subnets":
			resp := sdk.Response[[]*sdk.Subnet]{
				Data: []*sdk.Subnet{
					{
						ID:        "subnet-1",
						VpcID:     "vpc-1",
						Name:      "public",
						CIDRBlock: "10.0.1.0/24",
						AZ:        "us-east-1a",
						GatewayIP: "10.0.1.1",
						Status:    "available",
						CreatedAt: time.Now().UTC(),
					},
				},
			}
			_ = json.NewEncoder(w).Encode(resp)
		case r.Method == http.MethodPost && r.URL.Path == "/vpcs/vpc-1/subnets":
			resp := sdk.Response[*sdk.Subnet]{
				Data: &sdk.Subnet{
					ID:        "subnet-2",
					VpcID:     "vpc-1",
					Name:      "private",
					CIDRBlock: "10.0.2.0/24",
					AZ:        "us-east-1b",
					GatewayIP: "10.0.2.1",
					Status:    "available",
					CreatedAt: time.Now().UTC(),
				},
			}
			_ = json.NewEncoder(w).Encode(resp)
		case r.Method == http.MethodGet && r.URL.Path == "/volumes":
			volID := uuid.New()
			resp := sdk.Response[[]sdk.Volume]{
				Data: []sdk.Volume{
					{
						ID:        volID,
						Name:      "data",
						SizeGB:    20,
						Status:    "available",
						CreatedAt: time.Now().UTC(),
						UpdatedAt: time.Now().UTC(),
					},
				},
			}
			_ = json.NewEncoder(w).Encode(resp)
		case r.Method == http.MethodPost && r.URL.Path == "/volumes":
			volID := uuid.New()
			resp := sdk.Response[sdk.Volume]{
				Data: sdk.Volume{
					ID:        volID,
					Name:      "data",
					SizeGB:    20,
					Status:    "available",
					CreatedAt: time.Now().UTC(),
					UpdatedAt: time.Now().UTC(),
				},
			}
			_ = json.NewEncoder(w).Encode(resp)
		case r.Method == http.MethodDelete && r.URL.Path == "/volumes/vol-1":
			w.WriteHeader(http.StatusNoContent)
		case r.Method == http.MethodDelete && r.URL.Path == "/vpcs/vpc-1":
			w.WriteHeader(http.StatusNoContent)
		case r.Method == http.MethodDelete && r.URL.Path == "/subnets/subnet-1":
			w.WriteHeader(http.StatusNoContent)
		case r.Method == http.MethodPost && r.URL.Path == "/caches":
			vpcID := "vpc-1"
			resp := sdk.Response[sdk.Cache]{
				Data: sdk.Cache{
					ID:        "cache-1",
					Name:      "redis-main",
					Engine:    "redis",
					Version:   "7.2",
					Status:    "PROVISIONING",
					VpcID:     &vpcID,
					Port:      6379,
					MemoryMB:  128,
					CreatedAt: time.Now().UTC(),
					UpdatedAt: time.Now().UTC(),
				},
			}
			_ = json.NewEncoder(w).Encode(resp)
		case r.Method == http.MethodGet && r.URL.Path == "/caches":
			vpcID := "vpc-1"
			resp := sdk.Response[[]*sdk.Cache]{
				Data: []*sdk.Cache{
					{
						ID:       "cache-1",
						Name:     "redis-main",
						Engine:   "redis",
						Version:  "7.2",
						Status:   "RUNNING",
						VpcID:    &vpcID,
						Port:     6379,
						MemoryMB: 128,
					},
				},
			}
			_ = json.NewEncoder(w).Encode(resp)
		case r.Method == http.MethodGet && r.URL.Path == "/caches/cache-1":
			vpcID := "vpc-1"
			resp := sdk.Response[sdk.Cache]{
				Data: sdk.Cache{
					ID:       "cache-1",
					Name:     "redis-main",
					Engine:   "redis",
					Version:  "7.2",
					Status:   "RUNNING",
					VpcID:    &vpcID,
					Port:     6379,
					MemoryMB: 128,
				},
			}
			_ = json.NewEncoder(w).Encode(resp)
		case r.Method == http.MethodDelete && r.URL.Path == "/caches/cache-1":
			w.WriteHeader(http.StatusNoContent)
		case r.Method == http.MethodGet && r.URL.Path == "/caches/cache-1/connection":
			resp := sdk.Response[map[string]string]{
				Data: map[string]string{"connection_string": "redis://cache-1:6379"},
			}
			_ = json.NewEncoder(w).Encode(resp)
		case r.Method == http.MethodPost && r.URL.Path == "/caches/cache-1/flush":
			w.WriteHeader(http.StatusNoContent)
		case r.Method == http.MethodGet && r.URL.Path == "/caches/cache-1/stats":
			resp := sdk.Response[sdk.CacheStats]{
				Data: sdk.CacheStats{
					UsedMemoryBytes:  1024,
					MaxMemoryBytes:   2048,
					ConnectedClients: 5,
					TotalKeys:        10,
				},
			}
			_ = json.NewEncoder(w).Encode(resp)
		case r.Method == http.MethodPost && r.URL.Path == "/queues":
			resp := sdk.Response[sdk.Queue]{
				Data: sdk.Queue{
					ID:     "queue-1",
					Name:   "jobs",
					ARN:    "arn:queue:1",
					Status: "ACTIVE",
				},
			}
			_ = json.NewEncoder(w).Encode(resp)
		case r.Method == http.MethodGet && r.URL.Path == "/queues":
			resp := sdk.Response[[]sdk.Queue]{
				Data: []sdk.Queue{
					{
						ID:     "queue-1",
						Name:   "jobs",
						ARN:    "arn:queue:1",
						Status: "ACTIVE",
					},
				},
			}
			_ = json.NewEncoder(w).Encode(resp)
		case r.Method == http.MethodDelete && r.URL.Path == "/queues/queue-1":
			w.WriteHeader(http.StatusNoContent)
		case r.Method == http.MethodPost && r.URL.Path == "/queues/queue-1/messages":
			resp := sdk.Response[sdk.Message]{
				Data: sdk.Message{
					ID:            "msg-1",
					QueueID:       "queue-1",
					Body:          "hello",
					ReceiptHandle: "rh-1",
					CreatedAt:     time.Now().UTC(),
				},
			}
			_ = json.NewEncoder(w).Encode(resp)
		case r.Method == http.MethodGet && r.URL.Path == "/queues/queue-1/messages":
			resp := sdk.Response[[]sdk.Message]{
				Data: []sdk.Message{
					{
						ID:            "msg-1",
						QueueID:       "queue-1",
						Body:          "hello",
						ReceiptHandle: "rh-1",
						CreatedAt:     time.Now().UTC(),
					},
				},
			}
			_ = json.NewEncoder(w).Encode(resp)
		case r.Method == http.MethodDelete && r.URL.Path == "/queues/queue-1/messages/rh-1":
			w.WriteHeader(http.StatusNoContent)
		case r.Method == http.MethodPost && r.URL.Path == "/queues/queue-1/purge":
			w.WriteHeader(http.StatusNoContent)
		case r.Method == http.MethodPost && r.URL.Path == "/notify/topics":
			resp := sdk.Response[sdk.Topic]{
				Data: sdk.Topic{
					ID:   "topic-1",
					Name: "updates",
					ARN:  "arn:topic:1",
				},
			}
			_ = json.NewEncoder(w).Encode(resp)
		case r.Method == http.MethodGet && r.URL.Path == "/notify/topics":
			resp := sdk.Response[[]sdk.Topic]{
				Data: []sdk.Topic{{ID: "topic-1", Name: "updates", ARN: "arn:topic:1"}},
			}
			_ = json.NewEncoder(w).Encode(resp)
		case r.Method == http.MethodPost && r.URL.Path == "/notify/topics/topic-1/subscriptions":
			resp := sdk.Response[sdk.Subscription]{
				Data: sdk.Subscription{
					ID:       "sub-1",
					TopicID:  "topic-1",
					Protocol: "webhook",
					Endpoint: "https://example.com",
				},
			}
			_ = json.NewEncoder(w).Encode(resp)
		case r.Method == http.MethodPost && r.URL.Path == "/notify/topics/topic-1/publish":
			w.WriteHeader(http.StatusNoContent)
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
}

func setAPIContext(t *testing.T, server *httptest.Server) {
	t.Helper()

	oldURL := apiURL
	oldKey := apiKey
	oldJSON := outputJSON

	apiURL = server.URL
	apiKey = "test-key"
	outputJSON = false

	t.Cleanup(func() {
		apiURL = oldURL
		apiKey = oldKey
		outputJSON = oldJSON
	})
}

func TestSGListCommand_JSONOutput(t *testing.T) {
	server := setupAPIServer(t)
	defer server.Close()
	setAPIContext(t, server)

	_ = sgListCmd.Flags().Set(flagVPCID, "vpc-1")
	outputJSON = true
	t.Cleanup(func() {
		_ = sgListCmd.Flags().Set(flagVPCID, "")
		outputJSON = false
	})

	out := captureStdout(t, func() {
		sgListCmd.Run(sgListCmd, nil)
	})
	if !strings.Contains(out, "\"id\": \"sg-1\"") {
		t.Fatalf("expected security group in output, got: %s", out)
	}
}

func TestSGCreateCommand_MissingVPC(t *testing.T) {
	out := captureStdout(t, func() {
		sgCreateCmd.Run(sgCreateCmd, []string{"web"})
	})
	if !strings.Contains(out, "--vpc-id is required") {
		t.Fatalf("expected missing vpc-id error, got: %s", out)
	}
}

func TestVPCCreateCommand_Success(t *testing.T) {
	server := setupAPIServer(t)
	defer server.Close()
	setAPIContext(t, server)

	out := captureStdout(t, func() {
		vpcCreateCmd.Run(vpcCreateCmd, []string{"demo"})
	})
	if !strings.Contains(out, "VPC demo created successfully") {
		t.Fatalf("expected success message, got: %s", out)
	}
}

func TestVPCListCommand_JSONOutput(t *testing.T) {
	server := setupAPIServer(t)
	defer server.Close()
	setAPIContext(t, server)

	outputJSON = true
	t.Cleanup(func() { outputJSON = false })

	out := captureStdout(t, func() {
		vpcListCmd.Run(vpcListCmd, nil)
	})
	if !strings.Contains(out, "\"id\": \"vpc-1\"") {
		t.Fatalf("expected VPC in output, got: %s", out)
	}
}

func TestSubnetListCommand_JSONOutput(t *testing.T) {
	server := setupAPIServer(t)
	defer server.Close()
	setAPIContext(t, server)

	outputJSON = true
	t.Cleanup(func() { outputJSON = false })

	out := captureStdout(t, func() {
		subnetListCmd.Run(subnetListCmd, []string{"vpc-1"})
	})
	if !strings.Contains(out, "\"id\": \"subnet-1\"") {
		t.Fatalf("expected subnet in output, got: %s", out)
	}
}

func TestVolumeCreateCommand_JSONOutput(t *testing.T) {
	server := setupAPIServer(t)
	defer server.Close()
	setAPIContext(t, server)

	_ = volumeCreateCmd.Flags().Set("name", "data")
	_ = volumeCreateCmd.Flags().Set("size", "20")
	defer func() {
		_ = volumeCreateCmd.Flags().Set("name", "")
		_ = volumeCreateCmd.Flags().Set("size", "1")
	}()

	out := captureStdout(t, func() {
		volumeCreateCmd.Run(volumeCreateCmd, nil)
	})
	if !strings.Contains(out, "Volume created") {
		t.Fatalf("expected volume create output, got: %s", out)
	}
	if !strings.Contains(out, "\"name\": \"data\"") {
		t.Fatalf("expected volume json output, got: %s", out)
	}
}

func TestVolumeListCommand_JSONOutput(t *testing.T) {
	server := setupAPIServer(t)
	defer server.Close()
	setAPIContext(t, server)

	outputJSON = true
	t.Cleanup(func() { outputJSON = false })

	out := captureStdout(t, func() {
		volumeListCmd.Run(volumeListCmd, nil)
	})
	if !strings.Contains(out, "\"name\": \"data\"") {
		t.Fatalf("expected volume list output, got: %s", out)
	}
}

func TestVolumeDeleteCommand_Success(t *testing.T) {
	server := setupAPIServer(t)
	defer server.Close()
	setAPIContext(t, server)

	out := captureStdout(t, func() {
		volumeDeleteCmd.Run(volumeDeleteCmd, []string{"vol-1"})
	})
	if !strings.Contains(out, "Volume vol-1 deleted") {
		t.Fatalf("expected volume delete output, got: %s", out)
	}
}

func TestVPCRmCommand_Success(t *testing.T) {
	server := setupAPIServer(t)
	defer server.Close()
	setAPIContext(t, server)

	out := captureStdout(t, func() {
		vpcRmCmd.Run(vpcRmCmd, []string{"vpc-1"})
	})
	if !strings.Contains(out, "VPC vpc-1 removed") {
		t.Fatalf("expected VPC remove output, got: %s", out)
	}
}

func TestSubnetDeleteCommand_Success(t *testing.T) {
	server := setupAPIServer(t)
	defer server.Close()
	setAPIContext(t, server)

	out := captureStdout(t, func() {
		subnetDeleteCmd.Run(subnetDeleteCmd, []string{"subnet-1"})
	})
	if !strings.Contains(out, "deleted successfully") {
		t.Fatalf("expected subnet delete output, got: %s", out)
	}
}

func TestCacheCreateCommand_Success(t *testing.T) {
	server := setupAPIServer(t)
	defer server.Close()
	setAPIContext(t, server)

	_ = createCacheCmd.Flags().Set("name", "redis-main")
	_ = createCacheCmd.Flags().Set("version", "7.2")
	_ = createCacheCmd.Flags().Set("memory", "128")
	_ = createCacheCmd.Flags().Set("vpc", "vpc-1")
	_ = createCacheCmd.Flags().Set("wait", "false")
	defer func() {
		_ = createCacheCmd.Flags().Set("name", "")
		_ = createCacheCmd.Flags().Set("version", "7.2")
		_ = createCacheCmd.Flags().Set("memory", "128")
		_ = createCacheCmd.Flags().Set("vpc", "")
		_ = createCacheCmd.Flags().Set("wait", "false")
	}()

	out := captureStdout(t, func() {
		createCacheCmd.Run(createCacheCmd, nil)
	})
	if !strings.Contains(out, "Cache created with ID: cache-1") {
		t.Fatalf("expected cache create output, got: %s", out)
	}
}

func TestCacheListCommand_Output(t *testing.T) {
	server := setupAPIServer(t)
	defer server.Close()
	setAPIContext(t, server)

	out := captureStdout(t, func() {
		listCacheCmd.Run(listCacheCmd, nil)
	})
	if !strings.Contains(out, "cache-1") {
		t.Fatalf("expected cache list output, got: %s", out)
	}
}

func TestCacheShowCommand_Output(t *testing.T) {
	server := setupAPIServer(t)
	defer server.Close()
	setAPIContext(t, server)

	out := captureStdout(t, func() {
		getCacheCmd.Run(getCacheCmd, []string{"cache-1"})
	})
	if !strings.Contains(out, "Name:      redis-main") {
		t.Fatalf("expected cache show output, got: %s", out)
	}
}

func TestCacheDeleteCommand_Output(t *testing.T) {
	server := setupAPIServer(t)
	defer server.Close()
	setAPIContext(t, server)

	out := captureStdout(t, func() {
		deleteCacheCmd.Run(deleteCacheCmd, []string{"cache-1"})
	})
	if !strings.Contains(out, "Cache deleted successfully") {
		t.Fatalf("expected cache delete output, got: %s", out)
	}
}

func TestCacheConnectionCommand_Output(t *testing.T) {
	server := setupAPIServer(t)
	defer server.Close()
	setAPIContext(t, server)

	out := captureStdout(t, func() {
		connectionCacheCmd.Run(connectionCacheCmd, []string{"cache-1"})
	})
	if !strings.Contains(out, "redis://cache-1:6379") {
		t.Fatalf("expected cache connection output, got: %s", out)
	}
}

func TestCacheStatsCommand_Output(t *testing.T) {
	server := setupAPIServer(t)
	defer server.Close()
	setAPIContext(t, server)

	out := captureStdout(t, func() {
		statsCacheCmd.Run(statsCacheCmd, []string{"cache-1"})
	})
	if !strings.Contains(out, "Used Memory: 1.0 KB") {
		t.Fatalf("expected cache stats output, got: %s", out)
	}
}

func TestCacheFlushCommand_Output(t *testing.T) {
	server := setupAPIServer(t)
	defer server.Close()
	setAPIContext(t, server)

	_ = flushCacheCmd.Flags().Set("yes", "true")
	t.Cleanup(func() {
		_ = flushCacheCmd.Flags().Set("yes", "false")
	})

	out := captureStdout(t, func() {
		flushCacheCmd.Run(flushCacheCmd, []string{"cache-1"})
	})
	if !strings.Contains(out, "Cache flushed successfully") {
		t.Fatalf("expected cache flush output, got: %s", out)
	}
}

func TestQueueListCommand_JSONOutput(t *testing.T) {
	server := setupAPIServer(t)
	defer server.Close()
	setAPIContext(t, server)

	outputJSON = true
	t.Cleanup(func() { outputJSON = false })

	out := captureStdout(t, func() {
		listQueuesCmd.Run(listQueuesCmd, nil)
	})
	if !strings.Contains(out, "\"id\": \"queue-1\"") {
		t.Fatalf("expected queue list output, got: %s", out)
	}
}

func TestQueueCreateCommand_Output(t *testing.T) {
	server := setupAPIServer(t)
	defer server.Close()
	setAPIContext(t, server)

	out := captureStdout(t, func() {
		createQueueCmd.Run(createQueueCmd, []string{"jobs"})
	})
	if !strings.Contains(out, "Queue created successfully") {
		t.Fatalf("expected queue create output, got: %s", out)
	}
}

func TestQueueDeleteCommand_Output(t *testing.T) {
	server := setupAPIServer(t)
	defer server.Close()
	setAPIContext(t, server)

	out := captureStdout(t, func() {
		deleteQueueCmd.Run(deleteQueueCmd, []string{"queue-1"})
	})
	if !strings.Contains(out, "Queue deleted successfully") {
		t.Fatalf("expected queue delete output, got: %s", out)
	}
}

func TestQueueSendMessage_Output(t *testing.T) {
	server := setupAPIServer(t)
	defer server.Close()
	setAPIContext(t, server)

	out := captureStdout(t, func() {
		sendMessageCmd.Run(sendMessageCmd, []string{"queue-1", "hello"})
	})
	if !strings.Contains(out, "Message sent") {
		t.Fatalf("expected send message output, got: %s", out)
	}
}

func TestQueueReceiveMessages_JSONOutput(t *testing.T) {
	server := setupAPIServer(t)
	defer server.Close()
	setAPIContext(t, server)

	_ = receiveMessagesCmd.Flags().Set("max", "1")
	outputJSON = true
	t.Cleanup(func() {
		_ = receiveMessagesCmd.Flags().Set("max", "1")
		outputJSON = false
	})

	out := captureStdout(t, func() {
		receiveMessagesCmd.Run(receiveMessagesCmd, []string{"queue-1"})
	})
	if !strings.Contains(out, "\"id\": \"msg-1\"") {
		t.Fatalf("expected receive messages output, got: %s", out)
	}
}

func TestQueueAckMessage_Output(t *testing.T) {
	server := setupAPIServer(t)
	defer server.Close()
	setAPIContext(t, server)

	out := captureStdout(t, func() {
		ackMessageCmd.Run(ackMessageCmd, []string{"queue-1", "rh-1"})
	})
	if !strings.Contains(out, "Message acknowledged") {
		t.Fatalf("expected ack output, got: %s", out)
	}
}

func TestQueuePurge_Output(t *testing.T) {
	server := setupAPIServer(t)
	defer server.Close()
	setAPIContext(t, server)

	out := captureStdout(t, func() {
		purgeQueueCmd.Run(purgeQueueCmd, []string{"queue-1"})
	})
	if !strings.Contains(out, "Queue purged") {
		t.Fatalf("expected queue purge output, got: %s", out)
	}
}

func TestNotifyCreateTopic_Output(t *testing.T) {
	server := setupAPIServer(t)
	defer server.Close()
	setAPIContext(t, server)

	out := captureStdout(t, func() {
		createTopicCmd.Run(createTopicCmd, []string{"updates"})
	})
	if !strings.Contains(out, "Topic created") {
		t.Fatalf("expected create topic output, got: %s", out)
	}
}

func TestNotifyListTopics_Output(t *testing.T) {
	server := setupAPIServer(t)
	defer server.Close()
	setAPIContext(t, server)

	out := captureStdout(t, func() {
		listTopicsCmd.Run(listTopicsCmd, nil)
	})
	if !strings.Contains(out, "updates") {
		t.Fatalf("expected list topics output, got: %s", out)
	}
}

func TestNotifySubscribe_Output(t *testing.T) {
	server := setupAPIServer(t)
	defer server.Close()
	setAPIContext(t, server)

	_ = subscribeCmd.Flags().Set("protocol", "webhook")
	_ = subscribeCmd.Flags().Set("endpoint", "https://example.com")
	defer func() {
		_ = subscribeCmd.Flags().Set("endpoint", "")
	}()

	out := captureStdout(t, func() {
		subscribeCmd.Run(subscribeCmd, []string{"topic-1"})
	})
	if !strings.Contains(out, "Subscription created") {
		t.Fatalf("expected subscribe output, got: %s", out)
	}
}

func TestNotifyPublish_Output(t *testing.T) {
	server := setupAPIServer(t)
	defer server.Close()
	setAPIContext(t, server)

	out := captureStdout(t, func() {
		publishCmd.Run(publishCmd, []string{"topic-1", "hello"})
	})
	if !strings.Contains(out, "Message published") {
		t.Fatalf("expected publish output, got: %s", out)
	}
}
