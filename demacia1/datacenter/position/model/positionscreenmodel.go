package model

import (
	"demacia/common/cachemodel"
	"fmt"
	"github.com/go-xorm/builder"
	"gitlab.u-jy.cn/xiaoyang/go-zero/core/stores/redis"
	"gitlab.u-jy.cn/xiaoyang/go-zero/core/stores/sqlx"
	"strings"
)

var (
	positionScreenFieldNames    = strings.Join(cachemodel.RawFieldNames(&PositionScreen{}, true), ",")
	cachePositionScreenIdPrefix = "cache:organization:organization:id:"
)

type (
	PositionScreenModel struct {
		*cachemodel.CachedModel
	}

	PositionScreen struct {
		Id         int64  `db:"id"`
		PositionId int64  `db:"position_id"`
		ScreenId   int64  `db:"screen_id"`
		ScreenName string `db:"screen_name"`
	}
	PositionScreens []*PositionScreen
)

func NewPositionScreenModel(conn sqlx.SqlConn, cache *redis.Redis) *PositionScreenModel {
	return &PositionScreenModel{
		CachedModel: cachemodel.NewCachedModel(conn, `"position"."position_screen"`, cache),
	}
}

func (m *PositionScreenModel) SelectByPidList(ids []int64) (PositionScreens, error) {

	var positionScreens PositionScreens

	sql, args, _ := builder.Postgres().Select(positionScreenFieldNames).From(m.Table).Where(builder.In("position_id", ids)).ToSQL()

	err := m.Conn.QueryRows(&positionScreens, sql, args...)

	return positionScreens, err
}

func (m *PositionScreenModel) InsertPositionScreen(ps []*PositionScreen) error {

	sql := fmt.Sprintf(`INSERT INTO%s("position_id", "screen_id",  "screen_name")`, m.Table)

	for i, p := range ps {
		if i == 0 {
			sql += fmt.Sprintf("VALUES ( %d, %d, '%s')", p.PositionId, p.ScreenId, p.ScreenName)
		} else {
			sql += fmt.Sprintf(",( %d, %d, '%s')", p.PositionId, p.ScreenId, p.ScreenName)
		}
	}

	_, err := m.Conn.Exec(sql)

	return err
}

func (m *PositionScreenModel) UpdatePositionScreen(ps []*PositionScreen, pid int64) error {
	return m.Conn.Transact(func(session sqlx.Session) error {
		sqlD, args, _ := builder.Postgres().Delete().From(m.Table).Where(builder.Eq{"position_id": pid}).ToSQL()
		_, err := session.Exec(sqlD, args...)
		if err != nil {
			return err
		}
		sql := fmt.Sprintf(`INSERT INTO%s("position_id", "screen_id",  "screen_name")`, m.Table)

		for i, p := range ps {
			if i == 0 {
				sql += fmt.Sprintf("VALUES ( %d, %d, '%s')", p.PositionId, p.ScreenId, p.ScreenName)
			} else {
				sql += fmt.Sprintf(",( %d, %d, '%s')", p.PositionId, p.ScreenId, p.ScreenName)
			}
		}

		_, err = m.Conn.Exec(sql)
		return err
	})
}

func (m *PositionScreenModel) SelectByScreenIds(ids []int64) (PositionScreens, error) {

	var positionScreens PositionScreens

	sql, args, _ := builder.Postgres().Select(positionScreenFieldNames).From(m.Table).Where(builder.In("screen_id", ids)).ToSQL()

	err := m.Conn.QueryRows(&positionScreens, sql, args...)

	return positionScreens, err
}

func (m *PositionScreenModel) DeleteByPid(pid int64) error {

	sql, args, _ := builder.Postgres().Delete(builder.Eq{"position_id": pid}).From(m.Table).ToSQL()

	_, err := m.Conn.Exec(sql, args...)

	return err
}
