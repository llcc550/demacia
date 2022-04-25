package model

import (
	"database/sql"
	"fmt"
	"strings"
	"xorm.io/builder"

	"demacia/common/cachemodel"

	"gitlab.u-jy.cn/xiaoyang/go-zero/core/stores/redis"
	"gitlab.u-jy.cn/xiaoyang/go-zero/core/stores/sqlx"
)

var (
	parentFieldNames                = strings.Join(cachemodel.RawFieldNames(&Parent{}, true), ",")
	cacheParentInfoByParentIdPrefix = "cache:parent:parent:id:"
	cacheParentInfoByMobilePrefix   = "cache:parent:parent:mobile:"
	cacheParentInfoByTrueNamePrefix = "cache:parent:parent:true-name:"
)

type (
	ParentModel struct {
		*cachemodel.CachedModel
	}
	Parent struct {
		Id         int64  `db:"id"`
		TrueName   string `db:"true_name"`
		UserName   string `db:"user_name"`
		Password   string `db:"password"`
		Face       string `db:"face"`
		FaceStatus int8   `db:"face_status"`
		IdNumber   string `db:"id_number"`
		Address    string `db:"address"`
		Mobile     string `db:"mobile"`
		Pinyin     string `db:"pinyin"`
	}
)

func NewParentModel(conn sqlx.SqlConn, cache *redis.Redis) *ParentModel {
	return &ParentModel{
		CachedModel: cachemodel.NewCachedModel(conn, `"parent".parent`, cache),
	}
}

func (m *ParentModel) InsertOne(data *Parent) (int64, error) {
	sql, args, _ := builder.Postgres().Insert(builder.Eq{
		"user_name":   data.UserName,
		"true_name":   data.TrueName,
		"password":    data.Password,
		"face":        data.Face,
		"face_status": data.FaceStatus,
		"id_number":   data.IdNumber,
		"address":     data.Address,
		"mobile":      data.Mobile,
		"pinyin":      data.Pinyin,
	}).Into(m.Table).ToSQL()
	var res int64
	err := m.Conn.QueryRow(&res, sql+" returning id", args...)
	if err != nil {
		return 0, err
	}
	_ = m.DelCache(m.keys(data)...)
	return res, nil
}

func (m *ParentModel) UpdateOne(data *Parent) error {
	_, err := m.Exec(func(conn sqlx.SqlConn) (sql.Result, error) {
		query, args, _ := builder.Postgres().Update(builder.Eq{
			"user_name": data.UserName,
			"true_name": data.TrueName,
			"face":      data.Face,
			"id_number": data.IdNumber,
			"address":   data.Address,
			"mobile":    data.Mobile,
			"pinyin":    data.Pinyin,
		}).From(m.Table).Where(builder.Eq{"id": data.Id}).ToSQL()
		return conn.Exec(query, args...)
	}, m.keys(data)...)
	return err
}

func (m *ParentModel) DeleteByOrgId(orgId int64) error {
	query, args, _ := builder.Postgres().Update(builder.Eq{"deleted": true}).From(m.Table).Where(builder.Eq{"org_id": orgId, "deleted": false}).ToSQL()
	_, err := m.ExecNoCache(query, args...)
	return err
}

func (m *ParentModel) DeleteByClassId(classId int64) error {
	query, args, _ := builder.Postgres().Update(builder.Eq{"deleted": true}).From(m.Table).Where(builder.Eq{"class_id": classId, "deleted": false}).ToSQL()
	_, err := m.ExecNoCache(query, args...)
	return err
}

func (m *ParentModel) DeleteByParentIds(ids []int64) error {
	_, err := m.Exec(func(conn sqlx.SqlConn) (sql.Result, error) {
		query, args, _ := builder.Postgres().Update(builder.Eq{"deleted": true}).From(m.Table).Where(builder.In("id", ids)).ToSQL()
		return conn.Exec(query, args...)
	}, m.idKeys(ids)...)
	return err
}

