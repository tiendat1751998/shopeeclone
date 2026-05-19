package decisionlog

import (
	"context"
	"math"
)

type Logger struct {
	repo Repository
}

func NewLogger(repo Repository) *Logger {
	return &Logger{repo: repo}
}

func (l *Logger) LogDecision(ctx context.Context, log *DecisionLog) error {
	return l.repo.Save(ctx, log)
}

func (l *Logger) GetDecisionHistory(ctx context.Context) ([]*DecisionLog, error) {
	return l.repo.List(ctx)
}

func (l *Logger) GetDecisionByID(ctx context.Context, id string) (*DecisionLog, error) {
	return l.repo.Get(ctx, id)
}

func (l *Logger) GetStats(ctx context.Context) (*DecisionStats, error) {
	logs, err := l.repo.List(ctx)
	if err != nil {
		return nil, err
	}

	stats := &DecisionStats{
		TotalDecisions: len(logs),
	}

	var totalScore float64
	for _, log := range logs {
		switch log.Decision {
		case DecisionAllow:
			stats.AllowCount++
		case DecisionBlock:
			stats.BlockCount++
		case DecisionReview:
			stats.ReviewCount++
		}
		totalScore += log.RiskScore
	}

	if stats.TotalDecisions > 0 {
		stats.AvgRiskScore = math.Round((totalScore/float64(stats.TotalDecisions))*100) / 100
	}

	return stats, nil
}
