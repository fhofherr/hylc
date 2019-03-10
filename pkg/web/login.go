package web

import (
	"net/http"

	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type renderLoginPageHandler struct {
	logger   *zap.Logger
	renderer *templateRenderer
}

func (h *renderLoginPageHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	logger := h.logger.Sugar()
	err := h.renderer.execute(w, nil)
	if err != nil {
		logger.Errorf("%+v", errors.Wrap(err, "render login page"))
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusInternalServerError)
	}
}
