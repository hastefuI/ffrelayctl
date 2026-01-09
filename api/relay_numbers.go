package api

import (
	"encoding/json"
	"io"
	"net/http"
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
