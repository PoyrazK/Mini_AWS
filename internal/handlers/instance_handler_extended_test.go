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
)

// =============================================================================
// INPUT VALIDATION BOUNDARY TESTS
// =============================================================================
//
// These tests verify handler input validation at boundary conditions.
// Validation rules (from instance_handler.go):
// - Name: 1-64 chars, alphanumeric + hyphen + underscore only
// - Image: max 256 chars, non-empty after trim
// - Mount path: must start with "/" (absolute path)
// - Volume ID: non-empty after trim

func TestLaunchHandler_InputValidation_Boundaries(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name           string
		requestBody    string
		expectedStatus int
		expectedError  string
		needsMock      bool // true if request should reach service layer
	}{
		// =================================================================
		// NAME FIELD VALIDATION (1-64 chars, alphanumeric/hyphen/underscore)
		// =================================================================
		{
			name:           "name_exactly_64_chars_valid",
			requestBody:    fmt.Sprintf(`{"name": "%s", "image": "alpine"}`, strings.Repeat("a", 64)),
			expectedStatus: http.StatusAccepted,
			needsMock:      true,
		},
		{
			name:           "name_65_chars_invalid",
			requestBody:    fmt.Sprintf(`{"name": "%s", "image": "alpine"}`, strings.Repeat("a", 65)),
			expectedStatus: http.StatusBadRequest,
			expectedError:  "name must be between",
		},
		{
			name:           "name_empty_string",
			requestBody:    `{"name": "", "image": "alpine"}`,
			expectedStatus: http.StatusBadRequest,
			expectedError:  "name",
		},
		{
			name:           "name_whitespace_only",
			requestBody:    `{"name": "   ", "image": "alpine"}`,
			expectedStatus: http.StatusBadRequest,
			expectedError:  "name",
		},
		{
			name:           "name_with_special_characters",
			requestBody:    `{"name": "test$name", "image": "alpine"}`,
			expectedStatus: http.StatusBadRequest,
			expectedError:  "alphanumeric",
		},

		// =================================================================
		// IMAGE FIELD VALIDATION (max 256 chars, non-empty)
		// =================================================================
		{
			name:           "image_exactly_256_chars_valid",
			requestBody:    fmt.Sprintf(`{"name": "valid-name", "image": "%s"}`, strings.Repeat("i", 256)),
			expectedStatus: http.StatusAccepted,
			needsMock:      true,
		},
		{
			name:           "image_257_chars_invalid",
			requestBody:    fmt.Sprintf(`{"name": "valid-name", "image": "%s"}`, strings.Repeat("i", 257)),
			expectedStatus: http.StatusBadRequest,
			expectedError:  "image name too long",
		},

		// =================================================================
		// VOLUME ATTACHMENT VALIDATION
		// =================================================================
		{
			name:           "volume_id_whitespace_only",
			requestBody:    `{"name": "valid-name", "image": "alpine", "volumes": [{"volume_id": "  ", "mount_path": "/data"}]}`,
			expectedStatus: http.StatusBadRequest,
			expectedError:  "volume_id is required",
		},
		{
			name:           "mount_path_empty",
			requestBody:    `{"name": "valid-name", "image": "alpine", "volumes": [{"volume_id": "vol-1", "mount_path": ""}]}`,
			expectedStatus: http.StatusBadRequest,
			expectedError:  "mount_path is required",
		},
		{
			name:           "mount_path_relative_invalid",
			requestBody:    `{"name": "valid-name", "image": "alpine", "volumes": [{"volume_id": "vol-1", "mount_path": "data/files"}]}`,
			expectedStatus: http.StatusBadRequest,
			expectedError:  "absolute path",
		},
		{
			name:           "mount_path_absolute_valid",
			requestBody:    `{"name": "valid-name", "image": "alpine", "volumes": [{"volume_id": "vol-1", "mount_path": "/mnt/data"}]}`,
			expectedStatus: http.StatusAccepted,
			needsMock:      true,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			mockSvc, handler, router := setupInstanceHandlerTest(t)
			router.POST(instancesPath, handler.Launch)

			if tc.needsMock {
				mockSvc.On("LaunchInstance",
					mock.Anything, mock.Anything, mock.Anything,
					mock.Anything, mock.Anything, mock.Anything, mock.Anything,
				).Return(&domain.Instance{ID: uuid.New(), Name: "test"}, nil)
			}

			req := httptest.NewRequest(http.MethodPost, instancesPath, strings.NewReader(tc.requestBody))
			req.Header.Set(contentType, applicationJSON)
			recorder := httptest.NewRecorder()

			router.ServeHTTP(recorder, req)

			assert.Equal(t, tc.expectedStatus, recorder.Code, "Test: %s", tc.name)

			if tc.expectedError != "" {
				assert.Contains(t, recorder.Body.String(), tc.expectedError)
			}

			if tc.expectedStatus == http.StatusBadRequest {
				mockSvc.AssertNotCalled(t, "LaunchInstance",
					mock.Anything, mock.Anything, mock.Anything,
					mock.Anything, mock.Anything, mock.Anything, mock.Anything)
			}
		})
	}
}

