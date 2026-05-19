package alerting

import (
	"fmt"
	"time"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateRule(rule *Rule) error {
	return s.repo.CreateRule(rule)
}

func (s *Service) ListRules() ([]*Rule, error) {
	return s.repo.ListRules()
}

func (s *Service) Evaluate(rules []*Rule, metrics []MetricValue) ([]*Alert, error) {
	var alerts []*Alert

	for _, rule := range rules {
		for _, metric := range metrics {
			if metric.Name != rule.MetricName {
				continue
			}

			var fired bool
			switch rule.Operator {
			case ">":
				fired = metric.Value > rule.Threshold
			case "<":
				fired = metric.Value < rule.Threshold
			case ">=":
				fired = metric.Value >= rule.Threshold
			case "<=":
				fired = metric.Value <= rule.Threshold
			}

			if fired {
				alert := &Alert{
					Name:         fmt.Sprintf("%s-alert", rule.Name),
					Condition:    fmt.Sprintf("%s %s %v", rule.MetricName, rule.Operator, rule.Threshold),
					Threshold:    rule.Threshold,
					CurrentValue: metric.Value,
					Status:       AlertFiring,
					Severity:     AlertSeverityWarning,
					Service:      rule.Name,
					TriggeredAt:  time.Now(),
				}

				existingAlerts, _ := s.repo.ListAlerts()
				var cooldownActive bool
				for _, ea := range existingAlerts {
					if ea.Name == alert.Name && ea.Status == AlertFiring {
						if time.Since(ea.TriggeredAt).Seconds() < float64(rule.CooldownSeconds) {
							cooldownActive = true
							break
						}
					}
				}

				if !cooldownActive {
					s.repo.CreateAlert(alert)
					alerts = append(alerts, alert)
				}
			}
		}
	}

	return alerts, nil
}

func (s *Service) ListAlerts() ([]*Alert, error) {
	return s.repo.ListAlerts()
}
