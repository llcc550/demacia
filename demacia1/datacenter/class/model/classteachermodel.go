package model

import (
	"demacia/common/basefunc"
	"demacia/common/cachemodel"
	"gitlab.u-jy.cn/xiaoyang/go-zero/core/stores/redis"
	"strings"
	"xorm.io/builder"

	"gitlab.u-jy.cn/xiaoyang/go-zero/core/stores/sqlx"
)

var (
	classTeacherFieldNames = strings.Join(cachemodel.RawFieldNames(&ClassTeacher{}, true), ",")

	cacheClassClassTeacherIdPrefix = "cache:class:classTeacher:id:"
)

type (
	ClassTeacherModel struct {
		*cachemodel.CachedModel
	}
	ClassTeacher struct {
		Id          int64  `db:"id"`
		ClassId     int64  `db:"class_id"`     // 班级id
		TeacherId   int64  `db:"teacher_id"`   // 教师id
		TeacherName string `db:"teacher_name"` // 教师昵称
	}
	ClassTeachers []*ClassTeacher
)

func NewClassTeacherModel(conn sqlx.SqlConn, cache *redis.Redis) *ClassTeacherModel {
	return &ClassTeacherModel{
		CachedModel: cachemodel.NewCachedModel(conn, `"class"."class_teacher"`, cache),
	}
}

func (m *ClassTeacherModel) GetTeacherWithClassIds(ClassIds []int64) (*ClassTeachers, error) {
	sql, args, _ := builder.Postgres().
		Select(classTeacherFieldNames).
		From(m.Table).
		Where(builder.In("class_id", ClassIds)).
		And(builder.Eq{"deleted": false}). //"class_id": ClassId, "deleted": false}).
		ToSQL()
	var res ClassTeachers
	err := m.Conn.QueryRows(&res, sql, args...)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func (m *ClassTeacherModel) Insert(classTeacher *ClassTeacher) (int64, error) {
	sqlString, args, _ := builder.Postgres().Insert(builder.Eq{
		"class_id":     classTeacher.ClassId,
		"teacher_id":   classTeacher.TeacherId,
		"teacher_name": classTeacher.TeacherName,
	}).Into(m.Table).ToSQL()
	var LastInsertId int64
	err := m.Conn.QueryRow(&LastInsertId, sqlString+" returning id", args...)
	return LastInsertId, err
}
func (m *ClassTeacherModel) BatchInsert(classTeachers ClassTeachers) error {
	//sqlString, args, _ := builder.Postgres().Insert(builder.Eq{
	//	"class_id":     classTeacher.ClassId,
	//	"teacher_id":   classTeacher.TeacherId,
	//	"teacher_name": classTeacher.TeacherName,
	//}).Into(m.Table).ToSQL()
	insertString, args, err := basefunc.BatchInsertString(m.Table, classTeachers)
	if err != nil {
		return err
	}
	_, err = m.Conn.Exec(insertString, args...)

	if err != nil {
		return err
	}
	return err
}

func (m *ClassTeacherModel) Update(classTeacher *ClassTeacher) error {
	sqlString, args, _ := builder.Postgres().Update(builder.Eq{
		"teacher_id":   classTeacher.TeacherId,
		"teacher_name": classTeacher.TeacherName,
	}).Where(builder.Eq{
		"class_id": classTeacher.ClassId,
	}).From(m.Table).ToSQL()
	_, err := m.ExecNoCache(sqlString, args...)
	if err != nil {
		return err
	}
	return nil
}
func (m *ClassTeacherModel) DeletedByClassId(id int64) error {
	sqlString, args, _ := builder.Postgres().Update(builder.Eq{
		"deleted": true,
	}).Where(builder.Eq{
		"class_id": id,
	}).From(m.Table).ToSQL()
	_, err := m.ExecNoCache(sqlString, args...)
	if err != nil {
		return err
	}
	return nil
}

func (m *ClassTeacherModel) GetTeacherByClassId(ClassId int64) (*ClassTeachers, error) {
	sql, args, _ := builder.Postgres().
		Select(classTeacherFieldNames).
		From(m.Table).
		Where(builder.Eq{"class_id": ClassId, "deleted": false}).
		ToSQL()
	var res ClassTeachers
	err := m.Conn.QueryRows(&res, sql, args...)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func (m *ClassTeacherModel) GetTeacherIdByClassId(ClassId int64) (*ClassTeachers, error) {
	sql, args, _ := builder.Postgres().
		Select("id").
		From(m.Table).
		Where(builder.Eq{"class_id": ClassId, "deleted": false}).
		ToSQL()
	var res ClassTeachers
	err := m.Conn.QueryRows(&res, sql, args...)
	if err != nil {
		return nil, err
	}
	return &res, nil
}
