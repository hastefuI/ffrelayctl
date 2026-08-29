package api

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
)

const usersPath = APIBasePath + "users/"

func (c *Client) ListUsers(ctx context.Context) ([]User, error) {
	resp, err := c.Get(ctx, usersPath)
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

	var users []User
	if err := json.Unmarshal(body, &users); err != nil {
		return nil, err
	}

	return users, nil
}
