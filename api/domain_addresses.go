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

// ListDomainAddresses returns every mask on the account's own subdomain.
func (c *Client) ListDomainAddresses(ctx context.Context) ([]DomainAddress, error) {
	return c.FilterDomainAddresses(ctx, MaskFilter{})
}

// FilterDomainAddresses returns the masks on the account's own subdomain that
// match filter. A zero MaskFilter returns every such mask, as
// ListDomainAddresses does. MaskFilter.GeneratedFor is ignored, because only
// random masks carry that field.
func (c *Client) FilterDomainAddresses(ctx context.Context, filter MaskFilter) ([]DomainAddress, error) {
	resp, err := c.Get(ctx, domainAddressesPath+filter.query(false))
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

// GetDomainAddress returns the subdomain mask with the given id.
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

// CreateDomainAddress creates a mask on the account's own subdomain using the
// address given in req, and returns it.
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

// UpdateDomainAddress applies the fields set in req to the subdomain mask with
// the given id and returns the updated mask.
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

// DeleteDomainAddress deletes the subdomain mask with the given id. Mail sent
// to a deleted mask is no longer forwarded.
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
