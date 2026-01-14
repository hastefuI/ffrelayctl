package api

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
)

func TestClient_ListDomainAddresses(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	tests := []struct {
		name           string
		mockResponse   string
		mockStatusCode int
		wantErr        bool
		wantLen        int
		validate       func(*testing.T, []DomainAddress)
	}{
		{
			name: "successful list domain addresses",
			mockResponse: `[{
				"id": 12345,
				"address": "shopping",
				"full_address": "shopping@mysubdomain.mozmail.com",
				"enabled": true,
				"description": "Shopping sites",
				"block_list_emails": false,
				"created_at": "2025-01-01T00:00:00Z",
				"last_used_at": "2025-01-02T00:00:00Z",
				"num_forwarded": 25,
				"num_blocked": 3,
				"num_replied": 1,
				"num_spam": 0
			}]`,
			mockStatusCode: http.StatusOK,
			wantErr:        false,
			wantLen:        1,
			validate: func(t *testing.T, addresses []DomainAddress) {
				if addresses[0].ID != 12345 {
					t.Errorf("DomainAddress ID = %d, want 12345", addresses[0].ID)
				}
				if addresses[0].Address != "shopping" {
					t.Errorf("DomainAddress Address = %s, want shopping", addresses[0].Address)
				}
				if addresses[0].FullAddress != "shopping@mysubdomain.mozmail.com" {
					t.Errorf("DomainAddress FullAddress = %s, want shopping@mysubdomain.mozmail.com", addresses[0].FullAddress)
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
		{
			name:           "forbidden - no premium",
			mockResponse:   `{"detail": "You do not have permission to perform this action."}`,
			mockStatusCode: http.StatusForbidden,
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
				DefaultBaseURL+domainAddressesPath,
				httpmock.NewStringResponder(tt.mockStatusCode, tt.mockResponse),
			)

			client := NewClient("test")
			addresses, err := client.ListDomainAddresses()

			if (err != nil) != tt.wantErr {
				t.Errorf("ListDomainAddresses() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if len(addresses) != tt.wantLen {
					t.Errorf("ListDomainAddresses() returned %d addresses, want %d", len(addresses), tt.wantLen)
				}
				if tt.validate != nil {
					tt.validate(t, addresses)
				}
			}
		})
	}
}

func TestClient_GetDomainAddress(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	tests := []struct {
		name           string
		id             int
		mockResponse   string
		mockStatusCode int
		wantErr        bool
		validate       func(*testing.T, *DomainAddress)
	}{
		{
			name: "successful get domain address",
			id:   12345,
			mockResponse: `{
				"id": 12345,
				"address": "work",
				"full_address": "work@mysubdomain.mozmail.com",
				"enabled": true,
				"description": "Work related",
				"block_list_emails": true,
				"created_at": "2025-01-01T00:00:00Z",
				"last_used_at": null,
				"num_forwarded": 100,
				"num_blocked": 10,
				"num_replied": 5,
				"num_spam": 2
			}`,
			mockStatusCode: http.StatusOK,
			wantErr:        false,
			validate: func(t *testing.T, address *DomainAddress) {
				if address.ID != 12345 {
					t.Errorf("DomainAddress ID = %d, want 12345", address.ID)
				}
				if address.Address != "work" {
					t.Errorf("DomainAddress Address = %s, want work", address.Address)
				}
				if !address.BlockListEmails {
					t.Error("DomainAddress BlockListEmails = false, want true")
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

			url := fmt.Sprintf("%s%s%d/", DefaultBaseURL, domainAddressesPath, tt.id)
			httpmock.RegisterResponder(
				http.MethodGet,
				url,
				httpmock.NewStringResponder(tt.mockStatusCode, tt.mockResponse),
			)

			client := NewClient("test")
			address, err := client.GetDomainAddress(tt.id)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetDomainAddress() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && tt.validate != nil {
				tt.validate(t, address)
			}
		})
	}
}

func TestClient_CreateDomainAddress(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	tests := []struct {
		name           string
		request        CreateDomainAddressRequest
		mockResponse   string
		mockStatusCode int
		wantErr        bool
		validate       func(*testing.T, *DomainAddress)
	}{
		{
			name: "successful create",
			request: CreateDomainAddressRequest{
				Address:         "newaddress",
				Enabled:         true,
				Description:     "New domain address",
				BlockListEmails: false,
			},
			mockResponse: `{
				"id": 99999,
				"address": "newaddress",
				"full_address": "newaddress@mysubdomain.mozmail.com",
				"enabled": true,
				"description": "New domain address",
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
			validate: func(t *testing.T, address *DomainAddress) {
				if address.ID != 99999 {
					t.Errorf("DomainAddress ID = %d, want 99999", address.ID)
				}
				if address.Address != "newaddress" {
					t.Errorf("DomainAddress Address = %s, want newaddress", address.Address)
				}
			},
		},
		{
			name: "forbidden - no premium",
			request: CreateDomainAddressRequest{
				Address: "test",
				Enabled: true,
			},
			mockResponse:   `{"detail": "You do not have permission to perform this action."}`,
			mockStatusCode: http.StatusForbidden,
			wantErr:        true,
			validate:       nil,
		},
		{
			name: "address already exists",
			request: CreateDomainAddressRequest{
				Address: "existing",
				Enabled: true,
			},
			mockResponse:   `{"address": ["This address already exists."]}`,
			mockStatusCode: http.StatusBadRequest,
			wantErr:        true,
			validate:       nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			httpmock.Reset()

			httpmock.RegisterResponder(
				http.MethodPost,
				DefaultBaseURL+domainAddressesPath,
				httpmock.NewStringResponder(tt.mockStatusCode, tt.mockResponse),
			)

			client := NewClient("test")
			address, err := client.CreateDomainAddress(tt.request)

			if (err != nil) != tt.wantErr {
				t.Errorf("CreateDomainAddress() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && tt.validate != nil {
				tt.validate(t, address)
			}
		})
	}
}

func TestClient_UpdateDomainAddress(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	enabled := false
	description := "Updated description"
	blockList := true

	tests := []struct {
		name           string
		id             int
		request        UpdateDomainAddressRequest
		mockResponse   string
		mockStatusCode int
		wantErr        bool
		validate       func(*testing.T, *DomainAddress)
	}{
		{
			name: "successful update",
			id:   12345,
			request: UpdateDomainAddressRequest{
				Enabled:         &enabled,
				Description:     &description,
				BlockListEmails: &blockList,
			},
			mockResponse: `{
				"id": 12345,
				"address": "work",
				"full_address": "work@mysubdomain.mozmail.com",
				"enabled": false,
				"description": "Updated description",
				"block_list_emails": true,
				"created_at": "2025-01-01T00:00:00Z",
				"last_used_at": null,
				"num_forwarded": 100,
				"num_blocked": 10,
				"num_replied": 5,
				"num_spam": 2
			}`,
			mockStatusCode: http.StatusOK,
			wantErr:        false,
			validate: func(t *testing.T, address *DomainAddress) {
				if address.Enabled {
					t.Error("DomainAddress Enabled = true, want false")
				}
				if address.Description != "Updated description" {
					t.Errorf("DomainAddress Description = %s, want Updated description", address.Description)
				}
				if !address.BlockListEmails {
					t.Error("DomainAddress BlockListEmails = false, want true")
				}
			},
		},
		{
			name: "not found",
			id:   99999,
			request: UpdateDomainAddressRequest{
				Enabled: &enabled,
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

			url := fmt.Sprintf("%s%s%d/", DefaultBaseURL, domainAddressesPath, tt.id)
			httpmock.RegisterResponder(
				http.MethodPatch,
				url,
				httpmock.NewStringResponder(tt.mockStatusCode, tt.mockResponse),
			)

			client := NewClient("test")
			address, err := client.UpdateDomainAddress(tt.id, tt.request)

			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateDomainAddress() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && tt.validate != nil {
				tt.validate(t, address)
			}
		})
	}
}

func TestClient_DeleteDomainAddress(t *testing.T) {
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

			url := fmt.Sprintf("%s%sdomainaddresses/%d/", DefaultBaseURL, APIBasePath, tt.id)
			httpmock.RegisterResponder(
				http.MethodDelete,
				url,
				httpmock.NewStringResponder(tt.mockStatusCode, tt.mockResponse),
			)

			client := NewClient("test")
			err := client.DeleteDomainAddress(tt.id)

			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteDomainAddress() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClient_ListDomainAddresses_InvalidJSON(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(
		http.MethodGet,
		DefaultBaseURL+domainAddressesPath,
		httpmock.NewStringResponder(http.StatusOK, `{invalid json`),
	)

	client := NewClient("test")
	_, err := client.ListDomainAddresses()

	if err == nil {
		t.Error("ListDomainAddresses() expected error for invalid JSON, got nil")
	}
}
