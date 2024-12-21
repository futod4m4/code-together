package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/futod4m4/m/config"
	"github.com/futod4m4/m/internal/models"
	"github.com/futod4m4/m/internal/session"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"time"
)

const basePrefix = "api-session:"

type sessionRepo struct {
	redisClient *redis.Client
	basePrefix  string
	cfg         *config.Config
}

func NewSessionRepository(redisClient *redis.Client, cfg *config.Config) session.SessRepository {
	return &sessionRepo{
		redisClient: redisClient,
		basePrefix:  basePrefix,
		cfg:         cfg,
	}
}

func (s *sessionRepo) CreateSession(ctx context.Context, sess *models.Session, expire int) (string, error) {
	sess.SessionID = uuid.New().String()
	sessionKey := s.createKey(sess.SessionID)

	sessBytes, err := json.Marshal(&sess)
	if err != nil {
		return "", errors.WithMessage(err, "sessionRepo.CreateSession.json.Marshal")
	}

	if err = s.redisClient.Set(ctx, sessionKey, sessBytes, time.Second*time.Duration(expire)).Err(); err != nil {
		return "", errors.Wrap(err, "sessionRepo.CreateSession.redisClient.Set")
	}

	return sessionKey, nil
}

func (s *sessionRepo) GetSessionByID(ctx context.Context, sessionID string) (*models.Session, error) {
	sessBytes, err := s.redisClient.Get(ctx, sessionID).Bytes()
	if err != nil {
		return nil, errors.Wrap(err, "sessionRepo.GetSessionByID.redisClient.Get")
	}

	sess := &models.Session{}
	if err = json.Unmarshal(sessBytes, &sess); err != nil {
		return nil, errors.Wrap(err, "sessionRepo.GetSessionByID.Unmarshall")
	}

	return sess, nil
}

func (s *sessionRepo) DeleteSessionByID(ctx context.Context, sessionID string) error {
	if err := s.redisClient.Del(ctx, sessionID).Err(); err != nil {
		return errors.Wrap(err, "sessionRepo.DeleteSessionByID")
	}

	return nil
}

func (s *sessionRepo) createKey(sessionID string) string {
	return fmt.Sprintf("%s: %s", s.basePrefix, sessionID)
}
