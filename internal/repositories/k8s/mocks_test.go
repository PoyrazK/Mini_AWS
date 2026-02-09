package k8s

import (
	"context"
	"io"
	"time"

	"github.com/google/uuid"
	"github.com/poyrazk/thecloud/internal/core/domain"
	"github.com/poyrazk/thecloud/internal/core/ports"
	"github.com/stretchr/testify/mock"
)

type mockInstanceService struct{ mock.Mock }

func (m *mockInstanceService) LaunchInstance(ctx context.Context, name, image, ports, instanceType string, vpcID, subnetID *uuid.UUID, volumes []domain.VolumeAttachment) (*domain.Instance, error) {
	args := m.Called(ctx, name, image, ports, instanceType, vpcID, subnetID, volumes)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Instance), args.Error(1)
}
func (m *mockInstanceService) LaunchInstanceWithOptions(ctx context.Context, opts ports.CreateInstanceOptions) (*domain.Instance, error) {
	args := m.Called(ctx, opts)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Instance), args.Error(1)
}
func (m *mockInstanceService) StartInstance(ctx context.Context, idOrName string) error { return nil }
func (m *mockInstanceService) StopInstance(ctx context.Context, idOrName string) error  { return nil }
func (m *mockInstanceService) ListInstances(ctx context.Context) ([]*domain.Instance, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Instance), args.Error(1)
}
func (m *mockInstanceService) GetInstance(ctx context.Context, idOrName string) (*domain.Instance, error) {
	args := m.Called(ctx, idOrName)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Instance), args.Error(1)
}
func (m *mockInstanceService) GetInstanceLogs(ctx context.Context, idOrName string) (string, error) {
	return "", nil
}
func (m *mockInstanceService) GetInstanceStats(ctx context.Context, idOrName string) (*domain.InstanceStats, error) {
	return nil, nil
}
func (m *mockInstanceService) GetConsoleURL(ctx context.Context, idOrName string) (string, error) {
	return "", nil
}
func (m *mockInstanceService) TerminateInstance(ctx context.Context, idOrName string) error {
	return nil
}
func (m *mockInstanceService) Exec(ctx context.Context, idOrName string, cmd []string) (string, error) {
	args := m.Called(ctx, idOrName, cmd)
	if args.Get(0) == nil {
		return "", args.Error(1)
	}
	return args.String(0), args.Error(1)
}

type mockClusterRepo struct{ mock.Mock }

func (m *mockClusterRepo) Create(ctx context.Context, c *domain.Cluster) error { return nil }
func (m *mockClusterRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.Cluster, error) {
	return nil, nil
}
func (m *mockClusterRepo) ListByUserID(ctx context.Context, userID uuid.UUID) ([]*domain.Cluster, error) {
	return nil, nil
}
func (m *mockClusterRepo) ListAll(ctx context.Context) ([]*domain.Cluster, error) {
	return nil, nil
}
func (m *mockClusterRepo) Update(ctx context.Context, c *domain.Cluster) error      { return nil }
func (m *mockClusterRepo) Delete(ctx context.Context, id uuid.UUID) error           { return nil }
func (m *mockClusterRepo) AddNode(ctx context.Context, n *domain.ClusterNode) error { return nil }
func (m *mockClusterRepo) GetNodes(ctx context.Context, clusterID uuid.UUID) ([]*domain.ClusterNode, error) {
	args := m.Called(ctx, clusterID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.ClusterNode), args.Error(1)
}
func (m *mockClusterRepo) DeleteNode(ctx context.Context, nodeID uuid.UUID) error      { return nil }
func (m *mockClusterRepo) UpdateNode(ctx context.Context, n *domain.ClusterNode) error { return nil }

type MockSecretService struct{ mock.Mock }

func (m *MockSecretService) CreateSecret(ctx context.Context, name, value, description string) (*domain.Secret, error) {
	return nil, nil
}
func (m *MockSecretService) GetSecret(ctx context.Context, id uuid.UUID) (*domain.Secret, error) {
	return nil, nil
}
func (m *MockSecretService) GetSecretByName(ctx context.Context, name string) (*domain.Secret, error) {
	return nil, nil
}
func (m *MockSecretService) ListSecrets(ctx context.Context) ([]*domain.Secret, error) {
	return nil, nil
}
func (m *MockSecretService) DeleteSecret(ctx context.Context, id uuid.UUID) error {
	return nil
}
func (m *MockSecretService) Encrypt(ctx context.Context, userID uuid.UUID, plainText string) (string, error) {
	return plainText, nil
}
func (m *MockSecretService) Decrypt(ctx context.Context, userID uuid.UUID, cipherText string) (string, error) {
	args := m.Called(ctx, userID, cipherText)
	return args.String(0), args.Error(1)
}

