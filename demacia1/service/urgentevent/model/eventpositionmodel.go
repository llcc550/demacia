package model

import (
	"demacia/common/cachemodel"
	"gitlab.u-jy.cn/xiaoyang/go-zero/core/stores/redis"
	"strings"
	"xorm.io/builder"

	"gitlab.u-jy.cn/xiaoyang/go-zero/core/stores/sqlx"
)

var (
	eventPositionFieldNames = strings.Join(cachemodel.RawFieldNames(&EventPosition{}, true), ",")
)

type (
	EventPositionModel struct {
		*cachemodel.CachedModel
	}
	EventPosition struct {
		Id            int64  `db:"id"`
		EventId       int64  `db:"event_id"`
		PositionId    int64  `db:"position_id"`
		PositionTitle string `db:"position_title"`
	}
	EventPositions []*EventPosition
)

func NewEventPositionModel(conn sqlx.SqlConn, cache *redis.Redis) *EventPositionModel {
	return &EventPositionModel{
		CachedModel: cachemodel.NewCachedModel(conn, `"urgentevent"."event_position"`, cache),
	}
}

func (m *EventPositionModel) Insert(data *EventPosition) error {
	query, args, _ := builder.Postgres().Insert(builder.Eq{
		"event_id":       data.EventId,
		"position_id":    data.PositionId,
		"position_title": data.PositionTitle,
	}).Into(m.Table).ToSQL()
	_, err := m.ExecNoCache(query, args...)
	return err
}

func (m *EventPositionModel) Delete(eventId int64) error {
	query, args, _ := builder.Postgres().Update(builder.Eq{"deleted": true}).From(m.Table).Where(builder.Eq{"event_id": eventId, "deleted": false}).ToSQL()
	_, err := m.Conn.Exec(query, args...)
	return err
}
func (m *EventPositionModel) FindListByEventIds(eventIds []int64) (EventPositions, error) {
	var res EventPositions
	query, args, _ := builder.Postgres().Select(eventPositionFieldNames).From(m.Table).Where(builder.In("event_id", eventIds)).Where(builder.Eq{"deleted": false}).ToSQL()
	err := m.Conn.QueryRows(&res, query, args...)
	if err != nil {
		return nil, err
	}
	return res, nil
}
func (m *EventPositionModel) FindListByPositionIds(positionIds []int64) (EventPositions, error) {
	var res EventPositions
	query, args, _ := builder.Postgres().Select(eventPositionFieldNames).From(m.Table).Where(builder.Eq{"deleted": false}).And(builder.In("position_id", positionIds)).ToSQL()
	err := m.Conn.QueryRows(&res, query, args...)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (m *EventPositionModel) FindListByOrgId(orgId int64) (EventPositions, error) {
	var res EventPositions
	query, args, _ := builder.Postgres().Select(eventPositionFieldNames).From(m.Table).Where(builder.Eq{"org_id": orgId, "deleted": false}).ToSQL()
	err := m.Conn.QueryRows(&res, query, args...)
	if err != nil {
		return nil, err
	}
	return res, nil
}
