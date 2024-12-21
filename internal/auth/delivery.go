package auth

import "github.com/labstack/echo"

type Handlers interface {
	Register() echo.HandlerFunc
	Login() echo.HandlerFunc
	Logout() echo.HandlerFunc
	Update() echo.HandlerFunc
	Delete() echo.HandlerFunc
	GetUserByID() echo.HandlerFunc
	GetMe() echo.HandlerFunc
	GetCSRFToken() echo.HandlerFunc
}
