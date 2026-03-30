package files

import "github.com/labstack/echo"

type Handlers interface {
	CreateFile() echo.HandlerFunc
	UpdateFile() echo.HandlerFunc
	DeleteFile() echo.HandlerFunc
	GetFilesByRoomID() echo.HandlerFunc
}
