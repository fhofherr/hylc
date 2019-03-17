package web

import (
	"net/http"

	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type renderLoginPageHandler struct {
	Logger      *zap.Logger
	LoginAction string
	renderer    *templateRenderer
}

func (h *renderLoginPageHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	logger := h.Logger.Sugar()
	model := loginPageModel{
		PageTitle:           "Login",
		Title:               "Login",
		Action:              h.LoginAction,
		UsernameLabel:       "Username",
		UsernamePlaceholder: "Username",
		PasswordLabel:       "Password",
		SubmitButtonLabel:   "Login",
	}
	err := h.renderer.execute(w, model)
	if err != nil {
		logger.Errorf("%+v", errors.Wrap(err, "render login page"))
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusInternalServerError)
	}
}

type loginPageModel struct {
	PageTitle           string
	Title               string
	Action              string
	UsernameLabel       string
	UsernamePlaceholder string
	PasswordLabel       string
	SubmitButtonLabel   string
}
