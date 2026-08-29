package api

import (
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
)

func TestClient_GetProfiles(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	tests := []struct {
		name           string
		mockResponse   string
		mockStatusCode int
		wantErr        bool
		wantLen        int
		validate       func(*testing.T, []Profile)
	}{
		{
			name: "successful get profiles",
			mockResponse: `[{
				"id": 123,
				"server_storage": true,
				"subdomain": "mysubdomain",
				"has_premium": true,
				"has_phone": false,
				"onboarding_state": 3,
				"date_subscribed": "2025-01-01T00:00:00Z",
				"avatar": "https://example.com/avatar.png",
				"emails_blocked": 100,
				"emails_forwarded": 500,
				"emails_replied": 10,
				"level_one_trackers_blocked": 50,
				"remove_level_one_email_trackers": true,
				"at_mask_limit": false,
				"bounce_status": [false, ""]
			}]`,
			mockStatusCode: http.StatusOK,
			wantErr:        false,
			wantLen:        1,
			validate: func(t *testing.T, profiles []Profile) {
				if profiles[0].ID != 123 {
					t.Errorf("Profile ID = %d, want 123", profiles[0].ID)
				}
				if !profiles[0].HasPremium {
					t.Error("Profile HasPremium = false, want true")
				}
				if profiles[0].Subdomain == nil || *profiles[0].Subdomain != "mysubdomain" {
					t.Error("Profile Subdomain not set correctly")
				}
				if profiles[0].EmailsForwarded != 500 {
					t.Errorf("Profile EmailsForwarded = %d, want 500", profiles[0].EmailsForwarded)
				}
			},
		},
		{
			name:           "empty profiles list",
			mockResponse:   `[]`,
			mockStatusCode: http.StatusOK,
			wantErr:        false,
			wantLen:        0,
			validate:       nil,
		},
		{
			name: "multiple profiles",
			mockResponse: `[
				{"id": 1, "server_storage": true, "has_premium": false, "has_phone": false, "onboarding_state": 0, "avatar": "", "emails_blocked": 0, "emails_forwarded": 0, "emails_replied": 0, "level_one_trackers_blocked": 0, "remove_level_one_email_trackers": false, "at_mask_limit": false, "bounce_status": [false, ""]},
				{"id": 2, "server_storage": true, "has_premium": true, "has_phone": true, "onboarding_state": 3, "avatar": "", "emails_blocked": 10, "emails_forwarded": 100, "emails_replied": 5, "level_one_trackers_blocked": 20, "remove_level_one_email_trackers": true, "at_mask_limit": false, "bounce_status": [false, ""]}
			]`,
			mockStatusCode: http.StatusOK,
			wantErr:        false,
			wantLen:        2,
			validate: func(t *testing.T, profiles []Profile) {
				if profiles[0].ID != 1 {
					t.Errorf("First Profile ID = %d, want 1", profiles[0].ID)
				}
				if profiles[1].ID != 2 {
					t.Errorf("Second Profile ID = %d, want 2", profiles[1].ID)
				}
			},
		},
		{
			name:           "unauthorized error",
			mockResponse:   `{"detail": "Authentication credentials were not provided."}`,
			mockStatusCode: http.StatusUnauthorized,
			wantErr:        true,
			wantLen:        0,
			validate:       nil,
		},
		{
			name:           "server error",
			mockResponse:   `{"detail": "Internal server error"}`,
			mockStatusCode: http.StatusInternalServerError,
			wantErr:        true,
			wantLen:        0,
			validate:       nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			httpmock.Reset()

			httpmock.RegisterResponder(
				http.MethodGet,
				DefaultBaseURL+profilesPath,
				httpmock.NewStringResponder(tt.mockStatusCode, tt.mockResponse),
			)

			client := NewClient("test")
			profiles, err := client.GetProfiles(t.Context())

			if (err != nil) != tt.wantErr {
				t.Errorf("GetProfiles() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if len(profiles) != tt.wantLen {
					t.Errorf("GetProfiles() returned %d profiles, want %d", len(profiles), tt.wantLen)
				}
				if tt.validate != nil {
					tt.validate(t, profiles)
				}
			}
		})
	}
}

func TestClient_GetProfiles_InvalidJSON(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(
		http.MethodGet,
		DefaultBaseURL+profilesPath,
		httpmock.NewStringResponder(http.StatusOK, `{invalid json`),
	)

	client := NewClient("test")
	_, err := client.GetProfiles(t.Context())

	if err == nil {
		t.Error("GetProfiles() expected error for invalid JSON, got nil")
	}
}
