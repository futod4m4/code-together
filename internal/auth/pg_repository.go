package auth

import (
	"context"
	"github.com/futod4m4/m/internal/models"
	"github.com/google/uuid"
)

type Repository interface {
	Register(ctx context.Context, user *models.User) (*models.User, error)
	Update(ctx context.Context, user *models.User) (*models.User, error)
	Delete(ctx context.Context, userID uuid.UUID) error
	GetByID(ctx context.Context, userID uuid.UUID) (*models.User, error)
	FindUserByEmail(ctx context.Context, user *models.User) (*models.User, error)
}
