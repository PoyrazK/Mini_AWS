# Testing Guide

This document describes how to run tests locally and what each suite covers.

## Unit Tests (Go)
Run all unit tests:
```bash
go test ./...
```

## Integration Tests (Go)
Some packages use the `integration` build tag.
```bash
go test -tags=integration -v ./internal/repositories/postgres/...
```

## System Tests (CLI + API)
The system test runs the API and exercises the CLI against it.

Prerequisites:
- Docker running
- PostgreSQL running locally via Docker Compose

Steps:
```bash
docker compose up -d postgres
make build
DATABASE_URL=postgres://cloud:cloud@localhost:5433/thecloud ./tests/system_test.sh
```

## E2E Multi-Tenancy Test
The E2E test expects a running API at `http://localhost:8080`.

Steps:
```bash
docker compose up -d postgres
go run cmd/api/main.go
```
Then, in another terminal:
```bash
go test -tags=e2e ./tests/...
```
