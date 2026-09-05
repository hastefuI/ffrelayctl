package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const realPhonePath = APIBasePath + "realphone/"

// GetRealPhone returns the phone numbers registered on the account, whether or
// not they are verified.
func (c *Client) GetRealPhone(ctx context.Context) ([]RealPhone, error) {
	resp, err := c.Get(ctx, realPhonePath)
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

	var phones []RealPhone
	if err := json.Unmarshal(body, &phones); err != nil {
		return nil, err
	}

	return phones, nil
}

// RegisterRealPhone registers a real phone number and asks Relay to text it a
// verification code. The number returned is unverified until VerifyRealPhone
// is called with that code.
func (c *Client) RegisterRealPhone(ctx context.Context, req RegisterRealPhoneRequest) (*RealPhone, error) {
	data, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	resp, err := c.Post(ctx, realPhonePath, bytes.NewReader(data))
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

	var phone RealPhone
	if err := json.Unmarshal(body, &phone); err != nil {
		return nil, err
	}

	return &phone, nil
}

// VerifyRealPhone confirms the registration with the given id using the code
// Relay sent by text.
func (c *Client) VerifyRealPhone(ctx context.Context, id int, req VerifyRealPhoneRequest) (*RealPhone, error) {
	data, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	verifyPath := fmt.Sprintf("%s%d/", realPhonePath, id)
	resp, err := c.Patch(ctx, verifyPath, bytes.NewReader(data))
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

	var phone RealPhone
	if err := json.Unmarshal(body, &phone); err != nil {
		return nil, err
	}

	return &phone, nil
}

// DeleteRealPhone removes the registered phone number with the given id.
func (c *Client) DeleteRealPhone(ctx context.Context, id int) error {
	path := fmt.Sprintf("%s%d/", realPhonePath, id)
	resp, err := c.Delete(ctx, path)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode >= http.StatusBadRequest {
		return &APIError{StatusCode: resp.StatusCode, Body: string(body)}
	}

	return nil
}
