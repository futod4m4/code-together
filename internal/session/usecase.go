package session

import (
	"context"
	"github.com/futod4m4/m/internal/models"
)

type UCSession interface {
	CreateSession(ctx context.Context, session *models.Session, expire int) (string, error)
	GetSessionByID(ctx context.Context, sessionID string) (*models.Session, error)
	DeleteSessionByID(ctx context.Context, sessionID string) error
	RefreshSession(ctx context.Context, sessionID string, expire int) error
}
