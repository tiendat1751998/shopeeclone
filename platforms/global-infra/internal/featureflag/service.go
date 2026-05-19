package featureflag

import (
	"context"
	"crypto/sha256"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(ctx context.Context, flag *FeatureFlag) (*FeatureFlag, error) {
	if flag.Name == "" {
		return nil, fmt.Errorf("flag name is required")
	}
	if flag.PercentageRollout < 0 || flag.PercentageRollout > 100 {
		return nil, fmt.Errorf("percentage_rollout must be between 0 and 100")
	}
	if flag.UserSegment == "" {
		flag.UserSegment = SegmentAll
	}
	flag.UpdatedAt = time.Now()
	flag.CreatedAt = time.Now()
	if err := s.repo.Create(ctx, flag); err != nil {
		return nil, err
	}
	return flag, nil
}

func (s *Service) Get(ctx context.Context, name string) (*FeatureFlag, error) {
	return s.repo.Get(ctx, name)
}

func (s *Service) Update(ctx context.Context, flag *FeatureFlag) error {
	existing, err := s.repo.Get(ctx, flag.Name)
	if err != nil {
		return err
	}
	if existing == nil {
		return fmt.Errorf("flag not found: %s", flag.Name)
	}
	if flag.PercentageRollout < 0 || flag.PercentageRollout > 100 {
		return fmt.Errorf("percentage_rollout must be between 0 and 100")
	}
	flag.CreatedAt = existing.CreatedAt
	return s.repo.Update(ctx, flag)
}

func (s *Service) Delete(ctx context.Context, name string) error {
	return s.repo.Delete(ctx, name)
}

func (s *Service) List(ctx context.Context) ([]*FeatureFlag, error) {
	return s.repo.List(ctx)
}

func (s *Service) Evaluate(ctx context.Context, flagName, userID string) (*EvaluateResponse, error) {
	flag, err := s.repo.Get(ctx, flagName)
	if err != nil {
		return nil, err
	}
	if flag == nil {
		return &EvaluateResponse{
			FlagName: flagName,
			UserID:   userID,
			Enabled:  false,
			Reason:   "flag not found",
		}, nil
	}
	if !flag.Enabled {
		return &EvaluateResponse{
			FlagName: flagName,
			UserID:   userID,
			Enabled:  false,
			Reason:   "flag disabled",
		}, nil
	}
	if flag.UserSegment != SegmentAll && !isUserInSegment(userID, flag.UserSegment, uuid.New().String()) {
		return &EvaluateResponse{
			FlagName: flagName,
			UserID:   userID,
			Enabled:  false,
			Reason:   fmt.Sprintf("user not in segment: %s", flag.UserSegment),
		}, nil
	}
	hash := sha256.Sum256([]byte(userID + flag.Name))
	percentage := int(hash[0]) % 100
	if percentage >= flag.PercentageRollout {
		return &EvaluateResponse{
			FlagName: flagName,
			UserID:   userID,
			Enabled:  false,
			Reason:   fmt.Sprintf("percentage rollout miss: %d%%", flag.PercentageRollout),
		}, nil
	}
	return &EvaluateResponse{
		FlagName: flagName,
		UserID:   userID,
		Enabled:  true,
		Reason:   "flag active",
	}, nil
}

func isUserInSegment(userID string, segment Segment, _ string) bool {
	hash := sha256.Sum256([]byte(userID + string(segment)))
	val := int(hash[0]) % 100
	switch segment {
	case SegmentBetaUsers:
		return val < 20
	case SegmentStaff:
		return val < 5
	case SegmentInternal:
		return val < 10
	case SegmentAll:
		return true
	default:
		return false
	}
}

func (s *Service) EvaluateFlag(ctx context.Context, req EvaluateRequest) (*EvaluateResponse, error) {
	if req.FlagName == "" {
		return nil, fmt.Errorf("flag_name is required")
	}
	if req.UserID == "" {
		return nil, fmt.Errorf("user_id is required")
	}
	return s.Evaluate(ctx, req.FlagName, req.UserID)
}


