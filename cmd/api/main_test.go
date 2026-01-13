package main

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"os"
	"syscall"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/poyrazk/thecloud/internal/api/setup"
	"github.com/poyrazk/thecloud/internal/platform"
	"github.com/poyrazk/thecloud/internal/repositories/postgres"
	"github.com/redis/go-redis/v9"
)

func TestInitInfrastructure_MigrateOnlyStopsAfterMigrations(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	fakeDB := &stubDB{}

	resetLoadConfig := stubLoadConfig(func(*slog.Logger) (*platform.Config, error) {
		return &platform.Config{}, nil
	})
	defer resetLoadConfig()

	resetInitDatabase := stubInitDatabase(func(context.Context, *platform.Config, *slog.Logger) (postgres.DB, error) {
		return fakeDB, nil
	})
	defer resetInitDatabase()

	migrationsRan := false
	resetRunMigrations := stubRunMigrations(func(context.Context, postgres.DB, *slog.Logger) error {
		migrationsRan = true
		return nil
	})
	defer resetRunMigrations()

	resetInitRedis := stubInitRedis(func(context.Context, *platform.Config, *slog.Logger) (*redis.Client, error) {
		t.Fatalf("initRedis should not be called when migrate-only completes")
		return nil, nil
	})
	defer resetInitRedis()

	cfg, db, rdb, err := initInfrastructure(logger, true)

	if !errors.Is(err, ErrMigrationDone) {
		t.Fatalf("expected ErrMigrationDone, got %v", err)
	}
	if cfg != nil || db != nil || rdb != nil {
		t.Fatalf("expected nil resources after migrate-only run, got cfg=%v db=%v rdb=%v", cfg, db, rdb)
	}
	if !migrationsRan {
		t.Fatalf("expected migrations to run")
	}
	if !fakeDB.closed {
		t.Fatalf("expected database to be closed")
	}
}

func TestInitInfrastructure_ConfigError(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	resetLoadConfig := stubLoadConfig(func(*slog.Logger) (*platform.Config, error) {
		return nil, errors.New("boom")
	})
	defer resetLoadConfig()

	if _, _, _, err := initInfrastructure(logger, false); err == nil {
		t.Fatalf("expected error when config loading fails")
	}
}

func TestRunApplication_ApiRoleStartsAndShutsDown(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	t.Setenv("ROLE", "api")

	started := make(chan struct{})
	shutdownCalled := make(chan struct{})

	resetNewHTTPServer := stubNewHTTPServer(func(addr string, handler http.Handler) *http.Server {
		return &http.Server{Addr: addr, Handler: handler}
	})
	defer resetNewHTTPServer()

	resetStartServer := stubStartHTTPServer(func(*http.Server) error {
		close(started)
		return http.ErrServerClosed
	})
	defer resetStartServer()

	resetShutdown := stubShutdownHTTPServer(func(context.Context, *http.Server) error {
		close(shutdownCalled)
		return nil
	})
	defer resetShutdown()

	resetNotify := stubNotifySignals(func(c chan<- os.Signal, _ ...os.Signal) {
		go func() {
			<-started
			c <- syscall.SIGTERM
		}()
	})
	defer resetNotify()

	runApplication(&platform.Config{Port: "0"}, logger, gin.New(), &setup.Workers{})

	select {
	case <-shutdownCalled:
	case <-time.After(time.Second):
		t.Fatalf("expected server shutdown to be called")
	}
}

// Stub helpers below keep main.go testable without altering production behavior.

type stubDB struct{ closed bool }

func (s *stubDB) Exec(context.Context, string, ...interface{}) (pgconn.CommandTag, error) {
	return pgconn.CommandTag{}, nil
}
func (s *stubDB) Query(context.Context, string, ...interface{}) (pgx.Rows, error) { return nil, nil }
func (s *stubDB) QueryRow(context.Context, string, ...interface{}) pgx.Row        { return nil }
func (s *stubDB) Begin(context.Context) (pgx.Tx, error)                           { return nil, nil }
func (s *stubDB) Close()                                                          { s.closed = true }
func (s *stubDB) Ping(context.Context) error                                      { return nil }

func stubLoadConfig(fn func(*slog.Logger) (*platform.Config, error)) func() {
	prev := loadConfigFunc
	loadConfigFunc = fn
	return func() { loadConfigFunc = prev }
}

func stubInitDatabase(fn func(context.Context, *platform.Config, *slog.Logger) (postgres.DB, error)) func() {
	prev := initDatabaseFunc
	initDatabaseFunc = fn
	return func() { initDatabaseFunc = prev }
}

func stubRunMigrations(fn func(context.Context, postgres.DB, *slog.Logger) error) func() {
	prev := runMigrationsFunc
	runMigrationsFunc = fn
	return func() { runMigrationsFunc = prev }
}

func stubInitRedis(fn func(context.Context, *platform.Config, *slog.Logger) (*redis.Client, error)) func() {
	prev := initRedisFunc
	initRedisFunc = fn
	return func() { initRedisFunc = prev }
}

func stubNewHTTPServer(fn func(string, http.Handler) *http.Server) func() {
	prev := newHTTPServer
	newHTTPServer = fn
	return func() { newHTTPServer = prev }
}

func stubStartHTTPServer(fn func(*http.Server) error) func() {
	prev := startHTTPServer
	startHTTPServer = fn
	return func() { startHTTPServer = prev }
}

func stubShutdownHTTPServer(fn func(context.Context, *http.Server) error) func() {
	prev := shutdownHTTPServer
	shutdownHTTPServer = fn
	return func() { shutdownHTTPServer = prev }
}

func stubNotifySignals(fn func(chan<- os.Signal, ...os.Signal)) func() {
	prev := notifySignals
	notifySignals = fn
	return func() { notifySignals = prev }
}
