package api

import (
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

func TestListUsers(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	client := NewClient("test")

	t.Run("successful list users", func(t *testing.T) {
		responder := httpmock.NewJsonResponderOrPanic(200, []User{
			{Email: "ffrelayctl@domain.tld"},
			{Email: "test@example.com"},
		})
		httpmock.RegisterResponder("GET", DefaultBaseURL+APIBasePath+"users/", responder)

		users, err := client.ListUsers()

		assert.NoError(t, err)
		assert.Len(t, users, 2)
		assert.Equal(t, "ffrelayctl@domain.tld", users[0].Email)
		assert.Equal(t, "test@example.com", users[1].Email)
	})

	t.Run("empty list", func(t *testing.T) {
		responder := httpmock.NewJsonResponderOrPanic(200, []User{})
		httpmock.RegisterResponder("GET", DefaultBaseURL+APIBasePath+"users/", responder)

		users, err := client.ListUsers()

		assert.NoError(t, err)
		assert.Empty(t, users)
	})

	t.Run("error response", func(t *testing.T) {
		responder := httpmock.NewStringResponder(500, `{"error": "internal server error"}`)
		httpmock.RegisterResponder("GET", DefaultBaseURL+APIBasePath+"users/", responder)

		users, err := client.ListUsers()

		assert.Error(t, err)
		assert.Nil(t, users)
	})
}
