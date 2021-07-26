package api

import (
	"jwtGoApi/pkg/models"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (app App) Login(c echo.Context) error {
	user, err := models.ValidateLoginRequest(c)
	if err != nil {
		return c.JSON(err.Code, err)
	}

	token, err := app.userService.Login(user)
	if err != nil {
		return c.JSON(err.Code, err)
	}

	resp := &models.LoginResponse{Token: token}
	return c.JSON(http.StatusOK, resp)
}

func (app App) Register(c echo.Context) error {
	user, err := models.ValidateRegisterRequest(c)
	if err != nil {
		return c.JSON(err.Code, err)
	}

	err = app.userService.CreateAccount(user)
	if err != nil {
		return c.JSON(err.Code, err)
	}

	return c.JSON(http.StatusCreated, "")
}