package service

import (
	"backend/internal/config"
	"backend/internal/models"
	"backend/internal/repository"
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"net/url"
	"strings"
)

var ErrSubValidationFailed = errors.New("subscription validation failed")
var ErrPushNotConfigured = errors.New("push not configured")

type PushService struct {
	config config.PushConfig
	repo   repository.Repository
}

func NewPushService(config config.PushConfig, repo repository.Repository) *PushService {
	return &PushService{config: config, repo: repo}
}

func (p *PushService) Register(ctx context.Context, deviceID string, sub models.PushSubscription) error {
	if err := p.validateSubscription(sub); err != nil {
		return err
	}

	return p.repo.CreatePushSubscription(ctx, deviceID, sub)
}

func (p *PushService) Unregister(ctx context.Context, deviceID string) error {
	return p.repo.DeletePushSubscriptions(ctx, deviceID)
}

func (p *PushService) GetVAPIDPublicKey() (string, error) {
	if p.config.VAPID.PublicKey == "" {
		return "", ErrPushNotConfigured
	}

	return p.config.VAPID.PublicKey, nil
}

func (p *PushService) validateSubscription(sub models.PushSubscription) error {
	err := p.validateEndpoint(sub.Endpoint)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrSubValidationFailed, err)
	}

	err = p.validateKeys(sub.Keys)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrSubValidationFailed, err)
	}

	return nil
}

func (p *PushService) validateEndpoint(endpoint string) error {
	u, err := url.Parse(endpoint)
	if err != nil {
		return fmt.Errorf("invalid endpoint URL: %w", err)
	}

	if u.Scheme != "https" {
		return fmt.Errorf("endpoint must use https: %s", endpoint)
	}

	for _, host := range p.config.AllowedHosts {
		if u.Hostname() == host || strings.HasSuffix(u.Hostname(), "."+host) {
			return nil
		}
	}

	return fmt.Errorf("endpoint '%s' is not whitelisted", endpoint)
}

func (p *PushService) validateKeys(keys models.PushSubscriptionKeys) error {
	if len(keys.P256dh) == 0 || len(keys.Auth) == 0 {
		return fmt.Errorf("missing subscription keys")
	}
	if _, err := base64.RawURLEncoding.DecodeString(keys.P256dh); err != nil {
		return fmt.Errorf("invalid p256dh key")
	}
	if _, err := base64.RawURLEncoding.DecodeString(keys.Auth); err != nil {
		return fmt.Errorf("invalid auth key")
	}
	return nil
}
