package model

import (
	"database/sql"
	"fmt"
	"strings"

	"demacia/common/cachemodel"

	"gitlab.u-jy.cn/xiaoyang/go-zero/core/stores/redis"
	"gitlab.u-jy.cn/xiaoyang/go-zero/core/stores/sqlx"
	"xorm.io/builder"
)

var (
	deviceFieldNames    = strings.Join(cachemodel.RawFieldNames(&Device{}, true), ",")
	cacheDeviceIdPrefix = "cache:device:id:"
	cacheOrgIdPrefix    = "cache:device:org-id:"
)

type (
	DeviceModel struct {
		*cachemodel.CachedModel
	}
	Device struct {
		Id      int64  `db:"id"`
		Sn      string `db:"sn"`
		OrgId   int64  `db:"org_id"`
		Title   string `db:"title"`
		Network int8   `json:"network"`
	}
	Devices []*Device
)

func NewDeviceModel(conn sqlx.SqlConn, cache *redis.Redis) *DeviceModel {
	return &DeviceModel{
		CachedModel: cachemodel.NewCachedModel(conn, `"device"."device"`, cache),
	}
}

func (m *DeviceModel) GetDeviceInfoById(id int64) (*Device, error) {
	studentIdKey := fmt.Sprintf("%s%d", cacheDeviceIdPrefix, id)
	var resp Device
	err := m.QueryRow(&resp, studentIdKey, func(conn sqlx.SqlConn, v interface{}) error {
		query, args, _ := builder.Postgres().From(m.Table).Select(deviceFieldNames).Where(builder.Eq{"id": id, "deleted": false}).Limit(1).ToSQL()
		return conn.QueryRow(v, query, args...)
	})
	if err != nil {
		return nil, cachemodel.ErrNotFound
	}
	return &resp, nil
}

func (m *DeviceModel) InsertOne(data *Device) (int64, error) {
	eq := builder.Eq{
		"sn":     data.Sn,
		"org_id": data.OrgId,
		"title":  data.Title,
	}
	query, args, _ := builder.Postgres().Insert(eq).Into(m.Table).ToSQL()
	var res int64
	err := m.Conn.QueryRow(&res, query+" returning id", args...)
	if err != nil {
		return 0, err
	}
	return res, nil
}

func (m *DeviceModel) UpdateTitle(id int64, title string) error {
	key := fmt.Sprintf("%s%d", cacheDeviceIdPrefix, id)
	_, err := m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query, args, _ := builder.Postgres().Update(builder.Eq{"title": title}).From(m.Table).Where(builder.Eq{"id": id, "deleted": false}).ToSQL()
		return conn.Exec(query, args...)
	}, key)
	return err
}

func (m *DeviceModel) DeleteOneById(id int64) error {
	key := fmt.Sprintf("%s%d", cacheDeviceIdPrefix, id)
	_, err := m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query, args, _ := builder.Postgres().Update(builder.Eq{"deleted": true}).From(m.Table).Where(builder.Eq{"id": id}).ToSQL()
		return conn.Exec(query, args...)
	}, key)
	_ = m.DelCache(key)
	return err
}

func (m *DeviceModel) ListByConditions(ids []int64, orgId int64, title, sn string, network int8, page, limit int) (Devices, int, error) {
	var res Devices
	var count int
	eq := builder.Eq{"org_id": orgId, "deleted": false}.And()
	if len(ids) > 0 {
		eq = eq.And(builder.In("id", ids))
	}
	if title != "" {
		eq = eq.And(builder.Like{"title", title})
	}
	if sn != "" {
		eq = eq.And(builder.Like{"sn", sn})
	}
	if network != 0 {
		eq = eq.And(builder.Eq{"network": network})
	}
	query := builder.Postgres().Select(deviceFieldNames).From(m.Table).Where(eq)
	queryC, argsC, _ := builder.Postgres().Select("COUNT(*)").From(m.Table).Where(eq).ToSQL()
	if page != 0 && limit != 0 {
		query = query.Limit(limit, (page-1)*limit)
		err := m.Conn.QueryRow(&count, queryC, argsC...)
		if err != nil {
			return nil, 0, err
		}
	}
	queryQ, argsQ, _ := query.OrderBy("id").ToSQL()
	err := m.Conn.QueryRows(&res, queryQ, argsQ...)
	if err != nil {
		return nil, 0, err
	}
	return res, count, nil
}

func (m *DeviceModel) ListByOrgId(orgId int64) (Devices, error) {
	var res Devices
	key := fmt.Sprintf("%s%d", cacheOrgIdPrefix, orgId)
	eq := builder.Eq{"org_id": orgId, "deleted": false}.And()
	err := m.QueryRow(&res, key, func(conn sqlx.SqlConn, v interface{}) error {
		query, args, _ := builder.Postgres().Select(deviceFieldNames).From(m.Table).Where(eq).OrderBy("id ").ToSQL()
		return conn.QueryRows(v, query, args...)
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (m *DeviceModel) ListByIds(ids []int64) (Devices, error) {
	var res Devices
	eq := builder.Eq{"deleted": false}.And(builder.In("id", ids))
	query := builder.Postgres().Select(deviceFieldNames).From(m.Table).Where(eq)
	queryQ, argsQ, _ := query.OrderBy("id DESC").ToSQL()
	err := m.Conn.QueryRows(&res, queryQ, argsQ...)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (m *DeviceModel) GetOneByTitle(orgId int64, title string) (*Device, error) {
	var res Device
	eq := builder.Eq{"deleted": false}.And(builder.Eq{"title": title, "org_id": orgId})
	query, args, _ := builder.Postgres().Select(deviceFieldNames).From(m.Table).Where(eq).Limit(1).ToSQL()
	err := m.Conn.QueryRow(&res, query, args...)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func (m *DeviceModel) GetOneBySn(sn string) (*Device, error) {
	var res Device
	eq := builder.Eq{"deleted": false}.And(builder.Eq{"sn": sn})
	query, args, _ := builder.Postgres().Select(deviceFieldNames).From(m.Table).Where(eq).Limit(1).ToSQL()
	err := m.Conn.QueryRow(&res, query, args...)
	if err != nil {
		return nil, err
	}
	return &res, nil
}
