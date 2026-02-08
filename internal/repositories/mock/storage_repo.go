package mock

import (
	"context"

	"github.com/google/uuid"
	"github.com/poyrazk/thecloud/internal/core/domain"
	"github.com/poyrazk/thecloud/internal/core/ports"
	"github.com/poyrazk/thecloud/internal/errors"
)

type MockStorageRepository struct {
	Buckets map[string]*domain.Bucket
	Objects map[string]map[string][]*domain.Object // bucket -> key -> versions
	Uploads map[uuid.UUID]*domain.MultipartUpload
	Parts   map[uuid.UUID][]*domain.Part
}

func NewMockStorageRepository() *MockStorageRepository {
	return &MockStorageRepository{
		Buckets: make(map[string]*domain.Bucket),
		Objects: make(map[string]map[string][]*domain.Object),
		Uploads: make(map[uuid.UUID]*domain.MultipartUpload),
		Parts:   make(map[uuid.UUID][]*domain.Part),
	}
}

func (m *MockStorageRepository) SaveMeta(ctx context.Context, obj *domain.Object) error {
	if m.Objects[obj.Bucket] == nil {
		m.Objects[obj.Bucket] = make(map[string][]*domain.Object)
	}
	m.Objects[obj.Bucket][obj.Key] = append(m.Objects[obj.Bucket][obj.Key], obj)
	return nil
}

func (m *MockStorageRepository) GetMeta(ctx context.Context, bucket, key string) (*domain.Object, error) {
	if b, ok := m.Objects[bucket]; ok {
		if versions, ok := b[key]; ok && len(versions) > 0 {
			// return latest?
			return versions[len(versions)-1], nil
		}
	}
	return nil, errors.New(errors.NotFound, "object not found")
}

func (m *MockStorageRepository) List(ctx context.Context, bucket string) ([]*domain.Object, error) {
	var list []*domain.Object
	if b, ok := m.Objects[bucket]; ok {
		for _, versions := range b {
			if len(versions) > 0 {
				list = append(list, versions[len(versions)-1])
			}
		}
	}
	return list, nil
}

func (m *MockStorageRepository) SoftDelete(ctx context.Context, bucket, key string) error {
	return nil
}

func (m *MockStorageRepository) DeleteVersion(ctx context.Context, bucket, key, ver string) error {
	if b, ok := m.Objects[bucket]; ok {
		if versions, ok := b[key]; ok {
			var newVersions []*domain.Object
			for _, v := range versions {
				if v.VersionID != ver {
					newVersions = append(newVersions, v)
				}
			}
			m.Objects[bucket][key] = newVersions
		}
	}
	return nil
}

func (m *MockStorageRepository) GetMetaByVersion(ctx context.Context, bucket, key, ver string) (*domain.Object, error) {
	if b, ok := m.Objects[bucket]; ok {
		if versions, ok := b[key]; ok {
			for _, v := range versions {
				if v.VersionID == ver {
					return v, nil
				}
			}
		}
	}
	return nil, errors.New(errors.NotFound, "version not found")
}

func (m *MockStorageRepository) ListVersions(ctx context.Context, bucket, key string) ([]*domain.Object, error) {
	if b, ok := m.Objects[bucket]; ok {
		return b[key], nil
	}
	return nil, nil
}

func (m *MockStorageRepository) CreateBucket(ctx context.Context, b *domain.Bucket) error {
	m.Buckets[b.Name] = b
	return nil
}

func (m *MockStorageRepository) GetBucket(ctx context.Context, name string) (*domain.Bucket, error) {
	if b, ok := m.Buckets[name]; ok {
		return b, nil
	}
	return nil, errors.New(errors.NotFound, "bucket not found")
}

func (m *MockStorageRepository) DeleteBucket(ctx context.Context, name string) error {
	delete(m.Buckets, name)
	return nil
}

func (m *MockStorageRepository) ListBuckets(ctx context.Context, uid string) ([]*domain.Bucket, error) {
	var list []*domain.Bucket
	for _, b := range m.Buckets {
		if b.UserID.String() == uid {
			list = append(list, b)
		}
	}
	return list, nil
}

func (m *MockStorageRepository) SetBucketVersioning(ctx context.Context, name string, enabled bool) error {
	if b, ok := m.Buckets[name]; ok {
		b.VersioningEnabled = enabled
	}
	return nil
}

func (m *MockStorageRepository) SaveMultipartUpload(ctx context.Context, u *domain.MultipartUpload) error {
	m.Uploads[u.ID] = u
	return nil
}

func (m *MockStorageRepository) GetMultipartUpload(ctx context.Context, id uuid.UUID) (*domain.MultipartUpload, error) {
	if u, ok := m.Uploads[id]; ok {
		return u, nil
	}
	return nil, errors.New(errors.NotFound, "upload not found")
}

func (m *MockStorageRepository) DeleteMultipartUpload(ctx context.Context, id uuid.UUID) error {
	delete(m.Uploads, id)
	delete(m.Parts, id)
	return nil
}

func (m *MockStorageRepository) SavePart(ctx context.Context, p *domain.Part) error {
	m.Parts[p.UploadID] = append(m.Parts[p.UploadID], p)
	return nil
}

func (m *MockStorageRepository) ListParts(ctx context.Context, uid uuid.UUID) ([]*domain.Part, error) {
	return m.Parts[uid], nil
}

var _ ports.StorageRepository = (*MockStorageRepository)(nil)
