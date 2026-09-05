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
	relayAddressesPath = APIBasePath + "relayaddresses/"
)

// ListRelayAddresses returns every random email mask on the account.
func (c *Client) ListRelayAddresses(ctx context.Context) ([]RelayAddress, error) {
	return c.FilterRelayAddresses(ctx, MaskFilter{})
}

// FilterRelayAddresses returns the random email masks that match filter. A
// zero MaskFilter returns every random mask, as ListRelayAddresses does.
func (c *Client) FilterRelayAddresses(ctx context.Context, filter MaskFilter) ([]RelayAddress, error) {
	resp, err := c.Get(ctx, relayAddressesPath+filter.query(true))
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

	var addresses []RelayAddress
	if err := json.Unmarshal(body, &addresses); err != nil {
		return nil, err
	}

	return addresses, nil
}

// GetRelayAddress returns the random email mask with the given id.
func (c *Client) GetRelayAddress(ctx context.Context, id int) (*RelayAddress, error) {
	path := fmt.Sprintf("%s%d/", relayAddressesPath, id)
	resp, err := c.Get(ctx, path)
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

	var address RelayAddress
	if err := json.Unmarshal(body, &address); err != nil {
		return nil, err
	}

	return &address, nil
}

// CreateRelayAddress creates a random email mask and returns it. Relay picks
// the address, so req only carries the settings for the new mask.
func (c *Client) CreateRelayAddress(ctx context.Context, req CreateRelayAddressRequest) (*RelayAddress, error) {
	jsonBody, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	resp, err := c.Post(ctx, relayAddressesPath, bytes.NewReader(jsonBody))
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

	var address RelayAddress
	if err := json.Unmarshal(body, &address); err != nil {
		return nil, err
	}

	return &address, nil
}

// UpdateRelayAddress applies the fields set in req to the random email mask
// with the given id and returns the updated mask.
func (c *Client) UpdateRelayAddress(ctx context.Context, id int, req UpdateRelayAddressRequest) (*RelayAddress, error) {
	jsonBody, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	path := fmt.Sprintf("%s%d/", relayAddressesPath, id)
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

	var address RelayAddress
	if err := json.Unmarshal(body, &address); err != nil {
		return nil, err
	}

	return &address, nil
}

// DeleteRelayAddress deletes the random email mask with the given id. Mail
// sent to a deleted mask is no longer forwarded.
func (c *Client) DeleteRelayAddress(ctx context.Context, id int) error {
	path := fmt.Sprintf("%s%d/", relayAddressesPath, id)
	resp, err := c.Delete(ctx, path)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= http.StatusBadRequest {
		body, _ := io.ReadAll(resp.Body)
		return &APIError{StatusCode: resp.StatusCode, Body: string(body)}
	}

	return nil
}
