package repository

import (
	"context"
	"encoding/json"
	"github.com/futod4m4/m/internal/models"
	"github.com/futod4m4/m/internal/rooms"
	"github.com/go-redis/redis/v8"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"time"
)

type roomRedisRepo struct {
	redisClient *redis.Client
}

func NewRoomRedisRepo(redisClient *redis.Client) rooms.RedisRepository {
	return &roomRedisRepo{redisClient: redisClient}
}

// Get new by id
func (n *roomRedisRepo) GetRoomByIDCtx(ctx context.Context, key string) (*models.Room, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "roomRedisRepo.GetRoomByIDCtx")
	defer span.Finish()

	roomBytes, err := n.redisClient.Get(ctx, key).Bytes()
	if err != nil {
		return nil, errors.Wrap(err, "roomRedisRepo.GetRoomByIDCtx.redisClient.Get")
	}
	roomBase := &models.Room{}
	if err = json.Unmarshal(roomBytes, roomBase); err != nil {
		return nil, errors.Wrap(err, "roomRedisRepo.GetRoomByIDCtx.json.Unmarshal")
	}

	return roomBase, nil
}

// Cache room item
func (n *roomRedisRepo) SetRoomCtx(ctx context.Context, key string, seconds int, news *models.Room) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "roomRedisRepo.SetRoomCtx")
	defer span.Finish()

	roomBytes, err := json.Marshal(news)
	if err != nil {
		return errors.Wrap(err, "roomRedisRepo.SetRoomCtx.json.Marshal")
	}
	if err = n.redisClient.Set(ctx, key, roomBytes, time.Second*time.Duration(seconds)).Err(); err != nil {
		return errors.Wrap(err, "roomRedisRepo.SetRoomCtx.redisClient.Set")
	}
	return nil
}

// Delete new item from cache
func (n *roomRedisRepo) DeleteRoomCtx(ctx context.Context, key string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "roomRedisRepo.DeleteRoomCtx")
	defer span.Finish()

	if err := n.redisClient.Del(ctx, key).Err(); err != nil {
		return errors.Wrap(err, "roomRedisRepo.DeleteRoomCtx.redisClient.Del")
	}
	return nil
}
