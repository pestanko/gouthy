package store

import (
	"context"
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"github.com/pestanko/gouthy/app/shared"
	log "github.com/sirupsen/logrus"
	"time"
)

const RedisStoresSize = 10

func NewRedisStoresFromConfig(cfg *shared.RedisConfig) Stores {
	return NewRedisStores(cfg.Address, cfg.Password)
}

func NewRedisStores(addr, passwd string) Stores {
	return &redisStores{
		addr:      addr,
		password:  passwd,
		instances: make([]Store, RedisStoresSize),
	}
}

type redisStores struct {
	addr      string
	password  string
	instances []Store
}

func (store *redisStores) GetStore(storeId int) Store {
	if storeId > len(store.instances) {
		return nil
	}

	if store.instances[storeId] == nil {
		store.instances[storeId] = store.initializeInstance(storeId)
	}

	return nil
}

func (store *redisStores) initializeInstance(id int) Store {
	client := redis.NewClient(&redis.Options{
		Addr:     store.addr,
		Password: store.password,
		DB:       id,
	})
	return newRedisStoreDB(client, id)
}

func newRedisStoreDB(client *redis.Client, id int) Store {
	return &redisStoreDB{client: client, id: id}
}

type redisStoreDB struct {
	client *redis.Client
	id     int
}

func (r *redisStoreDB) ListKeys(ctx context.Context, pattern string) ([]string, error) {
	entry := shared.GetLogger(ctx).WithFields(log.Fields{
		"pattern":  pattern,
		"store_id": r.id,
	})
	entry.Debug("Store - list fields")
	return r.client.Keys(ctx, pattern).Result()
}

func (r *redisStoreDB) Get(ctx context.Context, key string) (string, error) {
	entry := shared.GetLogger(ctx).WithFields(log.Fields{
		"key":      key,
		"store_id": r.id,
	})
	str, err := r.client.Get(ctx, key).Result()
	if err != nil {
		entry.WithError(err).Warning("Store - get field failed")
	} else {
		entry.Debug("Store - get field")
	}

	return str, err
}

func (r *redisStoreDB) Set(ctx context.Context, key string, value interface{}, exp time.Duration) error {
	entry := shared.GetLogger(ctx).WithFields(log.Fields{
		"key":      key,
		"store_id": r.id,
		"duration": exp.Seconds(),
	})
	err := r.client.Set(ctx, key, value, exp).Err()
	if err != nil {
		entry.WithError(err).Warning("Store - set field failed")
	} else {
		entry.Debug("Store - set field")
	}

	return err
}

func (r *redisStoreDB) Delete(ctx context.Context, key string) error {
	entry := shared.GetLogger(ctx).WithFields(log.Fields{
		"key":      key,
		"store_id": r.id,
	})
	err := r.client.Del(ctx, key).Err()
	if err != nil {
		entry.WithError(err).Warning("Store - delete field failed")
	} else {
		entry.Debug("Store - delete field")
	}
	return err
}

func (r *redisStoreDB) GetJson(ctx context.Context, key string, value interface{}) error {
	data, err := r.Get(ctx, key)
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(data), value)
}

func (r *redisStoreDB) SetJson(ctx context.Context, key string, value interface{}, exp time.Duration) error {
	bytes, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return r.Set(ctx, key, string(bytes), exp)
}
