package httphandlers

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/poyrazk/thecloud/internal/core/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// =============================================================================
// TEST DOCUMENTATION
// =============================================================================
//
// This file contains extended boundary case tests for the Instance Handler.
// These tests complement the main test file by focusing on edge cases and
// input validation boundaries that are critical for security and stability.
//
// Test Categories:
// 1. Input Length Validation - Prevents buffer overflows and DoS attacks
// 2. Format Validation - Ensures data integrity
// 3. Security Validation - Prevents injection and path traversal attacks
//
// Testing Philosophy:
// - Tests should be independent and isolated
// - Each test creates fresh mocks to prevent state leakage
// - Tests are parallelized for faster execution
// - Error messages should be actionable for API consumers
//
// Reference: https://google.github.io/styleguide/go/best-practices#test-structure

// =============================================================================
// TEST FIXTURES
// =============================================================================

// boundaryTestCase defines a structured test case for input validation testing.
// Using a struct provides better IDE support, refactoring safety, and documentation.
type boundaryTestCase struct {
	name           string // Descriptive name for test identification
	requestBody    string // JSON request payload
	expectedStatus int    // Expected HTTP status code
	expectedError  string // Expected error message substring
	description    string // Why this test case matters
}

// generateLongString creates a string of specified length for boundary testing.
// This is preferred over inline strings.Repeat for clarity and potential future enhancements.
func generateLongString(char string, length int) string {
	return strings.Repeat(char, length)
}

// =============================================================================
// INPUT VALIDATION BOUNDARY TESTS
// =============================================================================
//
// These tests verify that the handler properly validates input at the boundaries
// defined by the API specification. Proper boundary validation is critical for:
//
// 1. SECURITY: Prevents buffer overflow attacks and resource exhaustion
// 2. DATA INTEGRITY: Ensures stored data meets system constraints
// 3. USER EXPERIENCE: Provides clear, actionable error messages
// 4. SYSTEM STABILITY: Prevents cascading failures from invalid data
//
// Test Pattern: Table-driven tests with isolated mock setup per case

func TestLaunchHandler_InputValidation_Boundaries(t *testing.T) {
	t.Parallel()

	testCases := []boundaryTestCase{
		// =====================================================================
		// NAME FIELD VALIDATION
		// =====================================================================
		// Instance names have specific constraints:
		// - Minimum: 1 character
		// - Maximum: 64 characters
		// - Pattern: alphanumeric, hyphens, underscores only
		{
			name:           "name_at_maximum_boundary_65_chars",
			requestBody:    fmt.Sprintf(`{"name": "%s", "image": "alpine"}`, generateLongString("a", 65)),
			expectedStatus: http.StatusBadRequest,
			expectedError:  "name must be between",
			description:    "Names exceeding 64 chars should be rejected to prevent database field overflow",
		},
		{
			name:           "name_empty_string",
			requestBody:    `{"name": "", "image": "alpine"}`,
			expectedStatus: http.StatusBadRequest,
			expectedError:  "name",
			description:    "Empty names should be rejected",
		},
		{
			name:           "name_whitespace_only",
			requestBody:    `{"name": "   ", "image": "alpine"}`,
			expectedStatus: http.StatusBadRequest,
			expectedError:  "name",
			description:    "Whitespace-only names should be rejected after trimming",
		},

		// =====================================================================
		// IMAGE FIELD VALIDATION
		// =====================================================================
		// Image references follow Docker/OCI naming conventions:
		// - Maximum: 256 characters (registry/namespace/name:tag)
		// - Cannot be empty or whitespace-only
		{
			name:           "image_exceeds_maximum_257_chars",
			requestBody:    fmt.Sprintf(`{"name": "valid-name", "image": "%s"}`, generateLongString("i", 257)),
			expectedStatus: http.StatusBadRequest,
			expectedError:  "image name too long",
			description:    "Image names exceeding 256 chars violate OCI spec and should be rejected",
		},

		// =====================================================================
		// VOLUME ATTACHMENT VALIDATION
		// =====================================================================
		// Volume attachments require:
		// - Non-empty volume_id (trimmed)
		// - Absolute mount_path starting with /
		// - Valid mount path characters (no .. for security)
		{
			name:           "volume_id_whitespace_only",
			requestBody:    `{"name": "valid-name", "image": "alpine", "volumes": [{"volume_id": "  ", "mount_path": "/data"}]}`,
			expectedStatus: http.StatusBadRequest,
			expectedError:  "volume_id is required",
			description:    "Whitespace-only volume IDs should be rejected after trimming",
		},
		{
			name:           "volume_mount_path_empty",
			requestBody:    `{"name": "valid-name", "image": "alpine", "volumes": [{"volume_id": "vol-1", "mount_path": ""}]}`,
			expectedStatus: http.StatusBadRequest,
			expectedError:  "mount_path is required",
			description:    "Empty mount paths should be rejected",
		},
		{
			name:           "volume_mount_path_relative",
			requestBody:    `{"name": "valid-name", "image": "alpine", "volumes": [{"volume_id": "vol-1", "mount_path": "data/files"}]}`,
			expectedStatus: http.StatusBadRequest,
			expectedError:  "absolute path",
			description:    "Relative paths are security risks and should be rejected",
		},
	}

	for _, tc := range testCases {
		tc := tc // Capture range variable for parallel execution

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			// Arrange - Fresh setup for each test ensures complete isolation
			mockSvc, handler, router := setupInstanceHandlerTest(t)
			router.POST(instancesPath, handler.Launch)

			// Act
			req := httptest.NewRequest(http.MethodPost, instancesPath, strings.NewReader(tc.requestBody))
			req.Header.Set(contentType, applicationJSON)
			recorder := httptest.NewRecorder()

			router.ServeHTTP(recorder, req)

			// Assert - Status code
			assert.Equal(t, tc.expectedStatus, recorder.Code,
				"Test: %s\nDescription: %s\nBody: %s",
				tc.name, tc.description, tc.requestBody)

			// Assert - Error message content (if expected)
			if tc.expectedError != "" {
				assert.Contains(t, recorder.Body.String(), tc.expectedError,
					"Response should contain error message: %s", tc.expectedError)
			}

			// Assert - Service isolation (validation failures should never reach service)
			if tc.expectedStatus == http.StatusBadRequest {
				mockSvc.AssertNotCalled(t, "LaunchInstance",
					mock.Anything, mock.Anything, mock.Anything,
					mock.Anything, mock.Anything, mock.Anything, mock.Anything)
			}
		})
	}
}

