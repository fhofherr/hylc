package web_test

import (
	"fmt"
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
		challenge  string
		bodyPred   func(*testing.T, string)
	}{
		{
			name: "render login page",
			cfg: web.PublicRouterConfig{
				Logger:      log.NewNOPLogger(),
				TemplateDir: "./template",
			},
			challenge:  "12345",
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
			challenge:  "54321",
			statusCode: http.StatusInternalServerError,
			bodyPred: func(t *testing.T, body string) {
				assert.Contains(t, body, "Internal server error")
			},
		},
		{
			name: "missing login_challenge parameter",
			cfg: web.PublicRouterConfig{
				Logger:      log.NewNOPLogger(),
				TemplateDir: "./template",
			},
			statusCode: http.StatusBadRequest,
			bodyPred: func(t *testing.T, body string) {
				assert.Contains(t, body, "Missing login_challenge parameter")
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := web.NewPublicRouter(tt.cfg)

			path := "/login"
			if tt.challenge != "" {
				path = fmt.Sprintf("%s?login_challenge=%s", path, tt.challenge)
			}
			req := httptest.NewRequest(http.MethodGet, path, nil)
			rr := httptest.NewRecorder()

			handler.ServeHTTP(rr, req)

			assert.Equal(t, tt.statusCode, rr.Code)
			assert.Equal(t, "text/html; charset=utf-8", rr.Header().Get("Content-Type"))
			tt.bodyPred(t, rr.Body.String())
		})
	}

}
