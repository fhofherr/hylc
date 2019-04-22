package web

import (
	"fmt"
	"net/http"

	"github.com/fhofherr/golf/log"
	"github.com/pkg/errors"
)

type renderLoginPageHandler struct {
	Logger      log.Logger
	LoginAction string
	renderer    *templateRenderer
}

func (h *renderLoginPageHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
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
		log.Log(h.Logger,
			"level", "error",
			"message", fmt.Sprintf("%+v", errors.Wrap(err, "render login page")))
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
