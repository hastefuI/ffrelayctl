package api

import (
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
