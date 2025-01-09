package ws

import (
	"github.com/futod4m4/m/internal/middleware"
	"github.com/futod4m4/m/internal/rooms"
	"github.com/labstack/echo"
)

func MapRoomRoutes(roomGroup *echo.Group, r rooms.WSHandlers, mw *middleware.MiddlewareManager) {
	roomGroup.GET("/join/:join_code/monaco-react-2", r.Join())
	roomGroup.GET("/leave", r.Leave(), mw.AuthSessionMiddleware, mw.CSRF)
}
