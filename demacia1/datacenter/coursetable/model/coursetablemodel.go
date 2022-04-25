package model

import (
	"database/sql"
	"demacia/common/cachemodel"
	"fmt"
	"gitlab.u-jy.cn/xiaoyang/go-zero/core/stores/redis"
	"gitlab.u-jy.cn/xiaoyang/go-zero/core/stores/sqlx"
	"strings"
	"time"
	"xorm.io/builder"
)

var (
	courseTableFieldNames               = strings.Join(cachemodel.RawFieldNames(&CourseTable{}, true), ",")
	cacheCourseTableIdPrefix            = "cache:courseTable:courseTable:id:"
	cacheCourseTableListByOrgIdPrefix   = "cache:courseTable:courseTable:org-id:"
	cacheCourseTableListByClassIdPrefix = "cache:courseTable:courseTable:class-id:"
)

type (
	CourseTableModel struct {
		*cachemodel.CachedModel
	}

	CourseTable struct {
		Id             int64     `db:"id"`
		OrganizationId int64     `db:"organization_id"` // 学校id
		PositionId     int64     `db:"position_id"`     // 教室id
		PositionName   string    `db:"position_name"`   // 教室名称
		SubjectId      int64     `db:"subject_id"`      // 科目id
		SubjectName    string    `db:"subject_name"`    // 科目名称
		ClassId        int64     `db:"class_id"`        // 班级id
		ClassName      string    `db:"class_name"`      // 班级名称
		ClassType      int64     `db:"class_type"`      // 班级类型 1行政 2兴趣
		TeacherId      int64     `db:"teacher_id"`      // 教师id
		TeacherName    string    `db:"teacher_name"`    // 教师姓名
		WeekDay        int8      `db:"week_day"`        // 星期几的课程
		StartTime      string    `db:"start_time"`      // 课程开始时间
		EndTime        string    `db:"end_time"`        // 课程结束时间
		CreateTime     time.Time `db:"create_time"`
		UpdateTime     time.Time `db:"update_time"`
		CourseSort     int8      `db:"course_sort"`
	}

	CourseTables []*CourseTable
)

func NewCourseTableModel(conn sqlx.SqlConn, cache *redis.Redis) *CourseTableModel {
	return &CourseTableModel{
		CachedModel: cachemodel.NewCachedModel(conn, `"coursetable"."course_table"`, cache),
	}
}

func (m *CourseTableModel) SelectCourseTableByOrgId(oid int64) (CourseTables, error) {

	var courseTables CourseTables

	courseTableListByOrgIdKey := fmt.Sprintf("%s%d", cacheCourseTableListByOrgIdPrefix, oid)

	err := m.QueryRow(&courseTables, courseTableListByOrgIdKey, func(conn sqlx.SqlConn, v interface{}) error {
		sqlQ, args, _ := builder.Postgres().Select(courseTableFieldNames).From(m.Table).Where(builder.Eq{"organization_id": oid, "deleted": false}).OrderBy("week_day,course_sort").ToSQL()
		return m.Conn.QueryRows(&courseTables, sqlQ, args...)
	})

	if err != nil {
		return nil, err
	}
	return courseTables, err
}

func (m *CourseTableModel) SelectByCid(classId int64) (CourseTables, error) {

	var courseTables CourseTables

	courseTableListByClassIdKey := fmt.Sprintf("%s%d", cacheCourseTableListByClassIdPrefix, classId)

	err := m.QueryRow(&courseTables, courseTableListByClassIdKey, func(conn sqlx.SqlConn, v interface{}) error {
		sqlQ, args, _ := builder.Postgres().Select(courseTableFieldNames).From(m.Table).Where(builder.Eq{`"class_id"`: classId, "deleted": false}).OrderBy("week_day,course_sort").ToSQL()
		return m.Conn.QueryRows(&courseTables, sqlQ, args...)
	})

	if err != nil {
		return nil, err
	}

	return courseTables, err
}

func (m *CourseTableModel) SelectByPositionId(positionId int64) (CourseTables, error) {
	var courseTables CourseTables

	toSQL, args, _ := builder.Postgres().Select(courseTableFieldNames).From(m.Table).Where(builder.Eq{"position_id": positionId, "deleted": false}).ToSQL()

	err := m.Conn.QueryRows(&courseTables, toSQL, args...)

	return courseTables, err
}

func (m *CourseTableModel) InsertCourseTable(course *CourseTable) error {

	toSQL, args, _ := builder.Postgres().Insert(builder.Eq{"course_sort": course.CourseSort, "week_day": course.WeekDay, "position_name": course.PositionName, "position_id": course.PositionId, "subject_id": course.SubjectId, "subject_name": course.SubjectName, "class_id": course.ClassId, "class_name": course.ClassName, "teacher_id": course.TeacherId, "teacher_name": course.TeacherName, "class_type": 1, "organization_id": course.OrganizationId, "start_time": course.StartTime, "end_time": course.EndTime}).Into(m.Table).ToSQL()

	_, err := m.Conn.Exec(toSQL, args...)

	_ = m.DelCache(m.keys(course)...)

	return err
}

func (m *CourseTableModel) UpdateCourseTable(course *CourseTable) error {

	_, err := m.Exec(func(conn sqlx.SqlConn) (sql.Result, error) {

		toSQL, args, _ := builder.Postgres().Update(builder.Eq{"position_name": course.PositionName, "position_id": course.PositionId, "subject_id": course.SubjectId, "subject_name": course.SubjectName, "class_id": course.ClassId, "class_name": course.ClassName, "teacher_id": course.TeacherId, "teacher_name": course.TeacherName, "class_type": 1, "organization_id": course.OrganizationId}).From(m.Table).Where(builder.Eq{"id": course.Id, "deleted": false}).ToSQL()

		return conn.Exec(toSQL, args...)

	}, m.keys(course)...)

	return err
}

