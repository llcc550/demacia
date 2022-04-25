package cachemodel

import (
	"time"

	"gitlab.u-jy.cn/xiaoyang/go-zero/core/stores/cache"
	"gitlab.u-jy.cn/xiaoyang/go-zero/core/stores/redis"
	"gitlab.u-jy.cn/xiaoyang/go-zero/core/stores/sqlc"
	"gitlab.u-jy.cn/xiaoyang/go-zero/core/stores/sqlx"
)

const expireCacheKey = 60 * 60 * 2

type CachedModel struct {
	Table string
	Conn  sqlx.SqlConn
	cache *redis.Redis
	sqlc.CachedConn
}

func NewCachedModel(conn sqlx.SqlConn, table string, rds *redis.Redis) *CachedModel {
	return &CachedModel{
		Table:      table,
		Conn:       conn,
		cache:      rds,
		CachedConn: sqlc.NewNodeConn(conn, rds, cache.WithExpiry(time.Second*expireCacheKey)),
	}
}

var ErrNotFound = sqlx.ErrNotFound
