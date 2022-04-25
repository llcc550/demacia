package model

import (
	"demacia/common/basefunc"
	"demacia/common/cachemodel"
	"demacia/common/errlist"
	"fmt"
	"gitlab.u-jy.cn/xiaoyang/go-zero/core/stores/redis"
	"gitlab.u-jy.cn/xiaoyang/go-zero/core/stores/sqlx"
	"strings"
	"xorm.io/builder"
)

var (
	courseTableDeployNames                  = strings.Join(cachemodel.RawFieldNames(&CourseTableDeploy{}, true), ",")
	cacheCourseTableDeployListByOrgIdPrefix = "cache:courseTable:courseTableDeploy:org-id:"
	cacheCourseTableDeployIdPrefix          = "cache:courseTable:courseTable:id:"
)

type (
	CourseTableDeployModel struct {
		*cachemodel.CachedModel
	}

	CourseTableDeploy struct {
		Id         int64  `db:"id"`
		OrgId      int64  `db:"org_id"`
		CourseSort int8   `db:"course_sort"`
		Weekday    int8   `db:"week_day"`
		Note       string `db:"note"`
		Grouping   int8   `db:"grouping"`
		StartTime  string `db:"start_time"`
		EndTime    string `db:"end_time"`
		CourseFlag int8   `db:"course_flag"`
	}

	CourseTableDeploys []*CourseTableDeploy
)

func NewCourseTableDeployModel(conn sqlx.SqlConn, cache *redis.Redis) *CourseTableDeployModel {
	return &CourseTableDeployModel{
		CachedModel: cachemodel.NewCachedModel(conn, `"coursetable"."course_table_deploy"`, cache),
	}
}

func (m *CourseTableDeployModel) SelectByOrgId(oid int64) (CourseTableDeploys, error) {

	var courseTableDeploys CourseTableDeploys

	courseTableDeployListByOrgIdKey := fmt.Sprintf("%s%d", cacheCourseTableDeployListByOrgIdPrefix, oid)

	err := m.QueryRow(&courseTableDeploys, courseTableDeployListByOrgIdKey, func(conn sqlx.SqlConn, v interface{}) error {
		_sql, args, _ := builder.Postgres().Select(courseTableDeployNames).From(m.Table).Where(builder.Eq{"org_id": oid}).OrderBy("week_day,course_sort").ToSQL()
		return m.Conn.QueryRows(&courseTableDeploys, _sql, args...)
	})

	return courseTableDeploys, err
}

func (m *CourseTableDeployModel) InsertCourseTableDeploy(deploys []*CourseTableDeploy) error {

	if len(deploys) == 0 {
		return errlist.Unknown
	}
	_sql, args, err := basefunc.BatchInsertString(m.Table, deploys)
	if err != nil {
		return err
	}
	_, err = m.Conn.Exec(_sql, args...)
	_ = m.DelCache(m.keys(deploys[0])...)
	return err
}

func (m *CourseTableDeployModel) DeleteByOrgId(oid int64) error {

	sql, args, _ := builder.Postgres().Delete(builder.Eq{"org_id": oid}).From(m.Table).ToSQL()

	_, err := m.Conn.Exec(sql, args...)

	return err
}

func (m *CourseTableDeployModel) SelectByOrgIdAndWeekdayAndSort(oid int64, weekday, sort int8) (CourseTableDeploy, error) {

	var courseTableDeploy CourseTableDeploy

	sql, args, _ := builder.Postgres().Select(courseTableDeployNames).From(m.Table).Where(builder.Eq{"org_id": oid, "week_day": weekday, "course_sort": sort}).Limit(1).ToSQL()

	err := m.Conn.QueryRow(&courseTableDeploy, sql, args...)

	return courseTableDeploy, err

}

func (m *CourseTableDeployModel) keys(deploy *CourseTableDeploy) []string {

	res := make([]string, 0, 1)
	if deploy.OrgId != 0 {
		res = append(res, fmt.Sprintf("%s%d", cacheCourseTableDeployListByOrgIdPrefix, deploy.OrgId))
	}
	return res
}
