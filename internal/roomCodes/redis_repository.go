package roomCodes

import (
	"context"
	"github.com/futod4m4/m/internal/models"
)

type RedisRepository interface {
	GetRoomCodeByIDCtx(ctx context.Context, key string) (*models.RoomCode, error)
	SetRoomCodeCtx(ctx context.Context, key string, seconds int, room *models.RoomCode) error
	DeleteRoomCodeCtx(ctx context.Context, key string) error
}
