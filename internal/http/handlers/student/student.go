package student

import (
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"

	"github.com/MadhavKrishanGoswami/students-api/internal/types"
	"github.com/MadhavKrishanGoswami/students-api/internal/utils/response"
	"github.com/go-playground/validator/v10"
)

func New() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slog.Info("Creatating Student")
		var student types.Student
		err := json.NewDecoder(r.Body).Decode(&student)
		if errors.Is(err, io.EOF) {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(errors.New("request body is empty")))
			return
		}
		if err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}
		// Validate the student data
		if err := validator.New().Struct(student); err != nil {
			slog.Error("Validation error", "error", err)
			response.WriteJson(w, http.StatusBadRequest, response.ValidationError(err.(validator.ValidationErrors)))
			return
		}

		response.WriteJson(w, http.StatusCreated, map[string]string{"success": "OK"})
	}
}
