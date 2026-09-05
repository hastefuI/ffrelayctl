package api

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
)

// Profile is a Firefox Relay account profile. It carries the subscription
// state and the counters totalled across every mask on the account. APIToken
// is a live credential for the whole account.
type Profile struct {
	ID                          int          `json:"id"`
	ServerStorage               bool         `json:"server_storage"`
	StorePhoneLog               bool         `json:"store_phone_log"`
	Subdomain                   *string      `json:"subdomain"`
	HasPremium                  bool         `json:"has_premium"`
	HasPhone                    bool         `json:"has_phone"`
	HasVPN                      bool         `json:"has_vpn"`
	HasMegabundle               bool         `json:"has_megabundle"`
	OnboardingState             int          `json:"onboarding_state"`
	OnboardingFreeState         int          `json:"onboarding_free_state"`
	DatePhoneRegistered         *string      `json:"date_phone_registered"`
	DateSubscribed              *string      `json:"date_subscribed"`
	AvatarURL                   string       `json:"avatar"`
	NextEmailTry                string       `json:"next_email_try"`
	BounceStatus                BounceStatus `json:"bounce_status"`
	APIToken                    string       `json:"api_token,omitempty"`
	EmailsBlocked               int          `json:"emails_blocked"`
	EmailsForwarded             int          `json:"emails_forwarded"`
	EmailsReplied               int          `json:"emails_replied"`
	LevelOneTrackersBlocked     int          `json:"level_one_trackers_blocked"`
	RemoveLevelOneEmailTrackers bool         `json:"remove_level_one_email_trackers"`
	TotalMasks                  int          `json:"total_masks"`
	AtMaskLimit                 bool         `json:"at_mask_limit"`
	MetricsEnabled              bool         `json:"metrics_enabled"`
}

// BounceStatus reports whether email forwarding is paused because a message
// bounced, and which kind of bounce paused it. Relay sends the pair as a two
// element array rather than an object.
type BounceStatus struct {
	Paused bool
	Type   string
}

