package http

import (
	"github.com/futod4m4/m/internal/middleware"
	"github.com/futod4m4/m/internal/rooms"
	"github.com/labstack/echo"
)

func MapRoomRoutes(roomGroup *echo.Group, r rooms.Handlers, mw *middleware.MiddlewareManager) {
	roomGroup.POST("/create", r.Create(), mw.AuthSessionMiddleware, mw.CSRF)
	roomGroup.PATCH("/update", r.Update(), mw.AuthSessionMiddleware, mw.CSRF)
	roomGroup.DELETE("/delete", r.Delete(), mw.AuthSessionMiddleware, mw.CSRF)
	roomGroup.GET("/:room_id", r.GetRoomByID())
	roomGroup.GET("/join/:join_code/monaco-react-2", r.Join())
	roomGroup.GET("/leave", r.Leave(), mw.AuthSessionMiddleware, mw.CSRF)
}
