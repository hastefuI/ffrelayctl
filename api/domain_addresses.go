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
	domainAddressesPath = APIBasePath + "domainaddresses/"
)

func (c *Client) ListDomainAddresses(ctx context.Context) ([]DomainAddress, error) {
	resp, err := c.Get(ctx, domainAddressesPath)
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

	var addresses []DomainAddress
	if err := json.Unmarshal(body, &addresses); err != nil {
		return nil, err
	}

	return addresses, nil
}

func (c *Client) GetDomainAddress(ctx context.Context, id int) (*DomainAddress, error) {
	path := fmt.Sprintf("%s%d/", domainAddressesPath, id)
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

	var address DomainAddress
	if err := json.Unmarshal(body, &address); err != nil {
		return nil, err
	}

	return &address, nil
}

func (c *Client) CreateDomainAddress(ctx context.Context, req CreateDomainAddressRequest) (*DomainAddress, error) {
	jsonBody, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	resp, err := c.Post(ctx, domainAddressesPath, bytes.NewReader(jsonBody))
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

	var address DomainAddress
	if err := json.Unmarshal(body, &address); err != nil {
		return nil, err
	}

	return &address, nil
}

func (c *Client) UpdateDomainAddress(ctx context.Context, id int, req UpdateDomainAddressRequest) (*DomainAddress, error) {
	jsonBody, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	path := fmt.Sprintf("%s%d/", domainAddressesPath, id)
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

	var address DomainAddress
	if err := json.Unmarshal(body, &address); err != nil {
		return nil, err
	}

	return &address, nil
}

func (c *Client) DeleteDomainAddress(ctx context.Context, id int) error {
	path := fmt.Sprintf("%s%d/", domainAddressesPath, id)
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
