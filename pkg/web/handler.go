package web

import (
	"fmt"
	"net/http"

	"github.com/fhofherr/golf/log"
)

type baseHandler struct {
	Logger   log.Logger
	Renderer *templateRenderer
}

func (h *baseHandler) badRequest(msg string, w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusBadRequest)
	_, err := w.Write([]byte(msg))
	if err != nil {
		log.Log(h.Logger, "level", "error", "message", fmt.Sprintf("%+v", err))
	}
}

func (h *baseHandler) internalServerError(err error, w http.ResponseWriter, req *http.Request) {
	log.Log(h.Logger,
		"level", "error",
		"message", fmt.Sprintf("%+v", err))
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusInternalServerError)
	_, _ = w.Write([]byte("Internal server error"))
}
