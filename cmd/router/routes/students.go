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

func handleRequests(w http.ResponseWriter, r *http.Request, student *domain.Student) error {
	if err := json.NewDecoder(r.Body).Decode(student); err != nil {
		slog.Error(fmt.Sprintf("Json decoding error: %s", err.Error()))
		httputils.SendResponse(
			w,
			httputils.Response{Error: "Invalid request: body malformed"},
			http.StatusBadRequest,
		)
		return fmt.Errorf("invalid Json: %v", err)
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

	httputils.SendResponse(
		w,
		httputils.Response{Data: student},
		http.StatusOK,
	)
}

var update = func(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var student domain.Student
	if err := handleRequests(w, r, &student); err == nil {
		studentID := chi.URLParam(r, "id")
		studentIDInt, err := strconv.Atoi(studentID)

		if err != nil {
			slog.Error(err.Error())
			httputils.SendResponse(
				w,
				httputils.Response{Error: "Could not update user."},
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
					httputils.Response{Error: "Could not update user."},
					http.StatusInternalServerError,
				)
			}
			return
		} else {
			successMsg := "Student updated successfully."
			slog.Info(successMsg)
			httputils.SendResponse(w,
				httputils.Response{Data: successMsg},
				http.StatusOK,
			)
		}
	}

}

var create = func(w http.ResponseWriter, r *http.Request) {

}

var delete = func(w http.ResponseWriter, r *http.Request) {

}
