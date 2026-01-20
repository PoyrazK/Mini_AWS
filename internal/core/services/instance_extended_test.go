package services_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/poyrazk/thecloud/internal/core/domain"
	"github.com/poyrazk/thecloud/internal/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// =============================================================================
// TEST FIXTURES & CONSTANTS
// =============================================================================
//
// Following Go testing best practices, we define test fixtures as constants
// and factory functions to ensure:
// 1. Single source of truth for test data
// 2. Immutability of test constants
// 3. Clear separation between test setup and assertions
//
// Reference: https://dave.cheney.net/2019/05/07/prefer-table-driven-tests

const (
	// testImageAlpine is the default lightweight image for tests.
	// Using a real image name improves test readability and debugging.
	testImageAlpine = "alpine:latest"

	// testInstanceNamePrefix provides consistent naming for test instances.
	testInstanceNamePrefix = "test-instance"
)

// portValidationTestCase defines a single test case for port validation.
// Using a dedicated struct improves IDE support and refactoring safety.
type portValidationTestCase struct {
	name        string
	ports       string
	expectError bool
	errorMsg    string // Optional: expected error message substring
}

// provisionErrorTestCase defines test cases for provision failure scenarios.
type provisionErrorTestCase struct {
	name          string
	setupMocks    func(*MockInstanceRepo, *MockVpcRepo, *MockSubnetRepo, *MockVolumeRepo, uuid.UUID)
	attachments   []domain.VolumeAttachment
	expectedError string
}

// =============================================================================
// PORT VALIDATION TESTS
// =============================================================================
//
// These tests verify the port parsing and validation logic in LaunchInstance.
// Port format: "hostPort:containerPort" (e.g., "8080:80")
//
// Business Rules:
// - Ports must be in range [0, 65535] (0 means auto-assign)
// - Format must be "host:container" with colon separator
// - Multiple ports separated by comma
// - Maximum 10 port mappings per instance
//
// Test Strategy: Table-driven tests for comprehensive edge case coverage

func TestLaunchInstance_PortValidation(t *testing.T) {
	t.Parallel() // Enable parallel execution for faster test runs

	testCases := []portValidationTestCase{
		// =====================================================================
		// VALID CASES - Should succeed
		// =====================================================================
		{
			name:        "single port mapping",
			ports:       "80:80",
			expectError: false,
		},
		{
			name:        "multiple port mappings",
			ports:       "80:80,443:443,8080:8080",
			expectError: false,
		},
		{
			name:        "high port numbers",
			ports:       "65535:65535",
			expectError: false,
		},
		{
			name:        "auto-assign host port (port 0)",
			ports:       "0:80",
			expectError: false,
		},
		{
			name:        "different host and container ports",
			ports:       "8080:80",
			expectError: false,
		},
		{
			name:        "empty ports string (no port mapping)",
			ports:       "",
			expectError: false,
		},

		// =====================================================================
		// INVALID CASES - Should fail with descriptive errors
		// =====================================================================
		{
			name:        "missing colon separator",
			ports:       "8080",
			expectError: true,
			errorMsg:    "invalid port format",
		},
		{
			name:        "missing host port",
			ports:       ":80",
			expectError: true,
			errorMsg:    "invalid",
		},
		{
			name:        "missing container port",
			ports:       "80:",
			expectError: true,
			errorMsg:    "invalid",
		},
		{
			name:        "port exceeds maximum (65536)",
			ports:       "70000:80",
			expectError: true,
			errorMsg:    "out of range",
		},
		{
			name:        "non-numeric host port",
			ports:       "abc:80",
			expectError: true,
			errorMsg:    "invalid",
		},
		{
			name:        "non-numeric container port",
			ports:       "80:xyz",
			expectError: true,
			errorMsg:    "invalid",
		},
		{
			name:        "negative port number",
			ports:       "-1:80",
			expectError: true,
			errorMsg:    "invalid",
		},
		{
			name:        "whitespace in port mapping",
			ports:       "80 : 80",
			expectError: true,
			errorMsg:    "invalid",
		},
	}

	for _, tc := range testCases {
		tc := tc // Capture range variable for parallel execution
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			// Arrange
			repo, _, _, _, _, _, _, _, _, svc := setupInstanceServiceTest(t)

			if !tc.expectError {
				repo.On("Create", mock.Anything, mock.Anything).Return(nil)
			}

			// Act
			_, err := svc.LaunchInstance(
				context.Background(),
				testInstanceNamePrefix,
				testImageAlpine,
				tc.ports,
				nil, nil, nil,
			)

			// Assert
			if tc.expectError {
				require.Error(t, err, "expected error for ports: %q", tc.ports)
				if tc.errorMsg != "" {
					assert.Contains(t, err.Error(), tc.errorMsg,
						"error message should contain %q", tc.errorMsg)
				}
			} else {
				assert.NoError(t, err, "unexpected error for ports: %q", tc.ports)
			}
		})
	}
}

