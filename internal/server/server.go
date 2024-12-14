package server

import (
	"context"
	"github.com/futod4m4/m/config"
	"github.com/futod4m4/m/pkg/logger"
	"github.com/labstack/echo"
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
	echo *echo.Echo
	cfg  *config.Config

	logger logger.Logger
}

func NewServer(cfg *config.Config, logger logger.Logger) *Server {
	return &Server{
		echo:   echo.New(),
		cfg:    cfg,
		logger: logger,
	}
}

func (s *Server) Run() error {
	if s.cfg.Server.SSL {
		// TODO: MapHandlers
		//if err := s.MapHandlers(s.echo); err != nil {
		//	return err
		//}

		s.echo.Server.ReadTimeout = time.Second * s.cfg.Server.ReadTimeout
		s.echo.Server.WriteTimeout = time.Second * s.cfg.Server.WriteTimeout

		go func() {
			s.logger.Infof("Server is listening on PORT: %s", s.cfg.Server.Port)
			s.echo.Server.ReadTimeout = time.Second * s.cfg.Server.ReadTimeout
			s.echo.Server.WriteTimeout = time.Second * s.cfg.Server.WriteTimeout
			s.echo.Server.MaxHeaderBytes = maxHeaderBytes
			if err := s.echo.StartTLS(s.cfg.Server.Port, certFile, keyFile); err != nil {
				s.logger.Fatalf("Error starting TLS Server: ", err)
			}
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
			s.logger.Fatalf("Error starting Server: ", err)
		}
	}()

	//if err := s.MapHandlers(s.echo); err != nil {
	//	return err
	//}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	<-quit

	ctx, shutdown := context.WithTimeout(context.Background(), ctxTimeout*time.Second)
	defer shutdown()

	s.logger.Info("Server Exited Properly")
	return s.echo.Server.Shutdown(ctx)
}
