package server

import (
	authHttp "github.com/futod4m4/m/internal/auth/delivery/http"
	authRepository "github.com/futod4m4/m/internal/auth/repository"
	authUseCase "github.com/futod4m4/m/internal/auth/usecase"
	apiMiddlewares "github.com/futod4m4/m/internal/middleware"
	roomCodeHttp "github.com/futod4m4/m/internal/roomCodes/delivery/http"
	roomCodeRepository "github.com/futod4m4/m/internal/roomCodes/repository"
	roomCodeUseCase "github.com/futod4m4/m/internal/roomCodes/usecase"
	roomHttp "github.com/futod4m4/m/internal/rooms/delivery/http"
	roomWS "github.com/futod4m4/m/internal/rooms/delivery/ws"
	roomRepository "github.com/futod4m4/m/internal/rooms/repository"
	roomUseCase "github.com/futod4m4/m/internal/rooms/usecase"
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
	rRepo := roomRepository.NewRoomRepository(s.db)
	rcRepo := roomCodeRepository.NewRoomCodeRepository(s.db)
	sRepo := sessionRepository.NewSessionRepository(s.redisClient, s.cfg)
	authRedisRepo := authRepository.NewAuthRedisRepository(s.redisClient)
	roomRedisRepo := roomRepository.NewRoomRedisRepo(s.redisClient)
	roomCodeRedisRepo := roomCodeRepository.NewRoomCodeRedisRepo(s.redisClient)

	// Init useCases
	authUC := authUseCase.NewAuthUseCase(s.cfg, aRepo, authRedisRepo, s.logger)
	roomUC := roomUseCase.NewRoomUseCase(s.cfg, rRepo, roomRedisRepo, s.logger)
	roomCodeUC := roomCodeUseCase.NewRoomCodeUseCase(s.cfg, rcRepo, roomCodeRedisRepo, s.logger)
	sessUC := usecase.NewSessionUseCase(sRepo, s.cfg)

	// Init handlers
	authHandlers := authHttp.NewAuthHandlers(s.cfg, authUC, sessUC, s.logger)
	roomHandlers := roomHttp.NewRoomHandlers(s.cfg, roomUC, s.logger)
	roomWSHandlers := roomWS.NewRoomWSHandlers(s.cfg, roomUC, s.logger)
	roomCodeHandlers := roomCodeHttp.NewRoomCodeHandlers(s.cfg, roomCodeUC, s.logger)

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

	authGroup := e.Group("/auth")
	roomGroup := e.Group("/room")
	roomCodeGroup := e.Group("/room_code")

	authHttp.MapAuthRoutes(authGroup, authHandlers, mw)
	roomHttp.MapRoomRoutes(roomGroup, roomHandlers, mw)
	roomWS.MapRoomRoutes(roomGroup, roomWSHandlers, mw)
	roomCodeHttp.MapRoomRoutes(roomCodeGroup, roomCodeHandlers, mw)

	return nil
}
