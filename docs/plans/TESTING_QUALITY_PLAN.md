# Testing & Quality Improvement Plan

## Current State Assessment

### Test Coverage Summary
| Package | Coverage | Priority |
|---------|----------|----------|
| `internal/workers` | 31.1% | 🔴 Critical |
| `internal/repositories/k8s` | 28.8% | 🔴 Critical |
| `pkg/sshutil` | 34.8% | 🔴 Critical |
| `internal/core/context` | 50.0% | 🟡 Medium |
| `pkg/httputil` | 57.1% | 🟡 Medium |
| `internal/storage/node` | 61.6% | 🟡 Medium |
| `internal/repositories/postgres` | 63.5% | 🟡 Medium |
| `internal/repositories/libvirt` | 66.6% | 🟡 Medium |
| `pkg/sdk` | 67.5% | 🟡 Medium |
| `pkg/crypto` | 71.4% | 🟢 Good |
| `internal/repositories/ovs` | 72.5% | 🟢 Good |
| `internal/repositories/docker` | 74.4% | 🟢 Good |
| `internal/handlers` | 76.1% | 🟢 Good |
| `internal/core/services` | 77.7% | 🟢 Good |
| `internal/api/setup` | 78.5% | 🟢 Good |
| `internal/repositories/lvm` | 85.0% | 🟢 Good |
| `pkg/ratelimit` | 92.6% | ✅ Excellent |
| `internal/core/domain` | 100.0% | ✅ Excellent |
| `internal/errors` | 100.0% | ✅ Excellent |
| `pkg/audit` | 100.0% | ✅ Excellent |

### E2E Tests: 18 test suites in `tests/`
### Total Test Files: 217

---

## Phase 1: Critical Coverage Gaps (Week 1-2)

### 1.1 Workers Package (31.1% → 70%)
**File**: `internal/workers/`

**Tasks**:
- [ ] Add unit tests for `autoscaling_worker.go`
- [ ] Add unit tests for `health_check_worker.go`
- [ ] Add unit tests for `metrics_worker.go`
- [ ] Add unit tests for `cron_worker.go`
- [ ] Mock external dependencies (Redis, Docker stats)

**Test Cases**:
```go
// autoscaling_worker_test.go
- TestScaleUp_WhenCPUExceedsThreshold
- TestScaleDown_WhenCPUIdle
- TestCooldownPeriod_PreventsRapidScaling
- TestWorker_GracefulShutdown
```

### 1.2 Kubernetes Repository (28.8% → 60%)
**File**: `internal/repositories/k8s/`

**Tasks**:
- [ ] Add mock Kubernetes client tests
- [ ] Test cluster provisioning logic
- [ ] Test node join/leave operations
- [ ] Test kubeconfig generation

**Test Cases**:
```go
// cluster_repo_test.go
- TestCreateCluster_Success
- TestCreateCluster_HAMode
- TestDeleteCluster_CleansUpResources
- TestGetKubeconfig_ReturnsValidConfig
```

### 1.3 SSH Utility (34.8% → 70%)
**File**: `pkg/sshutil/`

**Tasks**:
- [ ] Create mock SSH server for testing
- [ ] Test connection handling
- [ ] Test command execution
- [ ] Test file transfer operations

---

## Phase 2: Medium Priority (Week 3-4)

### 2.1 Context Package (50% → 80%)
**Tasks**:
- [ ] Test tenant context propagation
- [ ] Test user context extraction
- [ ] Test context cancellation scenarios

### 2.2 HTTP Utilities (57.1% → 80%)
**Tasks**:
- [ ] Test tenant middleware thoroughly
- [ ] Test error response formatting
- [ ] Test auth middleware edge cases
- [ ] Test rate limiting integration

### 2.3 Storage Node (61.6% → 75%)
**Tasks**:
- [ ] Test RPC handlers
- [ ] Test replication logic
- [ ] Test failure scenarios

### 2.4 Postgres Repository (63.5% → 80%)
**Tasks**:
- [ ] Add missing repository tests (tenant, rbac)
- [ ] Test transaction rollback scenarios
- [ ] Test concurrent access patterns
- [ ] Improve migration test coverage

---

## Phase 3: E2E Test Expansion (Week 5-6)

### 3.1 Multi-Tenancy E2E
**File**: `tests/multitenancy_e2e_test.go`

**New Test Cases**:
- [ ] `TestTenantIsolation_CannotAccessOtherTenantResources`
- [ ] `TestTenantQuota_BlocksExcessiveCreation`
- [ ] `TestTenantMembership_InviteAcceptFlow`
- [ ] `TestTenantDeletion_CascadesResources`
- [ ] `TestCrossTenantAPI_Returns403`

