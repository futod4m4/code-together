package rooms

import "github.com/labstack/echo"

type Handlers interface {
	Create() echo.HandlerFunc
	Update() echo.HandlerFunc
	Delete() echo.HandlerFunc
	GetRoomByID() echo.HandlerFunc
	Join() echo.HandlerFunc
	Leave() echo.HandlerFunc
}
