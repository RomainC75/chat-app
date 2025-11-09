package user_unit_test

import (
	"chat/internal/api/dto/requests"
	shared_infra "chat/internal/modules/shared/infra"
	user_management_app "chat/internal/modules/user-management/application"
	user_management_infra "chat/internal/modules/user-management/infra"
	user_repos "chat/internal/modules/user-management/repos"
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

var (
	fakeUuid       = uuid.MustParse("f1f4f647-02f1-47d8-988a-1f9788008b47")
	hashedPassword = "hashedPassword"
)

type TestDriver struct {
	UserSrv     *user_management_app.UserSrv
	fakeBcryp   *user_management_infra.FakeBcrypt
	fakeJwt     *user_management_infra.FakeJWT
	fakeUuidGen *shared_infra.FakeUUIDGenerator
	fakeClock   *shared_infra.FakeClock
}

func NewTestDriver() *TestDriver {
	fakeBcrypt := user_management_infra.NewFakeBcrypt()
	fakeJwt := user_management_infra.NewFakeJWT()
	fakeUuidGenerator := shared_infra.NewFakeUUIDGenerator()
	fakeClock := shared_infra.NewFakeClock()

	userSrv := user_management_app.NewUserSrv(
		user_repos.NewFakeUserRepo(),
		fakeUuidGenerator,
		fakeClock,
		fakeBcrypt,
		fakeJwt,
	)
	return &TestDriver{
		UserSrv:     userSrv,
		fakeBcryp:   fakeBcrypt,
		fakeJwt:     fakeJwt,
		fakeUuidGen: fakeUuidGenerator,
		fakeClock:   fakeClock,
	}
}

func (td *TestDriver) CreateUser(email string, password string) (user_management_app.CreateUserResponse, error) {
	td.fakeBcryp.ExpectedHash = hashedPassword
	td.fakeUuidGen.ExpectedUUID = fakeUuid

	ctx := context.Background()
	return td.UserSrv.CreateUserSrv(ctx, requests.SignupRequest{
		Email:    email,
		Password: password,
	})
}

func TestUser(t *testing.T) {
	t.Run("Test Create User", func(t *testing.T) {
		td := NewTestDriver()
		email := "test@example.com"
		password := "azerty"

		resp, err := td.CreateUser(email, password)
		assert.Nil(t, err)

		assert.Equal(t, fakeUuid, resp.Id)
		assert.Equal(t, email, resp.Email)
	})
}
