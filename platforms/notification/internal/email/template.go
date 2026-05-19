package email

import (
	"bytes"
	"html/template"
)

type TemplateEngine struct {
	templates map[string]*template.Template
}

func NewTemplateEngine() *TemplateEngine {
	return &TemplateEngine{
		templates: make(map[string]*template.Template),
	}
}

func (e *TemplateEngine) Register(name, tmpl string) error {
	t, err := template.New(name).Parse(tmpl)
	if err != nil {
		return err
	}
	e.templates[name] = t
	return nil
}

func (e *TemplateEngine) Render(name string, data interface{}) (string, error) {
	t, ok := e.templates[name]
	if !ok {
		return "", nil
	}

	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}
