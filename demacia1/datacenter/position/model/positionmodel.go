package model

import (
	"demacia/common/cachemodel"
	"fmt"
	"gitlab.u-jy.cn/xiaoyang/go-zero/core/stores/redis"
	"gitlab.u-jy.cn/xiaoyang/go-zero/core/stores/sqlx"
	"strings"
	"time"
	"xorm.io/builder"
)

var (
	positionFieldNames    = strings.Join(cachemodel.RawFieldNames(&Position{}, true), ",")
	cachePositionIdPrefix = "cache:organization:organization:id:"
)

type (
	PositionModel struct {
		*cachemodel.CachedModel
	}

	Position struct {
		Id           int64  `db:"id"`
		Oid          int64  `db:"org_id"`
		PositionName string `db:"position_name"`
		ClassId      int64  `db:"class_id"`
		ClassName    string `db:"class_name"`
	}

	Positions []*Position
)

func NewPositionModel(conn sqlx.SqlConn, cache *redis.Redis) *PositionModel {
	return &PositionModel{
		CachedModel: cachemodel.NewCachedModel(conn, `"position"."position"`, cache),
	}
}

func (m *PositionModel) SelectList(oid int64, page, limit int, positionName string) (Positions, int, error) {

	var count int

	var positions Positions

	and := builder.NewCond()
	if positionName != "" {
		and = builder.And(builder.Like{"position_name", positionName})
	}

	sql, args, _ := builder.Postgres().Select(positionFieldNames).From(m.Table).Where(builder.Eq{"org_id": oid}).And(and).Limit(limit, (page-1)*limit).ToSQL()
	sqlC, args, _ := builder.Postgres().Select("count(id)").From(m.Table).Where(builder.Eq{"org_id": oid}).And(and).ToSQL()

	err := m.Conn.QueryRows(&positions, sql, args...)
	if err != nil {
		return positions, count, err
	}

	err = m.Conn.QueryRow(&count, sqlC, args...)
	return positions, count, err
}

func (m *PositionModel) InsertPosition(p *Position) (int64, error) {

	sql, args, _ := builder.Postgres().Insert(builder.Eq{"org_id": p.Oid, "position_name": p.PositionName, "class_id": p.ClassId, "class_name": p.ClassName}).Into(m.Table).ToSQL()

	var res int64

	err := m.Conn.QueryRow(&res, sql+" returning id", args...)

	return res, err
}

func (m *PositionModel) InsertPositions(positions []*Position) error {

	sql := fmt.Sprintf(`INSERT INTO %s("position_name",  "class_name", "class_id", "org_id")`, m.Table)

	for i, position := range positions {
		if i == 0 {
			sql += fmt.Sprintf(` VALUES ('%s', '%s', %d, %d)`, position.PositionName, position.ClassName, position.ClassId, position.Oid)
		} else {
			sql += fmt.Sprintf(`,('%s', '%s', %d, %d)`, position.PositionName, position.ClassName, position.ClassId, position.Oid)
		}
	}

	_, err := m.Conn.Exec(sql)

	return err
}

func (m *PositionModel) SelectByIsBind(oid int64, cIds []int64) (Positions, error) {

	var positions Positions

	sql, args, _ := builder.Postgres().Select(positionFieldNames).From(m.Table).Where(builder.Eq{"org_id": oid}).And(builder.In("class_id", cIds)).ToSQL()
	err := m.Conn.QueryRows(&positions, sql, args...)

	return positions, err
}

func (m *PositionModel) SelectByNames(oid int64, names []string) (Positions, error) {
	var positions Positions

	sql, args, _ := builder.Postgres().Select(positionFieldNames).From(m.Table).Where(builder.Eq{"org_id": oid, "class_id": 0}).And(builder.In("position_name", names)).ToSQL()
	err := m.Conn.QueryRows(&positions, sql, args...)

	return positions, err
}

func (m *PositionModel) SelectByClassId(cid int64) (Position, error) {

	var position Position

	sql, args, _ := builder.Postgres().Select(positionFieldNames).From(m.Table).Where(builder.Eq{"class_id": cid}).Limit(1).ToSQL()

	err := m.Conn.QueryRow(&position, sql, args...)

	return position, err
}

func (m *PositionModel) UpdatePosition(position *Position) error {

	sql, args, _ := builder.Postgres().Update(builder.Eq{"position_name": position.PositionName, "class_id": position.ClassId, "class_name": position.ClassName, "update_time": time.Now().Format("2006-01-02 15:04:05")}).Where(builder.Eq{"id": position.Id}).From(m.Table).ToSQL()

	_, err := m.Conn.Exec(sql, args...)

	return err
}

func (m *PositionModel) SelectByOidAndId(oid, pid int64) (Position, error) {

	var position Position

	sql, args, _ := builder.Postgres().Select(positionFieldNames).From(m.Table).Where(builder.Eq{"org_id": oid, "id": pid}).Limit(1).ToSQL()

	err := m.Conn.QueryRow(&position, sql, args...)

	return position, err
}

func (m *PositionModel) SelectById(id int64) (Position, error) {
	var position Position

	sql, args, _ := builder.Postgres().Select(positionFieldNames).From(m.Table).Where(builder.Eq{"id": id}).Limit(1).ToSQL()

	err := m.Conn.QueryRow(&position, sql, args...)

	return position, err
}

func (m *PositionModel) DeleteById(oid, pid int64) error {

	toSQL, args, _ := builder.Postgres().Delete(builder.Eq{"org_id": oid, "id": pid}).From(m.Table).ToSQL()

	_, err := m.Conn.Exec(toSQL, args...)

	return err
}

func (m *PositionModel) DeleteByIds(oid, pid []int64) error {

	toSQL, args, _ := builder.Postgres().Delete(builder.Eq{"org_id": oid}).And(builder.In("id", pid)).From(m.Table).ToSQL()

	_, err := m.Conn.Exec(toSQL, args...)

	return err
}

func (m *PositionModel) DeleteByClassId(id int64) error {

	toSQL, args, _ := builder.Postgres().Delete(builder.Eq{"class_id": id}).From(m.Table).ToSQL()

	_, err := m.Conn.Exec(toSQL, args...)

	return err
}

func (m *PositionModel) SelectByClassFullName(name string, oid int64) (Position, error) {

	var position Position

	sql, args, _ := builder.Postgres().Select(positionFieldNames).From(m.Table).Where(builder.Eq{"position_name": name, "org_id": oid}).Limit(1).ToSQL()

	err := m.Conn.QueryRow(&position, sql, args...)

	return position, err

}

func (m *PositionModel) SelectExistByPositionName(oid int64, pName string) (int, error) {

	var count int

	sql, args, _ := builder.Postgres().Select("count(id)").From(m.Table).Where(builder.Eq{"org_id": oid, "position_name": pName}).ToSQL()

	err := m.Conn.QueryRow(&count, sql, args...)

	return count, err
}

func (m *PositionModel) UnbindClass(cid int64) error {

	sql, args, _ := builder.Postgres().Update(builder.Eq{"class_id": 0, "class_name": ""}).From(m.Table).Where(builder.Eq{"class_id": cid}).ToSQL()

	_, err := m.Conn.Exec(sql, args...)

	return err
}
