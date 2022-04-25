package model

import (
	"database/sql"
	"demacia/common/cachemodel"
	"fmt"
	"gitlab.u-jy.cn/xiaoyang/go-zero/core/stores/redis"
	"strings"
	"xorm.io/builder"

	"gitlab.u-jy.cn/xiaoyang/go-zero/core/stores/sqlx"
)

var (
	groupFieldNames    = strings.Join(cachemodel.RawFieldNames(&Group{}, true), ",")
	cacheGroupIdPrefix = "cache:device:group:id:"
)

type (
	GroupModel struct {
		*cachemodel.CachedModel
	}

	Group struct {
		Id    int64  `db:"id"`
		OrgId int64  `db:"org_id"`
		Name  string `db:"name"`
	}
	Groups []*Group
)

func NewGroupModel(conn sqlx.SqlConn, cache *redis.Redis) *GroupModel {
	return &GroupModel{
		CachedModel: cachemodel.NewCachedModel(conn, `"device"."group"`, cache),
	}
}

func (m *GroupModel) Insert(data *Group) (int64, error) {
	eq := builder.Eq{
		"name":   data.Name,
		"org_id": data.OrgId,
	}
	var groupId int64
	query, args, _ := builder.Postgres().Insert(eq).Into(m.Table).ToSQL()
	err := m.Conn.QueryRow(&groupId, query+" returning id", args...)
	if err != nil {
		return 0, err
	}
	_ = m.DelCache(m.formatPrimary(groupId))
	return groupId, nil
}

func (m *GroupModel) FindOne(id int64) (*Group, error) {
	deviceGroupIdKey := fmt.Sprintf("%s%v", cacheGroupIdPrefix, id)
	var resp Group
	err := m.QueryRow(&resp, deviceGroupIdKey, func(conn sqlx.SqlConn, v interface{}) error {
		query := fmt.Sprintf("select %s from %s where id = $1 limit 1", groupFieldNames, m.Table)
		return conn.QueryRow(v, query, id)
	})
	if err != nil {
		return nil, cachemodel.ErrNotFound
	}
	return &resp, nil
}

func (m *GroupModel) Update(data *Group) error {
	deviceGroupIdKey := fmt.Sprintf("%s%v", cacheGroupIdPrefix, data.Id)
	_, err := m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query, args, _ := builder.Postgres().Update(builder.Eq{"name": data.Name}).From(m.Table).Where(builder.Eq{"id": data.Id, "deleted": false}).ToSQL()
		return conn.Exec(query, args...)
	}, deviceGroupIdKey)
	return err
}

func (m *GroupModel) Delete(id int64) error {
	deviceGroupIdKey := fmt.Sprintf("%s%v", cacheGroupIdPrefix, id)
	_, err := m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query, args, _ := builder.Postgres().Update(builder.Eq{"deleted": true}).From(m.Table).Where(builder.Eq{"id": id}).ToSQL()
		return conn.Exec(query, args...)
	}, deviceGroupIdKey)
	return err
}

func (m *GroupModel) formatPrimary(primary interface{}) string {
	return fmt.Sprintf("%s%v", cacheGroupIdPrefix, primary)
}

func (m *GroupModel) queryPrimary(conn sqlx.SqlConn, v, primary interface{}) error {
	query := fmt.Sprintf("select %s from %s where id = $1 limit 1", groupFieldNames, m.Table)
	return conn.QueryRow(v, query, primary)
}

func (m *GroupModel) FindListByIds(ids []int64) (Groups, error) {
	var res Groups
	query, args, _ := builder.Postgres().Select(groupFieldNames).From(m.Table).Where(builder.In("id", ids)).And(builder.Eq{"deleted": false}).ToSQL()
	err := m.Conn.QueryRows(&res, query, args...)
	if err != nil {
		return nil, err
	}
	return res, err
}

func (m *GroupModel) FindListByOrgId(orgId int64) (Groups, error) {
	var res Groups
	query, args, _ := builder.Postgres().Select(groupFieldNames).From(m.Table).Where(builder.Eq{"org_id": orgId}).And(builder.Eq{"deleted": false}).OrderBy("id DESC").ToSQL()
	err := m.Conn.QueryRows(&res, query, args...)
	if err != nil {
		return nil, err
	}
	return res, err
}

func (m *GroupModel) FindListByGroupName(orgId int64, name string) (Groups, error) {
	var res Groups
	eq := builder.And()
	if name != "" {
		eq = eq.And(builder.Like{"name", name})
	}
	query, args, _ := builder.Postgres().Select(groupFieldNames).From(m.Table).Where(eq).And(builder.Eq{"org_id": orgId, "deleted": false}).ToSQL()
	err := m.Conn.QueryRows(&res, query, args...)
	if err != nil {
		return nil, err
	}
	return res, err
}
func (m *GroupModel) FindOneByGroupName(orgId int64, name string) (*Group, error) {
	var res Group
	query, args, _ := builder.Postgres().Select(groupFieldNames).From(m.Table).Where(builder.Eq{"org_id": orgId, "name": name, "deleted": false}).Limit(1).ToSQL()
	err := m.Conn.QueryRow(&res, query, args...)
	if err != nil {
		return nil, err
	}
	return &res, err
}
