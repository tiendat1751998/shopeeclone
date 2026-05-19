package audience

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

var ErrNotFound = errors.New("audience: not found")

type Service interface {
	CreateSegment(ctx context.Context, req *CreateSegmentRequest) (*Segment, error)
	ListSegments(ctx context.Context) ([]*Segment, error)
	EvaluateUser(ctx context.Context, segmentID string, userID string) (bool, error)
	EstimateSegmentSize(ctx context.Context, segmentID string) (int, error)
	GetSegmentUsers(ctx context.Context, segmentID string) ([]*UserProfile, error)
	AddToSegment(ctx context.Context, userID string, segmentID string) error
	CreateUser(ctx context.Context, u *UserProfile) (*UserProfile, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) CreateSegment(ctx context.Context, req *CreateSegmentRequest) (*Segment, error) {
	seg := &Segment{
		Name:     req.Name,
		Criteria: req.Criteria,
	}
	if err := s.repo.CreateSegment(ctx, seg); err != nil {
		return nil, err
	}
	return seg, nil
}

func (s *service) ListSegments(ctx context.Context) ([]*Segment, error) {
	return s.repo.ListSegments(ctx)
}

func (s *service) EvaluateUser(ctx context.Context, segmentID string, userID string) (bool, error) {
	seg, err := s.repo.GetSegmentByID(ctx, segmentID)
	if err != nil {
		return false, err
	}
	if seg == nil {
		return false, ErrNotFound
	}

	user, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		return false, err
	}
	if user == nil {
		return false, nil
	}

	c := seg.Criteria
	if c.AgeRange != nil {
		ageStr, ok := user.Attributes["age"]
		if !ok {
			return false, nil
		}
		age, err := strconv.Atoi(ageStr)
		if err != nil {
			return false, nil
		}
		if age < c.AgeRange.Min || age > c.AgeRange.Max {
			return false, nil
		}
	}

	if c.Gender != nil {
		gender, ok := user.Attributes["gender"]
		if !ok || !strings.EqualFold(gender, *c.Gender) {
			return false, nil
		}
	}

	if c.Location != nil {
		loc, ok := user.Attributes["location"]
		if !ok || !strings.EqualFold(loc, *c.Location) {
			return false, nil
		}
	}

	if c.LastActive != nil {
		lastActiveStr, ok := user.Attributes["last_active"]
		if !ok {
			return false, nil
		}
		lastActive, err := time.Parse(time.RFC3339, lastActiveStr)
		if err != nil {
			return false, nil
		}
		dur, err := time.ParseDuration(*c.LastActive)
		if err != nil {
			return false, nil
		}
		if time.Since(lastActive) > dur {
			return false, nil
		}
	}

	if len(c.Tags) > 0 {
		tagSet := make(map[string]bool)
		for _, t := range user.Tags {
			tagSet[strings.ToLower(t)] = true
		}
		for _, t := range c.Tags {
			if !tagSet[strings.ToLower(t)] {
				return false, nil
			}
		}
	}

	return true, nil
}

func (s *service) EstimateSegmentSize(ctx context.Context, segmentID string) (int, error) {
	seg, err := s.repo.GetSegmentByID(ctx, segmentID)
	if err != nil {
		return 0, err
	}
	if seg == nil {
		return 0, ErrNotFound
	}

	users, err := s.repo.ListUsers(ctx)
	if err != nil {
		return 0, err
	}

	count := 0
	for _, u := range users {
		match, err := s.EvaluateUser(ctx, segmentID, u.ID)
		if err != nil {
			continue
		}
		if match {
			count++
		}
	}
	return count, nil
}

func (s *service) GetSegmentUsers(ctx context.Context, segmentID string) ([]*UserProfile, error) {
	users, err := s.repo.ListUsers(ctx)
	if err != nil {
		return nil, err
	}
	var result []*UserProfile
	for _, u := range users {
		match, err := s.EvaluateUser(ctx, segmentID, u.ID)
		if err != nil {
			continue
		}
		if match {
			result = append(result, u)
		}
	}
	return result, nil
}

func (s *service) AddToSegment(ctx context.Context, userID string, segmentID string) error {
	return s.repo.AddToSegment(ctx, userID, segmentID)
}

func (s *service) CreateUser(ctx context.Context, u *UserProfile) (*UserProfile, error) {
	if err := s.repo.CreateUser(ctx, u); err != nil {
		return nil, err
	}
	return u, nil
}

func strPtr(s string) *string { return &s }

type UserProfilePredicate func(u *UserProfile) bool

func MatchesCriteria(c Criteria) UserProfilePredicate {
	return func(u *UserProfile) bool {
		if c.AgeRange != nil {
			ageStr := u.Attributes["age"]
			age, err := strconv.Atoi(ageStr)
			if err != nil || age < c.AgeRange.Min || age > c.AgeRange.Max {
				return false
			}
		}
		if c.Gender != nil {
			g, ok := u.Attributes["gender"]
			if !ok || !strings.EqualFold(g, *c.Gender) {
				return false
			}
		}
		if c.Location != nil {
			l, ok := u.Attributes["location"]
			if !ok || !strings.EqualFold(l, *c.Location) {
				return false
			}
		}
		if len(c.Tags) > 0 {
			tagSet := make(map[string]bool)
			for _, t := range u.Tags {
				tagSet[strings.ToLower(t)] = true
			}
			for _, t := range c.Tags {
				if !tagSet[strings.ToLower(t)] {
					return false
				}
			}
		}
		return true
	}
}

func (s *EvaluateLogger) Debug(msg string, args ...interface{}) {
	fmt.Printf("[EVAL] "+msg+"\n", args...)
}

type EvaluateLogger struct{}

var DefaultLogger = &EvaluateLogger{}
