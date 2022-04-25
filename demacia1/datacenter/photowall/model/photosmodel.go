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
	photosFieldNames = strings.Join(cachemodel.RawFieldNames(&Photos{}, true), ",")

	cachePhotowallPhotosIdPrefix = "cache:photowall:photos:id:"
)

type (
	PhotosModel struct {
		*cachemodel.CachedModel
	}

	Photos struct {
		Id            int64  `db:"id"`
		Title         string `db:"title"`           // 资源名称
		Url           string `db:"url"`             // 资源地址
		OrgId         int64  `db:"org_id"`          // 机构id
		PhotoFolderId int64  `db:"photo_folder_id"` // 相册id
		CreatedTime   string `db:"created_time"`    // 上传时间
	}
	Photose []Photos

	ListReq struct {
		PhotoFolderId int64
		Title         string // 资源名称
		OrgId         int64  // 机构id
		Page          int
		Limit         int
	}
)

func NewPhotosModel(conn sqlx.SqlConn, cache *redis.Redis) *PhotosModel {
	return &PhotosModel{
		cachemodel.NewCachedModel(conn, `"photowall"."photos"`, cache),
	}
}

func (m *PhotosModel) Insert(data *Photos) (int64, error) {
	toSQL, i, err := builder.Postgres().Insert(builder.Eq{
		"org_id":          data.OrgId,
		"title":           data.Title,
		"url":             data.Url,
		"photo_folder_id": data.PhotoFolderId,
		"created_time":    data.CreatedTime,
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

func (m *PhotosModel) List(data *ListReq) (Photose, int, error) {
	var res Photose
	var total int
	var offset int
	if data.Limit > 0 && data.Limit < 100 {
		offset = (data.Page - 1) * data.Limit
	}
	eq := builder.And(builder.Eq{"deleted": false, "org_id": data.OrgId})
	if data.Title != "" {
		eq = eq.And(builder.Like{"title", data.Title})
	}
	if data.PhotoFolderId > 0 {
		eq = eq.And(builder.Eq{"photo_folder_id": data.PhotoFolderId})
	}
	query, args, _ := builder.Postgres().
		Select(photosFieldNames).
		From(m.Table).
		Where(eq).
		Limit(data.Limit, offset).
		OrderBy("created_time DESC ").
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

func (m *PhotosModel) Update(data *Photos) error {
	photowallPhotosIdKey := fmt.Sprintf("%s%v", cachePhotowallPhotosIdPrefix, data.Id)
	_, err := m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set %s where id = $1", m.Table, photosFieldNames)
		return conn.Exec(query, data.Id, data.Title, data.Url, data.OrgId, data.CreatedTime)
	}, photowallPhotosIdKey)
	return err
}

func (m *PhotosModel) Delete(id int64) error {

	photowallPhotosIdKey := fmt.Sprintf("%s%v", cachePhotowallPhotosIdPrefix, id)
	_, err := m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("delete from %s where id = $1", m.Table)
		return conn.Exec(query, id)
	}, photowallPhotosIdKey)
	return err
}

func (m *PhotosModel) formatPrimary(primary interface{}) string {
	return fmt.Sprintf("%s%v", cachePhotowallPhotosIdPrefix, primary)
}

func (m *PhotosModel) queryPrimary(conn sqlx.SqlConn, v, primary interface{}) error {
	query := fmt.Sprintf("select %s from %s where id = $1 limit 1", photosFieldNames, m.Table)
	return conn.QueryRow(v, query, primary)
}
