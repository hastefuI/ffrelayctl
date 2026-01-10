package api

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
)

func TestClient_ListRelayNumbers(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	tests := []struct {
		name           string
		mockResponse   string
		mockStatusCode int
		wantErr        bool
		wantCount      int
	}{
		{
			name: "successful list relay numbers",
			mockResponse: `[
				{
					"id": 1,
					"number": "+18001234567",
					"enabled": true,
					"location": "US",
					"vendor_id": "vendor1",
					"country_code": "US",
					"created_at": "2026-01-01T00:00:00Z",
					"remaining_texts": 75,
					"remaining_minutes": 50,
					"calls_forwarded": 10,
					"calls_blocked": 2,
					"texts_forwarded": 25,
					"texts_blocked": 1
				},
				{
					"id": 2,
					"number": "+18002345678",
					"enabled": false,
					"location": "US",
					"vendor_id": "vendor2",
					"country_code": "US",
					"created_at": "2024-02-01T00:00:00Z",
					"remaining_texts": 100,
					"remaining_minutes": 50,
					"calls_forwarded": 0,
					"calls_blocked": 0,
					"texts_forwarded": 0,
					"texts_blocked": 0
				}
			]`,
			mockStatusCode: http.StatusOK,
			wantErr:        false,
			wantCount:      2,
		},
		{
			name:           "empty list",
			mockResponse:   `[]`,
			mockStatusCode: http.StatusOK,
			wantErr:        false,
			wantCount:      0,
		},
		{
			name: "null created_at field",
			mockResponse: `[
				{
					"id": 3,
					"number": "+18003456789",
					"enabled": true,
					"location": "US",
					"vendor_id": "vendor3",
					"country_code": "US",
					"created_at": null,
					"remaining_texts": 100,
					"remaining_minutes": 50,
					"calls_forwarded": 0,
					"calls_blocked": 0,
					"texts_forwarded": 0,
					"texts_blocked": 0
				}
			]`,
			mockStatusCode: http.StatusOK,
			wantErr:        false,
			wantCount:      1,
		},
		{
			name:           "unauthorized",
			mockResponse:   `{"detail": "Invalid token"}`,
			mockStatusCode: http.StatusUnauthorized,
			wantErr:        true,
		},
		{
			name:           "forbidden - no phone subscription",
			mockResponse:   `{"detail": "Phone subscription required"}`,
			mockStatusCode: http.StatusForbidden,
			wantErr:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			httpmock.Reset()

			client := NewClient("test")

			httpmock.RegisterResponder("GET", DefaultBaseURL+"/api/v1/relaynumber/",
				httpmock.NewStringResponder(tt.mockStatusCode, tt.mockResponse))

			numbers, err := client.ListRelayNumbers()

			if (err != nil) != tt.wantErr {
				t.Errorf("ListRelayNumbers() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if len(numbers) != tt.wantCount {
					t.Errorf("ListRelayNumbers() got %d numbers, want %d", len(numbers), tt.wantCount)
				}

				if tt.wantCount > 0 {
					if numbers[0].Number == "" {
						t.Error("ListRelayNumbers() first number has empty Number field")
					}
					if numbers[0].Location == "" {
						t.Error("ListRelayNumbers() first number has empty Location field")
					}
				}

				if tt.name == "null created_at field" && tt.wantCount > 0 {
					if numbers[0].CreatedAt != nil {
						t.Errorf("ListRelayNumbers() expected nil CreatedAt, got %v", *numbers[0].CreatedAt)
					}
				}
			}
		})
	}
}

func TestClient_ListRelayNumbers_InvalidJSON(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	client := NewClient("test")

	httpmock.RegisterResponder("GET", DefaultBaseURL+"/api/v1/relaynumber/",
		httpmock.NewStringResponder(http.StatusOK, `invalid json`))

	_, err := client.ListRelayNumbers()
	if err == nil {
		t.Error("ListRelayNumbers() expected error for invalid JSON, got nil")
	}
}

