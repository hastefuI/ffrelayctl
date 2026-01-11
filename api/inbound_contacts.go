package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

const (
	inboundContactsPath = APIBasePath + "inboundcontact/"
)

func (c *Client) ListInboundContacts() ([]InboundContact, error) {
	resp, err := c.Get(inboundContactsPath)
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

	var contacts []InboundContact
	if err := json.Unmarshal(body, &contacts); err != nil {
		return nil, err
	}

	return contacts, nil
}

func (c *Client) UpdateInboundContact(id int, req UpdateInboundContactRequest) (*InboundContact, error) {
	jsonBody, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	path := fmt.Sprintf("%s%d/", inboundContactsPath, id)
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

	var contact InboundContact
	if err := json.Unmarshal(body, &contact); err != nil {
		return nil, err
	}

	return &contact, nil
}
