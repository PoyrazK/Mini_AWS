package coordinator

import (
	"bytes"
	"context"
	"errors"
	"io"
	"strings"
	"testing"
	"time"

	pb "github.com/poyrazk/thecloud/internal/storage/protocol"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
)

// MockStorageNodeClient
type MockStorageNodeClient struct {
	mock.Mock
}

func (m *MockStorageNodeClient) Store(ctx context.Context, in *pb.StoreRequest, opts ...grpc.CallOption) (*pb.StoreResponse, error) {
	args := m.Called(ctx, in)
	return args.Get(0).(*pb.StoreResponse), args.Error(1)
}

func (m *MockStorageNodeClient) Retrieve(ctx context.Context, in *pb.RetrieveRequest, opts ...grpc.CallOption) (*pb.RetrieveResponse, error) {
	args := m.Called(ctx, in)
	return args.Get(0).(*pb.RetrieveResponse), args.Error(1)
}

func (m *MockStorageNodeClient) Delete(ctx context.Context, in *pb.DeleteRequest, opts ...grpc.CallOption) (*pb.DeleteResponse, error) {
	args := m.Called(ctx, in)
	return args.Get(0).(*pb.DeleteResponse), args.Error(1)
}

func (m *MockStorageNodeClient) Gossip(ctx context.Context, in *pb.GossipMessage, opts ...grpc.CallOption) (*pb.GossipResponse, error) {
	args := m.Called(ctx, in)
	return args.Get(0).(*pb.GossipResponse), args.Error(1)
}

func (m *MockStorageNodeClient) GetClusterStatus(ctx context.Context, in *pb.Empty, opts ...grpc.CallOption) (*pb.ClusterStatusResponse, error) {
	args := m.Called(ctx, in)
	return args.Get(0).(*pb.ClusterStatusResponse), args.Error(1)
}

func TestCoordinator_Write_Quorum(t *testing.T) {
	ring := NewConsistentHashRing(10)
	ring.AddNode("node-1")
	ring.AddNode("node-2")
	ring.AddNode("node-3")

	client1 := new(MockStorageNodeClient)
	client2 := new(MockStorageNodeClient)
	client3 := new(MockStorageNodeClient)

	clients := map[string]pb.StorageNodeClient{
		"node-1": client1,
		"node-2": client2,
		"node-3": client3,
	}

	coord := NewCoordinator(ring, clients, 3)
	defer coord.Stop()

	// Expect Store calls on all nodes
	// Assume N=3, W=2.
	data := []byte("hello")

	// Setup expectations
	// Note: StoreRequest includes timestamp which changes, so use mock.MatchedBy or ignore it.
	client1.On("Store", mock.Anything, mock.MatchedBy(func(req *pb.StoreRequest) bool {
		return req.Bucket == "b" && req.Key == "k" && string(req.Data) == "hello"
	})).Return(&pb.StoreResponse{Success: true}, nil)

	client2.On("Store", mock.Anything, mock.Anything).Return(&pb.StoreResponse{Success: true}, nil)
	client3.On("Store", mock.Anything, mock.Anything).Return(&pb.StoreResponse{Success: true}, nil)

	n, err := coord.Write(context.Background(), "b", "k", bytes.NewReader(data))
	assert.NoError(t, err)
	assert.Equal(t, int64(5), n)
}

func TestCoordinator_Write_QuorumFailure(t *testing.T) {
	ring := NewConsistentHashRing(10)
	ring.AddNode("node-1")
	ring.AddNode("node-2")
	ring.AddNode("node-3")

	client1 := new(MockStorageNodeClient)
	client2 := new(MockStorageNodeClient)
	client3 := new(MockStorageNodeClient)

	clients := map[string]pb.StorageNodeClient{
		"node-1": client1,
		"node-2": client2,
		"node-3": client3,
	}

	coord := NewCoordinator(ring, clients, 3) // W=2
	defer coord.Stop()

	// 2 nodes fail
	client1.On("Store", mock.Anything, mock.Anything).Return(&pb.StoreResponse{Success: false}, errors.New("failed"))
	client2.On("Store", mock.Anything, mock.Anything).Return(&pb.StoreResponse{Success: false}, errors.New("failed"))
	client3.On("Store", mock.Anything, mock.Anything).Return(&pb.StoreResponse{Success: true}, nil)

	_, err := coord.Write(context.Background(), "b", "k", strings.NewReader("hello"))
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "write quorum failed")
}

func TestCoordinator_Read_Repair(t *testing.T) {
	ring := NewConsistentHashRing(10)
	ring.AddNode("node-1")
	ring.AddNode("node-2")
	ring.AddNode("node-3")

	c1 := new(MockStorageNodeClient)
	c2 := new(MockStorageNodeClient)
	c3 := new(MockStorageNodeClient)

	clients := map[string]pb.StorageNodeClient{"node-1": c1, "node-2": c2, "node-3": c3}
	coord := NewCoordinator(ring, clients, 3)
	defer coord.Stop()

	tsNew := time.Now().UnixNano()
	tsOld := tsNew - 1000

	// Node 1: Latest data
	c1.On("Retrieve", mock.Anything, mock.Anything).Return(&pb.RetrieveResponse{
		Found: true, Data: []byte("new"), Timestamp: tsNew,
	}, nil)

	// Node 2: Old data (needs repair)
	c2.On("Retrieve", mock.Anything, mock.Anything).Return(&pb.RetrieveResponse{
		Found: true, Data: []byte("old"), Timestamp: tsOld,
	}, nil)

	// Node 3: Not found (needs repair)
	c3.On("Retrieve", mock.Anything, mock.Anything).Return(&pb.RetrieveResponse{
		Found: false,
	}, nil)

	// Expect repair writes to Node 2 and Node 3
	c2.On("Store", mock.Anything, mock.MatchedBy(func(req *pb.StoreRequest) bool {
		return string(req.Data) == "new" && req.Timestamp == tsNew
	})).Return(&pb.StoreResponse{Success: true}, nil)

	c3.On("Store", mock.Anything, mock.MatchedBy(func(req *pb.StoreRequest) bool {
		return string(req.Data) == "new" && req.Timestamp == tsNew
	})).Return(&pb.StoreResponse{Success: true}, nil)

	r, err := coord.Read(context.Background(), "b", "k")
	assert.NoError(t, err)

	data, _ := io.ReadAll(r)
	assert.Equal(t, "new", string(data))

	// Wait for async repair
	time.Sleep(100 * time.Millisecond)
	c2.AssertCalled(t, "Store", mock.Anything, mock.Anything)
	c3.AssertCalled(t, "Store", mock.Anything, mock.Anything)
}