// =============================================================================
// REPOSITORY FAILURE TESTS
// =============================================================================
//
// These tests verify that the service layer properly handles and propagates
// repository-level failures. This is critical for:
// 1. Debugging production issues
// 2. Ensuring proper error messages reach the client
// 3. Verifying no partial state is left after failures
//
// Pattern: Dependency Injection with Mock Repositories

func TestLaunchInstance_RepositoryFailure(t *testing.T) {
	t.Parallel()

	t.Run("database connection failure", func(t *testing.T) {
		t.Parallel()

		// Arrange
		repo, _, _, _, _, _, _, _, _, svc := setupInstanceServiceTest(t)

		dbError := errors.New(errors.Internal, "connection refused: database unavailable")
		repo.On("Create", mock.Anything, mock.Anything).Return(dbError)

		// Act
		instance, err := svc.LaunchInstance(
			context.Background(),
			testInstanceNamePrefix,
			testImageAlpine,
			"", nil, nil, nil,
		)

		// Assert
		require.Error(t, err)
		assert.Nil(t, instance, "instance should be nil on error")
		assert.Contains(t, err.Error(), "database unavailable")

		repo.AssertExpectations(t)
	})

	t.Run("constraint violation error", func(t *testing.T) {
		t.Parallel()

		// Arrange
		repo, _, _, _, _, _, _, _, _, svc := setupInstanceServiceTest(t)

		constraintError := errors.New(errors.AlreadyExists, "instance name already exists")
		repo.On("Create", mock.Anything, mock.Anything).Return(constraintError)

		// Act
		_, err := svc.LaunchInstance(
			context.Background(),
			testInstanceNamePrefix,
			testImageAlpine,
			"", nil, nil, nil,
		)

		// Assert
		require.Error(t, err)
		assert.Contains(t, err.Error(), "already exists")
	})
}

// =============================================================================
// TASK QUEUE INTEGRATION TESTS
// =============================================================================
//
// Note: TaskQueueStub is a concrete implementation that cannot be mocked
// without interface extraction. This is documented as a future refactoring
// opportunity to improve testability.
//
// TODO: Extract TaskQueue interface from TaskQueueStub for better testability
// See: internal/core/services/taskqueue.go

func TestLaunchInstance_TaskQueueFailure(t *testing.T) {
	t.Skip(`
		SKIP REASON: TaskQueueStub is a concrete struct, not an interface.
		
		IMPACT: Cannot test async task queue failure scenarios.
		
		RECOMMENDATION: Extract interface from TaskQueueStub:
		
		type TaskQueue interface {
			Enqueue(ctx context.Context, task Task) error
		}
		
		This would allow mocking queue failures in tests.
	`)
}

// =============================================================================
// PROVISION FLOW ERROR HANDLING TESTS
// =============================================================================
//
// The Provision method orchestrates multiple operations:
// 1. Fetch instance from repository
// 2. Resolve VPC/Subnet networking
// 3. Resolve volume attachments
// 4. Create compute instance
// 5. Configure networking
// 6. Update instance status
//
// These tests verify that failures at each step properly:
// - Update instance status to ERROR
// - Return descriptive error messages
// - Clean up partial state where possible
//
// Test Pattern: Arrange-Act-Assert with explicit mock verification

