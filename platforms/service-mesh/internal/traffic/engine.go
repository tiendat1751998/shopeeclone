package traffic

import (
	"context"
	"math/rand"
	"strings"

	"github.com/google/uuid"
)

type Engine struct {
	repo Repository
}

func NewEngine(repo Repository) *Engine {
	return &Engine{repo: repo}
}

func (e *Engine) CreateRule(ctx context.Context, rule *TrafficRule) (*TrafficRule, error) {
	if rule.ID == "" {
		rule.ID = uuid.New().String()
	}
	if err := e.repo.CreateRule(ctx, rule); err != nil {
		return nil, err
	}
	return rule, nil
}

func (e *Engine) ListRules(ctx context.Context) ([]*TrafficRule, error) {
	return e.repo.ListRules(ctx)
}

func (e *Engine) EvaluateRoute(ctx context.Context, source, destination string, headers map[string]string, path string, method string) (*TrafficRule, error) {
	rules, err := e.repo.ListRules(ctx)
	if err != nil {
		return nil, err
	}

	var matchingRules []*TrafficRule
	for _, rule := range rules {
		if rule.SourceService != "" && rule.SourceService != source {
			continue
		}
		if rule.DestinationService != "" && rule.DestinationService != destination {
			continue
		}

		if !matchConditions(rule.MatchConditions, headers, path, method) {
			continue
		}
		matchingRules = append(matchingRules, rule)
	}

	if len(matchingRules) == 0 {
		return nil, nil
	}

	totalWeight := 0
	for _, r := range matchingRules {
		totalWeight += r.Weight
	}

	if totalWeight == 0 {
		return matchingRules[0], nil
	}

	roll := rand.Intn(totalWeight)
	cumulative := 0
	for _, r := range matchingRules {
		cumulative += r.Weight
		if roll < cumulative {
			return r, nil
		}
	}

	return matchingRules[0], nil
}

func matchConditions(m MatchCondition, headers map[string]string, path, method string) bool {
	if len(m.Methods) > 0 {
		methodMatch := false
		for _, mtd := range m.Methods {
			if strings.EqualFold(mtd, method) {
				methodMatch = true
				break
			}
		}
		if !methodMatch {
			return false
		}
	}

	if m.PathPrefix != "" && !strings.HasPrefix(path, m.PathPrefix) {
		return false
	}

	for k, v := range m.Headers {
		if headers[k] != v {
			return false
		}
	}

	return true
}
