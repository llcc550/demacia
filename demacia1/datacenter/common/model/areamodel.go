package model

import (
	"fmt"
	"strings"

	"demacia/common/cachemodel"

	"gitlab.u-jy.cn/xiaoyang/go-zero/core/stores/redis"
	"gitlab.u-jy.cn/xiaoyang/go-zero/core/stores/sqlx"
	"xorm.io/builder"
)

var (
	areaFieldNames     = strings.Join(cachemodel.RawFieldNames(&Area{}, true), ",")
	cacheAreaPidPrefix = "cache:common:area:pid:"
	cacheAreaIdPrefix  = "cache:common:area:id:"
)

type (
	AreaModel struct {
		*cachemodel.CachedModel
	}
	Area struct {
		Id   int64  `db:"id"`
		Pid  int64  `db:"pid"`
		Name string `db:"name"`
	}
	Areas []*Area
)

func NewAreaModel(conn sqlx.SqlConn, cache *redis.Redis) *AreaModel {
	return &AreaModel{
		CachedModel: cachemodel.NewCachedModel(conn, `"common"."area"`, cache),
	}
}

func (m *AreaModel) FindListByPid(pid int64) (Areas, error) {
	key := fmt.Sprintf("%s%d", cacheAreaPidPrefix, pid)
	var resp Areas
	err := m.QueryRow(&resp, key, func(conn sqlx.SqlConn, v interface{}) error {
		query, args, _ := builder.Postgres().From(m.Table).Select(areaFieldNames).Where(builder.Eq{"pid": pid}).ToSQL()
		return conn.QueryRows(v, query, args...)
	})
	if err != nil {
		return nil, cachemodel.ErrNotFound
	}
	return resp, nil
}

func (m *AreaModel) FindOneById(id int64) (*Area, error) {
	key := fmt.Sprintf("%s%d", cacheAreaIdPrefix, id)
	var resp Area
	err := m.QueryRow(&resp, key, func(conn sqlx.SqlConn, v interface{}) error {
		query, args, _ := builder.Postgres().From(m.Table).Select(areaFieldNames).Where(builder.Eq{"id": id}).ToSQL()
		return conn.QueryRow(v, query, args...)
	})
	if err != nil {
		return nil, cachemodel.ErrNotFound
	}
	return &resp, nil
}
