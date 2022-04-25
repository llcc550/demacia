package model

import (
	"strings"

	"demacia/common/basefunc"
	"demacia/common/cachemodel"

	"gitlab.u-jy.cn/xiaoyang/go-zero/core/stores/redis"
	"gitlab.u-jy.cn/xiaoyang/go-zero/core/stores/sqlx"
	"xorm.io/builder"
)

var (
	ticketFieldNames = strings.Join(cachemodel.RawFieldNames(&Ticket{}, true), ",")
)

type (
	TicketModel struct {
		*cachemodel.CachedModel
	}
	Ticket struct {
		Id         int64  `db:"id"`
		UserId     int64  `db:"user_id"`
		OrgId      int64  `db:"org_id"`
		TicketDate string `db:"ticket_date"`
	}
	Tickets []*Ticket
)

func NewTicketModel(conn sqlx.SqlConn, cache *redis.Redis) *TicketModel {
	return &TicketModel{
		CachedModel: cachemodel.NewCachedModel(conn, `"ticket"."ticket"`, cache),
	}
}

func (m *TicketModel) FindListByUserIdAndDate(userId int64, date string) (Tickets, error) {
	query, args, _ := builder.Postgres().Select(ticketFieldNames).From(m.Table).Where(builder.Eq{"user_id": userId, "ticket_date": date}).ToSQL()
	var res Tickets
	err := m.Conn.QueryRows(&res, query, args...)
	if err != nil {
		return nil, cachemodel.ErrNotFound
	}
	return res, nil
}

func (m *TicketModel) BatchInsert(data Tickets) error {
	if len(data) == 0 {
		return nil
	}
	query, args, err := basefunc.BatchInsertString(m.Table, data)
	if err != nil {
		return err
	}
	_, err = m.ExecNoCache(query, args...)
	return err
}