func (m *CourseTableModel) SelectExistByClassIdAndWeekdayAndSort(oid, cid int64, weekday, sort int8) (int64, error) {

	var id int64

	toSQL, args, _ := builder.Postgres().Select("id").From(m.Table).Where(builder.Eq{"organization_id": oid, "class_id": cid, "week_day": weekday, "course_sort": sort, "deleted": false}).Limit(1).ToSQL()

	err := m.Conn.QueryRow(&id, toSQL, args...)

	return id, err
}

func (m *CourseTableModel) SelectByClassIdAndWeekdayAndSort(oid, cid int64, weekday, sort int8) (CourseTable, error) {

	var course CourseTable

	toSQL, args, _ := builder.Postgres().Select(courseTableFieldNames).From(m.Table).Where(builder.Eq{"organization_id": oid, "class_id": cid, "week_day": weekday, "course_sort": sort, "deleted": false}).Limit(1).ToSQL()

	err := m.Conn.QueryRow(&course, toSQL, args...)

	return course, err
}

func (m *CourseTableModel) DeleteByClassId(cid int64) error {

	courseTableListByClassIdKey := fmt.Sprintf("%s%d", cacheCourseTableListByClassIdPrefix, cid)

	toSQL, args, _ := builder.Postgres().Update(builder.Eq{"deleted": true}).From(m.Table).Where(builder.Eq{"class_id": cid}).ToSQL()

	_, err := m.ExecNoCache(toSQL, args...)

	_ = m.DelCache(courseTableListByClassIdKey)
	return err
}

func (m *CourseTableModel) DeleteByOrgId(oid int64) error {

	courseTableListByClassIdKey := fmt.Sprintf("%s%d", cacheCourseTableListByOrgIdPrefix, oid)

	toSQL, args, _ := builder.Postgres().Update(builder.Eq{"deleted": true}).From(m.Table).Where(builder.Eq{"organization_id": oid}).ToSQL()

	_, err := m.ExecNoCache(toSQL, args...)

	_ = m.DelCache(courseTableListByClassIdKey)
	return err
}

func (m *CourseTableModel) SelectByMemberId(userId int64) (CourseTables, error) {
	var courseTables CourseTables

	toSQL, args, _ := builder.Postgres().Select(courseTableFieldNames).From(m.Table).Where(builder.Eq{"teacher_id": userId, "deleted": false}).ToSQL()

	err := m.Conn.QueryRows(&courseTables, toSQL, args...)

	return courseTables, err
}

func (m *CourseTableModel) SelectByClassIdAndSubjectId(classId, subjectId int64) (CourseTables, error) {
	var courseTables CourseTables

	toSQL, args, _ := builder.Postgres().Select(courseTableFieldNames).From(m.Table).Where(builder.Eq{"class_id": classId, "subject_id": subjectId, "deleted": false}).ToSQL()

	err := m.Conn.QueryRows(&courseTables, toSQL, args...)

	return courseTables, err
}

func (m *CourseTableModel) UnBindTeacherByClassIdAndMemberId(classId, memberId int64) error {

	courseTableListByClassIdKey := fmt.Sprintf("%s%d", cacheCourseTableListByClassIdPrefix, classId)

	toSQL, args, _ := builder.Postgres().Update(builder.Eq{"teacher_id": 0, "teacher_name": ""}).From(m.Table).Where(builder.Eq{"class_id": classId, "teacher_id": memberId, "deleted": false}).ToSQL()

	_, err := m.ExecNoCache(toSQL, args...)

	_ = m.DelCache(courseTableListByClassIdKey)
	return err
}

func (m *CourseTableModel) ChangeTeacherByClassId(classId, memberId, subjectId int64, truename string) error {

	courseTableListByClassIdKey := fmt.Sprintf("%s%d", cacheCourseTableListByClassIdPrefix, classId)

	toSQL, args, _ := builder.Postgres().Update(builder.Eq{"teacher_id": memberId, "teacher_name": truename}).From(m.Table).Where(builder.Eq{"class_id": classId, "subject_id": subjectId, "deleted": false}).ToSQL()

	_, err := m.ExecNoCache(toSQL, args...)

	_ = m.DelCache(courseTableListByClassIdKey)

	return err
}

func (m *CourseTableModel) SqlDB() (*sql.DB, error) {
	return m.Conn.RawDB()
}

func (m *CourseTableModel) keys(courseTable *CourseTable) []string {

	res := make([]string, 0, 2)
	if courseTable.OrganizationId != 0 {
		res = append(res, fmt.Sprintf("%s%d", cacheCourseTableListByOrgIdPrefix, courseTable.OrganizationId))
	}

	if courseTable.ClassId != 0 {
		res = append(res, fmt.Sprintf("%s%d", cacheCourseTableListByClassIdPrefix, courseTable.ClassId))
	}
	return res
}

func (m *CourseTableModel) UnbindPosition(positionId int64) error {

	toSQL, args, _ := builder.Postgres().Update(builder.Eq{"position_id": 0, "position_name": ""}).From(m.Table).Where(builder.Eq{"position_id": positionId}).ToSQL()

	_, err := m.Conn.Exec(toSQL, args...)

	return err
}

func (m *CourseTableModel) DeleteTeacherInfo(memberId int64) error {

	toSQL, args, _ := builder.Postgres().Update(builder.Eq{"teacher_id": 0, "teacher_name": ""}).From(m.Table).Where(builder.Eq{"teacher_id": memberId}).ToSQL()

	_, err := m.Conn.Exec(toSQL, args...)

	return err
}
