package repository

import (
	"context"
	"database/sql"
	"github.com/futod4m4/m/internal/auth"
	"github.com/futod4m4/m/internal/models"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
)

type authRepo struct {
	db *sqlx.DB
}

func NewAuthRepository(db *sqlx.DB) auth.Repository {
	return &authRepo{db: db}
}

// Register Create new user
func (r *authRepo) Register(ctx context.Context, user *models.User) (*models.User, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "authRepo.Register")
	defer span.Finish()

	u := &models.User{}
	if err := r.db.QueryRowxContext(ctx, createUserQuery, &user.Email,
		&user.Password, &user.Nickname, &user.FirstName, &user.LastName,
	).StructScan(u); err != nil {
		return nil, errors.Wrap(err, "authRepo.Register.StructScan")
	}

	return u, nil
}

// Update Updating existing user
func (r *authRepo) Update(ctx context.Context, user *models.User) (*models.User, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "authRepo.Update")
	defer span.Finish()

	u := &models.User{}
	if err := r.db.QueryRowxContext(ctx, updateUserQuery, &user.FirstName, &user.LastName, &user.Nickname, &user.Email, &user.UserID).StructScan(u); err != nil {
		return nil, errors.Wrap(err, "authRepo.Update.StructScan")
	}

	return u, nil
}

// Delete Deleting existing user
func (r *authRepo) Delete(ctx context.Context, userID uuid.UUID) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "authRepo.Delete")
	defer span.Finish()

	result, err := r.db.ExecContext(ctx, deleteUserQuery, userID)
	if err != nil {
		return errors.WithMessage(err, "authRepo Delete ExecContext")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return errors.Wrap(err, "authRepo.Delete.RowsAffected")
	}

	if rowsAffected == 0 {
		return errors.Wrap(sql.ErrNoRows, "authRepo.Delete.rowsAffected")
	}

	return nil
}

func (r *authRepo) GetByID(ctx context.Context, userID uuid.UUID) (*models.User, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "authRepo.GetByID")
	defer span.Finish()

	u := &models.User{}
	if err := r.db.QueryRowxContext(ctx, getByIDQuery, userID).StructScan(u); err != nil {
		return nil, errors.Wrap(err, "authRepo.GetByID.StructScan")
	}

	return u, nil
}

func (r *authRepo) FindUserByEmail(ctx context.Context, user *models.User) (*models.User, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "authRepo.FindUserByEmail")
	defer span.Finish()

	foundUser := &models.User{}
	if err := r.db.QueryRowxContext(ctx, findUserByEmailQuery, user.Email).StructScan(foundUser); err != nil {
		return nil, errors.Wrap(err, "authRepo.FindByEmail.QueryRowxContext")
	}
	return foundUser, nil
}
