package web_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/fhofherr/golf/log"
	"github.com/fhofherr/hylc/pkg/web"
	"github.com/stretchr/testify/assert"
)

func TestRenderLoginPageHandler(t *testing.T) {
	tests := []struct {
		name       string
		cfg        web.PublicRouterConfig
		statusCode int
		bodyPred   func(*testing.T, string)
	}{
		{
			name: "render login page",
			cfg: web.PublicRouterConfig{
				Logger:      log.NewNOPLogger(),
				TemplateDir: "./template",
			},
			statusCode: http.StatusOK,
			bodyPred: func(t *testing.T, body string) {
				assert.Contains(t, body, `id="login"`)
			},
		},
		{
			name: "rendering error",
			cfg: web.PublicRouterConfig{
				Logger:      log.NewNOPLogger(),
				TemplateDir: "./missing-template-dir",
			},
			statusCode: http.StatusInternalServerError,
			bodyPred:   func(t *testing.T, s string) {},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := web.NewPublicRouter(tt.cfg)

			req := httptest.NewRequest(http.MethodGet, "/login", nil)
			rr := httptest.NewRecorder()

			handler.ServeHTTP(rr, req)

			assert.Equal(t, tt.statusCode, rr.Code)
			assert.Equal(t, "text/html; charset=utf-8", rr.Header().Get("Content-Type"))
			tt.bodyPred(t, rr.Body.String())
		})
	}

}
