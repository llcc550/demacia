package model

import (
	"demacia/common/basefunc"
	"demacia/common/cachemodel"
	"fmt"
	"gitlab.u-jy.cn/xiaoyang/go-zero/core/stores/redis"
	"strings"

	"github.com/go-xorm/builder"
	"gitlab.u-jy.cn/xiaoyang/go-zero/core/stores/sqlx"
)

var (
	subjectGradeFieldNames = strings.Join(cachemodel.RawFieldNames(&SubjectGrade{}, true), ",")

	cacheSubjectSubjectGradeIdPrefix        = "cache:subject:subjectGrade:id:"
	cacheSubjectSubjectGradeOrgIdPrefix     = "cache:subject:subjectGrade:org-id:"
	cacheSubjectSubjectGradeSubjectIdPrefix = "cache:subject:subjectGrade:subject-id:"
	cacheSubjectSubjectGradeGradeIdPrefix   = "cache:subject:subjectGrade:grade-id:"
)

type (
	SubjectGradeModel struct {
		*cachemodel.CachedModel
	}

	SubjectGrade struct {
		Id           int64  `db:"id"`
		OrgId        int64  `db:"org_id"`        // 机构ID
		GradeId      int64  `db:"grade_id"`      // 年级ID
		SubjectId    int64  `db:"subject_id"`    // 学科ID
		SubjectTitle string `db:"subject_title"` // 学科名称
		GradeTitle   string `db:"grade_title"`   // 年级名称
	}
	Grade struct {
		Id         int64  `db:"id"`
		OrgId      int64  `db:"org_id"`      // 机构ID
		GradeId    int64  `db:"grade_id"`    // 年级ID
		GradeTitle string `db:"grade_title"` // 年级名称
	}
	SubjectGrades []*SubjectGrade
)

func NewSubjectGradeModel(conn sqlx.SqlConn, cache *redis.Redis) *SubjectGradeModel {
	return &SubjectGradeModel{
		CachedModel: cachemodel.NewCachedModel(conn, `"subject"."subject_grade"`, cache),
	}
}

