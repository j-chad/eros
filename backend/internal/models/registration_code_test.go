package models

import (
	"backend/internal/testutil"
	"testing"
	"time"
)

func TestRegistrationCode_IsExpired_Past(t *testing.T) {
	rc := RegistrationCode{ExpiresAt: time.Now().Add(-time.Minute)}
	testutil.True(t, rc.IsExpired(), "code with past expiry should be expired")
}

func TestRegistrationCode_IsExpired_Future(t *testing.T) {
	rc := RegistrationCode{ExpiresAt: time.Now().Add(time.Hour)}
	testutil.False(t, rc.IsExpired(), "code with future expiry should not be expired")
}
