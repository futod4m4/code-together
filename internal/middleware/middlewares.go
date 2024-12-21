package middleware

import (
	"github.com/futod4m4/m/config"
	"github.com/futod4m4/m/internal/auth"
	"github.com/futod4m4/m/internal/session"
	"github.com/futod4m4/m/pkg/logger"
)

type MiddlewareManager struct {
	sessUC  session.UCSession
	authUC  auth.UseCase
	cfg     *config.Config
	origins []string
	logger  logger.Logger
}

// NewMiddlewareManager Middleware manager constructor
func NewMiddlewareManager(sessUC session.UCSession, authUC auth.UseCase, cfg *config.Config, origins []string, logger logger.Logger) *MiddlewareManager {
	return &MiddlewareManager{sessUC: sessUC, authUC: authUC, cfg: cfg, origins: origins, logger: logger}
}
