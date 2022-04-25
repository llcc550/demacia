package model

import (
	"demacia/common/cachemodel"
	"gitlab.u-jy.cn/xiaoyang/go-zero/core/stores/redis"
	"gitlab.u-jy.cn/xiaoyang/go-zero/core/stores/sqlx"
	"strings"
	"xorm.io/builder"
)

var (
	courseRecordConfigFieldNames = strings.Join(cachemodel.RawFieldNames(&CourseRecordConfig{}, true), ",")
)

type (
	CourseRecordConfigModel struct {
		*cachemodel.CachedModel
	}
	CourseRecordConfig struct {
		Id             int64 `db:"id"`
		OrgId          int64 `db:"org_id"`
		Enable         bool  `db:"enable"`
		SignPerson     int8  `db:"sign_person"`
		SignBeforeTime int8  `db:"sign_before_time"`
		SignHoliday    bool  `db:"sign_holiday"`
	}
	CourseRecordConfigs []*CourseRecordConfig
)

func NewCourseRecordConfigModel(conn sqlx.SqlConn, cache *redis.Redis) *CourseRecordConfigModel {
	return &CourseRecordConfigModel{
		CachedModel: cachemodel.NewCachedModel(conn, `"courserecord"."course_record_config"`, cache),
	}
}

func (m *CourseRecordConfigModel) SelectList() (CourseRecordConfigs, error) {

	var configs CourseRecordConfigs

	sql, args, _ := builder.Postgres().Select(courseRecordConfigFieldNames).From(m.Table).Where(builder.Eq{"enable": true}).ToSQL()

	err := m.Conn.QueryRows(&configs, sql, args...)

	return configs, err

}

func (m *CourseRecordConfigModel) SelectByOrgId(oid int64) (CourseRecordConfig, error) {

	var courseRecordConfig CourseRecordConfig

	sql, args, _ := builder.Postgres().Select(courseRecordConfigFieldNames).From(m.Table).Where(builder.Eq{"org_id": oid, "deleted": false}).ToSQL()

	err := m.Conn.QueryRow(&courseRecordConfig, sql, args...)

	return courseRecordConfig, err
}

func (m *CourseRecordConfigModel) SelectDefaultConfig() (CourseRecordConfig, error) {
	var courseRecordConfig CourseRecordConfig

	sql, args, _ := builder.Postgres().Select(courseRecordConfigFieldNames).From(m.Table).Where(builder.Eq{"org_id": 0, "deleted": false}).ToSQL()

	err := m.Conn.QueryRow(&courseRecordConfig, sql, args...)

	return courseRecordConfig, err
}

func (m *CourseRecordConfigModel) InsertCourseRecordConfig(config *CourseRecordConfig) error {

	sql, args, _ := builder.Postgres().Insert(builder.Eq{
		"org_id":           config.OrgId,
		"enable":           config.Enable,
		"sign_person":      config.SignPerson,
		"sign_before_time": config.SignBeforeTime,
		"sign_holiday":     config.SignHoliday,
	}).Into(m.Table).ToSQL()

	_, err := m.Conn.Exec(sql, args...)

	return err
}

func (m *CourseRecordConfigModel) UpdateCourseRecordConfig(config *CourseRecordConfig) error {
	sql, args, _ := builder.Postgres().Update(builder.Eq{
		"org_id":           config.OrgId,
		"enable":           config.Enable,
		"sign_person":      config.SignPerson,
		"sign_before_time": config.SignBeforeTime,
		"sign_holiday":     config.SignHoliday,
	}).From(m.Table).Where(builder.Eq{"id": config.Id}).ToSQL()

	_, err := m.Conn.Exec(sql, args...)

	return err
}
