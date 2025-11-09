package user_unit_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUser(t *testing.T) {
	t.Run("Test Create User", func(t *testing.T) {
		td := NewTestDriver()
		email := "test@example.com"
		password := "azerty"

		resp, err := td.CreateUser(email, password)
		assert.Nil(t, err)

		assert.Equal(t, FakeUuid, resp.Id)
		assert.Equal(t, email, resp.Email)
	})
	t.Run("Return error if user exists", func(t *testing.T) {
		td := NewTestDriver()
		_, err := td.CreateUser(existingEmail, "qsdf")
		assert.NotNil(t, err)
		assert.EqualError(t, err, "email already used")
	})

}