func TestProvision_NetworkFailure_SetsStatusError(t *testing.T) {
	t.Parallel()

	// Arrange
	repo, vpcRepo, subnetRepo, _, _, _, _, _, _, svc := setupInstanceServiceTest(t)

	ctx := context.Background()
	instanceID := uuid.New()
	vpcID := uuid.New()
	subnetID := uuid.New()

	instance := &domain.Instance{
		ID:       instanceID,
		Name:     fmt.Sprintf("%s-network-failure", testInstanceNamePrefix),
		Image:    testImageAlpine,
		VpcID:    &vpcID,
		SubnetID: &subnetID,
		Status:   domain.StatusStarting,
	}

	// Setup mock expectations in execution order
	repo.On("GetByID", mock.Anything, instanceID).Return(instance, nil)
	vpcRepo.On("GetByID", mock.Anything, vpcID).Return(&domain.VPC{
		ID:        vpcID,
		NetworkID: "vpc-network-1",
	}, nil)

	// Simulate subnet lookup failure (network resolution error)
	subnetRepo.On("GetByID", mock.Anything, subnetID).Return(
		nil,
		errors.New(errors.NotFound, "subnet not found: may have been deleted"),
	)

	// Verify status is updated to ERROR on failure
	repo.On("Update", mock.Anything, mock.MatchedBy(func(i *domain.Instance) bool {
		return i.ID == instanceID && i.Status == domain.StatusError
	})).Return(nil)

	// Act
	provisioner := assertProvisionInterface(t, svc)
	err := provisioner.Provision(ctx, instanceID, nil)

	// Assert
	require.Error(t, err)
	assert.Contains(t, err.Error(), "subnet")

	// Verify all mock expectations were met
	repo.AssertExpectations(t)
	vpcRepo.AssertExpectations(t)
	subnetRepo.AssertExpectations(t)
}

func TestProvision_VolumeResolutionFailure_SetsStatusError(t *testing.T) {
	t.Parallel()

	// Arrange
	repo, _, _, volumeRepo, _, _, _, _, _, svc := setupInstanceServiceTest(t)

	ctx := context.Background()
	instanceID := uuid.New()
	missingVolumeName := "critical-data-volume"

	instance := &domain.Instance{
		ID:     instanceID,
		Name:   fmt.Sprintf("%s-volume-failure", testInstanceNamePrefix),
		Image:  testImageAlpine,
		Status: domain.StatusStarting,
		// No VPC/Subnet - direct volume attachment test
	}

	attachments := []domain.VolumeAttachment{
		{
			VolumeIDOrName: missingVolumeName,
			MountPath:      "/data",
		},
	}

	// Setup mock expectations
	repo.On("GetByID", mock.Anything, instanceID).Return(instance, nil)

	// Volume resolution attempts name lookup when not a valid UUID
	volumeRepo.On("GetByName", mock.Anything, missingVolumeName).Return(
		nil,
		errors.New(errors.NotFound, "volume not found"),
	)

	// Status should be updated to ERROR
	repo.On("Update", mock.Anything, mock.MatchedBy(func(i *domain.Instance) bool {
		return i.ID == instanceID && i.Status == domain.StatusError
	})).Return(nil)

	// Act
	provisioner := assertProvisionInterface(t, svc)
	err := provisioner.Provision(ctx, instanceID, attachments)

	// Assert
	require.Error(t, err)

	repo.AssertExpectations(t)
	volumeRepo.AssertExpectations(t)
}

func TestProvision_InstanceNotFound_ReturnsError(t *testing.T) {
	t.Parallel()

	// Arrange
	repo, _, _, _, _, _, _, _, _, svc := setupInstanceServiceTest(t)

	ctx := context.Background()
	nonExistentID := uuid.New()

	repo.On("GetByID", mock.Anything, nonExistentID).Return(
		nil,
		errors.New(errors.NotFound, "instance not found"),
	)

	// Act
	provisioner := assertProvisionInterface(t, svc)
	err := provisioner.Provision(ctx, nonExistentID, nil)

	// Assert
	require.Error(t, err)
	assert.Contains(t, err.Error(), "not found")

	repo.AssertExpectations(t)
}

// =============================================================================
// HELPER FUNCTIONS
// =============================================================================
//
// Helper functions improve test readability and reduce duplication.
// They should be:
// 1. Focused on a single responsibility
// 2. Well-documented
// 3. Panic-free (use t.Helper() and require for fatal assertions)

// assertProvisionInterface asserts that the service implements the Provision
// method and returns it as a typed interface for testing.
//
// This pattern allows testing of methods that may not be part of the public
// interface but are critical for the service's operation.
func assertProvisionInterface(t *testing.T, svc interface{}) interface {
	Provision(context.Context, uuid.UUID, []domain.VolumeAttachment) error
} {
	t.Helper()

	provisioner, ok := svc.(interface {
		Provision(context.Context, uuid.UUID, []domain.VolumeAttachment) error
	})
	require.True(t, ok, "service must implement Provision method")

	return provisioner
}
