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
	studentParentFieldNames                   = strings.Join(cachemodel.RawFieldNames(&StudentParent{}, true), ",")
	cacheStudentParentInfoByIdPrefix          = "cache:parent:student_parent:id:"
	cacheStudentParentInfoByClassIdPrefix     = "cache:parent:student_parent:class-id:"
	cacheStudentParentInfoByStudentNamePrefix = "cache:parent:student_parent:student-name:"
	cacheStudentParentInfoByParentIdPrefix    = "cache:parent:student_parent:parent-id:"
)

type (
	StudentParentModel struct {
		*cachemodel.CachedModel
	}
	StudentParent struct {
		Id          int64  `db:"id"`
		OrgId       int64  `db:"org_id"`
		ClassId     int64  `db:"class_id"`
		ClassName   string `db:"class_name"`
		ParentId    int64  `db:"parent_id"`
		StudentId   int64  `db:"student_id"`
		StudentName string `db:"student_name"`
		Relation    int8   `db:"relation"`
	}
)

func NewStudentParentModel(conn sqlx.SqlConn, cache *redis.Redis) *StudentParentModel {
	return &StudentParentModel{
		CachedModel: cachemodel.NewCachedModel(conn, `"parent".student_parent`, cache),
	}
}

func (m *StudentParentModel) InsertOne(data *StudentParent) (int64, error) {
	query, args, _ := builder.Postgres().Insert(builder.Eq{
		"org_id":       data.OrgId,
		"class_id":     data.ClassId,
		"class_name":   data.ClassName,
		"parent_id":    data.ParentId,
		"student_id":   data.StudentId,
		"student_name": data.StudentName,
		"relation":     data.Relation,
	}).Into(m.Table).ToSQL()
	var res int64
	err := m.Conn.QueryRow(&res, query+" returning id", args...)
	if err != nil {
		return 0, err
	}
	_ = m.DelCache(m.keys(data)...)
	return res, nil
}

func (m *StudentParentModel) UpdateOne(data *StudentParent) error {
	_, err := m.Exec(func(conn sqlx.SqlConn) (sql.Result, error) {
		query, args, _ := builder.Postgres().Update(builder.Eq{
			"org_id":       data.OrgId,
			"class_id":     data.ClassId,
			"class_name":   data.ClassName,
			"parent_id":    data.ParentId,
			"student_id":   data.StudentId,
			"student_name": data.StudentName,
			"relation":     data.Relation,
		}).From(m.Table).Where(builder.Eq{"id": data.Id}).ToSQL()
		return conn.Exec(query, args...)
	}, m.keys(data)...)
	return err
}

func (m *StudentParentModel) DeleteByParentIds(ids []int64) error {
	_, err := m.Exec(func(conn sqlx.SqlConn) (sql.Result, error) {
		query, args, _ := builder.Postgres().Update(builder.Eq{"deleted": true}).From(m.Table).Where(builder.In("parent_id", ids)).ToSQL()
		return conn.Exec(query, args...)
	}, m.idKeys(ids)...)
	return err
}

func (m *StudentParentModel) FindOneById(id int64) (*StudentParent, error) {
	var resp StudentParent
	err := m.QueryRow(&resp, m.formatPrimary(id), func(conn sqlx.SqlConn, v interface{}) error {
		query, args, _ := builder.Postgres().Select(studentParentFieldNames).From(m.Table).Where(builder.Eq{"id": id, "deleted": false}).Limit(1).ToSQL()
		return conn.QueryRow(v, query, args...)
	})
	if err != nil {
		return nil, cachemodel.ErrNotFound
	}
	return &resp, nil
}

func (m *StudentParentModel) FindListByParentId(id int64) (*[]StudentParent, error) {
	var resp []StudentParent
	err := m.QueryRow(&resp, fmt.Sprintf("%s%d", cacheStudentParentInfoByParentIdPrefix, id), func(conn sqlx.SqlConn, v interface{}) error {
		query, args, _ := builder.Postgres().Select(studentParentFieldNames).From(m.Table).Where(builder.Eq{"parent_id": id, "deleted": false}).ToSQL()
		return conn.QueryRows(v, query, args...)
	})
	if err != nil {
		return nil, cachemodel.ErrNotFound
	}
	return &resp, nil
}

func (m *StudentParentModel) FindListByConditions(orgId, claasId int64, studentName string, page, limit int) ([]StudentParent, int, error) {
	var res []StudentParent
	var count int
	eq := builder.And()
	if claasId != -1 {
		eq = eq.And(builder.Eq{"class_id": claasId})
	}
	if orgId != -1 {
		eq = eq.And(builder.Eq{"org_id": orgId})
	}
	if studentName != "" {
		eq = eq.And(builder.Eq{"student_name": studentName})
	}
	query := builder.Postgres().Select(studentParentFieldNames).From(m.Table).Where(eq).And(builder.Eq{"deleted": false})
	if page > 0 || limit > 0 {
		query = query.Limit(limit, (page-1)*limit).OrderBy("id")
	}
	sqlQ, argsQ, _ := query.ToSQL()
	sqlC, argsC, _ := builder.Postgres().Select("COUNT(*)").From(m.Table).Where(eq).And(builder.Eq{"deleted": false}).ToSQL()
	err := m.Conn.QueryRows(&res, sqlQ, argsQ...)
	err = m.Conn.QueryRow(&count, sqlC, argsC...)
	return res, count, err
}

func (m *StudentParentModel) FindOneByParentIdAndStudentId(parentId, studentId int64) (*StudentParent, error) {
	var resp StudentParent
	query, args, _ := builder.Postgres().Select(studentParentFieldNames).From(m.Table).Where(builder.Eq{"parent_id": parentId, "student_id": studentId, "deleted": false}).Limit(1).ToSQL()
	err := m.Conn.QueryRow(&resp, query, args...)
	if err != nil {
		return nil, cachemodel.ErrNotFound
	}
	return &resp, nil
}

func (m *StudentParentModel) formatPrimary(primary interface{}) string {
	return fmt.Sprintf("%s%v", cacheStudentParentInfoByIdPrefix, primary)
}

func (m *StudentParentModel) queryPrimary(conn sqlx.SqlConn, v, primary interface{}) error {
	query := fmt.Sprintf("select %s from %s where `id` = $1 and deleted = 0 limit 1", studentParentFieldNames, m.Table)
	return conn.QueryRow(v, query, primary)
}

func (m *StudentParentModel) keys(data *StudentParent) []string {
	res := make([]string, 0, 4)
	if data.ClassId != 0 {
		res = append(res, fmt.Sprintf("%s%d", cacheStudentParentInfoByClassIdPrefix, data.ClassId))
	}
	if data.StudentName != "" {
		res = append(res, fmt.Sprintf("%s%s", cacheStudentParentInfoByStudentNamePrefix, data.StudentName))
	}
	if data.Id != 0 {
		res = append(res, m.formatPrimary(data.Id))
	}
	if data.ParentId != 0 {
		res = append(res, fmt.Sprintf("%s%d", cacheStudentParentInfoByParentIdPrefix, data.ParentId))
	}
	return res
}

func (m *StudentParentModel) idKeys(ids []int64) []string {
	res := make([]string, 0, len(ids))
	for _, id := range ids {
		res = append(res, m.formatPrimary(id))
	}
	return res
}
