package httputils

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
)

type Response struct {
	Error string `json:"error,omitempty"`
	Data  any    `json:"data,omitempty"`
}

func SendResponse(w http.ResponseWriter, res Response, status int) {
	w.Header().Set("Content-Type", "application/json")

	data, err := json.Marshal(res)

	if err != nil {
		slog.Error(fmt.Sprintf("error %s parsing response.", err.Error()))
		SendResponse(
			w,
			Response{Error: "somenthing went wrong"},
			http.StatusInternalServerError,
		)
		return
	}

	w.WriteHeader(status)
	if _, err := w.Write(data); err != nil {
		slog.Error(fmt.Sprintf("error %s writing response.", err.Error()))
		return
	}
}
