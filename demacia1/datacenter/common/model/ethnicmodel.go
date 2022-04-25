package model

import (
	"strings"

	"demacia/common/cachemodel"

	"gitlab.u-jy.cn/xiaoyang/go-zero/core/stores/redis"
	"gitlab.u-jy.cn/xiaoyang/go-zero/core/stores/sqlx"
	"xorm.io/builder"
)

var (
	ethnicFieldNames = strings.Join(cachemodel.RawFieldNames(&Ethnic{}, true), ",")
	cacheEthnic      = "cache:common:ethnic:all"
)

type (
	EthnicModel struct {
		*cachemodel.CachedModel
	}
	Ethnic struct {
		Id   int64  `db:"id"`
		Name string `db:"name"`
	}
	Ethnics []*Ethnic
)

func NewEthnicModel(conn sqlx.SqlConn, cache *redis.Redis) *EthnicModel {
	return &EthnicModel{
		CachedModel: cachemodel.NewCachedModel(conn, `"common"."ethnic"`, cache),
	}
}

func (m *EthnicModel) FindList() (Ethnics, error) {
	key := cacheEthnic
	var resp Ethnics
	err := m.QueryRow(&resp, key, func(conn sqlx.SqlConn, v interface{}) error {
		query, args, _ := builder.Postgres().From(m.Table).Select(ethnicFieldNames).OrderBy("id").ToSQL()
		return conn.QueryRows(v, query, args...)
	})
	if err != nil {
		return nil, cachemodel.ErrNotFound
	}
	return resp, nil
}
