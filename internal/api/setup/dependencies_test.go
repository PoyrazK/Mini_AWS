package setup

import (
	"context"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/poyrazk/thecloud/internal/repositories/postgres"
	redisv9 "github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockDB
type MockDB struct {
	mock.Mock
}

func (m *MockDB) Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error) {
	args := m.Called(ctx, sql, arguments)
	return args.Get(0).(pgconn.CommandTag), args.Error(1)
}
func (m *MockDB) Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error) {
	a := m.Called(ctx, sql, args)
	return a.Get(0).(pgx.Rows), a.Error(1)
}
func (m *MockDB) QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row {
	a := m.Called(ctx, sql, args)
	return a.Get(0).(pgx.Row)
}
func (m *MockDB) Begin(ctx context.Context) (pgx.Tx, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(pgx.Tx), args.Error(1)
}
func (m *MockDB) Close() {
	m.Called()
}
func (m *MockDB) Ping(ctx context.Context) error {
	return m.Called(ctx).Error(0)
}

func TestInitRepositories(t *testing.T) {
	mockDB := new(MockDB)
	// mockRed := redismock.NewClientMock() // Or just pass nil if not used in this specific test

	// The InitRepositories function uses the DB to create struct instances.
	// It doesn't actually call DB methods during initialization, so verification is trivial (no panic).
	// We can assert fields are not nil.

	repos := InitRepositories(mockDB, &redisv9.Client{})

	assert.NotNil(t, repos)
	assert.NotNil(t, repos.User)
	assert.NotNil(t, repos.Instance)
	assert.NotNil(t, repos.Vpc)
	assert.NotNil(t, repos.DNS)
	assert.IsType(t, &postgres.UserRepo{}, repos.User)
	assert.IsType(t, &postgres.InstanceRepository{}, repos.Instance)
}

type mockContext struct{}

func (m *mockContext) Deadline() (deadline time.Time, ok bool) { return }
func (m *mockContext) Done() <-chan struct{}                   { return nil }
func (m *mockContext) Err() error                              { return nil }
func (m *mockContext) Value(key interface{}) interface{}       { return nil }

func TestInitServices(t *testing.T) {
	// mockDB := new(MockDB)
	// c := ServiceConfig {
	// 	DB: mockDB,
	// 	Repos: &Repositories{},
	// 	Logger: slog.New(slog.NewTextHandler(io.Discard, nil)),
	// }
	// // This is harder because InitServices creates a lot of things and might fail
}
