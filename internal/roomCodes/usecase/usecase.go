package usecase

import (
	"context"
	"fmt"
	"github.com/futod4m4/m/config"
	"github.com/futod4m4/m/internal/models"
	"github.com/futod4m4/m/internal/roomCodes"
	"github.com/futod4m4/m/pkg/logger"
	"github.com/google/uuid"
	"github.com/opentracing/opentracing-go"
)

const (
	baseprefix    = "api-rooms:"
	cacheDuration = 7200
)

type roomCodeUC struct {
	cfg          *config.Config
	roomCodeRepo roomCodes.Repository
	redisRepo    roomCodes.RedisRepository
	logger       logger.Logger
}

func NewRoomCodeUseCase(cfg *config.Config, roomCodeRepo roomCodes.Repository, redisRepo roomCodes.RedisRepository, logger logger.Logger) roomCodes.RoomCodeUseCase {
	return &roomCodeUC{
		cfg:          cfg,
		roomCodeRepo: roomCodeRepo,
		redisRepo:    redisRepo,
		logger:       logger,
	}
}

func (u *roomCodeUC) CreateRoomCode(ctx context.Context, roomCode *models.RoomCode) (*models.RoomCode, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "roomCodeUC.CreateRoomCode")
	defer span.Finish()

	r, err := u.roomCodeRepo.CreateRoomCode(ctx, roomCode)
	if err != nil {
		return nil, err
	}

	return r, err
}

func (u *roomCodeUC) UpdateRoomCode(ctx context.Context, roomCode *models.RoomCode) (*models.RoomCode, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "roomCodeUC.UpdateRoomCode")
	defer span.Finish()

	updatedRoomCode, err := u.roomCodeRepo.UpdateRoomCode(ctx, roomCode)
	if err != nil {
		u.logger.Errorf("Error updating room code: %v", err)
		return nil, err
	}

	if err = u.redisRepo.DeleteRoomCodeCtx(ctx, u.getKeyWithPrefix(roomCode.ID.String())); err != nil {
		u.logger.Errorf("Error deleting room code from Redis: %v", err)
	}

	return updatedRoomCode, nil
}

func (u *roomCodeUC) DeleteRoomCode(ctx context.Context, roomCodeID uuid.UUID) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "roomCodeUC.DeleteRoomCode")
	defer span.Finish()

	if err := u.roomCodeRepo.DeleteRoomCode(ctx, roomCodeID); err != nil {
		return err
	}

	if err := u.redisRepo.DeleteRoomCodeCtx(ctx, u.getKeyWithPrefix(roomCodeID.String())); err != nil {
		return err
	}

	return nil
}

func (u *roomCodeUC) GetRoomCodeByID(ctx context.Context, roomCodeID uuid.UUID) (*models.RoomCode, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "roomCodeUC.GetRoomCodeByID")
	defer span.Finish()

	roomCodeBase, err := u.redisRepo.GetRoomCodeByIDCtx(ctx, u.getKeyWithPrefix(roomCodeID.String()))
	if err != nil {
		u.logger.Errorf("roomCodeUC.GetRoomCodeByID.GetRoomByIDCtx: %v", err)
	}
	if roomCodeBase != nil {
		return roomCodeBase, nil
	}

	r, err := u.roomCodeRepo.GetRoomCodeByID(ctx, roomCodeID)
	if err != nil {
		return nil, err
	}

	if err = u.redisRepo.SetRoomCodeCtx(ctx, u.getKeyWithPrefix(roomCodeID.String()), cacheDuration, r); err != nil {
		u.logger.Errorf("roomCodeUC.GetRoomCodeByID.SetRoomCodeCtx: %s", err)
	}

	return r, err
}

func (u *roomCodeUC) GetRoomCodeByRoomID(ctx context.Context, roomID uuid.UUID) (uuid.UUID, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "roomCodeUC.GetRoomCodeByID")
	defer span.Finish()

	rc, err := u.roomCodeRepo.GetRoomCodeByRoomID(ctx, roomID)
	if err != nil {
		return uuid.UUID{}, err
	}

	return rc, err
}

func (u *roomCodeUC) getKeyWithPrefix(roomCodeID string) string {
	return fmt.Sprintf("%s: %s", baseprefix, roomCodeID)
}
