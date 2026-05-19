package webpush

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

type Subscription struct {
	Endpoint string
	P256dh   []byte
	Auth     []byte
}

type Options struct {
	PrivateKey *ecdsa.PrivateKey // server's VAPID signing key
	Subject    string            // "mailto:..." or "https://..." contact
	TTL        int               // seconds the push service should retain the message
	Urgency    string            // "very-low", "low", "normal", or "high"
	Topic      string            // an arbitrary string to identify the message topic
}

func SendNotification(ctx context.Context, payload []byte, sub Subscription, opts Options) error {
	body, err := encrypt(payload, sub.P256dh, sub.Auth)
	if err != nil {
		return fmt.Errorf("encrypting payload: %w", err)
	}

	authHeader, err := vapidAuthorization(sub.Endpoint, opts.PrivateKey, opts.Subject)
	if err != nil {
		return fmt.Errorf("building VAPID authorization: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, sub.Endpoint, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("creating request: %w", err)
	}

	req.Header.Set("Content-Type", "application/octet-stream")
	req.Header.Set("Content-Encoding", "aes128gcm")
	req.Header.Set("Content-Length", strconv.Itoa(len(body)))
	req.Header.Set("Authorization", authHeader)
	req.Header.Set("TTL", strconv.Itoa(opts.TTL))

	if opts.Urgency != "" {
		req.Header.Set("Urgency", opts.Urgency)
	}

	if opts.Topic != "" {
		req.Header.Set("Topic", opts.Topic)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("sending request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("sending request: %d %s", resp.StatusCode, string(body))
	}

	return nil
}
