package model

import (
	"database/sql"
	"demacia/common/basefunc"
	"demacia/common/cachemodel"
	"fmt"

	"github.com/go-xorm/builder"
	"gitlab.u-jy.cn/xiaoyang/go-zero/core/stores/redis"
	"strings"

	"gitlab.u-jy.cn/xiaoyang/go-zero/core/stores/sqlx"
)

var (
	subjectTeacherFieldNames = strings.Join(cachemodel.RawFieldNames(&SubjectTeacher{}, true), ",")

	cacheSubjectSubjectTeacherIdPrefix        = "cache:subject:subjectTeacher:id:"
	cacheSubjectSubjectTeacherSubjectIdPrefix = "cache:subject:subjectTeacher:subject-id:"
	cacheSubjectSubjectTeacherOrgIdPrefix     = "cache:subject:subjectTeacher:org-id:"
)

type (
	SubjectTeacherModel struct {
		*cachemodel.CachedModel
	}

	SubjectTeacher struct {
		Id        int64  `db:"id"`
		SubjectId int64  `db:"subject_id"` // 学科id
		MemberId  int64  `db:"member_id"`  // 教师(人员)id
		TrueName  string `db:"true_name"`  // 教师名称
		OrgId     int64  `db:"org_id"`     // 机构id
	}
	SubjectTeachers []*SubjectTeacher
)

func NewSubjectTeacherModel(conn sqlx.SqlConn, cache *redis.Redis) *SubjectTeacherModel {
	return &SubjectTeacherModel{
		CachedModel: cachemodel.NewCachedModel(conn, `"subject"."subject_teacher"`, cache),
	}
}

func (m *SubjectTeacherModel) ListSubjectTeacherByOrgIdAndSubject(OrgId, SubjectId int64) (*SubjectTeachers, error) {
	var res SubjectTeachers
	listSubjectTeacherBySubjectIdKey := fmt.Sprintf("%s%d", cacheSubjectSubjectTeacherSubjectIdPrefix, SubjectId)
	err := m.QueryRow(&res, listSubjectTeacherBySubjectIdKey, func(conn sqlx.SqlConn, v interface{}) error {
		query, args, _ := builder.Postgres().
			Select(subjectTeacherFieldNames).
			From(m.Table).
			Where(builder.Eq{"org_id": OrgId, "subject_id": SubjectId, "deleted": false}).
			ToSQL()
		return conn.QueryRows(v, query, args...)
	})
	if err != nil {
		return nil, err
	}
	return &res, nil
}
func (m *SubjectTeacherModel) DeletedBySubjectId(SubjectId int64) error {
	listSubjectTeacherBySubjectIdKey := fmt.Sprintf("%s%d", cacheSubjectSubjectTeacherSubjectIdPrefix, SubjectId)

	_, err := m.Exec(func(conn sqlx.SqlConn) (res sql.Result, err error) {
		sqlString, args, _ := builder.Postgres().Update(builder.Eq{
			"deleted": true,
		}).Where(builder.Eq{"subject_id": SubjectId}).From(m.Table).ToSQL()
		return conn.Exec(sqlString, args...)
	}, listSubjectTeacherBySubjectIdKey)
	if err != nil {
		return err
	}
	_ = m.DelCache(listSubjectTeacherBySubjectIdKey)
	return err
}

func (m *SubjectTeacherModel) BatchInsert(SubjectTeachers SubjectTeachers) error {

	insertString, args, err := basefunc.BatchInsertString(m.Table, SubjectTeachers)
	if err != nil {
		return err
	}
	_, err = m.ExecNoCache(insertString, args...)
	for _, v := range SubjectTeachers {
		batchInsertKey := fmt.Sprintf("%s%d", cacheSubjectSubjectTeacherSubjectIdPrefix, v.SubjectId)
		_ = m.DelCache(batchInsertKey)
	}
	return err
}

func (m *SubjectTeacherModel) formatPrimary(primary interface{}) string {
	return fmt.Sprintf("%s%v", cacheSubjectSubjectTeacherIdPrefix, primary)
}

func (m *SubjectTeacherModel) queryPrimary(conn sqlx.SqlConn, v, primary interface{}) error {
	query := fmt.Sprintf("select %s from %s where `id` = $1 and deleted = 0 limit 1", subjectTeacherFieldNames, m.Table)
	return conn.QueryRow(v, query, primary)
}

func (m *SubjectTeacherModel) keys(subjectTeacher *SubjectTeacher) []string {
	res := make([]string, 0)
	if subjectTeacher.OrgId != 0 {
		res = append(res, fmt.Sprintf("%s%d", cacheSubjectSubjectTeacherOrgIdPrefix, subjectTeacher.OrgId))
	}
	if subjectTeacher.SubjectId != 0 {
		res = append(res, fmt.Sprintf("%s%d", cacheSubjectSubjectTeacherOrgIdPrefix, subjectTeacher.SubjectId))
	}
	if subjectTeacher.Id != 0 {
		res = append(res, m.formatPrimary(subjectTeacher.Id))
	}
	return res
}
