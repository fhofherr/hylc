package web

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/fhofherr/golf/log"
	"github.com/pkg/errors"
)

// Loginer performs a login using the provided challenge, username, and password.
//
// The username and password params may be the empty string. In this case the
// Loginer checks if the login can be skipped.
//
// If the login was successful or if the Loginer chose to deny the login request,
// Login returns an URL and nil as an error. In all other cases error is non-nil
// and the returned URL must be ignored. If Login wants to signal that a login
// page has to be displayed, the returned error implements LoginRequired.
type Loginer interface {
	Login(challenge, username, password string) (*url.URL, error)
}

type renderLoginPageHandler struct {
	Logger      log.Logger
	LoginAction string
	Renderer    *templateRenderer
	Loginer     Loginer
}

func (h *renderLoginPageHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	query := req.URL.Query()
	challenge, ok := query["login_challenge"]
	if !ok || len(challenge) != 1 {
		h.renderMissingChallenge(w, req)
		return
	}
	redirectUrl, err := h.Loginer.Login(challenge[0], "", "")
	if isLoginRequired(err) {
		h.renderLoginPage(w, req)
		return
	}
	if err != nil {
		log.Log(h.Logger,
			"level", "error",
			"message", fmt.Sprintf("%+v", errors.Wrap(err, "login error")))
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("Internal server error"))
		return
	}
	http.Redirect(w, req, redirectUrl.String(), http.StatusFound)
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
	err := h.Renderer.execute(w, model)
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

// LoginRequired wraps the LoginRequired method which returns if a login
// is actually required.
type LoginRequired interface {
	LoginRequired() bool
}

func isLoginRequired(err error) bool {
	if lerr, ok := err.(LoginRequired); ok {
		return lerr.LoginRequired()
	}
	return false
}
