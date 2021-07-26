package api

import (
	"fmt"
	"jwtGoApi/pkg/config"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
)

type Middleware struct {
	config *config.Settings
}

//Auth is a middleware that performs a validation
//over the sent JWT token. If the token is ok, we set the 
//user ID into the context and move on with the pipeline of handlers.
func (m Middleware) Auth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {

		//request must send in the header the assigned token
		token := c.Request().Header.Get("Authorization")
		if token == ""{
			return echo.ErrUnauthorized
		}

		type Claims struct {
			Id string `json:"id"`
			Exp int `json:"exp"`
			jwt.StandardClaims 
		}

		keyFunc := func(token *jwt.Token)(interface{}, error) {
			return []byte(m.config.JwtSecret), nil
		}
		//from given jwt token, parse it into a valid Claims structure.
		//otherwhise jwt is invalid.
		jwtToken, err := jwt.ParseWithClaims(token, &Claims{}, keyFunc)
		if err != nil {
			fmt.Println(err)
			return echo.ErrUnauthorized
		}

		claims, ok := jwtToken.Claims.(*Claims)
		if ok && jwtToken.Valid {
			c.Set("user", claims.Id)
			return next(c)
		}else {
			return echo.ErrUnauthorized
		}
	}
}