func TestClient_UpdateRelayNumber(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	trueVal := true
	falseVal := false

	tests := []struct {
		name           string
		relayID        int
		request        UpdateRelayNumberRequest
		mockResponse   string
		mockStatusCode int
		wantErr        bool
		wantEnabled    bool
	}{
		{
			name:    "successfully disable relay number",
			relayID: 1,
			request: UpdateRelayNumberRequest{
				Enabled: &falseVal,
			},
			mockResponse: `{
				"id": 1,
				"number": "+18001234567",
				"enabled": false,
				"location": "US",
				"vendor_id": "vendor1",
				"country_code": "US",
				"created_at": "2026-01-01T00:00:00Z",
				"remaining_texts": 75,
				"remaining_minutes": 50,
				"calls_forwarded": 10,
				"calls_blocked": 2,
				"texts_forwarded": 25,
				"texts_blocked": 1
			}`,
			mockStatusCode: http.StatusOK,
			wantErr:        false,
			wantEnabled:    false,
		},
		{
			name:    "successfully enable relay number",
			relayID: 2,
			request: UpdateRelayNumberRequest{
				Enabled: &trueVal,
			},
			mockResponse: `{
				"id": 2,
				"number": "+18002345678",
				"enabled": true,
				"location": "US",
				"vendor_id": "vendor2",
				"country_code": "US",
				"created_at": "2024-02-01T00:00:00Z",
				"remaining_texts": 100,
				"remaining_minutes": 50,
				"calls_forwarded": 0,
				"calls_blocked": 0,
				"texts_forwarded": 0,
				"texts_blocked": 0
			}`,
			mockStatusCode: http.StatusOK,
			wantErr:        false,
			wantEnabled:    true,
		},
		{
			name:    "relay number not found",
			relayID: 999,
			request: UpdateRelayNumberRequest{
				Enabled: &falseVal,
			},
			mockResponse:   `{"detail": "Not found"}`,
			mockStatusCode: http.StatusNotFound,
			wantErr:        true,
		},
		{
			name:    "unauthorized",
			relayID: 1,
			request: UpdateRelayNumberRequest{
				Enabled: &falseVal,
			},
			mockResponse:   `{"detail": "Invalid token"}`,
			mockStatusCode: http.StatusUnauthorized,
			wantErr:        true,
		},
		{
			name:    "forbidden - no phone subscription",
			relayID: 1,
			request: UpdateRelayNumberRequest{
				Enabled: &falseVal,
			},
			mockResponse:   `{"detail": "Phone subscription required"}`,
			mockStatusCode: http.StatusForbidden,
			wantErr:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			httpmock.Reset()

			client := NewClient("test")

			url := fmt.Sprintf("%s%s%d/", DefaultBaseURL, relayNumbersPath, tt.relayID)
			httpmock.RegisterResponder("PATCH", url,
				httpmock.NewStringResponder(tt.mockStatusCode, tt.mockResponse))

			number, err := client.UpdateRelayNumber(tt.relayID, tt.request)

			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateRelayNumber() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if number == nil {
					t.Error("UpdateRelayNumber() returned nil number")
					return
				}

				if number.ID != tt.relayID {
					t.Errorf("UpdateRelayNumber() got ID = %d, want %d", number.ID, tt.relayID)
				}

				if number.Enabled != tt.wantEnabled {
					t.Errorf("UpdateRelayNumber() got Enabled = %v, want %v", number.Enabled, tt.wantEnabled)
				}

				if number.Number == "" {
					t.Error("UpdateRelayNumber() returned number with empty Number field")
				}
			}
		})
	}
}

func TestClient_UpdateRelayNumber_InvalidJSON(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	client := NewClient("test")
	enabled := false

	httpmock.RegisterResponder("PATCH", DefaultBaseURL+"/api/v1/relaynumber/1/",
		httpmock.NewStringResponder(http.StatusOK, `invalid json`))

	_, err := client.UpdateRelayNumber(1, UpdateRelayNumberRequest{Enabled: &enabled})
	if err == nil {
		t.Error("UpdateRelayNumber() expected error for invalid JSON, got nil")
	}
}
