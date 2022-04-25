package model

import (
	"database/sql"
	"fmt"
	"strings"
	"xorm.io/builder"

	"demacia/common/cachemodel"

	"gitlab.u-jy.cn/xiaoyang/go-zero/core/stores/redis"
	"gitlab.u-jy.cn/xiaoyang/go-zero/core/stores/sqlc"
	"gitlab.u-jy.cn/xiaoyang/go-zero/core/stores/sqlx"
)

var (
	deviceGroupFieldNames          = strings.Join(cachemodel.RawFieldNames(&DeviceGroup{}, true), ",")
	cacheDeviceDeviceGroupIdPrefix = "cache:device:deviceGroup:id:"
	cacheDeviceDeviceIdPrefix      = "cache:device:deviceGroup:device-id:"
	cacheDeviceGroupIdPrefix       = "cache:device:deviceGroup:group-id:"
)

type (
	DeviceGroupModel struct {
		*cachemodel.CachedModel
	}

	DeviceGroup struct {
		Id       int64 `db:"id"`
		DeviceId int64 `db:"device_id"`
		GroupId  int64 `db:"group_id"`
	}
	DeviceGroups []*DeviceGroup
)

func NewDeviceGroupModel(conn sqlx.SqlConn, cache *redis.Redis) *DeviceGroupModel {
	return &DeviceGroupModel{
		CachedModel: cachemodel.NewCachedModel(conn, `"device"."device_group"`, cache),
	}
}

func (m *DeviceGroupModel) Insert(data *DeviceGroup) (sql.Result, error) {
	eq := builder.Eq{
		"device_id": data.DeviceId,
		"group_id":  data.GroupId,
	}
	query, args, _ := builder.Postgres().Insert(eq).Into(m.Table).ToSQL()
	ret, err := m.ExecNoCache(query, args...)
	_ = m.DelCache(m.keys(data)...)
	return ret, err
}

func (m *DeviceGroupModel) FindOne(id int64) (*DeviceGroup, error) {
	deviceDeviceGroupIdKey := fmt.Sprintf("%s%v", cacheDeviceDeviceGroupIdPrefix, id)
	var resp DeviceGroup
	err := m.QueryRow(&resp, deviceDeviceGroupIdKey, func(conn sqlx.SqlConn, v interface{}) error {
		query := fmt.Sprintf("select %s from %s where id = $1 limit 1", deviceGroupFieldNames, m.Table)
		return conn.QueryRow(v, query, id)
	})
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, cachemodel.ErrNotFound
	default:
		return nil, err
	}
}

func (m *DeviceGroupModel) Update(data *DeviceGroup) error {
	deviceDeviceGroupIdKey := fmt.Sprintf("%s%v", cacheDeviceDeviceGroupIdPrefix, data.Id)
	_, err := m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set %s where id = $1", m.Table, deviceGroupFieldNames)
		return conn.Exec(query, data.Id, data.DeviceId, data.GroupId)
	}, deviceDeviceGroupIdKey)
	return err
}

func (m *DeviceGroupModel) Delete(id int64) error {
	deviceDeviceGroupIdKey := fmt.Sprintf("%s%v", cacheDeviceDeviceGroupIdPrefix, id)
	_, err := m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("delete from %s where id = $1", m.Table)
		return conn.Exec(query, id)
	}, deviceDeviceGroupIdKey)
	return err
}

func (m *DeviceGroupModel) formatPrimary(primary interface{}) string {
	return fmt.Sprintf("%s%v", cacheDeviceDeviceGroupIdPrefix, primary)
}

func (m *DeviceGroupModel) queryPrimary(conn sqlx.SqlConn, v, primary interface{}) error {
	query := fmt.Sprintf("select %s from %s where id = $1 limit 1", deviceGroupFieldNames, m.Table)
	return conn.QueryRow(v, query, primary)
}

func (m *DeviceGroupModel) FindListByGroupId(groupId []int64) (*DeviceGroups, error) {
	var resp DeviceGroups
	query, args, _ := builder.Postgres().Select(deviceGroupFieldNames).From(m.Table).Where(builder.In("group_id", groupId)).And(builder.Eq{"deleted": false}).ToSQL()
	err := m.Conn.QueryRows(&resp, query, args...)
	if err != nil {
		return nil, err
	}
	return &resp, err
}

func (m *DeviceGroupModel) FindListByDeviceId(deviceId []int64) (*DeviceGroups, error) {
	var resp DeviceGroups
	query, args, _ := builder.Postgres().Select(deviceGroupFieldNames).From(m.Table).Where(builder.In("device_id", deviceId)).And(builder.Eq{"deleted": false}).ToSQL()
	err := m.Conn.QueryRows(&resp, query, args...)
	if err != nil {
		return nil, err
	}
	return &resp, err
}

func (m *DeviceGroupModel) DeleteByGroupId(groupId int64) error {
	deviceDeviceGroupIdKey := fmt.Sprintf("%s%v", cacheDeviceGroupIdPrefix, groupId)
	_, err := m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query, args, _ := builder.Postgres().Update(builder.Eq{"deleted": true}).From(m.Table).Where(builder.Eq{"group_id": groupId}).ToSQL()
		return conn.Exec(query, args...)
	}, deviceDeviceGroupIdKey)
	return err
}

func (m *DeviceGroupModel) DeleteByDeviceId(deviceId int64) error {
	deviceDeviceIdKey := fmt.Sprintf("%s%v", cacheDeviceDeviceIdPrefix, deviceId)
	_, err := m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query, args, _ := builder.Postgres().Update(builder.Eq{"deleted": true}).From(m.Table).Where(builder.Eq{"device_id": deviceId}).ToSQL()
		return conn.Exec(query, args...)
	}, deviceDeviceIdKey)
	return err
}

func (m *DeviceGroupModel) keys(data *DeviceGroup) []string {
	res := make([]string, 0, 3)
	if data.DeviceId != 0 {
		res = append(res, fmt.Sprintf("%s%d", cacheDeviceDeviceIdPrefix, data.DeviceId))
	}
	if data.GroupId != 0 {
		res = append(res, fmt.Sprintf("%s%d", cacheDeviceGroupIdPrefix, data.GroupId))
	}
	if data.Id != 0 {
		res = append(res, m.formatPrimary(data.Id))
	}
	return res
}
