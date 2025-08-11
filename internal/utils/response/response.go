package response

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
)

type Responnse struct {
	Status string
	Error  string
}

const (
	StatusOK    = "OK"
	StatusError = "ERROR"
)

func WriteJson(w http.ResponseWriter, status int, data interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	return json.NewEncoder(w).Encode(data)
}

func GeneralError(err error) Responnse {
	return Responnse{
		Status: StatusError,
		Error:  err.Error(),
	}
}

func ValidationError(errs validator.ValidationErrors) Responnse {
	var errMessages []string
	for _, err := range errs {
		switch err.ActualTag() {
		case "required":
			errMessages = append(errMessages, "Field "+err.Field()+" is required")
		case "email":
			errMessages = append(errMessages, "Field "+err.Field()+" must be a valid email address")
		case "gte":
			errMessages = append(errMessages, "Field "+err.Field()+" must be greater than or equal to "+err.Param())
		case "lte":
			errMessages = append(errMessages, "Field "+err.Field()+" must be less than or equal to "+err.Param())
		default:
			errMessages = append(errMessages, "Field "+err.Field()+" failed validation for tag "+err.ActualTag())
		}
	}
	return Responnse{
		Status: StatusError,
		Error:  "Validation failed: " + errMessages[0],
	}
}
