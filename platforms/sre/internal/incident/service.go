package incident

import "time"

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(title string, severity Severity, service, region, description, assignee string) (*Incident, error) {
	if severity != SeverityCritical && severity != SeverityMajor && severity != SeverityMinor {
		return nil, ErrInvalidSeverity
	}
	inc := &Incident{
		Title:       title,
		Severity:    severity,
		Service:     service,
		Region:      region,
		Description: description,
		Assignee:    assignee,
	}
	if err := s.repo.Create(inc); err != nil {
		return nil, err
	}
	return inc, nil
}

func (s *Service) Acknowledge(id string) (*Incident, error) {
	inc, err := s.repo.Get(id)
	if err != nil {
		return nil, err
	}
	if inc.Status == StatusResolved {
		return nil, ErrInvalidStatus
	}
	inc.Status = StatusTriaging
	if err := s.repo.Update(inc); err != nil {
		return nil, err
	}
	return inc, nil
}

func (s *Service) Resolve(id string) (*Incident, error) {
	inc, err := s.repo.Get(id)
	if err != nil {
		return nil, err
	}
	now := time.Now()
	inc.Status = StatusResolved
	inc.ResolvedAt = &now
	if err := s.repo.Update(inc); err != nil {
		return nil, err
	}
	return inc, nil
}

func (s *Service) Update(id string, updates map[string]interface{}) (*Incident, error) {
	inc, err := s.repo.Get(id)
	if err != nil {
		return nil, err
	}
	if v, ok := updates["status"]; ok {
		s, ok := v.(string)
		if !ok {
			return nil, ErrInvalidStatus
		}
		inc.Status = Status(s)
	}
	if v, ok := updates["severity"]; ok {
		s, ok := v.(string)
		if !ok {
			return nil, ErrInvalidSeverity
		}
		inc.Severity = Severity(s)
	}
	if v, ok := updates["assignee"]; ok {
		s, ok := v.(string)
		if !ok {
			return nil, ErrInvalidAssignee
		}
		inc.Assignee = s
	}
	if v, ok := updates["description"]; ok {
		s, ok := v.(string)
		if !ok {
			return nil, ErrInvalidDescription
		}
		inc.Description = s
	}
	if v, ok := updates["title"]; ok {
		s, ok := v.(string)
		if !ok {
			return nil, ErrInvalidTitle
		}
		inc.Title = s
	}
	if err := s.repo.Update(inc); err != nil {
		return nil, err
	}
	return inc, nil
}

func (s *Service) List(filter Filter) ([]*Incident, error) {
	return s.repo.List(filter)
}
