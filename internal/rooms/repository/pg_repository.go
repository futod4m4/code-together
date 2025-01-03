package repository

import (
	"context"
	"database/sql"
	"github.com/futod4m4/m/internal/models"
	"github.com/futod4m4/m/internal/rooms"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
)

type roomRepo struct {
	db *sqlx.DB
}

func NewAuthRepository(db *sqlx.DB) rooms.Repository {
	return &roomRepo{db: db}
}

func (r *roomRepo) CreateRoom(ctx context.Context, room *models.Room) (*models.Room, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "roomRepo.CreateRoom")
	defer span.Finish()

	var ro models.Room
	if err := r.db.QueryRowxContext(
		ctx,
		createRoom,
		&room.Name,
		&room.JoinCode,
		&room.Language,
		&room.OwnerID,
	).StructScan(&ro); err != nil {
		return nil, errors.Wrap(err, "roomRepo.Create.QueryRowxContext")
	}

	return &ro, nil
}

func (r *roomRepo) UpdateRoom(ctx context.Context, room *models.Room) (*models.Room, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "roomRepo.UpdateRoom")
	defer span.Finish()

	var ro models.Room
	if err := r.db.QueryRowxContext(
		ctx,
		updateRoom,
		&room.Name,
		&room.Language,
		&room.OwnerID,
	).StructScan(&ro); err != nil {
		return nil, errors.Wrap(err, "roomRepo.Update.QueryRowxContext")
	}

	return &ro, nil
}

func (r *roomRepo) DeleteRoom(ctx context.Context, roomID uuid.UUID) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "roomRepo.DeleteRoom")
	defer span.Finish()

	result, err := r.db.ExecContext(ctx, deleteRoomByID, roomID)
	if err != nil {
		return errors.Wrap(err, "roomRepo.DeleteRoom.ExecContext")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return errors.Wrap(err, "roomRepo.DeleteRoom.RowsAffected")
	}
	if rowsAffected == 0 {
		return errors.Wrap(sql.ErrNoRows, "roomRepo.DeleteRoom.rowsAffected")
	}

	return nil
}

func (r *roomRepo) GetRoomByID(ctx context.Context, roomID uuid.UUID) (*models.Room, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "roomRepo.GetRoomByID")
	defer span.Finish()

	ro := &models.Room{}
	if err := r.db.GetContext(ctx, ro, getRoomByID, roomID); err != nil {
		return nil, errors.Wrap(err, "roomRepo.GetRoomByID.GetContext")
	}

	return ro, nil
}

func (r *roomRepo) GetRoomByJoinCode(ctx context.Context, joinCode string) (*models.Room, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "roomRepo.GetRoomByID")
	defer span.Finish()

	ro := &models.Room{}
	if err := r.db.GetContext(ctx, ro, getRoomByID, joinCode); err != nil {
		return nil, errors.Wrap(err, "roomRepo.GetRoomByID.GetContext")
	}

	return ro, nil
}

func (r *roomRepo) GetRoomCode(roomID int) (string, error) {
	//TODO implement me
	return "", nil
}

func (r *roomRepo) UpdateRoomCode(roomID int, code string, userID string) error {
	//TODO implement me
	return nil
}