### 3.2 Security E2E
**File**: `tests/security_edge_test.go`

**New Test Cases**:
- [ ] `TestSQLInjection_AllEndpoints`
- [ ] `TestXSS_InputSanitization`
- [ ] `TestRateLimiting_BlocksBruteForce`
- [ ] `TestJWTExpiry_DeniesAccess`
- [ ] `TestAPIKeyRotation_InvalidatesOld`

### 3.3 Chaos Engineering
**File**: `tests/chaos_test.go`

**New Test Cases**:
- [ ] `TestDatabaseFailover_RecoveryTime`
- [ ] `TestNetworkPartition_GracefulDegradation`
- [ ] `TestWorkerCrash_JobRecovery`
- [ ] `TestHighLoad_NoDataLoss`

---

## Phase 4: Quality Infrastructure (Week 7-8)

### 4.1 CI/CD Enhancements
**Tasks**:
- [ ] Add coverage threshold check (fail if < 70%)
- [ ] Add benchmark regression tests
- [ ] Enable parallel test execution
- [ ] Add test result caching

**GitHub Actions Update**:
```yaml
# .github/workflows/test.yml additions
- name: Check Coverage Threshold
  run: |
    COVERAGE=$(go test -cover ./... | grep total | awk '{print $3}' | tr -d '%')
    if [ "$COVERAGE" -lt "70" ]; then
      echo "Coverage $COVERAGE% is below 70% threshold"
      exit 1
    fi

- name: Run Benchmarks
  run: go test -bench=. -benchmem ./... > bench.txt

- name: Compare Benchmarks
  uses: benchmark-action/github-action-benchmark@v1
```

### 4.2 Test Utilities
**File**: `pkg/testutil/`

**Create**:
- [ ] `testutil/fixtures.go` - Common test data generators
- [ ] `testutil/mocks.go` - Shared mock implementations
- [ ] `testutil/assertions.go` - Custom assertions
- [ ] `testutil/db.go` - Test database helpers

### 4.3 Load Testing
**Directory**: `tests/load/`

**Tasks**:
- [ ] Create k6 scripts for API endpoints
- [ ] Define performance baselines
- [ ] Add load test to CI (nightly)

**k6 Script Example**:
```javascript
// tests/load/instances.js
import http from 'k6/http';
import { check, sleep } from 'k6';

export const options = {
  stages: [
    { duration: '1m', target: 50 },
    { duration: '3m', target: 100 },
    { duration: '1m', target: 0 },
  ],
  thresholds: {
    http_req_duration: ['p(95)<500'],
    http_req_failed: ['rate<0.01'],
  },
};

export default function () {
  const res = http.get('http://localhost:8080/api/v1/instances');
  check(res, { 'status is 200': (r) => r.status === 200 });
  sleep(1);
}
```

---

## Phase 5: Documentation & Standards (Ongoing)

### 5.1 Test Documentation
- [ ] Add test README explaining test categories
- [ ] Document how to run specific test suites
- [ ] Create test naming conventions guide

### 5.2 Code Quality Tools
- [ ] Enable `golangci-lint` with strict config
- [ ] Add `go vet` to CI
- [ ] Enable race detector in CI tests
- [ ] Add deadcode detection

**golangci-lint config** (`.golangci.yml`):
```yaml
linters:
  enable:
    - errcheck
    - gosimple
    - govet
    - ineffassign
    - staticcheck
    - unused
    - bodyclose
    - gocritic
    - gocyclo
    - gosec

linters-settings:
  gocyclo:
    min-complexity: 15
  gocritic:
    enabled-tags:
      - diagnostic
      - performance
```

---

## Success Metrics

| Metric | Current | Target | Timeline |
|--------|---------|--------|----------|
| Overall Coverage | ~65% | 80% | 8 weeks |
| Critical Packages | 30-35% | 70% | 2 weeks |
| E2E Test Count | 18 | 30 | 6 weeks |
| CI Pipeline Time | ~3min | <5min | 4 weeks |
| Flaky Test Rate | Unknown | <1% | Ongoing |

---

## Quick Wins (Start Today)

1. **Add coverage badge** to README
2. **Fix the auth.go TODO** - implement user rollback
3. **Add test for tenant middleware** in `pkg/httputil`
4. **Create shared test fixtures** for tenant/user creation

---

## Commands Reference

```bash
# Run all tests with coverage
make test

# Run specific package tests
go test -v -cover ./internal/workers/...

# Run with race detector
go test -race ./...

# Generate coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html

# Run benchmarks
go test -bench=. -benchmem ./internal/core/services/...

# Run integration tests only
go test -tags=integration -v ./internal/repositories/postgres/...

# Run E2E tests
go test -v ./tests/...
```
