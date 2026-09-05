package api

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
)

const (
	profilesPath = APIBasePath + "profiles/"
)

// GetProfiles returns the profiles for the account. Relay returns one.
func (c *Client) GetProfiles(ctx context.Context) ([]Profile, error) {
	resp, err := c.Get(ctx, profilesPath)
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

	var profiles []Profile
	if err := json.Unmarshal(body, &profiles); err != nil {
		return nil, err
	}

	return profiles, nil
}
