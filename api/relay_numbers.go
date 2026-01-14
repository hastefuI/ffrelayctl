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

func (c *Client) GetRelayNumberSuggestions() (*RelayNumberSuggestions, error) {
	path := relayNumbersPath + "suggestions/"
	resp, err := c.Get(path)
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

func (c *Client) SearchRelayNumbers(areaCode string) ([]PhoneNumberOption, error) {
	path := fmt.Sprintf("%ssearch/?area_code=%s", relayNumbersPath, areaCode)
	resp, err := c.Get(path)
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
