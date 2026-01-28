package platform

import "testing"

const testBucketLabel = "bucket-1"

func TestMetricsCollectorsAreUsable(t *testing.T) {
	HTTPRequestsTotal.WithLabelValues("GET", "/health", "200").Inc()
	HTTPRequestDuration.WithLabelValues("GET", "/health").Observe(0.1)
	LBRequestsTotal.WithLabelValues("lb-1").Inc()
	InstancesTotal.WithLabelValues("running", "docker").Set(1)
	AutoScalingCurrentInstances.WithLabelValues("group-1").Set(2)
	AuthAttemptsTotal.WithLabelValues("success").Inc()
	VolumesTotal.WithLabelValues("available").Set(3)
	QueueMessagesTotal.WithLabelValues("queue-1", "send").Inc()
	StorageOperations.WithLabelValues("upload", testBucketLabel, "success").Inc()
	StorageLatency.WithLabelValues("upload", testBucketLabel).Observe(0.05)
	StorageBytesTransferred.WithLabelValues("upload").Add(123)
	StorageBucketObjects.WithLabelValues(testBucketLabel).Set(10)
	StorageBucketBytes.WithLabelValues(testBucketLabel).Set(2048)
}
