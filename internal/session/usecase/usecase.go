package usecase

import (
	"context"
	"github.com/futod4m4/m/config"
	"github.com/futod4m4/m/internal/models"
	"github.com/futod4m4/m/internal/session"
	"github.com/opentracing/opentracing-go"
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
	span, ctx := opentracing.StartSpanFromContext(ctx, "sessionUC.CreateSession")
	defer span.Finish()

	return u.sessionRepo.CreateSession(ctx, session, expire)
}

func (u *sessionUC) GetSessionByID(ctx context.Context, sessionID string) (*models.Session, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "sessionUC.GetSessionByID")
	defer span.Finish()

	return u.GetSessionByID(ctx, sessionID)
}

func (u *sessionUC) DeleteSessionByID(ctx context.Context, sessionID string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "sessionUC.DeleteSessionByID")
	defer span.Finish()

	return u.DeleteSessionByID(ctx, sessionID)
}
