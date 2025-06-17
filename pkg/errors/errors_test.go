package errors

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewApplicationError(t *testing.T) {
	err := NewApplicationError("app error", 1001)
	assert.Equal(t, "app error", err.Message)
	assert.Equal(t, 1001, err.Code)
	assert.Equal(t, "Application Error: app error (1001)", err.Error())
	assert.True(t, IsApplication(err))
}

func TestNewDomainError(t *testing.T) {
	err := NewDomainError("domain error", 2002)
	assert.Equal(t, "domain error", err.Message)
	assert.Equal(t, 2002, err.Code)
	assert.Equal(t, "Domain Error: domain error (2002)", err.Error())
	assert.True(t, IsDomain(err))
}

func TestNewGeneralError(t *testing.T) {
	err := New("general error", 3003)
	genErr := err.(GeneralError)
	assert.Equal(t, "general error", genErr.Message)
	assert.Equal(t, 3003, genErr.Code)
	assert.Equal(t, "General Error: general error (3003)", genErr.Error())
	assert.True(t, IsGeneral(err))
}

func TestNewValidationError(t *testing.T) {
	fieldErr := []error{New("field1 required", 0), New("field2 invalid", 0)}
	err := NewValidationError("validation failed", fieldErr, 4004)

	assert.Equal(t, "validation failed", err.Message)
	assert.Equal(t, 4004, err.Code)
	assert.Len(t, err.FieldErrs, 2)
	assert.Equal(t, "Validation Error", err.Error())
	assert.True(t, IsValidation(err))
}

func TestEncodeError(t *testing.T) {
	tests := []struct {
		name       string
		err        error
		wantStatus int
	}{
		{"ValidationError", NewValidationError("fail", nil, 400), http.StatusBadRequest},
		{"DomainError", NewDomainError("fail", 400), http.StatusBadRequest},
		{"ApplicationError", NewApplicationError("fail", 500), http.StatusInternalServerError},
		{"GeneralError", New("fail", 500), http.StatusInternalServerError},
		{"UnknownError", assert.AnError, http.StatusInternalServerError},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			EncodeError(recorder, tt.err)
			assert.Equal(t, "application/json", recorder.Header().Get("Content-Type"))
			assert.Equal(t, tt.wantStatus, recorder.Code)

			var body map[string]interface{}
			_ = json.NewDecoder(bytes.NewReader(recorder.Body.Bytes())).Decode(&body)
			assert.NotEmpty(t, body["message"])
		})
	}
}
