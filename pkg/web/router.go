package web

import (
	"net/http"

	"github.com/fhofherr/golf/log"
	"github.com/gorilla/mux"
)

// PublicRouterConfig configures the public router.
type PublicRouterConfig struct {
	Logger      log.Logger
	Loginer     Loginer
	TemplateDir string
}

// NewPublicRouter creates a new http.Handler serving all of hylc's publicly
// available routes.
func NewPublicRouter(cfg PublicRouterConfig) http.Handler {
	router := mux.NewRouter()

	login := router.PathPrefix("/login").Subrouter()
	login.Methods(http.MethodGet).Handler(&renderLoginPageHandler{
		Logger:      cfg.Logger,
		LoginAction: "/login",
		Loginer:     cfg.Loginer,
		Renderer: &templateRenderer{
			TemplateDir: cfg.TemplateDir,
		},
	})

	return router
}
