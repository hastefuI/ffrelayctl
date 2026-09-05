package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const (
	relayNumbersPath = APIBasePath + "relaynumber/"
)

// ListRelayNumbers returns the phone masks on the account.
func (c *Client) ListRelayNumbers(ctx context.Context) ([]RelayNumber, error) {
	resp, err := c.Get(ctx, relayNumbersPath)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= http.StatusBadRequest {
		return nil, &APIError{StatusCode: resp.StatusCode, Body: string(body)}
	}

	var numbers []RelayNumber
	if err := json.Unmarshal(body, &numbers); err != nil {
		return nil, err
	}

	return numbers, nil
}

// GetRelayNumberSuggestions returns the numbers Relay offers for a new phone
// mask, grouped by how closely each one matches the real number. Relay answers
// with 400 when it has nothing to offer, which is reported as empty lists
// rather than as an error.
func (c *Client) GetRelayNumberSuggestions(ctx context.Context) (*RelayNumberSuggestions, error) {
	path := relayNumbersPath + "suggestions/"
	resp, err := c.Get(ctx, path)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == http.StatusBadRequest {
		return &RelayNumberSuggestions{
			SamePrefixOptions: []PhoneNumberOption{},
			OtherAreasOptions: []PhoneNumberOption{},
			SameAreaOptions:   []PhoneNumberOption{},
			RandomOptions:     []PhoneNumberOption{},
		}, nil
	}

	if resp.StatusCode >= http.StatusBadRequest {
		return nil, &APIError{StatusCode: resp.StatusCode, Body: string(body)}
	}

	var suggestions RelayNumberSuggestions
	if err := json.Unmarshal(body, &suggestions); err != nil {
		return nil, err
	}

	return &suggestions, nil
}

// SearchRelayNumbers returns the numbers free to be claimed in areaCode. Relay
// answers with 400 when nothing matches, which is reported as an empty slice
// rather than as an error.
func (c *Client) SearchRelayNumbers(ctx context.Context, areaCode string) ([]PhoneNumberOption, error) {
	path := fmt.Sprintf("%ssearch/?area_code=%s", relayNumbersPath, areaCode)
	resp, err := c.Get(ctx, path)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == http.StatusBadRequest {
		return []PhoneNumberOption{}, nil
	}

	if resp.StatusCode > http.StatusBadRequest {
		return nil, &APIError{StatusCode: resp.StatusCode, Body: string(body)}
	}

	var numbers []PhoneNumberOption
	if err := json.Unmarshal(body, &numbers); err != nil {
		return nil, err
	}

	return numbers, nil
}

// UpdateRelayNumber applies the fields set in req to the phone mask with the
// given id and returns the updated mask.
func (c *Client) UpdateRelayNumber(ctx context.Context, id int, req UpdateRelayNumberRequest) (*RelayNumber, error) {
	jsonBody, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	path := fmt.Sprintf("%s%d/", relayNumbersPath, id)
	resp, err := c.Patch(ctx, path, bytes.NewReader(jsonBody))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= http.StatusBadRequest {
		return nil, &APIError{StatusCode: resp.StatusCode, Body: string(body)}
	}

	var number RelayNumber
	if err := json.Unmarshal(body, &number); err != nil {
		return nil, err
	}

	return &number, nil
}
