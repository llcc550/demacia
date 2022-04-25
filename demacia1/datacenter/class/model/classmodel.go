package model

import (
	"database/sql"
	"fmt"
	"strings"

	"demacia/common/basefunc"
	"demacia/common/cachemodel"

	"gitlab.u-jy.cn/xiaoyang/go-zero/core/stores/redis"
	"gitlab.u-jy.cn/xiaoyang/go-zero/core/stores/sqlx"
	"xorm.io/builder"
)

var (
	classFieldNames               = strings.Join(cachemodel.RawFieldNames(&Class{}, true), ",")
	cacheClassClassIdPrefix       = "cache:class:class:id:"
	cacheClassListByGradeIdPrefix = "cache:class:class:grade-id:"
	cacheClassListByOrgIdPrefix   = "cache:class:class:org-id:"
	cacheClassClassFullNamePrefix = "cache:class:class:full-name:"
)

type (
	ClassModel struct {
		*cachemodel.CachedModel
	}

	Class struct {
		Id             int64  `db:"id"`
		OrgId          int64  `db:"org_id"`           // 所属机构
		Title          string `db:"title"`            // 班级名称
		StageId        int64  `db:"stage_id"`         // 学段ID。冗余数据
		StageTitle     string `db:"stage_title"`      // 学段名称。冗余数据
		GradeId        int64  `db:"grade_id"`         // 年级ID。冗余数据
		GradeTitle     string `db:"grade_title"`      // 年级名称。冗余数据
		FullName       string `db:"full_name"`        // 班级全称，原则上全校唯一
		AliasName      string `db:"alias_name"`       // 班级全称，原则上全校唯一
		ClassMemberNum int64  `db:"class_member_num"` // 班级人数
		Desc           string `db:"description"`      // 班级宣传语
		Sort           int8   `db:"sort"`             // 班级序号
		Teachers       string `db:"teachers"`         // 班主任
	}

	Classs []*Class
)

func NewClassModel(conn sqlx.SqlConn, cache *redis.Redis) *ClassModel {
	return &ClassModel{
		CachedModel: cachemodel.NewCachedModel(conn, `"class"."class"`, cache),
	}
}

