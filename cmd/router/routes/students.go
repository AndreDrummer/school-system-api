package routes

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	schoolsystem "github.com/AndreDrummer/school-system-api/cmd/controller"
	"github.com/AndreDrummer/school-system-api/cmd/domain"
	apperrors "github.com/AndreDrummer/school-system-api/cmd/errors"
	httputils "github.com/AndreDrummer/school-system-api/cmd/http"

	"github.com/go-chi/chi/v5"
)

func handlePostRequests(w http.ResponseWriter, r *http.Request, student *domain.Student) error {
	if err := json.NewDecoder(r.Body).Decode(student); err != nil {
		jsonDecondingError := &apperrors.JsonDecodingError{
			Type: fmt.Sprintf("%T", student),
			Err:  err,
		}

		slog.Error(jsonDecondingError.Error())

		httputils.SendResponse(
			w,
			httputils.Response{Error: "Invalid request: body malformed"},
			http.StatusBadRequest,
		)

		return jsonDecondingError
	}

	return nil
}

func Students(r chi.Router) {
	r.Get("/students", listAll)
	r.Get("/students/{id:[0-9]+}", getByID)
	r.Delete("/students/{id:[0-9]+}", delete)
	r.Put("/students/{id:[0-9]+}", update)
	r.Post("/students", create)
}

var listAll = func(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	students, err := schoolsystem.AllStudents()

	if err != nil {
		slog.Error(err.Error())
		httputils.SendResponse(
			w,
			httputils.Response{Error: "Could not get the students list."},
			http.StatusInternalServerError,
		)
		return
	}

	slog.Info("Listing students...")
	httputils.SendResponse(
		w,
		httputils.Response{Data: students},
		http.StatusOK,
	)
}

var getByID = func(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	studentID := chi.URLParam(r, "id")
	studentIDInt, err := strconv.Atoi(studentID)

	if err != nil {
		slog.Error(err.Error())
		httputils.SendResponse(
			w,
			httputils.Response{Error: "Could not get student."},
			http.StatusBadRequest,
		)
		return
	}

	student, ok := schoolsystem.GetStudentByID(studentIDInt)

	if !ok {
		httputils.SendResponse(
			w,
			httputils.Response{Data: "Student not found"},
			http.StatusNotFound,
		)
		return
	}

	slog.Info("Student Found!")
	httputils.SendResponse(
		w,
		httputils.Response{Data: student},
		http.StatusOK,
	)
}

var update = func(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var student domain.Student
	if err := handlePostRequests(w, r, &student); err == nil {
		studentID := chi.URLParam(r, "id")
		studentIDInt, err := strconv.Atoi(studentID)

		if err != nil {
			slog.Error(err.Error())
			httputils.SendResponse(
				w,
				httputils.Response{Error: "Could not update student data."},
				http.StatusInternalServerError,
			)
			return
		}

		student.ID = studentIDInt
		_, err = schoolsystem.UpdateStudent(&student)

		if err != nil {
			notFoundError := &apperrors.NotFoundError{}

			if errors.As(err, &notFoundError) {
				errorMessage := fmt.Sprintf("Student of ID %v was not found.", studentID)
				slog.Error(errorMessage)

				httputils.SendResponse(
					w,
					httputils.Response{Error: errorMessage},
					http.StatusNotFound,
				)

			} else {
				slog.Error(err.Error())
				httputils.SendResponse(
					w,
					httputils.Response{Error: "Could not update student data."},
					http.StatusInternalServerError,
				)
			}
			return
		} else {
			successMsg := "Student updated successfully."
			slog.Info(successMsg)
			httputils.SendResponse(w,
				httputils.Response{Data: student},
				http.StatusOK,
			)
		}
	}
}

var create = func(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var student domain.Student
	if err := handlePostRequests(w, r, &student); err == nil {
		_, err := schoolsystem.AddStudent(&student)

		if err != nil {
			slog.Error(err.Error())
			httputils.SendResponse(
				w,
				httputils.Response{
					Error: "Was not possible to add a new student",
				},
				http.StatusBadRequest,
			)
			return
		}

		slog.Info("Student Added successfully!")
		httputils.SendResponse(
			w,
			httputils.Response{Data: student},
			http.StatusOK,
		)
	}
}

var delete = func(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	studentID := chi.URLParam(r, "id")
	studentIDInt, err := strconv.Atoi(studentID)

	if err != nil {
		slog.Error(err.Error())
		httputils.SendResponse(
			w,
			httputils.Response{
				Error: "Was not possible to delete student.",
			},
			http.StatusInternalServerError,
		)
		return
	}

	_, err = schoolsystem.RemoveStudent(studentIDInt)

	if err != nil {
		notFoundError := &apperrors.NotFoundError{}
		if errors.As(err, &notFoundError) {
			slog.Error(err.Error())
			httputils.SendResponse(
				w,
				httputils.Response{
					Error: fmt.Sprintf("Student of ID %d was not found", studentIDInt),
				},
				http.StatusNotFound,
			)
		} else {
			slog.Error(err.Error())
			httputils.SendResponse(
				w,
				httputils.Response{
					Error: "Was not possible to delete student.",
				},
				http.StatusInternalServerError,
			)
		}

		return
	}

	slog.Info("Student deleted successfully!")
	httputils.SendResponse(
		w,
		httputils.Response{Data: "Student deleted successfully!"},
		http.StatusOK,
	)
}
