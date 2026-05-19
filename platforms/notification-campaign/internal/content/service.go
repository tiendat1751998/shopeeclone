package content

import (
	"bytes"
	"context"
	"errors"
	"html/template"
	"math/rand"
	"strings"
)

var (
	ErrNotFound = errors.New("content: not found")
)

type Service interface {
	CreateTemplate(ctx context.Context, req *CreateTemplateRequest) (*ContentTemplate, error)
	ListTemplates(ctx context.Context) ([]*ContentTemplate, error)
	GetTemplate(ctx context.Context, id string) (*ContentTemplate, error)
	Render(ctx context.Context, req *RenderRequest) (string, string, error)
	CreateVariant(ctx context.Context, req *CreateVariantRequest) (*Variant, error)
	ListVariants(ctx context.Context, templateID string) ([]*Variant, error)
	SelectVariant(ctx context.Context, templateID string) (*Variant, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) CreateTemplate(ctx context.Context, req *CreateTemplateRequest) (*ContentTemplate, error) {
	t := &ContentTemplate{
		Name:        req.Name,
		Channel:     req.Channel,
		Subject:     req.Subject,
		Body:        req.Body,
		Variables:   req.Variables,
		PreviewText: req.PreviewText,
	}
	if err := s.repo.CreateTemplate(ctx, t); err != nil {
		return nil, err
	}
	return t, nil
}

func (s *service) ListTemplates(ctx context.Context) ([]*ContentTemplate, error) {
	return s.repo.ListTemplates(ctx)
}

func (s *service) GetTemplate(ctx context.Context, id string) (*ContentTemplate, error) {
	t, err := s.repo.GetTemplateByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if t == nil {
		return nil, ErrNotFound
	}
	return t, nil
}

func (s *service) Render(ctx context.Context, req *RenderRequest) (string, string, error) {
	t, err := s.repo.GetTemplateByID(ctx, req.TemplateID)
	if err != nil {
		return "", "", err
	}
	if t == nil {
		return "", "", ErrNotFound
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
	if err := subjectTmpl.Execute(&subjectBuf, req.Variables); err != nil {
		return "", "", err
	}

	var bodyBuf bytes.Buffer
	if err := bodyTmpl.Execute(&bodyBuf, req.Variables); err != nil {
		return "", "", err
	}

	return subjectBuf.String(), bodyBuf.String(), nil
}

func (s *service) CreateVariant(ctx context.Context, req *CreateVariantRequest) (*Variant, error) {
	t, err := s.repo.GetTemplateByID(ctx, req.TemplateID)
	if err != nil {
		return nil, err
	}
	if t == nil {
		return nil, ErrNotFound
	}

	v := &Variant{
		TemplateID:        req.TemplateID,
		Name:              req.Name,
		Modifications:     req.Modifications,
		TrafficPercentage: req.TrafficPercentage,
	}
	if v.TrafficPercentage < 0 {
		v.TrafficPercentage = 0
	}
	if v.TrafficPercentage > 100 {
		v.TrafficPercentage = 100
	}

	if err := s.repo.CreateVariant(ctx, v); err != nil {
		return nil, err
	}
	return v, nil
}

func (s *service) ListVariants(ctx context.Context, templateID string) ([]*Variant, error) {
	return s.repo.ListVariants(ctx, templateID)
}

func (s *service) SelectVariant(ctx context.Context, templateID string) (*Variant, error) {
	variants, err := s.repo.ListVariants(ctx, templateID)
	if err != nil {
		return nil, err
	}
	if len(variants) == 0 {
		return nil, nil
	}

	totalTraffic := 0
	for _, v := range variants {
		totalTraffic += v.TrafficPercentage
	}
	if totalTraffic <= 0 {
		return nil, nil
	}

	roll := rand.Intn(100) + 1
	cumulative := 0
	for _, v := range variants {
		if v.TrafficPercentage > 0 {
			cumulative += v.TrafficPercentage
			if roll <= cumulative {
				return v, nil
			}
		}
	}
	return variants[len(variants)-1], nil
}

func extractVariables(body string) []string {
	var vars []string
	seen := make(map[string]bool)
	parts := strings.Split(body, "{{")
	for _, part := range parts[1:] {
		end := strings.Index(part, "}}")
		if end < 0 {
			continue
		}
		v := strings.TrimSpace(part[:end])
		if !seen[v] && v != "" {
			seen[v] = true
			vars = append(vars, v)
		}
	}
	return vars
}
