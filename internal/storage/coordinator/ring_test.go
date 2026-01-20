package coordinator

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConsistentHashRing(t *testing.T) {
	ring := NewConsistentHashRing(10) // Small virtual nodes for predictable testing

	// 1. Add Nodes
	ring.AddNode("node-A")
	ring.AddNode("node-B")
	ring.AddNode("node-C")

	// 2. Test GetNodes distribution
	key1 := "bucket1/object1"
	nodes1 := ring.GetNodes(key1, 3)
	assert.Len(t, nodes1, 3)
	assert.Contains(t, nodes1, "node-A")
	assert.Contains(t, nodes1, "node-B")
	assert.Contains(t, nodes1, "node-C")

	// 3. Test Consistency (Same key -> Same nodes)
	nodes2 := ring.GetNodes(key1, 3)
	assert.Equal(t, nodes1, nodes2)

	// 4. Test Single Node Request
	primary := ring.GetNodes("some-key", 1)
	assert.Len(t, primary, 1)

	// 5. Remove Node
	ring.RemoveNode("node-B")
	nodesAfterRemoval := ring.GetNodes(key1, 3)
	assert.Len(t, nodesAfterRemoval, 2) // Only A and C left
	assert.NotContains(t, nodesAfterRemoval, "node-B")
}

func TestRingDistribution(t *testing.T) {
	ring := NewConsistentHashRing(100)
	nodes := []string{"node-1", "node-2", "node-3", "node-4", "node-5"}
	for _, n := range nodes {
		ring.AddNode(n)
	}

	counts := make(map[string]int)
	totalKeys := 1000

	for i := 0; i < totalKeys; i++ {
		key := "key-" + strconv.Itoa(i)
		primary := ring.GetNodes(key, 1)[0]
		counts[primary]++
	}

	// Ideally, with 100 virtual nodes, variance should be acceptable.
	// We just assert that every node got *some* keys.
	for _, n := range nodes {
		assert.Greater(t, counts[n], 0, "Node %s received 0 keys", n)
	}
}
