package model

import (
	"database/sql"
	"demacia/common/cachemodel"
	"fmt"
	"gitlab.u-jy.cn/xiaoyang/go-zero/core/stores/redis"
	"strings"
	"xorm.io/builder"

	"gitlab.u-jy.cn/xiaoyang/go-zero/core/stores/sqlc"
	"gitlab.u-jy.cn/xiaoyang/go-zero/core/stores/sqlx"
)

var (
	categoryFieldNames            = strings.Join(cachemodel.RawFieldNames(&Category{}, true), ",")
	cacheEventCategoryIdPrefix    = "cache:event:category:id:"
	cacheEventCategoryOrgIdPrefix = "cache:event:category:org-id:"
	cacheEventCategoryNamePrefix  = "cache:event:category:name:"
)

type (
	CategoryModel struct {
		*cachemodel.CachedModel
	}

	Category struct {
		Id    int64  `db:"id"`
		OrgId int64  `db:"org_id"`
		Name  string `db:"name"`
	}
	Categories []*Category
)

func NewCategoryModel(conn sqlx.SqlConn, cache *redis.Redis) *CategoryModel {
	return &CategoryModel{
		CachedModel: cachemodel.NewCachedModel(conn, `"event"."category"`, cache),
	}
}

func (m *CategoryModel) Insert(data *Category) (int64, error) {
	var res int64
	query, args, _ := builder.Postgres().Insert(builder.Eq{"org_id": data.OrgId, "name": data.Name}).Into(m.Table).ToSQL()
	err := m.Conn.QueryRow(&res, query+" returning id ", args...)
	return res, err
}

func (m *CategoryModel) Update(data *Category) error {
	_, err := m.Exec(func(conn sqlx.SqlConn) (sql.Result, error) {
		query, args, _ := builder.Postgres().Update(builder.Eq{"name": data.Name}).From(m.Table).Where(builder.Eq{"id": data.Id}).ToSQL()
		return conn.Exec(query, args...)
	}, m.keys(data)...)
	return err
}

func (m *CategoryModel) FindOneByName(orgId int64, name string) (Category, error) {
	var res Category
	query, args, _ := builder.Postgres().Select(categoryFieldNames).From(m.Table).Where(builder.Eq{"org_id": orgId, "name": name, "deleted": false}).Limit(1).ToSQL()
	err := m.Conn.QueryRow(&res, query, args...)
	return res, err
}

func (m *CategoryModel) FindOne(id int64) (*Category, error) {
	eventCategoryIdKey := fmt.Sprintf("%s%v", cacheEventCategoryIdPrefix, id)
	var resp Category
	err := m.QueryRow(&resp, eventCategoryIdKey, func(conn sqlx.SqlConn, v interface{}) error {
		query, args, _ := builder.Postgres().Select(categoryFieldNames).From(m.Table).Where(builder.Eq{"id": id, "deleted": false}).ToSQL()
		return conn.QueryRow(v, query, args...)
	})
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, cachemodel.ErrNotFound
	default:
		return nil, err
	}
}

func (m *CategoryModel) FindList(orgId int64, name string) (*Categories, error) {
	var resp Categories
	eq := builder.And()
	if name != "" {
		eq = eq.And(builder.Like{"name", name})
	}
	query, args, _ := builder.Postgres().Select(categoryFieldNames).From(m.Table).Where(builder.Eq{"org_id": orgId, "deleted": false}).And(eq).ToSQL()
	err := m.Conn.QueryRows(&resp, query, args...)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (m *CategoryModel) Delete(id int64) error {
	eventCategoryIdKey := fmt.Sprintf("%s%v", cacheEventCategoryIdPrefix, id)
	_, err := m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query, args, _ := builder.Postgres().Update(builder.Eq{"deleted": true}).From(m.Table).Where(builder.Eq{"id": id, "deleted": false}).ToSQL()
		return conn.Exec(query, args...)
	}, eventCategoryIdKey)
	return err
}

func (m *CategoryModel) formatPrimary(primary interface{}) string {
	return fmt.Sprintf("%s%v", cacheEventCategoryIdPrefix, primary)
}

func (m *CategoryModel) queryPrimary(conn sqlx.SqlConn, v, primary interface{}) error {
	query := fmt.Sprintf("select %s from %s where id = $1 limit 1", categoryFieldNames, m.Table)
	return conn.QueryRow(v, query, primary)
}

func (m *CategoryModel) keys(data *Category) []string {
	res := make([]string, 0, 3)
	if data.Id != 0 {
		res = append(res, fmt.Sprintf("%s%d", cacheEventCategoryIdPrefix, data.Id))
	}
	if data.OrgId != 0 {
		res = append(res, fmt.Sprintf("%s%d", cacheEventCategoryOrgIdPrefix, data.OrgId))
	}
	if data.Name != "" {
		res = append(res, fmt.Sprintf("%s%s", cacheEventCategoryNamePrefix, data.Name))
	}
	return res
}
