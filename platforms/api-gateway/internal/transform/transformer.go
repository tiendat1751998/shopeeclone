package transform

import (
	"fmt"
	"strings"
)

type Transformer struct {
	repo Repository
}

func NewTransformer(repo Repository) *Transformer {
	return &Transformer{repo: repo}
}

func (t *Transformer) Apply(ruleID string, req *TransformRequest) (*TransformResponse, error) {
	rule, err := t.repo.Get(ruleID)
	if err != nil {
		return nil, err
	}
	if rule == nil {
		return nil, fmt.Errorf("rule not found: %s", ruleID)
	}

	resp := &TransformResponse{
		Headers: copyMap(req.Headers),
		Query:   copyMap(req.Query),
		Path:    req.Path,
		Body:    req.Body,
	}

	for _, action := range rule.Actions {
		switch action.Type {
		case "set_header":
			resp.Headers[action.Name] = action.Value
		case "remove_header":
			delete(resp.Headers, action.Name)
		case "set_query":
			resp.Query[action.Name] = action.Value
		case "rewrite_path":
			resp.Path = rewritePath(action.Value, req.Path)
		case "transform_body":
			resp.Body = action.Value
		}
	}

	return resp, nil
}

func (t *Transformer) CreateRule(rule *Rule) error {
	if rule.ID == "" {
		return fmt.Errorf("id is required")
	}
	if len(rule.Actions) == 0 {
		return fmt.Errorf("at least one action is required")
	}
	return t.repo.Store(rule)
}

func copyMap(src map[string]string) map[string]string {
	dst := make(map[string]string, len(src))
	for k, v := range src {
		dst[k] = v
	}
	return dst
}

func rewritePath(pattern, originalPath string) string {
	if strings.Contains(pattern, "$path") {
		return strings.ReplaceAll(pattern, "$path", originalPath)
	}
	return pattern
}

type ResponseTransformRule struct {
	AddCORSHeaders bool   `json:"add_cors_headers"`
	WrapInEnvelope bool   `json:"wrap_in_envelope"`
	EnvelopeKey    string `json:"envelope_key"`
}

func NewDefaultResponseTransform() *ResponseTransformRule {
	return &ResponseTransformRule{
		AddCORSHeaders: true,
		WrapInEnvelope: false,
		EnvelopeKey:    "data",
	}
}

func ApplyResponseTransform(rule *ResponseTransformRule, originalBody interface{}) interface{} {
	if rule.WrapInEnvelope {
		envelopeKey := rule.EnvelopeKey
		if envelopeKey == "" {
			envelopeKey = "data"
		}
		return map[string]interface{}{
			envelopeKey: originalBody,
		}
	}
	return originalBody
}

type Composer struct {
	transformers []*Transformer
}

func NewComposer(transformers []*Transformer) *Composer {
	return &Composer{transformers: transformers}
}

func (c *Composer) Apply(ruleIDs []string, req *TransformRequest) (*TransformResponse, error) {
	current := req
	for i, ruleID := range ruleIDs {
		for _, t := range c.transformers {
			resp, err := t.Apply(ruleID, current)
			if err != nil {
				return nil, fmt.Errorf("transform %d (%s) failed: %w", i, ruleID, err)
			}
			if resp != nil {
				current = &TransformRequest{
					Path:    resp.Path,
					Method:  current.Method,
					Headers: resp.Headers,
					Query:   resp.Query,
					Body:    resp.Body,
				}
			}
		}
	}

	return &TransformResponse{
		Headers: current.Headers,
		Query:   current.Query,
		Path:    current.Path,
		Body:    current.Body,
	}, nil
}
