package pkg

import (
	"github.com/go-redis/redis"
	"github.com/kutty-kumar/charminder/pkg"
	"github.com/sirupsen/logrus"
	"time"
)


type RedisCache struct {
	*redis.Client
	logger *logrus.Logger
	entityCreator pkg.EntityCreator
}

func (r *RedisCache) Put(base pkg.Base) error {
	cmd := r.Client.Set(base.GetExternalId(), base, 0)
	return cmd.Err()
}

func (r *RedisCache) Get(externalId string) (pkg.Base, error) {
	cmd := r.Client.Get(externalId)
	if cmd.Err() != nil {
		return nil, cmd.Err()
	}
	entity := r.entityCreator()
	err := cmd.Scan(&entity)
	if err != nil {
		return nil, err
	}
	return entity, nil
}

func (r *RedisCache) MultiGet(externalIds []string) ([]pkg.Base, error) {
	var result []pkg.Base
	for _, externalId := range externalIds {
		base, err := r.Get(externalId)
		if err != nil {
			return nil, err
		}
		result = append(result, base)
	}
	return result, nil
}

func (r *RedisCache) Delete(externalId string) error {
	statusCmd := r.Client.Del(externalId)
	if statusCmd.Err() != nil {
		return statusCmd.Err()
	}
	return nil
}

func (r *RedisCache) MultiDelete(externalIds []string) error {
	statusCmd := r.Client.Del(externalIds...)
	if statusCmd.Err() != nil {
		return statusCmd.Err()
	}
	return nil
}

func (r *RedisCache) PutWithTtl(base pkg.Base, duration time.Duration) error {
	statusCmd := r.Client.Set(base.GetExternalId(), base, duration)
	if statusCmd.Err() != nil {
		return statusCmd.Err()
	}
	return nil
}

func (r *RedisCache) DeleteAll() error {
	cmd := r.Client.FlushDB()
	if cmd.Err() != nil {
		return cmd.Err()
	}
	return nil
}

func (r *RedisCache) Health() error {
	pong, err := r.Client.Ping().Result()
	if err != nil {
		return err
	}
	r.logger.Infof("Health check ping response <%v>", pong)
	return nil
}

func NewRedisCache(addr string, password string, db uint, logger *logrus.Logger, entityCreator pkg.EntityCreator) Cache {
	client := redis.NewClient(
		&redis.Options{
			Addr: addr,
			Password:  password,
			DB:  int(db),
		})
	return &RedisCache{
		client,
		logger,
		entityCreator,
	}
}

