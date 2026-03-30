package http

import (
	"github.com/futod4m4/m/internal/chat"
	"github.com/futod4m4/m/internal/middleware"
	"github.com/labstack/echo"
)

func MapChatRoutes(chatGroup *echo.Group, h chat.Handlers, mw *middleware.MiddlewareManager) {
	chatGroup.POST("/send", h.CreateMessage(), mw.AuthSessionMiddleware)
	chatGroup.GET("/:room_id", h.GetMessages())
}
