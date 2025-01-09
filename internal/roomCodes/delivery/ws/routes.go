package ws

import (
	"github.com/futod4m4/m/internal/middleware"
	"github.com/futod4m4/m/internal/roomCodes"
	"github.com/labstack/echo"
)

func MapRoomRoutes(roomCodeGroup *echo.Group, r roomCodes.WSHandlers, mw *middleware.MiddlewareManager) {
	roomCodeGroup.POST("/compile", r.Compile())
}
