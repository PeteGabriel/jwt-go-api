package userservicemocks

import (
	"jwtGoApi/pkg/domain"
	"jwtGoApi/pkg/models"
)

type UserServiceMock struct {
	CreateAccountMock func(user *domain.User) *models.Error
	LoginMock         func(user *domain.User) (string, *models.Error)
}

func (us UserServiceMock) CreateAccount(user *domain.User) *models.Error {
	return us.CreateAccountMock(user)
}

func (us UserServiceMock) Login(user *domain.User) (string, *models.Error) {
	return us.LoginMock(user)
}
