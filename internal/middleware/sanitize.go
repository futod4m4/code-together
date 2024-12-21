package middleware

import (
	"github.com/futod4m4/m/pkg/sanitize"
	"github.com/labstack/echo"
	"io"
	"net/http"
)

// Sanitize and read request body to ctx for next use in easy json
func (mw *MiddlewareManager) Sanitize(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		body, err := io.ReadAll(ctx.Request().Body)
		if err != nil {
			return ctx.NoContent(http.StatusBadRequest)
		}
		defer ctx.Request().Body.Close()

		sanBody, err := sanitize.SanitizeJSON(body)
		if err != nil {
			return ctx.NoContent(http.StatusBadRequest)
		}

		ctx.Set("body", sanBody)
		return next(ctx)
	}
}
