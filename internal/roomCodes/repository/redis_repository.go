package repository

import (
	"context"
	"encoding/json"
	"github.com/futod4m4/m/internal/models"
	"github.com/futod4m4/m/internal/roomCodes"
	"github.com/go-redis/redis/v8"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"time"
)

type roomRedisRepo struct {
	redisClient *redis.Client
}

func NewRoomCodeRedisRepo(redisClient *redis.Client) roomCodes.RedisRepository {
	return &roomRedisRepo{redisClient: redisClient}
}

// Get new by id
func (n *roomRedisRepo) GetRoomCodeByIDCtx(ctx context.Context, key string) (*models.RoomCode, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "roomCodeRedisRepo.GetRoomByIDCtx")
	defer span.Finish()

	roomCodeBytes, err := n.redisClient.Get(ctx, key).Bytes()
	if err != nil {
		return nil, errors.Wrap(err, "roomCodeRedisRepo.GetRoomByIDCtx.redisClient.Get")
	}
	roomCodeBase := &models.RoomCode{}
	if err = json.Unmarshal(roomCodeBytes, roomCodeBase); err != nil {
		return nil, errors.Wrap(err, "roomCodeRedisRepo.GetRoomByIDCtx.json.Unmarshal")
	}

	return roomCodeBase, nil
}

// Cache room item
func (n *roomRedisRepo) SetRoomCodeCtx(ctx context.Context, key string, seconds int, news *models.RoomCode) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "roomCodeRedisRepo.SetRoomCtx")
	defer span.Finish()

	roomCodeBytes, err := json.Marshal(news)
	if err != nil {
		return errors.Wrap(err, "roomCodeRedisRepo.SetRoomCtx.json.Marshal")
	}
	if err = n.redisClient.Set(ctx, key, roomCodeBytes, time.Second*time.Duration(seconds)).Err(); err != nil {
		return errors.Wrap(err, "roomCodeRedisRepo.SetRoomCtx.redisClient.Set")
	}
	return nil
}

// Delete new item from cache
func (n *roomRedisRepo) DeleteRoomCodeCtx(ctx context.Context, key string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "roomCodeRedisRepo.DeleteRoomCtx")
	defer span.Finish()

	if err := n.redisClient.Del(ctx, key).Err(); err != nil {
		return errors.Wrap(err, "roomCodeRedisRepo.DeleteRoomCtx.redisClient.Del")
	}
	return nil
}
