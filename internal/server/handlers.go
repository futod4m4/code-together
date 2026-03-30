package server

import (
	authHttp "github.com/futod4m4/m/internal/auth/delivery/http"
	authRepository "github.com/futod4m4/m/internal/auth/repository"
	authUseCase "github.com/futod4m4/m/internal/auth/usecase"
	chatHttp "github.com/futod4m4/m/internal/chat/delivery/http"
	chatRepository "github.com/futod4m4/m/internal/chat/repository"
	chatUseCase "github.com/futod4m4/m/internal/chat/usecase"
	fileHttp "github.com/futod4m4/m/internal/files/delivery/http"
	fileRepository "github.com/futod4m4/m/internal/files/repository"
	fileUseCase "github.com/futod4m4/m/internal/files/usecase"
	memberHttp "github.com/futod4m4/m/internal/members/delivery/http"
	memberRepository "github.com/futod4m4/m/internal/members/repository"
	memberUseCase "github.com/futod4m4/m/internal/members/usecase"
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
	"github.com/futod4m4/m/internal/sessions"
	"github.com/futod4m4/m/pkg/csrf"
	"github.com/futod4m4/m/pkg/metric"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"net/http"
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
	chRepo := chatRepository.NewChatRepository(s.db)
	fRepo := fileRepository.NewFileRepository(s.db)
	mRepo := memberRepository.NewMemberRepository(s.db)
	sRepo := sessionRepository.NewSessionRepository(s.redisClient, s.cfg)
	authRedisRepo := authRepository.NewAuthRedisRepository(s.redisClient)
	roomRedisRepo := roomRepository.NewRoomRedisRepo(s.redisClient)
	roomCodeRedisRepo := roomCodeRepository.NewRoomCodeRedisRepo(s.redisClient)

	// Init useCases
	authUC := authUseCase.NewAuthUseCase(s.cfg, aRepo, authRedisRepo, s.logger)
	roomUC := roomUseCase.NewRoomUseCase(s.cfg, rRepo, roomRedisRepo, s.logger)
	roomCodeUC := roomCodeUseCase.NewRoomCodeUseCase(s.cfg, rcRepo, roomCodeRedisRepo, s.logger)
	chatUC := chatUseCase.NewChatUseCase(chRepo, s.logger)
	fileUC := fileUseCase.NewFileUseCase(fRepo, s.logger)
	memberUC := memberUseCase.NewMemberUseCase(mRepo, s.logger)
	sessUC := usecase.NewSessionUseCase(sRepo, s.cfg)

	// Init handlers
	authHandlers := authHttp.NewAuthHandlers(s.cfg, authUC, sessUC, s.logger)
	roomHandlers := roomHttp.NewRoomHandlers(s.cfg, roomUC, memberUC, s.logger)
	roomWSHandlers := roomWS.NewRoomWSHandlers(s.cfg, roomUC, s.logger)
	roomCodeHandlers := roomCodeHttp.NewRoomCodeHandlers(s.cfg, roomCodeUC, s.logger)
	chatHandlers := chatHttp.NewChatHandlers(chatUC, s.logger)
	fileHandlers := fileHttp.NewFileHandlers(fileUC, s.logger)
	memberHandlers := memberHttp.NewMemberHandlers(memberUC, s.logger)

	mw := apiMiddlewares.NewMiddlewareManager(sessUC, authUC, s.cfg, []string{"*"}, s.logger)

	e.Use(mw.RequestLoggerMiddleware)
	e.Use(mw.MetricsMiddleware(metrics))

	if s.cfg.Server.SSL {
		e.Pre(middleware.HTTPSRedirect())
	}

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderXRequestID, csrf.CSRFHeader},
		AllowCredentials: true,
		AllowMethods:     []string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete, http.MethodOptions},
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
	chatGroup := e.Group("/chat")
	fileGroup := e.Group("/files")
	memberGroup := e.Group("/members")

	authHttp.MapAuthRoutes(authGroup, authHandlers, mw)
	roomHttp.MapRoomRoutes(roomGroup, roomHandlers, mw)
	roomWS.MapRoomRoutes(roomGroup, roomWSHandlers, mw)
	roomCodeHttp.MapRoomRoutes(roomCodeGroup, roomCodeHandlers, mw)
	chatHttp.MapChatRoutes(chatGroup, chatHandlers, mw)
	fileHttp.MapFileRoutes(fileGroup, fileHandlers)
	memberHttp.MapMemberRoutes(memberGroup, memberHandlers, mw)

	// Coding sessions
	sessHandlers := sessions.NewHandlers(s.db)
	sessGroup := e.Group("/sessions")
	sessGroup.POST("/start", sessHandlers.StartSession(), mw.AuthSessionMiddleware)
	sessGroup.POST("/stop/:session_id", sessHandlers.StopSession(), mw.AuthSessionMiddleware)
	sessGroup.POST("/snapshot", sessHandlers.AddSnapshot())
	sessGroup.POST("/viewers", sessHandlers.UpdateViewerCount())
	sessGroup.GET("/room/:room_id", sessHandlers.GetRoomSessions())
	sessGroup.GET("/:session_id/snapshots", sessHandlers.GetSessionSnapshots())

	// GitHub import
	ghHandlers := sessions.NewGitHubHandlers(s.db)
	e.POST("/github/import", ghHandlers.ImportRepo(), mw.AuthSessionMiddleware)

	// Ban system
	banHandlers := sessions.NewBanHandlers(s.db)
	banGroup := e.Group("/bans")
	banGroup.POST("/user", banHandlers.BanUser(), mw.AuthSessionMiddleware)
	banGroup.POST("/ip", banHandlers.BanIP(), mw.AuthSessionMiddleware)
	banGroup.GET("/check/:room_id", banHandlers.CheckBan())
	banGroup.GET("/list/:room_id", banHandlers.GetBannedList(), mw.AuthSessionMiddleware)
	banGroup.DELETE("/:ban_id", banHandlers.Unban(), mw.AuthSessionMiddleware)

	return nil
}
