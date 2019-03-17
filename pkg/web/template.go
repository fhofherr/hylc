package web

import (
	"html/template"
	"io"
	"os"
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
		layoutDir := filepath.Join(templateDir, "layout")
		layouts, err := findLayouts(layoutDir)
		if err != nil {
			r.err = errors.Wrap(err, "find layouts")
			return
		}
		filenames := append([]string{r.templatePath}, layouts...)
		tmpl, err := template.ParseFiles(filenames...)
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

func findLayouts(layoutsDir string) ([]string, error) {
	var layouts []string

	err := filepath.Walk(layoutsDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		filename := info.Name()
		if filepath.Ext(filename) != ".html" {
			return nil
		}
		layouts = append(layouts, path)
		return nil
	})
	return layouts, err
}
