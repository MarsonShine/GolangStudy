package dao

import (
	"context"

	"github.com/go-kratos/kratos/pkg/cache/redis"
	"github.com/go-kratos/kratos/pkg/conf/paladin"
	"github.com/go-kratos/kratos/pkg/log"
)

func NewRedis() (r *redis.Redis, cf func(), err error) {
	var (
		cfg redis.Config
		ct  paladin.Map
	)
	if err = paladin.Get("redis.toml").Unmarshal(&ct); err != nil {
		return
	}
	if err = ct.Get("Client").UnmarshalTOML(&cfg); err != nil {
		return
	}
	r = redis.NewRedis(&cfg)
	cf = func() { r.Close() }
	return
}

func (d *dao) PingRedis(ctx context.Context) (err error) {
	if _, err = d.redis.Do(ctx, "SET", "ping", "pong"); err != nil {
		log.Error("conn.Set(PING) error(%v)", err)
	}
	return
}

func (d *dao) GetDemo(c context.Context, key string) (string, error) {
	conn := d.redis.Conn(c)
	defer conn.Close()
	// 如果没有就去数据库取
	if s, _ := redis.String(conn.Do("GET", key)); s == "" {
		// 取数据库
		if _, err := d.redis.Do(c, "SET", key, "marsonshine"); err != nil {
			log.Error("conn.Set(%s) error(%v)", key, err)
		}
	}
	return redis.String(conn.Do("GET", key))
}
