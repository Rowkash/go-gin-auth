package sessions

import (
	"context"
	"fmt"
	"time"

	"github.com/Rowkash/go-gin-auth/internal/common"
	"github.com/redis/go-redis/v9"
)

type Service struct {
	rdb *redis.Client
}

type UserSession struct {
	UserId       string        `json:"userId" redis:"userId"`
	RefreshToken string        `json:"refreshToken" redis:"refreshToken"`
	ExpiresIn    time.Duration `json:"expiresIn" redis:"-"`
	CreatedAt    int64         `json:"createdAt" redis:"createdAt"`
}

const RedisNamespace = "api:session"

func NewService(rdb *redis.Client) *Service {
	return &Service{
		rdb: rdb,
	}
}

func (s *Service) Create(ctx context.Context, key string, value UserSession) error {
	pipe := s.rdb.TxPipeline()
	sessionKey := s.getSessionKey(key)
	userSetKey := s.getUserSessionsSetKey(value.UserId)

	pipe.HSet(ctx, sessionKey, value)
	pipe.Expire(ctx, sessionKey, value.ExpiresIn)
	pipe.SAdd(ctx, userSetKey, sessionKey)
	_, err := pipe.Exec(ctx)
	return err
}

func (s *Service) findOneByKey(ctx context.Context, key string) (UserSession, error) {
	sessionKey := s.getSessionKey(key)
	var session UserSession

	err := s.rdb.HGetAll(ctx, sessionKey).Scan(&session)
	if err != nil {
		return UserSession{}, err
	}

	if session.UserId == "" {
		return UserSession{}, common.NewNotFoundError("Session not found")
	}

	return session, nil
}

func (s *Service) Delete(ctx context.Context, key string) error {
	session, err := s.findOneByKey(ctx, key)
	if err != nil {
		return err
	}

	sessionKey := s.getSessionKey(key)
	userSetKey := s.getUserSessionsSetKey(session.UserId)

	pipe := s.rdb.TxPipeline()
	pipe.Del(ctx, sessionKey)
	pipe.SRem(ctx, userSetKey, sessionKey)

	_, err = pipe.Exec(ctx)
	return err
}

func (s *Service) getSessionKey(tokenOrID string) string {
	return fmt.Sprintf("%s:sessions:%s", RedisNamespace, tokenOrID)
}

func (s *Service) getUserSessionsSetKey(userID string) string {
	return fmt.Sprintf("%s:user:%s", RedisNamespace, userID)
}