// =============================================================================
// CHARACTER VALIDATION TESTS
// =============================================================================
//
// Tests for name character validation (alphanumeric, hyphen, underscore only).
// Note: These tests verify the isValidResourceName() function behavior.

func TestLaunchHandler_NameCharacterValidation(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name         string
		instanceName string
		shouldReject bool
		description  string
	}{
		{
			name:         "valid_alphanumeric_with_hyphen",
			instanceName: "my-instance-123",
			shouldReject: false,
			description:  "Hyphens are allowed in names",
		},
		{
			name:         "valid_alphanumeric_with_underscore",
			instanceName: "my_instance_123",
			shouldReject: false,
			description:  "Underscores are allowed in names",
		},
		{
			name:         "invalid_special_chars_dollar",
			instanceName: "test$name",
			shouldReject: true,
			description:  "Dollar sign is not allowed",
		},
		{
			name:         "invalid_special_chars_semicolon",
			instanceName: "test;name",
			shouldReject: true,
			description:  "Semicolon is not allowed",
		},
		{
			name:         "invalid_special_chars_quotes",
			instanceName: "test'name",
			shouldReject: true,
			description:  "Quotes are not allowed",
		},
		{
			name:         "invalid_spaces",
			instanceName: "test name",
			shouldReject: true,
			description:  "Spaces are not allowed in names",
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			mockSvc, handler, router := setupInstanceHandlerTest(t)
			router.POST(instancesPath, handler.Launch)

			if !tc.shouldReject {
				mockSvc.On("LaunchInstance",
					mock.Anything, mock.Anything, mock.Anything,
					mock.Anything, mock.Anything, mock.Anything, mock.Anything,
				).Return(&domain.Instance{ID: uuid.New()}, nil)
			}

			body := fmt.Sprintf(`{"name": "%s", "image": "alpine"}`, tc.instanceName)
			req := httptest.NewRequest(http.MethodPost, instancesPath, strings.NewReader(body))
			req.Header.Set(contentType, applicationJSON)
			recorder := httptest.NewRecorder()

			router.ServeHTTP(recorder, req)

			if tc.shouldReject {
				assert.Equal(t, http.StatusBadRequest, recorder.Code,
					"Expected rejection for: %s (%s)", tc.instanceName, tc.description)
				mockSvc.AssertNotCalled(t, "LaunchInstance",
					mock.Anything, mock.Anything, mock.Anything,
					mock.Anything, mock.Anything, mock.Anything, mock.Anything)
			} else {
				assert.Equal(t, http.StatusAccepted, recorder.Code,
					"Expected acceptance for: %s (%s)", tc.instanceName, tc.description)
			}
		})
	}
}
