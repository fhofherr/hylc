package web

import (
	"fmt"
	"net/http"

	"github.com/fhofherr/golf/log"
	"github.com/pkg/errors"
)

type baseHandler struct {
	Logger   log.Logger
	Renderer *templateRenderer
}

func (h *baseHandler) badRequest(msg string, w http.ResponseWriter, req *http.Request) {
	status := http.StatusBadRequest
	statusText := http.StatusText(status)
	model := errorModel{
		PageTitle:  fmt.Sprintf("%d - %s", status, statusText),
		Status:     status,
		StatusText: statusText,
		Message:    msg,
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(status)
	err := h.Renderer.Render(w, "error.html", model)
	logError(h.Logger, errors.Wrap(err, "rendering error"))
}

func (h *baseHandler) internalServerError(w http.ResponseWriter, req *http.Request) {
	status := http.StatusInternalServerError
	statusText := http.StatusText(status)
	model := errorModel{
		PageTitle:  fmt.Sprintf("%d - %s", status, statusText),
		Status:     status,
		StatusText: statusText,
		Message:    statusText,
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusInternalServerError)
	err := h.Renderer.Render(w, "error.html", model)
	logError(h.Logger, errors.Wrap(err, "rendering error"))
}

type errorModel struct {
	PageTitle  string
	Status     int
	StatusText string
	Message    string
}

func logError(logger log.Logger, err error) {
	if err == nil {
		return
	}
	log.Log(logger,
		"level", "error",
		"message", fmt.Sprintf("%+v", err))
}
