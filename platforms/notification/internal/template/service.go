package template

import (
	"bytes"
	"context"
	"html/template"
	"strings"
)

type Service interface {
	RenderTemplate(ctx context.Context, tmplID string, vars map[string]interface{}) (string, string, error)
	CreateTemplate(ctx context.Context, req *CreateTemplateRequest) (*Template, error)
	GetTemplate(ctx context.Context, id string) (*Template, error)
	GetTemplateByName(ctx context.Context, name string) (*Template, error)
	ListTemplates(ctx context.Context) ([]*Template, error)
	UpdateTemplate(ctx context.Context, id string, req *UpdateTemplateRequest) (*Template, error)
	DeleteTemplate(ctx context.Context, id string) error
	ListVersions(ctx context.Context, templateID string) ([]*TemplateVersion, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) RenderTemplate(ctx context.Context, tmplID string, vars map[string]interface{}) (string, string, error) {
	t, err := s.repo.GetByID(ctx, tmplID)
	if err != nil {
		return "", "", err
	}
	if t == nil {
		t, err = s.repo.GetByName(ctx, tmplID)
		if err != nil {
			return "", "", err
		}
	}
	if t == nil {
		return "", "", nil
	}

	subjectTmpl, err := template.New("subject").Parse(t.Subject)
	if err != nil {
		return "", "", err
	}

	bodyTmpl, err := template.New("body").Parse(t.Body)
	if err != nil {
		return "", "", err
	}

	var subjectBuf bytes.Buffer
	if err := subjectTmpl.Execute(&subjectBuf, vars); err != nil {
		return "", "", err
	}

	var bodyBuf bytes.Buffer
	if err := bodyTmpl.Execute(&bodyBuf, vars); err != nil {
		return "", "", err
	}

	return subjectBuf.String(), bodyBuf.String(), nil
}

func (s *service) CreateTemplate(ctx context.Context, req *CreateTemplateRequest) (*Template, error) {
	t := &Template{
		Name:      req.Name,
		Subject:   req.Subject,
		Body:      req.Body,
		Variables: req.Variables,
	}

	if err := s.repo.Create(ctx, t); err != nil {
		return nil, err
	}

	return t, nil
}

func (s *service) GetTemplate(ctx context.Context, id string) (*Template, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *service) GetTemplateByName(ctx context.Context, name string) (*Template, error) {
	return s.repo.GetByName(ctx, name)
}

func (s *service) ListTemplates(ctx context.Context) ([]*Template, error) {
	return s.repo.List(ctx)
}

func (s *service) UpdateTemplate(ctx context.Context, id string, req *UpdateTemplateRequest) (*Template, error) {
	t, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if t == nil {
		return nil, nil
	}

	if req.Subject != nil {
		t.Subject = *req.Subject
	}
	if req.Body != nil {
		t.Body = *req.Body
	}
	if req.Variables != nil {
		t.Variables = *req.Variables
	}

	if err := s.repo.Update(ctx, t); err != nil {
		return nil, err
	}

	return t, nil
}

func (s *service) DeleteTemplate(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

func (s *service) ListVersions(ctx context.Context, templateID string) ([]*TemplateVersion, error) {
	return s.repo.ListVersions(ctx, templateID)
}

func validateTemplateBody(body string) bool {
	return strings.Contains(body, "{{") && strings.Contains(body, "}}")
}
