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

func (c *Client) ListRelayAddresses(ctx context.Context) ([]RelayAddress, error) {
	resp, err := c.Get(ctx, relayAddressesPath)
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
