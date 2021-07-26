package services

import (
	"errors"
	"jwtGoApi/pkg/config"
	"jwtGoApi/pkg/domain"
	usermocks "jwtGoApi/pkg/mocks/data/users"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateAccount_UsernameAlreadyExists(t *testing.T){
	cfg := &config.Settings{}
	providerMock := &usermocks.UserProviderMock{}
	providerMock.UsernameExistsMock = func(username string) (bool, error) {
		return true, nil
	}

	svc := NewUserService(cfg, providerMock)
	newUser := &domain.User{Username: "Test", Password: "12345Test"}

	response := svc.CreateAccount(newUser)
	assert.Equal(t, "USERNAME_TAKEN", response.Name)
}

func TestCreateAccount_ServerErrorThrown(t *testing.T){
	cfg := &config.Settings{}
	providerMock := &usermocks.UserProviderMock{}
	providerMock.UsernameExistsMock = func(username string) (bool, error) {
		return false, nil
	}
	providerMock.CreateAccountMock = func(user *domain.User) error {
		return errors.New("Error inserting user")
	}

	svc := NewUserService(cfg, providerMock)
	newUser := &domain.User{Username: "Test2", Password: "12345Test"}

	response := svc.CreateAccount(newUser)
	assert.Equal(t, "SERVER_ERROR", response.Name)
	msg := "error happened while trying to create a new account."
	assert.Equal(t, msg, response.Message)
	assert.Equal(t, 500, response.Code)
}

func TestCreateAccount_UsernameSearchThrowsError(t *testing.T){
	cfg := &config.Settings{}
	providerMock := &usermocks.UserProviderMock{}
	providerMock.UsernameExistsMock = func(username string) (bool, error) {
		return false, errors.New("username search found an error")
	}

	svc := NewUserService(cfg, providerMock)
	newUser := &domain.User{Username: "Test2", Password: "12345Test"}

	response := svc.CreateAccount(newUser)
	assert.Equal(t, "SERVER_ERROR", response.Name)
	msg := "error happened while trying to check if account/user already exists"
	assert.Equal(t, msg, response.Message)
	assert.Equal(t, 500, response.Code)
}