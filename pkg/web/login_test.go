package web_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/fhofherr/hylc/pkg/web"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRenderLoginPageHandler(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		challenge  string
		loginURL   string
		loginError error
		bodyPred   func(*testing.T, string)
	}{
		{
			name:       "render login page",
			challenge:  "12345",
			statusCode: http.StatusOK,
			loginError: mockLoginRequiredError{},
			bodyPred: func(t *testing.T, body string) {
				assert.Contains(t, body, `id="login"`)
			},
		},
		{
			name:       "render login page",
			challenge:  "12345",
			statusCode: http.StatusFound,
			loginURL:   "https://login/successful",
			bodyPred: func(t *testing.T, body string) {
				assert.Contains(t, body, `id="login"`)
			},
		},
		{
			name:       "missing login_challenge parameter",
			statusCode: http.StatusBadRequest,
			bodyPred: func(t *testing.T, body string) {
				assert.Contains(t, body, "Missing login_challenge parameter")
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			loginURL := &url.URL{}
			ml := &mockLoginer{}
			if tt.loginURL != "" {
				var err error
				loginURL, err = url.Parse(tt.loginURL)
				assert.NoError(t, err)
			}
			ml.On("Login", tt.challenge, "", "").
				Return(loginURL, tt.loginError)

			cfg := web.PublicRouterConfig{
				TemplateDir: "./template",
				Loginer:     ml,
			}
			handler := web.NewPublicRouter(cfg)

			path := "/login"
			if tt.challenge != "" {
				path = fmt.Sprintf("%s?login_challenge=%s", path, tt.challenge)
			}
			req := httptest.NewRequest(http.MethodGet, path, nil)
			rr := httptest.NewRecorder()

			handler.ServeHTTP(rr, req)

			assert.Equal(t, tt.statusCode, rr.Code)
			if tt.loginURL != "" {
				assert.Equal(t, tt.loginURL, rr.Header().Get("Location"))
			} else {
				assert.Equal(t, "text/html; charset=utf-8", rr.Header().Get("Content-Type"))
				tt.bodyPred(t, rr.Body.String())
			}
		})
	}

}

func TestRenderLoginPageHandler_RenderingError(t *testing.T) {
	ml := &mockLoginer{}
	ml.Test(t)
	cfg := web.PublicRouterConfig{
		TemplateDir: "./missing-template-dir",
		Loginer:     ml,
	}
	handler := web.NewPublicRouter(cfg)
	path := "/login?login_challenge=12345"
	req := httptest.NewRequest(http.MethodGet, path, nil)
	rr := httptest.NewRecorder()

	ml.On("Login", "12345", "", "").
		Return(nil, mockLoginRequiredError{})
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	assert.Equal(t, "text/html; charset=utf-8", rr.Header().Get("Content-Type"))
	assert.Contains(t, rr.Body.String(), "Internal server error")
}

func TestRenderLoginPageHandler_UnexpectedLoginError(t *testing.T) {
	ml := &mockLoginer{}
	ml.Test(t)
	cfg := web.PublicRouterConfig{
		TemplateDir: "./missing-template-dir",
		Loginer:     ml,
	}
	handler := web.NewPublicRouter(cfg)
	path := "/login?login_challenge=12345"
	req := httptest.NewRequest(http.MethodGet, path, nil)
	rr := httptest.NewRecorder()
	ml.On("Login", "12345", "", "").
		Return(nil, errors.New("something went wrong"))
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	assert.Equal(t, "text/html; charset=utf-8", rr.Header().Get("Content-Type"))
	assert.Contains(t, rr.Body.String(), "Internal server error")
}

type mockLoginer struct {
	mock.Mock
}

func (ml *mockLoginer) Login(challenge, username, password string) (*url.URL, error) {
	args := ml.Called(challenge, username, password)
	loginUrl := args.Get(0)
	if loginUrl != nil {
		return loginUrl.(*url.URL), args.Error(1)
	}
	return nil, args.Error(1)
}

type mockLoginRequiredError struct{}

func (mockLoginRequiredError) LoginRequired() bool {
	return true
}

func (mockLoginRequiredError) Error() string {
	return ""
}
