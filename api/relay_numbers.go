package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

const (
	relayNumbersPath = APIBasePath + "relaynumber/"
)

func (c *Client) ListRelayNumbers() ([]RelayNumber, error) {
	resp, err := c.Get(relayNumbersPath)
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

func (c *Client) UpdateRelayNumber(id int, req UpdateRelayNumberRequest) (*RelayNumber, error) {
	jsonBody, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	path := fmt.Sprintf("%s%d/", relayNumbersPath, id)
	resp, err := c.Patch(path, strings.NewReader(string(jsonBody)))
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
