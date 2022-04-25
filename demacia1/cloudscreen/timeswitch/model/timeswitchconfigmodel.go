package model

import (
	"demacia/common/basefunc"
	"demacia/common/cachemodel"
	"gitlab.u-jy.cn/xiaoyang/go-zero/core/stores/redis"
	"gitlab.u-jy.cn/xiaoyang/go-zero/core/stores/sqlx"
	"strings"
	"xorm.io/builder"
)

var (
	timeSwitchConfigFieldNames             = strings.Join(cachemodel.RawFieldNames(&TimeSwitchConfig{}, true), ",")
	cacheTimeSwitchConfigListByOrgIdPrefix = "cache:timeSwitchConfig:courseRecord:org-id:"
)

type (
	TimeSwitchConfigModel struct {
		*cachemodel.CachedModel
	}
	TimeSwitchConfig struct {
		Id          int64 `db:"id"`
		OrgId       int64 `db:"org_id"`
		DeviceId    int64 `db:"device_id"`
		HolidayFlag int8  `db:"holiday_flag"`
	}
	TimeSwitchConfigs []*TimeSwitchConfig
)

func NewTimeSwitchConfigModel(conn sqlx.SqlConn, cache *redis.Redis) *TimeSwitchConfigModel {
	return &TimeSwitchConfigModel{
		CachedModel: cachemodel.NewCachedModel(conn, `"timeswitch"."time_switch_config"`, cache),
	}
}

func (m *TimeSwitchConfigModel) SelectByDeviceId(oid int64) (*TimeSwitchConfig, error) {
	var timeSwitchConfig TimeSwitchConfig
	sql, args, _ := builder.Postgres().Select(timeSwitchConfigFieldNames).From(m.Table).Where(builder.Eq{"device_id": oid}).Limit(1).ToSQL()
	err := m.Conn.QueryRow(&timeSwitchConfig, sql, args...)
	return &timeSwitchConfig, err
}

func (m *TimeSwitchConfigModel) InsertTimeSwitchConfigs(deviceIds []int64, timeSwitchConfigs TimeSwitchConfigs) error {
	return m.Transact(func(session sqlx.Session) error {
		sql, args, _ := builder.Postgres().Delete(builder.In("device_id", deviceIds)).From(m.Table).ToSQL()
		_, err := session.Exec(sql, args...)
		if err != nil {
			return err
		}
		if len(timeSwitchConfigs) == 0 {
			return nil
		}
		sql, args, err = basefunc.BatchInsertString(m.Table, timeSwitchConfigs)
		if err != nil {
			return err
		}
		_, err = session.Exec(sql, args...)
		if err != nil {
			return err
		}
		return nil
	})
}

func (m *TimeSwitchConfigModel) DeleteByDeviceId(deviceId int64) error {

	sql, args, _ := builder.Postgres().Delete(builder.Eq{"device_id": deviceId}).ToSQL()

	_, err := m.Conn.Exec(sql, args...)

	return err
}
