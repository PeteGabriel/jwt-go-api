package api

import (
	"jwtGoApi/pkg/config"
	"jwtGoApi/pkg/domain"
	userservicemocks "jwtGoApi/pkg/mocks/data/services"
	"jwtGoApi/pkg/models"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)


func TestLogin_Ok(t *testing.T){
	
	mock := &userservicemocks.UserServiceMock{}
	mock.LoginMock = func(user *domain.User) (string, *models.Error){
		return "some_valid_token", nil
	}

	e := echo.New()
	app := App{
		server: e,
		cfg: &config.Settings{},
		userService: mock,
	}

	credentials := `{"username": "Test2", "password": "12345Test"}`
	req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(credentials))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	loginErr := app.Login(e.NewContext(req, rec))
	assert.Equal(t, nil, loginErr)

	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestLogin_UserNotFound(t *testing.T){
	
	mock := &userservicemocks.UserServiceMock{}
	mock.LoginMock = func(user *domain.User) (string, *models.Error){
		return "", &models.Error{
			Message: "your username and/or password was not found",
			Code:    400,
			Name:    "INVALID_LOGIN",
		}
	}

	e := echo.New()
	app := App{
		server: e,
		cfg: &config.Settings{},
		userService: mock,
	}

	credentials := `{"username": "Test2", "password": "12345Test"}`
	req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(credentials))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	loginErr := app.Login(e.NewContext(req, rec))
	assert.Equal(t, nil, loginErr)
	
	assert.Equal(t, http.StatusBadRequest, rec.Code)

	expBody := "{\"message\":\"your username and/or password was not found\",\"code\":400,\"name\":\"INVALID_LOGIN\"}\n"
	assert.Equal(t, expBody, rec.Body.String())
}

func TestLogin_PasswordNotValid(t *testing.T){
	mock := &userservicemocks.UserServiceMock{}
	mock.LoginMock = func(user *domain.User) (string, *models.Error){
		return "", &models.Error{
			Code: 400,
			Message: "A validation error occurred.",
			Name: "VALIDATION",
		}
	}

	e := echo.New()
	app := App{
		server: e,
		cfg: &config.Settings{},
		userService: mock,
	}

	credentials := `{"username": "Test2", "password": "89"}`
	req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(credentials))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

loginErr := app.Login(e.NewContext(req, rec))
	assert.Equal(t, nil, loginErr)
	
	assert.Equal(t, http.StatusBadRequest, rec.Code)

	expBody := "{\"message\":\"A validation error occurred.\",\"code\":400,\"name\":\"VALIDATION\",\"validation\":[\"Password must be 8 characters\"]}\n"
	assert.Equal(t, expBody, rec.Body.String())
}