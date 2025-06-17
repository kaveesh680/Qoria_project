package errors

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type ApplicationError struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

func (e ApplicationError) Error() string {
	return fmt.Sprintf("Application Error: %s (%d)", e.Message, e.Code)
}

func NewApplicationError(message string, code int) ApplicationError {
	return ApplicationError{
		Message: message,
		Code:    code,
	}
}

type DomainError struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

func (e DomainError) Error() string {
	return fmt.Sprintf("Domain Error: %s (%d)", e.Message, e.Code)
}

func NewDomainError(message string, code int) DomainError {
	return DomainError{
		Message: message,
		Code:    code,
	}
}

type GeneralError struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

func (e GeneralError) Error() string {
	return fmt.Sprintf("General Error: %s (%d)", e.Message, e.Code)
}

func New(message string, code int) error {
	return GeneralError{
		Message: message,
		Code:    code,
	}
}

type ValidationError struct {
	Message   string   `json:"message"`
	Code      int      `json:"code"`
	FieldErrs []string `json:"fieldErrors"`
}

func (e ValidationError) Error() string {
	return "Validation Error"
}

func NewValidationError(message string, fieldErrs []error, code int) ValidationError {
	errStrings := make([]string, 0, len(fieldErrs))
	for _, err := range fieldErrs {
		errStrings = append(errStrings, err.Error())
	}
	return ValidationError{
		Message:   message,
		Code:      code,
		FieldErrs: errStrings,
	}
}

func IsDomain(err error) bool {
	_, ok := err.(DomainError)
	return ok
}

func IsApplication(err error) bool {
	_, ok := err.(ApplicationError)
	return ok
}

func IsValidation(err error) bool {
	_, ok := err.(ValidationError)
	return ok
}

func IsGeneral(err error) bool {
	_, ok := err.(GeneralError)
	return ok
}

func EncodeError(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")

	switch e := err.(type) {
	case ValidationError:
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(e)
	case DomainError:
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(e)
	case ApplicationError:
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(e)
	case GeneralError:
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(e)
	default:
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"message": "internal server error",
			"code":    "500",
		})
	}
}
