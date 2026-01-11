package httputil

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/poyrazk/thecloud/internal/errors"
)

type Meta struct {
	RequestID string `json:"request_id,omitempty"`
	Timestamp string `json:"timestamp"`
}

type Response struct {
	Data  interface{} `json:"data,omitempty"`
	Error interface{} `json:"error,omitempty"`
	Meta  *Meta       `json:"meta,omitempty"`
}

func Success(c *gin.Context, code int, data interface{}) {
	requestID, _ := c.Get("requestID")
	reqIDStr, _ := requestID.(string)

	c.JSON(code, Response{
		Data: data,
		Meta: &Meta{
			RequestID: reqIDStr,
			Timestamp: time.Now().UTC().Format(time.RFC3339),
		},
	})
}

func Error(c *gin.Context, err error) {
	var e errors.Error
	if apiErr, ok := err.(errors.Error); ok {
		e = apiErr
	} else {
		// Log unknown errors for debugging
		_ = c.Error(err)

		// Use GetCause to show we care about the underlying error for logs
		cause := errors.GetCause(err)
		if cause != nil {
			fmt.Printf("API ERROR CAUSE: %v\n", cause)
		} else {
			fmt.Printf("API ERROR: %v\n", err)
		}

		e = errors.Error{
			Type:    errors.Internal,
			Message: "An unexpected error occurred",
			Code:    string(errors.Internal),
			Cause:   err,
		}
	}

	statusCodeMap := map[errors.Type]int{
		errors.NotFound:              http.StatusNotFound,
		errors.InvalidInput:          http.StatusBadRequest,
		errors.Unauthorized:          http.StatusUnauthorized,
		errors.Forbidden:             http.StatusForbidden,
		errors.Conflict:              http.StatusConflict,
		errors.BucketNotFound:        http.StatusNotFound,
		errors.ObjectNotFound:        http.StatusNotFound,
		errors.ObjectTooLarge:        http.StatusRequestEntityTooLarge,
		errors.InstanceNotRunning:    http.StatusConflict,
		errors.PortConflict:          http.StatusConflict,
		errors.TooManyPorts:          http.StatusConflict,
		errors.ResourceLimitExceeded: http.StatusTooManyRequests,
	}

	statusCode := http.StatusInternalServerError
	if code, ok := statusCodeMap[e.Type]; ok {
		statusCode = code
	}

	requestID, _ := c.Get("requestID")
	reqIDStr, _ := requestID.(string)

	c.JSON(statusCode, Response{
		Error: e,
		Meta: &Meta{
			RequestID: reqIDStr,
			Timestamp: time.Now().UTC().Format(time.RFC3339),
		},
	})
}
