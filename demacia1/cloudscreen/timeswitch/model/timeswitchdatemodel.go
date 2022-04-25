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
	timeSwitchDateFieldNames             = strings.Join(cachemodel.RawFieldNames(&TimeSwitchDate{}, true), ",")
	cacheTimeSwitchDateListByOrgIdPrefix = "cache:timeSwitchConfig:courseRecord:org-id:"
)

type (
	TimeSwitchDateModel struct {
		*cachemodel.CachedModel
	}
	TimeSwitchDate struct {
		Id        int64     `db:"id"`
		DeviceId  int64     `db:"device_id"`
		StartDate time.Time `db:"start_date"`
		EndDate   time.Time `db:"end_date"`
	}
	TimeSwitchDates []*TimeSwitchDate
)

func NewTimeSwitchDateModel(conn sqlx.SqlConn, cache *redis.Redis) *TimeSwitchDateModel {
	return &TimeSwitchDateModel{
		CachedModel: cachemodel.NewCachedModel(conn, `"timeswitch"."time_switch_date"`, cache),
	}
}

func (m *TimeSwitchDateModel) SelectByDeviceId(deviceId int64) (TimeSwitchDates, error) {

	var timeSwitchDates TimeSwitchDates

	sql, args, _ := builder.Postgres().Select(timeSwitchDateFieldNames).From(m.Table).Where(builder.Eq{"device_id": deviceId}).ToSQL()

	err := m.Conn.QueryRows(&timeSwitchDates, sql, args...)

	return timeSwitchDates, err
}

func (m *TimeSwitchDateModel) InsertTimeSwitchDates(deviceIds []int64, timeSwitchDates TimeSwitchDates) error {
	return m.Transact(func(session sqlx.Session) error {
		sql, args, _ := builder.Postgres().Delete(builder.In("device_id", deviceIds)).From(m.Table).ToSQL()
		_, err := session.Exec(sql, args...)
		if err != nil {
			return err
		}
		if len(timeSwitchDates) == 0 {
			return nil
		}
		sql, args, err = basefunc.BatchInsertString(m.Table, timeSwitchDates)
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

func (m *TimeSwitchDateModel) DeleteByDeviceId(deviceId int64) error {

	sql, args, _ := builder.Postgres().Delete(builder.Eq{"device_id": deviceId}).ToSQL()

	_, err := m.Conn.Exec(sql, args...)

	return err
}