type MockLBService struct{ mock.Mock }

func (m *MockLBService) Create(ctx context.Context, name string, vpcID uuid.UUID, port int, algo string, idempotencyKey string) (*domain.LoadBalancer, error) {
	args := m.Called(ctx, name, vpcID, port, algo, idempotencyKey)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.LoadBalancer), args.Error(1)
}
func (m *MockLBService) Get(ctx context.Context, id uuid.UUID) (*domain.LoadBalancer, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.LoadBalancer), args.Error(1)
}
func (m *MockLBService) List(ctx context.Context) ([]*domain.LoadBalancer, error) { return nil, nil }
func (m *MockLBService) Delete(ctx context.Context, id uuid.UUID) error           { return nil }
func (m *MockLBService) AddTarget(ctx context.Context, lbID, instanceID uuid.UUID, port int, weight int) error {
	args := m.Called(ctx, lbID, instanceID, port, weight)
	return args.Error(0)
}
func (m *MockLBService) RemoveTarget(ctx context.Context, lbID, instanceID uuid.UUID) error {
	return nil
}
func (m *MockLBService) ListTargets(ctx context.Context, lbID uuid.UUID) ([]*domain.LBTarget, error) {
	return nil, nil
}

type MockStorageService struct{ mock.Mock }

func (m *MockStorageService) Upload(ctx context.Context, bucket, key string, r io.Reader) (*domain.Object, error) {
	return nil, nil
}
func (m *MockStorageService) Download(ctx context.Context, bucket, key string) (io.ReadCloser, *domain.Object, error) {
	return nil, nil, nil
}
func (m *MockStorageService) DownloadVersion(ctx context.Context, bucket, key, versionID string) (io.ReadCloser, *domain.Object, error) {
	return nil, nil, nil
}
func (m *MockStorageService) ListObjects(ctx context.Context, bucket string) ([]*domain.Object, error) {
	return nil, nil
}
func (m *MockStorageService) ListVersions(ctx context.Context, bucket, key string) ([]*domain.Object, error) {
	return nil, nil
}
func (m *MockStorageService) DeleteVersion(ctx context.Context, bucket, key, versionID string) error {
	return nil
}
func (m *MockStorageService) DeleteObject(ctx context.Context, bucket, key string) error { return nil }
func (m *MockStorageService) CreateBucket(ctx context.Context, name string, isPublic bool) (*domain.Bucket, error) {
	return nil, nil
}
func (m *MockStorageService) GetBucket(ctx context.Context, name string) (*domain.Bucket, error) {
	return nil, nil
}
func (m *MockStorageService) DeleteBucket(ctx context.Context, name string) error { return nil }
func (m *MockStorageService) SetBucketVersioning(ctx context.Context, name string, enabled bool) error {
	return nil
}
func (m *MockStorageService) ListBuckets(ctx context.Context) ([]*domain.Bucket, error) {
	return nil, nil
}
func (m *MockStorageService) GetClusterStatus(ctx context.Context) (*domain.StorageCluster, error) {
	return nil, nil
}
func (m *MockStorageService) CreateMultipartUpload(ctx context.Context, bucket, key string) (*domain.MultipartUpload, error) {
	return nil, nil
}
func (m *MockStorageService) UploadPart(ctx context.Context, uploadID uuid.UUID, partNumber int, r io.Reader) (*domain.Part, error) {
	return nil, nil
}
func (m *MockStorageService) CompleteMultipartUpload(ctx context.Context, uploadID uuid.UUID) (*domain.Object, error) {
	return nil, nil
}
func (m *MockStorageService) AbortMultipartUpload(ctx context.Context, uploadID uuid.UUID) error {
	return nil
}
func (m *MockStorageService) GeneratePresignedURL(ctx context.Context, bucket, key, method string, expiry time.Duration) (*domain.PresignedURL, error) {
	return nil, nil
}
