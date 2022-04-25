package model

import (
	"demacia/common/cachemodel"
	"fmt"
	"gitlab.u-jy.cn/xiaoyang/go-zero/core/stores/redis"
	"gitlab.u-jy.cn/xiaoyang/go-zero/core/stores/sqlx"
	"strings"
	"xorm.io/builder"
)

var (
	courseRecordCountFieldNames = strings.Join(cachemodel.RawFieldNames(&CourseRecordCount{}, true), ",")
)

type (
	CourseRecordCountModel struct {
		*cachemodel.CachedModel
	}
	CourseRecordCount struct {
		Id          int64  `db:"id"`
		UserId      int64  `db:"user_id"`
		Truename    string `db:"truename"`
		UserType    int8   `db:"user_type"`
		ShouldCount int64  `db:"should_count"`
		NormalCount int64  `db:"normal_count"`
		LateCount   int64  `db:"late_count"`
		CountDate   string `db:"count_date"`
		ClassId     int64  `db:"class_id"`
	}
	CourseRecordCounts []*CourseRecordCount

	StudentCourseRecordCount struct {
		ClassId     int64  `db:"class_id"`
		Truename    string `db:"truename"`
		ShouldCount int    `db:"should_count"`
		NormalCount int    `db:"normal_count"`
		LateCount   int    `db:"late_count"`
	}
	StudentCourseRecordCounts []*StudentCourseRecordCount
)

func NewCourseRecordCountModel(conn sqlx.SqlConn, cache *redis.Redis) *CourseRecordCountModel {
	return &CourseRecordCountModel{
		CachedModel: cachemodel.NewCachedModel(conn, `"courserecord"."course_record_count"`, cache),
	}
}

func (m *CourseRecordCountModel) InsertCourseRecordCount(courseRecordCount *CourseRecordCount) error {

	sql, args, _ := builder.Postgres().Insert(builder.Eq{"user_id": courseRecordCount.UserId, "truename": courseRecordCount.Truename, "user_type": courseRecordCount.UserType, "should_count": courseRecordCount.ShouldCount, "count_date": courseRecordCount.CountDate, "class_id": courseRecordCount.ClassId}).Into(m.Table).ToSQL()

	_, err := m.Conn.Exec(sql, args...)

	return err
}

func (m *CourseRecordCountModel) SelectByUserIdToToday(userId int64) (CourseRecordCount, error) {

	var recordCount CourseRecordCount

	sql, args, _ := builder.Postgres().Select(courseRecordCountFieldNames).From(m.Table).Where(builder.Eq{"user_id": userId}).ToSQL()

	err := m.Conn.QueryRow(&recordCount, sql, args...)

	return recordCount, err
}

func (m *CourseRecordCountModel) UpdateRecordCount(count *CourseRecordCount) error {

	sql, args, _ := builder.Postgres().Update(builder.Eq{"normal_count": count.NormalCount, "late_count": count.LateCount}).From(m.Table).Where(builder.Eq{"id": count.Id}).ToSQL()

	_, err := m.Conn.Exec(sql, args...)

	return err
}

func (m *CourseRecordCountModel) SelectByClassIds(classIds []int64, startDate, endDate, studentName string, page, limit int) (StudentCourseRecordCounts, int, error) {

	var count int

	var courseRecordCounts StudentCourseRecordCounts

	and := builder.NewCond()
	andC := builder.NewCond()

	column := "class_id,truename,sum(should_count) should_count,sum(normal_count) normal_count,sum(late_count) late_count"

	from := " (select class_id,truename from \"courserecord\".\"course_record_count\" where class_id IN (1,14) "
	if startDate != "" && endDate != "" {
		and = and.And(builder.Between{
			Col:     "count_date",
			LessVal: startDate,
			MoreVal: endDate,
		})
		from += fmt.Sprintf(" and count_date between '%s' and '%s' ", startDate, endDate)
	}
	if studentName != "" {
		and = and.And(builder.Like{"truename", studentName})
		andC = andC.And(builder.Like{"truename", studentName})
	}
	from += " group by class_id,truename) count "
	sql, args, _ := builder.Postgres().Select(column).From(m.Table).Where(builder.In("class_id", classIds)).And(and).GroupBy("class_id,truename").Limit(limit, (page-1)*limit).ToSQL()
	err := m.Conn.QueryRows(&courseRecordCounts, sql, args...)
	if err != nil {
		return nil, 0, err
	}
	sqlC, args, _ := builder.Postgres().Select("count(1)").And(andC).From(from).ToSQL()
	err = m.Conn.QueryRow(&count, sqlC, args...)

	return courseRecordCounts, count, err
}
