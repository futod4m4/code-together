package repository

import (
	"context"
	"github.com/futod4m4/m/internal/chat"
	"github.com/futod4m4/m/internal/models"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
)

type chatRepo struct {
	db *sqlx.DB
}

func NewChatRepository(db *sqlx.DB) chat.Repository {
	return &chatRepo{db: db}
}

func (r *chatRepo) CreateMessage(ctx context.Context, msg *models.RoomMessage) (*models.RoomMessage, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "chatRepo.CreateMessage")
	defer span.Finish()

	var m models.RoomMessage
	if err := r.db.QueryRowxContext(
		ctx,
		createMessage,
		&msg.RoomID,
		&msg.UserID,
		&msg.Nickname,
		&msg.Content,
	).StructScan(&m); err != nil {
		return nil, errors.Wrap(err, "chatRepo.CreateMessage.QueryRowxContext")
	}

	return &m, nil
}

func (r *chatRepo) GetMessagesByRoomID(ctx context.Context, roomID uuid.UUID, limit, offset int) ([]*models.RoomMessage, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "chatRepo.GetMessagesByRoomID")
	defer span.Finish()

	var messages []*models.RoomMessage
	if err := r.db.SelectContext(ctx, &messages, getMessagesByRoomID, roomID, limit, offset); err != nil {
		return nil, errors.Wrap(err, "chatRepo.GetMessagesByRoomID.SelectContext")
	}

	return messages, nil
}
