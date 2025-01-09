package repository

import (
	"context"
	"database/sql"
	"github.com/futod4m4/m/internal/models"
	"github.com/futod4m4/m/internal/roomCodes"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
)

type roomCodeRepo struct {
	db *sqlx.DB
}

func NewRoomCodeRepository(db *sqlx.DB) roomCodes.Repository {
	return &roomCodeRepo{db: db}
}

func (r *roomCodeRepo) CreateRoomCode(ctx context.Context, roomCode *models.RoomCode) (*models.RoomCode, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "roomCodeRepo.CreateRoomCode")
	defer span.Finish()

	var rc models.RoomCode
	if err := r.db.QueryRowxContext(
		ctx,
		createRoomCode,
		&roomCode.RoomID,
	).StructScan(&rc); err != nil {
		return nil, errors.Wrap(err, "roomCodeRepo.CreateRoomCode.QueryRowxContext")
	}

	return &rc, nil
}

func (r *roomCodeRepo) UpdateRoomCode(ctx context.Context, roomCode *models.RoomCode) (*models.RoomCode, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "roomRepo.UpdateRoomCode")
	defer span.Finish()

	var rc models.RoomCode
	if err := r.db.QueryRowxContext(
		ctx,
		updateRoomCode,
		&roomCode.RoomID,
		&roomCode.Code,
		&roomCode.ID,
	).StructScan(&rc); err != nil {
		return nil, errors.Wrap(err, "roomCodeRepo.Update.QueryRowxContext")
	}

	return &rc, nil
}

func (r *roomCodeRepo) DeleteRoomCode(ctx context.Context, roomCodeID uuid.UUID) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "roomCodeRepo.DeleteRoomCode")
	defer span.Finish()

	result, err := r.db.ExecContext(ctx, deleteRoomCodeByID, roomCodeID)
	if err != nil {
		return errors.Wrap(err, "roomCodeRepo.DeleteRoomCode.ExecContext")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return errors.Wrap(err, "roomRepo.DeleteRoomCode.RowsAffected")
	}
	if rowsAffected == 0 {
		return errors.Wrap(sql.ErrNoRows, "roomRepo.DeleteRoomCode.rowsAffected")
	}

	return nil
}

func (r *roomCodeRepo) GetRoomCodeByID(ctx context.Context, roomCodeID uuid.UUID) (*models.RoomCode, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "roomCodeRepo.GetRoomCodeByID")
	defer span.Finish()

	rc := &models.RoomCode{}
	if err := r.db.GetContext(ctx, rc, getRoomCodeByID, roomCodeID); err != nil {
		return nil, errors.Wrap(err, "roomCodeRepo.GetRoomCodeByID.GetContext")
	}

	return rc, nil
}

func (r *roomCodeRepo) GetRoomCodeByRoomID(ctx context.Context, roomID uuid.UUID) (uuid.UUID, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "roomCodeRepo.GetRoomCodeByRoomID")
	defer span.Finish()

	var rc uuid.UUID
	if err := r.db.GetContext(ctx, &rc, getRoomCodeByRoomID, roomID); err != nil {
		return uuid.UUID{}, errors.Wrap(err, "roomCodeRepo.GetRoomCodeByRoomID.GetContext")
	}

	return rc, nil
}
