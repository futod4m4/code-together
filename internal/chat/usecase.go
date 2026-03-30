package chat

import (
	"context"
	"github.com/futod4m4/m/internal/models"
	"github.com/google/uuid"
)

type UseCase interface {
	CreateMessage(ctx context.Context, msg *models.RoomMessage) (*models.RoomMessage, error)
	GetMessagesByRoomID(ctx context.Context, roomID uuid.UUID, limit, offset int) ([]*models.RoomMessage, error)
}
