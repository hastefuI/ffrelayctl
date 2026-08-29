package api

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
)

func TestClient_ListRelayAddresses(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	tests := []struct {
		name           string
		mockResponse   string
		mockStatusCode int
		wantErr        bool
		wantLen        int
		validate       func(*testing.T, []RelayAddress)
	}{
		{
			name: "successful list relay addresses",
			mockResponse: `[{
				"id": 12345,
				"address": "abc123",
				"domain": 1,
				"full_address": "abc123@relay.firefox.com",
				"enabled": true,
				"description": "Shopping",
				"generated_for": "amazon.com",
				"used_on": "",
				"block_list_emails": false,
				"created_at": "2025-01-01T00:00:00Z",
				"last_used_at": "2025-01-02T00:00:00Z",
				"num_forwarded": 50,
				"num_blocked": 5,
				"num_replied": 2,
				"num_spam": 1
			}]`,
			mockStatusCode: http.StatusOK,
			wantErr:        false,
			wantLen:        1,
			validate: func(t *testing.T, addresses []RelayAddress) {
				if addresses[0].ID != 12345 {
					t.Errorf("RelayAddress ID = %d, want 12345", addresses[0].ID)
				}
				if addresses[0].FullAddress != "abc123@relay.firefox.com" {
					t.Errorf("RelayAddress FullAddress = %s, want abc123@relay.firefox.com", addresses[0].FullAddress)
				}
				if !addresses[0].Enabled {
					t.Error("RelayAddress Enabled = false, want true")
				}
			},
		},
		{
			name:           "empty list",
			mockResponse:   `[]`,
			mockStatusCode: http.StatusOK,
			wantErr:        false,
			wantLen:        0,
			validate:       nil,
		},
		{
			name:           "unauthorized",
			mockResponse:   `{"detail": "Authentication credentials were not provided."}`,
			mockStatusCode: http.StatusUnauthorized,
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
				DefaultBaseURL+relayAddressesPath,
				httpmock.NewStringResponder(tt.mockStatusCode, tt.mockResponse),
			)

			client := NewClient("test")
			addresses, err := client.ListRelayAddresses(t.Context())

			if (err != nil) != tt.wantErr {
				t.Errorf("ListRelayAddresses() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if len(addresses) != tt.wantLen {
					t.Errorf("ListRelayAddresses() returned %d addresses, want %d", len(addresses), tt.wantLen)
				}
				if tt.validate != nil {
					tt.validate(t, addresses)
				}
			}
		})
	}
}

func TestClient_GetRelayAddress(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	tests := []struct {
		name           string
		id             int
		mockResponse   string
		mockStatusCode int
		wantErr        bool
		validate       func(*testing.T, *RelayAddress)
	}{
		{
			name: "successful get relay address",
			id:   12345,
			mockResponse: `{
				"id": 12345,
				"address": "abc123",
				"domain": 1,
				"full_address": "abc123@relay.firefox.com",
				"enabled": true,
				"description": "Test address",
				"generated_for": "example.com",
				"used_on": "",
				"block_list_emails": false,
				"created_at": "2025-01-01T00:00:00Z",
				"last_used_at": null,
				"num_forwarded": 10,
				"num_blocked": 0,
				"num_replied": 0,
				"num_spam": 0
			}`,
			mockStatusCode: http.StatusOK,
			wantErr:        false,
			validate: func(t *testing.T, address *RelayAddress) {
				if address.ID != 12345 {
					t.Errorf("RelayAddress ID = %d, want 12345", address.ID)
				}
				if address.Description != "Test address" {
					t.Errorf("RelayAddress Description = %s, want Test address", address.Description)
				}
			},
		},
		{
			name:           "not found",
			id:             99999,
			mockResponse:   `{"detail": "Not found."}`,
			mockStatusCode: http.StatusNotFound,
			wantErr:        true,
			validate:       nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			httpmock.Reset()

			url := fmt.Sprintf("%s%s%d/", DefaultBaseURL, relayAddressesPath, tt.id)
			httpmock.RegisterResponder(
				http.MethodGet,
				url,
				httpmock.NewStringResponder(tt.mockStatusCode, tt.mockResponse),
			)

			client := NewClient("test")
			address, err := client.GetRelayAddress(t.Context(), tt.id)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetRelayAddress() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && tt.validate != nil {
				tt.validate(t, address)
			}
		})
	}
}

