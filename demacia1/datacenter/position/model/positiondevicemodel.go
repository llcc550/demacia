package model

import (
	"demacia/common/basefunc"
	"demacia/common/cachemodel"
	"github.com/go-xorm/builder"
	"gitlab.u-jy.cn/xiaoyang/go-zero/core/stores/redis"
	"gitlab.u-jy.cn/xiaoyang/go-zero/core/stores/sqlx"
	"strings"
)

var (
	positionDeviceFieldNames    = strings.Join(cachemodel.RawFieldNames(&PositionDevice{}, true), ",")
	cachePositionDeviceIdPrefix = "cache:organization:organization:id:"
)

type (
	PositionDeviceModel struct {
		*cachemodel.CachedModel
	}

	PositionDevice struct {
		Id         int64  `db:"id"`
		PositionId int64  `db:"position_id"`
		DeviceId   int64  `db:"device_id"`
		DeviceName string `db:"device_name"`
	}
	PositionDevices []*PositionDevice
)

func NewPositionDeviceModel(conn sqlx.SqlConn, cache *redis.Redis) *PositionDeviceModel {
	return &PositionDeviceModel{
		CachedModel: cachemodel.NewCachedModel(conn, `"position"."position_device"`, cache),
	}
}

func (m *PositionDeviceModel) SelectByPidList(ids []int64) (PositionDevices, error) {

	var positionDevices PositionDevices

	sql, args, _ := builder.Postgres().Select(positionDeviceFieldNames).From(m.Table).Where(builder.In("position_id", ids)).ToSQL()

	err := m.Conn.QueryRows(&positionDevices, sql, args...)

	return positionDevices, err
}

func (m *PositionDeviceModel) InsertPositionDevice(ps []*PositionDevice) error {
	if len(ps) == 0 {
		return nil
	}
	query, args, err := basefunc.BatchInsertString(m.Table, ps)

	if err != nil {
		return err
	}
	_, err = m.Conn.Exec(query, args...)

	return err
}

func (m *PositionDeviceModel) UpdatePositionDevice(ps []*PositionDevice, pid int64) error {
	return m.Conn.Transact(func(session sqlx.Session) error {
		if len(ps) == 0 {
			return nil
		}

		sqlD, args, _ := builder.Postgres().Delete().From(m.Table).Where(builder.Eq{"position_id": pid}).ToSQL()
		_, err := session.Exec(sqlD, args...)
		if err != nil {
			return err
		}

		query, args, err := basefunc.BatchInsertString(m.Table, ps)
		if err != nil {
			return err
		}
		_, err = session.Exec(query, args...)
		return err
	})
}

func (m *PositionDeviceModel) SelectByDeviceIds(ids []int64) (PositionDevices, error) {

	var positionDevices PositionDevices

	sql, args, _ := builder.Postgres().Select(positionDeviceFieldNames).From(m.Table).Where(builder.In("device_id", ids)).ToSQL()

	err := m.Conn.QueryRows(&positionDevices, sql, args...)

	return positionDevices, err
}

func (m *PositionDeviceModel) DeleteByPid(pid int64) error {

	sql, args, _ := builder.Postgres().Delete(builder.Eq{"position_id": pid}).From(m.Table).ToSQL()

	_, err := m.Conn.Exec(sql, args...)

	return err
}
