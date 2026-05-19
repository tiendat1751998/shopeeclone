package preferences

import (
	"context"
	"time"
)

type Service interface {
	GetPreferences(ctx context.Context, userID string) (*UserPreference, error)
	UpdatePreferences(ctx context.Context, userID string, req *UpdatePreferenceRequest) (*UserPreference, error)
	ShouldSend(ctx context.Context, userID string, channel string, category string) (bool, error)
	AddSuppression(ctx context.Context, userID, email, phone, reason string) error
	IsSuppressed(ctx context.Context, userID, email, phone string) (bool, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) GetPreferences(ctx context.Context, userID string) (*UserPreference, error) {
	return s.repo.Get(ctx, userID)
}

func (s *service) UpdatePreferences(ctx context.Context, userID string, req *UpdatePreferenceRequest) (*UserPreference, error) {
	pref, err := s.repo.Get(ctx, userID)
	if err != nil {
		return nil, err
	}

	if req.ChannelOptIn != nil {
		pref.ChannelOptIn = *req.ChannelOptIn
	}
	if req.Categories != nil {
		pref.Categories = *req.Categories
	}
	if req.QuietHours != nil {
		pref.QuietHours = *req.QuietHours
	}
	if req.EmailDigest != nil {
		pref.EmailDigest = *req.EmailDigest
	}
	if req.PushEnabled != nil {
		pref.PushEnabled = *req.PushEnabled
	}
	if req.SMSPromotions != nil {
		pref.SMSPromotions = *req.SMSPromotions
	}

	if err := s.repo.Upsert(ctx, pref); err != nil {
		return nil, err
	}

	return pref, nil
}

func (s *service) ShouldSend(ctx context.Context, userID string, channel string, category string) (bool, error) {
	pref, err := s.repo.Get(ctx, userID)
	if err != nil {
		return false, err
	}

	suppressed, err := s.repo.IsSuppressed(ctx, userID, "", "")
	if err != nil {
		return false, err
	}
	if suppressed {
		return false, nil
	}

	switch channel {
	case "push":
		if !pref.ChannelOptIn.Push || !pref.PushEnabled {
			return false, nil
		}
	case "email":
		if !pref.ChannelOptIn.Email {
			return false, nil
		}
	case "sms":
		if !pref.ChannelOptIn.SMS {
			return false, nil
		}
	case "inapp":
		if !pref.ChannelOptIn.InApp {
			return false, nil
		}
	}

	if pref.QuietHours.Enabled {
		now := time.Now()
		currentTime := now.Format("15:04")
		if currentTime >= pref.QuietHours.Start && currentTime < pref.QuietHours.End {
			return false, nil
		}
	}

	return true, nil
}

func (s *service) AddSuppression(ctx context.Context, userID, email, phone, reason string) error {
	return s.repo.AddSuppression(ctx, &SuppressionEntry{
		UserID: userID,
		Email:  email,
		Phone:  phone,
		Reason: reason,
	})
}

func (s *service) IsSuppressed(ctx context.Context, userID, email, phone string) (bool, error) {
	return s.repo.IsSuppressed(ctx, userID, email, phone)
}
