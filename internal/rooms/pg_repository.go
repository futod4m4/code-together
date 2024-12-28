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

	// Code management
	GetRoomCode(roomID int) (string, error)
	UpdateRoomCode(roomID int, code string, userID string) error
}
