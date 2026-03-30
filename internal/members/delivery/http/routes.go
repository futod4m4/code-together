package http

import (
	"github.com/futod4m4/m/internal/members"
	"github.com/futod4m4/m/internal/middleware"
	"github.com/labstack/echo"
)

func MapMemberRoutes(memberGroup *echo.Group, h members.Handlers, mw *middleware.MiddlewareManager) {
	memberGroup.Use(mw.AuthSessionMiddleware)
	memberGroup.POST("/add", h.AddMember())
	memberGroup.PUT("/:room_id/role", h.UpdateRole())
	memberGroup.DELETE("/:room_id/:user_id", h.RemoveMember())
	memberGroup.GET("/:room_id", h.GetMembers())
}
