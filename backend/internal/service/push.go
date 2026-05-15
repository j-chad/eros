package service

import (
	"backend/internal/config"
	"backend/internal/models"
	"backend/internal/repository"
	"encoding/base64"
	"fmt"
	"net/url"
	"strings"
)

type PushService struct {
	config config.PushConfig
	repo   repository.Repository
}

func NewPushService(config config.PushConfig, repo repository.Repository) *PushService {
	return &PushService{config: config, repo: repo}
}

func (p *PushService) validateSubscription(sub *models.PushSubscription) error {
	err := p.validateEndpoint(sub.Endpoint)
	if err != nil {
		return err
	}

	return p.validateKeys(sub.Keys)
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
