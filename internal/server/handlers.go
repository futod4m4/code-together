package server

import (
	"fmt"
	authHttp "github.com/futod4m4/m/internal/auth/delivery/http"
	authRepository "github.com/futod4m4/m/internal/auth/repository"
	authUseCase "github.com/futod4m4/m/internal/auth/usecase"
	apiMiddlewares "github.com/futod4m4/m/internal/middleware"
	sessionRepository "github.com/futod4m4/m/internal/session/repository"
	"github.com/futod4m4/m/internal/session/usecase"
	"github.com/futod4m4/m/pkg/csrf"
	"github.com/futod4m4/m/pkg/metric"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

// Map Server Handlers
func (s *Server) MapHandlers(e *echo.Echo) error {
	metrics, err := metric.CreateMetrics(s.cfg.Metrics.URL, s.cfg.Metrics.ServiceName)
	if err != nil {
		s.logger.Errorf("CreateMetrics Error: %s", err)
	}
	s.logger.Info(
		"Metrics available URL: %s, ServiceName: %s",
		s.cfg.Metrics.URL,
		s.cfg.Metrics.ServiceName,
	)

	// Init repositories
	aRepo := authRepository.NewAuthRepository(s.db)
	sRepo := sessionRepository.NewSessionRepository(s.redisClient, s.cfg)
	authRedisRepo := authRepository.NewAuthRedisRepository(s.redisClient)

	// Init useCases
	authUC := authUseCase.NewAuthUseCase(s.cfg, aRepo, authRedisRepo, s.logger)
	sessUC := usecase.NewSessionUseCase(sRepo, s.cfg)

	// Init handlers
	authHandlers := authHttp.NewAuthHandlers(s.cfg, authUC, sessUC, s.logger)

	mw := apiMiddlewares.NewMiddlewareManager(sessUC, authUC, s.cfg, []string{"*"}, s.logger)

	e.Use(mw.RequestLoggerMiddleware)
	e.Use(mw.MetricsMiddleware(metrics))

	if s.cfg.Server.SSL {
		e.Pre(middleware.HTTPSRedirect())
	}

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderXRequestID, csrf.CSRFHeader},
	}))
	e.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{
		StackSize:         1 << 10, // 1 KB
		DisablePrintStack: true,
		DisableStackAll:   true,
	}))
	e.Use(middleware.RequestID())

	//e.Use(middleware.GzipWithConfig(middleware.GzipConfig{
	//	Level: 5,
	//	Skipper: func(c echo.Context) bool {
	//		return strings.Contains(c.Request().URL.Path, "swagger")
	//	},
	//}))
	e.Use(middleware.Secure())
	e.Use(middleware.BodyLimit("2M"))

	v1 := e.Group("/api/v1")

	authGroup := v1.Group("/auth")

	authHttp.MapAuthRoutes(authGroup, authHandlers, mw)

	fmt.Println("handlersMapped")

	return nil
}
