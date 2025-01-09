package usecase

import (
	"context"
	"fmt"
	"github.com/futod4m4/m/config"
	"github.com/futod4m4/m/internal/models"
	"github.com/futod4m4/m/internal/rooms"
	"github.com/futod4m4/m/pkg/httpErrors"
	"github.com/futod4m4/m/pkg/logger"
	"github.com/futod4m4/m/pkg/utils"
	"github.com/google/uuid"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"net/http"
)

const (
	baseprefix    = "api-rooms:"
	cacheDuration = 7200
)

type roomUC struct {
	cfg       *config.Config
	roomRepo  rooms.Repository
	redisRepo rooms.RedisRepository
	logger    logger.Logger
}

func NewRoomUseCase(cfg *config.Config, roomRepo rooms.Repository, redisRepo rooms.RedisRepository, logger logger.Logger) rooms.RoomUseCase {
	return &roomUC{
		cfg:       cfg,
		roomRepo:  roomRepo,
		redisRepo: redisRepo,
		logger:    logger,
	}
}

func (u *roomUC) CreateRoom(ctx context.Context, room *models.Room) (*models.Room, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "roomUC.CreateRoom")
	defer span.Finish()

	user, err := utils.GetUserFromCtx(ctx)
	if err != nil {
		return nil, httpErrors.NewUnauthorizedError(errors.WithMessage(err, "roomsUC.CreateRoom.GetUserFromCtx"))
	}

	room.OwnerID = user.UserID

	err = room.GenJoinCode()
	if err != nil {
		return nil, err
	}

	if room.Name == "" {
		room.Name = user.Nickname + "'s Room"
	}

	if err = utils.ValidateStruct(ctx, room); err != nil {
		return nil, httpErrors.NewBadRequestError(errors.WithMessage(err, "roomsUC.CreateRoom.ValidateStruct"))
	}

	r, err := u.roomRepo.CreateRoom(ctx, room)
	if err != nil {
		return nil, err
	}

	return r, err
}

func (u *roomUC) UpdateRoom(ctx context.Context, room *models.Room) (*models.Room, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "roomUC.UpdateRoom")
	defer span.Finish()

	roomByID, err := u.roomRepo.GetRoomByID(ctx, room.ID)
	if err != nil {
		return nil, err
	}

	if err = utils.ValidateIsOwner(ctx, roomByID.OwnerID.String(), u.logger); err != nil {
		return nil, httpErrors.NewRestError(http.StatusForbidden, "Forbidden", errors.Wrap(err, "roomUC.UpdateRoom.ValidateIsOwner"))
	}

	updatedRoom, err := u.roomRepo.UpdateRoom(ctx, room)
	if err != nil {
		return nil, err
	}

	if err = u.redisRepo.DeleteRoomCtx(ctx, u.getKeyWithPrefix(room.ID.String())); err != nil {
		u.logger.Errorf("roomUC.Update.DeleteRoomCtx: %v", err)
	}

	return updatedRoom, nil
}

func (u *roomUC) DeleteRoom(ctx context.Context, roomID uuid.UUID) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "roomUC.DeleteRoom")
	defer span.Finish()

	roomByID, err := u.roomRepo.GetRoomByID(ctx, roomID)
	if err != nil {
		return err
	}

	if err = utils.ValidateIsOwner(ctx, roomByID.OwnerID.String(), u.logger); err != nil {
		return httpErrors.NewRestError(http.StatusForbidden, "Forbidden", errors.Wrap(err, "roomUC.DeleteRoom.ValidateIsOwner"))
	}

	if err = u.roomRepo.DeleteRoom(ctx, roomID); err != nil {
		return err
	}

	if err = u.redisRepo.DeleteRoomCtx(ctx, u.getKeyWithPrefix(roomID.String())); err != nil {
		return err
	}

	return nil
}

func (u *roomUC) GetRoomByID(ctx context.Context, roomID uuid.UUID) (*models.Room, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "roomUC.GetRoomByID")
	defer span.Finish()

	roomBase, err := u.redisRepo.GetRoomByIDCtx(ctx, u.getKeyWithPrefix(roomID.String()))
	if err != nil {
		u.logger.Errorf("roomUC.GetRoomByID.GetRoomByIDCtx: %v", err)
	}
	if roomBase != nil {
		return roomBase, nil
	}

	r, err := u.roomRepo.GetRoomByID(ctx, roomID)
	if err != nil {
		return nil, err
	}

	if err = u.redisRepo.SetRoomCtx(ctx, u.getKeyWithPrefix(roomID.String()), cacheDuration, r); err != nil {
		u.logger.Errorf("roomUC.GetRoomByID.SetRoomCtx: %s", err)
	}

	return r, err
}

func (u *roomUC) GetRoomByJoinCode(ctx context.Context, joinCode string) (*models.Room, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "roomUC.GetRoomByJoinCode")
	defer span.Finish()

	r, err := u.roomRepo.GetRoomByJoinCode(ctx, joinCode)
	if err != nil {
		return nil, err
	}

	return r, nil
}

func (u *roomUC) getKeyWithPrefix(roomID string) string {
	return fmt.Sprintf("%s: %s", baseprefix, roomID)
}
