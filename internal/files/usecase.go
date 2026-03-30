package files

import (
	"context"
	"github.com/futod4m4/m/internal/models"
	"github.com/google/uuid"
)

type UseCase interface {
	CreateFile(ctx context.Context, file *models.RoomFile) (*models.RoomFile, error)
	UpdateFile(ctx context.Context, file *models.RoomFile) (*models.RoomFile, error)
	DeleteFile(ctx context.Context, fileID uuid.UUID) error
	GetFileByID(ctx context.Context, fileID uuid.UUID) (*models.RoomFile, error)
	GetFilesByRoomID(ctx context.Context, roomID uuid.UUID) ([]*models.RoomFile, error)
}
