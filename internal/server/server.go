package server

import (
	"context"
	"fmt"
	"github.com/futod4m4/m/config"
	"github.com/futod4m4/m/internal/websocket"
	"github.com/futod4m4/m/pkg/logger"
	"github.com/go-redis/redis/v8"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	certFile       = "ssl/Server.crt"
	keyFile        = "ssl/Server.pem"
	maxHeaderBytes = 1 << 20
	ctxTimeout     = 5
)

// Server struct
type Server struct {
	echo        *echo.Echo
	cfg         *config.Config
	logger      logger.Logger
	redisClient *redis.Client
	db          *sqlx.DB
}

func NewServer(cfg *config.Config, logger logger.Logger, db *sqlx.DB, redisClient *redis.Client) *Server {
	return &Server{
		echo:        echo.New(),
		cfg:         cfg,
		logger:      logger,
		db:          db,
		redisClient: redisClient,
	}
}

func (s *Server) Run() error {
	if s.cfg.Server.SSL {
		if err := s.MapHandlers(s.echo); err != nil {
			return err
		}

		s.echo.Server.ReadTimeout = time.Second * s.cfg.Server.ReadTimeout
		s.echo.Server.WriteTimeout = time.Second * s.cfg.Server.WriteTimeout

		go func() {
			s.logger.Infof("Server is listening on PORT: %s", s.cfg.Server.Port)
			s.echo.Server.ReadTimeout = time.Second * s.cfg.Server.ReadTimeout
			s.echo.Server.WriteTimeout = time.Second * s.cfg.Server.WriteTimeout
			s.echo.Server.MaxHeaderBytes = maxHeaderBytes
			if err := s.echo.StartTLS(s.cfg.Server.Port, certFile, keyFile); err != nil {
				s.logger.Fatalf("Error starting TLS Server: %v", err)
			}
		}()

		go func() {
			fmt.Println(s.cfg.WebSocketConfig.SocketPort)
			s.logger.Infof("WebSocket server is listening in PORT: %s", ":8081")
			upgrader := websocket.NewUpgrader(s.cfg)
			http.HandleFunc("/ws", websocket.HandleWebSocket(upgrader))
			log.Fatal(http.ListenAndServe(":8081", nil))
		}()

		quit := make(chan os.Signal, 1)
		signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

		<-quit

		ctx, shutdown := context.WithTimeout(context.Background(), ctxTimeout*time.Second)
		defer shutdown()

		s.logger.Info("Server Exited Properly")
		return s.echo.Server.Shutdown(ctx)
	}

	server := &http.Server{
		Addr:           s.cfg.Server.Port,
		ReadTimeout:    time.Second * s.cfg.Server.ReadTimeout,
		WriteTimeout:   time.Second * s.cfg.Server.WriteTimeout,
		MaxHeaderBytes: maxHeaderBytes,
	}

	go func() {
		s.logger.Infof("Server is listening on PORT: %s", s.cfg.Server.Port)
		if err := s.echo.StartServer(server); err != nil {
			s.logger.Fatalf("Error starting Server: %v", err)
		}
	}()

	go func() {
		fmt.Println(s.cfg.WebSocketConfig.SocketPort)
		s.logger.Infof("WebSocket server is listening in PORT: %s", ":8081")
		upgrader := websocket.NewUpgrader(s.cfg)
		http.HandleFunc("/ws", websocket.HandleWebSocket(upgrader))
		log.Fatal(http.ListenAndServe(":8081", nil))
	}()

	if err := s.MapHandlers(s.echo); err != nil {
		return err
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	<-quit

	ctx, shutdown := context.WithTimeout(context.Background(), ctxTimeout*time.Second)
	defer shutdown()

	s.logger.Info("Server Exited Properly")
	return s.echo.Server.Shutdown(ctx)
}
