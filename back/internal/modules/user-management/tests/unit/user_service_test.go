package user_unit_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUser(t *testing.T) {
	t.Run("Signup : Create User", func(t *testing.T) {
		td := NewTestDriver()
		email := "test@example.com"
		password := "azerty"

		resp, err := td.CreateUser(email, password)
		assert.Nil(t, err)

		assert.Equal(t, FakeUuid, resp.Id)
		assert.Equal(t, email, resp.Email)
	})
	t.Run("Signup : Return error if user exists", func(t *testing.T) {
		td := NewTestDriver()
		_, err := td.CreateUser(existingEmail, "qsdf")
		assert.NotNil(t, err)
		assert.EqualError(t, err, "email already used")
	})

	t.Run("Login User", func(t *testing.T) {
		td := NewTestDriver()
		res, err := td.LoginUser(existingEmail, "azerty", true)
		assert.Nil(t, err)
		assert.Equal(t, ExistingUser.GetID(), res.Id)
		assert.Equal(t, fakeToken, res.Token)
	})

	t.Run("Login User: good email & wrong password", func(t *testing.T) {
		td := NewTestDriver()
		_, err := td.LoginUser(existingEmail, "azerty", false)
		assert.NotNil(t, err)
		assert.EqualError(t, err, "wrong email or password")

	})
	t.Run("Login User: wrong email", func(t *testing.T) {
		td := NewTestDriver()
		_, err := td.LoginUser("wrong@example.com", "azerty", true)
		assert.NotNil(t, err)
		assert.EqualError(t, err, "wrong email or password")
	})
}
