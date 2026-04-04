package models

import (
	"backend/internal/testutil"
	"testing"
	"time"
)

func TestRegisterDeviceRequest_Validate_Empty(t *testing.T) {
	req := RegisterDeviceRequest{RegistrationCode: ""}
	err := req.Validate()
	testutil.NotNil(t, err)
	testutil.Equal(t, err.Code, "BAD_REQUEST")
}

func TestRegisterDeviceRequest_Validate_Valid(t *testing.T) {
	req := RegisterDeviceRequest{
		RegistrationCode: "ABCD-1234-EFGH",
		DeviceInfo:       "iPhone",
	}
	testutil.Nil(t, req.Validate())
}

func TestRegisterDeviceRequest_Validate_NoDeviceInfo(t *testing.T) {
	req := RegisterDeviceRequest{RegistrationCode: "CODE"}
	testutil.Nil(t, req.Validate())
}

func TestDevice_IsExpired_Past(t *testing.T) {
	d := Device{ExpiresAt: time.Now().Add(-time.Hour)}
	testutil.True(t, d.IsExpired(), "device with past expiry should be expired")
}

func TestDevice_IsExpired_Future(t *testing.T) {
	d := Device{ExpiresAt: time.Now().Add(time.Hour)}
	testutil.False(t, d.IsExpired(), "device with future expiry should not be expired")
}
