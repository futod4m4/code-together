package http

import (
	"github.com/futod4m4/m/internal/middleware"
	"github.com/futod4m4/m/internal/rooms"
	"github.com/labstack/echo"
)

func MapRoomRoutes(roomGroup *echo.Group, r rooms.HttpHandlers, mw *middleware.MiddlewareManager) {
	roomGroup.POST("/create", r.Create(), mw.AuthSessionMiddleware)
	roomGroup.PUT("/:room_id", r.Update(), mw.AuthSessionMiddleware, mw.CSRF)
	roomGroup.DELETE("/:room_id", r.Delete(), mw.AuthSessionMiddleware, mw.CSRF)
	roomGroup.GET("/:room_id", r.GetRoomByID())
	roomGroup.GET("/code/:join_code", r.GetRoomByJoinCode())
}
