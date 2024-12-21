package http

import (
	"github.com/futod4m4/m/internal/auth"
	"github.com/futod4m4/m/internal/middleware"
	"github.com/labstack/echo"
)

func MapAuthRoutes(authGroup *echo.Group, h auth.Handlers, mw *middleware.MiddlewareManager) {
	authGroup.POST("/register", h.Register())
	authGroup.POST("/login", h.Login())
	authGroup.POST("/logout", h.Logout())
	authGroup.GET("/:user_id", h.GetUserByID())
	//authGroup.Use(middleware.AuthJWTMiddleware(authUC, cfg))
	authGroup.Use(mw.AuthSessionMiddleware)
	authGroup.GET("/me", h.GetMe())
	authGroup.GET("/token", h.GetCSRFToken())
}