func (m *SubjectGradeModel) ListSubjectGradeByOrgId(orgId int64) (*SubjectGrades, error) {
	var res SubjectGrades
	subjectGradeByOrgIdKey := fmt.Sprintf("%s%d", cacheSubjectSubjectGradeOrgIdPrefix, orgId)
	err := m.QueryRow(&res, subjectGradeByOrgIdKey, func(conn sqlx.SqlConn, v interface{}) error {
		sql, args, _ := builder.Postgres().
			Select(subjectGradeFieldNames).
			From(m.Table).
			Where(builder.Eq{"org_id": orgId, "deleted": false}).
			ToSQL()
		return conn.QueryRows(v, sql, args...)
	})
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func (m *SubjectGradeModel) ListSubjectByGradeId(GradeId int64) (*SubjectGrades, error) {
	var res SubjectGrades
	subjectByGradeIdKey := fmt.Sprintf("%s%d", cacheSubjectSubjectGradeGradeIdPrefix, GradeId)
	err := m.QueryRow(&res, subjectByGradeIdKey, func(conn sqlx.SqlConn, v interface{}) error {
		sql, args, _ := builder.Postgres().
			Select(subjectGradeFieldNames).
			From(m.Table).
			Where(builder.Eq{"grade_id": GradeId, "deleted": false}).
			ToSQL()
		return conn.QueryRows(v, sql, args...)
	})
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func (m *SubjectGradeModel) ListGradeBySubjectId(SubjectId int64) (*SubjectGrades, error) {
	var res SubjectGrades
	listGradeBySubjectIdKey := fmt.Sprintf("%s%d", cacheSubjectSubjectGradeSubjectIdPrefix, SubjectId)
	err := m.QueryRow(&res, listGradeBySubjectIdKey, func(conn sqlx.SqlConn, v interface{}) error {
		query, args, _ := builder.Postgres().
			Select(subjectGradeFieldNames).
			From(m.Table).
			Where(builder.Eq{"subject_id": SubjectId, "deleted": false}).
			ToSQL()
		return conn.QueryRows(v, query, args...)
	})
	if err != nil {
		return nil, err
	}
	return &res, nil
}
func (m *SubjectGradeModel) FindOneById(id int64) (*SubjectGrade, error) {
	var res SubjectGrade
	findOneIdKey := fmt.Sprintf("%s%d", cacheSubjectSubjectGradeIdPrefix, id)
	err := m.QueryRow(&res, findOneIdKey, func(conn sqlx.SqlConn, v interface{}) error {
		query, args, _ := builder.Postgres().
			Select(subjectGradeFieldNames).
			From(m.Table).
			Where(builder.Eq{"id": id, "deleted": false}).
			ToSQL()
		return conn.QueryRows(v, query, args...)
	})
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func (m *SubjectGradeModel) Insert(SubjectGrade *SubjectGrade) (int64, error) {
	sqlString, args, _ := builder.Postgres().Insert(builder.Eq{
		"org_id":        SubjectGrade.OrgId,
		"grade_id":      SubjectGrade.GradeId,
		"subject_id":    SubjectGrade.SubjectId,
		"subject_title": SubjectGrade.SubjectTitle,
		"grade_title":   SubjectGrade.GradeTitle,
	}).Into(m.Table).ToSQL()
	var LastInsertId int64
	err := m.Conn.QueryRow(&LastInsertId, sqlString+" returning id", args...)
	if err != nil {
		return 0, err
	}
	_ = m.DelCache(m.Keys(SubjectGrade)...)
	return LastInsertId, nil
}

func (m *SubjectGradeModel) BatchInsert(SubjectGrades SubjectGrades) error {
	insertString, args, err := basefunc.BatchInsertString(m.Table, SubjectGrades)
	if err != nil {
		return err
	}
	_, err = m.ExecNoCache(insertString, args...)
	for _, v := range SubjectGrades {
		_ = m.DelCache(m.Keys(v)...)
	}
	return err
}

func (m *SubjectGradeModel) DeletedByOrgId(OrgId int64) error {
	deletedByOrgIdKey := fmt.Sprintf("%s%d", cacheSubjectSubjectGradeOrgIdPrefix, OrgId)
	sqlString, args, _ := builder.Postgres().Update(builder.Eq{
		"deleted": true,
	}).Where(builder.Eq{
		"org_id": OrgId,
	}).From(m.Table).ToSQL()
	_, err := m.ExecNoCache(sqlString, args...)
	_ = m.DelCache(deletedByOrgIdKey)
	return err

}
func (m *SubjectGradeModel) UpdateSubjectGradeTitleByGradeId(gradeId int64, gradeTitle string) error {
	//deletedByOrgIdKey := fmt.Sprintf("%s%d", cacheSubjectSubjectGradeOrgIdPrefix, OrgId)
	sqlString, args, _ := builder.Postgres().Update(builder.Eq{
		"grade_title": gradeTitle,
	}).Where(builder.Eq{
		"grade_id": gradeId,
	}).From(m.Table).ToSQL()
	_, err := m.ExecNoCache(sqlString, args...)
	//_ = m.DelCache(deletedByOrgIdKey)
	return err

}

func (m *SubjectGradeModel) Deleted(SubjectGrade *SubjectGrade) error {
	where := builder.Eq{"deleted": false}.And()
	if SubjectGrade.SubjectId > 0 {
		where = where.And(builder.Eq{"subject_id": SubjectGrade.SubjectId})
	} else if SubjectGrade.OrgId > 0 {
		where = where.And(builder.Eq{"org_id": SubjectGrade.OrgId})
	} else if SubjectGrade.GradeId > 0 {
		where = where.And(builder.Eq{"grade_id": SubjectGrade.GradeId})
	}
	sqlString, args, _ := builder.Postgres().Update(builder.Eq{
		"deleted": true,
	}).Where(where).From(m.Table).ToSQL()

	_, err := m.ExecNoCache(sqlString, args...)
	_ = m.DelCache(m.Keys(SubjectGrade)...)
	return err
}

// ================

func (m *SubjectGradeModel) RenameSubjectBySubjectIdAndOrgId(OrgId, SubjectId int64, subjectTitle string) error {
	sqlString, args, _ := builder.Postgres().Update(builder.Eq{
		"subject_title": subjectTitle,
	}).Where(builder.Eq{
		"subject_id": SubjectId,
		"org_id":     OrgId,
		"deleted":    false,
	}).From(m.Table).ToSQL()
	_, err := m.ExecNoCache(sqlString, args...)
	_ = m.DelCache(fmt.Sprintf("%s%d", cacheSubjectSubjectGradeOrgIdPrefix, OrgId))
	_ = m.DelCache(fmt.Sprintf("%s%d", cacheSubjectSubjectGradeSubjectIdPrefix, SubjectId))
	return err
}

func (m *SubjectGradeModel) formatPrimary(primary interface{}) string {
	return fmt.Sprintf("%s%v", cacheSubjectSubjectGradeIdPrefix, primary)
}

func (m *SubjectGradeModel) queryPrimary(conn sqlx.SqlConn, v, primary interface{}) error {
	query := fmt.Sprintf("select %s from %s where `id` = $1 and deleted = 0 limit 1", subjectGradeFieldNames, m.Table)
	return conn.QueryRow(v, query, primary)
}

func (m *SubjectGradeModel) Keys(SubjectGrade *SubjectGrade) []string {
	res := make([]string, 0)
	if SubjectGrade.OrgId > 0 {
		res = append(res, fmt.Sprintf("%s%d", cacheSubjectSubjectGradeOrgIdPrefix, SubjectGrade.OrgId))
	}
	if SubjectGrade.SubjectId > 0 {
		res = append(res, fmt.Sprintf("%s%d", cacheSubjectSubjectGradeSubjectIdPrefix, SubjectGrade.SubjectId))
	}
	if SubjectGrade.Id != 0 {
		res = append(res, m.formatPrimary(SubjectGrade.Id))
	}
	return res
}
