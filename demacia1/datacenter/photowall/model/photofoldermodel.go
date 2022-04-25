package model

import (
	"demacia/common/cachemodel"
	"fmt"
	"gitlab.u-jy.cn/xiaoyang/go-zero/core/logx"
	"gitlab.u-jy.cn/xiaoyang/go-zero/core/stores/redis"
	"gitlab.u-jy.cn/xiaoyang/go-zero/core/stores/sqlx"
	"strings"
	"xorm.io/builder"
)

var (
	photoFolderFieldNames             = strings.Join(cachemodel.RawFieldNames(&PhotoFolder{}, true), ",")
	cachePhotowallPhotoFolderIdPrefix = "cache:photowall:photoFolder:id:"
)

type (
	PhotoFolderModel struct {
		*cachemodel.CachedModel
	}

	PhotoFolder struct {
		Id    int64  `db:"id"`
		OrgId int64  `db:"org_id"` // 机构id
		Title string `db:"title"`  // 相册夹标题
	}
	PhotoFolders []PhotoFolder
)

func NewPhotoFolderModel(conn sqlx.SqlConn, cache *redis.Redis) *PhotoFolderModel {
	return &PhotoFolderModel{
		cachemodel.NewCachedModel(conn, `"photowall"."photo_folder"`, cache),
	}
}

func (m *PhotoFolderModel) Insert(data *PhotoFolder) (int64, error) {
	toSQL, i, err := builder.Postgres().Insert(builder.Eq{
		"org_id": data.OrgId,
		"title":  data.Title,
	}).Into(m.Table).ToSQL()
	if err != nil {
		return 0, err
	}
	var LastInsertId int64
	err = m.Conn.QueryRow(&LastInsertId, toSQL+" returning id", i...)
	if err != nil {
		return 0, err
	}
	return LastInsertId, nil
}

func (m *PhotoFolderModel) List(title string, orgId int64, page, limit int) (PhotoFolders, int, error) {
	var res PhotoFolders
	var total int
	var offset int
	if limit > 0 && limit < 100 {
		offset = (page - 1) * limit
	}
	eq := builder.And(builder.Eq{"deleted": false, "org_id": orgId})
	if title != "" {
		eq = eq.And(builder.Like{"title", title})
	}
	query, args, _ := builder.Postgres().
		Select(photoFolderFieldNames).
		From(m.Table).
		Where(eq).
		Limit(limit, offset).
		OrderBy("id DESC ").
		ToSQL()
	err := m.QueryRowsNoCache(&res, query, args...)
	if err != nil {
		return nil, 0, err
	}
	totalQuery, totalArgs, _ := builder.Postgres().
		Select("COUNT(id)").
		From(m.Table).
		Where(eq).
		ToSQL()
	err = m.QueryRowNoCache(&total, totalQuery, totalArgs...)
	if err != nil {
		return nil, 0, err
	}
	return res, total, nil
}

func (m *PhotoFolderModel) FindOneByName(title string, orgId int64) (*PhotoFolder, error) {
	toSQL, i, err := builder.Postgres().Select(photoFolderFieldNames).From(m.Table).Where(builder.Eq{"title": title, "org_id": orgId, "deleted": false}).Limit(1).ToSQL()
	if err != nil {
		return nil, err
	}
	var resp PhotoFolder
	err = m.QueryRowNoCache(&resp, toSQL, i...)
	if err != nil {
		if err != cachemodel.ErrNotFound {
			logx.Errorf("PhotoFolder findOneByName err :%s", err.Error())
		}
		return nil, err
	}
	return &resp, nil
}

func (m *PhotoFolderModel) FindOneById(id int64) (*PhotoFolder, error) {
	toSQL, i, err := builder.Postgres().Select(photoFolderFieldNames).From(m.Table).Where(builder.Eq{"id": id, "deleted": false}).Limit(1).ToSQL()
	if err != nil {
		return nil, err
	}
	var resp PhotoFolder
	err = m.QueryRowNoCache(&resp, toSQL, i...)
	if err != nil {
		if err != cachemodel.ErrNotFound {
			logx.Errorf("PhotoFolder FindOneById err :%s", err.Error())
		}
		return nil, err
	}
	return &resp, nil
}

func (m *PhotoFolderModel) Rename(data *PhotoFolder) error {
	sqlString, args, _ := builder.Postgres().Update(builder.Eq{
		"title": data.Title,
	}).Where(builder.Eq{
		"id":      data.Id,
		"org_id":  data.OrgId,
		"deleted": false,
	}).From(m.Table).ToSQL()
	_, err := m.ExecNoCache(sqlString, args...)
	return err
}

func (m *PhotoFolderModel) Delete(id, orgId int64) error {

	sqlString, args, _ := builder.Postgres().Update(builder.Eq{
		"deleted": true,
	}).Where(builder.Eq{"id": id, "org_id": orgId}).From(m.Table).ToSQL()

	_, err := m.ExecNoCache(sqlString, args...)
	return err
}

func (m *PhotoFolderModel) formatPrimary(primary interface{}) string {
	return fmt.Sprintf("%s%v", cachePhotowallPhotoFolderIdPrefix, primary)
}

func (m *PhotoFolderModel) queryPrimary(conn sqlx.SqlConn, v, primary interface{}) error {
	query := fmt.Sprintf("select %s from %s where id = $1 limit 1", photoFolderFieldNames, m.Table)
	return conn.QueryRow(v, query, primary)
}
