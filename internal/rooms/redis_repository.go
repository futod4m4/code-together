package rooms

import (
	"context"
	"github.com/futod4m4/m/internal/models"
)

type RedisRepository interface {
	GetRoomByIDCtx(ctx context.Context, key string) (*models.Room, error)
	SetRoomCtx(ctx context.Context, key string, seconds int, room *models.Room) error
	DeleteRoomCtx(ctx context.Context, key string) error
}
