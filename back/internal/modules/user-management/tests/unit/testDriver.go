package user_unit_test

import (
	"chat/internal/api/dto/requests"
	shared_infra "chat/internal/modules/shared/infra"
	user_management_app "chat/internal/modules/user-management/application"
	user_management_domain "chat/internal/modules/user-management/domain"
	user_management_infra "chat/internal/modules/user-management/infra"
	user_repos "chat/internal/modules/user-management/repos"
	"context"
	"time"

	"github.com/google/uuid"
)

var (
	FakeUuid       = uuid.MustParse("f1f4f647-02f1-47d8-988a-1f9788008b47")
	HashedPassword = "hashedPassword"
	existingEmail  = "existing@example.com"
	ExistingUser   = user_management_domain.NewUser(uuid.MustParse("d61142ad-ed2d-4db1-9c3f-4f0f7cf1c565"), existingEmail, "otherHashedPassword", time.Date(2025, time.November, 10, 23, 0, 0, 0, time.UTC))
	fakeToken      = "fakeToken"
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

	fakeUserRepo := user_repos.NewFakeUserRepo()
	fakeUserRepo.Users = []*user_management_domain.User{ExistingUser}

	userSrv := user_management_app.NewUserSrv(
		fakeUserRepo,
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
	td.fakeBcryp.ExpectedHash = HashedPassword
	td.fakeUuidGen.ExpectedUUID = FakeUuid

	ctx := context.Background()
	return td.UserSrv.CreateUserSrv(ctx, requests.SignupRequest{
		Email:    email,
		Password: password,
	})
}

func (td *TestDriver) LoginUser(email string, password string, expectedBcrypCompareResult bool) (user_management_app.LogResponse, error) {
	ctx := context.Background()
	td.fakeJwt.ExpectedToken = fakeToken
	td.fakeBcryp.ExpectedCompareResult = expectedBcrypCompareResult

	return td.UserSrv.LogUserSrv(ctx, requests.LoginRequest{
		Email:    email,
		Password: password,
	})
}
