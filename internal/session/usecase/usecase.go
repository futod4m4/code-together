package usecase

import (
	"context"
	"github.com/futod4m4/m/config"
	"github.com/futod4m4/m/internal/models"
	"github.com/futod4m4/m/internal/session"
)

// Session use case
type sessionUC struct {
	sessionRepo session.SessRepository
	cfg         *config.Config
}

func NewSessionUseCase(sessionRepo session.SessRepository, cfg *config.Config) session.UCSession {
	return &sessionUC{
		sessionRepo: sessionRepo,
		cfg:         cfg,
	}
}

func (u *sessionUC) CreateSession(ctx context.Context, session *models.Session, expire int) (string, error) {
	return u.sessionRepo.CreateSession(ctx, session, expire)
}

func (u *sessionUC) GetSessionByID(ctx context.Context, sessionID string) (*models.Session, error) {
	return u.GetSessionByID(ctx, sessionID)
}

func (u *sessionUC) DeleteSessionByID(ctx context.Context, sessionID string) error {
	return u.DeleteSessionByID(ctx, sessionID)
}
