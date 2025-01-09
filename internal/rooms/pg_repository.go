package rooms

import (
	"context"
	"github.com/futod4m4/m/internal/models"
	"github.com/google/uuid"
)

type Repository interface {
	CreateRoom(ctx context.Context, room *models.Room) (*models.Room, error)
	UpdateRoom(ctx context.Context, room *models.Room) (*models.Room, error)
	DeleteRoom(ctx context.Context, roomID uuid.UUID) error
	GetRoomByID(ctx context.Context, roomID uuid.UUID) (*models.Room, error)
	GetRoomByJoinCode(ctx context.Context, joinCode string) (*models.Room, error)
}
