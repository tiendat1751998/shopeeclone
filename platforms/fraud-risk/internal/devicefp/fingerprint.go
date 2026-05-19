package devicefp

import (
	"context"
	"crypto/sha256"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) IdentifyDevice(ctx context.Context, fp *Fingerprint) (*DeviceProfile, bool, error) {
	hash := s.calculateHash(fp)
	fp.Hash = hash

	existing, err := s.repo.GetByHash(ctx, hash)
	if err == nil {
		now := time.Now().UTC()
		existing.LastSeen = now
		if !contains(existing.UserAgents, fp.UserAgent) {
			existing.UserAgents = append(existing.UserAgents, fp.UserAgent)
		}
		if err := s.repo.Save(ctx, existing, hash); err != nil {
			return nil, false, err
		}
		return existing, false, nil
	}

	deviceID := uuid.New().String()
	fp.DeviceID = deviceID
	now := time.Now().UTC()
	profile := &DeviceProfile{
		DeviceID:     deviceID,
		UserAgents:   []string{fp.UserAgent},
		IPs:          []string{},
		FirstSeen:    now,
		LastSeen:     now,
		RiskScore:    0,
		IsSuspicious: false,
	}

	if err := s.repo.Save(ctx, profile, hash); err != nil {
		return nil, false, err
	}

	return profile, true, nil
}

func (s *Service) GetDeviceHistory(ctx context.Context, deviceID string) (*DeviceProfile, error) {
	return s.repo.Get(ctx, deviceID)
}

func (s *Service) MarkSuspicious(ctx context.Context, deviceID string) (*DeviceProfile, error) {
	profile, err := s.repo.Get(ctx, deviceID)
	if err != nil {
		return nil, err
	}
	profile.IsSuspicious = true
	profile.RiskScore = 100
	if err := s.repo.Save(ctx, profile, ""); err != nil {
		return nil, err
	}
	return profile, nil
}

func (s *Service) calculateHash(fp *Fingerprint) string {
	data := strings.Join([]string{
		fp.UserAgent,
		fmt.Sprintf("%dx%d", fp.ScreenWidth, fp.ScreenHeight),
		fmt.Sprintf("%d", fp.ColorDepth),
		fp.Platform,
		fp.Language,
		fp.Timezone,
	}, "|")
	hash := sha256.Sum256([]byte(data))
	return fmt.Sprintf("%x", hash)
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