func TestClient_CreateRelayAddress(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	tests := []struct {
		name           string
		request        CreateRelayAddressRequest
		mockResponse   string
		mockStatusCode int
		wantErr        bool
		validate       func(*testing.T, *RelayAddress)
	}{
		{
			name: "successful create",
			request: CreateRelayAddressRequest{
				Enabled:         true,
				Description:     "New mask",
				GeneratedFor:    "example.com",
				BlockListEmails: false,
			},
			mockResponse: `{
				"id": 99999,
				"address": "newmask123",
				"domain": 1,
				"full_address": "newmask123@relay.firefox.com",
				"enabled": true,
				"description": "New mask",
				"generated_for": "example.com",
				"used_on": "",
				"block_list_emails": false,
				"created_at": "2025-01-01T00:00:00Z",
				"last_used_at": null,
				"num_forwarded": 0,
				"num_blocked": 0,
				"num_replied": 0,
				"num_spam": 0
			}`,
			mockStatusCode: http.StatusCreated,
			wantErr:        false,
			validate: func(t *testing.T, address *RelayAddress) {
				if address.ID != 99999 {
					t.Errorf("RelayAddress ID = %d, want 99999", address.ID)
				}
				if address.Description != "New mask" {
					t.Errorf("RelayAddress Description = %s, want New mask", address.Description)
				}
			},
		},
		{
			name: "at mask limit",
			request: CreateRelayAddressRequest{
				Enabled: true,
			},
			mockResponse:   `{"detail": "You have reached the maximum number of email masks."}`,
			mockStatusCode: http.StatusForbidden,
			wantErr:        true,
			validate:       nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			httpmock.Reset()

			httpmock.RegisterResponder(
				http.MethodPost,
				DefaultBaseURL+relayAddressesPath,
				httpmock.NewStringResponder(tt.mockStatusCode, tt.mockResponse),
			)

			client := NewClient("test")
			address, err := client.CreateRelayAddress(t.Context(), tt.request)

			if (err != nil) != tt.wantErr {
				t.Errorf("CreateRelayAddress() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && tt.validate != nil {
				tt.validate(t, address)
			}
		})
	}
}

func TestClient_UpdateRelayAddress(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	tests := []struct {
		name           string
		id             int
		request        UpdateRelayAddressRequest
		mockResponse   string
		mockStatusCode int
		wantErr        bool
		validate       func(*testing.T, *RelayAddress)
	}{
		{
			name: "successful update",
			id:   12345,
			request: UpdateRelayAddressRequest{
				Enabled:     new(false),
				Description: new("Updated description"),
			},
			mockResponse: `{
				"id": 12345,
				"address": "abc123",
				"domain": 1,
				"full_address": "abc123@relay.firefox.com",
				"enabled": false,
				"description": "Updated description",
				"generated_for": "example.com",
				"used_on": "",
				"block_list_emails": false,
				"created_at": "2025-01-01T00:00:00Z",
				"last_used_at": null,
				"num_forwarded": 10,
				"num_blocked": 0,
				"num_replied": 0,
				"num_spam": 0
			}`,
			mockStatusCode: http.StatusOK,
			wantErr:        false,
			validate: func(t *testing.T, address *RelayAddress) {
				if address.Enabled {
					t.Error("RelayAddress Enabled = true, want false")
				}
				if address.Description != "Updated description" {
					t.Errorf("RelayAddress Description = %s, want Updated description", address.Description)
				}
			},
		},
		{
			name: "not found",
			id:   99999,
			request: UpdateRelayAddressRequest{
				Enabled: new(false),
			},
			mockResponse:   `{"detail": "Not found."}`,
			mockStatusCode: http.StatusNotFound,
			wantErr:        true,
			validate:       nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			httpmock.Reset()

			url := fmt.Sprintf("%s%s%d/", DefaultBaseURL, relayAddressesPath, tt.id)
			httpmock.RegisterResponder(
				http.MethodPatch,
				url,
				httpmock.NewStringResponder(tt.mockStatusCode, tt.mockResponse),
			)

			client := NewClient("test")
			address, err := client.UpdateRelayAddress(t.Context(), tt.id, tt.request)

			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateRelayAddress() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && tt.validate != nil {
				tt.validate(t, address)
			}
		})
	}
}

func TestClient_DeleteRelayAddress(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	tests := []struct {
		name           string
		id             int
		mockStatusCode int
		mockResponse   string
		wantErr        bool
	}{
		{
			name:           "successful delete",
			id:             12345,
			mockStatusCode: http.StatusNoContent,
			mockResponse:   "",
			wantErr:        false,
		},
		{
			name:           "not found",
			id:             99999,
			mockStatusCode: http.StatusNotFound,
			mockResponse:   `{"detail": "Not found."}`,
			wantErr:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			httpmock.Reset()

			url := fmt.Sprintf("%s%s%d/", DefaultBaseURL, relayAddressesPath, tt.id)
			httpmock.RegisterResponder(
				http.MethodDelete,
				url,
				httpmock.NewStringResponder(tt.mockStatusCode, tt.mockResponse),
			)

			client := NewClient("test")
			err := client.DeleteRelayAddress(t.Context(), tt.id)

			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteRelayAddress() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClient_ListRelayAddresses_InvalidJSON(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(
		http.MethodGet,
		DefaultBaseURL+relayAddressesPath,
		httpmock.NewStringResponder(http.StatusOK, `{invalid json`),
	)

	client := NewClient("test")
	_, err := client.ListRelayAddresses(t.Context())

	if err == nil {
		t.Error("ListRelayAddresses() expected error for invalid JSON, got nil")
	}
}
