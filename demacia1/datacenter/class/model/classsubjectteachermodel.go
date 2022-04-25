package model

import (
	"demacia/common/cachemodel"
	"gitlab.u-jy.cn/xiaoyang/go-zero/core/stores/redis"
	"strings"

	"gitlab.u-jy.cn/xiaoyang/go-zero/core/stores/sqlx"
	"xorm.io/builder"
)

var (
	classSubjectTeacherFieldNames = strings.Join(cachemodel.RawFieldNames(&ClassSubjectTeacher{}, true), ",")

	cacheClassClassSubjectTeacherIdPrefix = "cache:class:classSubjectTeacher:id:"
)

type (
	ClassSubjectTeacherModel struct {
		*cachemodel.CachedModel
	}

	ClassSubjectTeacher struct {
		Id           int64  `db:"id"`
		ClassId      int64  `db:"class_id"`      // 班级id
		SubjectId    int64  `db:"subject_id"`    // 学科id
		SubjectTitle string `db:"subject_title"` // 学科名称
		MemberId     int64  `db:"member_id"`     // 教师id
		TrueName     string `db:"true_name"`     // 教师名称
	}

	ClassSubjectTeachers []ClassSubjectTeacher
)

func NewClassSubjectTeacherModel(conn sqlx.SqlConn, cache *redis.Redis) *ClassSubjectTeacherModel {
	return &ClassSubjectTeacherModel{
		CachedModel: cachemodel.NewCachedModel(conn, `"class"."class_subject_teacher"`, cache),
	}

}

func (m *ClassSubjectTeacherModel) ListByClassId(ClassId int64) (ClassSubjectTeachers, error) {
	sql, args, _ := builder.Postgres().
		Select(classSubjectTeacherFieldNames).
		From(m.Table).
		Where(builder.Eq{"class_id": ClassId, "deleted": false}).
		ToSQL()
	var res ClassSubjectTeachers
	err := m.Conn.QueryRows(&res, sql, args...)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (m *ClassSubjectTeacherModel) FindSubjectBySubjectIdAndClassId(subjectId, ClassId int64) (*ClassSubjectTeacher, error) {
	sql, args, _ := builder.Postgres().
		Select(classSubjectTeacherFieldNames).
		From(m.Table).
		Where(builder.Eq{"class_id": ClassId, "subject_id": subjectId, "deleted": false}).
		ToSQL()
	var res ClassSubjectTeacher
	err := m.Conn.QueryRow(&res, sql, args...)
	if err != nil {
		return nil, err
	}
	return &res, nil
}
func (m *ClassSubjectTeacherModel) FindTeacherByMemberIdAndClassId(memberId, classId int64) (*ClassSubjectTeacher, error) {
	sql, args, _ := builder.Postgres().
		Select(classSubjectTeacherFieldNames).
		From(m.Table).
		Where(builder.Eq{"member_id": memberId, "class_id": classId, "deleted": false}).
		ToSQL()
	var res ClassSubjectTeacher
	err := m.Conn.QueryRow(&res, sql, args...)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func (m *ClassSubjectTeacherModel) Insert(teach *ClassSubjectTeacher) (int64, error) {
	sqlString, args, _ := builder.Postgres().Insert(builder.Eq{
		"class_id":      teach.ClassId,
		"subject_id":    teach.SubjectId,
		"subject_title": teach.SubjectTitle,
		"member_id":     teach.MemberId,
		"true_name":     teach.TrueName,
	}).Into(m.Table).ToSQL()
	var LastInsertId int64
	err := m.Conn.QueryRow(&LastInsertId, sqlString+" returning id", args...)
	return LastInsertId, err
}

func (m *ClassSubjectTeacherModel) Update(teach *ClassSubjectTeacher) error {
	sqlString, args, _ := builder.Postgres().Update(builder.Eq{
		"member_id": teach.MemberId,
		"true_name": teach.TrueName,
	}).Where(builder.Eq{
		"class_id":   teach.ClassId,
		"subject_id": teach.SubjectId,
		"deleted":    false,
	}).
		From(m.Table).ToSQL()
	_, err := m.ExecNoCache(sqlString, args...)
	return err
}

func (m *ClassSubjectTeacherModel) DeleteBySubjectIdAndClassId(subjectId, classId int64) error {
	sqlString, args, _ := builder.Postgres().Update(builder.Eq{
		"deleted": true,
	}).Where(builder.Eq{
		"subject_id": subjectId,
		"class_id":   classId,
	}).From(m.Table).ToSQL()
	_, err := m.ExecNoCache(sqlString, args...)
	if err != nil {
		return err
	}
	return nil
}

func (m *ClassSubjectTeacherModel) FindTeacherByClassIdAndSubjectId(classId, subjectId int64) (*ClassSubjectTeacher, error) {
	sql, args, _ := builder.Postgres().
		Select(classSubjectTeacherFieldNames).
		From(m.Table).
		Where(builder.Eq{"subject_id": subjectId, "class_id": classId, "deleted": false}).
		ToSQL()
	var res ClassSubjectTeacher
	err := m.Conn.QueryRow(&res, sql, args...)
	if err != nil {
		return nil, err
	}
	return &res, nil
}
