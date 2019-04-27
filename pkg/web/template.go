package web

import (
	"html/template"
	"io"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"

	"github.com/pkg/errors"
)

// DefaultTemplateDirectory to use if nothing else is configured.
var DefaultTemplateDirectory = "./pkg/web/template"

type templateRenderer struct {
	TemplateDir string
	templates   map[string]*template.Template
	layouts     []string
	mu          sync.RWMutex
	initialized uint32
}

// Render renders the template with the file name filename. If Render returns
// an error the template was not rendered.
func (r *templateRenderer) Render(w io.Writer, filename string, data interface{}) error {
	if err := r.initialize(); err != nil {
		return errors.Wrap(err, "initialize renderer")
	}
	tmpl, ok := r.template(filename)
	if !ok {
		if err := r.loadTemplate(filename); err != nil {
			return errors.Wrapf(err, "load template: %s", filename)
		}
		tmpl, _ = r.template(filename)
	}
	if err := tmpl.Execute(w, data); err != nil {
		return errors.Wrapf(err, "Render template: %s", filename)
	}
	return nil
}

func (r *templateRenderer) initialize() error {
	// Don't do anything if we have been initialized already.
	if atomic.LoadUint32(&r.initialized) == 1 {
		return nil
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	// Another go-routine has acquired the lock before us
	// and initialized us. Abort.
	if r.initialized == 1 {
		return nil
	}

	r.templates = make(map[string]*template.Template)
	if r.TemplateDir == "" {
		r.TemplateDir = DefaultTemplateDirectory
	}
	layoutDir := filepath.Join(r.TemplateDir, "layout")
	layouts, err := findLayouts(layoutDir)
	if err != nil {
		return errors.Wrap(err, "find layouts")
	}
	r.layouts = layouts
	r.initialized = 1
	return nil
}

func (r *templateRenderer) template(filename string) (*template.Template, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	tmpl, ok := r.templates[filename]
	return tmpl, ok
}

func (r *templateRenderer) loadTemplate(filename string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Another go-routine has loaded the template already. Abort.
	if _, ok := r.templates[filename]; ok {
		return nil
	}

	tmplPath := filepath.Join(r.TemplateDir, filename)
	filenames := append([]string{tmplPath}, r.layouts...)
	tmpl, err := template.ParseFiles(filenames...)
	if err != nil {
		return errors.Wrapf(err, "parse template file %s", tmplPath)
	}
	r.templates[filename] = tmpl
	return nil
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
