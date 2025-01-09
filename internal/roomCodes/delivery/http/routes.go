package http

import (
	"github.com/futod4m4/m/internal/middleware"
	"github.com/futod4m4/m/internal/roomCodes"
	"github.com/labstack/echo"
)

func MapRoomRoutes(roomCodeGroup *echo.Group, r roomCodes.HttpHandlers, mw *middleware.MiddlewareManager) {
	roomCodeGroup.POST("/create", r.Create())
	roomCodeGroup.PUT("/:room_code_id", r.Update())
	roomCodeGroup.DELETE("/:room_code_id", r.Delete())
	roomCodeGroup.GET("/:room_code_id", r.GetRoomCodeByID())
	roomCodeGroup.GET("/code/:room_id", r.GetRoomCodeByRoomID())
}
