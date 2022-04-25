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
	timeSwitchFieldNames             = strings.Join(cachemodel.RawFieldNames(&TimeSwitch{}, true), ",")
	cacheTimeSwitchListByOrgIdPrefix = "cache:timeSwitchConfig:courseRecord:org-id:"
)

type (
	TimeSwitchModel struct {
		*cachemodel.CachedModel
	}
	TimeSwitch struct {
		Id        int64  `db:"id"`
		DeviceId  int64  `db:"device_id"`
		StartTime string `db:"start_time"`
		EndTime   string `db:"end_time"`
		Weekday   int8   `db:"weekday"`
	}
	TimeSwitches []*TimeSwitch
)

func NewTimeSwitchModel(conn sqlx.SqlConn, cache *redis.Redis) *TimeSwitchModel {
	return &TimeSwitchModel{
		CachedModel: cachemodel.NewCachedModel(conn, `"timeswitch"."time_switch"`, cache),
	}
}

func (m *TimeSwitchModel) SelectByDeviceId(deviceId int64) (TimeSwitches, error) {

	var timeSwitches TimeSwitches

	sql, args, _ := builder.Postgres().Select(timeSwitchFieldNames).From(m.Table).Where(builder.Eq{"device_id": deviceId}).ToSQL()

	err := m.Conn.QueryRows(&timeSwitches, sql, args...)

	return timeSwitches, err
}

func (m *TimeSwitchModel) InsertTimeSwitches(deviceIds []int64, timeSwitches TimeSwitches) error {
	return m.Transact(func(session sqlx.Session) error {
		sql, args, _ := builder.Postgres().Delete(builder.In("device_id", deviceIds)).From(m.Table).ToSQL()
		_, err := session.Exec(sql, args...)
		if err != nil {
			return err
		}
		if len(timeSwitches) == 0 {
			return nil
		}
		sql, args, err = basefunc.BatchInsertString(m.Table, timeSwitches)
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

func (m *TimeSwitchModel) DeleteByDeviceId(deviceId int64) error {

	sql, args, _ := builder.Postgres().Delete(builder.Eq{"device_id": deviceId}).ToSQL()

	_, err := m.Conn.Exec(sql, args...)

	return err
}
