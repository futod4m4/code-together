package members

import "github.com/labstack/echo"

type Handlers interface {
	AddMember() echo.HandlerFunc
	UpdateRole() echo.HandlerFunc
	RemoveMember() echo.HandlerFunc
	GetMembers() echo.HandlerFunc
}
