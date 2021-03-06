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

func TestLogin_Ok(t *testing.T) {

	mock := &userservicemocks.UserServiceMock{}
	mock.LoginMock = func(user *domain.User) (string, *models.Error) {
		return "some_valid_token", nil
	}

	e := echo.New()
	app := App{
		server:      e,
		cfg:         &config.Settings{},
		userService: mock,
	}

	credentials := `{"username": "Test2", "password": "12345Test"}`
	req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(credentials))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	loginErr := app.Login(e.NewContext(req, rec))
	assert.Equal(t, nil, loginErr)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "{\"Token\":\"some_valid_token\"}\n", rec.Body.String())
}

func TestLogin_UserNotFound(t *testing.T) {

	mock := &userservicemocks.UserServiceMock{}
	mock.LoginMock = func(user *domain.User) (string, *models.Error) {
		return "", &models.Error{
			Message: "your username and/or password was not found",
			Code:    400,
			Name:    "INVALID_LOGIN",
		}
	}

	e := echo.New()
	app := App{
		server:      e,
		cfg:         &config.Settings{},
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

func TestLogin_PasswordNotValid(t *testing.T) {
	mock := &userservicemocks.UserServiceMock{}
	mock.LoginMock = func(user *domain.User) (string, *models.Error) {
		return "", &models.Error{
			Code:    400,
			Message: "A validation error occurred.",
			Name:    "VALIDATION",
		}
	}

	e := echo.New()
	app := App{
		server:      e,
		cfg:         &config.Settings{},
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

func TestLogin_UsernameNotValid(t *testing.T) {
	mock := &userservicemocks.UserServiceMock{}
	mock.LoginMock = func(user *domain.User) (string, *models.Error) {
		return "", &models.Error{
			Code:    400,
			Message: "A validation error occurred.",
			Name:    "VALIDATION",
		}
	}

	e := echo.New()
	app := App{
		server:      e,
		cfg:         &config.Settings{},
		userService: mock,
	}

	credentials := `{"username": "T", "password": "89wqrBasfa2"}`
	req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(credentials))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	loginErr := app.Login(e.NewContext(req, rec))
	assert.Equal(t, nil, loginErr)

	assert.Equal(t, http.StatusBadRequest, rec.Code)

	expBody := "{\"message\":\"A validation error occurred.\",\"code\":400,\"name\":\"VALIDATION\",\"validation\":[\"Username must be longer than 2 characters\"]}\n"
	assert.Equal(t, expBody, rec.Body.String())
}

func TestRegister_Ok(t *testing.T) {

	e := echo.New()
	app := App{
		server: e,
		cfg:    &config.Settings{},
		userService: &userservicemocks.UserServiceMock{
			CreateAccountMock: func(user *domain.User) *models.Error {
				return nil
			},
		},
	}

	credentials := `{"username": "Test1", "password": "TestNumber1"}`
	req := httptest.NewRequest(http.MethodPost, "/register", strings.NewReader(credentials))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	registerRequest := app.Register(e.NewContext(req, rec))
	assert.Nil(t, registerRequest)

	assert.Equal(t, http.StatusCreated, rec.Code)
}

func TestRegister_UsernameNotValid(t *testing.T) {

	e := echo.New()
	app := App{
		server: e,
		cfg:    &config.Settings{},
		userService: &userservicemocks.UserServiceMock{
			CreateAccountMock: func(user *domain.User) *models.Error {
				return nil
			},
		},
	}

	credentials := `{"username": "1", "password": "TestNumber1"}`
	req := httptest.NewRequest(http.MethodPost, "/register", strings.NewReader(credentials))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	registerRequest := app.Register(e.NewContext(req, rec))
	assert.Nil(t, registerRequest)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	expBody := "{\"message\":\"A validation error occurred.\",\"code\":400,\"name\":\"VALIDATION\",\"validation\":[\"Username must be longer than 2 characters\"]}\n"
	assert.Equal(t, expBody, rec.Body.String())
}

func TestRegister_PasswordNotValid(t *testing.T) {

	e := echo.New()
	app := App{
		server: e,
		cfg:    &config.Settings{},
		userService: &userservicemocks.UserServiceMock{
			CreateAccountMock: func(user *domain.User) *models.Error {
				return nil
			},
		},
	}

	credentials := `{"username": "Test3", "password": "1"}`
	req := httptest.NewRequest(http.MethodPost, "/register", strings.NewReader(credentials))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	registerRequest := app.Register(e.NewContext(req, rec))
	assert.Nil(t, registerRequest)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	expBody := "{\"message\":\"A validation error occurred.\",\"code\":400,\"name\":\"VALIDATION\",\"validation\":[\"Password must be 8 characters\"]}\n"
	assert.Equal(t, expBody, rec.Body.String())
}

func TestRegister_AccountAlreadyExists(t *testing.T) {
	e := echo.New()
	app := App{
		server: e,
		cfg:    &config.Settings{},
		userService: &userservicemocks.UserServiceMock{
			CreateAccountMock: func(user *domain.User) *models.Error {
				return &models.Error{
					Message: "account/user already exists with the same username",
					Code:    400,
					Name:    "USERNAME_TAKEN",
				}
			},
		},
	}

  credentials := `{"username": "Test3", "password": "Password123"}`
	req := httptest.NewRequest(http.MethodPost, "/register", strings.NewReader(credentials))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	registerRequest := app.Register(e.NewContext(req, rec))
	assert.Nil(t, registerRequest)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	expBody := "{\"message\":\"account/user already exists with the same username\",\"code\":400,\"name\":\"USERNAME_TAKEN\"}\n"
	assert.Equal(t, expBody, rec.Body.String())
}