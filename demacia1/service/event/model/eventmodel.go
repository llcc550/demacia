package model

import (
	"database/sql"
	"demacia/common/cachemodel"
	"fmt"
	"gitlab.u-jy.cn/xiaoyang/go-zero/core/stores/redis"
	"gitlab.u-jy.cn/xiaoyang/go-zero/core/stores/sqlx"
	"strings"
	"xorm.io/builder"
)

var (
	eventFieldNames = strings.Join(cachemodel.RawFieldNames(&Event{}, true), ",")

	cacheEventEventIdPrefix  = "cache:event:event:id:"
	cacheEventOrgIdPrefix    = "cache:event:event:org-id:"
	cacheEventCateIdPrefix   = "cache:event:event:category-id:"
	cacheEventPushTypePrefix = "cache:event:event:push-type:"
)

type (
	EventModel struct {
		*cachemodel.CachedModel
	}

	Event struct {
		Id         int64  `db:"id"`
		OrgId      int64  `db:"org_id"`
		Name       string `db:"name"`
		CategoryId int64  `db:"category_id"`
		Content    string `db:"content"`
		PushType   int8   `db:"push_type"`
		StartTime  int64  `db:"start_time"`
		EndTime    int64  `db:"end_time"`
		MemberId   int64  `db:"member_id"`
		MemberName string `db:"member_name"`
		CreatedAt  int64  `db:"created_at"`
	}
)

func NewEventModel(conn sqlx.SqlConn, cache *redis.Redis) *EventModel {
	return &EventModel{
		CachedModel: cachemodel.NewCachedModel(conn, `"event"."event"`, cache),
	}
}

func (m *EventModel) Insert(data *Event) (int64, error) {
	eq := builder.Eq{
		"org_id":      data.OrgId,
		"name":        data.Name,
		"category_id": data.CategoryId,
		"content":     data.Content,
		"push_type":   data.PushType,
		"start_time":  data.StartTime,
		"end_time":    data.EndTime,
		"member_id":   data.MemberId,
		"member_name": data.MemberName,
		"created_at":  data.CreatedAt,
	}
	var eventId int64
	query, args, _ := builder.Postgres().Insert(eq).Into(m.Table).ToSQL()
	err := m.Conn.QueryRow(&eventId, query+" returning id", args...)
	if err != nil {
		return 0, err
	}
	_ = m.DelCache(m.keys(data)...)
	return eventId, err
}

func (m *EventModel) Update(data *Event) error {
	eq := builder.Eq{
		"org_id":      data.OrgId,
		"name":        data.Name,
		"category_id": data.CategoryId,
		"content":     data.Content,
		"push_type":   data.PushType,
		"start_time":  data.StartTime,
		"end_time":    data.EndTime,
	}
	_, err := m.Exec(func(conn sqlx.SqlConn) (sql.Result, error) {
		query, args, _ := builder.Postgres().Update(eq).From(m.Table).Where(builder.Eq{"id": data.Id, "deleted": false}).ToSQL()
		return conn.Exec(query, args...)
	}, m.keys(data)...)
	return err
}

func (m *EventModel) FindOne(id int64) (*Event, error) {
	eventEventIdKey := fmt.Sprintf("%s%v", cacheEventEventIdPrefix, id)
	var resp Event
	err := m.QueryRow(&resp, eventEventIdKey, func(conn sqlx.SqlConn, v interface{}) error {
		query, args, _ := builder.Postgres().Select(eventFieldNames).From(m.Table).Where(builder.Eq{"id": id, "deleted": false}).ToSQL()
		return conn.QueryRow(v, query, args...)
	})
	if err != nil {
		return nil, cachemodel.ErrNotFound
	}
	return &resp, nil
}

func (m *EventModel) Delete(id int64) error {
	eventEventIdKey := fmt.Sprintf("%s%v", cacheEventEventIdPrefix, id)
	_, err := m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query, args, _ := builder.Postgres().Update(builder.Eq{"deleted": true}).From(m.Table).Where(builder.Eq{"id": id, "deleted": false}).ToSQL()
		return conn.Exec(query, args...)
	}, eventEventIdKey)
	return err
}

func (m *EventModel) formatPrimary(primary interface{}) string {
	return fmt.Sprintf("%s%v", cacheEventEventIdPrefix, primary)
}

func (m *EventModel) queryPrimary(conn sqlx.SqlConn, v, primary interface{}) error {
	query := fmt.Sprintf("select %s from %s where id = $1 limit 1", eventFieldNames, m.Table)
	return conn.QueryRow(v, query, primary)
}

func (m *EventModel) keys(event *Event) []string {
	res := make([]string, 0, 4)
	if event.OrgId != 0 {
		res = append(res, fmt.Sprintf("%s%d", cacheEventOrgIdPrefix, event.OrgId))
	}
	if event.CategoryId != 0 {
		res = append(res, fmt.Sprintf("%s%d", cacheEventCateIdPrefix, event.CategoryId))
	}
	if event.PushType != 0 {
		res = append(res, fmt.Sprintf("%s%d", cacheEventPushTypePrefix, event.PushType))
	}
	if event.Id != 0 {
		res = append(res, m.formatPrimary(event.Id))
	}
	return res
}

func (m *EventModel) FindListByConditions(eventIds []int64, eventName, memberName string, orgId, startTime, endTime, categoryId int64, page, limit int) ([]*Event, int, error) {
	var resp []*Event
	var count int
	eq := builder.Eq{"deleted": false, "org_id": orgId}.And()
	if len(eventIds) != 0 {
		eq = eq.And(builder.In("id", eventIds))
	}
	if eventName != "" {
		eq = eq.And(builder.Like{"name", eventName})
	}
	if memberName != "" {
		eq = eq.And(builder.Like{"member_name", memberName})
	}
	if startTime != 0 {
		eq = eq.And(builder.Gte{"start_time": startTime})
	}
	if endTime != 0 {
		eq = eq.And(builder.Lte{"end_time": endTime})
	}
	if categoryId != 0 {
		eq = eq.And(builder.Eq{"category_id": categoryId})
	}
	query, args, _ := builder.Postgres().Select(eventFieldNames).From(m.Table).Where(eq).Limit(limit, (page-1)*limit).OrderBy(" created_at DESC ").ToSQL()
	queryC, argsC, _ := builder.Postgres().Select("COUNT( * )").From(m.Table).Where(eq).ToSQL()
	err := m.Conn.QueryRows(&resp, query, args...)
	err = m.Conn.QueryRow(&count, queryC, argsC...)
	if err != nil {
		return nil, 0, err
	}
	return resp, count, nil
}
