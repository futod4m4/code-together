package repository

import (
	"context"
	"database/sql"
	"github.com/futod4m4/m/internal/auth"
	"github.com/futod4m4/m/internal/models"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type AuthRepo struct {
	db *sqlx.DB
}

func NewAuthRepository(db *sqlx.DB) auth.Repository {
	return &AuthRepo{db: db}
}

// Register Create new user
func (r AuthRepo) Register(ctx context.Context, user *models.User) (*models.User, error) {
	u := &models.User{}
	if err := r.db.QueryRowxContext(ctx, createUserQuery, &user.Email, &user.Password,
		&user.Nickname, &user.FirstName, &user.LastName,
	).StructScan(u); err != nil {
		return nil, errors.Wrap(err, "authRepo.Register.StructScan")
	}

	return u, nil
}

// Update Updating existing user
func (r AuthRepo) Update(ctx context.Context, user *models.User) (*models.User, error) {
	u := &models.User{}
	if err := r.db.QueryRowxContext(ctx, updateUserQuery, &user.FirstName, &user.LastName, &user.Nickname, &user.Email, &user.UserID).StructScan(u); err != nil {
		return nil, errors.Wrap(err, "authRepo.Update.StructScan")
	}

	return u, nil
}

// Delete Deleting existing user
func (r AuthRepo) Delete(ctx context.Context, userID uuid.UUID) error {
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

func (r AuthRepo) GetByID(ctx context.Context, userID uuid.UUID) (*models.User, error) {
	u := &models.User{}
	if err := r.db.QueryRowxContext(ctx, getByIDQuery, userID).StructScan(u); err != nil {
		return nil, errors.Wrap(err, "authRepo.GetByID.StructScan")
	}

	return u, nil
}

func (r AuthRepo) FindUserByEmail(ctx context.Context, user *models.User) (*models.User, error) {
	foundUser := &models.User{}
	if err := r.db.QueryRowxContext(ctx, findUserByEmailQuery, user.Email).StructScan(foundUser); err != nil {
		return nil, errors.Wrap(err, "authRepo.FindUserByEmail.StructScan")
	}

	return foundUser, nil
}
