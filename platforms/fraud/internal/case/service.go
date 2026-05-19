package fraudcase

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateCase(ctx context.Context, alertID, userID, title, description string, riskScore float64, priority CasePriority) (*FraudCase, error) {
	c := &FraudCase{
		ID:          uuid.New().String(),
		AlertID:     alertID,
		UserID:      userID,
		Title:       title,
		Description: description,
		Status:      StatusOpen,
		Priority:    priority,
		RiskScore:   riskScore,
		Evidence:    []Evidence{},
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	if err := s.repo.Create(ctx, c); err != nil {
		return nil, err
	}
	return c, nil
}

func (s *Service) GetCase(ctx context.Context, id string) (*FraudCase, error) {
	return s.repo.Get(ctx, id)
}

func (s *Service) ListCases(ctx context.Context, status CaseStatus, priority CasePriority, offset, limit int) ([]*FraudCase, int, error) {
	return s.repo.List(ctx, status, priority, offset, limit)
}

func (s *Service) UpdateCase(ctx context.Context, c *FraudCase) error {
	return s.repo.Update(ctx, c)
}

func (s *Service) AssignInvestigator(ctx context.Context, caseID, investigator string) error {
	c, err := s.repo.Get(ctx, caseID)
	if err != nil {
		return err
	}
	c.Investigator = investigator
	if c.Status == StatusOpen {
		c.Status = StatusInvestigating
	}
	return s.repo.Update(ctx, c)
}

func (s *Service) AddEvidence(ctx context.Context, caseID, eType, description, data, addedBy string) error {
	c, err := s.repo.Get(ctx, caseID)
	if err != nil {
		return err
	}
	evidence := Evidence{
		ID:          uuid.New().String(),
		Type:        eType,
		Description: description,
		Data:        data,
		AddedBy:     addedBy,
		AddedAt:     time.Now(),
	}
	c.Evidence = append(c.Evidence, evidence)
	return s.repo.Update(ctx, c)
}

func (s *Service) UpdateStatus(ctx context.Context, caseID string, status CaseStatus, resolution string) error {
	c, err := s.repo.Get(ctx, caseID)
	if err != nil {
		return err
	}

	if !isValidTransition(c.Status, status) {
		return ErrInvalidTransition
	}

	c.Status = status
	if status == StatusResolved || status == StatusClosed {
		now := time.Now()
		c.ResolvedAt = &now
		c.Resolution = resolution
	}
	return s.repo.Update(ctx, c)
}

func (s *Service) Escalate(ctx context.Context, caseID string, priority CasePriority) error {
	c, err := s.repo.Get(ctx, caseID)
	if err != nil {
		return err
	}
	c.Priority = priority
	c.Status = StatusEscalated
	return s.repo.Update(ctx, c)
}

func isValidTransition(current, next CaseStatus) bool {
	transitions := map[CaseStatus][]CaseStatus{
		StatusOpen:         {StatusInvestigating, StatusClosed},
		StatusInvestigating: {StatusEscalated, StatusResolved, StatusClosed},
		StatusEscalated:   {StatusInvestigating, StatusResolved, StatusClosed},
		StatusResolved:    {StatusClosed},
		StatusClosed:      {},
	}
	allowed, ok := transitions[current]
	if !ok {
		return false
	}
	for _, s := range allowed {
		if s == next {
			return true
		}
	}
	return false
}
