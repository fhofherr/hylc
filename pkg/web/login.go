package web

import (
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

type loginPageHandler struct {
	baseHandler
	LoginAction string
	Loginer     Loginer
}

func (h *loginPageHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	query := req.URL.Query()
	challenge, ok := query["login_challenge"]
	if !ok || len(challenge) != 1 {
		h.missingChallenge(w, req)
		return
	}
	redirectUrl, err := h.Loginer.Login(challenge[0], "", "")
	if isLoginRequired(err) {
		h.loginPage(w, req)
		return
	}
	if err != nil {
		logError(h.Logger, errors.Wrap(err, "login error"))
		h.internalServerError(w, req)
		return
	}
	http.Redirect(w, req, redirectUrl.String(), http.StatusFound)
}

func (h *loginPageHandler) missingChallenge(w http.ResponseWriter, req *http.Request) {
	log.Log(h.Logger, "level", "info", "message", "/login called without login_challenge")
	h.badRequest("Missing login_challenge parameter", w, req)
}

func (h *loginPageHandler) loginPage(w http.ResponseWriter, req *http.Request) {
	model := loginPageModel{
		PageTitle:           "Login",
		Title:               "Login",
		Action:              h.LoginAction,
		UsernameLabel:       "Username",
		UsernamePlaceholder: "Username",
		PasswordLabel:       "Password",
		SubmitButtonLabel:   "Login",
	}
	err := h.Renderer.Render(w, "login.html", model)
	if err != nil {
		logError(h.Logger, errors.Wrap(err, "render login page"))
		h.internalServerError(w, req)
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
