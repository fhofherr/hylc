package web

import (
	"net/http"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

// PublicRouterConfig configures the public router.
type PublicRouterConfig struct {
	Logger      *zap.Logger
	TemplateDir string
}

// NewPublicRouter creates a new http.Handler serving all of hylc's publicly
// available routes.
func NewPublicRouter(cfg PublicRouterConfig) http.Handler {
	router := mux.NewRouter()

	login := router.PathPrefix("/login").Subrouter()
	login.Methods(http.MethodGet).Handler(&renderLoginPageHandler{
		logger: cfg.Logger,
		renderer: &templateRenderer{
			Filename:    "login.html",
			TemplateDir: cfg.TemplateDir,
		},
	})

	return router
}