// =============================================================================
// SECURITY-FOCUSED VALIDATION TESTS
// =============================================================================
//
// These tests specifically target security-sensitive validation scenarios.
// They are separated for clarity and to ensure security requirements are
// explicitly documented and tested.

func TestLaunchHandler_SecurityValidation(t *testing.T) {
	t.Parallel()

	securityTestCases := []struct {
		name         string
		requestBody  string
		shouldReject bool
		securityRisk string // Documented security risk being tested
	}{
		{
			name:         "sql_injection_attempt_in_name",
			requestBody:  `{"name": "'; DROP TABLE instances; --", "image": "alpine"}`,
			shouldReject: true,
			securityRisk: "SQL injection via instance name",
		},
		{
			name:         "path_traversal_in_mount_path",
			requestBody:  `{"name": "valid-name", "image": "alpine", "volumes": [{"volume_id": "v1", "mount_path": "/../../../etc/shadow"}]}`,
			shouldReject: true,
			securityRisk: "Path traversal to access sensitive host files",
		},
	}

	for _, tc := range securityTestCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			// Arrange
			mockSvc, handler, router := setupInstanceHandlerTest(t)
			router.POST(instancesPath, handler.Launch)

			if !tc.shouldReject {
				mockSvc.On("LaunchInstance",
					mock.Anything, mock.Anything, mock.Anything,
					mock.Anything, mock.Anything, mock.Anything, mock.Anything,
				).Return(&domain.Instance{ID: uuid.New()}, nil)
			}

			// Act
			req := httptest.NewRequest(http.MethodPost, instancesPath, strings.NewReader(tc.requestBody))
			req.Header.Set(contentType, applicationJSON)
			recorder := httptest.NewRecorder()

			router.ServeHTTP(recorder, req)

			// Assert
			if tc.shouldReject {
				require.Equal(t, http.StatusBadRequest, recorder.Code,
					"Security risk not properly handled: %s", tc.securityRisk)
				mockSvc.AssertNotCalled(t, "LaunchInstance",
					mock.Anything, mock.Anything, mock.Anything,
					mock.Anything, mock.Anything, mock.Anything, mock.Anything)
			}
		})
	}
}
