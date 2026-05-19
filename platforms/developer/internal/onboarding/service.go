package onboarding

import (
	"context"
	"sync"
)

type Service struct {
	repo   Repository
	mu     sync.RWMutex
	tasks  map[string]map[string]bool
}

func NewService(repo Repository) *Service {
	return &Service{
		repo:  repo,
		tasks: make(map[string]map[string]bool),
	}
}

func (s *Service) ListTemplates(ctx context.Context) ([]*Template, error) {
	return s.repo.ListTemplates(ctx)
}

func (s *Service) GetTemplate(ctx context.Context, name string) (*Template, error) {
	t, err := s.repo.GetTemplate(ctx, name)
	if err != nil || t == nil {
		return nil, err
	}

	s.mu.RLock()
	completed := s.tasks[name]
	s.mu.RUnlock()

	for i := range t.Tasks {
		if completed != nil && completed[t.Tasks[i].ID] {
			t.Tasks[i].IsCompleted = true
		}
	}
	return t, nil
}

func (s *Service) CompleteTask(ctx context.Context, taskID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// find which template has this task
	templates, err := s.repo.ListTemplates(ctx)
	if err != nil {
		return err
	}

	for _, t := range templates {
		for _, task := range t.Tasks {
			if task.ID == taskID {
				if s.tasks[t.Name] == nil {
					s.tasks[t.Name] = make(map[string]bool)
				}
				s.tasks[t.Name][taskID] = true
				return nil
			}
		}
	}
	return nil
}

func (s *Service) GetProgress(ctx context.Context) (*Progress, error) {
	templates, err := s.repo.ListTemplates(ctx)
	if err != nil {
		return nil, err
	}

	s.mu.RLock()
	completed := s.tasks
	s.mu.RUnlock()

	var allTasks []OnboardingTask
	completedCount := 0

	for _, t := range templates {
		for _, task := range t.Tasks {
			taskCopy := task
			if completed[t.Name] != nil && completed[t.Name][task.ID] {
				taskCopy.IsCompleted = true
				completedCount++
			}
			allTasks = append(allTasks, taskCopy)
		}
	}

	total := len(allTasks)
	var percentage float64
	if total > 0 {
		percentage = float64(completedCount) / float64(total) * 100
	}

	return &Progress{
		TotalTasks:     total,
		CompletedTasks: completedCount,
		Percentage:     percentage,
		Tasks:          allTasks,
	}, nil
}
