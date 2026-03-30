package usecase

import (
	"context"
	"github.com/futod4m4/m/internal/chat"
	"github.com/futod4m4/m/internal/models"
	"github.com/futod4m4/m/pkg/logger"
	"github.com/google/uuid"
	"github.com/opentracing/opentracing-go"
)

type chatUC struct {
	chatRepo chat.Repository
	logger   logger.Logger
}

func NewChatUseCase(chatRepo chat.Repository, logger logger.Logger) chat.UseCase {
	return &chatUC{
		chatRepo: chatRepo,
		logger:   logger,
	}
}

func (u *chatUC) CreateMessage(ctx context.Context, msg *models.RoomMessage) (*models.RoomMessage, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "chatUC.CreateMessage")
	defer span.Finish()

	return u.chatRepo.CreateMessage(ctx, msg)
}

func (u *chatUC) GetMessagesByRoomID(ctx context.Context, roomID uuid.UUID, limit, offset int) ([]*models.RoomMessage, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "chatUC.GetMessagesByRoomID")
	defer span.Finish()

	return u.chatRepo.GetMessagesByRoomID(ctx, roomID, limit, offset)
}
