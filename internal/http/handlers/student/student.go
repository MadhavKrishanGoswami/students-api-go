package student

import (
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/MadhavKrishanGoswami/students-api/internal/storage"
	"github.com/MadhavKrishanGoswami/students-api/internal/types"
	"github.com/MadhavKrishanGoswami/students-api/internal/utils/response"
	"github.com/go-playground/validator/v10"
)

func New(storage storage.Storage) http.HandlerFunc {
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

		// Create the student in the storage
		id, err := storage.CreateStudent(student.Name, student.Email, student.Age)

		slog.Info("Student created", "id", id, "name", student.Name)

		if err != nil {
			slog.Error("Failed to create student", "error", err)
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		response.WriteJson(w, http.StatusCreated, map[string]string{"id": string(id)})
	}
}

func GetById(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slog.Info("Getting Student by ID")
		id := r.PathValue("id")
		if id == "" {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(errors.New("id is required")))
			return
		}
		intid, err := strconv.Atoi(id)
		if err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(errors.New("invalid intid")))
			return
		}
		student, err := storage.GetStudentById(int64(intid))
		if err != nil {
			response.WriteJson(w, http.StatusNotFound, response.GeneralError(errors.New("student not found")))
			return
		}
		// Handle other errors
		if err != nil {
			slog.Error("Failed to get student", "error", err)
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		response.WriteJson(w, http.StatusOK, student)
	}
}

func GetList(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slog.Info("Getting Student List")
		// For simplicity, we are not implementing pagination or filtering here
		students, err := storage.GetListStudents(100, 0) // Get first 100 students
		if err != nil {
			slog.Error("Failed to get students", "error", err)
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		response.WriteJson(w, http.StatusOK, students)
	}
}
