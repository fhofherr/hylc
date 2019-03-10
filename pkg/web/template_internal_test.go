package web

import (
	"bytes"
	"path/filepath"
	"strings"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRenderTemplate(t *testing.T) {
	tests := []struct {
		name      string
		filename  string
		data      interface{}
		rendered  string
		expectErr bool
	}{
		{"missing template", "missing.html", nil, "", true},
		{"empty template", "empty.html", nil, "empty template", false},
		{"not empty template", "not_empty.html", map[string]string{"Variable": "Value"}, "Value", false},
		{"syntax error in template", "syntax_error.html", map[string]string{"Variable": "Value"}, "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			var buf bytes.Buffer
			renderer := &templateRenderer{
				TemplateDir: "testdata/template",
				Filename:    tt.filename,
			}

			err := renderer.execute(&buf, tt.data)
			if !tt.expectErr && err != nil {
				t.Fatal(err)
			}
			assert.Equal(t, tt.expectErr, err != nil)
			assert.Equal(t, err, renderer.err)
			assert.Equal(t, tt.rendered, strings.TrimSpace(buf.String()))
		})
	}
}

func TestUseDefaultTemplateDirectory(t *testing.T) {
	var buf bytes.Buffer
	renderer := &templateRenderer{
		Filename: "missing.html",
	}
	expectedPath := filepath.Join(DefaultTemplateDirectory, "missing.html")
	err := renderer.execute(&buf, nil)
	// We expect an error here since the DefaultTemplateDirector is not
	// reachable from within the execution directory of the tests.
	assert.Error(t, err)
	assert.Equal(t, expectedPath, renderer.templatePath)
}

func TestRenderTemplateWriteFails(t *testing.T) {
	expectedErr := errors.New("write error")
	w := &mockWriter{}
	w.On("Write", mock.Anything).
		Return(0, expectedErr)

	renderer := &templateRenderer{
		TemplateDir: "testdata/template",
		Filename:    "empty.html",
	}

	err := renderer.execute(w, nil)
	assert.Error(t, expectedErr, errors.Cause(err))
}

type mockWriter struct {
	mock.Mock
}

func (w *mockWriter) Write(bs []byte) (int, error) {
	args := w.Called(bs)
	return args.Int(0), args.Error(1)
}
