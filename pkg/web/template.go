package web

import (
	"html/template"
	"io"
	"path/filepath"
	"sync"

	"github.com/pkg/errors"
)

// DefaultTemplateDirectory to use if nothing else is configured.
var DefaultTemplateDirectory = "./pkg/web/template"

type templateRenderer struct {
	TemplateDir  string
	Filename     string
	templatePath string
	template     *template.Template
	once         sync.Once
	err          error
}

func (r *templateRenderer) execute(w io.Writer, data interface{}) error {
	r.once.Do(func() {
		templateDir := r.TemplateDir
		if templateDir == "" {
			templateDir = DefaultTemplateDirectory
		}
		r.templatePath = filepath.Join(templateDir, r.Filename)
		tmpl, err := template.ParseFiles(r.templatePath)
		if err != nil {
			r.err = errors.Wrapf(err, "parse template file %s", r.templatePath)
			return
		}
		r.template = tmpl
	})
	if r.err != nil {
		return r.err
	}
	err := r.template.Execute(w, data)
	if err != nil {
		r.err = errors.Wrapf(err, "execute template %s", r.templatePath)
	}
	return r.err
}
