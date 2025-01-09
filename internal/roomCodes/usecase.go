package roomCodes

import (
	"context"
	"github.com/futod4m4/m/internal/models"
	"github.com/google/uuid"
)

type RoomCodeUseCase interface {
	CreateRoomCode(ctx context.Context, roomCode *models.RoomCode) (*models.RoomCode, error)
	UpdateRoomCode(ctx context.Context, roomCode *models.RoomCode) (*models.RoomCode, error)
	DeleteRoomCode(ctx context.Context, roomCodeID uuid.UUID) error
	GetRoomCodeByID(ctx context.Context, roomCodeID uuid.UUID) (*models.RoomCode, error)
	GetRoomCodeByRoomID(ctx context.Context, roomID uuid.UUID) (uuid.UUID, error)
}
