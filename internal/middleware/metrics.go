package middleware

import (
	"github.com/futod4m4/m/pkg/metric"
	"github.com/labstack/echo"
	"github.com/pkg/errors"
	"time"
)

// Prometheus metrics middleware
func (mw *MiddlewareManager) MetricsMiddleware(metrics metric.Metrics) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()
			err := next(c)
			var status int
			if err != nil {
				var HTTPError *echo.HTTPError
				errors.As(err, &HTTPError)
				status = HTTPError.Code
			} else {
				status = c.Response().Status
			}
			metrics.ObserveResponseTime(status, c.Request().Method, c.Path(), time.Since(start).Seconds())
			metrics.IncHits(status, c.Request().Method, c.Path())
			return err
		}
	}
}
