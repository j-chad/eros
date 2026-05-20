package service

import (
	"backend/internal/config"
	"backend/internal/logging"
	"backend/internal/models"
	"backend/internal/repository"
	"backend/internal/webpush"
	"context"
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
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

func (p *PushService) Register(ctx context.Context, sub models.PushSubscription) error {
	logger := logging.FromContext(ctx)

	deviceID := sub.DeviceID
	if deviceID == "" {
		return fmt.Errorf("missing device ID in push subscription")
	}

	logger.DebugContext(ctx, "registering push subscription", "device_id", deviceID)

	if err := p.validateSubscription(sub); err != nil {
		logger.WarnContext(ctx, "invalid push subscription", "error", err)
		return err
	}

	err := p.repo.CreatePushSubscription(ctx, deviceID, sub)
	if err != nil {
		logger.ErrorContext(ctx, "failed to create push subscription", "error", err)
		return err
	}

	return nil
}

func (p *PushService) Unregister(ctx context.Context, deviceID string) error {
	logger := logging.FromContext(ctx)

	logger.DebugContext(ctx, "removing push subscription", "device_id", deviceID)

	err := p.repo.DeletePushSubscriptions(ctx, deviceID)
	if err != nil {
		logger.ErrorContext(ctx, "failed to delete push subscriptions", "error", err)
		return err
	}

	return nil
}

func (p *PushService) GetVAPIDPublicKey() (string, error) {
	if p.config.VAPID.PublicKey == "" {
		return "", ErrPushNotConfigured
	}

	return p.config.VAPID.PublicKey, nil
}

func (p *PushService) SendMessage(ctx context.Context, request models.PushRequest) (models.PushSendResult, error) {
	logger := logging.FromContext(ctx)
	logger.DebugContext(ctx, "sending push", "request", request)

	result := models.PushSendResult{}

	subscriptions, err := p.repo.GetPushSubscriptions(ctx)
	if err != nil {
		logger.ErrorContext(ctx, "failed to get push subscriptions", "error", err)
		return result, fmt.Errorf("failed to get push subscriptions: %w", err)
	}

	payload, err := json.Marshal(request.Message)
	if err != nil {
		logger.ErrorContext(ctx, "failed to marshal push message", "error", err)
		return result, fmt.Errorf("failed to marshal push message: %w", err)
	}

	privateKey, err := p.decodePrivateKey()
	if err != nil {
		logger.ErrorContext(ctx, "failed to decode private key", "error", err)
		return result, fmt.Errorf("failed to decode private key: %w", err)
	}

	options := webpush.Options{
		PrivateKey: privateKey,
		Subject:    p.config.VAPID.Subject,
		TTL:        int(request.TTL.Seconds()),
		Urgency:    string(request.Urgency),
		Topic:      request.Topic,
	}

	for _, subscription := range subscriptions {
		p256dh, err := decodeSubscriptionKey(subscription.Keys.P256dh)
		if err != nil {
			logger.ErrorContext(ctx, "failed to decode subscription key", "error", err)
			_ = p.Unregister(ctx, subscription.DeviceID)
			result.Cleaned += 1
			continue
		}

		auth, err := decodeSubscriptionKey(subscription.Keys.Auth)
		if err != nil {
			logger.ErrorContext(ctx, "failed to decode subscription key", "error", err)
			_ = p.Unregister(ctx, subscription.DeviceID)
			result.Cleaned += 1
			continue
		}

		err = webpush.SendNotification(ctx, payload, webpush.Subscription{
			Endpoint: subscription.Endpoint,
			P256dh:   p256dh,
			Auth:     auth,
		}, options)

		if err != nil {
			logger.ErrorContext(ctx, "failed to send push notification", "error", err)
			result.Failed += 1
			continue
		}

		result.Sent += 1
	}

	return result, nil
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

func (p *PushService) validateKeys(keys models.PushKeys) error {
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

func (p *PushService) decodePrivateKey() (*ecdsa.PrivateKey, error) {
	privateKeyString := p.config.VAPID.PrivateKey
	if privateKeyString == "" {
		return nil, fmt.Errorf("missing private key")
	}
	privateKeyBytes, err := base64.RawURLEncoding.DecodeString(privateKeyString)
	if err != nil {
		return nil, fmt.Errorf("failed to decode private key: %w", err)
	}

	privateKey, err := x509.ParseECPrivateKey(privateKeyBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %w", err)
	}

	return privateKey, nil
}

func decodeSubscriptionKey(key string) ([]byte, error) {
	return base64.RawURLEncoding.DecodeString(key)
}
