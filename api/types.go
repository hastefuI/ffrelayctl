package api

import (
	"encoding/json"
	"fmt"
)

// Firefox Relay Profile
type Profile struct {
	ID                          int          `json:"id"`
	ServerStorage               bool         `json:"server_storage"`
	Subdomain                   *string      `json:"subdomain"`
	HasPremium                  bool         `json:"has_premium"`
	HasPhone                    bool         `json:"has_phone"`
	OnboardingState             int          `json:"onboarding_state"`
	DateSubscribed              *string      `json:"date_subscribed"`
	AvatarURL                   string       `json:"avatar"`
	EmailsBlocked               int          `json:"emails_blocked"`
	EmailsForwarded             int          `json:"emails_forwarded"`
	EmailsReplied               int          `json:"emails_replied"`
	LevelOneTrackersBlocked     int          `json:"level_one_trackers_blocked"`
	RemoveLevelOneEmailTrackers bool         `json:"remove_level_one_email_trackers"`
	AtMaskLimit                 bool         `json:"at_mask_limit"`
	BounceStatus                BounceStatus `json:"bounce_status"`
}

type BounceStatus struct {
	Paused bool
	Type   string
}

func (b *BounceStatus) UnmarshalJSON(data []byte) error {
	var tuple []interface{}
	if err := json.Unmarshal(data, &tuple); err != nil {
		return err
	}

	if len(tuple) != 2 {
		return fmt.Errorf("bounce_status: expected array of length 2, got %d", len(tuple))
	}

	paused, ok := tuple[0].(bool)
	if !ok {
		return fmt.Errorf("bounce_status[0]: expected bool, got %T", tuple[0])
	}
	b.Paused = paused

	bounceType, ok := tuple[1].(string)
	if !ok {
		return fmt.Errorf("bounce_status[1]: expected string, got %T", tuple[1])
	}
	b.Type = bounceType

	return nil
}

func (b BounceStatus) MarshalJSON() ([]byte, error) {
	return json.Marshal([]interface{}{b.Paused, b.Type})
}

type RelayAddress struct {
	ID              int     `json:"id"`
	Address         string  `json:"address"`
	Domain          int     `json:"domain"`
	FullAddress     string  `json:"full_address"`
	Enabled         bool    `json:"enabled"`
	Description     string  `json:"description"`
	GeneratedFor    string  `json:"generated_for"`
	UsedOn          string  `json:"used_on"`
	BlockListEmails bool    `json:"block_list_emails"`
	CreatedAt       string  `json:"created_at"`
	LastUsedAt      *string `json:"last_used_at"`
	NumForwarded    int     `json:"num_forwarded"`
	NumBlocked      int     `json:"num_blocked"`
	NumReplied      int     `json:"num_replied"`
	NumSpam         int     `json:"num_spam"`
}

type CreateRelayAddressRequest struct {
	Enabled         bool   `json:"enabled"`
	Description     string `json:"description,omitempty"`
	GeneratedFor    string `json:"generated_for,omitempty"`
	BlockListEmails bool   `json:"block_list_emails"`
	UsedOn          string `json:"used_on,omitempty"`
}

type UpdateRelayAddressRequest struct {
	Enabled         *bool   `json:"enabled,omitempty"`
	Description     *string `json:"description,omitempty"`
	BlockListEmails *bool   `json:"block_list_emails,omitempty"`
	UsedOn          *string `json:"used_on,omitempty"`
}

type DomainAddress struct {
	ID              int     `json:"id"`
	Address         string  `json:"address"`
	FullAddress     string  `json:"full_address"`
	Enabled         bool    `json:"enabled"`
	Description     string  `json:"description"`
	BlockListEmails bool    `json:"block_list_emails"`
	CreatedAt       string  `json:"created_at"`
	LastUsedAt      *string `json:"last_used_at"`
	NumForwarded    int     `json:"num_forwarded"`
	NumBlocked      int     `json:"num_blocked"`
	NumReplied      int     `json:"num_replied"`
	NumSpam         int     `json:"num_spam"`
}

type CreateDomainAddressRequest struct {
	Address         string `json:"address"`
	Enabled         bool   `json:"enabled"`
	Description     string `json:"description,omitempty"`
	BlockListEmails bool   `json:"block_list_emails"`
}

type UpdateDomainAddressRequest struct {
	Enabled         *bool   `json:"enabled,omitempty"`
	Description     *string `json:"description,omitempty"`
	BlockListEmails *bool   `json:"block_list_emails,omitempty"`
}

type RelayNumber struct {
	ID             int     `json:"id"`
	Number         string  `json:"number"`
	Enabled        bool    `json:"enabled"`
	Location       string  `json:"location"`
	VendorID       string  `json:"vendor_id"`
	CountryCode    string  `json:"country_code"`
	CreatedAt      *string `json:"created_at"`
	RemainingText  int     `json:"remaining_texts"`
	RemainingMin   int     `json:"remaining_minutes"`
	CallsForwarded int     `json:"calls_forwarded"`
	CallsBlocked   int     `json:"calls_blocked"`
	TextsForwarded int     `json:"texts_forwarded"`
	TextsBlocked   int     `json:"texts_blocked"`
}

type APIError struct {
	StatusCode int
	Body       string
}

func (e *APIError) Error() string {
	return e.Body
}
