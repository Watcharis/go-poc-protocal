package pkg

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-playground/validator"
	"github.com/google/uuid"
)

func GenerateUUID() string {
	u := uuid.New().String()
	return strings.ReplaceAll(u, "-", "")
}

func ValidateStruct(payload interface{}) error {
	validate := validator.New()
	err := validate.Struct(payload)
	if err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			fmt.Println(err)
		}

		for _, err := range err.(validator.ValidationErrors) {
			return errors.New(err.Field() + " is " + err.Tag())
		}
	}
	return nil
}

func HealthCheck(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("invalid http methods"))
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("service is running!!"))
}

func SetContentType(w http.ResponseWriter, contentType string) http.ResponseWriter {
	w.Header().Set("Content-Type", contentType)
	return w
}

func SetHttpStatusCode(w http.ResponseWriter, code int) http.ResponseWriter {
	w.WriteHeader(code)
	return w
}
