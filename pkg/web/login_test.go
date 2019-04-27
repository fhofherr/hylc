package web_test

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/fhofherr/hylc/pkg/web"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRenderLoginPageHandler(t *testing.T) {
	tests := []struct {
		name           string
		loginChallenge string
		loginURL       *url.URL
		loginErr       error
		status         int
		header         http.Header
	}{
		{
			name:           "redirect to login url",
			loginChallenge: "12345",
			loginURL:       web.MustParseURL(t, "https://login/successful"),
			status:         http.StatusFound,
			header:         web.NewHTTPHeader(t, "Location", "https://login/successful"),
		},
		{
			name:           "render login page",
			loginChallenge: "12345",
			loginErr:       mockLoginRequiredError{},
			status:         http.StatusOK,
		},
		{
			name:     "fail on missing login challenge",
			loginURL: web.MustParseURL(t, "https://login/successful"),
			status:   http.StatusBadRequest,
		},
		{
			name:           "unexpected login error",
			loginChallenge: "12345",
			loginErr:       errors.New("something went wrong"),
			status:         http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			var values url.Values
			ml := NewMockLoginer(t)
			ml.On("Login", tt.loginChallenge, "", "").
				Return(tt.loginURL, tt.loginErr)
			cfg := web.PublicRouterConfig{
				TemplateDir: "./template",
				Loginer:     ml,
			}
			handler := web.NewPublicRouter(cfg)
			if tt.loginChallenge != "" {
				values = web.NewURLValues(t, "login_challenge", tt.loginChallenge)
			}
			header, body := web.AssertHTTP(
				t, handler.ServeHTTP, http.MethodGet, "/login", values, nil, tt.status)
			web.AssertHTTPGoldenFile(t, header, body)
		})
	}
}

func TestRenderLoginPageHandler_RenderingError(t *testing.T) {
	ml := NewMockLoginer(t)
	ml.On("Login", mock.Anything, mock.Anything, mock.Anything).
		Return(nil, mockLoginRequiredError{})
	cfg := web.PublicRouterConfig{
		TemplateDir: "./missing-template-dir",
		Loginer:     ml,
	}
	handler := web.NewPublicRouter(cfg)
	values := web.NewURLValues(t, "login_challenge", "12345")
	assert.HTTPError(t, handler.ServeHTTP, http.MethodGet, "/login", values)
}

type mockLoginer struct {
	mock.Mock
}

func NewMockLoginer(t *testing.T) *mockLoginer {
	ml := &mockLoginer{}
	ml.Test(t)
	return ml
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
