package utils

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
)

// ============================================================================
// RESPONSE HELPERS
// ============================================================================

// writeJSONError sends a JSON error response
func WriteJSONError(w http.ResponseWriter, code int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "error",
		"message": msg,
	})
}

// writeJSONSuccess sends a JSON success response with data
func WriteJSONSuccess(w http.ResponseWriter, code int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "success",
		"data":   data,
	})
}

// writeJSONList sends a JSON success response with a list and count
func WriteJSONList(w http.ResponseWriter, code int, data interface{}, count int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "success",
		"count":  count,
		"data":   data,
	})
}

// formatValidationErrors converts validator errors to readable message
func FormatValidationErrors(errs validator.ValidationErrors) string {
	var messages []string
	for _, err := range errs {
		switch err.Tag() {
		case "required":
			messages = append(messages, err.Field()+" is required")
		case "email":
			messages = append(messages, err.Field()+" must be a valid email")
		case "min":
			messages = append(messages, err.Field()+" must be at least "+err.Param()+" characters")
		case "max":
			messages = append(messages, err.Field()+" must be at most "+err.Param()+" characters")
		default:
			messages = append(messages, err.Field()+" validation failed")
		}
	}
	return strings.Join(messages, "; ")
}
