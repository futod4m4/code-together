package roomCodes

import "github.com/labstack/echo"

type HttpHandlers interface {
	Create() echo.HandlerFunc
	Update() echo.HandlerFunc
	Delete() echo.HandlerFunc
	GetRoomCodeByID() echo.HandlerFunc
	GetRoomCodeByRoomID() echo.HandlerFunc
}

type WSHandlers interface {
	Compile() echo.HandlerFunc
}
