package middleware

import (
	"github.com/futod4m4/m/pkg/utils"
	"github.com/labstack/echo"
	"time"
)

// Request logger middleware
func (mw *MiddlewareManager) RequestLoggerMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		start := time.Now()
		err := next(ctx)

		req := ctx.Request()
		res := ctx.Response()
		status := res.Status
		size := res.Size
		s := time.Since(start).String()
		requestID := utils.GetRequestID(ctx)

		// Проверяем, не является ли ошибка nil, если да, то выводим "nil"
		errorMsg := "nil"
		if err != nil {
			errorMsg = err.Error() // Преобразуем ошибку в строку
		}

		mw.logger.Infof("RequestID: %s, Method: %s, URI: %s, Status: %v, Size: %v, Time: %s, Error: %s",
			requestID, req.Method, req.URL, status, size, s, errorMsg,
		)
		return err
	}
}
