package api

import (
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

func TestGetRealPhone(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	client := NewClient("test-token")

	t.Run("successful get real phones", func(t *testing.T) {
		responder := httpmock.NewJsonResponderOrPanic(200, []RealPhone{
			{
				ID:                   12040,
				Number:               "+18001234567",
				VerificationSentDate: stringPtr("2026-01-01T00:00:00Z"),
				Verified:             true,
				VerifiedDate:         stringPtr("2026-01-01T00:10:00Z"),
				CountryCode:          "US",
			},
		})
		httpmock.RegisterResponder("GET", DefaultBaseURL+APIBasePath+"realphone/", responder)

		phones, err := client.GetRealPhone()

		assert.NoError(t, err)
		assert.Len(t, phones, 1)
		assert.Equal(t, 12040, phones[0].ID)
		assert.Equal(t, "+18001234567", phones[0].Number)
		assert.True(t, phones[0].Verified)
		assert.Equal(t, "US", phones[0].CountryCode)
	})

	t.Run("error response", func(t *testing.T) {
		responder := httpmock.NewStringResponder(400, `{"error": "bad request"}`)
		httpmock.RegisterResponder("GET", DefaultBaseURL+APIBasePath+"realphone/", responder)

		phones, err := client.GetRealPhone()

		assert.Error(t, err)
		assert.Nil(t, phones)
	})
}

func TestRegisterRealPhone(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	client := NewClient("test-token")

	t.Run("successful register", func(t *testing.T) {
		responder := httpmock.NewJsonResponderOrPanic(201, RealPhone{
			ID:                   12040,
			Number:               "+18001234567",
			VerificationSentDate: stringPtr("2026-01-01T00:00:00Z"),
			Verified:             false,
			VerifiedDate:         nil,
			CountryCode:          "US",
		})
		httpmock.RegisterResponder("POST", DefaultBaseURL+APIBasePath+"realphone/", responder)

		req := RegisterRealPhoneRequest{
			Number: "+18001234567",
		}
		phone, err := client.RegisterRealPhone(req)

		assert.NoError(t, err)
		assert.NotNil(t, phone)
		assert.Equal(t, 12040, phone.ID)
		assert.Equal(t, "+18001234567", phone.Number)
		assert.False(t, phone.Verified)
		assert.Nil(t, phone.VerifiedDate)
	})

	t.Run("error response", func(t *testing.T) {
		responder := httpmock.NewStringResponder(400, `{"error": "invalid phone number"}`)
		httpmock.RegisterResponder("POST", DefaultBaseURL+APIBasePath+"realphone/", responder)

		req := RegisterRealPhoneRequest{
			Number: "invalid",
		}
		phone, err := client.RegisterRealPhone(req)

		assert.Error(t, err)
		assert.Nil(t, phone)
	})
}

func TestVerifyRealPhone(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	client := NewClient("test-token")

	t.Run("successful verify", func(t *testing.T) {
		responder := httpmock.NewJsonResponderOrPanic(200, RealPhone{
			ID:                   12040,
			Number:               "+18001234567",
			VerificationSentDate: stringPtr("2026-01-01T00:00:00Z"),
			Verified:             true,
			VerifiedDate:         stringPtr("2026-01-01T00:10:00Z"),
			CountryCode:          "US",
		})
		httpmock.RegisterResponder("PATCH", DefaultBaseURL+APIBasePath+"realphone/12040/", responder)

		req := VerifyRealPhoneRequest{
			Number:           "+18001234567",
			VerificationCode: "123456",
		}
		phone, err := client.VerifyRealPhone(12040, req)

		assert.NoError(t, err)
		assert.NotNil(t, phone)
		assert.Equal(t, 12040, phone.ID)
		assert.True(t, phone.Verified)
		assert.NotNil(t, phone.VerifiedDate)
	})

	t.Run("error response", func(t *testing.T) {
		responder := httpmock.NewStringResponder(400, `{"error": "invalid verification code"}`)
		httpmock.RegisterResponder("PATCH", DefaultBaseURL+APIBasePath+"realphone/12040/", responder)

		req := VerifyRealPhoneRequest{
			Number:           "+18001234567",
			VerificationCode: "000000",
		}
		phone, err := client.VerifyRealPhone(12040, req)

		assert.Error(t, err)
		assert.Nil(t, phone)
	})
}

func TestDeleteRealPhone(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	client := NewClient("test")

	t.Run("successful delete", func(t *testing.T) {
		responder := httpmock.NewStringResponder(204, "")
		httpmock.RegisterResponder("DELETE", DefaultBaseURL+APIBasePath+"realphone/12040/", responder)

		err := client.DeleteRealPhone(12040)

		assert.NoError(t, err)
	})

	t.Run("error response", func(t *testing.T) {
		responder := httpmock.NewStringResponder(404, `{"error": "not found"}`)
		httpmock.RegisterResponder("DELETE", DefaultBaseURL+APIBasePath+"realphone/99999/", responder)

		err := client.DeleteRealPhone(99999)

		assert.Error(t, err)
	})
}

func stringPtr(s string) *string {
	return &s
}
