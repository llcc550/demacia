package model

import (
	"fmt"
	"strings"

	"demacia/common/basefunc"
	"demacia/common/cachemodel"

	"gitlab.u-jy.cn/xiaoyang/go-zero/core/stores/redis"
	"gitlab.u-jy.cn/xiaoyang/go-zero/core/stores/sqlx"
	"xorm.io/builder"
)

var (
	cardFieldNames  = strings.Join(cachemodel.RawFieldNames(&Card{}, true), ",")
	cacheCardPrefix = "cache:card:list:"
)

type (
	CardModel struct {
		*cachemodel.CachedModel
	}
	Card struct {
		Id         int64  `db:"id"`
		OrgId      int64  `db:"org_id"`
		ObjectId   int64  `db:"object_id"`
		CardNum    string `db:"card_num"`
		ObjectRole int8   `db:"object_role"` // 1：教师，2：学生

	}
	Cards []*Card
)

func NewCardModel(conn sqlx.SqlConn, cache *redis.Redis) *CardModel {
	return &CardModel{
		CachedModel: cachemodel.NewCachedModel(conn, `"card"."card"`, cache),
	}
}

func (m *CardModel) List(orgId, objectId int64, objectRole int8) (Cards, error) {
	key := m.formatPrimary(orgId, objectId, objectRole)
	var resp Cards
	err := m.QueryRow(&resp, key, func(conn sqlx.SqlConn, v interface{}) error {
		query, args, _ := builder.Postgres().Select(cardFieldNames).From(m.Table).Where(builder.Eq{"org_id": orgId, "object_id": objectId, "object_role": objectRole}).ToSQL()
		return conn.QueryRows(v, query, args...)
	})
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (m *CardModel) Check(orgId, objectId int64, objectRole int8, cardNum string) error {
	var res Card
	query, args, _ := builder.Postgres().Select(cardFieldNames).From(m.Table).Where(builder.Eq{"org_id": orgId, "card_num": cardNum, "object_role": objectRole}).And(builder.Neq{"object_id": objectId}).ToSQL()
	err := m.Conn.QueryRow(&res, query, args...)
	return err
}

func (m *CardModel) Delete(orgId, objectId int64, objectRole int8) error {
	query, args, _ := builder.Postgres().Delete(builder.Eq{"org_id": orgId, "object_id": objectId, "object_role": objectRole}).From(m.Table).ToSQL()
	_, err := m.ExecNoCache(query, args...)
	return err
}

func (m *CardModel) BatchInsert(data Cards) error {
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

func (m *CardModel) RemoveCache(orgId, objectId int64, objectRole int8) error {
	return m.CachedModel.DelCache(m.formatPrimary(orgId, objectId, objectRole))
}

func (m *CardModel) formatPrimary(orgId, objectId int64, objectRole int8) string {
	return fmt.Sprintf("%s%d-%d-%d", cacheCardPrefix, orgId, objectRole, objectId)
}
