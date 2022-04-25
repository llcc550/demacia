package model

import (
	"demacia/common/basefunc"
	"demacia/common/cachemodel"
	"demacia/common/errlist"
	"gitlab.u-jy.cn/xiaoyang/go-zero/core/stores/redis"
	"gitlab.u-jy.cn/xiaoyang/go-zero/core/stores/sqlx"
	"strings"
	"time"
	"xorm.io/builder"
)

var (
	courseRecordFieldNames             = strings.Join(cachemodel.RawFieldNames(&CourseRecord{}, true), ",")
	cacheCourseRecordListByOrgIdPrefix = "cache:courseRecord:courseRecord:org-id:"
)

type (
	CourseRecordModel struct {
		*cachemodel.CachedModel
	}
	CourseRecord struct {
		Id           int64  `db:"id"`
		OrgId        int64  `db:"org_id"`
		UserId       int64  `db:"user_id"`
		Truename     string `db:"truename"`
		UserType     int8   `db:"user_type"`
		SignDate     string `db:"sign_date"`
		SignTime     string `db:"sign_time"`
		SubjectName  string `db:"subject_name"`
		CourseNote   string `db:"course_note"`
		Status       int8   `db:"status"`
		Photo        string `db:"photo"`
		ClassName    string `db:"class_name"`
		ClassId      int64  `db:"class_id"`
		PositionName string `db:"position_name"`
		StartTime    string `db:"start_time"`
		EndTime      string `db:"end_time"`
	}
	CourseRecords []*CourseRecord
)

func NewCourseRecordModel(conn sqlx.SqlConn, cache *redis.Redis) *CourseRecordModel {
	return &CourseRecordModel{
		CachedModel: cachemodel.NewCachedModel(conn, `"courserecord"."course_record"`, cache),
	}
}

func (m *CourseRecordModel) InsertCourseRecords(records CourseRecords) error {
	if len(records) == 0 {
		return errlist.Unknown
	}

	sql, args, err := basefunc.BatchInsertString(m.Table, records)
	if err != nil {
		return err
	}

	_, err = m.Conn.Exec(sql, args...)

	return err
}

func (m *CourseRecordModel) SelectByUserIdAndDate(userId int64, date string) (CourseRecords, error) {

	var records CourseRecords

	sql, args, _ := builder.Postgres().Select(courseRecordFieldNames).From(m.Table).Where(builder.Eq{"user_id": userId}).And(builder.Eq{"sign_date": date}).ToSQL()
	err := m.Conn.QueryRows(&records, sql, args...)

	return records, err
}

func (m *CourseRecordModel) UpdateCourseRecord(record *CourseRecord) error {

	sql, args, _ := builder.Postgres().Update(builder.Eq{"sign_time": record.SignTime, "status": record.Status}).Where(builder.Eq{"id": record.Id}).From(m.Table).ToSQL()

	_, err := m.Conn.Exec(sql, args...)

	return err
}

func (m *CourseRecordModel) SelectByParam(startTime, endTime, truename string, userType, state int8, page, limit int, oid int64) (CourseRecords, int, error) {

	var count int

	var courseRecords CourseRecords

	eq := builder.NewCond()

	if startTime != "" && endTime != "" {
		eq = eq.And(builder.Between{
			Col:     "sign_date",
			LessVal: startTime,
			MoreVal: endTime,
		})
	}

	if truename != "" {
		eq = eq.And(builder.Like{"truename", truename})
	}

	if state != -1 {
		eq = eq.And(builder.Eq{"status": state})
	}

	if userType != -1 {
		eq = eq.And(builder.Eq{"user_type": userType})
	}

	sql, args, err := builder.Postgres().Select(courseRecordFieldNames).From(m.Table).Where(builder.Eq{"org_id": oid}).And(eq).Limit(limit, (page-1)*limit).ToSQL()
	err = m.Conn.QueryRows(&courseRecords, sql, args...)
	if err != nil {
		return nil, 0, err
	}
	sqlC, argsC, _ := builder.Postgres().Select("count(1)").From(m.Table).Where(builder.Eq{"org_id": oid}).And(eq).ToSQL()
	err = m.Conn.QueryRow(&count, sqlC, argsC...)
	return courseRecords, count, err
}

func (m *CourseRecordModel) SelectBetweenStartTime(cid int64, startTime string) (CourseRecords, error) {

	var courseRecords CourseRecords

	sql, args, _ := builder.Postgres().Select(courseRecordFieldNames).From(m.Table).Where(builder.Eq{"class_id": cid}).And(builder.Eq{"sign_date": time.Now().Format("2006-01-02")}).And(builder.Expr("start_time <= ? and end_time > ?", startTime, startTime)).ToSQL()

	err := m.Conn.QueryRows(&courseRecords, sql, args...)

	return courseRecords, err
}

func (m *CourseRecordModel) SelectByClassIdAndParam(cid int64, truename, subjectName, queryDate string, status int8, page, limit int) (CourseRecords, int, error) {

	var courseRecords CourseRecords

	var count int

	and := builder.NewCond()
	if truename != "" {
		and = and.And(builder.Like{"truename", truename})
	}

	if subjectName != "" {
		and = and.And(builder.Eq{"subject_name": subjectName})
	}

	if queryDate != "" {
		and = and.And(builder.Eq{"sign_date": queryDate})
	}

	if status != -1 {
		and = and.And(builder.Eq{"status": status})
	}
	sql, args, _ := builder.Postgres().Select(courseRecordFieldNames).From(m.Table).Where(builder.Eq{"class_id": cid}).And(and).Limit(limit, (page-1)*limit).ToSQL()
	err := m.Conn.QueryRows(&courseRecords, sql, args...)
	if err != nil {
		return nil, 0, err
	}

	sqlC, argsC, _ := builder.Postgres().Select("count(1)").From(m.Table).Where(builder.Eq{"class_id": cid}).And(and).ToSQL()

	err = m.Conn.QueryRow(&count, sqlC, argsC...)

	return courseRecords, count, err
}

func (m *CourseRecordModel) SelectByStudentId(userId int64, userType, page, limit int) (CourseRecords, int, error) {
	var courseRecords CourseRecords
	var count int

	sql, args, _ := builder.Postgres().Select(courseRecordFieldNames).From(m.Table).Where(builder.Eq{"user_id": userId, "user_type": userType}).And(builder.Neq{"status": 0}).OrderBy("create_time desc").Limit(limit, (page-1)*limit).ToSQL()

	err := m.Conn.QueryRows(&courseRecords, sql, args...)
	if err != nil {
		return nil, 0, err
	}

	sqlC, argsC, err := builder.Postgres().Select("count(1)").From(m.Table).Where(builder.Eq{"user_id": userId, "user_type": userType}).And(builder.Neq{"status": 0}).ToSQL()

	err = m.Conn.QueryRow(&count, sqlC, argsC...)

	return courseRecords, count, err
}

func (m *CourseRecordModel) UpdatePhoto(record *CourseRecord) error {

	sql, args, _ := builder.Postgres().Update(builder.Eq{"photo": record.Photo}).From(m.Table).Where(builder.Eq{"id": record.Id}).ToSQL()

	_, err := m.Conn.Exec(sql, args...)

	return err
}