// UnmarshalJSON decodes the paused and type pair Relay sends as an array.
func (b *BounceStatus) UnmarshalJSON(data []byte) error {
	var tuple []any
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

// MarshalJSON encodes the paused and type pair back into the array Relay expects.
func (b BounceStatus) MarshalJSON() ([]byte, error) {
	return json.Marshal([]any{b.Paused, b.Type})
}

// RelayAddress is a random email mask and its forwarding counters.
type RelayAddress struct {
	ID                         int     `json:"id"`
	MaskType                   string  `json:"mask_type"`
	Address                    string  `json:"address"`
	Domain                     int     `json:"domain"`
	FullAddress                string  `json:"full_address"`
	Enabled                    bool    `json:"enabled"`
	Description                string  `json:"description"`
	GeneratedFor               string  `json:"generated_for"`
	UsedOn                     string  `json:"used_on"`
	BlockListEmails            bool    `json:"block_list_emails"`
	CreatedAt                  string  `json:"created_at"`
	LastModifiedAt             string  `json:"last_modified_at"`
	LastUsedAt                 *string `json:"last_used_at"`
	NumForwarded               int     `json:"num_forwarded"`
	NumBlocked                 int     `json:"num_blocked"`
	NumLevelOneTrackersBlocked *int    `json:"num_level_one_trackers_blocked"`
	NumReplied                 int     `json:"num_replied"`
	NumSpam                    int     `json:"num_spam"`
}

// CreateRelayAddressRequest is the body for Client.CreateRelayAddress.
type CreateRelayAddressRequest struct {
	Enabled         bool   `json:"enabled"`
	Description     string `json:"description,omitempty"`
	GeneratedFor    string `json:"generated_for,omitempty"`
	BlockListEmails bool   `json:"block_list_emails"`
	UsedOn          string `json:"used_on,omitempty"`
}

// UpdateRelayAddressRequest is the body for Client.UpdateRelayAddress.
// A nil field leaves that value unchanged.
type UpdateRelayAddressRequest struct {
	Enabled         *bool   `json:"enabled,omitempty"`
	Description     *string `json:"description,omitempty"`
	GeneratedFor    *string `json:"generated_for,omitempty"`
	BlockListEmails *bool   `json:"block_list_emails,omitempty"`
	UsedOn          *string `json:"used_on,omitempty"`
}

// DomainAddress is an email mask on the account's own subdomain, which is a
// premium feature, and its forwarding counters.
type DomainAddress struct {
	ID                         int     `json:"id"`
	MaskType                   string  `json:"mask_type"`
	Address                    string  `json:"address"`
	Domain                     int     `json:"domain"`
	FullAddress                string  `json:"full_address"`
	Enabled                    bool    `json:"enabled"`
	Description                string  `json:"description"`
	UsedOn                     string  `json:"used_on"`
	BlockListEmails            bool    `json:"block_list_emails"`
	CreatedAt                  string  `json:"created_at"`
	LastModifiedAt             string  `json:"last_modified_at"`
	LastUsedAt                 *string `json:"last_used_at"`
	NumForwarded               int     `json:"num_forwarded"`
	NumBlocked                 int     `json:"num_blocked"`
	NumLevelOneTrackersBlocked *int    `json:"num_level_one_trackers_blocked"`
	NumReplied                 int     `json:"num_replied"`
	NumSpam                    int     `json:"num_spam"`
}

// CreateDomainAddressRequest is the body for Client.CreateDomainAddress.
type CreateDomainAddressRequest struct {
	Address         string `json:"address"`
	Enabled         bool   `json:"enabled"`
	Description     string `json:"description,omitempty"`
	BlockListEmails bool   `json:"block_list_emails"`
	UsedOn          string `json:"used_on,omitempty"`
}

// UpdateDomainAddressRequest is the body for Client.UpdateDomainAddress.
// A nil field leaves that value unchanged.
type UpdateDomainAddressRequest struct {
	Enabled         *bool   `json:"enabled,omitempty"`
	Description     *string `json:"description,omitempty"`
	BlockListEmails *bool   `json:"block_list_emails,omitempty"`
	UsedOn          *string `json:"used_on,omitempty"`
}

// MaskFilter narrows a mask listing to the masks that match every field set on
// it. The zero value matches every mask. Relay matches UsedOn as a
// case-insensitive substring and every other field exactly. GeneratedFor
// applies to random masks only.
type MaskFilter struct {
	Enabled         *bool
	BlockListEmails *bool
	Address         string
	Description     string
	GeneratedFor    string
	UsedOn          string
}

// query encodes the filter as a query string prefixed with "?", or an empty
// string when no field is set. generatedFor reports whether the endpoint has a
// generated_for field to filter on.
func (f MaskFilter) query(generatedFor bool) string {
	values := url.Values{}
	if f.Enabled != nil {
		values.Set("enabled", strconv.FormatBool(*f.Enabled))
	}
	if f.BlockListEmails != nil {
		values.Set("block_list_emails", strconv.FormatBool(*f.BlockListEmails))
	}
	if f.Address != "" {
		values.Set("address", f.Address)
	}
	if f.Description != "" {
		values.Set("description", f.Description)
	}
	if generatedFor && f.GeneratedFor != "" {
		values.Set("generated_for", f.GeneratedFor)
	}
	if f.UsedOn != "" {
		values.Set("used_on", f.UsedOn)
	}

	if len(values) == 0 {
		return ""
	}
	return "?" + values.Encode()
}

// RelayNumber is a phone mask, with its call and text counters and the quota
// left for the current month.
type RelayNumber struct {
	ID                     int     `json:"id"`
	Number                 string  `json:"number"`
	Enabled                bool    `json:"enabled"`
	Location               string  `json:"location"`
	CountryCode            string  `json:"country_code"`
	CreatedAt              *string `json:"created_at"`
	RemainingText          int     `json:"remaining_texts"`
	RemainingMin           int     `json:"remaining_minutes"`
	CallsForwarded         int     `json:"calls_forwarded"`
	CallsBlocked           int     `json:"calls_blocked"`
	TextsForwarded         int     `json:"texts_forwarded"`
	TextsBlocked           int     `json:"texts_blocked"`
	CallsAndTextsForwarded int     `json:"calls_and_texts_forwarded"`
	CallsAndTextsBlocked   int     `json:"calls_and_texts_blocked"`
}

// UpdateRelayNumberRequest is the body for Client.UpdateRelayNumber.
// A nil field leaves that value unchanged.
type UpdateRelayNumberRequest struct {
	Enabled *bool `json:"enabled,omitempty"`
}

// PhoneNumberOption is a phone number that is free to be claimed as a mask.
type PhoneNumberOption struct {
	FriendlyName string  `json:"friendly_name"`
	ISOCountry   string  `json:"iso_country"`
	Locality     *string `json:"locality"`
	PhoneNumber  string  `json:"phone_number"`
	PostalCode   *string `json:"postal_code"`
	Region       string  `json:"region"`
}

// RelayNumberSuggestions groups the numbers Relay offers for a new phone mask
// by how closely each one matches the real number.
type RelayNumberSuggestions struct {
	RealNum           *string             `json:"real_num"`
	SamePrefixOptions []PhoneNumberOption `json:"same_prefix_options"`
	OtherAreasOptions []PhoneNumberOption `json:"other_areas_options"`
	SameAreaOptions   []PhoneNumberOption `json:"same_area_options"`
	RandomOptions     []PhoneNumberOption `json:"random_options"`
}

// InboundContact is a number that has called or texted a phone mask, with the
// counters kept for that one contact and whether it is blocked.
type InboundContact struct {
	ID              int     `json:"id"`
	RelayNumber     int     `json:"relay_number"`
	InboundNumber   string  `json:"inbound_number"`
	LastInboundDate string  `json:"last_inbound_date"`
	LastInboundType string  `json:"last_inbound_type"`
	NumCalls        int     `json:"num_calls"`
	NumCallsBlocked int     `json:"num_calls_blocked"`
	LastCallDate    *string `json:"last_call_date"`
	NumTexts        int     `json:"num_texts"`
	NumTextsBlocked int     `json:"num_texts_blocked"`
	LastTextDate    *string `json:"last_text_date"`
	Blocked         bool    `json:"blocked"`
}

// UpdateInboundContactRequest is the body for Client.UpdateInboundContact.
// A nil field leaves that value unchanged.
type UpdateInboundContactRequest struct {
	Blocked *bool `json:"blocked,omitempty"`
}

// User is the account the API token belongs to.
type User struct {
	Email string `json:"email"`
}

// RealPhone is a phone number registered on the account, which is the number a
// phone mask forwards to once it is verified.
type RealPhone struct {
	ID                   int     `json:"id"`
	Number               string  `json:"number"`
	VerificationSentDate *string `json:"verification_sent_date"`
	Verified             bool    `json:"verified"`
	VerifiedDate         *string `json:"verified_date"`
	CountryCode          string  `json:"country_code"`
}

// RegisterRealPhoneRequest is the body for Client.RegisterRealPhone.
type RegisterRealPhoneRequest struct {
	Number string `json:"number"`
}

// VerifyRealPhoneRequest is the body for Client.VerifyRealPhone. It carries
// the code Relay sent by text.
type VerifyRealPhoneRequest struct {
	Number           string `json:"number"`
	VerificationCode string `json:"verification_code"`
}

// APIError is returned when Relay answers with a status of 400 or above.
type APIError struct {
	StatusCode int
	Body       string
}

// Error returns the raw response body.
func (e *APIError) Error() string {
	return e.Body
}