func (m *ParentModel) FindOneByMobile(mobile string) (*Parent, error) {
	var resp Parent
	parentInfoByMobileKey := fmt.Sprintf("%s%s", cacheParentInfoByMobilePrefix, mobile)
	err := m.QueryRow(&resp, parentInfoByMobileKey, func(conn sqlx.SqlConn, v interface{}) error {
		query, args, _ := builder.Postgres().Select(parentFieldNames).From(m.Table).Where(builder.Eq{"mobile": mobile, "deleted": false}).Limit(1).ToSQL()
		return conn.QueryRow(v, query, args...)
	})
	if err != nil {
		return nil, cachemodel.ErrNotFound
	}
	return &resp, nil
}

func (m *ParentModel) FindOneByIdNumber(idNumber string) (*Parent, error) {
	var resp Parent
	query, args, _ := builder.Postgres().Select(parentFieldNames).From(m.Table).Where(builder.Eq{"id_number": idNumber, "deleted": false}).Limit(1).ToSQL()
	err := m.Conn.QueryRow(&resp, query, args...)
	if err != nil {
		return nil, cachemodel.ErrNotFound
	}
	return &resp, nil
}

func (m *ParentModel) FindOneById(id int64) (*Parent, error) {
	var resp Parent
	studentInfoByStudentIdKey := m.formatPrimary(id)
	err := m.QueryRow(&resp, studentInfoByStudentIdKey, func(conn sqlx.SqlConn, v interface{}) error {
		query, args, _ := builder.Postgres().Select(parentFieldNames).From(m.Table).Where(builder.Eq{"id": id, "deleted": false}).Limit(1).ToSQL()
		return conn.QueryRow(v, query, args...)
	})
	if err != nil {
		return nil, cachemodel.ErrNotFound
	}
	return &resp, nil
}

func (m *ParentModel) FindListByConditions(parentIds []int64, parentName string, faceStatus int8, page, limit int) ([]Parent, int, error) {
	var res []Parent
	var count int
	eq := builder.And()
	eq = eq.And(builder.In("id", parentIds))
	if parentName != "" {
		eq = eq.And(builder.Eq{"true_name": parentName})
	}
	if faceStatus != -1 {
		eq = eq.And(builder.Eq{"face_status": faceStatus})
	}
	query := builder.Postgres().Select(parentFieldNames).From(m.Table).Where(eq).And(builder.Eq{"deleted": false})
	if page > 0 || limit > 0 {
		query = query.Limit(limit, (page-1)*limit).OrderBy("pinyin")
	}
	sqlQ, argsQ, _ := query.ToSQL()
	sqlC, argsC, _ := builder.Postgres().Select("COUNT(*)").From(m.Table).Where(eq).And(builder.Eq{"deleted": false}).ToSQL()
	err := m.Conn.QueryRows(&res, sqlQ, argsQ...)
	err = m.Conn.QueryRow(&count, sqlC, argsC...)
	return res, count, err
}

func (m *ParentModel) formatPrimary(primary interface{}) string {
	return fmt.Sprintf("%s%v", cacheParentInfoByParentIdPrefix, primary)
}

func (m *ParentModel) queryPrimary(conn sqlx.SqlConn, v, primary interface{}) error {
	query := fmt.Sprintf("select %s from %s where `id` = $1 and deleted = 0 limit 1", parentFieldNames, m.Table)
	return conn.QueryRow(v, query, primary)
}

func (m *ParentModel) keys(data *Parent) []string {
	res := make([]string, 0, 3)
	if data.UserName != "" {
		res = append(res, fmt.Sprintf("%s%s", cacheParentInfoByMobilePrefix, data.UserName))
	}
	if data.TrueName != "" {
		res = append(res, fmt.Sprintf("%s%s", cacheParentInfoByTrueNamePrefix, data.TrueName))
	}
	if data.Id != 0 {
		res = append(res, m.formatPrimary(data.Id))
	}
	return res
}

func (m *ParentModel) idKeys(ids []int64) []string {
	res := make([]string, 0, len(ids))
	for _, id := range ids {
		res = append(res, m.formatPrimary(id))
	}
	return res
}
