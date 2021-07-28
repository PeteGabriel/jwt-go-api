package models

import (
	"jwtGoApi/pkg/domain"

	"github.com/labstack/echo/v4"
)



type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}


type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

const (
	passwordValidationError = "Password must be 8 characters"
	usernameValidationError = "Username must be longer than 2 characters"
)

func ValidateRegisterRequest(c echo.Context) (*domain.User, *Error){

	regRequest := new(RegisterRequest)
	//echo provides this in order to validate expected data automatically
	if err := c.Bind(regRequest); err != nil {
		return nil, BindError()
	}

	var validationErrors []string

	if len(regRequest.Password) < 8 {
		validationErrors = append(validationErrors, passwordValidationError)
	}

	if len(regRequest.Username) < 3 {
		validationErrors = append(validationErrors, usernameValidationError)
	}

	if len(validationErrors) > 0 {
		return nil, ValidationError(validationErrors)
	}

	return &domain.User{
		Username: regRequest.Username,
		Password: regRequest.Password,
	}, nil
}

func ValidateLoginRequest(c echo.Context) (*domain.User, *Error){

	logRequest := new(LoginRequest)
	//echo provides this in order to validate expected data automatically
	//will unmarshal the request body into a struct, if that fails we return the BindError function created previously.
	if err := c.Bind(logRequest); err != nil {
		return nil, BindError()
	}

	var validationErrors []string

	if len(logRequest.Password) < 8 {
		validationErrors = append(validationErrors, passwordValidationError)
	}

	if len(logRequest.Username) < 3 {
		validationErrors = append(validationErrors, usernameValidationError)
	}

	if len(validationErrors) > 0 {
		return nil, ValidationError(validationErrors)
	}

	return &domain.User{
		Username: logRequest.Username,
		Password: logRequest.Password,
	}, nil
}