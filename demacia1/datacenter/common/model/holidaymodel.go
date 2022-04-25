package model

import (
	"demacia/common/cachemodel"
	"gitlab.u-jy.cn/xiaoyang/go-zero/core/stores/redis"
	"gitlab.u-jy.cn/xiaoyang/go-zero/core/stores/sqlx"
	"strings"
	"xorm.io/builder"
)

var (
	holidayFieldNames = strings.Join(cachemodel.RawFieldNames(&Holiday{}, true), ",")
	cacheHoliday      = "cache:common:holiday:all"
)

type (
	HolidayModel struct {
		*cachemodel.CachedModel
	}
	Holiday struct {
		Id          int64  `db:"id"`
		SpecialDate string `db:"special_date"`
		Year        int64  `db:"year"`
	}
	Holidays []*Holiday
)

func NewHolidayModel(conn sqlx.SqlConn, cache *redis.Redis) *HolidayModel {
	return &HolidayModel{
		CachedModel: cachemodel.NewCachedModel(conn, `"common"."holiday"`, cache),
	}
}

func (m *HolidayModel) SelectList(year int64) (Holidays, error) {
	var holidays Holidays
	key := cacheHoliday
	sql, args, _ := builder.Postgres().Select(holidayFieldNames).From(m.Table).Where(builder.Eq{"year": year, "deleted": false}).ToSQL()

	err := m.QueryRow(&holidays, key, func(conn sqlx.SqlConn, v interface{}) error {
		return conn.QueryRows(v, sql, args...)
	})

	if err != nil {
		return nil, cachemodel.ErrNotFound
	}

	return holidays, err
}
