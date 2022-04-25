package model

import (
	"demacia/common/basefunc"
	"demacia/common/cachemodel"
	"gitlab.u-jy.cn/xiaoyang/go-zero/core/stores/redis"
	"gitlab.u-jy.cn/xiaoyang/go-zero/core/stores/sqlx"
	"strings"
	"time"
	"xorm.io/builder"
)

var (
	courseRecordDateFieldNames = strings.Join(cachemodel.RawFieldNames(&CourseRecordDate{}, true), ",")
)

type (
	CourseRecordDateModel struct {
		*cachemodel.CachedModel
	}

	CourseRecordDate struct {
		Id          int64  `db:"id"`
		OrgId       int64  `db:"org_id"`
		SpecialDate string `db:"special_date"`
		Type        int8   `db:"type"`
		Year        int64  `db:"year"`
		IsHoliday   int8   `db:"is_holiday"`
	}

	CourseRecordDates []*CourseRecordDate
)

func NewCourseRecordDateModel(conn sqlx.SqlConn, cache *redis.Redis) *CourseRecordDateModel {
	return &CourseRecordDateModel{
		CachedModel: cachemodel.NewCachedModel(conn, `"courserecord"."course_record_date"`, cache),
	}
}

func (m *CourseRecordDateModel) SelectHolidayList() (CourseRecordDates, error) {
	var dates CourseRecordDates
	sql, args, _ := builder.Postgres().Select(courseRecordDateFieldNames).From(m.Table).Where(builder.Eq{"org_id": 0, "deleted": false}).ToSQL()

	err := m.Conn.QueryRows(&dates, sql, args...)

	return dates, err
}

func (m *CourseRecordDateModel) SelectTodayIsHoliday() (int, error) {

	var count int

	sql, args, _ := builder.Postgres().Select("count(1)").From(m.Table).Where(builder.Eq{"year": time.Now().Year(), "special_date": time.Now().Format("2006-01-02")}).ToSQL()

	err := m.Conn.QueryRow(&count, sql, args...)

	return count, err
}

func (m *CourseRecordDateModel) SelectTodayIsMust(orgId int64) (*CourseRecordDate, error) {

	var courseRecordDate CourseRecordDate

	sql, args, _ := builder.Postgres().Select(courseRecordDateFieldNames).From(m.Table).Where(builder.Eq{"year": time.Now().Year(), "special_date": time.Now().Format("2006-01-02"), "org_id": orgId}).ToSQL()

	err := m.Conn.QueryRow(&courseRecordDate, sql, args...)

	return &courseRecordDate, err
}

func (m *CourseRecordDateModel) SelectByOrgId(orgId int64) (CourseRecordDates, error) {

	var courseRecordDates CourseRecordDates

	sql, args, _ := builder.Postgres().Select(courseRecordDateFieldNames).From(m.Table).Where(builder.Eq{"org_id": orgId}).ToSQL()

	err := m.Conn.QueryRows(&courseRecordDates, sql, args...)

	return courseRecordDates, err
}

func (m *CourseRecordDateModel) InsertDates(oid int64, dates CourseRecordDates) error {
	return m.Conn.Transact(func(session sqlx.Session) error {
		sql, args, _ := builder.Postgres().Delete(builder.Eq{"org_id": oid}).From(m.Table).ToSQL()
		_, err := m.Conn.Exec(sql, args...)
		if err != nil {
			return err
		}
		sqlI, args, err := basefunc.BatchInsertString(m.Table, dates)
		if err != nil {
			return err
		}
		_, err = m.Conn.Exec(sqlI, args...)
		if err != nil {
			return err
		}
		return nil
	})
}
