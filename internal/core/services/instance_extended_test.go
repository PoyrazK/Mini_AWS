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
// TEST CONSTANTS
// =============================================================================

const (
	testImageAlpine        = "alpine:latest"
	testInstanceNamePrefix = "test-instance"
)

// =============================================================================
// PORT VALIDATION TESTS
// =============================================================================
//
// Business Rules (from domain/instance.go):
// - Ports must be in range [0, 65535] (0 means auto-assign)
// - Format: "host:container" with colon separator
// - parsePort() trims whitespace before parsing

func TestLaunchInstance_PortValidation(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name        string
		ports       string
		expectError bool
		errorMsg    string
	}{
		// Valid cases
		{"single_port_mapping", "80:80", false, ""},
		{"multiple_port_mappings", "80:80,443:443,8080:8080", false, ""},
		{"high_port_numbers", "65535:65535", false, ""},
		{"auto_assign_host_port", "0:80", false, ""},
		{"different_host_container_ports", "8080:80", false, ""},
		{"empty_ports_string", "", false, ""},
		// Note: "80 : 80" is VALID because parsePort() calls strings.TrimSpace

		// Invalid cases
		{"missing_colon_separator", "8080", true, "invalid port format"},
		{"missing_host_port", ":80", true, "invalid"},
		{"missing_container_port", "80:", true, "invalid"},
		{"port_exceeds_maximum", "70000:80", true, "out of range"},
		{"non_numeric_host_port", "abc:80", true, "invalid"},
		{"non_numeric_container_port", "80:xyz", true, "invalid"},
		{"negative_port_number", "-1:80", true, "invalid"},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			repo, _, _, _, _, _, _, _, _, svc := setupInstanceServiceTest(t)

			if !tc.expectError {
				repo.On("Create", mock.Anything, mock.Anything).Return(nil)
			}

			_, err := svc.LaunchInstance(
				context.Background(),
				testInstanceNamePrefix,
				testImageAlpine,
				tc.ports,
				nil, nil, nil,
			)

			if tc.expectError {
				require.Error(t, err, "expected error for ports: %q", tc.ports)
				if tc.errorMsg != "" {
					assert.Contains(t, err.Error(), tc.errorMsg)
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

func TestLaunchInstance_RepositoryFailure(t *testing.T) {
	t.Parallel()

	t.Run("database_connection_failure", func(t *testing.T) {
		t.Parallel()

		repo, _, _, _, _, _, _, _, _, svc := setupInstanceServiceTest(t)
		repo.On("Create", mock.Anything, mock.Anything).Return(
			errors.New(errors.Internal, "connection refused: database unavailable"),
		)

		instance, err := svc.LaunchInstance(
			context.Background(),
			testInstanceNamePrefix,
			testImageAlpine,
			"", nil, nil, nil,
		)

		require.Error(t, err)
		assert.Nil(t, instance)
		assert.Contains(t, err.Error(), "database unavailable")
		repo.AssertExpectations(t)
	})

	t.Run("constraint_violation", func(t *testing.T) {
		t.Parallel()

		repo, _, _, _, _, _, _, _, _, svc := setupInstanceServiceTest(t)
		repo.On("Create", mock.Anything, mock.Anything).Return(
			errors.New(errors.AlreadyExists, "instance name already exists"),
		)

		_, err := svc.LaunchInstance(
			context.Background(),
			testInstanceNamePrefix,
			testImageAlpine,
			"", nil, nil, nil,
		)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "already exists")
	})
}

// =============================================================================
// PROVISION FLOW ERROR HANDLING TESTS
// =============================================================================

func TestProvision_NetworkFailure_SetsStatusError(t *testing.T) {
	t.Parallel()

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

	repo.On("GetByID", mock.Anything, instanceID).Return(instance, nil)
	vpcRepo.On("GetByID", mock.Anything, vpcID).Return(&domain.VPC{
		ID:        vpcID,
		NetworkID: "vpc-network-1",
	}, nil)
	subnetRepo.On("GetByID", mock.Anything, subnetID).Return(
		nil,
		errors.New(errors.NotFound, "subnet not found"),
	)
	repo.On("Update", mock.Anything, mock.MatchedBy(func(i *domain.Instance) bool {
		return i.ID == instanceID && i.Status == domain.StatusError
	})).Return(nil)

	provisioner := assertProvisionInterface(t, svc)
	err := provisioner.Provision(ctx, instanceID, nil)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "subnet")
	repo.AssertExpectations(t)
}

func TestProvision_VolumeResolutionFailure_SetsStatusError(t *testing.T) {
	t.Parallel()

	repo, _, _, volumeRepo, _, _, _, _, _, svc := setupInstanceServiceTest(t)

	ctx := context.Background()
	instanceID := uuid.New()

	instance := &domain.Instance{
		ID:     instanceID,
		Name:   fmt.Sprintf("%s-volume-failure", testInstanceNamePrefix),
		Image:  testImageAlpine,
		Status: domain.StatusStarting,
	}

	attachments := []domain.VolumeAttachment{
		{VolumeIDOrName: "missing-volume", MountPath: "/data"},
	}

	repo.On("GetByID", mock.Anything, instanceID).Return(instance, nil)
	volumeRepo.On("GetByName", mock.Anything, "missing-volume").Return(
		nil,
		errors.New(errors.NotFound, "volume not found"),
	)
	repo.On("Update", mock.Anything, mock.MatchedBy(func(i *domain.Instance) bool {
		return i.ID == instanceID && i.Status == domain.StatusError
	})).Return(nil)

	provisioner := assertProvisionInterface(t, svc)
	err := provisioner.Provision(ctx, instanceID, attachments)

	require.Error(t, err)
	repo.AssertExpectations(t)
}

func TestProvision_InstanceNotFound_ReturnsError(t *testing.T) {
	t.Parallel()

	repo, _, _, _, _, _, _, _, _, svc := setupInstanceServiceTest(t)

	ctx := context.Background()
	nonExistentID := uuid.New()

	repo.On("GetByID", mock.Anything, nonExistentID).Return(
		nil,
		errors.New(errors.NotFound, "instance not found"),
	)

	provisioner := assertProvisionInterface(t, svc)
	err := provisioner.Provision(ctx, nonExistentID, nil)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
	repo.AssertExpectations(t)
}

// =============================================================================
// HELPER FUNCTIONS
// =============================================================================

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
