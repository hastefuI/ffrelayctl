package api

import (
	"encoding/json"
	"testing"
)

func TestProfile_BounceStatusUnmarshal(t *testing.T) {
	tests := []struct {
		name       string
		json       string
		wantErr    bool
		wantPaused bool
		wantType   string
	}{
		{
			name: "bounce_status false with empty string",
			json: `{
				"id": 123,
				"server_storage": true,
				"has_premium": false,
				"has_phone": false,
				"onboarding_state": 0,
				"avatar": "",
				"emails_blocked": 10,
				"emails_forwarded": 50,
				"emails_replied": 5,
				"level_one_trackers_blocked": 20,
				"remove_level_one_email_trackers": true,
				"at_mask_limit": false,
				"bounce_status": [false, ""]
			}`,
			wantErr:    false,
			wantPaused: false,
			wantType:   "",
		},
		{
			name: "bounce_status true with hard type",
			json: `{
				"id": 123,
				"server_storage": true,
				"has_premium": false,
				"has_phone": false,
				"onboarding_state": 0,
				"avatar": "",
				"emails_blocked": 10,
				"emails_forwarded": 50,
				"emails_replied": 5,
				"level_one_trackers_blocked": 20,
				"remove_level_one_email_trackers": true,
				"at_mask_limit": false,
				"bounce_status": [true, "hard"]
			}`,
			wantErr:    false,
			wantPaused: true,
			wantType:   "hard",
		},
		{
			name: "bounce_status true with soft type",
			json: `{
				"id": 123,
				"server_storage": true,
				"has_premium": false,
				"has_phone": false,
				"onboarding_state": 0,
				"avatar": "",
				"emails_blocked": 10,
				"emails_forwarded": 50,
				"emails_replied": 5,
				"level_one_trackers_blocked": 20,
				"remove_level_one_email_trackers": true,
				"at_mask_limit": false,
				"bounce_status": [true, "soft"]
			}`,
			wantErr:    false,
			wantPaused: true,
			wantType:   "soft",
		},
		{
			name: "actual API response from problem statement",
			json: `{
				"id": 123456789,
				"server_storage": true,
				"store_phone_log": true,
				"subdomain": "test",
				"has_premium": true,
				"has_phone": true,
				"has_vpn": false,
				"has_megabundle": false,
				"onboarding_state": 3,
				"onboarding_free_state": 0,
				"date_phone_registered": "2023-10-22T04:30:14.053613Z",
				"date_subscribed": "2025-10-29T07:38:40.706011Z",
				"avatar": "https://profile.accounts.firefox.com/v1/avatar/d",
				"next_email_try": "2026-01-04T00:57:15.765116Z",
				"bounce_status": [false, ""],
				"api_token": "abc-123-123-abc-asdf",
				"emails_blocked": 0,
				"emails_forwarded": 5943,
				"emails_replied": 4,
				"level_one_trackers_blocked": 0,
				"remove_level_one_email_trackers": false,
				"total_masks": 97,
				"at_mask_limit": false,
				"metrics_enabled": true
			}`,
			wantErr:    false,
			wantPaused: false,
			wantType:   "",
		},
		{
			name: "invalid type for first element",
			json: `{
				"id": 123,
				"server_storage": true,
				"has_premium": false,
				"has_phone": false,
				"onboarding_state": 0,
				"avatar": "",
				"emails_blocked": 10,
				"emails_forwarded": 50,
				"emails_replied": 5,
				"level_one_trackers_blocked": 20,
				"remove_level_one_email_trackers": true,
				"at_mask_limit": false,
				"bounce_status": ["not-a-bool", ""]
			}`,
			wantErr: true,
		},
		{
			name: "invalid type for second element",
			json: `{
				"id": 123,
				"server_storage": true,
				"has_premium": false,
				"has_phone": false,
				"onboarding_state": 0,
				"avatar": "",
				"emails_blocked": 10,
				"emails_forwarded": 50,
				"emails_replied": 5,
				"level_one_trackers_blocked": 20,
				"remove_level_one_email_trackers": true,
				"at_mask_limit": false,
				"bounce_status": [false, 123]
			}`,
			wantErr: true,
		},
		{
			name: "array with wrong length (empty)",
			json: `{
				"id": 123,
				"server_storage": true,
				"has_premium": false,
				"has_phone": false,
				"onboarding_state": 0,
				"avatar": "",
				"emails_blocked": 10,
				"emails_forwarded": 50,
				"emails_replied": 5,
				"level_one_trackers_blocked": 20,
				"remove_level_one_email_trackers": true,
				"at_mask_limit": false,
				"bounce_status": []
			}`,
			wantErr: true,
		},
		{
			name: "array with wrong length (too many)",
			json: `{
				"id": 123,
				"server_storage": true,
				"has_premium": false,
				"has_phone": false,
				"onboarding_state": 0,
				"avatar": "",
				"emails_blocked": 10,
				"emails_forwarded": 50,
				"emails_replied": 5,
				"level_one_trackers_blocked": 20,
				"remove_level_one_email_trackers": true,
				"at_mask_limit": false,
				"bounce_status": [false, "soft", "extra"]
			}`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var profile Profile
			err := json.Unmarshal([]byte(tt.json), &profile)
			if (err != nil) != tt.wantErr {
				t.Errorf("json.Unmarshal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil {
				if profile.BounceStatus.Paused != tt.wantPaused {
					t.Errorf("BounceStatus.Paused = %v, want %v", profile.BounceStatus.Paused, tt.wantPaused)
				}
				if profile.BounceStatus.Type != tt.wantType {
					t.Errorf("BounceStatus.Type = %v, want %v", profile.BounceStatus.Type, tt.wantType)
				}
				t.Logf("Successfully unmarshaled profile with bounce_status: paused=%v, type=%q", profile.BounceStatus.Paused, profile.BounceStatus.Type)
			}
		})
	}
}

func TestProfile_BounceStatusMarshal(t *testing.T) {
	profile := Profile{
		ID:                          123,
		ServerStorage:               true,
		HasPremium:                  false,
		HasPhone:                    false,
		OnboardingState:             0,
		EmailsBlocked:               10,
		EmailsForwarded:             50,
		EmailsReplied:               5,
		LevelOneTrackersBlocked:     20,
		RemoveLevelOneEmailTrackers: true,
		AtMaskLimit:                 false,
		BounceStatus: BounceStatus{
			Paused: true,
			Type:   "hard",
		},
	}

	data, err := json.Marshal(profile)
	if err != nil {
		t.Fatalf("json.Marshal() error = %v", err)
	}

	var result map[string]any
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatalf("json.Unmarshal() error = %v", err)
	}

	bounceStatus, ok := result["bounce_status"]
	if !ok {
		t.Fatal("bounce_status field not found in marshaled JSON")
	}

	bounceArray, ok := bounceStatus.([]any)
	if !ok {
		t.Fatalf("bounce_status is not an array, got type %T", bounceStatus)
	}

	if len(bounceArray) != 2 {
		t.Errorf("Expected 2 elements in bounce_status array, got %d", len(bounceArray))
	}

	if paused, ok := bounceArray[0].(bool); !ok || !paused {
		t.Errorf("bounce_status[0] = %v, want true", bounceArray[0])
	}

	if bounceType, ok := bounceArray[1].(string); !ok || bounceType != "hard" {
		t.Errorf("bounce_status[1] = %v, want \"hard\"", bounceArray[1])
	}

	t.Logf("Successfully marshaled profile with bounce_status as array [bool, string]")
}
