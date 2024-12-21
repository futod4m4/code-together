package session

import (
	"context"
	"github.com/futod4m4/m/internal/models"
)

type SessRepository interface {
	CreateSession(ctx context.Context, sess *models.Session, expire int) (string, error)
	GetSessionByID(ctx context.Context, sessionID string) (*models.Session, error)
	DeleteSessionByID(ctx context.Context, sessionID string) error
}
