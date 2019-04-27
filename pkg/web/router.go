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
	bh := baseHandler{
		Logger: cfg.Logger,
		Renderer: &templateRenderer{
			TemplateDir: cfg.TemplateDir,
		},
	}
	login := router.PathPrefix("/login").Subrouter()
	login.Methods(http.MethodGet).Handler(&loginPageHandler{
		baseHandler: bh,
		LoginAction: "/login",
		Loginer:     cfg.Loginer,
	})

	return router
}
