package usecase

import (
	"context"
	"fmt"
	"github.com/futod4m4/m/internal/files"
	"github.com/futod4m4/m/internal/models"
	"github.com/futod4m4/m/pkg/httpErrors"
	"github.com/futod4m4/m/pkg/logger"
	"github.com/google/uuid"
	"github.com/opentracing/opentracing-go"
	"net/http"
)

const maxFilesPerRoom = 10

type fileUC struct {
	fileRepo files.Repository
	logger   logger.Logger
}

func NewFileUseCase(fileRepo files.Repository, logger logger.Logger) files.UseCase {
	return &fileUC{
		fileRepo: fileRepo,
		logger:   logger,
	}
}

func (u *fileUC) CreateFile(ctx context.Context, file *models.RoomFile) (*models.RoomFile, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "fileUC.CreateFile")
	defer span.Finish()

	count, err := u.fileRepo.CountFilesByRoomID(ctx, file.RoomID)
	if err != nil {
		return nil, err
	}
	if count >= maxFilesPerRoom {
		return nil, httpErrors.NewRestError(http.StatusBadRequest,
			fmt.Sprintf("Maximum %d files per room", maxFilesPerRoom), nil)
	}

	return u.fileRepo.CreateFile(ctx, file)
}

func (u *fileUC) UpdateFile(ctx context.Context, file *models.RoomFile) (*models.RoomFile, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "fileUC.UpdateFile")
	defer span.Finish()

	return u.fileRepo.UpdateFile(ctx, file)
}

func (u *fileUC) DeleteFile(ctx context.Context, fileID uuid.UUID) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "fileUC.DeleteFile")
	defer span.Finish()

	return u.fileRepo.DeleteFile(ctx, fileID)
}

func (u *fileUC) GetFileByID(ctx context.Context, fileID uuid.UUID) (*models.RoomFile, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "fileUC.GetFileByID")
	defer span.Finish()

	return u.fileRepo.GetFileByID(ctx, fileID)
}

func (u *fileUC) GetFilesByRoomID(ctx context.Context, roomID uuid.UUID) ([]*models.RoomFile, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "fileUC.GetFilesByRoomID")
	defer span.Finish()

	return u.fileRepo.GetFilesByRoomID(ctx, roomID)
}
