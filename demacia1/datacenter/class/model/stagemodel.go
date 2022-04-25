package model

import (
	"demacia/common/cachemodel"

	"strings"

	"gitlab.u-jy.cn/xiaoyang/go-zero/core/stores/redis"
	"gitlab.u-jy.cn/xiaoyang/go-zero/core/stores/sqlx"
	"xorm.io/builder"
)

var (
	stageFieldNames         = strings.Join(cachemodel.RawFieldNames(&Stage{}, true), ",")
	cacheClassStageIdPrefix = "cache:class:stage:id:"
)

type (
	StageModel struct {
		*cachemodel.CachedModel
	}

	Stage struct {
		Id    int64  `db:"id"`
		OrgId int64  `db:"org_id"` // 所属机构ID
		Title string `db:"title"`  // 学段名称
		Year  int64  `db:"year"`   // 学制
	}
	Stages []*Stage
)

func NewStageModel(conn sqlx.SqlConn, cache *redis.Redis) *StageModel {
	return &StageModel{
		CachedModel: cachemodel.NewCachedModel(conn, `"class"."stage"`, cache),
	}
}

func (m *StageModel) ListByOrgId(orgId int64) (Stages, error) {
	sql, args, _ := builder.Postgres().Select(stageFieldNames).From(m.Table).Where(builder.Eq{"org_id": orgId, "deleted": "false"}).ToSQL()
	var res Stages
	err := m.Conn.QueryRows(&res, sql, args...)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (m *StageModel) InsertStageOfOrgId(stage *Stage) (int64, error) {
	sqlString, args, _ := builder.Postgres().Insert(builder.Eq{
		"org_id": stage.OrgId,
		"title":  stage.Title,
		"year":   stage.Year,
	}).Into(m.Table).ToSQL()
	var LastInsertId int64
	err := m.Conn.QueryRow(&LastInsertId, sqlString+" returning id", args...)
	return LastInsertId, err
}

func (m *StageModel) GetStageById(Id int64) (*Stage, error) {
	sql, args, _ := builder.Postgres().Select(stageFieldNames).From(m.Table).Where(builder.Eq{"id": Id, "deleted": "false"}).ToSQL()
	var res Stage
	err := m.Conn.QueryRow(&res, sql, args...)
	if err != nil {
		return nil, err
	}
	return &res, nil
}
func (m *StageModel) GetStageByTitleAndOrg(stageTitle string, orgId int64) (*Stage, error) {
	sql, args, _ := builder.Postgres().Select(stageFieldNames).From(m.Table).Where(builder.Eq{"title": stageTitle, "org_id": orgId, "deleted": "false"}).ToSQL()
	var res Stage
	err := m.Conn.QueryRow(&res, sql, args...)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func (m *StageModel) UpdateStageById(orgId int64, stageId int64, stageTitle string) error {
	sqlString, args, _ := builder.Postgres().Update(builder.Eq{
		"title": stageTitle,
	}).From(m.Table).Where(builder.Eq{"id": stageId, "org_id": orgId, "deleted": "false"}).ToSQL()
	_, err := m.ExecNoCache(sqlString, args...)
	if err != nil {
		return err
	}
	return nil
}
func (m *StageModel) UpdateStageYear(orgId int64, stageId int64, Year int64) error {
	sqlString, args, _ := builder.Postgres().Update(builder.Eq{
		"year": Year,
	}).From(m.Table).Where(builder.Eq{"id": stageId, "org_id": orgId, "deleted": false}).ToSQL()
	_, err := m.ExecNoCache(sqlString, args...)
	if err != nil {
		return err
	}
	return nil
}

func (m *StageModel) DeleteStageById(orgId int64, stageId int64) error {
	sqlString, args, _ := builder.Postgres().Update(builder.Eq{
		"deleted": "true",
	}).From(m.Table).Where(builder.Eq{"id": stageId, "org_id": orgId}).ToSQL()
	_, err := m.ExecNoCache(sqlString, args...)
	if err != nil {
		return err
	}
	return nil
}
