package pipeline

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
)

type Service interface {
	CreatePipeline(ctx context.Context, name, indexName string, stages []PipelineStage) (*Pipeline, error)
	GetPipeline(ctx context.Context, id string) (*Pipeline, error)
	ListPipelines(ctx context.Context) ([]*Pipeline, error)
	ProcessDocument(ctx context.Context, pipelineID string, doc *Document) (*Document, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) CreatePipeline(ctx context.Context, name, indexName string, stages []PipelineStage) (*Pipeline, error) {
	for _, st := range stages {
		switch st.Processor {
		case ProcessorTokenize, ProcessorAnalyze, ProcessorTransform, ProcessorEnrich:
		default:
			return nil, ErrUnknownProcessor
		}
	}
	p := &Pipeline{
		ID:        uuid.New().String(),
		Name:      name,
		IndexName: indexName,
		Stages:    stages,
		IsActive:  true,
	}
	if err := s.repo.CreatePipeline(ctx, p); err != nil {
		return nil, err
	}
	return p, nil
}

func (s *service) GetPipeline(ctx context.Context, id string) (*Pipeline, error) {
	return s.repo.GetPipeline(ctx, id)
}

func (s *service) ListPipelines(ctx context.Context) ([]*Pipeline, error) {
	return s.repo.ListPipelines(ctx)
}

func (s *service) ProcessDocument(ctx context.Context, pipelineID string, doc *Document) (*Document, error) {
	p, err := s.repo.GetPipeline(ctx, pipelineID)
	if err != nil {
		return nil, err
	}

	result := &Document{
		ID:     doc.ID,
		Fields: make(map[string]interface{}),
	}
	for k, v := range doc.Fields {
		result.Fields[k] = v
	}

	if p.IndexName != "" {
		result.Fields["_index"] = p.IndexName
	}

	for _, stage := range p.Stages {
		switch stage.Processor {
		case ProcessorTokenize:
			tokenize(result, stage)
		case ProcessorAnalyze:
			analyze(result, stage)
		case ProcessorTransform:
			transform(result, stage)
		case ProcessorEnrich:
			enrich(result, stage)
		}
	}

	return result, nil
}

func tokenize(doc *Document, stage PipelineStage) {
	for k, v := range doc.Fields {
		str, ok := v.(string)
		if !ok {
			continue
		}
		tokens := strings.Fields(strings.ToLower(str))
		doc.Fields[k+"_tokens"] = tokens
	}
	if stage.Config != nil {
		if separator, ok := stage.Config["separator"]; ok {
			doc.Fields["token_separator"] = separator
		}
	}
}

func analyze(doc *Document, stage PipelineStage) {
	if tokens, ok := doc.Fields["content_tokens"].([]string); ok {
		stopwords := []string{"the", "a", "an", "is", "are", "was", "were", "in", "on", "at", "to", "for", "of", "and", "or"}
		stopSet := make(map[string]bool, len(stopwords))
		for _, sw := range stopwords {
			stopSet[sw] = true
		}
		filtered := make([]string, 0, len(tokens))
		for _, t := range tokens {
			if !stopSet[t] {
				filtered = append(filtered, t)
			}
		}
		doc.Fields["analyzed_tokens"] = filtered
	}
	if _, ok := doc.Fields["content_tokens"]; !ok {
		for k, v := range doc.Fields {
			if tokens, ok := v.([]string); ok && strings.HasSuffix(k, "_tokens") {
				stopwords := []string{"the", "a", "an", "is", "are", "was", "were", "in", "on", "at", "to", "for", "of", "and", "or"}
				stopSet := make(map[string]bool, len(stopwords))
				for _, sw := range stopwords {
					stopSet[sw] = true
				}
				filtered := make([]string, 0, len(tokens))
				for _, t := range tokens {
					if !stopSet[t] {
						filtered = append(filtered, t)
					}
				}
				doc.Fields["analyzed_tokens"] = filtered
				break
			}
		}
	}
	if stage.Config != nil {
		if lang, ok := stage.Config["language"]; ok {
			doc.Fields["analyze_language"] = lang
		}
	}
}

func transform(doc *Document, stage PipelineStage) {
	for k, v := range doc.Fields {
		switch val := v.(type) {
		case string:
			doc.Fields[k] = strings.TrimSpace(val)
		case []string:
			trimmed := make([]string, len(val))
			for i, s := range val {
				trimmed[i] = strings.TrimSpace(s)
			}
			doc.Fields[k] = trimmed
		}
	}
	if stage.Config != nil {
		if rename, ok := stage.Config["rename"]; ok {
			if renameMap, ok := rename.(map[string]interface{}); ok {
				for oldKey, newKey := range renameMap {
					if newKeyStr, ok := newKey.(string); ok {
						if v, exists := doc.Fields[oldKey]; exists {
							doc.Fields[newKeyStr] = v
							delete(doc.Fields, oldKey)
						}
					}
				}
			}
		}
	}
}

func enrich(doc *Document, stage PipelineStage) {
	if stage.Config != nil {
		for k, v := range stage.Config {
			doc.Fields["enriched_"+k] = v
		}
	}
	doc.Fields["enriched_at"] = fmt.Sprintf("pipeline_processed")
}