func (m *ClassModel) ListByOrgIdAndGradeId(orgId, gradeId int64) (Classs, error) {
	var res Classs
	classListByGradeIdKey := fmt.Sprintf("%s%d", cacheClassListByGradeIdPrefix, gradeId)
	err := m.QueryRow(&res, classListByGradeIdKey, func(conn sqlx.SqlConn, v interface{}) error {
		query, args, _ := builder.Postgres().
			Select(classFieldNames).
			From(m.Table).
			Where(builder.Eq{"org_id": orgId, "grade_id": gradeId, "deleted": false}).
			OrderBy("sort ASC ").
			ToSQL()
		return conn.QueryRows(v, query, args...)
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (m *ClassModel) ListByOrgId(orgId int64) (Classs, error) {
	var res Classs
	classListByOrgIdKey := fmt.Sprintf("%s%d", cacheClassListByOrgIdPrefix, orgId)
	err := m.QueryRow(&res, classListByOrgIdKey, func(conn sqlx.SqlConn, v interface{}) error {
		query, args, _ := builder.Postgres().
			Select(classFieldNames).
			From(m.Table).
			Where(builder.Eq{"org_id": orgId, "deleted": false}).
			OrderBy("sort ASC ").
			ToSQL()
		return conn.QueryRows(v, query, args...)
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (m *ClassModel) PageListByOrgId(orgId, stageId, gradeId, classId int64, teacherName string, page, limit int) (Classs, int, error) {
	var res Classs
	var total int
	var offset int
	if page == 0 {
		page = 1
	}
	if limit > 0 && limit < 100 {
		offset = (page - 1) * limit
	}
	eq := builder.And(builder.Eq{"deleted": false})
	if orgId > 0 {
		eq = eq.And(builder.Eq{"org_id": orgId})
	}
	if stageId > 0 {
		eq = eq.And(builder.Eq{"stage_id": stageId})
	}
	if gradeId > 0 {
		eq = eq.And(builder.Eq{"grade_id": gradeId})
	}
	if classId > 0 {
		eq = eq.And(builder.Eq{"grade_id": gradeId})
	}
	if teacherName != "" {
		eq = eq.And(builder.Like{"teachers", teacherName})
	}
	query, args, _ := builder.Postgres().
		Select(classFieldNames).
		From(m.Table).
		Where(eq).
		Limit(limit, offset).
		OrderBy("id DESC,sort Asc ").
		ToSQL()
	err := m.QueryRowsNoCache(&res, query, args...)
	if err != nil {
		return nil, 0, err
	}
	totalQuery, totalArgs, _ := builder.Postgres().
		Select("COUNT(id)").
		From(m.Table).
		Where(eq).
		ToSQL()
	err = m.QueryRowNoCache(&total, totalQuery, totalArgs...)
	if err != nil {
		return nil, 0, err
	}
	return res, total, nil
}

func (m *ClassModel) GetClassById(Id int64) (*Class, error) {
	var res Class
	classListByOrgIdKey := m.formatPrimary(Id)
	err := m.QueryRow(&res, classListByOrgIdKey, func(conn sqlx.SqlConn, v interface{}) error {
		query, args, _ := builder.Postgres().Select(classFieldNames).From(m.Table).Where(builder.Eq{"id": Id, "deleted": false}).ToSQL()
		return conn.QueryRow(v, query, args...)
	})
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func (m *ClassModel) Insert(class *Class) (int64, error) {
	query, args, _ := builder.Postgres().Insert(builder.Eq{
		"org_id":           class.OrgId,
		"title":            class.Title,
		"stage_id":         class.StageId,
		"stage_title":      class.StageTitle,
		"grade_id":         class.GradeId,
		"grade_title":      class.GradeTitle,
		"full_name":        class.FullName,
		"alias_name":       class.AliasName,
		"class_member_num": 0,
		"description":      class.Desc,
		"sort":             class.Sort,
		"teachers":         class.Teachers,
	}).Into(m.Table).ToSQL()
	var classId int64
	err := m.Conn.QueryRow(&classId, query+" returning id", args...)
	if err != nil {
		return 0, err
	}
	_ = m.DelCache(m.keys(class)...)
	return classId, nil
}

func (m *ClassModel) Update(class *Class) error {
	classInfo, err := m.GetClassById(class.Id)
	if err != nil || classInfo.OrgId != class.OrgId {
		return cachemodel.ErrNotFound
	}
	//if class.Teachers != "" {
	//	upEq := builder.Eq{
	//		"teachers": class.Teachers,
	//	}
	//}
	_, err = m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query, args, _ := builder.Postgres().Update(builder.Eq{
			"alias_name":  class.AliasName,
			"teachers":    class.Teachers,
			"description": class.Desc,
		}).Where(builder.Eq{
			"id":      class.Id,
			"deleted": false,
		}).From(m.Table).ToSQL()
		return conn.Exec(query, args...)
	}, m.keys(classInfo)...)
	return err
}

func (m *ClassModel) ChangeStudentNumByClassId(classInfo *Class, Num int64) error {
	_, err := m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query, args, _ := builder.Postgres().Update(builder.Eq{"class_member_num": Num}).Where(builder.Eq{"id": classInfo.Id, "deleted": false}).From(m.Table).ToSQL()
		return conn.Exec(query, args...)
	}, m.keys(classInfo)...)
	return err
}

func (m *ClassModel) GetStudentTotalNumByGradeId(GradeId int64) (int64, error) {
	query, args, _ := builder.Postgres().Select("SUM(class_member_num) as total").Where(builder.Eq{"grade_id": GradeId, "deleted": false}).From(m.Table).ToSQL()
	var total int64
	err := m.Conn.QueryRow(&total, query, args...)
	return total, err
}

func (m *ClassModel) BatchInsert(classes *Classs) error {
	query, args, err := basefunc.BatchInsertString(m.Table, classes)
	if err != nil {
		return err
	}
	_, err = m.ExecNoCache(query, args...)
	return err
}

func (m *ClassModel) DeleteById(Id int64) error {
	classInfo, err := m.GetClassById(Id)
	if err != nil {
		return cachemodel.ErrNotFound
	}
	_, err = m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query, args, _ := builder.Postgres().Update(builder.Eq{"deleted": true}).Where(builder.Eq{"id": Id}).From(m.Table).ToSQL()
		return conn.Exec(query, args...)
	}, m.keys(classInfo)...)
	return err
}

func (m *ClassModel) DeleteByOrgId(orgId int64) error {
	_, err := m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query, args, _ := builder.Postgres().Update(builder.Eq{"deleted": true}).Where(builder.Eq{"org_id": orgId}).From(m.Table).ToSQL()
		return conn.Exec(query, args...)
	}, fmt.Sprintf("%s%d", cacheClassListByOrgIdPrefix, orgId))
	return err
}

func (m *ClassModel) GetClassByFullNameAndOrgId(orgId int64, fullName string) (*Class, error) {
	var resp Class
	key := fmt.Sprintf("%s%d:%s", cacheClassClassFullNamePrefix, orgId, fullName)
	err := m.QueryRowIndex(&resp, key, m.formatPrimary, func(conn sqlx.SqlConn, v interface{}) (i interface{}, e error) {
		query, args, _ := builder.Postgres().From(m.Table).Select(classFieldNames).Where(builder.Eq{"org_id": orgId, "full_name": fullName, "deleted": false}).Limit(1).ToSQL()
		if err := conn.QueryRow(&resp, query, args...); err != nil {
			return nil, err
		}
		return resp.Id, nil
	}, m.queryPrimary)
	if err != nil {
		return nil, cachemodel.ErrNotFound
	}
	return &resp, nil
}

func (m *ClassModel) formatPrimary(primary interface{}) string {
	return fmt.Sprintf("%s%v", cacheClassClassIdPrefix, primary)
}

func (m *ClassModel) queryPrimary(conn sqlx.SqlConn, v, primary interface{}) error {
	query := fmt.Sprintf("select %s from %s where `id` = ? and deleted = 0 limit 1", classFieldNames, m.Table)
	return conn.QueryRow(v, query, primary)
}

func (m *ClassModel) keys(class *Class) []string {
	res := make([]string, 0, 4)
	if class.OrgId != 0 {
		res = append(res, fmt.Sprintf("%s%d", cacheClassListByOrgIdPrefix, class.OrgId))
	}
	if class.GradeId != 0 {
		res = append(res, fmt.Sprintf("%s%d", cacheClassListByGradeIdPrefix, class.GradeId))
	}
	if class.OrgId != 0 && class.FullName != "" {
		res = append(res, fmt.Sprintf("%s%d:%s", cacheClassClassFullNamePrefix, class.OrgId, class.FullName))
	}
	if class.Id != 0 {
		res = append(res, m.formatPrimary(class.Id))
	}
	return res
}
