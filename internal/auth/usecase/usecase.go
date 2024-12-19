package usecase

import (
	"context"
	"fmt"
	"github.com/futod4m4/m/config"
	"github.com/futod4m4/m/internal/auth"
	"github.com/futod4m4/m/internal/models"
	"github.com/futod4m4/m/pkg/httpErrors"
	"github.com/futod4m4/m/pkg/logger"
	"github.com/futod4m4/m/pkg/utils"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"net/http"
)

const (
	basePrefix    = "api-auth:"
	cacheDuration = 3600
)

type AuthUC struct {
	cfg       *config.Config
	authRepo  auth.Repository
	redisRepo auth.RedisRepository
	logger    logger.Logger
}

func NewAuthUseCase(cfg *config.Config, authRepo auth.Repository, redisRepo auth.RedisRepository, logger logger.Logger) *AuthUC {
	return &AuthUC{
		cfg:       cfg,
		authRepo:  authRepo,
		redisRepo: redisRepo,
		logger:    logger,
	}
}

func (u *AuthUC) Register(ctx context.Context, user *models.User) (*models.UserWithToken, error) {
	existsUser, err := u.authRepo.FindByEmail(ctx, user)
	if existsUser != nil || err == nil {
		return nil, httpErrors.NewRestError(http.StatusBadRequest, httpErrors.ErrEmailAlreadyExists, nil)
	}

	if err := user.PrepareCreate(); err != nil {
		return nil, httpErrors.NewBadRequestError(errors.Wrap(err, "authUC.Register.PrepareCreate"))
	}

	createdUser, err := u.authRepo.Register(ctx, user)
	if err != nil {
		return nil, err
	}

	createdUser.SanitizePassword()

	token, err := utils.GenerateJWTToken(createdUser, u.cfg)
	if err != nil {
		return nil, httpErrors.NewInternalServerError(errors.Wrap(err, "authUC.Register.GenerateJWTToken"))
	}

	return &models.UserWithToken{
		User:  createdUser,
		Token: token,
	}, nil
}

func (u *AuthUC) Login(ctx context.Context, user *models.User) (*models.UserWithToken, error) {
	foundUser, err := u.authRepo.FindByEmail(ctx, user)
	if err != nil {
		return nil, err
	}

	if err = foundUser.ComparePasswords(user.Password); err != nil {
		return nil, httpErrors.NewUnauthorizedError(errors.Wrap(err, "authUC.GetUsers.ComparePasswords"))
	}

	foundUser.SanitizePassword()

	token, err := utils.GenerateJWTToken(user, u.cfg)
	if err != nil {
		return nil, httpErrors.NewInternalServerError(errors.Wrap(err, "authUC.GetUsers.GenerateJWTToken"))
	}

	return &models.UserWithToken{
		User:  foundUser,
		Token: token,
	}, nil
}

func (u *AuthUC) Update(ctx context.Context, user *models.User) (*models.User, error) {
	if err := user.PrepareUpdate(); err != nil {
		return nil, httpErrors.NewBadRequestError(errors.Wrap(err, "authUC.Register.PrepareUpdate"))
	}

	updatedUser, err := u.authRepo.Update(ctx, user)
	if err != nil {
		return nil, err
	}

	if err := u.redisRepo.DeleteUserCtx(ctx, u.GenerateUserKey(user.UserID.String())); err != nil {
		u.logger.Errorf("AuthUC.Update.DeleteUserCtx: %s", err)
	}

	updatedUser.SanitizePassword()

	return updatedUser, nil
}

func (u *AuthUC) Delete(ctx context.Context, userID uuid.UUID) error {
	if err := u.authRepo.Delete(ctx, userID); err != nil {
		return err
	}

	if err := u.redisRepo.DeleteUserCtx(ctx, u.GenerateUserKey(userID.String())); err != nil {
		u.logger.Errorf("AuthUC.Delete.DeleteUserCtx: %s", err)
	}

	return nil
}

func (u *AuthUC) GetByID(ctx context.Context, userID uuid.UUID) (*models.User, error) {

	cachedUser, err := u.redisRepo.GetByIDCtx(ctx, u.GenerateUserKey(userID.String()))
	if err != nil {
		u.logger.Errorf("authUC.GetByID.GetByIDCtx: %v", err)
	}
	if cachedUser != nil {
		return cachedUser, nil
	}

	user, err := u.authRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if err := u.redisRepo.SetUserCtx(ctx, u.GenerateUserKey(user.UserID.String()), cacheDuration, user); err != nil {
		u.logger.Errorf("authUC.GetByID.SetUserCtx: %v", err)
	}

	user.SanitizePassword()

	return user, nil
}

func (u *AuthUC) GenerateUserKey(userID string) string {
	return fmt.Sprintf("%s: %s", basePrefix, userID)
}
