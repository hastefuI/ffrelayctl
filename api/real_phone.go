package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const realPhonePath = APIBasePath + "realphone/"

func (c *Client) GetRealPhone() ([]RealPhone, error) {
	resp, err := c.Get(realPhonePath)
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

func (c *Client) RegisterRealPhone(req RegisterRealPhoneRequest) (*RealPhone, error) {
	data, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	resp, err := c.Post(realPhonePath, bytes.NewReader(data))
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

func (c *Client) VerifyRealPhone(id int, req VerifyRealPhoneRequest) (*RealPhone, error) {
	data, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	verifyPath := fmt.Sprintf("%s%d/", realPhonePath, id)
	resp, err := c.Patch(verifyPath, bytes.NewReader(data))
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

func (c *Client) DeleteRealPhone(id int) error {
	path := fmt.Sprintf("%s%d/", realPhonePath, id)
	resp, err := c.Delete(path)
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
