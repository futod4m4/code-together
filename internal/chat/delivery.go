package chat

import "github.com/labstack/echo"

type Handlers interface {
	CreateMessage() echo.HandlerFunc
	GetMessages() echo.HandlerFunc
}
