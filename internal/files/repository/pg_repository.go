package repository

import (
	"context"
	"database/sql"
	"github.com/futod4m4/m/internal/files"
	"github.com/futod4m4/m/internal/models"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
)

type fileRepo struct {
	db *sqlx.DB
}

func NewFileRepository(db *sqlx.DB) files.Repository {
	return &fileRepo{db: db}
}

func (r *fileRepo) CreateFile(ctx context.Context, file *models.RoomFile) (*models.RoomFile, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "fileRepo.CreateFile")
	defer span.Finish()

	var f models.RoomFile
	if err := r.db.QueryRowxContext(ctx, createFile,
		&file.RoomID, &file.Filename, &file.Language, &file.Content, &file.IsEntryPoint,
	).StructScan(&f); err != nil {
		return nil, errors.Wrap(err, "fileRepo.CreateFile.QueryRowxContext")
	}
	return &f, nil
}

func (r *fileRepo) UpdateFile(ctx context.Context, file *models.RoomFile) (*models.RoomFile, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "fileRepo.UpdateFile")
	defer span.Finish()

	var f models.RoomFile
	if err := r.db.QueryRowxContext(ctx, updateFile,
		&file.Content, &file.Language, &file.Filename, &file.ID,
	).StructScan(&f); err != nil {
		return nil, errors.Wrap(err, "fileRepo.UpdateFile.QueryRowxContext")
	}
	return &f, nil
}

func (r *fileRepo) DeleteFile(ctx context.Context, fileID uuid.UUID) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "fileRepo.DeleteFile")
	defer span.Finish()

	result, err := r.db.ExecContext(ctx, deleteFile, fileID)
	if err != nil {
		return errors.Wrap(err, "fileRepo.DeleteFile.ExecContext")
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return errors.Wrap(err, "fileRepo.DeleteFile.RowsAffected")
	}
	if rowsAffected == 0 {
		return errors.Wrap(sql.ErrNoRows, "fileRepo.DeleteFile.rowsAffected")
	}
	return nil
}

func (r *fileRepo) GetFileByID(ctx context.Context, fileID uuid.UUID) (*models.RoomFile, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "fileRepo.GetFileByID")
	defer span.Finish()

	f := &models.RoomFile{}
	if err := r.db.GetContext(ctx, f, getFileByID, fileID); err != nil {
		return nil, errors.Wrap(err, "fileRepo.GetFileByID.GetContext")
	}
	return f, nil
}

func (r *fileRepo) GetFilesByRoomID(ctx context.Context, roomID uuid.UUID) ([]*models.RoomFile, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "fileRepo.GetFilesByRoomID")
	defer span.Finish()

	var filesList []*models.RoomFile
	if err := r.db.SelectContext(ctx, &filesList, getFilesByRoomID, roomID); err != nil {
		return nil, errors.Wrap(err, "fileRepo.GetFilesByRoomID.SelectContext")
	}
	return filesList, nil
}

func (r *fileRepo) CountFilesByRoomID(ctx context.Context, roomID uuid.UUID) (int, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "fileRepo.CountFilesByRoomID")
	defer span.Finish()

	var count int
	if err := r.db.GetContext(ctx, &count, countFilesByRoomID, roomID); err != nil {
		return 0, errors.Wrap(err, "fileRepo.CountFilesByRoomID.GetContext")
	}
	return count, nil
}
