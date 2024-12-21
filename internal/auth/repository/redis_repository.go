package repository

import (
	"context"
	"encoding/json"
	"github.com/futod4m4/m/internal/auth"
	"github.com/futod4m4/m/internal/models"
	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
	"time"
)

type authRedisRepo struct {
	redisClient *redis.Client
}

func NewAuthRedisRepository(redisClient *redis.Client) auth.RedisRepository {
	return &authRedisRepo{redisClient: redisClient}
}

func (a *authRedisRepo) GetByIDCtx(ctx context.Context, key string) (*models.User, error) {
	userBytes, err := a.redisClient.Get(ctx, key).Bytes()
	if err != nil {
		return nil, errors.Wrap(err, "authRedisRepo.GetByIDCtx.Get")
	}

	user := &models.User{}
	if err = json.Unmarshal(userBytes, &user); err != nil {
		return nil, errors.Wrap(err, "authRedisRepo.GetByIDCtx.Unmarshall")
	}

	return user, nil
}

func (a *authRedisRepo) SetUserCtx(ctx context.Context, key string, seconds int, user *models.User) error {

	userBytes, err := json.Marshal(user)
	if err != nil {
		return errors.Wrap(err, "authRedisRepo.SetUserCtx.json.Marshal")
	}

	if err = a.redisClient.Set(ctx, key, userBytes, time.Second*time.Duration(seconds)).Err(); err != nil {
		return errors.Wrap(err, "authRedisRepo.SetUserCtx.Set")
	}

	return nil
}

func (a *authRedisRepo) DeleteUserCtx(ctx context.Context, key string) error {
	if err := a.redisClient.Del(ctx, key).Err(); err != nil {
		return errors.Wrap(err, "authRedisRepo.DeleteUserCtx")
	}

	return nil
}
