package api

import (
	"fmt"
	"io"
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
			addresses, err := client.ListDomainAddresses(t.Context())

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
			address, err := client.GetDomainAddress(t.Context(), tt.id)

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
			address, err := client.CreateDomainAddress(t.Context(), tt.request)

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
				Enabled:         new(false),
				Description:     new("Updated description"),
				BlockListEmails: new(true),
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

			url := fmt.Sprintf("%s%s%d/", DefaultBaseURL, domainAddressesPath, tt.id)
			httpmock.RegisterResponder(
				http.MethodPatch,
				url,
				httpmock.NewStringResponder(tt.mockStatusCode, tt.mockResponse),
			)

			client := NewClient("test")
			address, err := client.UpdateDomainAddress(t.Context(), tt.id, tt.request)

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
			err := client.DeleteDomainAddress(t.Context(), tt.id)

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
	_, err := client.ListDomainAddresses(t.Context())

	if err == nil {
		t.Error("ListDomainAddresses() expected error for invalid JSON, got nil")
	}
}

func TestClient_FilterDomainAddresses(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	tests := []struct {
		name      string
		filter    MaskFilter
		wantQuery string
	}{
		{
			name:      "empty filter sends no query",
			filter:    MaskFilter{},
			wantQuery: "",
		},
		{
			name:      "generated for is dropped",
			filter:    MaskFilter{GeneratedFor: "go.dev"},
			wantQuery: "",
		},
		{
			name:      "generated for is dropped from a filter that keeps its other fields",
			filter:    MaskFilter{Enabled: new(true), GeneratedFor: "go.dev", UsedOn: "go"},
			wantQuery: "enabled=true&used_on=go",
		},
		{
			name:      "address and description",
			filter:    MaskFilter{Address: "shopping", Description: "Shopping"},
			wantQuery: "address=shopping&description=Shopping",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			httpmock.Reset()

			var gotQuery string
			httpmock.RegisterResponder(
				http.MethodGet,
				DefaultBaseURL+domainAddressesPath,
				func(req *http.Request) (*http.Response, error) {
					gotQuery = req.URL.RawQuery
					return httpmock.NewStringResponse(http.StatusOK, `[]`), nil
				},
			)

			client := NewClient("test")
			if _, err := client.FilterDomainAddresses(t.Context(), tt.filter); err != nil {
				t.Fatalf("FilterDomainAddresses() error = %v", err)
			}

			if gotQuery != tt.wantQuery {
				t.Errorf("FilterDomainAddresses() query = %q, want %q", gotQuery, tt.wantQuery)
			}
		})
	}
}

func TestClient_DomainAddress_UsedOnRequestBody(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	t.Run("create sends used on", func(t *testing.T) {
		httpmock.Reset()

		var gotBody string
		httpmock.RegisterResponder(
			http.MethodPost,
			DefaultBaseURL+domainAddressesPath,
			func(req *http.Request) (*http.Response, error) {
				body, err := io.ReadAll(req.Body)
				if err != nil {
					return nil, err
				}
				gotBody = string(body)
				return httpmock.NewStringResponse(http.StatusCreated, `{"id": 12345}`), nil
			},
		)

		client := NewClient("test")
		req := CreateDomainAddressRequest{Address: "godocs", Enabled: true, UsedOn: "go.dev"}
		if _, err := client.CreateDomainAddress(t.Context(), req); err != nil {
			t.Fatalf("CreateDomainAddress() error = %v", err)
		}

		want := `{"address":"godocs","enabled":true,"block_list_emails":false,"used_on":"go.dev"}`
		if gotBody != want {
			t.Errorf("CreateDomainAddress() body = %s, want %s", gotBody, want)
		}
	})

	t.Run("update sends used on", func(t *testing.T) {
		httpmock.Reset()

		var gotBody string
		httpmock.RegisterResponder(
			http.MethodPatch,
			fmt.Sprintf("%s%s12345/", DefaultBaseURL, domainAddressesPath),
			func(req *http.Request) (*http.Response, error) {
				body, err := io.ReadAll(req.Body)
				if err != nil {
					return nil, err
				}
				gotBody = string(body)
				return httpmock.NewStringResponse(http.StatusOK, `{"id": 12345}`), nil
			},
		)

		client := NewClient("test")
		req := UpdateDomainAddressRequest{UsedOn: new("go.dev")}
		if _, err := client.UpdateDomainAddress(t.Context(), 12345, req); err != nil {
			t.Fatalf("UpdateDomainAddress() error = %v", err)
		}

		want := `{"used_on":"go.dev"}`
		if gotBody != want {
			t.Errorf("UpdateDomainAddress() body = %s, want %s", gotBody, want)
		}
	})
}
