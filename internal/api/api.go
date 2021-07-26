package api

import (
	"jwtGoApi/pkg/config"
	"jwtGoApi/pkg/data/providers"
	"jwtGoApi/pkg/services"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.mongodb.org/mongo-driver/mongo"
)

type App struct {
	server      *echo.Echo
	userService services.IUserService
	cfg         *config.Settings
}

func New(settings *config.Settings, client *mongo.Client) *App {
	server := echo.New()

	//middleware which recovers from panics anywhere in the chain
	server.Use(middleware.Recover())
	server.Use(middleware.RequestID())

	userprovider := providers.NewUserProvider(settings, client)
	userSvc := services.NewUserService(settings, userprovider)

	return &App{
		server:      server,
		userService: userSvc,
		cfg:         settings,
	}
}

func (a App) ConfigureRoutes() {
	a.server.GET("/v1/public/healthy", a.HealthCheck)

	a.server.POST("/v1/public/account/login", a.Login)
	a.server.POST("/v1/public/account/register", a.Register)

	protected := a.server.Group("/v1/api")

	middleware := Middleware{config: a.cfg}
	protected.Use(middleware.Auth)
	protected.GET("/secret", func(c echo.Context) error {
		userId := c.Get("user").(string)
		return c.String(200, userId)
	}) 
	
}

func (a App) Start() {
	a.ConfigureRoutes()
	a.server.Start(":5000")
}
