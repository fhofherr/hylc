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
	query := req.URL.Query()
	_, ok := query["login_challenge"]
	if !ok {
		h.renderMissingChallenge(w, req)
		return
	}
	h.renderLoginPage(w, req)
}

func (h *renderLoginPageHandler) renderMissingChallenge(w http.ResponseWriter, req *http.Request) {
	log.Log(h.Logger, "level", "info", "message", "/login called without login_challenge")
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusBadRequest)
	_, _ = w.Write([]byte("Missing login_challenge parameter"))
}

func (h *renderLoginPageHandler) renderLoginPage(w http.ResponseWriter, req *http.Request) {
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
		_, _ = w.Write([]byte("Internal server error"))
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
