package rooms

import "github.com/labstack/echo"

type HttpHandlers interface {
	Create() echo.HandlerFunc
	Update() echo.HandlerFunc
	Delete() echo.HandlerFunc
	GetRoomByID() echo.HandlerFunc
	GetRoomByJoinCode() echo.HandlerFunc
	GetMyRooms() echo.HandlerFunc
}

type WSHandlers interface {
	Join() echo.HandlerFunc
	Leave() echo.HandlerFunc
}
