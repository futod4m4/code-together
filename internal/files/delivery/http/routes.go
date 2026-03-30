package http

import (
	"github.com/futod4m4/m/internal/files"
	"github.com/labstack/echo"
)

func MapFileRoutes(fileGroup *echo.Group, h files.Handlers) {
	fileGroup.POST("/create", h.CreateFile())
	fileGroup.PUT("/:file_id", h.UpdateFile())
	fileGroup.DELETE("/:file_id", h.DeleteFile())
	fileGroup.GET("/room/:room_id", h.GetFilesByRoomID())
}
